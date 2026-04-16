package tui

// Screen represents which screen is currently visible
type Screen int

const (
	ScreenWelcome      Screen = iota
	ScreenFramework           // screen 1 — pick framework (radio)
	ScreenStack               // screen 2 — pick stack version (radio)
	ScreenArchitecture        // screen 3 — pick architecture (checkboxes)
	ScreenDrivenDesign        // screen 4 — pick driven design (radio)
	ScreenPractices           // screen 5 — pick practices (checkboxes)
	ScreenInstalling          // installing...
	ScreenDone                // done
)

// Framework options
type Framework string

const (
	FrameworkDotNet     Framework = "dotnet"
	FrameworkSpringBoot Framework = "springboot"
	FrameworkAngular    Framework = "angular"
	FrameworkSQL        Framework = "sql"
)

// Option is a single selectable item in any screen
type Option struct {
	ID    string
	Label string
}

// Selection holds everything the user has picked across all screens
type Selection struct {
	Framework    Framework
	Stack        string   // e.g. "dotnet-core-8"
	Architecture []string // e.g. ["hexagonal", "onion"]
	DrivenDesign string   // e.g. "ddd"
	Practices    []string // e.g. ["solid", "clean-code"]
}

// Model holds the complete application state
type Model struct {
	Screen    Screen
	Selection Selection
	Options   []Option
	Cursor    int
	Width     int
	Height    int
	Results   map[string]bool
	Err       error
}

// ── Framework options ─────────────────────────────────────────────────────────

var FrameworkOptions = []Option{
	{ID: "dotnet", Label: ".NET  —  C#"},
	{ID: "springboot", Label: "Spring Boot  —  Java"},
	{ID: "angular", Label: "Angular  —  TypeScript"},
	{ID: "sql", Label: "SQL  —  MS SQL Server"},
}

// ── Stack options per framework ───────────────────────────────────────────────

var StackOptions = map[Framework][]Option{
	FrameworkDotNet: {
		{ID: "dotnet-framework-45", Label: ".NET Framework 4.5"},
		{ID: "dotnet-core-8", Label: ".NET Core 8"},
		{ID: "aspnet-core", Label: "ASP.NET Core"},
	},
	FrameworkSpringBoot: {
		{ID: "spring-boot-2-java17", Label: "Spring Boot 2  —  Java 17"},
		{ID: "spring-boot-3-java21", Label: "Spring Boot 3  —  Java 21"},
	},
	FrameworkAngular: {
		{ID: "angular-14", Label: "Angular 14"},
		{ID: "angular-17", Label: "Angular 17"},
	},
	FrameworkSQL: {
		{ID: "mssql-2019", Label: "MS SQL Server 2019"},
	},
}

// ── Architecture options ──────────────────────────────────────────────────────

func ArchitectureOptions(fw Framework) []Option {
	opts := []Option{
		{ID: "clean-architecture", Label: "Clean Architecture"},
		{ID: "hexagonal", Label: "Hexagonal Architecture"},
		{ID: "onion", Label: "Onion Architecture"},
	}
	if fw == FrameworkSpringBoot {
		opts = append(opts, Option{ID: "microservices", Label: "Microservices"})
	}
	return opts
}

// ArchitectureConflicts — selecting one blocks the other
var ArchitectureConflicts = map[string]string{
	"hexagonal":          "clean-architecture",
	"clean-architecture": "hexagonal",
}

// ── Driven design options ─────────────────────────────────────────────────────

var DrivenDesignOptions = []Option{
	{ID: "tdd", Label: "Test-Driven Design"},
	{ID: "edd", Label: "Event-Driven Design"},
	{ID: "ddd", Label: "Domain-Driven Design"},
}

// ── Practices options ─────────────────────────────────────────────────────────

var PracticesOptions = []Option{
	{ID: "solid", Label: "SOLID Principles"},
	{ID: "dependency-injection", Label: "Dependency Injection"},
	{ID: "clean-code", Label: "Clean Code"},
	{ID: "oop", Label: "OOP"},
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func remove(slice []string, val string) []string {
	out := []string{}
	for _, v := range slice {
		if v != val {
			out = append(out, v)
		}
	}
	return out
}

// Sword art in braille
const SwordArt = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀  ⣀⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀  ⠙⣿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀ ⣠⣶⣿⣿⣿⣶⣄⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿ ⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠙⣿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀
`

// Version of the app
const Version = "0.1.0"
