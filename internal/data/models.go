// Package data provides data models and parsing for Claude Code metrics.
package data

import (
	"fmt"
	"time"
)

// TimeRange represents a time filter for data.
type TimeRange int

const (
	TimeAll TimeRange = iota
	TimeToday
	TimeWeek
	TimeMonth
)

// VMStatus represents the status of the Claude Desktop VM.
type VMStatus struct {
	Running    bool
	PID        *int
	CPUPercent *float64
	MemoryMB   *float64
}

// SessionEntry represents a Claude Code session entry.
type SessionEntry struct {
	SessionID    string
	ProjectPath  string
	ProjectName  string
	Summary      string
	MessageCount int
	Created      time.Time
	Modified     time.Time
	GitBranch    *string
}

// Duration returns the session duration based on created and modified times.
func (s SessionEntry) Duration() time.Duration {
	return s.Modified.Sub(s.Created)
}

// FormatDuration returns a human-readable duration string.
func (s SessionEntry) FormatDuration() string {
	d := s.Duration()
	if d < time.Minute {
		return "<1m"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}

// DailyActivity represents activity stats for a single day.
type DailyActivity struct {
	Date          string
	MessageCount  int
	SessionCount  int
	ToolCallCount int
	TokenCount    int // Estimated tokens
}

// ProjectSummary represents aggregated stats for a project.
type ProjectSummary struct {
	ProjectName   string
	ProjectPath   string
	SessionCount  int
	TotalMessages int
	LastActivity  time.Time
}

// DashboardData aggregates all dashboard data.
type DashboardData struct {
	VMStatus      VMStatus
	Sessions      []SessionEntry
	DailyActivity []DailyActivity
	Projects      []ProjectSummary
}

// TotalSessions returns the total number of sessions.
func (d *DashboardData) TotalSessions() int {
	return len(d.Sessions)
}

// TotalMessages returns the total messages across all sessions.
func (d *DashboardData) TotalMessages() int {
	total := 0
	for _, s := range d.Sessions {
		total += s.MessageCount
	}
	return total
}

// TotalToolCalls returns the total tool calls from daily activity.
func (d *DashboardData) TotalToolCalls() int {
	total := 0
	for _, a := range d.DailyActivity {
		total += a.ToolCallCount
	}
	return total
}

// GetLastNDays returns the last n days of daily activity.
func (d *DashboardData) GetLastNDays(n int) []DailyActivity {
	if len(d.DailyActivity) <= n {
		return d.DailyActivity
	}
	return d.DailyActivity[len(d.DailyActivity)-n:]
}

// GetMessageTrend returns message counts for the last n days.
func (d *DashboardData) GetMessageTrend(n int) []int {
	days := d.GetLastNDays(n)
	values := make([]int, len(days))
	for i, day := range days {
		values[i] = day.MessageCount
	}
	return values
}

// GetSessionTrend returns session counts for the last n days.
func (d *DashboardData) GetSessionTrend(n int) []int {
	days := d.GetLastNDays(n)
	values := make([]int, len(days))
	for i, day := range days {
		values[i] = day.SessionCount
	}
	return values
}

// GetToolCallTrend returns tool call counts for the last n days.
func (d *DashboardData) GetToolCallTrend(n int) []int {
	days := d.GetLastNDays(n)
	values := make([]int, len(days))
	for i, day := range days {
		values[i] = day.ToolCallCount
	}
	return values
}

// TotalTokens returns the total estimated tokens from daily activity.
func (d *DashboardData) TotalTokens() int {
	total := 0
	for _, a := range d.DailyActivity {
		total += a.TokenCount
	}
	return total
}

// GetTokenTrend returns token counts for the last n days.
func (d *DashboardData) GetTokenTrend(n int) []int {
	days := d.GetLastNDays(n)
	values := make([]int, len(days))
	for i, day := range days {
		values[i] = day.TokenCount
	}
	return values
}

// TimeRangeName returns a display name for the time range.
func (tr TimeRange) String() string {
	switch tr {
	case TimeToday:
		return "Today"
	case TimeWeek:
		return "This Week"
	case TimeMonth:
		return "This Month"
	default:
		return "All Time"
	}
}

// StartTime returns the start time for filtering based on the time range.
func (tr TimeRange) StartTime() time.Time {
	now := time.Now()
	switch tr {
	case TimeToday:
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case TimeWeek:
		return now.AddDate(0, 0, -7)
	case TimeMonth:
		return now.AddDate(0, -1, 0)
	default:
		return time.Time{} // Zero time means no filter
	}
}

// FilterSessions filters sessions by time range.
func (d *DashboardData) FilterSessions(tr TimeRange) []SessionEntry {
	if tr == TimeAll {
		return d.Sessions
	}
	startTime := tr.StartTime()
	var filtered []SessionEntry
	for _, s := range d.Sessions {
		if s.Modified.After(startTime) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// FilterProjects filters projects by time range.
func (d *DashboardData) FilterProjects(tr TimeRange) []ProjectSummary {
	if tr == TimeAll {
		return d.Projects
	}
	startTime := tr.StartTime()
	var filtered []ProjectSummary
	for _, p := range d.Projects {
		if p.LastActivity.After(startTime) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// FilterDailyActivity filters daily activity by time range.
func (d *DashboardData) FilterDailyActivity(tr TimeRange) []DailyActivity {
	if tr == TimeAll {
		return d.DailyActivity
	}
	startTime := tr.StartTime()
	startDate := startTime.Format("2006-01-02")
	var filtered []DailyActivity
	for _, day := range d.DailyActivity {
		if day.Date >= startDate {
			filtered = append(filtered, day)
		}
	}
	return filtered
}
