package data

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const vmProcessPattern = "com.apple.Virtualization.VirtualMachine"

// GetVMStatus gets the status of the Claude Desktop VM.
func GetVMStatus() VMStatus {
	// Find VM process using pgrep
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pgrep", "-f", vmProcessPattern)
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return VMStatus{Running: false}
	}

	// Get first PID
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return VMStatus{Running: false}
	}

	pid, err := strconv.Atoi(lines[0])
	if err != nil {
		return VMStatus{Running: false}
	}

	// Get CPU and memory stats using ps
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	psCmd := exec.CommandContext(ctx2, "ps", "-p", strconv.Itoa(pid), "-o", "%cpu=,rss=")
	psOutput, err := psCmd.Output()
	if err != nil {
		pidPtr := pid
		return VMStatus{Running: true, PID: &pidPtr}
	}

	parts := strings.Fields(strings.TrimSpace(string(psOutput)))
	if len(parts) >= 2 {
		cpuPercent, err1 := strconv.ParseFloat(parts[0], 64)
		memoryKB, err2 := strconv.Atoi(parts[1])

		pidPtr := pid
		status := VMStatus{Running: true, PID: &pidPtr}

		if err1 == nil {
			status.CPUPercent = &cpuPercent
		}
		if err2 == nil {
			memoryMB := float64(memoryKB) / 1024.0
			status.MemoryMB = &memoryMB
		}

		return status
	}

	pidPtr := pid
	return VMStatus{Running: true, PID: &pidPtr}
}
