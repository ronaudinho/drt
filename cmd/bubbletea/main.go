package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dotabuff/manta"
	"github.com/ronaudinho/drt/cmd/bubbletea/internal"
)

const (
	// TODO: too small width and height make the scale not accurate
	// but too big width and height make the terminal window too big
	width  = 20
	height = 20
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var keys = internal.KeyMap{
	Networth: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Networth"),
	),
	Kill: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "Kill"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "a"),
		key.WithHelp("<-/a", "seek backward by -10k tick"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "d"),
		key.WithHelp("->/d", "seek forward by +10k tick"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

type sessionState int

type drtModel struct {
	tickPositions                          map[uint32]map[string]pos
	XYPositions                            map[string]pos
	help                                   help.Model
	printUnitNames                         string
	messageToUser                          string
	temporaryMessageToDisplayTickPositions string
	keys                                   internal.KeyMap
	table                                  table.Model
	spinner                                spinner.Model
	state                                  sessionState
	counter                                int
	secondsElapsed                         int
	currentTick                            uint32
}

const (
	showLoading sessionState = iota
	showReplay
)

// TODO: replace this dummy data examples with the actual parsed DotA2 replay json file
// TODO: add keyMap to show/hide networth, kill, lasthit, bb status, etc by dynamically changing the table column
var originalRows = []table.Row{
	{"1", "Tokyo", "Japan", "37,274,000"},
	{"2", "Delhi", "India", "32,065,760"},
	{"3", "Shanghai", "China", "28,516,904"},
	{"4", "Dhaka", "Bangladesh", "22,478,116"},
	{"5", "São Paulo", "Brazil", "22,429,800"},
	{"6", "Mexico City", "Mexico", "22,085,140"},
	{"7", "Cairo", "Egypt", "21,750,020"},
	{"8", "Beijing", "China", "21,333,332"},
	{"9", "Mumbai", "India", "20,961,472"},
	{"10", "Osaka", "Japan", "19,059,856"},
	{"11", "Chongqing", "China", "16,874,740"},
	{"12", "Karachi", "Pakistan", "16,839,950"},
	{"13", "Istanbul", "Turkey", "15,636,243"},
	{"14", "Kinshasa", "DR Congo", "15,628,085"},
	{"15", "Lagos", "Nigeria", "15,387,639"},
	{"16", "Buenos Aires", "Argentina", "15,369,919"},
	{"17", "Kolkata", "India", "15,133,888"},
	{"18", "Manila", "Philippines", "14,406,059"},
	{"19", "Tianjin", "China", "14,011,828"},
	{"20", "Guangzhou", "China", "13,964,637"},
	{"21", "Rio De Janeiro", "Brazil", "13,634,274"},
	{"22", "Lahore", "Pakistan", "13,541,764"},
	{"23", "Bangalore", "India", "13,193,035"},
	{"24", "Shenzhen", "China", "12,831,330"},
	{"25", "Moscow", "Russia", "12,640,818"},
	{"26", "Chennai", "India", "11,503,293"},
	{"27", "Bogota", "Colombia", "11,344,312"},
	{"28", "Paris", "France", "11,142,303"},
	{"29", "Jakarta", "Indonesia", "11,074,811"},
	{"30", "Lima", "Peru", "11,044,607"},
	{"31", "Bangkok", "Thailand", "10,899,698"},
	{"32", "Hyderabad", "India", "10,534,418"},
	{"33", "Seoul", "South Korea", "9,975,709"},
	{"34", "Nagoya", "Japan", "9,571,596"},
	{"35", "London", "United Kingdom", "9,540,576"},
	{"36", "Chengdu", "China", "9,478,521"},
	{"37", "Nanjing", "China", "9,429,381"},
	{"38", "Tehran", "Iran", "9,381,546"},
	{"39", "Ho Chi Minh City", "Vietnam", "9,077,158"},
	{"40", "Luanda", "Angola", "8,952,496"},
	{"41", "Wuhan", "China", "8,591,611"},
	{"42", "Xi An Shaanxi", "China", "8,537,646"},
	{"43", "Ahmedabad", "India", "8,450,228"},
	{"44", "Kuala Lumpur", "Malaysia", "8,419,566"},
	{"45", "New York City", "United States", "8,177,020"},
	{"46", "Hangzhou", "China", "8,044,878"},
	{"47", "Surat", "India", "7,784,276"},
	{"48", "Suzhou", "China", "7,764,499"},
	{"49", "Hong Kong", "Hong Kong", "7,643,256"},
	{"50", "Riyadh", "Saudi Arabia", "7,538,200"},
	{"51", "Shenyang", "China", "7,527,975"},
	{"52", "Baghdad", "Iraq", "7,511,920"},
	{"53", "Dongguan", "China", "7,511,851"},
	{"54", "Foshan", "China", "7,497,263"},
	{"55", "Dar Es Salaam", "Tanzania", "7,404,689"},
	{"56", "Pune", "India", "6,987,077"},
	{"57", "Santiago", "Chile", "6,856,939"},
	{"58", "Madrid", "Spain", "6,713,557"},
	{"59", "Haerbin", "China", "6,665,951"},
	{"60", "Toronto", "Canada", "6,312,974"},
	{"61", "Belo Horizonte", "Brazil", "6,194,292"},
	{"62", "Khartoum", "Sudan", "6,160,327"},
	{"63", "Johannesburg", "South Africa", "6,065,354"},
	{"64", "Singapore", "Singapore", "6,039,577"},
	{"65", "Dalian", "China", "5,930,140"},
	{"66", "Qingdao", "China", "5,865,232"},
	{"67", "Zhengzhou", "China", "5,690,312"},
	{"68", "Ji Nan Shandong", "China", "5,663,015"},
	{"69", "Barcelona", "Spain", "5,658,472"},
	{"70", "Saint Petersburg", "Russia", "5,535,556"},
	{"71", "Abidjan", "Ivory Coast", "5,515,790"},
	{"72", "Yangon", "Myanmar", "5,514,454"},
	{"73", "Fukuoka", "Japan", "5,502,591"},
	{"74", "Alexandria", "Egypt", "5,483,605"},
	{"75", "Guadalajara", "Mexico", "5,339,583"},
	{"76", "Ankara", "Turkey", "5,309,690"},
	{"77", "Chittagong", "Bangladesh", "5,252,842"},
	{"78", "Addis Ababa", "Ethiopia", "5,227,794"},
	{"79", "Melbourne", "Australia", "5,150,766"},
	{"80", "Nairobi", "Kenya", "5,118,844"},
	{"81", "Hanoi", "Vietnam", "5,067,352"},
	{"82", "Sydney", "Australia", "5,056,571"},
	{"83", "Monterrey", "Mexico", "5,036,535"},
	{"84", "Changsha", "China", "4,809,887"},
	{"85", "Brasilia", "Brazil", "4,803,877"},
	{"86", "Cape Town", "South Africa", "4,800,954"},
	{"87", "Jiddah", "Saudi Arabia", "4,780,740"},
	{"88", "Urumqi", "China", "4,710,203"},
	{"89", "Kunming", "China", "4,657,381"},
	{"90", "Changchun", "China", "4,616,002"},
	{"91", "Hefei", "China", "4,496,456"},
	{"92", "Shantou", "China", "4,490,411"},
	{"93", "Xinbei", "Taiwan", "4,470,672"},
	{"94", "Kabul", "Afghanistan", "4,457,882"},
	{"95", "Ningbo", "China", "4,405,292"},
	{"96", "Tel Aviv", "Israel", "4,343,584"},
	{"97", "Yaounde", "Cameroon", "4,336,670"},
	{"98", "Rome", "Italy", "4,297,877"},
	{"99", "Shijiazhuang", "China", "4,285,135"},
	{"100", "Montreal", "Canada", "4,276,526"},
}

func initTable(columnName string) table.Model {
	columns := []table.Column{}
	var rows []table.Row

	// Default column is city and country
	columns = append(columns, table.Column{Title: "City", Width: 10})
	for _, row := range originalRows {
		rows = append(rows, table.Row{row[1]})
	}

	columns = append(columns, table.Column{Title: "Country", Width: 10})
	for i, row := range originalRows {
		rows[i] = append(rows[i], row[2])
	}

	// Toggle columns
	switch columnName {
	// TODO: replace hardcoded columnName with enum
	case "Population":
		columns = append(columns, table.Column{Title: "Population", Width: 10})
		for i, row := range originalRows {
			rows[i] = append(rows[i], row[3])
		}
	case "Rank":
		columns = append(columns, table.Column{Title: "Rank", Width: 4})
		for i, row := range originalRows {
			rows[i] = append(rows[i], row[0])
		}
	default:
		return table.New()
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func newModel() drtModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return drtModel{
		keys:          keys,
		help:          help.New(),
		tickPositions: map[uint32]map[string]pos{},
		// Rank or Population
		table:       initTable("Rank"),
		XYPositions: map[string]pos{},
		spinner:     s,
	}
}

// TODO: this currently only holds x-y pos but we can
// incorporate other units related data into this
type pos struct {
	CX, CY, VX, VY uint32
}

// type that returned for tea.Cmd
type (
	tickMsg           struct{}
	loadReplayDataMsg struct{ mapPositions map[uint32]map[string]pos }
)

func tick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return tickMsg{}
	}
}

func loadReplayData() tea.Cmd {
	return func() tea.Msg {
		mapPositions, err := parse("7569667371")
		if err != nil {
			log.Fatal(err)
		}
		return loadReplayDataMsg{mapPositions}
	}
}

// Init implements tea.Model.
func (m drtModel) Init() tea.Cmd {
	// TODO: add features to load replay data from list of downloaded replays
	// TODO: maybe integrate cmd/drt feature to download replay and list downloaded replay
	return loadReplayData()
}

// Update implements tea.Model.
func (m drtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case loadReplayDataMsg:
		m.tickPositions = msg.mapPositions

		if len(m.tickPositions) != 0 {
			m.state = showReplay
		}
		// TODO: replace tick with minutes and seconds like in DotA2 (h:m), maybe add night and day cycle
		return m, tick()
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
		// keys is ticks that sorted in ascending order
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
			m.currentTick = k

			units := m.tickPositions[k]
			xmax, ymax := width, height
			m.XYPositions = scaleCXCY(units, xmax, ymax)

			tickCounter++
			if tickCounter >= tickInSeconds {
				tickCounter = 0
				break
			}
		}
		return m, tick()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Networth):
			m.messageToUser = "Network toggled"
			newTable := initTable("Population")
			m.table = newTable
		case key.Matches(msg, m.keys.Kill):
			m.messageToUser = "Kill toggled"
			newTable := initTable("Rank")
			m.table = newTable
		case key.Matches(msg, m.keys.Left):
			m.messageToUser = "a pressed"
			m.currentTick -= 10000
		case key.Matches(msg, m.keys.Right):
			m.messageToUser = "d pressed"
			m.currentTick += 10000
		case key.Matches(msg, m.keys.Quit):
			m.messageToUser = "q pressed"
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		// TODO: use qwer, etc shortcuts to show networth, kill, lasthit, bb status, etc
		default:
			m.messageToUser = "Key pressed: " + msg.String()
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// TODO: code from rsc, maybe use the code directly from rsc as a package
func scaleCXCY(units map[string]pos, xmax int, ymax int) map[string]pos {
	scaledUnits := make(map[string]pos)

	for name, unitPos := range units {
		xp := float32((int(unitPos.CX) - 68)) / float32(120)
		yp := float32((int(unitPos.CY) - 68)) / float32(120)
		xx := int(xp * float32(xmax))
		yy := int((1 - yp) * float32(ymax))
		scaledUnits[name] = pos{
			CX: uint32(xx),
			CY: uint32(yy),
		}
	}

	return scaledUnits
}

func (m drtModel) replayView(s string) string {
	// TODO: parse the timer into minutes and seconds like in DotA2
	s += fmt.Sprintf("\t\t\t\t\t%v seconds elapsed\n", m.secondsElapsed)

	// Main Content
	for x := 0; x < width; x++ {
		// BUG: this append on s for debugging to show row numbers, but a critical bug found when this commented
		// bug unitName displayed multiple times
		s += strconv.Itoa(x + 11)
		// using empty string didn't fix the bug
		// s += "  "
		for y := 0; y < height; y++ {
			sb := strings.Builder{}

			sortedUnitName := make([]string, 0)
			for k := range m.XYPositions {
				sortedUnitName = append(sortedUnitName, k)
			}
			sort.Strings(sortedUnitName)

			// TODO: find a better way to drawn multiple units in the same position instead just ignoring it
			anUnitAlreadyDrawnThere := false
			sb.WriteString("\n")

			for i, unitName := range sortedUnitName {
				// TODO: find a way to groups the unit by team instead of sorted by unitName
				sb.WriteString(unitName[:3] + ": " + unitName + "\t")
				if i == 4 {
					sb.WriteString("\n")
				}

				if anUnitAlreadyDrawnThere {
					continue
				}
				if m.XYPositions[unitName].CX == uint32(x) && m.XYPositions[unitName].CY == uint32(y) {
					s += unitName[:3]
					anUnitAlreadyDrawnThere = true
				}
			}

			if !anUnitAlreadyDrawnThere {
				s += "---"
			}
			m.printUnitNames = sb.String()
		}
		s += "\n"
	}

	s += fmt.Sprintf("tick: %s", strconv.Itoa(int(m.currentTick)))
	s += fmt.Sprintf("%s\n", m.temporaryMessageToDisplayTickPositions)
	s += fmt.Sprintf("%s\n", m.printUnitNames)

	s += "\n\n"
	s += m.table.View()
	return s
}

// View implements tea.Model.
func (m drtModel) View() string {
	// Header
	s := "drt v0.0.1\n\n"

	switch m.state {
	case showLoading:
		s += fmt.Sprintf("\n\n   %s Loading forever...press \tesc/ctrl+c\t to quit\n\n", m.spinner.View())
	case showReplay:
		s = m.replayView(s)
	}
	s += m.messageToUser

	// Footer
	s += "\n\n"
	s += m.help.View(m.keys)
	return baseStyle.Render(s) + "\n"
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

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
