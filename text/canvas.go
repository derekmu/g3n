package text

import (
	"github.com/derekmu/g3n/math32"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"strings"
)

// Canvas is an image to draw on.
type Canvas struct {
	RGBA    *image.RGBA
	BgColor math32.Color4
}

// NewCanvas creates a new Canvas with the specified width, height, and background color.
func NewCanvas(width, height int, bgColor math32.Color4) *Canvas {
	return &Canvas{
		RGBA:    image.NewRGBA(image.Rect(0, 0, width, height)),
		BgColor: bgColor,
	}
}

// DrawText draws text at the specified position of this canvas, using the specified font.
// The supplied text string can contain line breaks
func (c *Canvas) DrawText(x, y int, text string, f *Font) {
	// fill with background color
	draw.Draw(c.RGBA, c.RGBA.Bounds(), image.NewUniform(c.BgColor), image.Point{}, draw.Src)
	f.DrawTextOnImage(text, x, y, c.RGBA)
}

// DrawTextCaret draws text and a caret at the specified position, in pixels, of this canvas.
// The supplied text string can contain line breaks.
func (c *Canvas) DrawTextCaret(x, y int, text string, f *Font, drawCaret bool, line, col, selStart, selEnd int) {
	// Creates drawer
	f.updateFace()
	d := font.Drawer{Dst: c.RGBA, Src: &f.fg, Face: f.face}

	// fill with background color
	draw.Draw(c.RGBA, c.RGBA.Bounds(), image.NewUniform(c.BgColor), image.Point{}, draw.Src)

	// Draw text
	actualPointSize := int(f.attrib.PointSize * f.attrib.ScaleY)
	metrics := f.face.Metrics()
	py := y + metrics.Ascent.Round()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.attrib.LineSpacing - float64(1)) * float64(lineHeight))
	lines := strings.Split(text, "\n")
	for l, s := range lines {
		d.Dot = fixed.P(x, py)
		if selStart != selEnd && l == line && selEnd <= StrCount(s) {
			width, _ := f.MeasureText(StrPrefix(s, selStart))
			widthEnd, _ := f.MeasureText(StrPrefix(s, selEnd))
			// Draw selection caret
			// This will not work when the selection spans multiple lines
			// Currently there is no multiline edit text
			// Once there is, this needs to change
			caretH := actualPointSize + 2
			caretY := int(d.Dot.Y>>6) - actualPointSize + 2
			for w := width; w < widthEnd; w++ {
				for j := caretY; j < caretY+caretH; j++ {
					c.RGBA.Set(x+w, j, math32.Color4{0, 0, 1, 0.5}) // blue
				}
			}
		}
		d.DrawString(s)
		// Checks for caret position
		if drawCaret && l == line && col <= StrCount(s) {
			width, _ := f.MeasureText(StrPrefix(s, col))
			// Draw caret vertical line
			caretH := actualPointSize + 2
			caretY := int(d.Dot.Y>>6) - actualPointSize + 2
			for i := 0; i < int(f.attrib.ScaleX); i++ {
				for j := caretY; j < caretY+caretH; j++ {
					c.RGBA.Set(x+width+i, j, math32.Color4{0, 0, 0, 1}) // black
				}
			}
		}
		py += lineHeight
		if l > 1 {
			py += lineGap
		}
	}
}
