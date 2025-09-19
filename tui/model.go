package tui

import (
	"fmt"
	"time"

	"github.com/Abhinandan-Khurana/go-watch-processes/monitor"
	"github.com/Abhinandan-Khurana/go-watch-processes/notifications"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Bold   = "\033[1m"
)

const (
	toolName = `
	                ▐     ▌                                
	▞▀▌▞▀▖▄▄▖▌  ▌▝▀▖▜▀ ▞▀▖▛▀▖▄▄▖▛▀▖▙▀▖▞▀▖▞▀▖▞▀▖▞▀▘▞▀▘▞▀▖▞▀▘	 
	▚▄▌▌ ▌   ▐▐▐ ▞▀▌▐ ▖▌ ▖▌ ▌   ▙▄▘▌  ▌ ▌▌ ▖▛▀ ▝▀▖▝▀▖▛▀ ▝▀▖  
	▗▄▘▝▀     ▘▘ ▝▀▘ ▀ ▝▀ ▘ ▘   ▌  ▘  ▝▀ ▝▀ ▝▀▘▀▀ ▀▀ ▝▀▘▀▀   
																													 
	`
	authorName  = "Abhinandan Khurana"
	toolVersion = "v1.0.0"
)

// Define some styles for better visual output
var (
	styleCompleted = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))   // Green
	styleAlert     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))   // Red
	styleDetails   = lipgloss.NewStyle().Foreground(lipgloss.Color("244")) // Gray
)

// model represents the state of the TUI.
type model struct {
	longRunningEvents chan *monitor.ProcessEvent
	advancedEvents    chan *monitor.AdvancedDetectionEvent
	events            []string // A simple list to hold event strings for display
}

// InitialModel creates the initial model for the TUI.
func InitialModel() model {
	return model{
		longRunningEvents: make(chan *monitor.ProcessEvent),
		advancedEvents:    make(chan *monitor.AdvancedDetectionEvent),
		events:            make([]string, 0),
	}
}

// Init initializes the TUI.
func (m model) Init() tea.Cmd {
	go monitor.Start(m.longRunningEvents, m.advancedEvents)
	return tea.Batch(
		waitForLongRunningEvent(m.longRunningEvents),
		waitForAdvancedEvent(m.advancedEvents),
	)
}

// Update handles messages and updates the model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	// A long-running process has completed
	case *monitor.ProcessEvent:
		duration := time.Since(msg.StartTime).Round(time.Second)
		notifTitle := "Long-Running Task Completed"
		notifMsg := fmt.Sprintf("Name: %s (PID: %d)\nDuration: %s", msg.Name, msg.PID, duration)

		// Add to our display list
		logMsg := fmt.Sprintf("g %s %s",
			styleCompleted.Render("✔ Task Completed:"),
			msg.Name,
			styleDetails.Render(fmt.Sprintf("(PID: %d, Duration: g)", msg.PID, duration)),
		)
		m.addEvent(logMsg)

		notifications.Notify(notifTitle, notifMsg)
		return m, waitForLongRunningEvent(m.longRunningEvents)

	// An advanced detection event has occurred
	case *monitor.AdvancedDetectionEvent:
		notifTitle := msg.Reason
		notifMsg := fmt.Sprintf("Name: %s (PID: %d)\nDetails: %s", msg.Name, msg.PID, msg.Details)

		// Add to our display list
		logMsg := fmt.Sprintf("g %s %s",
			styleAlert.Render("❗ Alert: "+msg.Reason),
			msg.Name,
			styleDetails.Render(fmt.Sprintf("(PID: %d, g)", msg.PID, msg.Details)),
		)
		m.addEvent(logMsg)

		notifications.Notify(notifTitle, notifMsg)
		return m, waitForAdvancedEvent(m.advancedEvents)
	}
	return m, nil
}

// addEvent adds a new event to the list, keeping it from growing too large.
func (m *model) addEvent(event string) {
	m.events = append(m.events, fmt.Sprintf("[g] %s", time.Now().Format("15:04:05"), event))
	// Keep the list to a manageable size, e.g., max 20 entries
	if len(m.events) > 20 {
		m.events = m.events[1:]
	}
}

// View renders the TUI.
func (m model) View() string {
	s := banner()
	s += "\n\n"
	s += "Monitoring for long-running tasks and suspicious activity... Press 'q' to quit.\n\n"

	if len(m.events) == 0 {
		s += "No notable events detected yet."
	} else {
		s += lipgloss.JoinVertical(lipgloss.Left, m.events...)
	}

	return s
}

func banner() string {
	// (Banner can remain the same or be updated)
	return fmt.Sprintf(Blue+"%s \t \n"+Bold+Yellow+"\tBy: %s"+Reset+" | "+Red+"Version: %s"+Reset, toolName, authorName, toolVersion)
}

// --- Commands to wait for our custom events ---

func waitForLongRunningEvent(ch <-chan *monitor.ProcessEvent) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func waitForAdvancedEvent(ch <-chan *monitor.AdvancedDetectionEvent) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
