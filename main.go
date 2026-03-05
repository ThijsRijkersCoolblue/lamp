package main

import (
	"image"
	"image/color"
	"image/draw"
	"lamp/ansi"
	"lamp/window"
	"log"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/creack/pty"
	"github.com/gdamore/tcell/v2"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	fontSize = 15.0 // pt — change this to taste
	cols     = 200
	rows     = 50
)

var (
	charW, charH int
	face         xfont.Face
)

func initFont() {
	// Use the system monospace font on macOS
	fontPaths := []string{
		"/System/Library/Fonts/Menlo.ttc",
		"/Library/Fonts/Courier New.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf", // Linux fallback
	}
	var fontData []byte
	for _, p := range fontPaths {
		data, err := os.ReadFile(p)
		if err == nil {
			fontData = data
			break
		}
	}
	if fontData == nil {
		log.Fatal("no monospace font found — install Menlo or DejaVu Sans Mono")
	}

	collection, err := opentype.ParseCollection(fontData)
	var f *opentype.Font
	if err != nil {
		// Not a collection, try as single font
		f, err = opentype.Parse(fontData)
		if err != nil {
			log.Fatal("failed to parse font:", err)
		}
	} else {
		f, err = collection.Font(0)
		if err != nil {
			log.Fatal("failed to get font from collection:", err)
		}
	}

	const dpi = 144 // 2x retina (72dpi * 2)
	face, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		log.Fatal("failed to create font face:", err)
	}

	// Measure actual glyph dimensions
	metrics := face.Metrics()
	charH = (metrics.Ascent + metrics.Descent).Ceil()
	advance, _, _ := face.GlyphBounds('M')
	charW = (advance.Max.X - advance.Min.X).Ceil()
	if charW == 0 {
		charW = charH / 2
	}
}

func main() {
	initFont()

	cmd := exec.Command("bash")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer ptmx.Close()

	simScreen := tcell.NewSimulationScreen("UTF-8")
	simScreen.Init()
	simScreen.SetSize(cols, rows)
	pty.Setsize(ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})

	cursorX, cursorY := 0, 0
	ansiState := ansi.NewState(rows)

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

	// Pixel buffer at full retina resolution
	pixelW := cols * charW
	pixelH := rows * charH
	img := image.NewRGBA(image.Rect(0, 0, pixelW, pixelH))

	raster := canvas.NewRaster(func(w, h int) image.Image {
		draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

		cells, sw, sh := simScreen.GetContents()
		metrics := face.Metrics()
		ascent := metrics.Ascent.Ceil()

		for r := 0; r < sh && r < rows; r++ {
			for c := 0; c < sw && c < cols; c++ {
				cell := cells[r*sw+c]
				ch := rune(' ')
				if len(cell.Runes) > 0 && cell.Runes[0] != 0 {
					ch = cell.Runes[0]
				}
				fg, bg, _ := cell.Style.Decompose()

				// Draw background
				bgCol := tcellColorToRGBA(bg, false)
				if bgCol != color.Black {
					fillRect(img, c*charW, r*charH, charW, charH, bgCol)
				}

				// Draw glyph
				if ch != ' ' {
					d := &xfont.Drawer{
						Dst:  img,
						Src:  image.NewUniform(tcellColorToRGBA(fg, true)),
						Face: face,
						Dot:  fixed.P(c*charW, r*charH+ascent),
					}
					d.DrawString(string(ch))
				}
			}
		}

		// Cursor block
		cx, cy := cursorX, cursorY
		if cx >= 0 && cx < cols && cy >= 0 && cy < rows {
			fillRect(img, cx*charW, cy*charH, charW, charH,
				color.RGBA{R: 255, G: 255, B: 255, A: 80})
		}

		return img
	})

	// Logical size is half of pixel size (retina 2x)
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
		if ev := fyneKeyToTcell(e); ev != nil {
			window.HandleEvent(simScreen, ev, writeToPTY)
		}
	})
	w.Canvas().SetOnTypedRune(func(r rune) {
		ev := tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
		window.HandleEvent(simScreen, ev, writeToPTY)
	})

	w.ShowAndRun()
}

func fillRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	r, g, b, a := col.RGBA()
	c := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.SetRGBA(x+dx, y+dy, c)
		}
	}
}

func tcellColorToRGBA(c tcell.Color, isFg bool) color.Color {
	if c == tcell.ColorDefault {
		if isFg { return color.White }
		return color.Black
	}
	if c == tcell.ColorWhite {
		return color.White
	}
	if c == tcell.ColorBlack {
		return color.Black
	}
	r, g, b := c.RGB()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

func fyneKeyToTcell(e *fyne.KeyEvent) *tcell.EventKey {
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