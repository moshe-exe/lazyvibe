package ui

import (
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/data"
)

// Panel indices
const (
	PanelStats = iota
	PanelActivity
	PanelProjects
	PanelSessions
	panelCount = 4
)

// Panel grid layout for vim navigation
// Left: Stats (top), Activity (bottom)
// Right: Projects (top), Sessions (bottom)
var panelGrid = [][]int{
	{PanelStats, PanelProjects},
	{PanelActivity, PanelSessions},
}

// Messages for timer-based updates
type vmTickMsg struct{}
type sessionsTickMsg struct{}

// Model is the main application model.
type Model struct {
	dataManager *data.Manager
	dashData    *data.DashboardData

	width  int
	height int

	focused   int
	paused    bool
	showHelp  bool
	timeRange data.TimeRange

	// Flash message for status updates
	flashMessage string
	flashExpiry  time.Time

	// Sub-models
	header   HeaderModel
	stats    StatsModel
	activity ActivityModel
	projects ProjectsModel
	sessions SessionsModel
	help     HelpModel
	detail   DetailModal

	// Tickers
	vmTicker       *time.Ticker
	sessionsTicker *time.Ticker
}

// NewModel creates a new application model.
func NewModel(dataManager *data.Manager) Model {
	return Model{
		dataManager: dataManager,
		focused:     PanelStats,
		header:      NewHeaderModel(),
		stats:       NewStatsModel(),
		activity:    NewActivityModel(),
		projects:    NewProjectsModel(),
		sessions:    NewSessionsModel(),
		help:        NewHelpModel(),
		detail:      NewDetailModal(),
	}
}

// Init initializes the application.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadData(),
		m.vmTickCmd(),
		m.sessionsTickCmd(),
	)
}

func (m Model) loadData() tea.Cmd {
	return func() tea.Msg {
		return m.dataManager.GetDashboardData(false)
	}
}

func (m Model) vmTickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return vmTickMsg{}
	})
}

func (m Model) sessionsTickCmd() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return sessionsTickMsg{}
	})
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeypress(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
		return m, nil

	case data.DashboardData:
		m.dashData = &msg
		m.updateWidgets()
		return m, nil

	case vmTickMsg:
		if !m.paused {
			vmStatus := m.dataManager.GetVMStatus(false)
			m.header.Update(vmStatus, m.paused)
		}
		return m, m.vmTickCmd()

	case sessionsTickMsg:
		if !m.paused {
			return m, m.loadData()
		}
		return m, m.sessionsTickCmd()
	}

	return m, nil
}

func (m Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Detail modal intercepts keys when visible
	if m.detail.IsVisible() {
		switch msg.String() {
		case "esc", "enter", "q":
			m.detail.Hide()
		case "y":
			// Copy session ID from detail view
			if m.detail.session != nil {
				m.copySessionIDDirect(m.detail.session.SessionID)
			}
		}
		return m, nil
	}

	// Help modal intercepts keys when visible
	if m.help.IsVisible() {
		switch msg.String() {
		case "esc", "?":
			m.help.Hide()
		}
		return m, nil
	}

	// Handle filter mode input
	if m.isFilterMode() {
		return m.handleFilterInput(msg)
	}

	switch msg.String() {
	// Quit
	case "q", "ctrl+c":
		return m, tea.Quit

	// Refresh
	case "r":
		dashData := m.dataManager.RefreshAll()
		m.dashData = &dashData
		m.updateWidgets()
		return m, nil

	// Pause/resume
	case "p":
		m.paused = !m.paused
		if m.dashData != nil {
			m.header.Update(m.dashData.VMStatus, m.paused)
		}
		return m, nil

	// Panel jump
	case "1":
		m.focusPanel(PanelStats)
	case "2":
		m.focusPanel(PanelActivity)
	case "3":
		m.focusPanel(PanelProjects)
	case "4":
		m.focusPanel(PanelSessions)

	// Vim navigation between panels
	case "h":
		m.navLeft()
	case "l":
		m.navRight()

	// Tab navigation
	case "tab":
		m.focusNext()
	case "shift+tab":
		m.focusPrevious()

	// Activity panel: cycle metric
	case "m":
		if m.focused == PanelActivity {
			m.activity.CycleMetric()
		}

	// Navigation within lists (j=up, k=down)
	case "j", "up":
		m.cursorUp()
	case "k", "down":
		m.cursorDown()
	case "u":
		m.cursorUp5()
	case "i":
		m.cursorDown5()

	// Sort
	case "s":
		m.cycleSort()
	case "S":
		m.toggleSortDirection()

	// Time range
	case "t":
		m.cycleTimeRange()

	// Theme
	case "T":
		CycleTheme()

	// Filter
	case "/":
		m.enterFilterMode()

	// Copy session ID
	case "y":
		m.copySessionID()

	// Open detail modal
	case "enter":
		m.openDetailModal()

	// Help
	case "?":
		m.help.Toggle()
	}

	return m, nil
}

func (m *Model) focusPanel(index int) {
	m.focused = index
	m.updateFocusStates()
}

func (m *Model) focusNext() {
	m.focused = (m.focused + 1) % panelCount
	m.updateFocusStates()
}

func (m *Model) focusPrevious() {
	m.focused = (m.focused - 1 + panelCount) % panelCount
	m.updateFocusStates()
}

func (m *Model) navLeft() {
	row, col := m.getPanelPosition(m.focused)
	if col > 0 {
		m.focused = panelGrid[row][col-1]
		m.updateFocusStates()
	}
}

func (m *Model) navRight() {
	row, col := m.getPanelPosition(m.focused)
	if col < len(panelGrid[row])-1 {
		m.focused = panelGrid[row][col+1]
		m.updateFocusStates()
	}
}

func (m Model) getPanelPosition(index int) (row, col int) {
	for r, rowPanels := range panelGrid {
		for c, panelIdx := range rowPanels {
			if panelIdx == index {
				return r, c
			}
		}
	}
	return 0, 0
}

func (m *Model) cursorDown() {
	switch m.focused {
	case PanelProjects:
		m.projects.CursorDown()
	case PanelSessions:
		m.sessions.CursorDown()
	}
}

func (m *Model) cursorUp() {
	switch m.focused {
	case PanelProjects:
		m.projects.CursorUp()
	case PanelSessions:
		m.sessions.CursorUp()
	}
}

func (m *Model) cursorUp5() {
	switch m.focused {
	case PanelProjects:
		m.projects.CursorUpN(5)
	case PanelSessions:
		m.sessions.CursorUpN(5)
	}
}

func (m *Model) cursorDown5() {
	switch m.focused {
	case PanelProjects:
		m.projects.CursorDownN(5)
	case PanelSessions:
		m.sessions.CursorDownN(5)
	}
}

func (m *Model) cycleSort() {
	switch m.focused {
	case PanelProjects:
		m.projects.CycleSort()
	case PanelSessions:
		m.sessions.CycleSort()
	}
}

func (m *Model) toggleSortDirection() {
	switch m.focused {
	case PanelProjects:
		m.projects.ToggleSortDirection()
	case PanelSessions:
		m.sessions.ToggleSortDirection()
	}
}

func (m *Model) cycleTimeRange() {
	m.timeRange = (m.timeRange + 1) % 4
	m.updateWidgets()
}

func (m *Model) isFilterMode() bool {
	switch m.focused {
	case PanelProjects:
		return m.projects.IsFilterMode()
	case PanelSessions:
		return m.sessions.IsFilterMode()
	}
	return false
}

func (m *Model) enterFilterMode() {
	switch m.focused {
	case PanelProjects:
		m.projects.SetFilterMode(true)
	case PanelSessions:
		m.sessions.SetFilterMode(true)
	}
}

func (m Model) handleFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		// Clear filter and exit filter mode
		switch m.focused {
		case PanelProjects:
			m.projects.SetFilterMode(false)
		case PanelSessions:
			m.sessions.SetFilterMode(false)
		}
	case "enter":
		// Apply filter and exit filter mode (keep filter active)
		switch m.focused {
		case PanelProjects:
			m.projects.filterMode = false
		case PanelSessions:
			m.sessions.filterMode = false
		}
	case "backspace":
		switch m.focused {
		case PanelProjects:
			m.projects.HandleFilterBackspace()
		case PanelSessions:
			m.sessions.HandleFilterBackspace()
		}
	default:
		// Add character to filter if it's a printable character
		if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
			switch m.focused {
			case PanelProjects:
				m.projects.HandleFilterInput(key)
			case PanelSessions:
				m.sessions.HandleFilterInput(key)
			}
		}
	}

	return m, nil
}

func (m *Model) copySessionID() {
	if m.focused != PanelSessions {
		return
	}

	session := m.sessions.GetSelected()
	if session == nil {
		return
	}

	// Copy to clipboard using pbcopy (macOS)
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		m.flashMessage = "Copy failed"
		m.flashExpiry = time.Now().Add(2 * time.Second)
		return
	}

	if err := cmd.Start(); err != nil {
		m.flashMessage = "Copy failed"
		m.flashExpiry = time.Now().Add(2 * time.Second)
		return
	}

	stdin.Write([]byte(session.SessionID))
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		m.flashMessage = "Copy failed"
		m.flashExpiry = time.Now().Add(2 * time.Second)
		return
	}

	// Show truncated session ID in flash message
	id := session.SessionID
	if len(id) > 12 {
		id = id[:12] + "..."
	}
	m.flashMessage = "Copied: " + id
	m.flashExpiry = time.Now().Add(2 * time.Second)
}

func (m *Model) copySessionIDDirect(sessionID string) {
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		m.flashMessage = "Copy failed"
		m.flashExpiry = time.Now().Add(2 * time.Second)
		return
	}

	if err := cmd.Start(); err != nil {
		m.flashMessage = "Copy failed"
		m.flashExpiry = time.Now().Add(2 * time.Second)
		return
	}

	stdin.Write([]byte(sessionID))
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		m.flashMessage = "Copy failed"
		m.flashExpiry = time.Now().Add(2 * time.Second)
		return
	}

	id := sessionID
	if len(id) > 12 {
		id = id[:12] + "..."
	}
	m.flashMessage = "Copied: " + id
	m.flashExpiry = time.Now().Add(2 * time.Second)
}

func (m *Model) openDetailModal() {
	if m.focused != PanelSessions {
		return
	}

	session := m.sessions.GetSelected()
	if session != nil {
		m.detail.Show(session)
	}
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Ignore mouse in modal mode
	if m.help.IsVisible() || m.detail.IsVisible() {
		return m, nil
	}

	// Calculate panel boundaries
	headerHeight := 1
	leftWidth := m.width * 35 / 100
	contentHeight := m.height - headerHeight - 1 // -1 for footer
	topHeight := contentHeight * 35 / 100

	x, y := msg.X, msg.Y

	// Determine which panel was clicked
	var targetPanel int
	if x < leftWidth {
		// Left side: Stats (top), Activity (bottom)
		if y > headerHeight && y < headerHeight+topHeight {
			targetPanel = PanelStats
		} else if y >= headerHeight+topHeight {
			targetPanel = PanelActivity
		} else {
			return m, nil
		}
	} else {
		// Right side: Projects (top), Sessions (bottom)
		if y > headerHeight && y < headerHeight+topHeight {
			targetPanel = PanelProjects
		} else if y >= headerHeight+topHeight {
			targetPanel = PanelSessions
		} else {
			return m, nil
		}
	}

	switch msg.Button {
	case tea.MouseButtonLeft:
		// Click to focus panel
		if msg.Action == tea.MouseActionPress {
			m.focusPanel(targetPanel)
		}
	case tea.MouseButtonWheelUp:
		// Scroll up in focused panel
		if targetPanel == m.focused {
			m.cursorUp()
		}
	case tea.MouseButtonWheelDown:
		// Scroll down in focused panel
		if targetPanel == m.focused {
			m.cursorDown()
		}
	}

	return m, nil
}

func (m *Model) updateFocusStates() {
	m.stats.SetFocused(m.focused == PanelStats)
	m.activity.SetFocused(m.focused == PanelActivity)
	m.projects.SetFocused(m.focused == PanelProjects)
	m.sessions.SetFocused(m.focused == PanelSessions)
}

func (m *Model) updateSizes() {
	// Header takes 1 line
	// Footer takes 1 line
	// Main content gets the rest
	headerHeight := 1
	footerHeight := 1
	contentHeight := m.height - headerHeight - footerHeight

	// Left panel is 35% width, right panel is 65%
	leftWidth := m.width * 35 / 100
	rightWidth := m.width - leftWidth

	// Top row: 35% height, Bottom row: 65% height
	topHeight := contentHeight * 35 / 100
	bottomHeight := contentHeight - topHeight

	m.header.SetWidth(m.width)
	m.stats.SetSize(leftWidth, topHeight)
	m.activity.SetSize(leftWidth, bottomHeight)
	m.projects.SetSize(rightWidth, topHeight)
	m.sessions.SetSize(rightWidth, bottomHeight)
	m.help.SetSize(m.width, m.height)
	m.detail.SetSize(m.width, m.height)
}

func (m *Model) updateWidgets() {
	if m.dashData == nil {
		return
	}

	m.header.Update(m.dashData.VMStatus, m.paused)
	m.stats.Update(m.dashData, m.timeRange)
	m.activity.Update(m.dashData, m.timeRange)

	// Apply time range filter to projects and sessions
	filteredProjects := m.dashData.FilterProjects(m.timeRange)
	filteredSessions := m.dashData.FilterSessions(m.timeRange)

	m.projects.Update(filteredProjects, m.timeRange)
	m.sessions.Update(filteredSessions, m.timeRange)
	m.updateFocusStates()
}

// View renders the application.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var lines []string

	// Header
	lines = append(lines, m.header.View())

	// Main content
	// Left panel is 35% width, right panel is 65%
	leftWidth := m.width * 35 / 100
	rightWidth := m.width - leftWidth

	// Calculate heights
	headerHeight := 1
	footerHeight := 1
	contentHeight := m.height - headerHeight - footerHeight

	// Left column: Stats on top, Activity on bottom
	leftTop := m.stats.View()
	leftBottom := m.activity.View()
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, leftTop, leftBottom)
	leftColumn = lipgloss.NewStyle().Width(leftWidth).Height(contentHeight).Render(leftColumn)

	// Right column: Projects on top, Sessions on bottom
	rightTop := m.projects.View()
	rightBottom := m.sessions.View()
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, rightTop, rightBottom)
	rightColumn = lipgloss.NewStyle().Width(rightWidth).Height(contentHeight).Render(rightColumn)

	// Join columns
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
	lines = append(lines, mainContent)

	// Footer
	footer := m.renderFooter()
	lines = append(lines, footer)

	view := strings.Join(lines, "\n")

	// Overlay help if visible
	if m.help.IsVisible() {
		helpView := m.help.View()
		// Overlay the help on top of the view
		// Simple overlay: just return help view centered
		return helpView
	}

	// Overlay detail modal if visible
	if m.detail.IsVisible() {
		return m.detail.View()
	}

	return view
}

func (m Model) renderFooter() string {
	// Check for flash message
	if m.flashMessage != "" && time.Now().Before(m.flashExpiry) {
		flashStyle := lipgloss.NewStyle().
			Background(SurfaceDark).
			Foreground(Success).
			Width(m.width).
			Padding(0, 1)
		return flashStyle.Render(m.flashMessage)
	}

	var parts []string

	// Get context-specific keybindings from the focused panel
	var panelBindings []Keybinding
	switch m.focused {
	case PanelStats:
		panelBindings = m.stats.GetKeybindings()
	case PanelActivity:
		panelBindings = m.activity.GetKeybindings()
	case PanelProjects:
		panelBindings = m.projects.GetKeybindings()
	case PanelSessions:
		panelBindings = m.sessions.GetKeybindings()
	}

	// Add panel-specific bindings first (highlighted - these are dynamic)
	panelKeyStyle := lipgloss.NewStyle().Foreground(Primary).Bold(true)
	panelDescStyle := lipgloss.NewStyle().Foreground(Text)
	for _, b := range panelBindings {
		key := panelKeyStyle.Render(b.Key)
		desc := panelDescStyle.Render(b.Desc)
		parts = append(parts, key+" "+desc)
	}

	// Global key hints (muted - these are always available)
	globalHints := []Keybinding{
		{"q", "quit"},
		{"r", "refresh"},
		{"p", "pause"},
		{"t", m.timeRange.String()},
		{"1-4", "jump"},
		{"?", "help"},
	}

	globalKeyStyle := lipgloss.NewStyle().Foreground(TextMuted)
	globalDescStyle := lipgloss.NewStyle().Foreground(Text)
	for _, h := range globalHints {
		key := globalKeyStyle.Render(h.Key)
		desc := globalDescStyle.Render(h.Desc)
		parts = append(parts, key+" "+desc)
	}

	leftContent := strings.Join(parts, "  ")

	// Right side: theme indicator (same style as global shortcuts)
	themeKey := globalKeyStyle.Render("T")
	themeName := globalDescStyle.Render(CurrentTheme)
	rightContent := themeKey + " " + themeName

	// Calculate spacing to push theme to the right
	leftLen := lipgloss.Width(leftContent)
	rightLen := lipgloss.Width(rightContent)
	padding := m.width - leftLen - rightLen - 4 // -4 for outer padding
	if padding < 2 {
		padding = 2
	}

	footer := leftContent + strings.Repeat(" ", padding) + rightContent

	footerStyle := lipgloss.NewStyle().
		Background(SurfaceDark).
		Width(m.width).
		Padding(0, 1)

	return footerStyle.Render(footer)
}
