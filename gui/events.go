// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/window"
)

// Core events sent by the GUI manager.
// The target panel is the panel immediately under the mouse cursor.

// Events sent to target panel's lowest subscribed ancestor.

const (
	OnMouseDown = window.OnMouseDown
	OnMouseUp   = window.OnMouseUp
	OnScroll    = window.OnScroll
)

// Events sent to all panels except the ancestors of the target panel.

const (
	OnMouseDownOut = "gui.OnMouseDownOut"
	OnMouseUpOut   = "gui.OnMouseUpOut"
)

// Events sent to new target panel and all of its ancestors up to, but not including, the common ancestor of the new and old targets.

const OnCursorEnter = "gui.OnCursorEnter"

// Events sent to old target panel and all of its ancestors up to (not including) the common ancestor of the new and old targets.

const OnCursorLeave = "gui.OnCursorLeave"

// Events sent to the cursor-focused IDispatcher if any, else sent to target panel's lowest subscribed ancestor.

const OnCursor = window.OnCursor

// Events sent to the new key-focused IDispatcher, specified on a call to gui.Manager().SetKeyFocus(core.IDispatcher).

const (
	OnFocus     = "gui.OnFocus"
	OnFocusLost = "gui.OnFocusLost"
)

// Events sent to the key-focused IDispatcher.

const (
	OnKeyDown   = window.OnKeyDown
	OnKeyUp     = window.OnKeyUp
	OnKeyRepeat = window.OnKeyRepeat
	OnChar      = window.OnChar
)

// Events sent in other circumstances.

const (
	OnResize     = "gui.OnResize"     // Panel size changed (no parameters)
	OnEnable     = "gui.OnEnable"     // Panel enabled/disabled (no parameters)
	OnClick      = "gui.OnClick"      // Widget clicked by mouse left button or via key press
	OnChange     = "gui.OnChange"     // Value was changed. Emitted by List, DropDownList, CheckBox and Edit
	OnRadioGroup = "gui.OnRadioGroup" // Radio button within a group changed state
)
