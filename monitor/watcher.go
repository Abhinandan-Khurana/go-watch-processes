package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// --- Ignore List ---
// A list of substrings to filter out common, noisy background processes.
// The check is case-insensitive.
var ignoreList = []string{
	"Google Chrome Helper",
	"mdworker_shared",
	"mdworker",
	"Code Helper", // For VS Code
	"gopls",       // Go Language Server
	"kernel_task",
	"WindowServer",
	"loginwindow",
	"launchd",
	"distnoted",
}

// ProcessEvent holds all relevant data for a monitored process.
type ProcessEvent struct {
	PID       int32
	Name      string
	StartTime time.Time
	CPUHigh   bool // Flag to ensure we only notify once for high CPU
	cpuChecks int  // Counter for sustained high CPU usage
}

// AdvancedDetectionEvent represents a special alert.
type AdvancedDetectionEvent struct {
	PID     int32
	Name    string
	Reason  string // e.g., "High CPU Usage", "Suspicious Name"
	Details string // e.g., "CPU at 85%", "Matched keyword: miner"
}

const (
	cpuAlertThreshold    = 80.0 // Percent
	cpuSustainedDuration = 5    // Seconds
)

var suspiciousKeywords = []string{"miner", "keylog", "tor", "rat", "malware", "backdoor", "xmrig", "mimikatz", "lazagne", "cobaltstrike", "meterpreter", "remcos", "quasar", "darkcomet", "nanocore", "plugx", "gh0st", "njrat", "redline", "vidar", "azorult", "agenttesla", "lokibot", "formbook", "lockbit", "revil"}

// isIgnored checks if a process name should be ignored based on the ignoreList.
func isIgnored(processName string) bool {
	lowerCaseName := strings.ToLower(processName)
	for _, ignored := range ignoreList {
		if strings.Contains(lowerCaseName, strings.ToLower(ignored)) {
			return true
		}
	}
	return false
}

// Start spawns the new advanced process watcher.
func Start(
	longRunningEvents chan<- *ProcessEvent,
	advancedEvents chan<- *AdvancedDetectionEvent,
) {
	// Internal map to track running processes
	trackedProcs := make(map[int32]*ProcessEvent)

	// Initialize with currently running processes to avoid a flood of alerts on startup
	initialProcs, _ := process.Processes()
	for _, p := range initialProcs {
		trackedProcs[p.Pid] = &ProcessEvent{PID: p.Pid} // Just mark them as seen
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentProcs, err := process.Processes()
		if err != nil {
			continue
		}

		currentProcsMap := make(map[int32]*process.Process)
		for _, p := range currentProcs {
			currentProcsMap[p.Pid] = p
		}

		// --- Check for NEW processes and perform advanced checks ---
		for pid, p := range currentProcsMap {
			if _, exists := trackedProcs[pid]; !exists {
				// This is a new process
				name, _ := p.Name()

				// --- IGNORE CHECK ---
				if isIgnored(name) {
					continue // Skip this process entirely
				}

				startTime, _ := p.CreateTime()

				newEvent := &ProcessEvent{
					PID:       pid,
					Name:      name,
					StartTime: time.Unix(0, startTime*int64(time.Millisecond)),
				}
				trackedProcs[pid] = newEvent

				// Advanced Check 1: Suspicious Name
				for _, keyword := range suspiciousKeywords {
					if strings.Contains(strings.ToLower(name), keyword) {
						advancedEvents <- &AdvancedDetectionEvent{
							PID:     pid,
							Name:    name,
							Reason:  "Suspicious Process Detected",
							Details: "Matched keyword: '" + keyword + "'",
						}
						break // Only one notification per process
					}
				}
			} else {
				// This is an existing process, perform ongoing checks
				trackedProc := trackedProcs[pid]
				if trackedProc == nil || trackedProc.CPUHigh || trackedProc.StartTime.IsZero() {
					// Don't monitor ignored (StartTime is zero) or already alerted processes
					continue
				}

				// Advanced Check 2: High CPU Usage
				cpuPercent, err := p.CPUPercent()
				if err != nil {
					continue
				}

				if cpuPercent > cpuAlertThreshold {
					trackedProc.cpuChecks++
				} else {
					trackedProc.cpuChecks = 0 // Reset if usage drops
				}

				if trackedProc.cpuChecks >= cpuSustainedDuration {
					trackedProc.CPUHigh = true // Mark it so we don't alert again
					advancedEvents <- &AdvancedDetectionEvent{
						PID:     pid,
						Name:    trackedProc.Name,
						Reason:  "Sustained High CPU Usage",
						Details: fmt.Sprintf("CPU at %.2f%% for %d seconds", cpuPercent, cpuSustainedDuration),
					}
				}
			}
		}

		// --- Check for COMPLETED processes ---
		for pid, trackedProc := range trackedProcs {
			if _, exists := currentProcsMap[pid]; !exists {
				// Process has terminated
				if trackedProc != nil && !trackedProc.StartTime.IsZero() {
					duration := time.Since(trackedProc.StartTime)
					if duration.Seconds() > 10 {
						// This is a long-running process that just finished
						longRunningEvents <- trackedProc
					}
				}
				delete(trackedProcs, pid) // Clean up
			}
		}
	}
}
