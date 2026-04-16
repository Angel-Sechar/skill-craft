package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	switch m.Screen {
	case ScreenWelcome:
		return m.viewWelcome()
	case ScreenFramework:
		return m.viewRadio("Select your framework", "")
	case ScreenStack:
		return m.viewRadio("Select your stack", "")
	case ScreenArchitecture:
		return m.viewArchitecture()
	case ScreenDrivenDesign:
		return m.viewRadio("Select driven design", "")
	case ScreenPractices:
		return m.viewPractices()
	case ScreenInstalling:
		return m.viewInstalling()
	case ScreenDone:
		return m.viewDone()
	}
	return ""
}

// ── Welcome ───────────────────────────────────────────────────────────────────

func (m Model) viewWelcome() string {
	var b strings.Builder

	b.WriteString("\n")
	for _, line := range strings.Split(SwordArt, "\n") {
		b.WriteString(SwordStyle.Render(line) + "\n")
	}
	b.WriteString("\n")
	b.WriteString("  " + TitleStyle.Render("skill-craft") + "  " + VersionStyle.Render("v"+Version) + "\n")
	b.WriteString("\n")
	b.WriteString(TaglineStyle.Render("  Craft your development skills.") + "\n")
	b.WriteString("\n\n")
	b.WriteString(HintStyle.Render("  enter: start   q: quit") + "\n")

	return BoxStyle.Render(b.String())
}

// ── Radio screen (single selection) ──────────────────────────────────────────

func (m Model) viewRadio(title string, subtitle string) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  " + TitleStyle.Render(title) + "\n")
	if subtitle != "" {
		b.WriteString("  " + AuthorStyle.Render(subtitle) + "\n")
	}
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n\n")

	for i, opt := range m.Options {
		if i == m.Cursor {
			b.WriteString(SelectedItemStyle.Render("  ▸ "+opt.Label) + "\n")
		} else {
			b.WriteString(NormalItemStyle.Render("    "+opt.Label) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(HintStyle.Render("  ↑↓: navigate   enter: select   q: quit") + "\n")

	return BoxStyle.Render(b.String())
}

// ── Architecture screen (checkboxes with conflict) ────────────────────────────

func (m Model) viewArchitecture() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  " + TitleStyle.Render("Select architecture") + "\n")
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n\n")

	for i, opt := range m.Options {
		checked := contains(m.Selection.Architecture, opt.ID)
		blocked := m.isArchitectureBlocked(opt.ID)

		var checkbox string
		if checked {
			checkbox = CheckedStyle.Render("  [✓] ")
		} else if blocked {
			checkbox = UncheckedStyle.Render("  [✗] ")
		} else {
			checkbox = UncheckedStyle.Render("  [ ] ")
		}

		var label string
		if blocked {
			label = AuthorStyle.Render(opt.Label + "  (conflicts with selection)")
		} else if i == m.Cursor {
			label = SelectedItemStyle.Render("▸ " + opt.Label)
		} else {
			label = NormalItemStyle.Render("  " + opt.Label)
		}

		b.WriteString(checkbox + label + "\n")
	}

	count := len(m.Selection.Architecture)
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n")
	b.WriteString(AuthorStyle.Render(fmt.Sprintf("  %d selected", count)) + "\n\n")
	b.WriteString(HintStyle.Render("  ↑↓: navigate   space: toggle   enter: confirm   q: quit") + "\n")

	return BoxStyle.Render(b.String())
}

// ── Practices screen (checkboxes, no limit) ───────────────────────────────────

func (m Model) viewPractices() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  " + TitleStyle.Render("Select good practices") + "\n")
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n\n")

	for i, opt := range m.Options {
		checked := contains(m.Selection.Practices, opt.ID)

		var checkbox string
		if checked {
			checkbox = CheckedStyle.Render("  [✓] ")
		} else {
			checkbox = UncheckedStyle.Render("  [ ] ")
		}

		var label string
		if i == m.Cursor {
			label = SelectedItemStyle.Render("▸ " + opt.Label)
		} else {
			label = NormalItemStyle.Render("  " + opt.Label)
		}

		b.WriteString(checkbox + label + "\n")
	}

	count := len(m.Selection.Practices)
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n")
	b.WriteString(AuthorStyle.Render(fmt.Sprintf("  %d selected", count)) + "\n\n")
	b.WriteString(HintStyle.Render("  ↑↓: navigate   space: toggle   enter: install   q: quit") + "\n")

	return BoxStyle.Render(b.String())
}

// ── Installing screen ─────────────────────────────────────────────────────────

func (m Model) viewInstalling() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  " + TitleStyle.Render("Installing skills...") + "\n")
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n\n")

	agents := []struct {
		key  string
		name string
		path string
	}{
		{"claude", "Claude Code", "~/.claude/skills/"},
		{"opencode", "OpenCode", "~/.config/opencode/skills/"},
		{"codex", "Codex", "~/.codex/skills/"},
	}

	for _, agent := range agents {
		result, done := m.Results[agent.key]
		var line string
		if !done {
			line = AuthorStyle.Render(fmt.Sprintf("  ◌  %-12s  %s", agent.name, agent.path))
		} else if result {
			line = SuccessStyle.Render(fmt.Sprintf("  ✓  %-12s  %s", agent.name, agent.path))
		} else {
			line = AuthorStyle.Render(fmt.Sprintf("  –  %-12s  not found", agent.name))
		}
		b.WriteString(line + "\n")
	}

	b.WriteString("\n")
	b.WriteString(HintStyle.Render("  please wait...") + "\n")

	return BoxStyle.Render(b.String())
}

// ── Done screen ───────────────────────────────────────────────────────────────

func (m Model) viewDone() string {
	var b strings.Builder

	b.WriteString("\n")
	for _, line := range strings.Split(SwordArt, "\n") {
		b.WriteString(SwordStyle.Render(line) + "\n")
	}
	b.WriteString("\n")
	b.WriteString("  " + TitleStyle.Render("Done.") + "\n")
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n\n")

	// Summary
	b.WriteString(AuthorStyle.Render("  Stack installed:") + "\n\n")

	if m.Selection.Stack != "" {
		b.WriteString(NormalItemStyle.Render("  framework    "+m.Selection.Stack) + "\n")
	}
	for _, arch := range m.Selection.Architecture {
		b.WriteString(NormalItemStyle.Render("  architecture  "+arch) + "\n")
	}
	if m.Selection.DrivenDesign != "" {
		b.WriteString(NormalItemStyle.Render("  design        "+m.Selection.DrivenDesign) + "\n")
	}
	for _, p := range m.Selection.Practices {
		b.WriteString(NormalItemStyle.Render("  practice      "+p) + "\n")
	}

	// Agent results
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render("  ────────────────────────────────────") + "\n\n")

	installed := 0
	for _, ok := range m.Results {
		if ok {
			installed++
		}
	}

	b.WriteString(SuccessStyle.Render(fmt.Sprintf(
		"  Installed into %d agent(s).", installed)) + "\n\n")
	b.WriteString(TaglineStyle.Render("  Your agent is forged.") + "\n\n")
	b.WriteString(HintStyle.Render("  press q to exit") + "\n")

	return BoxStyle.Render(b.String())
}

// ── Helper — is this architecture option blocked by current selection ─────────

func (m Model) isArchitectureBlocked(id string) bool {
	conflict := ArchitectureConflicts[id]
	if conflict == "" {
		return false
	}
	return contains(m.Selection.Architecture, conflict)
}
