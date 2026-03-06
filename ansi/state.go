package ansi

import "github.com/gdamore/tcell/v2"

type State struct {
	ScrollTop    int
	ScrollBottom int
	Style        tcell.Style
	Leftover     []byte
}

func NewState(height int) *State {
	return &State{
		ScrollTop:    0,
		ScrollBottom: height - 1,
		Style:        tcell.StyleDefault,
	}
}
