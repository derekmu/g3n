// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/math32"
)

/***************************************

 ScrollBar Panel
 +--------------------------------+
 |         scroll button          |
 |        +--------------+        |
 |        |              |        |
 |        |              |        |
 |        +--------------+        |
 +--------------------------------+

**/

// ScrollBar is the scrollbar GUI element.
type ScrollBar struct {
	Panel                      // Embedded panel
	vertical   bool            // type of scrollbar
	button     scrollBarButton // scrollbar button
	cursorOver bool
}

type scrollBarButton struct {
	Panel              // Embedded panel
	sb      *ScrollBar // pointer to parent scroll bar
	pressed bool       // mouse button pressed flag
	mouseX  float32    // last mouse click x position
	mouseY  float32    // last mouse click y position
	Size    float32    // button size
	MinSize float32    // minimum button size
}

// NewVScrollBar creates and returns a pointer to a new vertical scroll bar
// with the specified dimensions.
func NewVScrollBar(width, height float32) *ScrollBar {
	return newScrollBar(width, height, true)
}

// NewHScrollBar creates and returns a pointer to a new horizontal scroll bar
// with the specified dimensions.
func NewHScrollBar(width, height float32) *ScrollBar {
	return newScrollBar(width, height, false)
}

// newScrollBar creates and returns a pointer to a new scroll bar panel
// with the specified width, height, orientation and target.
func newScrollBar(width, height float32, vertical bool) *ScrollBar {
	sb := new(ScrollBar)
	sb.InitScrollBar(width, height, vertical)
	return sb
}

// InitScrollBar initializes this scrollbar
func (sb *ScrollBar) InitScrollBar(width, height float32, vertical bool) {
	sb.vertical = vertical
	sb.Panel.InitPanel(sb, width, height)
	sb.Panel.Subscribe(OnMouseDown, sb.onMouse)

	// Initialize scrollbar button
	sb.button.Panel.InitPanel(&sb.button, 0, 0)
	sb.button.Panel.Subscribe(OnMouseDown, sb.button.onMouse)
	sb.button.Panel.Subscribe(OnMouseUp, sb.button.onMouse)
	sb.button.Panel.Subscribe(OnCursor, sb.button.onCursor)
	sb.button.SetMargins(RectBounds{1, 1, 1, 1})
	sb.button.Size = 32
	sb.button.sb = sb
	sb.Add(&sb.button)

	sb.recalc()
}

// SetButtonSize sets the button size
func (sb *ScrollBar) SetButtonSize(size float32) {
	// Clamp to minimum size if requested size smaller than minimum
	if size > sb.button.MinSize {
		sb.button.Size = size
	} else {
		sb.button.Size = sb.button.MinSize
	}
	sb.recalc()
}

// Value returns the current position of the button in the scrollbar
// The returned value is between 0.0 and 1.0
func (sb *ScrollBar) Value() float64 {
	if sb.vertical {
		den := float64(sb.content.Height) - float64(sb.button.height)
		if den == 0 {
			return 0
		}
		return float64(sb.button.Position().Y) / den
	}

	// horizontal
	den := float64(sb.content.Width) - float64(sb.button.width)
	if den == 0 {
		return 0
	}
	return float64(sb.button.Position().X) / den
}

// SetValue sets the position of the button of the scrollbar
// from 0.0 (minimum) to 1.0 (maximum).
func (sb *ScrollBar) SetValue(v float32) {
	v = math32.Clamp(v, 0.0, 1.0)
	if sb.vertical {
		pos := v * (float32(sb.content.Height) - float32(sb.button.height))
		sb.button.SetPositionY(pos)
	} else {
		pos := v * (float32(sb.content.Width) - float32(sb.button.width))
		sb.button.SetPositionX(pos)
	}
}

// onMouse receives subscribed mouse events over the scrollbar outer panel
func (sb *ScrollBar) onMouse(_ string, ev interface{}) {
	e := ev.(*core.MouseEvent)
	if e.Button != core.MouseButtonLeft {
		return
	}
	if sb.vertical {
		posy := e.Ypos - sb.pospix.Y
		newY := math32.Clamp(posy-(sb.button.height/2), 0, sb.content.Height-sb.button.height)
		sb.button.SetPositionY(newY)
	} else {
		posx := e.Xpos - sb.pospix.X
		newX := math32.Clamp(posx-(sb.button.width/2), 0, sb.content.Width-sb.button.width)
		sb.button.SetPositionX(newX)
	}
	sb.Dispatch(OnChange, nil)
}

// recalc recalculates sizes and positions
func (sb *ScrollBar) recalc() {
	if sb.vertical {
		sb.button.SetSize(sb.content.Width, sb.button.Size)
	} else {
		sb.button.SetSize(sb.button.Size, sb.content.Height)
	}
}

// onMouse receives subscribed mouse events for the scroll bar button
func (button *scrollBarButton) onMouse(evname string, ev interface{}) {
	e := ev.(*core.MouseEvent)
	if e.Button != core.MouseButtonLeft {
		return
	}
	switch evname {
	case OnMouseDown:
		button.pressed = true
		button.mouseX = e.Xpos
		button.mouseY = e.Ypos
		GetManager().SetCursorFocus(button)
	case OnMouseUp:
		button.pressed = false
		GetManager().SetCursorFocus(nil)
	default:
		return
	}
}

// onCursor receives subscribed cursor events for the scroll bar button
func (button *scrollBarButton) onCursor(_ string, ev interface{}) {
	e := ev.(*core.CursorEvent)
	if !button.pressed {
		return
	}
	if button.sb.vertical {
		dy := button.mouseY - e.Ypos
		py := button.Position().Y
		button.SetPositionY(math32.Clamp(py-dy, 0, button.sb.content.Height-button.Size))
	} else {
		dx := button.mouseX - e.Xpos
		px := button.Position().X
		button.SetPositionX(math32.Clamp(px-dx, 0, button.sb.content.Width-button.Size))
	}
	button.mouseX = e.Xpos
	button.mouseY = e.Ypos
	button.sb.Dispatch(OnChange, nil)
}
