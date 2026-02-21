package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/data"
	"github.com/moshe-exe/lazyvibe/internal/util"
)

// ProjectSortField represents the field to sort by.
type ProjectSortField int

const (
	ProjectSortByActivity ProjectSortField = iota
	ProjectSortByName
	ProjectSortBySessions
	ProjectSortByMessages
)

// ProjectsModel represents the projects table component.
type ProjectsModel struct {
	projects    []data.ProjectSummary
	allProjects []data.ProjectSummary // Unfiltered list
	cursor      int
	offset      int
	focused     bool
	width       int
	height      int
	sortField   ProjectSortField
	sortDesc    bool
	filterQuery string
	filterMode  bool
	timeRange   data.TimeRange
}

// NewProjectsModel creates a new projects model.
func NewProjectsModel() ProjectsModel {
	return ProjectsModel{
		sortField: ProjectSortByActivity, // Default sort by last activity
		sortDesc:  true,                  // Most recent first
	}
}

// Update updates the projects data.
func (p *ProjectsModel) Update(projects []data.ProjectSummary, timeRange data.TimeRange) {
	p.allProjects = make([]data.ProjectSummary, len(projects))
	copy(p.allProjects, projects)
	p.timeRange = timeRange
	p.applyFilter()
	p.sortProjects()

	// Reset cursor if out of bounds
	if p.cursor >= len(p.projects) {
		p.cursor = max(0, len(p.projects)-1)
	}
}

// applyFilter filters projects based on the current filter query.
func (p *ProjectsModel) applyFilter() {
	if p.filterQuery == "" {
		p.projects = make([]data.ProjectSummary, len(p.allProjects))
		copy(p.projects, p.allProjects)
		return
	}

	query := strings.ToLower(p.filterQuery)
	var filtered []data.ProjectSummary
	for _, project := range p.allProjects {
		if strings.Contains(strings.ToLower(project.ProjectName), query) {
			filtered = append(filtered, project)
		}
	}
	p.projects = filtered
}

// SetFilterMode enables or disables filter mode.
func (p *ProjectsModel) SetFilterMode(enabled bool) {
	p.filterMode = enabled
	if !enabled {
		p.filterQuery = ""
		p.applyFilter()
		p.sortProjects()
	}
}

// IsFilterMode returns whether filter mode is active.
func (p *ProjectsModel) IsFilterMode() bool {
	return p.filterMode
}

// HandleFilterInput handles a character input in filter mode.
func (p *ProjectsModel) HandleFilterInput(char string) {
	p.filterQuery += char
	p.applyFilter()
	p.sortProjects()
	p.cursor = 0
	p.offset = 0
}

// HandleFilterBackspace removes the last character from the filter query.
func (p *ProjectsModel) HandleFilterBackspace() {
	if len(p.filterQuery) > 0 {
		p.filterQuery = p.filterQuery[:len(p.filterQuery)-1]
		p.applyFilter()
		p.sortProjects()
	}
}

// GetFilterQuery returns the current filter query.
func (p *ProjectsModel) GetFilterQuery() string {
	return p.filterQuery
}

// GetFilteredCount returns filtered/total count string.
func (p *ProjectsModel) GetFilteredCount() string {
	if p.filterQuery == "" {
		return ""
	}
	return fmt.Sprintf("%d/%d", len(p.projects), len(p.allProjects))
}

// sortProjects sorts the projects based on current sort field and direction.
func (p *ProjectsModel) sortProjects() {
	sort.Slice(p.projects, func(i, j int) bool {
		var less bool
		switch p.sortField {
		case ProjectSortByName:
			less = p.projects[i].ProjectName < p.projects[j].ProjectName
		case ProjectSortBySessions:
			less = p.projects[i].SessionCount < p.projects[j].SessionCount
		case ProjectSortByMessages:
			less = p.projects[i].TotalMessages < p.projects[j].TotalMessages
		case ProjectSortByActivity:
			less = p.projects[i].LastActivity.Before(p.projects[j].LastActivity)
		}
		if p.sortDesc {
			return !less
		}
		return less
	})
}

// CycleSort cycles through sort fields.
func (p *ProjectsModel) CycleSort() {
	p.sortField = (p.sortField + 1) % 4
	p.sortProjects()
}

// ToggleSortDirection toggles between ascending and descending.
func (p *ProjectsModel) ToggleSortDirection() {
	p.sortDesc = !p.sortDesc
	p.sortProjects()
}

// SetFocused sets the focus state.
func (p *ProjectsModel) SetFocused(focused bool) {
	p.focused = focused
}

// SetSize sets the panel dimensions.
func (p *ProjectsModel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// CursorUp moves the cursor up.
func (p *ProjectsModel) CursorUp() {
	if p.cursor > 0 {
		p.cursor--
		p.ensureVisible()
	}
}

// CursorDown moves the cursor down.
func (p *ProjectsModel) CursorDown() {
	if p.cursor < len(p.projects)-1 {
		p.cursor++
		p.ensureVisible()
	}
}

// CursorUpN moves the cursor up by n items.
func (p *ProjectsModel) CursorUpN(n int) {
	p.cursor -= n
	if p.cursor < 0 {
		p.cursor = 0
	}
	p.ensureVisible()
}

// CursorDownN moves the cursor down by n items.
func (p *ProjectsModel) CursorDownN(n int) {
	p.cursor += n
	if p.cursor >= len(p.projects) {
		p.cursor = len(p.projects) - 1
	}
	if p.cursor < 0 {
		p.cursor = 0
	}
	p.ensureVisible()
}

func (p *ProjectsModel) ensureVisible() {
	visibleRows := p.visibleRows()
	if visibleRows <= 0 {
		return
	}

	if p.cursor < p.offset {
		p.offset = p.cursor
	} else if p.cursor >= p.offset+visibleRows {
		p.offset = p.cursor - visibleRows + 1
	}
}

func (p *ProjectsModel) visibleRows() int {
	// Account for border, title, header, separator
	if p.height > 0 {
		return p.height - 6
	}
	return 10
}

// sortFieldName returns the display name for the sort field.
func (p ProjectsModel) sortFieldName() string {
	switch p.sortField {
	case ProjectSortByName:
		return "Name"
	case ProjectSortBySessions:
		return "Sessions"
	case ProjectSortByMessages:
		return "Messages"
	case ProjectSortByActivity:
		return "Activity"
	}
	return "Activity"
}

// View renders the projects table.
func (p ProjectsModel) View() string {
	var lines []string

	// Title with panel number, sort indicator, and time range
	sortIndicator := "↓"
	if !p.sortDesc {
		sortIndicator = "↑"
	}

	title := PanelTitleStyle.Render("Projects")
	numKey := MutedStyle.Render(" 3")
	sortInfo := MutedStyle.Render(fmt.Sprintf(" [%s %s]", p.sortFieldName(), sortIndicator))
	timeRange := MutedStyle.Render(" [" + p.timeRange.String() + "]")

	titleLine := title + numKey + sortInfo
	// Show filter count if filtering
	if p.filterQuery != "" {
		titleLine += " " + MutedStyle.Render(p.GetFilteredCount())
	}
	titleLine += timeRange

	// Add nav hints on the right when focused
	if p.focused && !p.filterMode {
		navHints := MutedStyle.Render("j/k u/i")
		leftWidth := lipgloss.Width(titleLine)
		rightWidth := lipgloss.Width(navHints)
		contentWidth := p.width - 4 // Account for border and padding
		padding := contentWidth - leftWidth - rightWidth
		if padding < 2 {
			padding = 2
		}
		titleLine += strings.Repeat(" ", padding) + navHints
	}

	lines = append(lines, titleLine)

	// Show filter input if in filter mode
	if p.filterMode {
		filterLine := "/" + p.filterQuery + "█"
		lines = append(lines, lipgloss.NewStyle().Foreground(Primary).Render(filterLine))
	} else {
		lines = append(lines, "")
	}

	// Calculate column widths
	// Account for scrollbar (2 chars: "▓ " or "░ ")
	scrollbarW := 2
	contentWidth := p.width - 4 - scrollbarW // Account for border, padding, and scrollbar
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Column widths: Project (flex), Sessions (8), Messages (10), Last Active (12)
	// Account for selection indicator (2 chars: "▶ " or "  ")
	indicatorW := 2
	sessionsW := 8
	messagesW := 10
	lastActiveW := 12
	projectW := contentWidth - indicatorW - sessionsW - messagesW - lastActiveW - 6 // 6 for spacing
	if projectW < 10 {
		projectW = 10
	}

	// Header (with indicator spacing)
	header := fmt.Sprintf("%*s%-*s %*s %*s %*s",
		indicatorW, "",
		projectW, "Project",
		sessionsW, "Sessions",
		messagesW, "Messages",
		lastActiveW, "Last Active")
	lines = append(lines, MutedStyle.Render(header))
	lines = append(lines, MutedStyle.Render(strings.Repeat("-", contentWidth)))

	visibleRows := p.visibleRows()

	if len(p.projects) == 0 {
		lines = append(lines, MutedStyle.Render("No projects found"))
	} else {
		endIdx := min(p.offset+visibleRows, len(p.projects))

		for i := p.offset; i < endIdx; i++ {
			project := p.projects[i]
			isSelected := i == p.cursor && p.focused

			// Selection indicator
			indicator := "  "
			if isSelected {
				indicator = "▶ "
			}

			name := truncate(project.ProjectName, projectW)
			lastActive := util.FormatRelativeTime(project.LastActivity)

			row := fmt.Sprintf("%s%-*s %*d %*s %*s",
				indicator,
				projectW, name,
				sessionsW, project.SessionCount,
				messagesW, formatNumber(project.TotalMessages),
				lastActiveW, lastActive)

			if isSelected {
				row = HighlightStyle.Render(row)
			}
			lines = append(lines, row)
		}
	}

	// Build scrollbar
	scrollbar := RenderScrollbar(len(p.projects), visibleRows, p.offset, visibleRows)
	scrollbarLines := strings.Split(scrollbar, "\n")

	// Join content lines with scrollbar
	content := strings.Join(lines, "\n")

	// If we have a scrollbar, join it to the right of the content
	if len(scrollbarLines) > 0 && scrollbar != "" {
		contentLines := strings.Split(content, "\n")
		var combined []string
		// First lines are title, blank, header, separator (4 lines)
		headerLines := 4
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
	style := PanelStyle(p.focused)
	if p.width > 0 {
		style = style.Width(p.width - 2)
	}
	if p.height > 0 {
		style = style.Height(p.height - 2)
	}

	return style.Render(content)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetScrollInfo returns scroll position info for status display.
func (p ProjectsModel) GetScrollInfo() string {
	if len(p.projects) == 0 {
		return ""
	}
	return lipgloss.NewStyle().Foreground(TextMuted).Render(
		fmt.Sprintf("[%d/%d]", p.cursor+1, len(p.projects)))
}

// GetKeybindings returns context-specific keybindings for this panel.
func (p ProjectsModel) GetKeybindings() []Keybinding {
	if p.filterMode {
		return []Keybinding{
			{"esc", "clear"},
			{"enter", "apply"},
		}
	}
	return []Keybinding{
		{"s/S", "sort"},
		{"/", "filter"},
	}
}
