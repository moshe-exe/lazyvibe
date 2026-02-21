package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/data"
	"github.com/moshe-exe/lazyvibe/internal/util"
)

// SessionSortField represents the field to sort sessions by.
type SessionSortField int

const (
	SessionSortByTime SessionSortField = iota
	SessionSortByMessages
	SessionSortByProject
)

// SessionsModel represents the recent sessions component.
type SessionsModel struct {
	sessions    []data.SessionEntry
	allSessions []data.SessionEntry // Unfiltered list
	cursor      int
	offset      int
	focused     bool
	width       int
	height      int
	sortField   SessionSortField
	sortDesc    bool
	filterQuery string
	filterMode  bool
	timeRange   data.TimeRange
}

// NewSessionsModel creates a new sessions model.
func NewSessionsModel() SessionsModel {
	return SessionsModel{
		sortField: SessionSortByTime,
		sortDesc:  true, // Most recent first
	}
}

// Update updates the sessions data.
func (s *SessionsModel) Update(sessions []data.SessionEntry, timeRange data.TimeRange) {
	// Copy and keep only last 20
	sorted := make([]data.SessionEntry, len(sessions))
	copy(sorted, sessions)

	// Initial sort by time to get most recent
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Modified.After(sorted[j].Modified)
	})

	if len(sorted) > 20 {
		sorted = sorted[:20]
	}

	s.allSessions = sorted
	s.timeRange = timeRange
	s.applyFilter()
	s.sortSessions()

	// Reset cursor if out of bounds
	if s.cursor >= len(s.sessions) {
		s.cursor = max(0, len(s.sessions)-1)
	}
}

// applyFilter filters sessions based on the current filter query.
func (s *SessionsModel) applyFilter() {
	if s.filterQuery == "" {
		s.sessions = make([]data.SessionEntry, len(s.allSessions))
		copy(s.sessions, s.allSessions)
		return
	}

	query := strings.ToLower(s.filterQuery)
	var filtered []data.SessionEntry
	for _, session := range s.allSessions {
		if strings.Contains(strings.ToLower(session.Summary), query) ||
			strings.Contains(strings.ToLower(session.ProjectName), query) {
			filtered = append(filtered, session)
		}
	}
	s.sessions = filtered
}

// SetFilterMode enables or disables filter mode.
func (s *SessionsModel) SetFilterMode(enabled bool) {
	s.filterMode = enabled
	if !enabled {
		s.filterQuery = ""
		s.applyFilter()
		s.sortSessions()
	}
}

// IsFilterMode returns whether filter mode is active.
func (s *SessionsModel) IsFilterMode() bool {
	return s.filterMode
}

// HandleFilterInput handles a character input in filter mode.
func (s *SessionsModel) HandleFilterInput(char string) {
	s.filterQuery += char
	s.applyFilter()
	s.sortSessions()
	s.cursor = 0
	s.offset = 0
}

// HandleFilterBackspace removes the last character from the filter query.
func (s *SessionsModel) HandleFilterBackspace() {
	if len(s.filterQuery) > 0 {
		s.filterQuery = s.filterQuery[:len(s.filterQuery)-1]
		s.applyFilter()
		s.sortSessions()
	}
}

// GetFilterQuery returns the current filter query.
func (s *SessionsModel) GetFilterQuery() string {
	return s.filterQuery
}

// GetFilteredCount returns filtered/total count string.
func (s *SessionsModel) GetFilteredCount() string {
	if s.filterQuery == "" {
		return ""
	}
	return fmt.Sprintf("%d/%d", len(s.sessions), len(s.allSessions))
}

// sortSessions sorts the sessions based on current sort field and direction.
func (s *SessionsModel) sortSessions() {
	sort.Slice(s.sessions, func(i, j int) bool {
		var less bool
		switch s.sortField {
		case SessionSortByTime:
			less = s.sessions[i].Modified.Before(s.sessions[j].Modified)
		case SessionSortByMessages:
			less = s.sessions[i].MessageCount < s.sessions[j].MessageCount
		case SessionSortByProject:
			less = s.sessions[i].ProjectName < s.sessions[j].ProjectName
		}
		if s.sortDesc {
			return !less
		}
		return less
	})
}

// CycleSort cycles through sort fields.
func (s *SessionsModel) CycleSort() {
	s.sortField = (s.sortField + 1) % 3
	s.sortSessions()
}

// ToggleSortDirection toggles between ascending and descending.
func (s *SessionsModel) ToggleSortDirection() {
	s.sortDesc = !s.sortDesc
	s.sortSessions()
}

// sortFieldName returns the display name for the sort field.
func (s SessionsModel) sortFieldName() string {
	switch s.sortField {
	case SessionSortByTime:
		return "Time"
	case SessionSortByMessages:
		return "Messages"
	case SessionSortByProject:
		return "Project"
	}
	return "Time"
}

// SetFocused sets the focus state.
func (s *SessionsModel) SetFocused(focused bool) {
	s.focused = focused
}

// SetSize sets the panel dimensions.
func (s *SessionsModel) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// CursorUp moves the cursor up.
func (s *SessionsModel) CursorUp() {
	if s.cursor > 0 {
		s.cursor--
		s.ensureVisible()
	}
}

// CursorDown moves the cursor down.
func (s *SessionsModel) CursorDown() {
	if s.cursor < len(s.sessions)-1 {
		s.cursor++
		s.ensureVisible()
	}
}

// CursorUpN moves the cursor up by n items.
func (s *SessionsModel) CursorUpN(n int) {
	s.cursor -= n
	if s.cursor < 0 {
		s.cursor = 0
	}
	s.ensureVisible()
}

// CursorDownN moves the cursor down by n items.
func (s *SessionsModel) CursorDownN(n int) {
	s.cursor += n
	if s.cursor >= len(s.sessions) {
		s.cursor = len(s.sessions) - 1
	}
	if s.cursor < 0 {
		s.cursor = 0
	}
	s.ensureVisible()
}

func (s *SessionsModel) ensureVisible() {
	visibleRows := s.visibleRows()
	if visibleRows <= 0 {
		return
	}

	if s.cursor < s.offset {
		s.offset = s.cursor
	} else if s.cursor >= s.offset+visibleRows {
		s.offset = s.cursor - visibleRows + 1
	}
}

func (s *SessionsModel) visibleRows() int {
	// Each session takes 2 lines, account for border and title
	if s.height > 0 {
		return (s.height - 4) / 2
	}
	return 5
}

// View renders the sessions list.
func (s SessionsModel) View() string {
	var lines []string

	// Title with panel number, sort indicator, and time range
	sortIndicator := "↓"
	if !s.sortDesc {
		sortIndicator = "↑"
	}

	title := PanelTitleStyle.Render("Sessions")
	numKey := MutedStyle.Render(" 4")
	sortInfo := MutedStyle.Render(fmt.Sprintf(" [%s %s]", s.sortFieldName(), sortIndicator))
	timeRange := MutedStyle.Render(" [" + s.timeRange.String() + "]")

	titleLine := title + numKey + sortInfo
	// Show filter count if filtering
	if s.filterQuery != "" {
		titleLine += " " + MutedStyle.Render(s.GetFilteredCount())
	}
	titleLine += timeRange

	// Add nav hints on the right when focused
	if s.focused && !s.filterMode {
		navHints := MutedStyle.Render("j/k u/i")
		leftWidth := lipgloss.Width(titleLine)
		rightWidth := lipgloss.Width(navHints)
		contentWidth := s.width - 4 // Account for border and padding
		padding := contentWidth - leftWidth - rightWidth
		if padding < 2 {
			padding = 2
		}
		titleLine += strings.Repeat(" ", padding) + navHints
	}

	lines = append(lines, titleLine)

	// Show filter input if in filter mode
	if s.filterMode {
		filterLine := "/" + s.filterQuery + "█"
		lines = append(lines, lipgloss.NewStyle().Foreground(Primary).Render(filterLine))
	} else {
		lines = append(lines, "")
	}

	visibleRows := s.visibleRows()

	// Account for scrollbar width
	scrollbarW := 2
	contentWidth := s.width - 4 - scrollbarW
	if contentWidth < 30 {
		contentWidth = 30
	}

	if len(s.sessions) == 0 {
		lines = append(lines, MutedStyle.Render("No sessions found"))
	} else {
		endIdx := min(s.offset+visibleRows, len(s.sessions))

		for i := s.offset; i < endIdx; i++ {
			session := s.sessions[i]
			isSelected := i == s.cursor && s.focused

			// Selection indicator
			indicator := "  "
			if isSelected {
				indicator = "▶ "
			}

			relativeTime := util.FormatRelativeTimeShort(session.Modified)

			// First line: time, summary, branch (account for indicator width)
			summary := session.Summary
			maxSummaryLen := contentWidth - len(relativeTime) - 12 // -2 more for indicator
			if maxSummaryLen < 10 {
				maxSummaryLen = 10
			}
			if len(summary) > maxSummaryLen {
				summary = summary[:maxSummaryLen-3] + "..."
			}

			branch := ""
			if session.GitBranch != nil && *session.GitBranch != "" {
				branch = fmt.Sprintf(" [%s]", truncate(*session.GitBranch, 15))
			}

			line1 := fmt.Sprintf("%s%s: %s%s", indicator, relativeTime, summary, branch)

			// Second line: project name, message count, duration (indented to align with content)
			line2 := fmt.Sprintf("    %s | %d msgs | %s", truncate(session.ProjectName, 18), session.MessageCount, session.FormatDuration())

			if isSelected {
				line1 = HighlightStyle.Render(line1)
				line2 = HighlightStyle.Render(line2)
			} else {
				line2 = MutedStyle.Render(line2)
			}

			lines = append(lines, line1)
			lines = append(lines, line2)
		}
	}

	// Build scrollbar - each session is 2 lines, so scrollbar height is visibleRows*2
	scrollbarHeight := visibleRows * 2
	scrollbar := RenderScrollbar(len(s.sessions), visibleRows, s.offset, scrollbarHeight)
	scrollbarLines := strings.Split(scrollbar, "\n")

	// Join content lines with scrollbar
	content := strings.Join(lines, "\n")

	// If we have a scrollbar, join it to the right of the content
	if len(scrollbarLines) > 0 && scrollbar != "" {
		contentLines := strings.Split(content, "\n")
		var combined []string
		// First lines are title, blank (2 lines)
		headerLines := 2
		for i, line := range contentLines {
			scrollChar := " "
			if i >= headerLines && i-headerLines < len(scrollbarLines) {
				scrollChar = scrollbarLines[i-headerLines]
			}
			combined = append(combined, line+" "+scrollChar)
		}
		content = strings.Join(combined, "\n")
	}

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

// GetScrollInfo returns scroll position info for status display.
func (s SessionsModel) GetScrollInfo() string {
	if len(s.sessions) == 0 {
		return ""
	}
	return lipgloss.NewStyle().Foreground(TextMuted).Render(
		fmt.Sprintf("[%d/%d]", s.cursor+1, len(s.sessions)))
}

// GetKeybindings returns context-specific keybindings for this panel.
func (s SessionsModel) GetKeybindings() []Keybinding {
	if s.filterMode {
		return []Keybinding{
			{"esc", "clear"},
			{"enter", "apply"},
		}
	}
	return []Keybinding{
		{"s/S", "sort"},
		{"/", "filter"},
		{"y", "copy"},
		{"enter", "details"},
	}
}

// GetSelected returns the currently selected session.
func (s SessionsModel) GetSelected() *data.SessionEntry {
	if len(s.sessions) == 0 || s.cursor >= len(s.sessions) {
		return nil
	}
	return &s.sessions[s.cursor]
}
