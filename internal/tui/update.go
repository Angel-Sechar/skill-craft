package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init starts the program
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all keyboard input and state transitions
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.Screen {
		case ScreenWelcome:
			return m.updateWelcome(msg)
		case ScreenSkills:
			return m.updateSkills(msg)
		case ScreenInstalling:
			return m, nil // no input during install
		case ScreenDone:
			return m.updateDone(msg)
		}
	}

	return m, nil
}

func (m Model) updateWelcome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.Screen = ScreenSkills
	}
	return m, nil
}

func (m Model) updateSkills(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	total := len(m.Skills)

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < total-1 {
			m.Cursor++
		}

	case " ": // spacebar toggles
		m.Skills[m.Cursor].Checked = !m.Skills[m.Cursor].Checked

	case "enter":
		// Only proceed if at least one skill is selected
		if m.checkedCount() > 0 {
			m.Screen = ScreenInstalling
			m.Results = make(map[string]bool)
			return m, doInstall(m.selectedSkills())
		}
	}

	return m, nil
}

func (m Model) updateDone(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "enter":
		return m, tea.Quit
	}
	return m, nil
}

// installDoneMsg is sent when installation finishes
type installDoneMsg struct {
	Results map[string]bool
}

// doInstall runs the installation in the background
func doInstall(skills []Skill) tea.Cmd {
	return func() tea.Msg {
		results := installSkills(skills)
		return installDoneMsg{Results: results}
	}
}

// Handle installDoneMsg, update model from any screen
func (m Model) handleInstallDone(msg installDoneMsg) (tea.Model, tea.Cmd) {
	m.Results = msg.Results
	m.Screen = ScreenDone
	return m, nil
}

// selectedSkills returns only checked skills
func (m Model) selectedSkills() []Skill {
	var out []Skill
	for _, s := range m.Skills {
		if s.Checked {
			out = append(out, s)
		}
	}
	return out
}
