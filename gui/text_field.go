package gui

import (
	"image"
	"slices"
	"time"
	"unicode/utf8"

	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/text"
	"github.com/derekmu/g3n/texture"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// TextField is an editable text GUI element.
type TextField struct {
	Panel
	label          *Panel
	mouse          MouseTracker
	font           *text.Font
	text           []rune
	textWidths     []fixed.Int26_6
	canvas         *text.Canvas
	color          math32.Color
	fontAttributes text.FontAttributes
	focused        bool
	cursorOwner    bool
	update         bool
	caret          struct {
		index     int
		speed     time.Duration
		rect      image.Rectangle
		state     bool
		blinkTime time.Time
	}
}

// NewTextField creates a TextField using the default font.
func NewTextField(txt string) *TextField {
	return NewTextFieldWithFont(txt, StyleDefault().Font)
}

// NewTextFieldWithFont creates a TextField with text and a font.
func NewTextFieldWithFont(txt string, fnt *text.Font) *TextField {
	l := new(TextField)
	l.InitTextField(txt, fnt)
	return l
}

// InitTextField initializes the TextField.
func (f *TextField) InitTextField(txt string, fnt *text.Font) {
	f.InitPanel(f, 0, 0)
	f.label = NewPanel(0, 0)
	f.label.SetEnabled(false)
	f.font = fnt
	f.color = math32.Color3{R: 1, G: 1, B: 1}
	f.fontAttributes = text.FontAttributes{
		PointSize: 14,
		DPI:       72,
		Hinting:   font.HintingFull,
	}
	f.SetPaddings(RectBounds{
		Top:    2,
		Right:  2,
		Bottom: 2,
		Left:   2,
	})
	f.SetText(txt)
	f.SetCaretIndex(len(txt))
	f.SetCaretSpeed(time.Second / 2)
	f.Add(f.label)
	f.mouse.Init(f)
	f.Subscribe(f.onGuiEvent)
}

// SetText sets the field's text.
func (f *TextField) SetText(txt string) {
	l := utf8.RuneCountInString(txt)
	update := l != len(f.text)
	if cap(f.text) < l {
		// Capacity isn't enough, reallocate
		f.text = make([]rune, l)
	} else if len(f.text) < l {
		// Longer, expand
		f.text = f.text[:l]
	} else if len(f.text) > l {
		// Shorter, trim
		f.text = f.text[:l]
	}
	i := 0
	for _, r := range txt {
		update = update || r != f.text[i]
		f.text[i] = r
		i++
	}
	if update {
		f.update = true
	}
}

// Text returns the field's text.
func (f *TextField) Text() string {
	return string(f.text)
}

// SetCaretIndex sets the caret index.
func (f *TextField) SetCaretIndex(index int) {
	index = max(0, min(len(f.text), index))
	f.update = f.update || index != f.caret.index
	f.caret.index = index
}

// CaretIndex returns the caret index.
func (f *TextField) CaretIndex() int {
	return f.caret.index
}

// SetCaretSpeed sets speed that the caret blinks.
func (f *TextField) SetCaretSpeed(speed time.Duration) {
	f.caret.speed = speed
}

// CaretSpeed is the speed that the caret blinks.
func (f *TextField) CaretSpeed() time.Duration {
	return f.caret.speed
}

// SetColor sets the text color.
func (f *TextField) SetColor(color math32.Color) {
	f.update = f.update || f.color != color
	f.color = color
}

// Color returns the text color.
func (f *TextField) Color() math32.Color {
	return f.color
}

// SetFont sets the font.
func (f *TextField) SetFont(fnt *text.Font) {
	f.update = f.update || f.font != fnt
	f.font = fnt
}

// Font returns the font.
func (f *TextField) Font() *text.Font {
	return f.font
}

// SetFontSize sets the point size of the font.
func (f *TextField) SetFontSize(size int32) {
	f.update = f.update || f.fontAttributes.PointSize != size
	f.fontAttributes.PointSize = size
}

// FontSize returns the point size of the font.
func (f *TextField) FontSize() int32 {
	return f.fontAttributes.PointSize
}

// RenderSetup updates the texture before rendering.
func (f *TextField) RenderSetup(gl *gls.GLS, ri *core.RenderInfo) {
	f.drawText()
	now := time.Now()
	if now.After(f.caret.blinkTime) {
		c := f.color.ToRGBA()
		if f.caret.state {
			c.A = 0
		}
		f.canvas.DrawRectangle(f.caret.rect, c)
		f.caret.state = !f.caret.state
		f.caret.blinkTime = now.Add(f.caret.speed)
		f.label.texture.SetFromRGBA(&f.canvas.RGBA)
	}
	f.Panel.RenderSetup(gl, ri)
}

// CaretWidth is the pixel width of the caret.
func (f *TextField) CaretWidth() int {
	if f.focused {
		return max(1, int(f.fontAttributes.PointSize)/40)
	}
	return 0
}

// measureText returns the caret width, width, and height of the text.
// It also updates textWidths
func (f *TextField) measureText() (cwidth, width, height int) {
	f.font.SetAttributes(f.fontAttributes)
	if cap(f.textWidths) < len(f.text) {
		f.textWidths = make([]fixed.Int26_6, len(f.text))
	} else {
		f.textWidths = f.textWidths[:len(f.text)]
	}
	width, height = f.font.MeasureRunes(f.text, f.textWidths)
	cwidth = f.CaretWidth()
	return cwidth, width + cwidth, height
}

// FitToText updates the size of the field to match the text.
func (f *TextField) FitToText() {
	_, width, height := f.measureText()
	f.SetContentSize(width, height)
}

// drawText redraws the field texture if needed.
func (f *TextField) drawText() {
	// Don't do anything if nothing has changed
	if !f.update {
		return
	}

	// Update the canvas size
	cwidth, width, height := f.measureText()
	if f.canvas == nil {
		f.canvas = text.NewCanvas(width, height)
	} else {
		f.canvas.Resize(width, height)
	}

	// Fill with text color with alpha zero to reduce blending artifacts
	bgColor := f.color.ToColor4()
	bgColor.A = 0
	f.canvas.Fill(bgColor.ToRGBA())

	// Draw the text
	f.font.SetColor(f.color)
	f.canvas.DrawRunes(0, 0, f.text, f.font)

	// Update the cursor
	lwidthf := fixed.I(0)
	for i := range f.caret.index {
		lwidthf += f.textWidths[i]
	}
	lwidth := lwidthf.Floor()
	f.caret.rect = image.Rect(lwidth, cwidth, lwidth+cwidth, height-cwidth)
	f.caret.state = true
	f.caret.blinkTime = time.Now().Add(f.caret.speed)
	f.canvas.DrawRectangle(f.caret.rect, f.color.ToRGBA())

	// Update texture
	tex := f.label.texture
	if tex == nil {
		tex = texture.NewTexture2DFromRGBA(&f.canvas.RGBA)
		tex.SetMagFilter(gls.NEAREST)
		tex.SetMinFilter(gls.NEAREST)
	} else {
		tex.SetFromRGBA(&f.canvas.RGBA)
	}
	f.label.SetTexture(tex)
	f.label.SetSize(tex.Width(), tex.Height())
	// TODO: position label such that the caret is visible

	f.update = false
}

func (f *TextField) onGuiEvent(event core.GuiEvent) bool {
	switch ev := event.(type) {
	case core.GuiClickEvent:
		sum := fixed.I(0)
		x := int(ev.X) - int(f.Position().X) - int(f.label.Position().X) - f.paddings.Left
		ci := -1
		for i, w := range f.textWidths {
			sum += w
			if sum.Ceil() >= x {
				ci = i
				break
			}
		}
		if ci == -1 {
			ci = len(f.text)
		}
		f.SetCaretIndex(ci)
		GetManager().SetKeyFocus(f)
	case core.GuiMouseEnterEvent:
		f.updateMouseIcon()
	case core.MouseEvent:
		// TODO: update caret range if mouse is down
	case core.GuiMouseLeaveEvent:
		f.updateMouseIcon()
	case core.GuiEnableEvent:
		f.updateMouseIcon()
		// TODO: change color when disabled?
	case core.KeyDownEvent:
		f.onKeyEvent(ev.Key, ev.Mods)
	case core.KeyRepeatEvent:
		f.onKeyEvent(ev.Key, ev.Mods)
	case core.CharEvent:
		if f.Enabled() {
			if cap(f.text) > len(f.text) {
				f.text = slices.Insert(f.text, f.caret.index, ev.Char)
			} else {
				newText := slices.Insert(f.text, f.caret.index, ev.Char)
				f.text = newText
			}
			f.SetCaretIndex(f.caret.index + 1)
			f.update = true
		}
	case core.GuiFocusEvent:
		f.focused = true
		f.update = true
	case core.GuiFocusLostEvent:
		f.focused = false
		f.update = true
	default:
		return false
	}
	return true
}

func (f *TextField) updateMouseIcon() {
	if f.mouse.Over && f.Enabled() {
		_ = GetManager().window.SetCursorIcon(core.CursorIBeam)
		f.cursorOwner = true
	} else if f.cursorOwner {
		_ = GetManager().window.SetCursorIcon(core.CursorArrow)
		f.cursorOwner = false
	}
}

func (f *TextField) onKeyEvent(key core.Key, _ core.ModifierKey) {
	if !f.Enabled() {
		return
	}
	switch key {
	case core.KeyLeft:
		f.SetCaretIndex(f.caret.index - 1)
	case core.KeyRight:
		f.SetCaretIndex(f.caret.index + 1)
	case core.KeyEscape:
		if f.focused {
			GetManager().SetKeyFocus(nil)
		}
	case core.KeyBackspace:
		if f.caret.index > 0 {
			copy(f.text[f.caret.index-1:], f.text[f.caret.index:])
			f.text = f.text[:len(f.text)-1]
			f.SetCaretIndex(f.caret.index - 1)
			f.update = true
		}
	case core.KeyDelete:
		if f.caret.index < len(f.text) {
			copy(f.text[f.caret.index:], f.text[f.caret.index+1:])
			f.text = f.text[:len(f.text)-1]
			f.update = true
		}
	case core.KeyHome, core.KeyPageUp:
		f.SetCaretIndex(0)
	case core.KeyEnd, core.KeyPageDown:
		f.SetCaretIndex(len(f.text))
	}
	// TODO: select all
	// TODO: undo/redo
	// TODO: copy/paste
}
