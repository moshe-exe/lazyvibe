package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/moshe-exe/lazyvibe/internal/data"
)

// HeaderModel represents the status header component.
type HeaderModel struct {
	vmStatus data.VMStatus
	paused   bool
	width    int
}

// NewHeaderModel creates a new header model.
func NewHeaderModel() HeaderModel {
	return HeaderModel{}
}

// Update updates the header with new VM status and pause state.
func (h *HeaderModel) Update(vmStatus data.VMStatus, paused bool) {
	h.vmStatus = vmStatus
	h.paused = paused
}

// SetWidth sets the header width.
func (h *HeaderModel) SetWidth(width int) {
	h.width = width
}

// View renders the header.
func (h HeaderModel) View() string {
	var parts []string

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)
	parts = append(parts, titleStyle.Render("lazyvibe"))

	// VM status
	if h.vmStatus.Running {
		vmInfo := fmt.Sprintf("VM: Running (PID %d)", *h.vmStatus.PID)
		parts = append(parts, SuccessStyle.Render(vmInfo))

		// CPU bar
		if h.vmStatus.CPUPercent != nil {
			cpuPercent := *h.vmStatus.CPUPercent
			cpuBar := RenderBar(cpuPercent, 10)
			cpuLabel := lipgloss.NewStyle().Foreground(TextMuted).Render("CPU ")
			cpuValue := fmt.Sprintf(" %.0f%%", cpuPercent)
			parts = append(parts, cpuLabel+"["+cpuBar+"]"+cpuValue)
		}

		// MEM bar
		if h.vmStatus.MemoryMB != nil {
			// Estimate percent based on typical VM allocation (assume 4GB max)
			memMB := *h.vmStatus.MemoryMB
			memPercent := memMB / 4096 * 100 // 4GB = 4096MB
			if memPercent > 100 {
				memPercent = 100
			}
			memBar := RenderBar(memPercent, 10)
			memLabel := lipgloss.NewStyle().Foreground(TextMuted).Render("MEM ")
			memValue := fmt.Sprintf(" %.0fMB", memMB)
			parts = append(parts, memLabel+"["+memBar+"]"+memValue)
		}
	} else {
		parts = append(parts, MutedStyle.Render("VM: Not Running"))
	}

	// Pause indicator
	if h.paused {
		pauseStyle := lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)
		parts = append(parts, pauseStyle.Render("[PAUSED]"))
	}

	content := strings.Join(parts, " | ")

	// Style the entire header
	style := HeaderStyle.Width(h.width)
	return style.Render(content)
}
