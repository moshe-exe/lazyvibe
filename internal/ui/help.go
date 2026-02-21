package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpModel represents the help modal component.
type HelpModel struct {
	visible bool
	width   int
	height  int
}

// NewHelpModel creates a new help model.
func NewHelpModel() HelpModel {
	return HelpModel{}
}

// Toggle toggles the help visibility.
func (h *HelpModel) Toggle() {
	h.visible = !h.visible
}

// Show shows the help modal.
func (h *HelpModel) Show() {
	h.visible = true
}

// Hide hides the help modal.
func (h *HelpModel) Hide() {
	h.visible = false
}

// IsVisible returns whether the help is visible.
func (h HelpModel) IsVisible() bool {
	return h.visible
}

// SetSize sets the available dimensions.
func (h *HelpModel) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// View renders the help modal.
func (h HelpModel) View() string {
	if !h.visible {
		return ""
	}

	var lines []string

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(Primary)
	lines = append(lines, titleStyle.Render("Claude Code Monitor - Keyboard Shortcuts"))
	lines = append(lines, "")

	// Panel Navigation
	sectionStyle := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(Text)
	lines = append(lines, sectionStyle.Render("Panel Navigation"))
	lines = append(lines, helpLine("1", "Jump to Stats panel"))
	lines = append(lines, helpLine("2", "Jump to Activity panel"))
	lines = append(lines, helpLine("3", "Jump to Projects panel"))
	lines = append(lines, helpLine("4", "Jump to Sessions panel"))
	lines = append(lines, helpLine("Tab", "Next panel"))
	lines = append(lines, helpLine("Shift+Tab", "Previous panel"))
	lines = append(lines, "")

	// Vim-Style Movement
	lines = append(lines, sectionStyle.Render("Vim-Style Movement"))
	lines = append(lines, helpLine("h", "Move left"))
	lines = append(lines, helpLine("l", "Move right"))
	lines = append(lines, helpLine("j", "Move down / Next item"))
	lines = append(lines, helpLine("k", "Move up / Previous item"))
	lines = append(lines, helpLine("g", "Go to top of list"))
	lines = append(lines, helpLine("G", "Go to bottom of list"))
	lines = append(lines, "")

	// General
	lines = append(lines, sectionStyle.Render("General"))
	lines = append(lines, helpLine("r", "Force refresh all data"))
	lines = append(lines, helpLine("p", "Pause/resume auto-refresh"))
	lines = append(lines, helpLine("?", "Toggle this help"))
	lines = append(lines, helpLine("q", "Quit"))
	lines = append(lines, "")

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(TextMuted)
	lines = append(lines, footerStyle.Render("Press Esc or ? to close this help"))

	content := strings.Join(lines, "\n")

	// Modal box style
	modalWidth := 50

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(Secondary).
		Background(Surface).
		Padding(1, 2).
		Width(modalWidth)

	modal := boxStyle.Render(content)

	// Center the modal
	if h.width > 0 && h.height > 0 {
		modalLines := strings.Split(modal, "\n")
		modalActualHeight := len(modalLines)

		// Vertical centering
		topPadding := (h.height - modalActualHeight) / 2
		if topPadding < 0 {
			topPadding = 0
		}

		// Horizontal centering
		leftPadding := (h.width - modalWidth - 4) / 2 // Account for border
		if leftPadding < 0 {
			leftPadding = 0
		}

		var centered []string
		for i := 0; i < topPadding; i++ {
			centered = append(centered, "")
		}
		for _, line := range modalLines {
			centered = append(centered, strings.Repeat(" ", leftPadding)+line)
		}

		return strings.Join(centered, "\n")
	}

	return modal
}

func helpLine(key, desc string) string {
	keyRendered := HelpKeyStyle.Render(padRight(key, 12))
	descRendered := HelpDescStyle.Render(desc)
	return "  " + keyRendered + descRendered
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}
