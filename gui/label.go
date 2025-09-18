// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
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
	redraw         bool
}

// NewLabel creates a Label using the default font.
func NewLabel(txt string) *Label {
	return NewLabelWithFont(txt, StyleDefault().Font)
}

// NewLabelWithFont creates a Label.
func NewLabelWithFont(txt string, fnt *text.Font) *Label {
	l := new(Label)
	l.InitLabel(txt, fnt)
	return l
}

// InitLabel initializes the Label.
func (l *Label) InitLabel(txt string, fnt *text.Font) {
	l.InitPanel(l, 0, 0)
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
	l.redraw = l.redraw || txt != l.text
	l.text = txt
}

// Text returns the label text.
func (l *Label) Text() string {
	return l.text
}

// SetColor sets the text color.
func (l *Label) SetColor(color math32.Color) {
	l.redraw = l.redraw || l.color != color
	l.color = color
}

// Color returns the text color.
func (l *Label) Color() math32.Color {
	return l.color
}

// SetFont sets the font.
func (l *Label) SetFont(f *text.Font) {
	l.redraw = l.redraw || l.font != f
	l.font = f
}

// Font returns the font.
func (l *Label) Font() *text.Font {
	return l.font
}

// SetFontSize sets the point size of the font.
func (l *Label) SetFontSize(size int32) {
	l.redraw = l.redraw || l.fontAttributes.PointSize != size
	l.fontAttributes.PointSize = size
}

// FontSize returns the point size of the font.
func (l *Label) FontSize() int32 {
	return l.fontAttributes.PointSize
}

// RenderSetup updates the texture before rendering.
func (l *Label) RenderSetup(gl *gls.GLS, ri *core.RenderInfo) {
	l.DrawText()
	l.Panel.RenderSetup(gl, ri)
}

// MeasureText returns the width and height of the text.
func (l *Label) MeasureText() (width, height int) {
	l.font.SetAttributes(l.fontAttributes)
	return l.font.MeasureText(l.text)
}

// FitSizeToText updates the size of the label to match the text.
func (l *Label) FitSizeToText() {
	w, h := l.MeasureText()
	l.SetSize(w, h)
}

// DrawText redraws the label texture if needed.
func (l *Label) DrawText() {
	// Don't do anything if nothing has changed
	if !l.redraw {
		return
	}

	// Update the canvas size
	width, height := l.MeasureText()
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

	l.redraw = false
}
