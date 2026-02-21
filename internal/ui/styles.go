// Package ui provides the TUI components using Bubbletea and Lipgloss.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/config"
)

// Color palette - dynamically set by theme
var (
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
)

// Common styles - regenerated when theme changes
var (
	HeaderStyle             lipgloss.Style
	PanelBorderStyle        lipgloss.Style
	PanelBorderFocusedStyle lipgloss.Style
	PanelTitleStyle         lipgloss.Style
	StatLabelStyle          lipgloss.Style
	StatValueStyle          lipgloss.Style
	SuccessStyle            lipgloss.Style
	WarningStyle            lipgloss.Style
	ErrorStyle              lipgloss.Style
	MutedStyle              lipgloss.Style
	HighlightStyle          lipgloss.Style
	HelpKeyStyle            lipgloss.Style
	HelpDescStyle           lipgloss.Style
)

// CurrentTheme holds the name of the current theme
var CurrentTheme string

func init() {
	// Initialize with default theme
	ApplyTheme("default")
}

// ApplyTheme applies a theme by name.
func ApplyTheme(name string) {
	theme := config.GetTheme(name)
	CurrentTheme = name

	// Apply colors
	Primary = theme.Primary
	Secondary = theme.Secondary
	Success = theme.Success
	Warning = theme.Warning
	Error = theme.Error
	Surface = theme.Surface
	SurfaceDark = theme.SurfaceDark
	Text = theme.Text
	TextMuted = theme.TextMuted
	TextBright = theme.TextBright

	// Regenerate styles
	HeaderStyle = lipgloss.NewStyle().
		Foreground(TextBright).
		Background(SurfaceDark).
		Bold(true).
		Padding(0, 1)

	PanelBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(TextMuted)

	PanelBorderFocusedStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary)

	PanelTitleStyle = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	StatLabelStyle = lipgloss.NewStyle().
		Foreground(TextMuted)

	StatValueStyle = lipgloss.NewStyle().
		Foreground(Text).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(Success)

	WarningStyle = lipgloss.NewStyle().
		Foreground(Warning)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(Error)

	MutedStyle = lipgloss.NewStyle().
		Foreground(TextMuted)

	HighlightStyle = lipgloss.NewStyle().
		Background(Primary).
		Foreground(SurfaceDark)

	HelpKeyStyle = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
		Foreground(TextMuted)
}

// CycleTheme cycles to the next available theme.
func CycleTheme() string {
	names := config.ThemeNames()
	for i, name := range names {
		if name == CurrentTheme {
			nextIdx := (i + 1) % len(names)
			ApplyTheme(names[nextIdx])
			return names[nextIdx]
		}
	}
	ApplyTheme("default")
	return "default"
}

// PanelStyle returns the appropriate border style based on focus state.
func PanelStyle(focused bool) lipgloss.Style {
	if focused {
		return PanelBorderFocusedStyle
	}
	return PanelBorderStyle
}

// RenderScrollbar renders a vertical scrollbar track.
// total: total items, visible: visible items, offset: first visible item, height: available height for scrollbar
func RenderScrollbar(total, visible, offset, height int) string {
	if total <= visible || height <= 0 {
		return ""
	}

	thumbSize := height * visible / total
	if thumbSize < 1 {
		thumbSize = 1
	}
	thumbPos := height * offset / total

	var sb strings.Builder
	for i := 0; i < height; i++ {
		if i >= thumbPos && i < thumbPos+thumbSize {
			sb.WriteString("▓")
		} else {
			sb.WriteString("░")
		}
		if i < height-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// Keybinding represents a key and its description for context-aware footer.
type Keybinding struct {
	Key  string
	Desc string
}

// GradientColor returns a color based on a ratio (0.0 to 1.0).
// Low values are green, medium are yellow, high are red.
func GradientColor(ratio float64) lipgloss.Color {
	switch {
	case ratio < 0.33:
		return Success // green
	case ratio < 0.66:
		return Warning // yellow
	default:
		return Error // red
	}
}

// RenderBar renders a progress bar with gradient coloring.
// percent: 0-100, width: character width of the bar
func RenderBar(percent float64, width int) string {
	if width <= 0 {
		return ""
	}
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	color := GradientColor(percent / 100)
	return lipgloss.NewStyle().Foreground(color).Render(bar)
}

// formatNumber formats a number with comma separators.
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%d,%03d,%03d", n/1000000, (n/1000)%1000, n%1000)
}
