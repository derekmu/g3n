// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
)

/***************************************

 Slider
 +--------------------------------+
 |  +--------------------------+  |
 |  |      +----------+        |  |
 |  |      |          |        |  |
 |  |      |          |        |  |
 |  |      +----------+        |  |
 |  +--------------------------+  |
 +--------------------------------+

**/

// Slider is the GUI element for sliders and progress bars
type Slider struct {
	Panel               // Embedded panel
	slider      Panel   // embedded slider panel
	label       *Label  // optional label
	horiz       bool    // orientation
	pos         float32 // current slider position
	posLast     float32 // last position of the mouse cursor when dragging
	pressed     bool    // mouse button is pressed and dragging
	cursorOver  bool    // mouse is over slider
	scaleFactor float32 // scale factor (default = 1.0)
}

// NewHSlider creates and returns a pointer to a new horizontal slider
// with the specified initial dimensions.
func NewHSlider(width, height float32) *Slider {
	return newSlider(true, width, height)
}

// NewVSlider creates and returns a pointer to a new vertical slider
// with the specified initial dimensions.
func NewVSlider(width, height float32) *Slider {
	return newSlider(false, width, height)
}

// NewSlider creates and returns a pointer to a new slider with the
// specified initial dimensions.
func newSlider(horiz bool, width, height float32) *Slider {
	s := new(Slider)
	s.horiz = horiz
	s.scaleFactor = 1.0

	// Initialize main panel
	s.Panel.InitPanel(s, width, height)
	s.Panel.Subscribe(OnMouseDown, s.onMouse)
	s.Panel.Subscribe(OnMouseUp, s.onMouse)
	s.Panel.Subscribe(OnCursor, s.onCursor)
	s.Panel.Subscribe(OnCursorEnter, s.onCursor)
	s.Panel.Subscribe(OnCursorLeave, s.onCursor)
	s.Panel.Subscribe(OnScroll, s.onScroll)
	s.Panel.Subscribe(OnKeyDown, s.onKey)
	s.Panel.Subscribe(OnKeyRepeat, s.onKey)
	s.Panel.Subscribe(OnResize, s.onResize)

	// Initialize slider panel
	s.slider.InitPanel(&s.slider, 0, 0)
	s.Panel.Add(&s.slider)

	s.recalc()
	return s
}

// SetText sets the text of the slider optional label
func (s *Slider) SetText(text string) *Slider {
	if s.label == nil {
		s.label = NewLabel(text)
		s.Panel.Add(s.label)
	} else {
		s.label.SetText(text)
	}
	s.recalc()
	return s
}

// SetValue sets the value of the slider considering the current scale factor
// and updates its visual appearance.
func (s *Slider) SetValue(value float32) *Slider {
	pos := value / s.scaleFactor
	s.setPos(pos)
	return s
}

// Value returns the current value of the slider considering the current scale factor
func (s *Slider) Value() float32 {
	return s.pos * s.scaleFactor
}

// SetScaleFactor set the slider scale factor (default = 1.0)
func (s *Slider) SetScaleFactor(factor float32) *Slider {
	s.scaleFactor = factor
	return s
}

// ScaleFactor returns  the slider current scale factor (default = 1.0)
func (s *Slider) ScaleFactor() float32 {
	return s.scaleFactor
}

// setPos sets the slider position from 0.0 to 1.0
// and updates its visual appearance.
func (s *Slider) setPos(pos float32) {
	const eps = 0.01
	if pos < 0 {
		pos = 0
	} else if pos > 1.0 {
		pos = 1
	}
	if pos > (s.pos+eps) && pos < (s.pos+eps) {
		return
	}
	s.pos = pos
	s.recalc()
	s.Dispatch(OnChange, nil)
}

// onMouse process subscribed mouse events over the outer panel
func (s *Slider) onMouse(evname string, ev interface{}) {
	if !s.Enabled() {
		return
	}

	mev := ev.(*core.MouseEvent)
	if mev.Button != core.MouseButtonLeft {
		return
	}
	switch evname {
	case OnMouseDown:
		s.pressed = true
		if s.horiz {
			s.posLast = mev.Xpos
		} else {
			s.posLast = mev.Ypos
		}
		GetManager().SetKeyFocus(s)
		GetManager().SetCursorFocus(s)
	case OnMouseUp:
		s.pressed = false
		GetManager().SetCursorFocus(nil)
	default:
		return
	}
}

// onCursor process subscribed cursor events
func (s *Slider) onCursor(evname string, ev interface{}) {
	if !s.Enabled() {
		return
	}

	if evname == OnCursorEnter {
		s.cursorOver = true
		if s.horiz {
			GetManager().window.SetCursor(core.HResizeCursor)
		} else {
			GetManager().window.SetCursor(core.VResizeCursor)
		}
	} else if evname == OnCursorLeave {
		s.cursorOver = false
		GetManager().window.SetCursor(core.ArrowCursor)
	} else if evname == OnCursor {
		if !s.pressed {
			return
		}
		cev := ev.(*core.CursorEvent)
		var pos float32
		if s.horiz {
			delta := cev.Xpos - s.posLast
			s.posLast = cev.Xpos
			newpos := s.slider.Width() + delta
			pos = newpos / s.Panel.ContentWidth()
		} else {
			delta := cev.Ypos - s.posLast
			s.posLast = cev.Ypos
			newpos := s.slider.Height() - delta
			pos = newpos / s.Panel.ContentHeight()
		}
		s.setPos(pos)
	}
}

// onScroll process subscribed scroll events
func (s *Slider) onScroll(_ string, ev interface{}) {
	if !s.Enabled() {
		return
	}

	sev := ev.(*core.ScrollEvent)
	v := s.pos
	v += sev.Yoffset * 0.01
	s.setPos(v)
}

// onKey process subscribed key events
func (s *Slider) onKey(_ string, ev interface{}) {
	if !s.Enabled() {
		return
	}

	kev := ev.(*core.KeyEvent)
	delta := float32(0.01)
	// Horizontal slider
	if s.horiz {
		switch kev.Key {
		case core.KeyLeft:
			s.setPos(s.pos - delta)
		case core.KeyRight:
			s.setPos(s.pos + delta)
		default:
			return
		}
		// Vertical slider
	} else {
		switch kev.Key {
		case core.KeyDown:
			s.setPos(s.pos - delta)
		case core.KeyUp:
			s.setPos(s.pos + delta)
		default:
			return
		}
	}
}

// onResize process subscribed resize events
func (s *Slider) onResize(_ string, _ interface{}) {
	s.recalc()
}

// recalc recalculates the dimensions and positions of the internal panels.
func (s *Slider) recalc() {
	if s.horiz {
		if s.label != nil {
			lx := (s.Panel.ContentWidth() - s.label.Width()) / 2
			if s.Panel.ContentHeight() < s.label.Height() {
				s.Panel.SetContentHeight(s.label.Height())
			}
			ly := (s.Panel.ContentHeight() - s.label.Height()) / 2
			s.label.SetPosition(lx, ly)
		}
		width := s.Panel.ContentWidth() * s.pos
		s.slider.SetSize(width, s.Panel.ContentHeight())
	} else {
		if s.label != nil {
			if s.Panel.ContentWidth() < s.label.Width() {
				s.Panel.SetContentWidth(s.label.Width())
			}
			lx := (s.Panel.ContentWidth() - s.label.Width()) / 2
			ly := (s.Panel.ContentHeight() - s.label.Height()) / 2
			s.label.SetPosition(lx, ly)
		}
		height := s.Panel.ContentHeight() * s.pos
		s.slider.SetPositionY(s.Panel.ContentHeight() - height)
		s.slider.SetSize(s.Panel.ContentWidth(), height)
	}
}
