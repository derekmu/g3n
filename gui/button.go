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
	ButtonStateMin   = ButtonNormal
	ButtonStateMax   = ButtonDisabled
	ButtonStateCount = ButtonStateMax - ButtonStateMin + 1
)

// Button is a UI element that dispatches click events and uses different textures for each ButtonState.
type Button struct {
	Panel
	mouseOver bool
	pressed   core.MouseState
	textures  [ButtonStateCount]*texture.Texture2D
}

// NewButton creates a new Button with the specified text the button label.
func NewButton() *Button {
	b := new(Button)
	b.InitButton()
	return b
}

// InitButton initializes the image and subscribes to events.
func (b *Button) InitButton() {
	b.InitPanel(b, 0, 0)
	b.Subscribe(b.onGuiEvent)
	b.updateTexture()
}

// Dispose disposes of the label and all button textures.
func (b *Button) Dispose() {
	b.Panel.Dispose()
	for i, tex := range b.textures {
		if tex != nil && tex.Decref() {
			tex.Dispose()
		}
		b.textures[i] = nil
	}
}

// SetStateTexture changes the texture used by the button in a given state.
func (b *Button) SetStateTexture(state ButtonState, tex *texture.Texture2D) {
	ptex := b.textures[state]
	if ptex != nil && ptex.Decref() {
		ptex.Dispose()
	}
	if tex != nil {
		tex.Incref()
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
	} else if b.pressed != 0 {
		return ButtonPressed
	} else if b.mouseOver {
		return ButtonOver
	} else {
		return ButtonNormal
	}
}

func (b *Button) onGuiEvent(event core.GuiEvent) bool {
	switch ev := event.(type) {
	case core.MouseUpEvent:
		if b.Enabled() {
			clicked := b.pressed.IsSet(ev.Button)
			b.pressed = b.pressed.Unset(ev.Button)
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
		if b.Enabled() {
			b.pressed = b.pressed.Set(ev.Button)
			b.updateTexture()
		}
	case core.GuiCursorEnterEvent:
		b.mouseOver = true
		b.updateTexture()
	case core.GuiCursorLeaveEvent:
		b.mouseOver = false
		// Pressing and dragging out cancels clicks
		b.pressed = 0
		b.updateTexture()
	case core.GuiEnableEvent:
		// Enabling or disabling a button cancels clicks
		b.pressed = 0
		b.updateTexture()
	default:
		return false
	}
	return true
}
