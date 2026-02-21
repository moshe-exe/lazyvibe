package ui

import (
	"fmt"
	"strings"

	"github.com/moshe-exe/lazyvibe/internal/data"
)

// StatsModel represents the stats panel showing metrics.
type StatsModel struct {
	data      *data.DashboardData
	focused   bool
	width     int
	height    int
	timeRange data.TimeRange
}

// NewStatsModel creates a new stats model.
func NewStatsModel() StatsModel {
	return StatsModel{}
}

// Update updates the stats data.
func (s *StatsModel) Update(d *data.DashboardData, timeRange data.TimeRange) {
	s.data = d
	s.timeRange = timeRange
}

// SetFocused sets the focus state.
func (s *StatsModel) SetFocused(focused bool) {
	s.focused = focused
}

// SetSize sets the panel dimensions.
func (s *StatsModel) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// View renders the stats panel.
func (s StatsModel) View() string {
	var lines []string

	// Title with panel number and time range
	title := PanelTitleStyle.Render("Stats")
	numKey := MutedStyle.Render(" 1")
	timeRange := MutedStyle.Render(" [" + s.timeRange.String() + "]")
	lines = append(lines, title+numKey+timeRange)
	lines = append(lines, "")

	if s.data == nil {
		lines = append(lines, MutedStyle.Render("Loading..."))
	} else {
		lines = append(lines, s.renderMetrics()...)
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

// renderMetrics renders the metrics as a single column.
func (s StatsModel) renderMetrics() []string {
	// Get filtered data based on time range
	filteredProjects := s.data.FilterProjects(s.timeRange)
	filteredSessions := s.data.FilterSessions(s.timeRange)
	filteredActivity := s.data.FilterDailyActivity(s.timeRange)

	// Calculate totals from filtered data
	totalMessages := 0
	for _, sess := range filteredSessions {
		totalMessages += sess.MessageCount
	}

	totalTools := 0
	totalTokens := 0
	for _, a := range filteredActivity {
		totalTools += a.ToolCallCount
		totalTokens += a.TokenCount
	}

	return []string{
		s.metricLine("Projects", fmt.Sprintf("%d", len(filteredProjects))),
		s.metricLine("Sessions", fmt.Sprintf("%d", len(filteredSessions))),
		s.metricLine("Messages", formatNumber(totalMessages)),
		s.metricLine("Tools", formatNumber(totalTools)),
		s.metricLine("Tokens", formatTokens(totalTokens)),
	}
}

// metricLine renders a single metric line with 2-char margin.
func (s StatsModel) metricLine(label, value string) string {
	l := StatLabelStyle.Render(fmt.Sprintf("%-10s", label))
	v := StatValueStyle.Render(value)
	return "  " + l + v
}

// formatTokens formats token count with K/M suffix.
func formatTokens(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("~%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("~%.0fK", float64(n)/1000)
	}
	return fmt.Sprintf("~%d", n)
}

// GetKeybindings returns context-specific keybindings for this panel.
func (s StatsModel) GetKeybindings() []Keybinding {
	return []Keybinding{}
}
