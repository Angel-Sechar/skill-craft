package tui

import (
	"os"
	"path/filepath"

	"github.com/Angel-Sechar/skill-craft/assets"
)

// Agent target directories
var agentDirs = map[string]string{
	"claude":   ".claude/skills",
	"opencode": ".config/opencode/skills",
	"codex":    ".codex/skills",
}

// installSkills writes selected skills into all detected agent directories
func installSkills(sel Selection, overwrite bool) map[string]bool {
	results := make(map[string]bool)

	home, err := os.UserHomeDir()
	if err != nil {
		return results
	}

	skillIDs := resolveSkillIDs(sel)

	for key, relPath := range agentDirs {
		dir := filepath.Join(home, relPath)

		if !dirExists(dir) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				results[key] = false
				continue
			}
		}

		success := true
		for _, id := range skillIDs {
			if err := writeSkill(dir, id, overwrite); err != nil {
				success = false
			}
		}
		results[key] = success
	}

	return results
}

// resolveSkillIDs builds the full list of skill IDs to install from the selection
func resolveSkillIDs(sel Selection) []string {
	ids := []string{}

	// Stack skill
	if sel.Stack != "" {
		ids = append(ids, sel.Stack)
	}

	// Architecture skills
	ids = append(ids, sel.Architecture...)

	// Driven design skill
	if sel.DrivenDesign != "" {
		ids = append(ids, sel.DrivenDesign)
	}

	// Practices skills
	ids = append(ids, sel.Practices...)

	return ids
}

// writeSkill writes a single SKILL.md file into the agent's skills directory
func writeSkill(agentDir string, skillID string, overwrite bool) error {
	content := readEmbeddedSkill(skillID)
	if content == "" {
		return nil
	}

	skillDir := filepath.Join(agentDir, skillID)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}

	dest := filepath.Join(skillDir, "SKILL.md")
	if !overwrite && fileExists(dest) {
		return nil
	}

	return os.WriteFile(dest, []byte(content), 0644)
}

// readEmbeddedSkill reads a SKILL.md from the embedded filesystem
func readEmbeddedSkill(id string) string {
	// try each category
	categories := []string{"framework", "architecture", "driven", "practices", "database"}
	for _, cat := range categories {
		path := "skills/" + cat + "/" + id + "/SKILL.md"
		content, err := assets.SkillsFS.ReadFile(path)
		if err == nil {
			return string(content)
		}
	}
	return ""
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// checkExistingSkills returns a list of skill IDs that already exist in any agent directory
func checkExistingSkills(sel Selection) []string {
	home, _ := os.UserHomeDir()
	existing := []string{}
	ids := resolveSkillIDs(sel)

	// check all three agents, not just claude
	for _, relPath := range agentDirs {
		dir := filepath.Join(home, relPath)
		if !dirExists(dir) {
			continue
		}
		for _, id := range ids {
			dest := filepath.Join(dir, id, "SKILL.md")
			if fileExists(dest) && !contains(existing, id) {
				existing = append(existing, id)
			}
		}
	}
	return existing
}
