package text

import (
	"github.com/derekmu/g3n/math32"
	"image"
	"image/draw"
)

// Canvas is an image to draw on.
type Canvas struct {
	RGBA    *image.RGBA
	BgColor math32.Color
}

// NewCanvas creates a new Canvas with the specified width, height, and background color.
func NewCanvas(width, height int, bgColor math32.Color) *Canvas {
	return &Canvas{
		RGBA:    image.NewRGBA(image.Rect(0, 0, width, height)),
		BgColor: bgColor,
	}
}

// DrawText draws text at the specified position of this canvas, using the specified font.
func (c *Canvas) DrawText(x, y int, text string, f *Font) {
	draw.Draw(c.RGBA, c.RGBA.Bounds(), image.NewUniform(c.BgColor.RGBA()), image.Point{}, draw.Src)
	f.DrawText(text, x, y, c.RGBA)
}
