// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"image"
	"image/draw"
	"math"
	"os"
	"strings"

	"github.com/derekmu/g3n/math32"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Font represents a TrueType font face.
type Font struct {
	ttf            *truetype.Font // The TrueType font
	face           font.Face      // The font face
	attrib         FontAttributes // Internal attribute cache
	fg             *image.Uniform // Text color cache
	bg             *image.Uniform // Background color cache
	scaleX, scaleY float64        // Scales of actual pixel/GL point, used for fix Retina Monitor
	changed        bool           // Whether attributes have changed and the font face needs to be recreated
}

// FontAttributes contains tunable attributes of a font.
type FontAttributes struct {
	PointSize   float64      // Point size of the font
	DPI         float64      // Resolution of the font in dots per inch
	LineSpacing float64      // Spacing between lines (in terms of font height)
	Hinting     font.Hinting // Font hinting
}

func (a *FontAttributes) newTTOptions(scaleX, scaleY float64) *truetype.Options {
	dpi := a.DPI
	if scaleX != 0 && scaleY != 0 {
		dpi *= math.Sqrt(scaleX * scaleY)
	}
	return &truetype.Options{
		Size:    a.PointSize,
		DPI:     dpi,
		Hinting: a.Hinting,
	}
}

// NewFont creates and returns a new font object using the specified TrueType font file.
func NewFont(ttfFile string) (*Font, error) {
	// Reads font bytes
	fontBytes, err := os.ReadFile(ttfFile)
	if err != nil {
		return nil, err
	}
	return NewFontFromData(fontBytes)
}

// NewFontFromData creates and returns a new font object from the specified TTF data.
func NewFontFromData(fontData []byte) (*Font, error) {
	// Parses the font data
	ttf, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	f := new(Font)
	f.ttf = ttf

	// InitPanel with default values
	f.attrib = FontAttributes{}
	f.attrib.PointSize = 12
	f.attrib.DPI = 72
	f.attrib.LineSpacing = 1.0
	f.attrib.Hinting = font.HintingNone
	f.SetColor(math32.Color4{0, 0, 0, 1})

	// Create font face
	f.face = truetype.NewFace(f.ttf, f.attrib.newTTOptions(f.scaleX, f.scaleY))

	return f, nil
}

// SetPointSize sets the point size of the font.
func (f *Font) SetPointSize(size float64) {
	if math.Abs(size-f.attrib.PointSize) < 0.00001 {
		return
	}
	f.attrib.PointSize = size
	f.changed = true
}

// SetDPI sets the resolution of the font in dots per inches (DPI).
func (f *Font) SetDPI(dpi float64) {
	if math.Abs(dpi-f.attrib.DPI) < 0.00001 {
		return
	}
	f.attrib.DPI = dpi
	f.changed = true
}

// SetLineSpacing sets the amount of spacing between lines (in terms of font height).
func (f *Font) SetLineSpacing(spacing float64) {
	if math.Abs(spacing-f.attrib.LineSpacing) < 0.00001 {
		return
	}
	f.attrib.LineSpacing = spacing
	f.changed = true
}

// SetHinting sets the hinting type.
func (f *Font) SetHinting(hinting font.Hinting) {
	if hinting == f.attrib.Hinting {
		return
	}
	f.attrib.Hinting = hinting
	f.changed = true
}

func (f *Font) ScaleXY() (x, y float64) {
	return f.scaleX, f.scaleY
}

func (f *Font) ScaleX() float64 {
	return f.scaleX
}

func (f *Font) ScaleY() float64 {
	return f.scaleY
}

// SetScaleXY sets the ratio of actual pixel/GL point.
func (f *Font) SetScaleXY(x, y float64) {
	if x == f.scaleX && y == f.scaleY {
		return
	}
	f.scaleX = x
	f.scaleY = y
	f.changed = true
}

// SetFgColor sets the text color.
func (f *Font) SetFgColor(color math32.Color4) {
	f.fg = image.NewUniform(color)
}

// SetBgColor sets the background color.
func (f *Font) SetBgColor(color math32.Color4) {
	f.bg = image.NewUniform(color)
}

// SetColor sets the text color to the specified value and makes the background color transparent.
// Note that for perfect transparency in the anti-aliased region it's important that the RGB components
// of the text and background colors match. This method handles that for the user.
func (f *Font) SetColor(fg math32.Color4) {
	f.fg = image.NewUniform(fg)
	f.bg = image.NewUniform(math32.Color4{fg.R, fg.G, fg.B, 0})
}

// SetAttributes sets the font attributes.
func (f *Font) SetAttributes(fa *FontAttributes) {
	f.SetPointSize(fa.PointSize)
	f.SetDPI(fa.DPI)
	f.SetLineSpacing(fa.LineSpacing)
	f.SetHinting(fa.Hinting)
}

// updateFace updates the font face if parameters have changed.
func (f *Font) updateFace() {
	if f.changed {
		f.face = truetype.NewFace(f.ttf, f.attrib.newTTOptions(f.scaleX, f.scaleY))
		f.changed = false
	}
}

// MeasureText returns the minimum width and height in pixels necessary for an image to contain the specified text.
// The supplied text string can contain line breaks.
func (f *Font) MeasureText(text string) (int, int) {
	// Create font drawer
	f.updateFace()
	d := font.Drawer{Dst: nil, Src: f.fg, Face: f.face}

	// Draw text
	var width, height int
	metrics := f.face.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.attrib.LineSpacing - float64(1)) * float64(lineHeight))

	lines := strings.Split(text, "\n")
	for i, s := range lines {
		d.Dot = fixed.P(0, height)
		lineWidth := d.MeasureString(s).Ceil()
		if lineWidth > width {
			width = lineWidth
		}
		height += lineHeight
		if i > 1 {
			height += lineGap
		}
	}
	return width, height
}

// Metrics returns the font metrics.
func (f *Font) Metrics() font.Metrics {
	f.updateFace()
	return f.face.Metrics()
}

// DrawText draws the specified text on a new, tightly fitting image, and returns a pointer to the image.
func (f *Font) DrawText(text string) *image.RGBA {
	width, height := f.MeasureText(text)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), f.bg, image.Point{}, draw.Src)
	f.DrawTextOnImage(text, 0, 0, img)
	return img
}

// DrawTextOnImage draws the specified text on the specified image at the specified coordinates.
func (f *Font) DrawTextOnImage(text string, x, y int, dst *image.RGBA) {
	f.updateFace()
	d := font.Drawer{Dst: dst, Src: f.fg, Face: f.face}

	// Draw text
	metrics := f.face.Metrics()
	py := y + metrics.Ascent.Round()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.attrib.LineSpacing - float64(1)) * float64(lineHeight))
	lines := strings.Split(text, "\n")
	for i, s := range lines {
		d.Dot = fixed.P(x, py)
		d.DrawString(s)
		py += lineHeight
		if i > 1 {
			py += lineGap
		}
	}
}

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
	d := font.Drawer{Dst: c.RGBA, Src: f.fg, Face: f.face}

	// fill with background color
	draw.Draw(c.RGBA, c.RGBA.Bounds(), image.NewUniform(c.BgColor), image.Point{}, draw.Src)

	// Draw text
	actualPointSize := int(f.attrib.PointSize * f.scaleY)
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
			for i := 0; i < int(f.scaleX); i++ {
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
