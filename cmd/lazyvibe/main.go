// Package main is the entry point for lazyvibe CLI.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moshe-exe/lazyvibe/internal/config"
	"github.com/moshe-exe/lazyvibe/internal/data"
	"github.com/moshe-exe/lazyvibe/internal/ui"
)

func main() {
	// Load configuration
	cfg, _ := config.Load()
	if cfg != nil && cfg.Theme != "" {
		ui.ApplyTheme(cfg.Theme)
	}

	// CLI flags
	dump := flag.Bool("dump", false, "Dump all dashboard data as JSON and exit")
	capture := flag.String("capture", "", "Capture visual output as ASCII text at specified terminal size (e.g., 120x40) and exit")
	flag.Parse()

	if *dump {
		dumpData()
		return
	}

	if *capture != "" {
		runCapture(*capture)
		return
	}

	// Normal TUI mode
	runTUI()
}

func dumpData() {
	manager := data.NewManager()
	dashData := manager.GetDashboardData(false)

	// Convert to JSON-friendly structure
	output := map[string]interface{}{
		"vm_status": map[string]interface{}{
			"running":     dashData.VMStatus.Running,
			"pid":         dashData.VMStatus.PID,
			"cpu_percent": dashData.VMStatus.CPUPercent,
			"memory_mb":   dashData.VMStatus.MemoryMB,
		},
		"sessions":       dashData.Sessions,
		"daily_activity": dashData.DailyActivity,
		"projects":       dashData.Projects,
		"totals": map[string]int{
			"sessions":   dashData.TotalSessions(),
			"messages":   dashData.TotalMessages(),
			"tool_calls": dashData.TotalToolCalls(),
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func runCapture(sizeStr string) {
	width, height, err := parseSize(sizeStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	manager := data.NewManager()
	model := ui.NewModel(manager)

	// Simulate window size and data load
	dashData := manager.GetDashboardData(false)

	// Create a new model with the size
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: width, Height: height})
	model = newModel.(ui.Model)

	// Update with data
	newModel, _ = model.Update(dashData)
	model = newModel.(ui.Model)

	// Render the view
	fmt.Println(model.View())
}

func parseSize(sizeStr string) (int, int, error) {
	parts := strings.Split(strings.ToLower(sizeStr), "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid size format '%s'. Use WxH (e.g., 120x40)", sizeStr)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width '%s'", parts[0])
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height '%s'", parts[1])
	}

	if width < 40 || height < 10 {
		return 0, 0, fmt.Errorf("minimum size is 40x10")
	}

	return width, height, nil
}

func runTUI() {
	manager := data.NewManager()
	model := ui.NewModel(manager)

	p := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
