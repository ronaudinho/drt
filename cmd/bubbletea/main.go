package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dotabuff/manta"
	"github.com/ronaudinho/drt/cmd/bubbletea/internal"
)

const (
	width  = 30
	height = 30
)

var keys = internal.KeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "a"),
		key.WithHelp("<-/a", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "d"),
		key.WithHelp("->/d", "move right"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type drtModel struct {
	tickPositions                          map[uint32]map[string]pos
	messageToUser                          string
	temporaryMessageToDisplayTickPositions string
	keys                                   internal.KeyMap
	counter                                int
	secondsElapsed                         int
	currentTick                            uint32
}

func newModel() drtModel {
	mapPositions, err := parse("7569667371")
	if err != nil {
		log.Fatal(err)
	}

	return drtModel{
		keys:          keys,
		tickPositions: mapPositions,
	}
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
func (m drtModel) Init() tea.Cmd {
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
		m.messageToUser = ""

		keys := make([]uint32, 0, len(m.tickPositions))
		for k := range m.tickPositions {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		// TODO: need trial and error to know how many ticks to show per second
		// TODO: separate tick for secondsElapsed and tickCounter
		// secondsElapsed have 1 seconds delay which make tickcounter got delayed too
		tickCounter := 0
		tickInSeconds := 50
		for _, k := range keys {
			// for current replay with ID 7569667371
			// the match started at tick 19k++, so tick below that will have static CX, CY
			if m.currentTick < 20000 {
				m.currentTick = k
				continue
			}
			if m.currentTick != 0 && m.currentTick >= k {
				continue
			}
			m.temporaryMessageToDisplayTickPositions = fmt.Sprintf("%v: %v", k, m.tickPositions[k])
			m.currentTick = k

			tickCounter++
			if tickCounter >= tickInSeconds {
				tickCounter = 0
				break
			}
		}
		return m, tick()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			m.messageToUser = "a pressed"
			m.currentTick -= 10000
		case key.Matches(msg, m.keys.Right):
			m.messageToUser = "d pressed"
			m.currentTick += 10000
		case key.Matches(msg, m.keys.Quit):
			m.messageToUser = "q pressed"
			return m, tea.Quit
		// TODO: use qwer, etc shortcuts to show networth, kill, lasthit, bb status, etc
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
	s += fmt.Sprintf("%s\n\n", m.temporaryMessageToDisplayTickPositions)
	s += fmt.Sprintf("%s\n\n", m.messageToUser)

	// Footer
	s += "Press 'q' to quit"
	s += "\nPress 'a' to seek backward by -10k tick"
	s += "\nPress 'd' to seek forward by +10k tick"
	return s
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}

func parse(id string) (map[uint32]map[string]pos, error) {
	f, err := os.Open(fmt.Sprintf("%s.dem", id))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p, err := manta.NewStreamParser(f)
	if err != nil {
		return nil, err
	}

	units := make(map[uint32]map[string]pos)
	p.OnEntity(func(e *manta.Entity, _ manta.EntityOp) error {
		c := e.GetClassName()
		if !strings.HasPrefix(c, "CDOTA_Unit_Hero_") {
			return nil
		}
		c = strings.TrimPrefix(c, "CDOTA_Unit_Hero_")
		// TODO: find a way to get timestamp instead
		// using tick does not seem like the best idea
		t := p.Tick
		cx, _ := e.GetUint32("CBodyComponent.m_cellX")
		cy, _ := e.GetUint32("CBodyComponent.m_cellY")
		// vx, _ := e.GetFloat32("CBodyComponent.m_vecX")
		// vy, _ := e.GetFloat32("CBodyComponent.m_vecY")

		if _, ok := units[t]; !ok {
			units[t] = make(map[string]pos)
		}
		units[t][c] = pos{
			CX: cx,
			CY: cy,
		}
		return nil
	})

	p.Start()

	for k, v := range units {
		if len(v) < 10 {
			delete(units, k)
		}
	}
	// TODO: json? really?
	b, _ := json.MarshalIndent(units, "", "  ")
	os.WriteFile(fmt.Sprintf("%s.json", id), b, 0666)
	return units, nil
}
