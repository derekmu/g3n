// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"image"
	"image/draw"
	"math"
	"strings"

	"github.com/derekmu/g3n/math32"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Font represents a TrueType font face.
type Font struct {
	ttf       *truetype.Font
	face      font.Face
	attrib    FontAttributes
	fg        image.Uniform
	bg        image.Uniform
	faceCache map[faceKey]font.Face
}

// faceKey is a FontAttributes struct suitable for being a key in a map.
type faceKey struct {
	PointSize      int64
	DPI            int64
	ScaleX, ScaleY int64
	Hinting        font.Hinting
}

// FontAttributes contains tunable attributes of a font.
type FontAttributes struct {
	PointSize      float64
	DPI            float64
	ScaleX, ScaleY float64
	Hinting        font.Hinting
	LineSpacing    float64
}

// newTrueTypeOptions returns a key object for these font attributes.
func (a *FontAttributes) makeFaceKey() faceKey {
	// Since float64s are finicky to compare, convert to int64s while keeping a few decimal places of accuracy
	// Line spacing doesn't pertain to the font face, so isn't included in the key
	return faceKey{
		PointSize: int64(a.PointSize * 10000),
		DPI:       int64(a.PointSize * 10000),
		ScaleX:    int64(a.PointSize * 10000),
		ScaleY:    int64(a.PointSize * 10000),
		Hinting:   a.Hinting,
	}
}

// newTrueTypeOptions returns true type font options with these font attributes.
func (a *FontAttributes) newTrueTypeOptions() *truetype.Options {
	dpi := a.DPI
	if a.ScaleX != 0 && a.ScaleY != 0 {
		dpi *= math.Sqrt(a.ScaleX * a.ScaleY)
	}
	return &truetype.Options{
		Size:    a.PointSize,
		DPI:     dpi,
		Hinting: a.Hinting,
	}
}

// NewFontFromData creates and returns a new font object from the specified TTF data.
func NewFontFromData(fontData []byte) (*Font, error) {
	ttf, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}
	f := new(Font)
	f.ttf = ttf
	f.attrib = FontAttributes{
		PointSize:   12,
		DPI:         72,
		ScaleX:      1.0,
		ScaleY:      1.0,
		Hinting:     font.HintingFull,
		LineSpacing: 1.0,
	}
	f.SetColor(math32.Color4{0, 0, 0, 1})
	return f, nil
}

// SetPointSize sets the point size of the font.
func (f *Font) SetPointSize(size float64) {
	f.attrib.PointSize = size
}

// SetDPI sets the resolution of the font in dots per inches (DPI).
func (f *Font) SetDPI(dpi float64) {
	f.attrib.DPI = dpi
}

// SetLineSpacing sets the amount of spacing between lines (in terms of font height).
func (f *Font) SetLineSpacing(spacing float64) {
	f.attrib.LineSpacing = spacing
}

// SetHinting sets the hinting type.
func (f *Font) SetHinting(hinting font.Hinting) {
	f.attrib.Hinting = hinting
}

// SetScaleXY sets the ratio of actual pixel/GL point.
func (f *Font) SetScaleXY(x, y float64) {
	f.attrib.ScaleX = x
	f.attrib.ScaleY = y
}

// SetAttributes sets the font attributes.
func (f *Font) SetAttributes(attrib FontAttributes) {
	f.attrib = attrib
}

// SetFgColor sets the text color.
func (f *Font) SetFgColor(color math32.Color4) {
	f.fg.C = color
}

// SetBgColor sets the background color.
func (f *Font) SetBgColor(color math32.Color4) {
	f.bg.C = color
}

// SetColor sets the text color to the specified value and makes the background color transparent.
// Note that for perfect transparency in the anti-aliased region it's important that the RGB components
// of the text and background colors match. This method handles that for the user.
func (f *Font) SetColor(fg math32.Color4) {
	f.fg.C = fg
	f.bg.C = math32.Color4{fg.R, fg.G, fg.B, 0}
}

// updateFace updates the font face if parameters have changed.
func (f *Font) updateFace() {
	key := f.attrib.makeFaceKey()
	var ok bool
	f.face, ok = f.faceCache[key]
	if !ok {
		f.face = truetype.NewFace(f.ttf, f.attrib.newTrueTypeOptions())
		f.faceCache[key] = f.face
	}
}

// MeasureText returns the minimum width and height in pixels necessary for an image to contain the specified text.
// The supplied text string can contain line breaks.
func (f *Font) MeasureText(text string) (int, int) {
	// Create font drawer
	f.updateFace()
	d := font.Drawer{Dst: nil, Src: &f.fg, Face: f.face}

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
	draw.Draw(img, img.Bounds(), &f.bg, image.Point{}, draw.Src)
	f.DrawTextOnImage(text, 0, 0, img)
	return img
}

// DrawTextOnImage draws the specified text on the specified image at the specified coordinates.
func (f *Font) DrawTextOnImage(text string, x, y int, dst *image.RGBA) {
	f.updateFace()
	d := font.Drawer{Dst: dst, Src: &f.fg, Face: f.face}

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
