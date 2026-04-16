package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current screen
func (m Model) View() string {
	switch m.Screen {
	case ScreenWelcome:
		return m.viewWelcome()
	case ScreenSkills:
		return m.viewSkills()
	case ScreenInstalling:
		return m.viewInstalling()
	case ScreenDone:
		return m.viewDone()
	default:
		return m.viewWelcome()
	}
}

func (m Model) viewWelcome() string {
	var b strings.Builder

	// Sword in bright gold
	b.WriteString(SwordStyle.Render(SwordArt))
	b.WriteString("\n\n")

	// Divider
	b.WriteString(DividerStyle.Render("────────────────────────"))
	b.WriteString("\n")

	// Title
	b.WriteString(TitleStyle.Render("  skill-craft"))
	b.WriteString("  ")
	b.WriteString(VersionStyle.Render("v" + Version))
	b.WriteString("\n")

	// Divider
	b.WriteString(DividerStyle.Render("────────────────────────"))
	b.WriteString("\n\n")

	// Tagline
	b.WriteString(TaglineStyle.Render("  forge your agent."))
	b.WriteString("\n\n")

	// Author
	b.WriteString(AuthorStyle.Render("  by Angel Sechar"))
	b.WriteString("\n\n\n")

	// Key hints
	b.WriteString(HintStyle.Render("  enter: continue   q: quit"))
	b.WriteString("\n")

	return BoxStyle.Render(b.String())
}

func (m Model) viewSkills() string {
	var b strings.Builder

	// Header
	b.WriteString(TitleStyle.Render("  Select skills to install"))
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render("────────────────────────────────────────"))
	b.WriteString("\n\n")

	// Group skills by category
	categories := m.groupByCategory()

	flatIdx := 0
	for _, cat := range categories {
		// Category header
		b.WriteString(CategoryStyle.Render("  " + cat))
		b.WriteString("\n")

		for _, skill := range m.skillsForCategory(cat) {
			// Checkbox
			var checkbox string
			if skill.Checked {
				checkbox = CheckedStyle.Render("  [✓] ")
			} else {
				checkbox = UncheckedStyle.Render("  [ ] ")
			}

			// Label — highlight if cursor is here
			var label string
			if flatIdx == m.Cursor {
				label = SelectedItemStyle.Render("▸ " + skill.Label)
			} else {
				label = NormalItemStyle.Render("  " + skill.Label)
			}

			b.WriteString(checkbox + label + "\n")
			flatIdx++
		}
		b.WriteString("\n")
	}

	// Selected count
	checked := m.checkedCount()
	b.WriteString(DividerStyle.Render("────────────────────────────────────────"))
	b.WriteString("\n")
	b.WriteString(AuthorStyle.Render(fmt.Sprintf("  %d skill(s) selected", checked)))
	b.WriteString("\n\n")

	// Key hints
	b.WriteString(HintStyle.Render("  ↑↓: navigate   space: toggle   enter: install   q: quit"))
	b.WriteString("\n")

	return BoxStyle.Render(b.String())
}

func (m Model) viewInstalling() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("  Installing skills..."))
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render("────────────────────────────────────────"))
	b.WriteString("\n\n")

	agents := []struct {
		name string
		path string
		key  string
	}{
		{"Claude Code", "~/.claude/skills/", "claude"},
		{"OpenCode", "~/.config/opencode/skills/", "opencode"},
		{"Codex", "~/.codex/skills/", "codex"},
	}

	for _, agent := range agents {
		result, exists := m.Results[agent.key]
		if !exists {
			b.WriteString(AuthorStyle.Render(fmt.Sprintf("  ◌  %-12s %s", agent.name, agent.path)))
		} else if result {
			b.WriteString(SuccessStyle.Render(fmt.Sprintf("  ✓  %-12s %s", agent.name, agent.path)))
		} else {
			b.WriteString(AuthorStyle.Render(fmt.Sprintf("  –  %-12s %s (not found)", agent.name, agent.path)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HintStyle.Render("  please wait..."))
	b.WriteString("\n")

	return BoxStyle.Render(b.String())
}

func (m Model) viewDone() string {
	var b strings.Builder

	b.WriteString(SwordStyle.Render("  †"))
	b.WriteString("\n\n")

	b.WriteString(TitleStyle.Render("  Done."))
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render("────────────────────────────────────────"))
	b.WriteString("\n\n")

	installed := 0
	for _, ok := range m.Results {
		if ok {
			installed++
		}
	}

	b.WriteString(SuccessStyle.Render(fmt.Sprintf(
		"  %d skill(s) installed into %d agent(s).",
		m.checkedCount(), installed,
	)))
	b.WriteString("\n\n")

	b.WriteString(TaglineStyle.Render("  Your agent is forged."))
	b.WriteString("\n\n")

	b.WriteString(HintStyle.Render("  press q to exit"))
	b.WriteString("\n")

	return BoxStyle.Render(b.String())
}

// ── helpers ──────────────────────────────────────────────────────────────────

func (m Model) groupByCategory() []string {
	seen := map[string]bool{}
	order := []string{}
	for _, s := range m.Skills {
		if !seen[s.Category] {
			seen[s.Category] = true
			order = append(order, s.Category)
		}
	}
	return order
}

func (m Model) skillsForCategory(cat string) []Skill {
	var out []Skill
	for _, s := range m.Skills {
		if s.Category == cat {
			out = append(out, s)
		}
	}
	return out
}

func (m Model) checkedCount() int {
	n := 0
	for _, s := range m.Skills {
		if s.Checked {
			n++
		}
	}
	return n
}

// Center a string within a given width
func center(s string, width int) string {
	pad := (width - lipgloss.Width(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s
}
