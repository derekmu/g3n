// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/text"
	"github.com/derekmu/g3n/texture"
	"golang.org/x/image/font"
	"image"
)

// Label is a text only UI element.
type Label struct {
	Image
	font           *text.Font
	text           string
	rgba           image.RGBA
	canvas         *text.Canvas
	color          math32.Color4
	bgColor        math32.Color4
	fontAttributes text.FontAttributes
}

// NewLabel creates a Label with the specified text using the default font.
func NewLabel(txt string) *Label {
	return NewLabelWithFont(txt, StyleDefault().Font)
}

// NewIconLabel creates a Label with the specified text using the default icon font.
func NewIconLabel(txt string) *Label {
	return NewLabelWithFont(txt, StyleDefault().FontIcon)
}

// NewLabelWithFont creates a Label with the specified text using the specified font.
func NewLabelWithFont(txt string, fnt *text.Font) *Label {
	l := new(Label)
	l.InitLabel(txt, fnt)
	return l
}

// InitLabel initializes this Label.
func (l *Label) InitLabel(txt string, fnt *text.Font) {
	l.InitImage()
	l.font = fnt
	l.color = math32.Color4{1, 1, 1, 1}
	l.bgColor = math32.Color4{1, 1, 1, 0}
	l.fontAttributes = text.FontAttributes{
		PointSize:   14,
		DPI:         72,
		LineSpacing: 1.0,
		Hinting:     font.HintingNone,
	}
	l.SetText(txt)
}

// SetText sets and redraws the label text.
func (l *Label) SetText(txt string) {
	if txt != l.text {
		l.text = txt
		l.drawText()
	}
}

// Text returns the label text.
func (l *Label) Text() string {
	return l.text
}

// SetColor sets the text color.
func (l *Label) SetColor(color math32.Color4) {
	if l.color != color {
		l.color = color
		l.drawText()
	}
}

// Color returns the text color.
func (l *Label) Color() math32.Color4 {
	return l.color
}

// SetBgColor sets the background color.
func (l *Label) SetBgColor(color math32.Color4) {
	if l.bgColor != color {
		l.bgColor = color
		l.Panel.SetColor(color)
		l.drawText()
	}
}

// BgColor returns the background color.
func (l *Label) BgColor() math32.Color4 {
	return l.bgColor
}

// SetFont sets the font.
func (l *Label) SetFont(f *text.Font) {
	if l.font != f {
		l.font = f
		l.drawText()
	}
}

// Font returns the font.
func (l *Label) Font() *text.Font {
	return l.font
}

// SetFontSize sets the point size of the font.
func (l *Label) SetFontSize(size float64) {
	if l.fontAttributes.PointSize != size {
		l.fontAttributes.PointSize = size
		l.drawText()
	}
}

// FontSize returns the point size of the font.
func (l *Label) FontSize() float64 {
	return l.fontAttributes.PointSize
}

// SetFontDPI sets the resolution of the font in dots per inch (DPI).
func (l *Label) SetFontDPI(dpi float64) {
	if l.fontAttributes.DPI != dpi {
		l.fontAttributes.DPI = dpi
		l.drawText()
	}
}

// FontDPI returns the resolution of the font in dots per inch (DPI).
func (l *Label) FontDPI() float64 {
	return l.fontAttributes.DPI
}

// SetLineSpacing sets the spacing between lines.
func (l *Label) SetLineSpacing(spacing float64) {
	if l.fontAttributes.LineSpacing != spacing {
		l.fontAttributes.LineSpacing = spacing
		l.drawText()
	}
}

// LineSpacing returns the spacing between lines.
func (l *Label) LineSpacing() float64 {
	return l.fontAttributes.LineSpacing
}

// drawText redraws the label texture.
func (l *Label) drawText() {
	// Set font properties
	l.font.SetAttributes(&l.fontAttributes)
	l.font.SetColor(l.color)

	scaleX, scaleY := GetManager().window.GetScale()
	l.font.SetScaleXY(scaleX, scaleY)

	// Create an image with the text
	txt := l.text
	// Need at least one character to get dimensions
	if txt == "" {
		txt = " "
	}
	width, height := l.font.MeasureText(txt)
	if l.canvas == nil || l.rgba.Rect.Dx() < width || l.rgba.Rect.Dy() < height {
		// Allocate a new canvas if the existing one can't hold the text
		l.canvas = text.NewCanvas(width, height, l.bgColor)
		// Keep a copy of the RGBA
		l.rgba = *l.canvas.RGBA
	} else {
		// Reuse part of the already allocated image
		l.canvas.RGBA.Pix = l.rgba.Pix[:4*width*height]
		l.canvas.RGBA.Stride = 4 * width
		l.canvas.RGBA.Rect = image.Rect(0, 0, width, height)
		// Update the color
		l.canvas.BgColor = l.bgColor
	}
	l.canvas.DrawText(0, 0, txt, l.font)

	if l.tex == nil {
		// Create texture if it doesn't exist yet
		l.tex = texture.NewTexture2DFromRGBA(l.canvas.RGBA)
		l.tex.SetMagFilter(gls.NEAREST)
		l.tex.SetMinFilter(gls.NEAREST)
		l.Panel.Material().AddTexture(l.tex)
	} else {
		// Otherwise updateTexture texture with new image
		l.tex.SetFromRGBA(l.canvas.RGBA)
	}

	// Change the image texture
	l.SetTexture(l.tex)
}
