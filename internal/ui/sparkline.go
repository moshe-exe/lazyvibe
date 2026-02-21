package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/data"
)

// SparklineModel represents the activity sparkline component.
type SparklineModel struct {
	activity []data.DailyActivity
	focused  bool
	width    int
	height   int
}

// NewSparklineModel creates a new sparkline model.
func NewSparklineModel() SparklineModel {
	return SparklineModel{}
}

// Update updates the activity data.
func (s *SparklineModel) Update(activity []data.DailyActivity) {
	s.activity = activity
}

// SetFocused sets the focus state.
func (s *SparklineModel) SetFocused(focused bool) {
	s.focused = focused
}

// SetSize sets the panel dimensions.
func (s *SparklineModel) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// View renders the sparkline panel.
func (s SparklineModel) View() string {
	var lines []string

	// Title
	lines = append(lines, PanelTitleStyle.Render("Activity (Last 30 Days)"))
	lines = append(lines, MutedStyle.Render(strings.Repeat("-", 23)))

	if len(s.activity) == 0 {
		lines = append(lines, MutedStyle.Render("No data available"))
	} else {
		// Get last 30 days
		recent := s.activity
		if len(recent) > 30 {
			recent = recent[len(recent)-30:]
		}

		// Build sparkline data
		values := make([]int, len(recent))
		total := 0
		maxVal := 0
		for i, day := range recent {
			values[i] = day.MessageCount
			total += day.MessageCount
			if day.MessageCount > maxVal {
				maxVal = day.MessageCount
			}
		}

		avg := 0
		if len(values) > 0 {
			avg = total / len(values)
		}

		// Create block sparkline (better visual)
		sparkline := s.blockSparkline(values)
		lines = append(lines, sparkline)
		lines = append(lines, "")
		lines = append(lines, s.statLine("Total:", fmt.Sprintf("%s messages", formatNumber(total))))
		lines = append(lines, s.statLine("Average:", fmt.Sprintf("%d/day", avg)))
		lines = append(lines, s.statLine("Peak:", fmt.Sprintf("%s/day", formatNumber(maxVal))))
	}

	content := strings.Join(lines, "\n")

	// Apply border
	style := PanelStyle(s.focused)
	if s.width > 0 {
		style = style.Width(s.width - 2)
	}
	if s.height > 0 {
		style = style.Height(s.height - 2)
	}

	return style.Render(content)
}

func (s SparklineModel) statLine(label, value string) string {
	labelRendered := StatLabelStyle.Render(fmt.Sprintf("%-10s", label))
	valueRendered := StatValueStyle.Render(value)
	return labelRendered + valueRendered
}

// blockSparkline renders values using block characters with gradient colors.
// Uses 8 block heights: ▁▂▃▄▅▆▇█
func (s SparklineModel) blockSparkline(values []int) string {
	if len(values) == 0 {
		return ""
	}

	// Block characters for 8 different heights
	blocks := []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	maxVal := 0
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}

	if maxVal == 0 {
		return strings.Repeat("▁", len(values))
	}

	var result strings.Builder
	for _, v := range values {
		// Normalize to 0-8 range
		idx := (v * 8) / maxVal
		if idx > 8 {
			idx = 8
		}
		char := string(blocks[idx])

		// Apply gradient color based on value intensity
		ratio := float64(v) / float64(maxVal)
		color := GradientColor(ratio)
		styled := lipgloss.NewStyle().Foreground(color).Render(char)
		result.WriteString(styled)
	}

	return result.String()
}

// GetKeybindings returns context-specific keybindings for this panel.
func (s SparklineModel) GetKeybindings() []Keybinding {
	return []Keybinding{}
}
