package tui

import (
	"os"
	"path/filepath"
)

// Agent target directories
var agentDirs = map[string]string{
	"claude":   ".claude/skills",
	"opencode": ".config/opencode/skills",
	"codex":    ".codex/skills",
}

// installSkills writes selected skills into all detected agent directories
func installSkills(skills []Skill) map[string]bool {
	results := make(map[string]bool)
	home, err := os.UserHomeDir()
	if err != nil {
		return results
	}

	for key, relPath := range agentDirs {
		dir := filepath.Join(home, relPath)

		// Skip agents that aren't installed
		if !dirExists(dir) {
			results[key] = false
			continue
		}

		// Write each selected skill
		success := true
		for _, skill := range skills {
			if err := writeSkill(dir, skill); err != nil {
				success = false
			}
		}
		results[key] = success
	}

	return results
}

// writeSkill writes a single SKILL.md file into the agent's skills directory
func writeSkill(agentDir string, skill Skill) error {
	skillDir := filepath.Join(agentDir, skill.ID)

	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}

	dest := filepath.Join(skillDir, "SKILL.md")

	// Don't overwrite existing skills
	if fileExists(dest) {
		return nil
	}

	content := skillContent(skill)
	return os.WriteFile(dest, []byte(content), 0644)
}

// skillContent returns the SKILL.md content for a given skill
// In the real version this will use go:embed to read from assets/skills/
func skillContent(skill Skill) string {
	return "---\nname: " + skill.ID + "\ndescription: " + skill.Label + "\n---\n\n# " + skill.Label + "\n\nSkill content coming soon.\n"
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
