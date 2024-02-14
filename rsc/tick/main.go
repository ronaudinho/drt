package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dotabuff/manta"
	"github.com/gdamore/tcell/v2"
)

// TODO: this currently only holds x-y pos but we can
// incorporate other units related data into this
type pos struct {
	CX, CY, VX, VY uint32
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("id not supplied")
	}
	id := os.Args[1]
	// TODO: inmem map is not the best way to approach this
	m, err := parse(id)
	if err != nil {
		log.Fatal(err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	xmax, ymax := s.Size()
	// TODO: make this seekable
	// this currently only shows position of hero* units
	// at whatever timestamp that is accessed
	for n, u := range m {
		drawText(s, 1, 1, 1, 10, boxStyle, fmt.Sprint(n))
		drawUnits(s, boxStyle, u)
		break
	}
	drawBox(s, 3, ymax-5, xmax-3, ymax-3, boxStyle, "drag to seek")

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	ox := -1
	for {
		s.Show()
		// TODO: PollEvent is probably not the best option
		// as we will have events coming in without input
		// when playing the replay
		// we can of course use s.PostEvent() but the amount
		// of events we can store is limited.
		// Consider using s.ChannelEvent() perhaps.
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
		case *tcell.EventMouse:
			// TODO: this does nothing but draw box
			x, y := ev.Position()
			switch ev.Buttons() {
			case tcell.Button1:
				if ox < 0 && ymax-5 < y && y < ymax-3 {
					ox = x
				}
			case tcell.ButtonNone:
				if ox >= 0 && ymax-5 < y && y < ymax-3 {
					label := fmt.Sprintf("%d /%d", x, xmax)
					drawBox(s, x-1, y-1, x, y, boxStyle, label)
					ox = -1
				}
			}
		}
	}
}

func drawUnits(s tcell.Screen, style tcell.Style, units map[string]pos) {
	xmax, ymax := s.Size()
	y := 2
	// TODO: xmax and ymax for cell, tentatively using 68 and 188
	for name, pos := range units {
		xp := float32((int(pos.CX) - 68)) / float32(120)
		yp := float32((int(pos.CY) - 68)) / float32(120)
		xx := int(xp * float32(xmax))
		yy := int((1 - yp) * float32(ymax))
		drawText(s, xx, yy, xx+1, yy+1, style, name)
		// drawText(s, 1, y, 100, y+1, style, fmt.Sprintf("%s %d %d %5.2f %5.2f %5.2f %5.2f", name, pos.CX, pos.CY, xp, yp, xx, yy))
		y++
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
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
	p.OnEntity(func(e *manta.Entity, op manta.EntityOp) error {
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
