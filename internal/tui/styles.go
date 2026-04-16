package tui

import "github.com/charmbracelet/lipgloss"

// Golden theme colors
const (
	ColorBackground = "#0d0d0d"
	ColorGold       = "#f5c542"
	ColorBrightGold = "#ffd700"
	ColorAmber      = "#e8a020"
	ColorMuted      = "#555550"
	ColorText       = "#f0ead6"
	ColorSuccess    = "#7ec98f"
	ColorError      = "#e05252"
	ColorBorder     = "#3a3420"
)

var (
	// Box border — the outer frame
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color(ColorAmber)).
			Padding(0, 1)

	// Logo sword — bright gold
	SwordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorBrightGold))

	// App name — skill-craft
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGold)).
			Bold(true)

	// Tagline — forge your agent.
	TaglineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))

	// Divider line ───────────
	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorAmber))

	// Author — by Angel Sechar
	AuthorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// Key hints at the bottom
	HintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// Menu item — selected
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorBrightGold)).
				Bold(true)

	// Menu item — unselected
	NormalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))

	// Category header in skill list
	CategoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorAmber)).
			Bold(true)

	// Skill checked
	CheckedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess))

	// Skill unchecked
	UncheckedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// Version string
	VersionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// Success message
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess)).
			Bold(true)

	// Error message
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError))
)
