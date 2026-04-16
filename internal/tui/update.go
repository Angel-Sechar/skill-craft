package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case installDoneMsg:
		m.Results = msg.Results
		m.Screen = ScreenDone
		return m, nil

	case tea.KeyMsg:
		switch m.Screen {
		case ScreenWelcome:
			return m.updateWelcome(msg)
		case ScreenFramework:
			return m.updateRadio(msg, m.onFrameworkConfirm)
		case ScreenStack:
			return m.updateRadio(msg, m.onStackConfirm)
		case ScreenArchitecture:
			return m.updateArchitecture(msg)
		case ScreenDrivenDesign:
			return m.updateRadio(msg, m.onDrivenDesignConfirm)
		case ScreenPractices:
			return m.updatePractices(msg)
		case ScreenConflicts:
			return m.updateConflicts(msg)
		case ScreenInstalling:
			return m, nil
		case ScreenDone:
			return m.updateDone(msg)
		}
	}
	return m, nil
}

// ── Welcome ───────────────────────────────────────────────────────────────────

func (m Model) updateWelcome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.Screen = ScreenFramework
		m.Options = FrameworkOptions
		m.Cursor = 0
	}
	return m, nil
}

// ── Generic radio (single selection) ─────────────────────────────────────────

func (m Model) updateRadio(msg tea.KeyMsg, onConfirm func(Model) (Model, tea.Cmd)) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(m.Options)-1 {
			m.Cursor++
		}
	case "enter":
		return onConfirm(m)
	}
	return m, nil
}

// ── Framework confirmed ───────────────────────────────────────────────────────

func (m Model) onFrameworkConfirm(m2 Model) (Model, tea.Cmd) {
	selected := m2.Options[m2.Cursor]
	m2.Selection.Framework = Framework(selected.ID)
	m2.Screen = ScreenStack
	m2.Options = StackOptions[Framework(selected.ID)]
	m2.Cursor = 0
	return m2, nil
}

// ── Stack confirmed ───────────────────────────────────────────────────────────

func (m Model) onStackConfirm(m2 Model) (Model, tea.Cmd) {
	selected := m2.Options[m2.Cursor]
	m2.Selection.Stack = selected.ID

	// SQL skips architecture, driven design, and practices — go straight to install
	if m2.Selection.Framework == FrameworkSQL {
		m2.Screen = ScreenInstalling
		m2.Results = make(map[string]bool)
		return m2, doInstall(m2.Selection, false)
	}

	m2.Screen = ScreenArchitecture
	m2.Options = ArchitectureOptions(m2.Selection.Framework)
	m2.Cursor = 0
	return m2, nil
}

// ── Architecture (checkboxes with conflict) ───────────────────────────────────

func (m Model) updateArchitecture(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Options)-1 {
			m.Cursor++
		}

	case " ":
		id := m.Options[m.Cursor].ID

		if contains(m.Selection.Architecture, id) {
			// deselect
			m.Selection.Architecture = remove(m.Selection.Architecture, id)
		} else {
			// check conflict before selecting
			conflict := ArchitectureConflicts[id]
			if conflict != "" && contains(m.Selection.Architecture, conflict) {
				// blocked — do nothing
				break
			}
			m.Selection.Architecture = append(m.Selection.Architecture, id)
		}

	case "enter":
		if len(m.Selection.Architecture) > 0 {
			m.Screen = ScreenDrivenDesign
			m.Options = DrivenOptions
			m.Cursor = 0
		}
	}
	return m, nil
}

// ── Driven design confirmed ───────────────────────────────────────────────────

func (m Model) onDrivenDesignConfirm(m2 Model) (Model, tea.Cmd) {
	selected := m2.Options[m2.Cursor]
	m2.Selection.DrivenDesign = selected.ID
	m2.Screen = ScreenPractices
	m2.Options = PracticesOptions
	m2.Cursor = 0
	return m2, nil
}

// ── Practices (checkboxes, no limit) ─────────────────────────────────────────

func (m Model) updatePractices(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Options)-1 {
			m.Cursor++
		}

	case " ":
		id := m.Options[m.Cursor].ID
		if contains(m.Selection.Practices, id) {
			m.Selection.Practices = remove(m.Selection.Practices, id)
		} else {
			m.Selection.Practices = append(m.Selection.Practices, id)
		}

	case "enter": // in updatePractices
		existing := checkExistingSkills(m.Selection)
		if len(existing) > 0 {
			m.ExistingSkills = existing
			m.Screen = ScreenConflicts
			m.Options = ConflictOptions
			m.Cursor = 0
		} else {
			m.Screen = ScreenInstalling
			m.Results = make(map[string]bool)
			return m, doInstall(m.Selection, false)
		}
	}
	return m, nil
}

// ── Done ──────────────────────────────────────────────────────────────────────

func (m Model) updateDone(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "enter":
		return m, tea.Quit
	}
	return m, nil
}

// ── Conflicts (radio) ─────────────────────────────────────────────────────────
func (m Model) updateConflicts(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(m.Options)-1 {
			m.Cursor++
		}
	case "enter":
		selected := m.Options[m.Cursor]
		switch selected.ID {
		case "overwrite":
			m.Screen = ScreenInstalling
			m.Results = make(map[string]bool)
			return m, doInstall(m.Selection, true)
		case "skip":
			m.Screen = ScreenInstalling
			m.Results = make(map[string]bool)
			return m, doInstall(m.Selection, false)
		case "cancel":
			m.Screen = ScreenPractices
			m.Options = PracticesOptions
			m.Cursor = 0
		}
	}
	return m, nil
}

// ── Install command ───────────────────────────────────────────────────────────

type installDoneMsg struct {
	Results map[string]bool
}

func doInstall(sel Selection, overwrite bool) tea.Cmd {
	return func() tea.Msg {
		results := installSkills(sel, overwrite)
		return installDoneMsg{Results: results}
	}
}
