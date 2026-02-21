package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/data"
)

// DetailModal represents a modal for showing session details.
type DetailModal struct {
	visible bool
	session *data.SessionEntry
	width   int
	height  int
}

// NewDetailModal creates a new detail modal.
func NewDetailModal() DetailModal {
	return DetailModal{}
}

// SetSize sets the modal dimensions.
func (d *DetailModal) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// Show displays the modal with the given session.
func (d *DetailModal) Show(session *data.SessionEntry) {
	d.session = session
	d.visible = true
}

// Hide hides the modal.
func (d *DetailModal) Hide() {
	d.visible = false
	d.session = nil
}

// IsVisible returns whether the modal is visible.
func (d *DetailModal) IsVisible() bool {
	return d.visible
}

// View renders the detail modal.
func (d DetailModal) View() string {
	if !d.visible || d.session == nil {
		return ""
	}

	// Modal dimensions
	modalWidth := d.width * 60 / 100
	if modalWidth < 50 {
		modalWidth = 50
	}
	if modalWidth > 80 {
		modalWidth = 80
	}

	modalHeight := 16

	var lines []string

	// Title
	lines = append(lines, PanelTitleStyle.Render("Session Details"))
	lines = append(lines, MutedStyle.Render(strings.Repeat("-", modalWidth-4)))
	lines = append(lines, "")

	// Session info
	session := d.session

	lines = append(lines, d.detailLine("Session ID:", session.SessionID))
	lines = append(lines, d.detailLine("Project:", session.ProjectName))
	lines = append(lines, d.detailLine("Path:", truncateMiddle(session.ProjectPath, modalWidth-20)))

	if session.GitBranch != nil && *session.GitBranch != "" {
		lines = append(lines, d.detailLine("Branch:", *session.GitBranch))
	}

	lines = append(lines, "")
	lines = append(lines, d.detailLine("Messages:", fmt.Sprintf("%d", session.MessageCount)))
	lines = append(lines, d.detailLine("Duration:", session.FormatDuration()))
	lines = append(lines, d.detailLine("Created:", session.Created.Format("2006-01-02 15:04")))
	lines = append(lines, d.detailLine("Modified:", session.Modified.Format("2006-01-02 15:04")))

	lines = append(lines, "")

	// Summary
	lines = append(lines, d.detailLine("Summary:", ""))
	summary := session.Summary
	if len(summary) > modalWidth-6 {
		summary = summary[:modalWidth-9] + "..."
	}
	lines = append(lines, "  "+MutedStyle.Render(summary))

	content := strings.Join(lines, "\n")

	// Modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight)

	modal := modalStyle.Render(content)

	// Center the modal
	paddingLeft := (d.width - modalWidth) / 2
	paddingTop := (d.height - modalHeight - 4) / 2

	if paddingLeft < 0 {
		paddingLeft = 0
	}
	if paddingTop < 0 {
		paddingTop = 0
	}

	// Create padding
	topPadding := strings.Repeat("\n", paddingTop)
	leftPadding := strings.Repeat(" ", paddingLeft)

	// Apply left padding to each line
	lines = strings.Split(modal, "\n")
	for i, line := range lines {
		lines[i] = leftPadding + line
	}

	return topPadding + strings.Join(lines, "\n")
}

func (d DetailModal) detailLine(label, value string) string {
	labelRendered := StatLabelStyle.Render(fmt.Sprintf("%-12s", label))
	valueRendered := StatValueStyle.Render(value)
	return labelRendered + valueRendered
}

// truncateMiddle truncates a string in the middle if too long.
func truncateMiddle(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 5 {
		return s[:maxLen]
	}
	half := (maxLen - 3) / 2
	return s[:half] + "..." + s[len(s)-half:]
}
