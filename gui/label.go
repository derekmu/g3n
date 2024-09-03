// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/text"
	"github.com/derekmu/g3n/texture"
)

// Label is a text only UI element.
type Label struct {
	Panel                    // Embedded Panel
	font  *text.Font         // TrueType font face
	tex   *texture.Texture2D // Texture with text
	style *LabelStyle        // The style of the panel and font attributes
	text  string             // Text being displayed
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
	l.Panel.InitPanel(l, 0, 0)
	l.Panel.mat.SetTransparent(true)
	if font != StyleDefault().FontIcon {
		l.Panel.SetPaddings(2, 0, 2, 0)
	}
	styleCopy := StyleDefault().Label
	l.style = &styleCopy
	l.SetText(text)
}

// SetText sets and redraws the label text.
func (l *Label) SetText(text string) {
	l.text = text
	// Need at least one character to get dimensions
	if text == "" {
		text = " "
	}

	// Set font properties
	l.font.SetAttributes(&l.style.FontAttributes)
	l.font.SetColor(l.style.FgColor)

	scaleX, scaleY := GetManager().window.GetScale()
	l.font.SetScaleXY(scaleX, scaleY)

	// Create an image with the text
	textImage := l.font.DrawText(text)

	// Create texture if it doesn't exist yet
	if l.tex == nil {
		l.tex = texture.NewTexture2DFromRGBA(textImage)
		l.tex.SetMagFilter(gls.NEAREST)
		l.tex.SetMinFilter(gls.NEAREST)
		l.Panel.Material().AddTexture(l.tex)
		// Otherwise update texture with new image
	} else {
		l.tex.SetFromRGBA(textImage)
	}

	// Update label panel dimensions
	width, height := float32(textImage.Rect.Dx()), float32(textImage.Rect.Dy())
	// since we enlarged the font texture for higher quality, we have to scale it back to its original point size
	width, height = width/float32(scaleX), height/float32(scaleY)
	l.Panel.SetContentSize(width, height)
}

// Text returns the label text.
func (l *Label) Text() string {
	return l.text
}

// SetColor sets the text color.
func (l *Label) SetColor(color math32.Color4) *Label {
	l.style.FgColor = color
	l.SetText(l.text)
	return l
}

// Color returns the text color.
func (l *Label) Color() math32.Color4 {
	return l.style.FgColor
}

// SetBgColor sets the background color.
func (l *Label) SetBgColor(color math32.Color4) *Label {
	l.style.BgColor = color
	l.Panel.SetColor(color)
	l.SetText(l.text)
	return l
}

// BgColor returns the background color.
func (l *Label) BgColor() math32.Color4 {
	return l.style.BgColor
}

// SetFont sets the font.
func (l *Label) SetFont(f *text.Font) {
	l.font = f
	l.SetText(l.text)
}

// Font returns the font.
func (l *Label) Font() *text.Font {
	return l.font
}

// SetFontSize sets the point size of the font.
func (l *Label) SetFontSize(size float64) *Label {
	l.style.PointSize = size
	l.SetText(l.text)
	return l
}

// FontSize returns the point size of the font.
func (l *Label) FontSize() float64 {
	return l.style.PointSize
}

// SetFontDPI sets the resolution of the font in dots per inch (DPI).
func (l *Label) SetFontDPI(dpi float64) *Label {
	l.style.DPI = dpi
	l.SetText(l.text)
	return l
}

// FontDPI returns the resolution of the font in dots per inch (DPI).
func (l *Label) FontDPI() float64 {
	return l.style.DPI
}

// SetLineSpacing sets the spacing between lines.
func (l *Label) SetLineSpacing(spacing float64) *Label {
	l.style.LineSpacing = spacing
	l.SetText(l.text)
	return l
}

// LineSpacing returns the spacing between lines.
func (l *Label) LineSpacing() float64 {
	return l.style.LineSpacing
}

// setTextCaret sets the label text and draws a caret at the specified line and column.
// It is normally used by the Edit widget.
func (l *Label) setTextCaret(msg string, mx, width int, drawCaret bool, line, col, selStart, selEnd int) {
	// Set font properties
	l.font.SetAttributes(&l.style.FontAttributes)
	l.font.SetColor(l.style.FgColor)

	scaleX, scaleY := GetManager().window.GetScale()
	l.font.SetScaleXY(scaleX, scaleY)

	// Create canvas and draw text
	_, height := l.font.MeasureText(msg)
	canvas := text.NewCanvas(width, height, l.style.BgColor)
	canvas.DrawTextCaret(mx, 0, msg, l.font, drawCaret, line, col, selStart, selEnd)

	// Creates texture if if doesnt exist.
	if l.tex == nil {
		l.tex = texture.NewTexture2DFromRGBA(canvas.RGBA)
		l.Panel.Material().AddTexture(l.tex)
		// Otherwise update texture with new image
	} else {
		l.tex.SetFromRGBA(canvas.RGBA)
	}
	// Set texture filtering parameters for text
	l.tex.SetMagFilter(gls.NEAREST)
	l.tex.SetMinFilter(gls.NEAREST)

	// Updates label panel dimensions
	l.Panel.SetContentSize(float32(width)/float32(scaleX), float32(height)/float32(scaleY))
	l.text = msg
}
