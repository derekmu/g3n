// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"github.com/derekmu/g3n/math32"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"log"
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
	f.color.C = c.NRGBA()
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

// MeasureText returns the minimum width and height in pixels necessary for an image to contain the specified text.
func (f *Font) MeasureText(text string) (int, int) {
	f.updateFace()
	d := font.Drawer{Face: f.face, Dot: fixed.P(0, 0)}
	width := d.MeasureString(text).Ceil()
	metrics := f.face.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()
	return width, height
}

// DrawText draws the specified text on the specified image at the specified coordinates.
func (f *Font) DrawText(text string, x, y int, dst *image.RGBA) {
	f.updateFace()
	metrics := f.face.Metrics()
	py := y + metrics.Ascent.Round()
	d := font.Drawer{Dst: dst, Src: &f.color, Face: f.face, Dot: fixed.P(x, py)}
	d.DrawString(text)
}
