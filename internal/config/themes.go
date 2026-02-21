package config

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme.
type Theme struct {
	Name        string
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Success     lipgloss.Color
	Warning     lipgloss.Color
	Error       lipgloss.Color
	Surface     lipgloss.Color
	SurfaceDark lipgloss.Color
	Text        lipgloss.Color
	TextMuted   lipgloss.Color
	TextBright  lipgloss.Color
}

// Themes contains all available themes.
var Themes = map[string]Theme{
	"default": {
		Name:        "One Dark",
		Primary:     lipgloss.Color("#61afef"),
		Secondary:   lipgloss.Color("#c678dd"),
		Success:     lipgloss.Color("#98c379"),
		Warning:     lipgloss.Color("#e5c07b"),
		Error:       lipgloss.Color("#e06c75"),
		Surface:     lipgloss.Color("#282c34"),
		SurfaceDark: lipgloss.Color("#21252b"),
		Text:        lipgloss.Color("#abb2bf"),
		TextMuted:   lipgloss.Color("#5c6370"),
		TextBright:  lipgloss.Color("#ffffff"),
	},
	"dracula": {
		Name:        "Dracula",
		Primary:     lipgloss.Color("#bd93f9"),
		Secondary:   lipgloss.Color("#ff79c6"),
		Success:     lipgloss.Color("#50fa7b"),
		Warning:     lipgloss.Color("#f1fa8c"),
		Error:       lipgloss.Color("#ff5555"),
		Surface:     lipgloss.Color("#282a36"),
		SurfaceDark: lipgloss.Color("#21222c"),
		Text:        lipgloss.Color("#f8f8f2"),
		TextMuted:   lipgloss.Color("#6272a4"),
		TextBright:  lipgloss.Color("#ffffff"),
	},
	"nord": {
		Name:        "Nord",
		Primary:     lipgloss.Color("#88c0d0"),
		Secondary:   lipgloss.Color("#b48ead"),
		Success:     lipgloss.Color("#a3be8c"),
		Warning:     lipgloss.Color("#ebcb8b"),
		Error:       lipgloss.Color("#bf616a"),
		Surface:     lipgloss.Color("#2e3440"),
		SurfaceDark: lipgloss.Color("#242933"),
		Text:        lipgloss.Color("#d8dee9"),
		TextMuted:   lipgloss.Color("#4c566a"),
		TextBright:  lipgloss.Color("#eceff4"),
	},
	"gruvbox": {
		Name:        "Gruvbox",
		Primary:     lipgloss.Color("#83a598"),
		Secondary:   lipgloss.Color("#d3869b"),
		Success:     lipgloss.Color("#b8bb26"),
		Warning:     lipgloss.Color("#fabd2f"),
		Error:       lipgloss.Color("#fb4934"),
		Surface:     lipgloss.Color("#282828"),
		SurfaceDark: lipgloss.Color("#1d2021"),
		Text:        lipgloss.Color("#ebdbb2"),
		TextMuted:   lipgloss.Color("#665c54"),
		TextBright:  lipgloss.Color("#fbf1c7"),
	},
	"catppuccin": {
		Name:        "Catppuccin Mocha",
		Primary:     lipgloss.Color("#89b4fa"),
		Secondary:   lipgloss.Color("#cba6f7"),
		Success:     lipgloss.Color("#a6e3a1"),
		Warning:     lipgloss.Color("#f9e2af"),
		Error:       lipgloss.Color("#f38ba8"),
		Surface:     lipgloss.Color("#1e1e2e"),
		SurfaceDark: lipgloss.Color("#181825"),
		Text:        lipgloss.Color("#cdd6f4"),
		TextMuted:   lipgloss.Color("#6c7086"),
		TextBright:  lipgloss.Color("#ffffff"),
	},
}

// GetTheme returns a theme by name, falling back to default.
func GetTheme(name string) Theme {
	if theme, ok := Themes[name]; ok {
		return theme
	}
	return Themes["default"]
}

// ThemeNames returns a list of available theme names.
func ThemeNames() []string {
	return []string{"default", "dracula", "nord", "gruvbox", "catppuccin"}
}
