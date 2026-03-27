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

// Resize updates the image to the size requested.
// If the buffer is large enough, Resize reuses it. Otherwise, it allocates a new buffer.
// Does not preserve pixel locations after resizing. The entire image should be redrawn.
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

// Width is the width of the image.
func (c *Canvas) Width() int {
	return c.Rect.Dx()
}

// Height is the height of the image.
func (c *Canvas) Height() int {
	return c.Rect.Dy()
}

// DrawBytes draws text on the image.
func (c *Canvas) DrawBytes(x, y int, text []byte, f *Font) *Canvas {
	f.DrawBytes(text, x, y, &c.RGBA)
	return c
}

// DrawString draws text on the image.
func (c *Canvas) DrawString(x, y int, text string, f *Font) *Canvas {
	f.DrawString(text, x, y, &c.RGBA)
	return c
}

// DrawRectangle draws a rectangle on the image.
func (c *Canvas) DrawRectangle(bounds image.Rectangle, co color.Color) *Canvas {
	draw.Draw(
		c,
		bounds,
		&image.Uniform{C: co},
		image.Point{},
		draw.Src,
	)
	return c
}

// Fill completely fills the image with a color.
func (c *Canvas) Fill(co color.Color) *Canvas {
	return c.DrawRectangle(c.Bounds(), co)
}
