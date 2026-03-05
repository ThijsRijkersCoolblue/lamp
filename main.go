package main

import (
	"image"
	"image/color"
	"image/draw"
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

	cmd := exec.Command("bash")
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
			ansi.ProcessOutput(simScreen, buf[:n], &cursorX, &cursorY, ansiState)
		}
	}()

	writeToPTY := func(b []byte) { ptmx.Write(b) }

	pixelW := window.Cols * window.CharW
	pixelH := window.Rows * window.CharH
	img := image.NewRGBA(image.Rect(0, 0, pixelW, pixelH))

	raster := canvas.NewRaster(func(w, h int) image.Image {
		draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

		cells, sw, sh := simScreen.GetContents()
		metrics := window.Face.Metrics()
		ascent := metrics.Ascent.Ceil()

		for r := 0; r < sh && r < window.Rows; r++ {
			for c := 0; c < sw && c < window.Cols; c++ {
				cell := cells[r*sw+c]
				ch := rune(' ')
				if len(cell.Runes) > 0 && cell.Runes[0] != 0 {
					ch = cell.Runes[0]
				}
				fg, bg, _ := cell.Style.Decompose()

				bgCol := window.TcellColorToRGBA(bg, false)
				if bgCol != color.Black {
					window.FillRect(img, c*window.CharW, r*window.CharH, window.CharW, window.CharH, bgCol)
				}

				if ch != ' ' {
					d := &xfont.Drawer{
						Dst:  img,
						Src:  image.NewUniform(window.TcellColorToRGBA(fg, true)),
						Face: window.Face,
						Dot:  fixed.P(c*window.CharW, r*window.CharH+ascent),
					}
					d.DrawString(string(ch))
				}
			}
		}

		cx, cy := cursorX, cursorY
		if cx >= 0 && cx < window.Cols && cy >= 0 && cy < window.Rows {
			window.FillRect(img, cx*window.CharW, cy*window.CharH, window.CharW, window.CharH,
				color.RGBA{R: 255, G: 255, B: 255, A: 80})
		}

		return img
	})

	logicalW := float32(pixelW) / 2
	logicalH := float32(pixelH) / 2
	raster.SetMinSize(fyne.NewSize(logicalW, logicalH))

	a := app.New()
	w := a.NewWindow("Lamp")
	w.SetContent(raster)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(logicalW, logicalH))

	go func() {
		ticker := time.NewTicker(time.Second / 30)
		for range ticker.C {
			raster.Refresh()
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

	w.ShowAndRun()
}
