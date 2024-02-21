package internal

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Networth key.Binding
	Kill     key.Binding
	Left     key.Binding
	Right    key.Binding
	Quit     key.Binding
	Help     key.Binding
}

// FullHelp implements help.KeyMap.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Networth, k.Kill, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                      // second column
	}
}

// ShortHelp implements help.KeyMap.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
