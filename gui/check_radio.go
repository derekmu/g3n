// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gui/assets/icon"
)

const (
	checkON  = icon.CheckBox
	checkOFF = icon.CheckBoxOutlineBlank
	radioON  = icon.RadioButtonChecked
	radioOFF = icon.RadioButtonUnchecked
)

// CheckRadio is a GUI element that can be either a checkbox or a radio button
type CheckRadio struct {
	Panel             // Embedded panel
	Label      *Label // Text label
	icon       *Label
	check      bool
	group      string // current group name
	cursorOver bool
	state      bool
	codeON     string
	codeOFF    string
	subroot    bool // indicates root subcription
}

// NewCheckBox creates and returns a pointer to a new CheckBox widget
// with the specified text
func NewCheckBox(text string) *CheckRadio {
	return newCheckRadio(true, text)
}

// NewRadioButton creates and returns a pointer to a new RadioButton widget
// with the specified text
func NewRadioButton(text string) *CheckRadio {
	return newCheckRadio(false, text)
}

// newCheckRadio creates and returns a pointer to a new CheckRadio widget
// with the specified type and text
func newCheckRadio(check bool, text string) *CheckRadio {
	cb := new(CheckRadio)

	// Adapts to specified type: CheckBox or RadioButton
	cb.check = check
	cb.state = false
	if cb.check {
		cb.codeON = checkON
		cb.codeOFF = checkOFF
	} else {
		cb.codeON = radioON
		cb.codeOFF = radioOFF
	}

	// Initialize panel
	cb.Panel.InitPanel(cb, 0, 0)

	// Subscribe to events
	cb.Panel.Subscribe(OnKeyDown, cb.onKey)
	cb.Panel.Subscribe(OnCursorEnter, cb.onCursor)
	cb.Panel.Subscribe(OnCursorLeave, cb.onCursor)
	cb.Panel.Subscribe(OnMouseDown, cb.onMouse)
	cb.Panel.Subscribe(OnEnable, func(evname string, ev interface{}) { cb.update() })

	// Creates label
	cb.Label = NewLabel(text)
	cb.Label.Subscribe(OnResize, func(evname string, ev interface{}) { cb.recalc() })
	cb.Panel.Add(cb.Label)

	// Creates icon label
	cb.icon = NewIconLabel(" ")
	cb.Panel.Add(cb.icon)

	cb.recalc()
	cb.update()
	return cb
}

// Value returns the current state of the checkbox
func (cb *CheckRadio) Value() bool {
	return cb.state
}

// SetValue sets the current state of the checkbox
func (cb *CheckRadio) SetValue(state bool) *CheckRadio {
	if state == cb.state {
		return cb
	}
	cb.state = state
	cb.update()
	cb.Dispatch(OnChange, nil)
	return cb
}

// Group returns the name of the radio group
func (cb *CheckRadio) Group() string {
	return cb.group
}

// SetGroup sets the name of the radio group
func (cb *CheckRadio) SetGroup(group string) *CheckRadio {
	cb.group = group
	return cb
}

// toggleState toggles the current state of the checkbox/radiobutton
func (cb *CheckRadio) toggleState() {
	// Subscribes once to the root panel for OnRadioGroup events
	// The root panel is used to dispatch events to all checkradios
	if !cb.subroot {
		GetManager().Subscribe(OnRadioGroup, func(name string, ev interface{}) {
			cb.onRadioGroup(ev.(*CheckRadio))
		})
		cb.subroot = true
	}

	if cb.check {
		cb.state = !cb.state
	} else {
		if len(cb.group) == 0 {
			cb.state = !cb.state
		} else {
			if cb.state {
				return
			}
			cb.state = !cb.state
		}
	}
	cb.update()
	cb.Dispatch(OnChange, nil)
	if !cb.check && len(cb.group) > 0 {
		GetManager().Dispatch(OnRadioGroup, cb)
	}
}

// onMouse process OnMouseDown events
func (cb *CheckRadio) onMouse(evname string, ev interface{}) {
	// Dispatch OnClick for left mouse button down
	if evname == OnMouseDown {
		mev := ev.(*core.MouseEvent)
		if mev.Button == core.MouseButtonLeft && cb.Enabled() {
			GetManager().SetKeyFocus(cb)
			cb.toggleState()
			cb.Dispatch(OnClick, nil)
		}
	}
}

// onCursor process OnCursor* events
func (cb *CheckRadio) onCursor(evname string, _ interface{}) {
	if evname == OnCursorEnter {
		cb.cursorOver = true
	} else {
		cb.cursorOver = false
	}
	cb.update()
}

// onKey receives subscribed key events
func (cb *CheckRadio) onKey(evname string, ev interface{}) {
	kev := ev.(*core.KeyEvent)
	if evname == OnKeyDown && kev.Key == core.KeyEnter {
		cb.toggleState()
		cb.update()
		cb.Dispatch(OnClick, nil)
		return
	}
	return
}

// onRadioGroup receives subscribed OnRadioGroup events
func (cb *CheckRadio) onRadioGroup(other *CheckRadio) {
	// If event is for this button, ignore
	if cb == other {
		return
	}
	// If other radio group is not the group of this button, ignore
	if cb.group != other.group {
		return
	}
	// Toggle this button state
	cb.SetValue(!other.Value())
}

// update updates the visual appearance of the checkbox
func (cb *CheckRadio) update() {
	if cb.state {
		cb.icon.SetText(cb.codeON)
	} else {
		cb.icon.SetText(cb.codeOFF)
	}
}

// recalc recalculates dimensions and position from inside out
func (cb *CheckRadio) recalc() {
	// Sets icon position
	cb.icon.SetFontSize(cb.Label.FontSize() * 1.3)
	cb.icon.SetPosition(0, 0)

	// Label position
	spacing := float32(4)
	cb.Label.SetPosition(cb.icon.Width()+spacing, 0)

	// Content width
	width := cb.icon.Width() + spacing + cb.Label.Width()
	cb.SetContentSize(width, cb.Label.Height())
}
