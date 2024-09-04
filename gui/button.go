// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/texture"
)

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

// Button is a button UI element that extends Image and uses different textures for each ButtonState.
// There is also a Label for text on top.
type Button struct {
	Image
	Label         *Label
	expandToLabel bool
	styles        *ButtonStyles
	mouseOver     bool
	pressed       bool
	textures      [ButtonDisabled + 1]*texture.Texture2D
}

// NewButton creates a new Button with the specified text the button label.
func NewButton(text string) *Button {
	b := new(Button)
	b.InitButton(text)
	return b
}

// InitButton initializes the image and subscribes to events.
func (b *Button) InitButton(text string) {
	b.InitImage()

	b.styles = &StyleDefault().Button

	// Create label
	b.Label = NewLabel(text)
	b.Add(b.Label)
	b.Label.Subscribe(OnResize, func(string, interface{}) { b.recalculateSize() })
	b.expandToLabel = true

	// subscribe to events
	b.Subscribe(OnMouseUp, b.onMouse)
	b.Subscribe(OnMouseDown, b.onMouse)
	b.Subscribe(OnCursorEnter, b.onCursor)
	b.Subscribe(OnCursorLeave, b.onCursor)
	b.Subscribe(OnEnable, b.onEnable)
	b.Subscribe(OnEnable, func(name string, ev interface{}) { b.updateStyle() })
	b.Subscribe(OnResize, func(name string, ev interface{}) { b.recalculateSize() })

	b.updateTexture()
	b.updateStyle()
}

// Dispose disposes of the label and all button textures.
func (b *Button) Dispose() {
	b.Image.Dispose()
	b.Label.Dispose()
	for _, tex := range b.textures {
		if tex != nil {
			tex.Dispose()
		}
	}
}

// onCursor handles cursor enter and leave events.
func (b *Button) onCursor(evname string, _ interface{}) {
	switch evname {
	case OnCursorEnter:
		b.mouseOver = true
		b.updateTexture()
		b.updateStyle()
	case OnCursorLeave:
		b.mouseOver = false
		// Pressing, dragging out, and releasing cancels the click
		b.pressed = false
		b.updateTexture()
		b.updateStyle()
	}
}

// onMouse handles mouse down and up events.
func (b *Button) onMouse(evname string, ev interface{}) {
	if !b.Enabled() {
		return
	}
	mev := ev.(*core.MouseEvent)
	switch evname {
	case OnMouseDown:
		if mev.Button == core.MouseButtonLeft {
			b.pressed = true
			b.updateTexture()
			b.updateStyle()
		}
	case OnMouseUp:
		if mev.Button == core.MouseButtonLeft {
			if b.pressed && b.mouseOver {
				b.Dispatch(OnClick, nil)
			}
			b.pressed = false
			b.updateTexture()
			b.updateStyle()
		}
	}
}

// onEnable handles enable and disable events.
func (b *Button) onEnable(evname string, _ interface{}) {
	switch evname {
	case OnEnable:
		// Enabling or disabling a button cancels the click
		b.pressed = false
		b.updateTexture()
		b.updateStyle()
	}
}

// SetStateTexture changes the texture used by the button in a given state.
// Any prior texture for the state is disposed.
func (b *Button) SetStateTexture(state ButtonState, tex *texture.Texture2D) {
	if b.textures[state] != nil {
		b.textures[state].Dispose()
	}
	b.textures[state] = tex
	b.updateTexture()
}

// GetStateTexture returns the texture used by the button in a given state.
func (b *Button) GetStateTexture(state ButtonState) *texture.Texture2D {
	return b.textures[state]
}

// updateTexture changes the texture of the button based on the present button state.
func (b *Button) updateTexture() {
	b.SetTexture(b.textures[b.GetButtonState()])
}

// GetButtonState returns present button state.
func (b *Button) GetButtonState() ButtonState {
	if !b.Enabled() {
		return ButtonDisabled
	} else if b.pressed {
		return ButtonPressed
	} else if b.mouseOver {
		return ButtonOver
	} else {
		return ButtonNormal
	}
}

// SetStyles set the button styles.
func (b *Button) SetStyles(bs *ButtonStyles) {
	b.styles = bs
	b.updateStyle()
}

// applyStyle applies a button style.
func (b *Button) applyStyle(bs *ButtonStyle) {
	b.Panel.ApplyStyle(&bs.PanelStyle)
	b.Label.SetColor(bs.FgColor)
}

// updateStyle applies a button style depending on the button state.
func (b *Button) updateStyle() {
	switch b.GetButtonState() {
	case ButtonNormal:
		b.applyStyle(&b.styles.Normal)
	case ButtonOver:
		b.applyStyle(&b.styles.Over)
	case ButtonPressed:
		b.applyStyle(&b.styles.Pressed)
	case ButtonDisabled:
		b.applyStyle(&b.styles.Disabled)
	}
}

// SetExpandToLabel sets whether this button resizes automatically make room for its label.
func (b *Button) SetExpandToLabel(expand bool) {
	b.expandToLabel = expand
}

// GetExpandToLabel returns whether this button resizes automatically make room for its label.
func (b *Button) GetExpandToLabel() bool {
	return b.expandToLabel
}

// recalculateSize recalculates all dimensions and position from inside out.
func (b *Button) recalculateSize() {
	// Current width and height of button content area
	width := b.ContentWidth()
	height := b.ContentHeight()

	labelWidth, labelHeight := b.Label.Size()
	if b.Label.Text() == "" {
		labelWidth, labelHeight = 0, 0
	}

	// Sets new content width and height if necessary
	if b.expandToLabel {
		resize := false
		if width < labelWidth {
			width = labelWidth
			resize = true
		}
		if height < labelHeight {
			height = labelHeight
			resize = true
		}
		if resize {
			b.SetContentSize(width, height)
		}
	}

	// Centralize or pin to 0, 0
	lx := max(0, (width-labelWidth)/2)
	ly := max(0, (height-labelHeight)/2)
	b.Label.SetPosition(lx, ly)
}
