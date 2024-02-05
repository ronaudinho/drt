package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type drtModel struct {
	messageToUser string
	counter       int
}

// Init implements tea.Model.
func (drtModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m drtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.messageToUser = "q pressed"
			return m, tea.Quit
		default:
			m.messageToUser = "Key pressed: " + msg.String()
		}
	}

	return m, nil
}

// View implements tea.Model.
func (m drtModel) View() string {
	// Header
	s := "drt v0.0.1\n\n"
	// Main Content
	s += fmt.Sprintf("%s\n\n", m.messageToUser)
	// Footer
	s += "Press 'q' to quit"
	return s
}

func main() {
	p := tea.NewProgram(drtModel{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
