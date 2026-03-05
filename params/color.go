package params

import "github.com/gdamore/tcell/v2"

func AnsiColor(n int) tcell.Color {
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
