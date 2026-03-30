// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"image"
	"log"

	"github.com/derekmu/g3n/math32"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Font represents a font face.
type Font struct {
	font       *opentype.Font
	faceCache  map[FontAttributes]font.Face
	face       font.Face
	attributes FontAttributes
	color      image.Uniform
}

// FontAttributes holds the tunable attributes of a font face.
type FontAttributes struct {
	PointSize int32
	DPI       int32
	Hinting   font.Hinting
}

// makeFaceOptions returns true type font options with these font attributes.
func (a *FontAttributes) makeFaceOptions() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    float64(a.PointSize),
		DPI:     float64(a.DPI),
		Hinting: a.Hinting,
	}
}

// NewFontFromData creates and returns a new font object from the specified TTF data.
func NewFontFromData(fontData []byte) (*Font, error) {
	ttf, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}
	f := new(Font)
	f.font = ttf
	f.attributes = FontAttributes{
		PointSize: 12,
		DPI:       72,
		Hinting:   font.HintingFull,
	}
	f.SetColor(math32.Color3{})
	f.faceCache = make(map[FontAttributes]font.Face)
	return f, nil
}

// SetPointSize sets the point size of the font.
func (f *Font) SetPointSize(size int32) {
	f.attributes.PointSize = size
}

// SetDPI sets the resolution of the font in dots per inches (DPI).
func (f *Font) SetDPI(dpi int32) {
	f.attributes.DPI = dpi
}

// SetHinting sets the hinting type.
func (f *Font) SetHinting(hinting font.Hinting) {
	f.attributes.Hinting = hinting
}

// SetAttributes sets the font attributes.
func (f *Font) SetAttributes(attributes FontAttributes) {
	f.attributes = attributes
}

// SetColor sets the text color.
func (f *Font) SetColor(c math32.Color) {
	f.color.C = c.ToNRGBA()
}

// Metrics returns the font metrics.
func (f *Font) Metrics() font.Metrics {
	f.updateFace()
	return f.face.Metrics()
}

// updateFace updates the font face, creating a new face if the current attributes aren't cached.
func (f *Font) updateFace() {
	var ok bool
	f.face, ok = f.faceCache[f.attributes]
	if !ok {
		face, err := opentype.NewFace(f.font, f.attributes.makeFaceOptions())
		if err != nil {
			log.Panic(err)
		}
		f.face = face
		f.faceCache[f.attributes] = face
	}
}

// MeasureBytes returns the dimensions necessary for an image to contain the text.
func (f *Font) MeasureBytes(text []byte) (width, height int) {
	d := font.Drawer(f.Drawer(0, 0, nil))
	width = d.MeasureBytes(text).Ceil()
	metrics := f.Metrics()
	height = (metrics.Ascent + metrics.Descent).Ceil()
	return width, height
}

// MeasureString returns the dimensions necessary for an image to contain the text.
func (f *Font) MeasureString(text string) (width, height int) {
	d := font.Drawer(f.Drawer(0, 0, nil))
	width = d.MeasureString(text).Ceil()
	metrics := f.Metrics()
	height = (metrics.Ascent + metrics.Descent).Ceil()
	return width, height
}

// MeasureRunes returns the dimensions necessary for an image to contain the text.
func (f *Font) MeasureRunes(text []rune, widths []fixed.Int26_6) (width, height int) {
	d := f.Drawer(0, 0, nil)
	width = d.MeasureRunes(text, widths).Ceil()
	metrics := f.Metrics()
	height = (metrics.Ascent + metrics.Descent).Ceil()
	return width, height
}

// DrawBytes draws the text on the image.
func (f *Font) DrawBytes(text []byte, x, y int, dst *image.RGBA) {
	d := font.Drawer(f.Drawer(x, y, dst))
	d.DrawBytes(text)
}

// DrawString draws the text on the image.
func (f *Font) DrawString(text string, x, y int, dst *image.RGBA) {
	d := font.Drawer(f.Drawer(x, y, dst))
	d.DrawString(text)
}

// DrawRunes draws the text on the image.
func (f *Font) DrawRunes(text []rune, x, y int, dst *image.RGBA) {
	d := f.Drawer(x, y, dst)
	d.DrawRunes(text)
}

// Drawer returns a Drawer with this Font's settings.
func (f *Font) Drawer(x, y int, dst *image.RGBA) Drawer {
	metrics := f.Metrics()
	return Drawer{
		Dst:  dst,
		Src:  &f.color,
		Face: f.face,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y) + metrics.Ascent,
		},
	}
}
