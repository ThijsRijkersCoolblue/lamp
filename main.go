package main

import (
	"image"
	"image/color"
	"lamp/ansi"
	"lamp/events"
	"lamp/window"
	"log"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/creack/pty"
	"github.com/gdamore/tcell/v2"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	window.InitFont()

	window.Cols = 188
	window.Rows = 45

	cellW := window.CharW / 2
	cellH := window.CharH / 2

	logicalW := float32(window.Cols * cellW)
	logicalH := float32(window.Rows * cellH)

	cmd := exec.Command("bash", "--login")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer ptmx.Close()

	simScreen := tcell.NewSimulationScreen("UTF-8")
	simScreen.Init()
	simScreen.SetSize(window.Cols, window.Rows)
	pty.Setsize(ptmx, &pty.Winsize{Cols: uint16(window.Cols), Rows: uint16(window.Rows)})

	cursorX, cursorY := 0, 0
	ansiState := ansi.NewState(window.Rows)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				return
			}
			data := buf[:n]
			if len(ansiState.Leftover) > 0 {
				data = append(ansiState.Leftover, data...)
				ansiState.Leftover = nil
			}
			ansi.ProcessOutput(simScreen, data, &cursorX, &cursorY, ansiState)
		}
	}()

	writeToPTY := func(b []byte) { ptmx.Write(b) }

	raster := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))

		cells, sw, sh := simScreen.GetContents()
		ascent := window.Face.Metrics().Ascent.Ceil()

		cw := w / window.Cols
		ch := h / window.Rows
		if cw == 0 {
			cw = 1
		}
		if ch == 0 {
			ch = 1
		}

		for row := 0; row < sh && row < window.Rows; row++ {
			for col := 0; col < sw && col < window.Cols; col++ {
				cell := cells[row*sw+col]
				r := rune(' ')
				if len(cell.Runes) > 0 && cell.Runes[0] != 0 {
					r = cell.Runes[0]
				}
				fg, bg, _ := cell.Style.Decompose()
				x, y := col*cw, row*ch

				window.FillRect(img, x, y, cw, ch, window.TcellColorToRGBA(bg, false))

				if r != ' ' {
					(&xfont.Drawer{
						Dst:  img,
						Src:  image.NewUniform(window.TcellColorToRGBA(fg, true)),
						Face: window.Face,
						Dot:  fixed.P(x, y+ascent),
					}).DrawString(string(r))
				}
			}
		}

		cx, cy := cursorX, cursorY
		if cx >= 0 && cx < window.Cols && cy >= 0 && cy < window.Rows {
			window.FillRect(img, cx*cw, cy*ch, cw, ch,
				color.RGBA{R: 255, G: 255, B: 255, A: 80})
		}

		return img
	})

	raster.SetMinSize(fyne.NewSize(logicalW, logicalH))

	a := app.New()
	w := a.NewWindow("Lamp")
	w.SetContent(raster)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(logicalW, logicalH))

	go func() {
		ticker := time.NewTicker(time.Second / 30)
		for range ticker.C {
			fyne.Do(func() {
				raster.Refresh()
			})
		}
	}()

	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if ev := window.FyneKeyToTcell(e); ev != nil {
			events.HandleEvent(simScreen, ev, writeToPTY)
		}
	})
	w.Canvas().SetOnTypedRune(func(r rune) {
		ev := tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
		events.HandleEvent(simScreen, ev, writeToPTY)
	})

	w.Show()
	go func() {
		time.Sleep(50 * time.Millisecond)
		fyne.Do(func() {
			w.Resize(fyne.NewSize(logicalW+1, logicalH+1))
			w.Resize(fyne.NewSize(logicalW, logicalH))
		})
	}()
	w.ShowAndRun()
}
