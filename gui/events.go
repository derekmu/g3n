// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
)

// Core events sent by the GUI eventManager.
// The target panel is the panel immediately under the mouse cursor.

// Events sent to target panel's lowest subscribed ancestor.

const (
	OnMouseDown = core.OnMouseDown
	OnMouseUp   = core.OnMouseUp
	OnScroll    = core.OnScroll
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

const OnCursor = core.OnCursor

// Events sent to the new key-focused IDispatcher, specified on a call to gui.eventManager().SetKeyFocus(core.IDispatcher).

const (
	OnFocus     = "gui.OnFocus"
	OnFocusLost = "gui.OnFocusLost"
)

// Events sent to the key-focused IDispatcher.

const (
	OnKeyDown   = core.OnKeyDown
	OnKeyUp     = core.OnKeyUp
	OnKeyRepeat = core.OnKeyRepeat
	OnChar      = core.OnChar
)

// Events sent in other circumstances.

const (
	OnResize     = "gui.OnResize"     // Panel size changed (no parameters)
	OnEnable     = "gui.OnEnable"     // Panel enabled/disabled (no parameters)
	OnClick      = "gui.OnClick"      // Widget clicked by mouse left button or via key press
	OnChange     = "gui.OnChange"     // Value was changed.
	OnRadioGroup = "gui.OnRadioGroup" // Radio button within a group changed state
)
