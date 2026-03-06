package ansi

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

func ProcessOutput(screen tcell.Screen, data []byte, cursorX, cursorY *int, state *State) {
	width, height := screen.Size()

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
			*cursorY++
		case '\t':
			next := (*cursorX + 8) &^ 7
			*cursorX = next
			if *cursorX >= width {
				*cursorX = width - 1
			}
		case 0x08, 0x7f:
			if *cursorX > 0 {
				*cursorX--
				screen.SetContent(*cursorX, *cursorY, ' ', nil, state.Style)
			}
		case 0x1b:
			if i+1 >= len(data) {
				state.Leftover = append([]byte{}, data[i:]...)
				screen.Show()
				return
			}
			switch data[i+1] {
			case '[':
				j := i + 2
				for j < len(data) && (data[j] < '@' || data[j] > '~') {
					j++
				}
				if j >= len(data) {
					state.Leftover = append([]byte{}, data[i:]...)
					screen.Show()
					return
				}
				cmd := data[j]
				params := string(data[i+2 : j])

				// alternate screen enter
				if cmd == 'h' && (params == "?1049" || params == "?1047" || params == "?47") {
					clearScreen(screen, state.Style)
					*cursorX = 0
					*cursorY = 0
					state.ScrollTop = 0
					state.ScrollBottom = height - 1
				}
				// alternate screen exit
				if cmd == 'l' && (params == "?1049" || params == "?1047" || params == "?47") {
					clearScreen(screen, state.Style)
					*cursorX = 0
					*cursorY = 0
					state.ScrollTop = 0
					state.ScrollBottom = height - 1
				}
				if params == "" || params[0] != '?' {
					handleCSI(screen, cmd, params, cursorX, cursorY, &state.Style, width, height, state)
				}
				i = j
			case ']':
				j := i + 2
				for j < len(data) {
					if data[j] == 0x07 {
						break
					}
					if data[j] == 0x1b && j+1 < len(data) && data[j+1] == '\\' {
						j++
						break
					}
					j++
				}
				if j >= len(data) {
					state.Leftover = append([]byte{}, data[i:]...)
					screen.Show()
					return
				}
				i = j
			case '(', ')', '*', '+', '-', '.', '/':
				i += 2
				if i < len(data) {
					i++
				}
				continue
			default:
				i++
			}
		default:
			if ch >= 32 {
				r, size := utf8.DecodeRune(data[i:])
				x, y := *cursorX, *cursorY
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

		if *cursorY > state.ScrollBottom {
			scrollUp(screen, state.Style, width, state.ScrollTop, state.ScrollBottom)
			*cursorY = state.ScrollBottom
		}

		i++
	}

	screen.Show()
}

func clearScreen(screen tcell.Screen, style tcell.Style) {
	w, h := screen.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}
}
