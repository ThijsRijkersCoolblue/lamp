package events

import (
	"github.com/gdamore/tcell/v2"
)

func HandleEvent(screen tcell.Screen, event tcell.Event, write func([]byte)) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			write([]byte{3})
		case tcell.KeyCtrlD:
			write([]byte{4})
		case tcell.KeyCtrlZ:
			write([]byte{26})
		case tcell.KeyCtrlL:
			write([]byte{12})

		case tcell.KeyEnter:
			write([]byte{'\r'})

		case tcell.KeyBackspace, tcell.KeyBackspace2:
			write([]byte{0x7f})

		case tcell.KeyTab:
			write([]byte{'\t'})

		case tcell.KeyUp:
			write([]byte{0x1b, '[', 'A'})
		case tcell.KeyDown:
			write([]byte{0x1b, '[', 'B'})
		case tcell.KeyRight:
			write([]byte{0x1b, '[', 'C'})
		case tcell.KeyLeft:
			write([]byte{0x1b, '[', 'D'})

		case tcell.KeyHome:
			write([]byte{0x1b, '[', 'H'})
		case tcell.KeyEnd:
			write([]byte{0x1b, '[', 'F'})
		case tcell.KeyPgUp:
			write([]byte{0x1b, '[', '5', '~'})
		case tcell.KeyPgDn:
			write([]byte{0x1b, '[', '6', '~'})
		case tcell.KeyDelete:
			write([]byte{0x1b, '[', '3', '~'})
		case tcell.KeyInsert:
			write([]byte{0x1b, '[', '2', '~'})

		case tcell.KeyEscape:
			write([]byte{0x1b})

		case tcell.KeyF1:
			write([]byte{0x1b, 'O', 'P'})
		case tcell.KeyF2:
			write([]byte{0x1b, 'O', 'Q'})
		case tcell.KeyF3:
			write([]byte{0x1b, 'O', 'R'})
		case tcell.KeyF4:
			write([]byte{0x1b, 'O', 'S'})
		case tcell.KeyF5:
			write([]byte{0x1b, '[', '1', '5', '~'})
		case tcell.KeyF6:
			write([]byte{0x1b, '[', '1', '7', '~'})
		case tcell.KeyF7:
			write([]byte{0x1b, '[', '1', '8', '~'})
		case tcell.KeyF8:
			write([]byte{0x1b, '[', '1', '9', '~'})
		case tcell.KeyF9:
			write([]byte{0x1b, '[', '2', '0', '~'})
		case tcell.KeyF10:
			write([]byte{0x1b, '[', '2', '1', '~'})
		case tcell.KeyF11:
			write([]byte{0x1b, '[', '2', '3', '~'})
		case tcell.KeyF12:
			write([]byte{0x1b, '[', '2', '4', '~'})

		default:
			r := ev.Rune()
			if r != 0 {
				write([]byte(string(r)))
			}
		}
	}
}
