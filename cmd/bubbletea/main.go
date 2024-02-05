package main

import (
	"fmt"
	"os"
	"sort"
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

// TODO: this currently only holds x-y pos but we can
// incorporate other units related data into this
type pos struct {
	CX, CY, VX, VY uint32
}

// TODO: replace this dummy data examples with the actual json file
var myUnits = map[string]pos{
	// H1, H2: spawn in bottom left radiant base
	// H3, H3: spawn in top right dire base
	// BUG: if CX and CY have same value, it should appear on top each other not draw another box
	// TODO: need more tests if CX, CY bigger than drawed map width and height
	"H1": {CX: 28, CY: 1, VX: 0, VY: 0},
	"H2": {CX: 29, CY: 2, VX: 0, VY: 0},
	"H3": {CX: 0, CY: 28, VX: 0, VY: 0},
	"H4": {CX: 1, CY: 29, VX: 0, VY: 0},
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
		// TODO: use qwer, etc shortcuts to show networth, kill, lasthit, bb status, etc
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
	// TODO: parse the timer into minutes and seconds like in DotA2
	s += fmt.Sprintf("\t\t%v seconds elapsed\n", m.secondsElapsed)

	// Main Content
	// TODO: draw map: wrap in function and add scaling for the width and height so the map can made bigger or smaller
	// XYPositions := drawUnits(myUnits)
	XYPositions := myUnits
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			keys := make([]string, 0)
			for k := range XYPositions {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				if XYPositions[k].CX == uint32(x) && XYPositions[k].CY == uint32(y) {
					// remove previous defautl block
					s = s[:len(s)-len("⬜️")]
					// TODO: replace yellow box with DotA2 hero icons
					s += "🟨"
				}
			}

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
