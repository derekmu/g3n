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
)

// Label is a text only UI element.
type Label struct {
	Panel
	font           *text.Font
	text           string
	canvas         *text.Canvas
	color          math32.Color
	fontAttributes text.FontAttributes
}

// NewLabel creates a Label with the specified text using the default font.
func NewLabel(txt string) *Label {
	return NewLabelWithFont(txt, StyleDefault().Font)
}

// NewLabelWithFont creates a Label with the specified text using the specified font.
func NewLabelWithFont(txt string, fnt *text.Font) *Label {
	l := new(Label)
	l.InitLabel(txt, fnt)
	return l
}

// InitLabel initializes this Label.
func (l *Label) InitLabel(txt string, fnt *text.Font) {
	l.InitPanel(l, 0, 0)
	l.SetResizeToTexture(true)
	l.font = fnt
	l.color = math32.Color3{R: 1, G: 1, B: 1}
	l.fontAttributes = text.FontAttributes{
		PointSize: 14,
		DPI:       72,
		Hinting:   font.HintingFull,
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
func (l *Label) SetColor(color math32.Color) {
	if l.color != color {
		l.color = color
		l.drawText()
	}
}

// Color returns the text color.
func (l *Label) Color() math32.Color {
	return l.color
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
func (l *Label) SetFontSize(size int32) {
	if l.fontAttributes.PointSize != size {
		l.fontAttributes.PointSize = size
		l.drawText()
	}
}

// FontSize returns the point size of the font.
func (l *Label) FontSize() int32 {
	return l.fontAttributes.PointSize
}

// SetFontDPI sets the resolution of the font in dots per inch (DPI).
func (l *Label) SetFontDPI(dpi int32) {
	if l.fontAttributes.DPI != dpi {
		l.fontAttributes.DPI = dpi
		l.drawText()
	}
}

// FontDPI returns the resolution of the font in dots per inch (DPI).
func (l *Label) FontDPI() int32 {
	return l.fontAttributes.DPI
}

// SetTextColor updates both the text and color and redraws the label.
//
// This reduces redraws compared to changing the text and color separately.
func (l *Label) SetTextColor(txt string, color math32.Color) {
	if txt != l.text || l.color != color {
		l.text = txt
		l.color = color
		l.drawText()
	}
}

// drawText redraws the label texture.
func (l *Label) drawText() {
	// Update the canvas
	width, height := l.font.MeasureText(l.text)
	if l.canvas == nil {
		l.canvas = text.NewCanvas(width, height)
	} else {
		l.canvas.Resize(width, height)
	}
	// Fill with text color with alpha zero to reduce blending artifacts
	bgColor := l.color.ToColor4()
	bgColor.A = 0
	l.canvas.Fill(bgColor.ToRGBA())
	// Draw the text
	l.font.SetAttributes(l.fontAttributes)
	l.font.SetColor(l.color)
	l.canvas.DrawText(0, 0, l.text, l.font)
	// Update texture
	tex := l.texture
	if tex == nil {
		tex = texture.NewTexture2DFromRGBA(&l.canvas.RGBA)
		tex.SetMagFilter(gls.NEAREST)
		tex.SetMinFilter(gls.NEAREST)
	} else {
		tex.SetFromRGBA(&l.canvas.RGBA)
	}
	l.SetTexture(tex)
}
