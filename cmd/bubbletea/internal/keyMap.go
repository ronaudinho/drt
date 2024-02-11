package internal

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
}
