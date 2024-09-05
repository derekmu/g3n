// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/texture"
)

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
	Label          *Label
	expandToLabel  bool
	labelAlignment Align
	mouseOver      bool
	pressed        bool
	textures       [ButtonDisabled + 1]*texture.Texture2D
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

	b.Label = NewLabel(text)
	b.Add(b.Label)
	b.Label.Subscribe(b.onLabelEvent)
	b.expandToLabel = true
	b.labelAlignment = AlignCenterCenter

	b.Subscribe(b.onGuiEvent)

	b.updateTexture()
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

// SetExpandToLabel sets whether this button resizes automatically make room for its label.
func (b *Button) SetExpandToLabel(expand bool) {
	b.expandToLabel = expand
}

// GetExpandToLabel returns whether this button resizes automatically make room for its label.
func (b *Button) GetExpandToLabel() bool {
	return b.expandToLabel
}

// SetLabelAlignment sets how this button aligns its label within its content area.
func (b *Button) SetLabelAlignment(align Align) {
	b.labelAlignment = align
}

// GetLabelAlignment returns how this button aligns its label within its content area.
func (b *Button) GetLabelAlignment() Align {
	return b.labelAlignment
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

	// Position the label as desired
	if b.labelAlignment != AlignNone {
		lx, ly := b.labelAlignment.CalculatePosition(width, height, labelWidth, labelHeight)
		b.Label.SetPosition(lx, ly)
	}
}

func (b *Button) onLabelEvent(event core.GuiEvent) bool {
	switch event.GuiEventType() {
	case core.GuiResize:
		b.recalculateSize()
	default:
		return false
	}
	return true
}

func (b *Button) onGuiEvent(event core.GuiEvent) bool {
	switch ev := event.(type) {
	case core.MouseUpEvent:
		if ev.Button == core.MouseButtonLeft {
			clicked := b.pressed && b.mouseOver
			b.pressed = false
			b.updateTexture()
			if clicked {
				b.Dispatch(core.GuiClickEvent{
					X:      ev.X,
					Y:      ev.Y,
					Button: ev.Button,
					Mods:   ev.Mods,
				})
			}
		}
	case core.MouseDownEvent:
		if ev.Button == core.MouseButtonLeft {
			b.pressed = true
			b.updateTexture()
		}
	case core.GuiCursorEnterEvent:
		b.mouseOver = true
		b.updateTexture()
	case core.GuiCursorLeaveEvent:
		b.mouseOver = false
		// Pressing, dragging out, and releasing cancels the click
		b.pressed = false
		b.updateTexture()
	case core.GuiEnableEvent:
		// Enabling or disabling a button cancels the click
		b.pressed = false
		b.updateTexture()
	case core.GuiResizeEvent:
		b.recalculateSize()
	default:
		return false
	}
	return true
}
