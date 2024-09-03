// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/text"
	"github.com/derekmu/g3n/texture"
	"image"
)

// Label is a text only UI element.
type Label struct {
	Image
	font   *text.Font
	style  LabelStyle
	text   string
	rgba   *image.RGBA
	canvas *text.Canvas
}

// LabelStyle contains the styling attributes of a Label.
type LabelStyle struct {
	PanelStyle
	text.FontAttributes
	FgColor math32.Color4
}

// NewLabel creates a Label with the specified text using the default font.
func NewLabel(text string) *Label {
	return NewLabelWithFont(text, StyleDefault().Font)
}

// NewIconLabel creates a Label with the specified text using the default icon font.
func NewIconLabel(text string) *Label {
	return NewLabelWithFont(text, StyleDefault().FontIcon)
}

// NewLabelWithFont creates a Label with the specified text using the specified font.
func NewLabelWithFont(text string, font *text.Font) *Label {
	l := new(Label)
	l.InitLabel(text, font)
	return l
}

// InitLabel initializes this Label.
func (l *Label) InitLabel(text string, font *text.Font) {
	l.font = font
	l.Image.InitImage()
	l.style = StyleDefault().Label
	l.SetText(text)
}

// SetText sets and redraws the label text.
func (l *Label) SetText(text string) {
	// Need at least one character to get dimensions
	if text == "" {
		text = " "
	}
	if text != l.text {
		l.text = text
		l.drawText()
	}
}

// Text returns the label text.
func (l *Label) Text() string {
	return l.text
}

// SetColor sets the text color.
func (l *Label) SetColor(color math32.Color4) {
	if l.style.FgColor != color {
		l.style.FgColor = color
		l.drawText()
	}
}

// Color returns the text color.
func (l *Label) Color() math32.Color4 {
	return l.style.FgColor
}

// SetBgColor sets the background color.
func (l *Label) SetBgColor(color math32.Color4) {
	if l.style.BgColor != color {
		l.style.BgColor = color
		l.Panel.SetColor(color)
		l.drawText()
	}
}

// BgColor returns the background color.
func (l *Label) BgColor() math32.Color4 {
	return l.style.BgColor
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
	if size != l.style.PointSize {
		l.style.PointSize = size
		l.drawText()
	}
}

// FontSize returns the point size of the font.
func (l *Label) FontSize() float64 {
	return l.style.PointSize
}

// SetFontDPI sets the resolution of the font in dots per inch (DPI).
func (l *Label) SetFontDPI(dpi float64) {
	if dpi != l.style.DPI {
		l.style.DPI = dpi
		l.drawText()
	}
}

// FontDPI returns the resolution of the font in dots per inch (DPI).
func (l *Label) FontDPI() float64 {
	return l.style.DPI
}

// SetLineSpacing sets the spacing between lines.
func (l *Label) SetLineSpacing(spacing float64) {
	if spacing != l.style.LineSpacing {
		l.style.LineSpacing = spacing
		l.drawText()
	}
}

// LineSpacing returns the spacing between lines.
func (l *Label) LineSpacing() float64 {
	return l.style.LineSpacing
}

// drawText redraws the label texture.
func (l *Label) drawText() {
	// Set font properties
	l.font.SetAttributes(&l.style.FontAttributes)
	l.font.SetColor(l.style.FgColor)

	scaleX, scaleY := GetManager().window.GetScale()
	l.font.SetScaleXY(scaleX, scaleY)

	// Create an image with the text
	width, height := l.font.MeasureText(l.text)
	if l.canvas == nil || l.rgba.Rect.Dx() < width || l.rgba.Rect.Dy() < height {
		// Allocate a new canvas if the existing one can't hold the text
		l.canvas = text.NewCanvas(width, height, l.style.BgColor)
		// Keep track of the full RGBA
		l.rgba = l.canvas.RGBA
	} else {
		// Reuse part of the already allocated image
		l.canvas.RGBA = l.rgba.SubImage(image.Rect(0, 0, width, height)).(*image.RGBA)
		// Update the color
		l.canvas.BgColor = l.style.BgColor
	}
	l.canvas.DrawText(0, 0, l.text, l.font)

	if l.tex == nil {
		// Create texture if it doesn't exist yet
		l.tex = texture.NewTexture2DFromRGBA(l.canvas.RGBA)
		l.tex.SetMagFilter(gls.NEAREST)
		l.tex.SetMinFilter(gls.NEAREST)
		l.Panel.Material().AddTexture(l.tex)
	} else {
		// Otherwise update texture with new image
		l.tex.SetFromRGBA(l.canvas.RGBA)
	}

	// Change the image texture
	l.SetTexture(l.tex)
}
