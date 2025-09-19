package main

import (
	"fmt"
	"os"

	"github.com/Abhinandan-Khurana/go-watch-processes/notifications"
	"github.com/Abhinandan-Khurana/go-watch-processes/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := notifications.Init(); err != nil {
		fmt.Printf("Error during notification setup: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tui.InitialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
