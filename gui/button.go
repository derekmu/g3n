// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
)

// Button is a button UI element.
type Button struct {
	Panel                   // Embedded Panel
	Label     *Label        // Label panel
	image     *Image        // pointer to button image (may be nil)
	styles    *ButtonStyles // pointer to current button styles
	mouseOver bool          // true if mouse is over button
	pressed   bool          // true if button is pressed
}

// ButtonStyle contains the styling of a Button.
type ButtonStyle BasicStyle

// ButtonStyles contains one ButtonStyle for each ButtonState.
type ButtonStyles struct {
	Normal   ButtonStyle
	Over     ButtonStyle
	Focus    ButtonStyle
	Pressed  ButtonStyle
	Disabled ButtonStyle
}

// ButtonState identifies the state of a Button.
type ButtonState int

const (
	ButtonNormal ButtonState = iota
	ButtonOver
	ButtonPressed
	ButtonDisabled
	// ButtonFocus
)

// NewButton creates a new Button with the specified text for the button label.
func NewButton(text string) *Button {
	b := new(Button)
	b.styles = &StyleDefault().Button

	// Initializes the button panel
	b.Panel.InitPanel(b, 0, 0)

	// Subscribe to panel events
	b.Subscribe(OnKeyDown, b.onKey)
	b.Subscribe(OnKeyUp, b.onKey)
	b.Subscribe(OnMouseUp, b.onMouse)
	b.Subscribe(OnMouseDown, b.onMouse)
	b.Subscribe(OnMouseUpOut, b.onMouse)
	b.Subscribe(OnCursor, b.onCursor)
	b.Subscribe(OnCursorEnter, b.onCursor)
	b.Subscribe(OnCursorLeave, b.onCursor)
	b.Subscribe(OnEnable, func(name string, ev interface{}) { b.update() })
	b.Subscribe(OnResize, func(name string, ev interface{}) { b.recalc() })

	// Creates label
	b.Label = NewLabel(text)
	b.Label.Subscribe(OnResize, func(name string, ev interface{}) { b.recalc() })
	b.Panel.Add(b.Label)

	b.recalc() // recalc first then update!
	b.update()
	return b
}

// SetImage sets the button image from the specified filename.
func (b *Button) SetImage(imgfile string) error {
	img, err := NewImage(imgfile)
	if err != nil {
		return err
	}
	if b.image != nil {
		b.Panel.Remove(b.image)
	}
	b.image = img
	b.Panel.Add(b.image)
	b.recalc()
	return nil
}

// SetStyles set the button styles overriding the default style.
func (b *Button) SetStyles(bs *ButtonStyles) {
	b.styles = bs
	b.update()
}

// onCursor processes subscribed cursor events.
func (b *Button) onCursor(evname string, _ interface{}) {
	switch evname {
	case OnCursorEnter:
		b.mouseOver = true
		b.update()
	case OnCursorLeave:
		b.mouseOver = false
		b.update()
	}
}

// onMouse processes subscribed mouse events.
func (b *Button) onMouse(evname string, _ interface{}) {
	if !b.Enabled() {
		return
	}
	switch evname {
	case OnMouseDown:
		GetManager().SetKeyFocus(b)
		b.pressed = true
		b.update()
	case OnMouseUpOut:
		fallthrough
	case OnMouseUp:
		if b.pressed && b.mouseOver {
			b.Dispatch(OnClick, nil)
		}
		b.pressed = false
		b.update()
	default:
		return
	}
}

// onKey processes subscribed key events.
func (b *Button) onKey(evname string, ev interface{}) {
	kev := ev.(*core.KeyEvent)
	if kev.Key != core.KeyEnter {
		return
	}
	switch evname {
	case OnKeyDown:
		b.pressed = true
		b.update()
		b.Dispatch(OnClick, nil)
	case OnKeyUp:
		b.pressed = false
		b.update()
	}
}

// update updates the button visual state.
func (b *Button) update() {
	if !b.Enabled() {
		b.applyStyle(&b.styles.Disabled)
		return
	}
	if b.pressed && b.mouseOver {
		b.applyStyle(&b.styles.Pressed)
		return
	}
	if b.mouseOver {
		b.applyStyle(&b.styles.Over)
		return
	}
	b.applyStyle(&b.styles.Normal)
}

// applyStyle applies the specified button style.
func (b *Button) applyStyle(bs *ButtonStyle) {
	b.Panel.ApplyStyle(&bs.PanelStyle)
	b.Label.SetColor(bs.FgColor)
}

// recalc recalculates all dimensions and position from inside out.
func (b *Button) recalc() {
	// Current width and height of button content area
	width := b.Panel.ContentWidth()
	height := b.Panel.ContentHeight()

	// Image width
	imageWidth := float32(0)
	spacing := float32(4)
	if b.image != nil {
		imageWidth = b.image.Width()
	}
	if imageWidth == 0 {
		spacing = 0
	}

	// Label width
	labelWidth := spacing + b.Label.Width()
	// If the label is empty and an image was defined, ignore the label to centralize the image
	if b.Label.Text() == "" && imageWidth > 0 {
		labelWidth = 0
	}

	// Sets new content width and height if necessary
	minWidth := imageWidth + labelWidth
	minHeight := b.Label.Height()
	resize := false
	if width < minWidth {
		width = minWidth
		resize = true
	}
	if height < minHeight {
		height = minHeight
		resize = true
	}
	if resize {
		b.SetContentSize(width, height)
	}

	// Centralize horizontally
	px := (width - minWidth) / 2

	// Set label position
	ly := (height - b.Label.Height()) / 2
	b.Label.SetPosition(px+imageWidth+spacing, ly)

	// Image position
	if b.image != nil {
		iy := (height - b.image.height) / 2
		b.image.SetPosition(px, iy)
	}
}
