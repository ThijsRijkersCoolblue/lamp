package ansi

import (
	"lamp/params"
	"lamp/sgr"

	"github.com/gdamore/tcell/v2"
)

func handleCSI(screen tcell.Screen, cmd byte, parameters string, cursorX, cursorY *int, style *tcell.Style, width, height int, state *State) {
	args := params.ParseParams(parameters)

	switch cmd {
	case 'A':
		n := params.GetArg(args, 0, 1)
		*cursorY -= n
		if *cursorY < 0 {
			*cursorY = 0
		}
	case 'B':
		n := params.GetArg(args, 0, 1)
		*cursorY += n
		if *cursorY >= height {
			*cursorY = height - 1
		}
	case 'C':
		n := params.GetArg(args, 0, 1)
		*cursorX += n
		if *cursorX >= width {
			*cursorX = width - 1
		}
	case 'D':
		n := params.GetArg(args, 0, 1)
		*cursorX -= n
		if *cursorX < 0 {
			*cursorX = 0
		}
	case 'E':
		n := params.GetArg(args, 0, 1)
		*cursorY += n
		*cursorX = 0
		if *cursorY >= height {
			*cursorY = height - 1
		}
	case 'F':
		n := params.GetArg(args, 0, 1)
		*cursorY -= n
		*cursorX = 0
		if *cursorY < 0 {
			*cursorY = 0
		}
	case 'G':
		col := params.GetArg(args, 0, 1) - 1
		if col < 0 {
			col = 0
		}
		if col >= width {
			col = width - 1
		}
		*cursorX = col
	case 'H', 'f':
		row := params.GetArg(args, 0, 1) - 1
		col := params.GetArg(args, 1, 1) - 1
		if row < 0 {
			row = 0
		}
		if col < 0 {
			col = 0
		}
		if row >= height {
			row = height - 1
		}
		if col >= width {
			col = width - 1
		}
		*cursorY = row
		*cursorX = col
	case 'J':
		n := params.GetArg(args, 0, 0)
		switch n {
		case 0:
			for x := *cursorX; x < width; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
			for y := *cursorY + 1; y < height; y++ {
				for x := 0; x < width; x++ {
					screen.SetContent(x, y, ' ', nil, *style)
				}
			}
		case 1:
			for y := 0; y < *cursorY; y++ {
				for x := 0; x < width; x++ {
					screen.SetContent(x, y, ' ', nil, *style)
				}
			}
			for x := 0; x <= *cursorX; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		case 2, 3:
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					screen.SetContent(x, y, ' ', nil, *style)
				}
			}
		}
	case 'K':
		n := params.GetArg(args, 0, 0)
		switch n {
		case 0:
			for x := *cursorX; x < width; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		case 1:
			for x := 0; x <= *cursorX; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		case 2:
			for x := 0; x < width; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		}
	case 'X':
		// erase n characters from cursor position
		n := params.GetArg(args, 0, 1)
		for x := *cursorX; x < *cursorX+n && x < width; x++ {
			screen.SetContent(x, *cursorY, ' ', nil, *style)
		}
	case '@':
		// insert n blank characters
		n := params.GetArg(args, 0, 1)
		for x := width - 1; x >= *cursorX+n; x-- {
			ch, comb, st, _ := screen.GetContent(x-n, *cursorY)
			screen.SetContent(x, *cursorY, ch, comb, st)
		}
		for x := *cursorX; x < *cursorX+n && x < width; x++ {
			screen.SetContent(x, *cursorY, ' ', nil, *style)
		}
	case 'L':
		n := params.GetArg(args, 0, 1)
		for i := 0; i < n; i++ {
			for y := state.ScrollBottom; y > *cursorY; y-- {
				for x := 0; x < width; x++ {
					ch, comb, st, _ := screen.GetContent(x, y-1)
					screen.SetContent(x, y, ch, comb, st)
				}
			}
			for x := 0; x < width; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		}
	case 'M':
		n := params.GetArg(args, 0, 1)
		for i := 0; i < n; i++ {
			for y := *cursorY; y < state.ScrollBottom; y++ {
				for x := 0; x < width; x++ {
					ch, comb, st, _ := screen.GetContent(x, y+1)
					screen.SetContent(x, y, ch, comb, st)
				}
			}
			for x := 0; x < width; x++ {
				screen.SetContent(x, state.ScrollBottom, ' ', nil, *style)
			}
		}
	case 'P':
		n := params.GetArg(args, 0, 1)
		for x := *cursorX; x < width-n; x++ {
			ch, comb, st, _ := screen.GetContent(x+n, *cursorY)
			screen.SetContent(x, *cursorY, ch, comb, st)
		}
		for x := width - n; x < width; x++ {
			screen.SetContent(x, *cursorY, ' ', nil, *style)
		}
	case 'S':
		n := params.GetArg(args, 0, 1)
		for i := 0; i < n; i++ {
			scrollUp(screen, *style, width, state.ScrollTop, state.ScrollBottom)
		}
	case 'T':
		n := params.GetArg(args, 0, 1)
		for i := 0; i < n; i++ {
			scrollDown(screen, *style, width, state.ScrollTop, state.ScrollBottom)
		}
	case 'd':
		row := params.GetArg(args, 0, 1) - 1
		if row < 0 {
			row = 0
		}
		if row >= height {
			row = height - 1
		}
		*cursorY = row
	case 'r':
		top := params.GetArg(args, 0, 1) - 1
		bot := params.GetArg(args, 1, height) - 1
		if top < 0 {
			top = 0
		}
		if bot >= height {
			bot = height - 1
		}
		if top < bot {
			state.ScrollTop = top
			state.ScrollBottom = bot
		}
		*cursorX = 0
		*cursorY = 0
	case 'm':
		sgr.HandleSGR(args, style)
	}
}
