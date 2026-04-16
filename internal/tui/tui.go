package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Start launches the TUI
func Start() error {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// initialModel builds the starting state with all curated skills
func initialModel() Model {
	return Model{
		Screen: ScreenWelcome,
		Cursor: 0,
		Skills: curatedSkills(),
	}
}

// curatedSkills returns the full catalog
func curatedSkills() []Skill {
	return []Skill{
		// Framework
		{ID: "dotnet-framework-45", Label: ".NET Framework 4.5", Category: "Framework"},
		{ID: "dotnet-core-8", Label: ".NET Core 8", Category: "Framework"},
		{ID: "aspnet-core", Label: "ASP.NET Core", Category: "Framework"},
		{ID: "spring-boot-2", Label: "Spring Boot 2 (Java 17)", Category: "Framework"},
		{ID: "spring-boot-3", Label: "Spring Boot 3 (Java 21)", Category: "Framework"},
		{ID: "angular", Label: "Angular", Category: "Framework"},

		// Architecture
		{ID: "clean-architecture", Label: "Clean Architecture", Category: "Architecture"},
		{ID: "hexagonal", Label: "Hexagonal Architecture", Category: "Architecture"},
		{ID: "onion", Label: "Onion Architecture", Category: "Architecture"},
		{ID: "monolith", Label: "Monolith", Category: "Architecture"},
		{ID: "scream", Label: "Scream Architecture", Category: "Architecture"},
		{ID: "microservices", Label: "Microservices", Category: "Architecture"},

		// Language
		{ID: "csharp", Label: "C#", Category: "Language"},
		{ID: "java", Label: "Java", Category: "Language"},
		{ID: "typescript", Label: "TypeScript", Category: "Language"},

		// Database
		{ID: "mssql", Label: "MS SQL Server", Category: "Database"},
		{ID: "mysql", Label: "MySQL", Category: "Database"},
		{ID: "postgresql", Label: "PostgreSQL", Category: "Database"},

		// Practices
		{ID: "event-driven", Label: "Event-Driven Design", Category: "Practices"},
		{ID: "ddd", Label: "Domain-Driven Design", Category: "Practices"},
		{ID: "tdd", Label: "Test-Driven Design", Category: "Practices"},
		{ID: "solid", Label: "SOLID", Category: "Practices"},
		{ID: "dependency-injection", Label: "Dependency Injection", Category: "Practices"},
		{ID: "clean-code", Label: "Clean Code", Category: "Practices"},
		{ID: "query-performance", Label: "Query Performance", Category: "Practices"},
	}
}

// Dummy render for debugging outside TUI
func Preview() {
	m := initialModel()
	fmt.Println(m.View())
}
