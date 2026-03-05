package window

import (
	"log"
	"os"

	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	FontSize = 15.0
	Cols     = 200
	Rows     = 50
)

var (
	CharW, CharH int
	Face         xfont.Face
)

func InitFont() {
	fontPaths := []string{
		"/System/Library/Fonts/Menlo.ttc",
		"/Library/Fonts/Courier New.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
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

	const dpi = 144
	Face, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    FontSize,
		DPI:     dpi,
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		log.Fatal("failed to create font face:", err)
	}

	metrics := Face.Metrics()
	CharH = (metrics.Ascent + metrics.Descent).Ceil()
	advance, _, _ := Face.GlyphBounds('M')
	CharW = (advance.Max.X - advance.Min.X).Ceil()
	if CharW == 0 {
		CharW = CharH / 2
	}
}
