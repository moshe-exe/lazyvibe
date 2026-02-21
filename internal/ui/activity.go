package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/data"
)

// HeatmapMetric represents which metric to display in the heatmap.
type HeatmapMetric int

const (
	MetricMessages HeatmapMetric = iota
	MetricSessions
	MetricTools
	MetricTokens
)

// Name returns the display name for the metric.
func (m HeatmapMetric) Name() string {
	switch m {
	case MetricMessages:
		return "Messages"
	case MetricSessions:
		return "Sessions"
	case MetricTools:
		return "Tools"
	case MetricTokens:
		return "Tokens"
	}
	return "Messages"
}

// ActivityModel represents the activity heatmap panel.
type ActivityModel struct {
	data          *data.DashboardData
	focused       bool
	width         int
	height        int
	heatmapMetric HeatmapMetric
	timeRange     data.TimeRange
}

// NewActivityModel creates a new activity model.
func NewActivityModel() ActivityModel {
	return ActivityModel{}
}

// Update updates the activity data.
func (a *ActivityModel) Update(d *data.DashboardData, timeRange data.TimeRange) {
	a.data = d
	a.timeRange = timeRange
}

// SetFocused sets the focus state.
func (a *ActivityModel) SetFocused(focused bool) {
	a.focused = focused
}

// SetSize sets the panel dimensions.
func (a *ActivityModel) SetSize(width, height int) {
	a.width = width
	a.height = height
}

// CycleMetric cycles through heatmap metrics.
func (a *ActivityModel) CycleMetric() {
	a.heatmapMetric = (a.heatmapMetric + 1) % 4
}

// View renders the activity panel.
func (a ActivityModel) View() string {
	var lines []string

	// Title with panel number, metric name, and time range
	title := PanelTitleStyle.Render("Activity")
	numKey := MutedStyle.Render(" 2")
	metric := MutedStyle.Render(" [" + a.heatmapMetric.Name() + "]")
	timeRange := MutedStyle.Render(" [" + a.timeRange.String() + "]")
	lines = append(lines, title+numKey+metric+timeRange)
	lines = append(lines, "")

	if a.data == nil {
		lines = append(lines, MutedStyle.Render("Loading..."))
	} else {
		// Render the heatmap
		heatmapLines := a.renderHeatmap()
		lines = append(lines, heatmapLines...)

		// Legend for heatmap colors
		lines = append(lines, "")
		legend := a.renderLegend()
		lines = append(lines, legend)

		// Summary stats below heatmap
		lines = append(lines, "")
		lines = append(lines, a.renderHeatmapSummary())
	}

	content := strings.Join(lines, "\n")

	// Apply border
	style := PanelStyle(a.focused)
	if a.width > 0 {
		style = style.Width(a.width - 2)
	}
	if a.height > 0 {
		style = style.Height(a.height - 2)
	}

	return style.Render(content)
}

// getHeatmapData returns the activity map based on current metric and time range.
func (a ActivityModel) getHeatmapData() (map[string]int, int) {
	filteredActivity := a.data.FilterDailyActivity(a.timeRange)
	activityMap := make(map[string]int)
	maxVal := 0
	for _, day := range filteredActivity {
		var val int
		switch a.heatmapMetric {
		case MetricMessages:
			val = day.MessageCount
		case MetricSessions:
			val = day.SessionCount
		case MetricTools:
			val = day.ToolCallCount
		case MetricTokens:
			val = day.TokenCount
		}
		activityMap[day.Date] = val
		if val > maxVal {
			maxVal = val
		}
	}
	return activityMap, maxVal
}

// renderHeatmapSummary renders the summary line below the heatmap.
func (a ActivityModel) renderHeatmapSummary() string {
	filteredActivity := a.data.FilterDailyActivity(a.timeRange)
	if len(filteredActivity) == 0 {
		return ""
	}

	total := 0
	maxVal := 0
	for _, day := range filteredActivity {
		var val int
		switch a.heatmapMetric {
		case MetricMessages:
			val = day.MessageCount
		case MetricSessions:
			val = day.SessionCount
		case MetricTools:
			val = day.ToolCallCount
		case MetricTokens:
			val = day.TokenCount
		}
		total += val
		if val > maxVal {
			maxVal = val
		}
	}

	avg := 0
	if len(filteredActivity) > 0 {
		avg = total / len(filteredActivity)
	}

	// 5-char margin to align with heatmap (month label column)
	return "     " + MutedStyle.Render(fmt.Sprintf("%s · %d/day · %s max",
		formatNumber(total), avg, formatNumber(maxVal)))
}

// renderHeatmap renders a GitHub-style activity heatmap with month labels.
// Shows weeks as rows, days as columns (Mon-Sun).
func (a ActivityModel) renderHeatmap() []string {
	if a.data == nil || len(a.data.DailyActivity) == 0 {
		return []string{MutedStyle.Render("No activity data")}
	}

	// Get activity data based on current metric
	activityMap, maxVal := a.getHeatmapData()

	// Calculate how many weeks to show based on panel height
	availableHeight := a.height - 8 // Account for borders, title, legend, summary
	weeksToShow := availableHeight - 2
	if weeksToShow > 12 {
		weeksToShow = 12
	}
	if weeksToShow < 4 {
		weeksToShow = 4
	}

	// Generate dates for the last N weeks
	now := time.Now()
	// Find the most recent Saturday to align weeks
	daysUntilSunday := int(now.Weekday())
	endDate := now.AddDate(0, 0, -daysUntilSunday+6) // End on Saturday
	startDate := endDate.AddDate(0, 0, -(weeksToShow*7 - 1))

	// Use consistent block character, vary color for intensity
	block := "█"

	// Day column headers (with 5-char margin for month labels)
	dayLabels := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	var lines []string

	// Header row with day labels
	var headerRow strings.Builder
	headerRow.WriteString("     ") // 5-char margin for month labels
	for _, label := range dayLabels {
		headerRow.WriteString(MutedStyle.Render(label + " "))
	}
	lines = append(lines, headerRow.String())

	// Build the grid: N rows (weeks), 7 columns (Mon-Sun)
	// Most recent week at top, oldest at bottom
	prevMonth := ""
	for week := weeksToShow - 1; week >= 0; week-- {
		var row strings.Builder

		// Week label (show month when it changes)
		weekStartDate := startDate.AddDate(0, 0, week*7)
		weekMonth := weekStartDate.Format("Jan")
		if weekMonth != prevMonth {
			row.WriteString(MutedStyle.Render(fmt.Sprintf("%-4s ", weekMonth[:3])))
			prevMonth = weekMonth
		} else {
			row.WriteString("     ")
		}

		// Iterate through days of the week
		for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
			// Calculate the date for this cell
			dayOffset := week*7 + dayOfWeek
			cellDate := startDate.AddDate(0, 0, dayOffset)
			dateStr := cellDate.Format("2006-01-02")

			// Get activity for this date
			count := activityMap[dateStr]

			// Determine intensity level (0-3)
			var intensity int
			if maxVal > 0 && count > 0 {
				ratio := float64(count) / float64(maxVal)
				if ratio < 0.25 {
					intensity = 1
				} else if ratio < 0.5 {
					intensity = 2
				} else {
					intensity = 3
				}
			}

			// Apply color based on intensity (same block char, different colors)
			var styled string
			var color lipgloss.Color
			switch intensity {
			case 0:
				color = SurfaceDark // Empty/no activity
			case 1:
				color = TextMuted // Low activity
			case 2:
				color = Success // Medium activity (green)
			case 3:
				color = Primary // High activity (bright)
			}
			styled = lipgloss.NewStyle().Foreground(color).Render(block)

			row.WriteString(styled + styled + " ")
		}

		lines = append(lines, row.String())
	}

	return lines
}

// renderLegend renders the color legend for the heatmap.
func (a ActivityModel) renderLegend() string {
	block := "██"
	none := lipgloss.NewStyle().Foreground(SurfaceDark).Render(block)
	low := lipgloss.NewStyle().Foreground(TextMuted).Render(block)
	med := lipgloss.NewStyle().Foreground(Success).Render(block)
	high := lipgloss.NewStyle().Foreground(Primary).Render(block)

	// 5-char margin to align with heatmap
	return "     " + MutedStyle.Render("Less ") + none + " " + low + " " + med + " " + high + MutedStyle.Render(" More")
}

// GetKeybindings returns context-specific keybindings for this panel.
func (a ActivityModel) GetKeybindings() []Keybinding {
	return []Keybinding{
		{"m", "metric"},
	}
}
