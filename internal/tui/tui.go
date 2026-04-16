package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Start launches the TUI
func Start() error {
	m := Model{
		Screen: ScreenWelcome,
		Cursor: 0,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
