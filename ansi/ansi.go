package ansi

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

type State struct {
	ScrollTop    int
	ScrollBottom int
	Style        tcell.Style
}

func NewState(height int) *State {
	return &State{
		ScrollTop:    0,
		ScrollBottom: height - 1,
		Style:        tcell.StyleDefault,
	}
}

func ProcessOutput(screen tcell.Screen, data []byte, cursorX, cursorY *int, state *State) {
	width, height := screen.Size()

	// Update scroll bottom if height changed
	if state.ScrollBottom >= height {
		state.ScrollBottom = height - 1
	}

	i := 0
	for i < len(data) {
		ch := data[i]

		switch ch {
		case '\r':
			*cursorX = 0
		case '\n':
			*cursorX = 0
			*cursorY++
		case '\t':
			*cursorX += 4
		case 0x08, 0x7f: // Backspace
			if *cursorX > 0 {
				*cursorX--
				screen.SetContent(*cursorX, *cursorY, ' ', nil, state.Style)
			}
		case 0x1b: // ESC
			if i+1 < len(data) {
				switch data[i+1] {
				case '[': // CSI
					i += 2
					start := i
					for i < len(data) && (data[i] < '@' || data[i] > '~') {
						i++
					}
					if i < len(data) {
						cmd := data[i]
						params := string(data[start:i])

						// Alternate screen handling
						if cmd == 'h' && params == "?1049" {
							screen.Clear()
							*cursorX = 0
							*cursorY = 0
							state.ScrollTop = 0
							state.ScrollBottom = height - 1
						}
						if cmd == 'l' && params == "?1049" {
							screen.Clear()
							*cursorX = 0
							*cursorY = 0
							state.ScrollTop = 0
							state.ScrollBottom = height - 1
						}

						handleCSI(screen, cmd, params, cursorX, cursorY, &state.Style, width, height, state)
					}

				case ']': // OSC (title)
					i += 2
					for i < len(data) && data[i] != 0x07 {
						if data[i] == 0x1b && i+1 < len(data) && data[i+1] == '\\' {
							i++
							break
						}
						i++
					}

				case '(', ')', '*', '+', '-', '.', '/':
					i += 2
					if i < len(data) && data[i] == 'B' {
						i++
					}
					continue
				default:
					i++
				}
			}

		default:
			if ch >= 32 {
				r, size := utf8.DecodeRune(data[i:])
				x := *cursorX
				y := *cursorY
				if x >= width {
					x = width - 1
				}
				if y >= height {
					y = height - 1
				}
				screen.SetContent(x, y, r, nil, state.Style)
				*cursorX++
				i += size
				continue
			}
		}

		// Scroll if needed (only within scroll region)
		if *cursorY > state.ScrollBottom {
			scrollUp(screen, state.Style, width, state.ScrollTop, state.ScrollBottom)
			*cursorY = state.ScrollBottom
		}

		i++
	}

	screen.Show()
}

func scrollUp(screen tcell.Screen, style tcell.Style, width, scrollTop, scrollBottom int) {
	for y := scrollTop + 1; y <= scrollBottom; y++ {
		for x := 0; x < width; x++ {
			ch, comb, st, _ := screen.GetContent(x, y)
			screen.SetContent(x, y-1, ch, comb, st)
		}
	}
	for x := 0; x < width; x++ {
		screen.SetContent(x, scrollBottom, ' ', nil, style)
	}
}

func scrollDown(screen tcell.Screen, style tcell.Style, width, scrollTop, scrollBottom int) {
	for y := scrollBottom - 1; y >= scrollTop; y-- {
		for x := 0; x < width; x++ {
			ch, comb, st, _ := screen.GetContent(x, y)
			screen.SetContent(x, y+1, ch, comb, st)
		}
	}
	for x := 0; x < width; x++ {
		screen.SetContent(x, scrollTop, ' ', nil, style)
	}
}

func handleCSI(screen tcell.Screen, cmd byte, params string, cursorX, cursorY *int, style *tcell.Style, width, height int, state *State) {
	args := parseParams(params)

	switch cmd {
	case 'A': // Cursor Up
		n := getArg(args, 0, 1)
		*cursorY -= n
		if *cursorY < 0 {
			*cursorY = 0
		}
	case 'B': // Cursor Down
		n := getArg(args, 0, 1)
		*cursorY += n
		if *cursorY >= height {
			*cursorY = height - 1
		}
	case 'C': // Cursor Forward
		n := getArg(args, 0, 1)
		*cursorX += n
		if *cursorX >= width {
			*cursorX = width - 1
		}
	case 'D': // Cursor Back
		n := getArg(args, 0, 1)
		*cursorX -= n
		if *cursorX < 0 {
			*cursorX = 0
		}
	case 'E': // Cursor Next Line
		n := getArg(args, 0, 1)
		*cursorY += n
		*cursorX = 0
		if *cursorY >= height {
			*cursorY = height - 1
		}
	case 'F': // Cursor Previous Line
		n := getArg(args, 0, 1)
		*cursorY -= n
		*cursorX = 0
		if *cursorY < 0 {
			*cursorY = 0
		}
	case 'G': // Cursor Horizontal Absolute
		col := getArg(args, 0, 1) - 1
		if col < 0 {
			col = 0
		}
		if col >= width {
			col = width - 1
		}
		*cursorX = col
	case 'H', 'f': // Cursor Position
		row := getArg(args, 0, 1) - 1
		col := getArg(args, 1, 1) - 1
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

	case 'J': // Erase in Display
		n := getArg(args, 0, 0)
		switch n {
		case 0: // Erase from cursor to end of screen
			for x := *cursorX; x < width; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
			for y := *cursorY + 1; y < height; y++ {
				for x := 0; x < width; x++ {
					screen.SetContent(x, y, ' ', nil, *style)
				}
			}
		case 1: // Erase from start of screen to cursor
			for y := 0; y < *cursorY; y++ {
				for x := 0; x < width; x++ {
					screen.SetContent(x, y, ' ', nil, *style)
				}
			}
			for x := 0; x <= *cursorX; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		case 2, 3: // Erase entire screen — do NOT reset cursor
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					screen.SetContent(x, y, ' ', nil, *style)
				}
			}
		}

	case 'K': // Erase in Line — do NOT reset cursorX
		n := getArg(args, 0, 0)
		switch n {
		case 0: // Erase from cursor to end of line
			for x := *cursorX; x < width; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		case 1: // Erase from start of line to cursor
			for x := 0; x <= *cursorX; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		case 2: // Erase entire line
			for x := 0; x < width; x++ {
				screen.SetContent(x, *cursorY, ' ', nil, *style)
			}
		}

	case 'L': // Insert Lines (shift lines down within scroll region)
		n := getArg(args, 0, 1)
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

	case 'M': // Delete Lines (shift lines up within scroll region)
		n := getArg(args, 0, 1)
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

	case 'P': // Delete Characters
		n := getArg(args, 0, 1)
		for x := *cursorX; x < width-n; x++ {
			ch, comb, st, _ := screen.GetContent(x+n, *cursorY)
			screen.SetContent(x, *cursorY, ch, comb, st)
		}
		for x := width - n; x < width; x++ {
			screen.SetContent(x, *cursorY, ' ', nil, *style)
		}

	case 'S': // Scroll Up
		n := getArg(args, 0, 1)
		for i := 0; i < n; i++ {
			scrollUp(screen, *style, width, state.ScrollTop, state.ScrollBottom)
		}

	case 'T': // Scroll Down
		n := getArg(args, 0, 1)
		for i := 0; i < n; i++ {
			scrollDown(screen, *style, width, state.ScrollTop, state.ScrollBottom)
		}

	case 'd': // Line Position Absolute
		row := getArg(args, 0, 1) - 1
		if row < 0 {
			row = 0
		}
		if row >= height {
			row = height - 1
		}
		*cursorY = row

	case 'r': // Set Scroll Region
		top := getArg(args, 0, 1) - 1
		bot := getArg(args, 1, height) - 1
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
		// Cursor moves to top-left on scroll region change
		*cursorX = 0
		*cursorY = 0

	case 'm': // SGR - colors and styles
		if len(args) == 0 {
			*style = tcell.StyleDefault
			return
		}

		for i := 0; i < len(args); i++ {
			a := args[i]

			switch {
			case a == 0:
				*style = tcell.StyleDefault
			case a == 1:
				*style = (*style).Bold(true)
			case a == 2:
				*style = (*style).Dim(true)
			case a == 3:
				*style = (*style).Italic(true)
			case a == 4:
				*style = (*style).Underline(true)
			case a == 5:
				*style = (*style).Blink(true)
			case a == 7:
				*style = (*style).Reverse(true)
			case a == 22:
				*style = (*style).Bold(false).Dim(false)
			case a == 23:
				*style = (*style).Italic(false)
			case a == 24:
				*style = (*style).Underline(false)
			case a == 25:
				*style = (*style).Blink(false)
			case a == 27:
				*style = (*style).Reverse(false)

			// basic fg/bg
			case 30 <= a && a <= 37:
				*style = (*style).Foreground(ansiColor(a - 30))
			case 40 <= a && a <= 47:
				*style = (*style).Background(ansiColor(a - 40))
			case 90 <= a && a <= 97:
				*style = (*style).Foreground(ansiColor(a - 90 + 8))
			case 100 <= a && a <= 107:
				*style = (*style).Background(ansiColor(a - 100 + 8))

			// reset fg/bg
			case a == 39:
				*style = (*style).Foreground(tcell.ColorDefault)
			case a == 49:
				*style = (*style).Background(tcell.ColorDefault)

			// 256-color fg/bg
			case a == 38 && i+2 < len(args) && args[i+1] == 5:
				*style = (*style).Foreground(tcell.Color(args[i+2]))
				i += 2
			case a == 48 && i+2 < len(args) && args[i+1] == 5:
				*style = (*style).Background(tcell.Color(args[i+2]))
				i += 2

			// truecolor fg/bg
			case a == 38 && i+4 < len(args) && args[i+1] == 2:
				*style = (*style).Foreground(tcell.NewRGBColor(
					int32(args[i+2]),
					int32(args[i+3]),
					int32(args[i+4]),
				))
				i += 4
			case a == 48 && i+4 < len(args) && args[i+1] == 2:
				*style = (*style).Background(tcell.NewRGBColor(
					int32(args[i+2]),
					int32(args[i+3]),
					int32(args[i+4]),
				))
				i += 4
			}
		}
	}
}

func parseParams(params string) []int {
	params = strings.TrimSpace(params)
	if params == "" {
		return []int{}
	}

	parts := strings.Split(params, ";")
	result := make([]int, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if n, err := strconv.Atoi(p); err == nil {
			result = append(result, n)
		} else {
			continue
		}
	}

	return result
}

func getArg(args []int, index int, def int) int {
	if index >= len(args) || args[index] == 0 {
		return def
	}
	return args[index]
}

func ansiColor(n int) tcell.Color {
	switch n {
	case 0:
		return tcell.ColorBlack
	case 1:
		return tcell.ColorMaroon
	case 2:
		return tcell.ColorGreen
	case 3:
		return tcell.ColorOlive
	case 4:
		return tcell.ColorNavy
	case 5:
		return tcell.ColorPurple
	case 6:
		return tcell.ColorTeal
	case 7:
		return tcell.ColorSilver
	case 8:
		return tcell.ColorGray
	case 9:
		return tcell.ColorRed
	case 10:
		return tcell.ColorLime
	case 11:
		return tcell.ColorYellow
	case 12:
		return tcell.ColorBlue
	case 13:
		return tcell.ColorFuchsia
	case 14:
		return tcell.ColorAqua
	case 15:
		return tcell.ColorWhite
	default:
		return tcell.ColorDefault
	}
}
