package window

import (
	"image"
	"image/color"

	"github.com/gdamore/tcell/v2"
)

func FillRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	r, g, b, a := col.RGBA()
	c := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.SetRGBA(x+dx, y+dy, c)
		}
	}
}

func TcellColorToRGBA(c tcell.Color, isFg bool) color.Color {
	if c == tcell.ColorDefault {
		if isFg {
			return color.White
		}
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
