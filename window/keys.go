package window

import (
	"fyne.io/fyne/v2"
	"github.com/gdamore/tcell/v2"
)

func FyneKeyToTcell(e *fyne.KeyEvent) *tcell.EventKey {
	keyMap := map[fyne.KeyName]tcell.Key{
		fyne.KeyReturn:    tcell.KeyEnter,
		fyne.KeyEnter:     tcell.KeyEnter,
		fyne.KeyTab:       tcell.KeyTab,
		fyne.KeyBackspace: tcell.KeyBackspace2,
		fyne.KeyEscape:    tcell.KeyEscape,
		fyne.KeyUp:        tcell.KeyUp,
		fyne.KeyDown:      tcell.KeyDown,
		fyne.KeyLeft:      tcell.KeyLeft,
		fyne.KeyRight:     tcell.KeyRight,
		fyne.KeyHome:      tcell.KeyHome,
		fyne.KeyEnd:       tcell.KeyEnd,
		fyne.KeyPageUp:    tcell.KeyPgUp,
		fyne.KeyPageDown:  tcell.KeyPgDn,
		fyne.KeyDelete:    tcell.KeyDelete,
		fyne.KeyF1:        tcell.KeyF1,
		fyne.KeyF2:        tcell.KeyF2,
		fyne.KeyF3:        tcell.KeyF3,
		fyne.KeyF4:        tcell.KeyF4,
		fyne.KeyF5:        tcell.KeyF5,
		fyne.KeyF6:        tcell.KeyF6,
		fyne.KeyF7:        tcell.KeyF7,
		fyne.KeyF8:        tcell.KeyF8,
		fyne.KeyF9:        tcell.KeyF9,
		fyne.KeyF10:       tcell.KeyF10,
		fyne.KeyF11:       tcell.KeyF11,
		fyne.KeyF12:       tcell.KeyF12,
	}
	if k, ok := keyMap[e.Name]; ok {
		return tcell.NewEventKey(k, 0, tcell.ModNone)
	}
	return nil
}
