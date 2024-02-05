package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	width  = 30
	height = 30
)

type drtModel struct {
	messageToUser  string
	counter        int
	secondsElapsed int
}

type tickMsg struct{}

// Init implements tea.Model.
func (drtModel) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return tickMsg{}
	}
}

// Update implements tea.Model.
func (m drtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.secondsElapsed++
		return m, tick()
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
	s := "drt v0.0.1\n"
	// TODO: parse the timer into minutes and seconds like in DotA2
	s += fmt.Sprintf("%v seconds elapsed\n\n", m.secondsElapsed)
	// Main Content
	// TODO: draw map: wrap in function and add scaling for the width and height so the map can made bigger or smaller
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			s += "⬜️"
		}
		s += "\n"
	}
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
