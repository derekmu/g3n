package text

import (
	"image"
	"image/color"
	"image/draw"
)

// Canvas is an image to draw on.
type Canvas struct {
	image.RGBA
	buffer []uint8
}

// NewCanvas creates a new Canvas with the specified width and height.
func NewCanvas(width, height int) *Canvas {
	c := &Canvas{}
	c.Resize(width, height)
	return c
}

// Resize updates the underlying image to the size requested.
// If the underlying buffer is large enough, Resize reuses it. Otherwise, it allocates a new buffer.
// Does not retain relative pixel values after resizing. The entire image should be redrawn after calling (see Fill).
func (c *Canvas) Resize(width, height int) *Canvas {
	pixLength := 4 * width * height
	if len(c.buffer) < pixLength {
		// Allocate a new image if the existing one is too small
		c.RGBA.Pix = make([]uint8, pixLength)
		c.buffer = c.RGBA.Pix
	} else {
		// Reuse the already allocated image
		c.RGBA.Pix = c.buffer[:pixLength]
	}
	c.RGBA.Stride = 4 * width
	c.RGBA.Rect = image.Rect(0, 0, width, height)
	return c
}

// DrawText draws text at the specified position of this canvas, using the specified font.
func (c *Canvas) DrawText(x, y int, text string, f *Font) *Canvas {
	f.DrawText(text, x, y, &c.RGBA)
	return c
}

// Fill completely fills the image with a color.
func (c *Canvas) Fill(co color.Color) *Canvas {
	draw.Draw(
		c,
		c.Bounds(),
		&image.Uniform{C: co},
		image.Point{},
		draw.Src,
	)
	return c
}
