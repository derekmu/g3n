// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"github.com/derekmu/g3n/core"
)

var gm *Manager

// GetManager returns the GUI Manager singleton or panics if InitManager hasn't been called.
func GetManager() *Manager {
	// Return singleton if already created
	if gm != nil {
		return gm
	}
	panic(fmt.Errorf("gui.InitManager not called"))
}

type IWindow interface {
	core.IDispatcher
	GetScale() (x float64, y float64)
	SetCursor(cursor core.Cursor)
}

// Manager routes events to the appropriate places
type Manager struct {
	core.Dispatcher   // Embedded Dispatcher
	core.TimerManager // Embedded TimerManager
	window            IWindow
	scene             core.INode        // INode containing IPanels to dispatch events to (can contain non-IPanels as well)
	modal             IPanel            // Panel which along its descendants will exclusively receive all events
	target            IPanel            // Panel immediately under the cursor
	keyFocus          core.IDispatcher  // IDispatcher which will exclusively receive all key and char events
	cursorFocus       core.IDispatcher  // IDispatcher which will exclusively receive all OnCursor events
	cev               *core.CursorEvent // IDispatcher which will exclusively receive all OnCursor events
}

// InitManager creates the GUI Manager or panics if it's already been called.
func InitManager(window IWindow) {
	if gm != nil {
		panic(fmt.Errorf("gui.InitManager already called"))
	}
	gm = new(Manager)
	gm.Dispatcher.Initialize()
	gm.TimerManager.Initialize()
	gm.window = window

	// Subscribe to window events
	window.Subscribe(core.OnKeyUp, gm.onKeyboard)
	window.Subscribe(core.OnKeyDown, gm.onKeyboard)
	window.Subscribe(core.OnKeyRepeat, gm.onKeyboard)
	window.Subscribe(core.OnChar, gm.onKeyboard)
	window.Subscribe(core.OnCursor, gm.onCursor)
	window.Subscribe(core.OnMouseUp, gm.onMouse)
	window.Subscribe(core.OnMouseDown, gm.onMouse)
	window.Subscribe(core.OnScroll, gm.onScroll)
}

// Set sets the INode to watch for events.
// It's usually a scene containing a hierarchy of INodes.
// The Manager only cares about IPanels inside that hierarchy.
func (gm *Manager) Set(scene core.INode) {
	gm.scene = scene
}

// SetModal sets the specified panel and its descendants to be the exclusive receivers of events.
func (gm *Manager) SetModal(ipan IPanel) {
	gm.modal = ipan
	gm.SetKeyFocus(nil)
	gm.SetCursorFocus(nil)
}

// SetKeyFocus sets the key-focused IDispatcher, which will exclusively receive key and char events.
func (gm *Manager) SetKeyFocus(disp core.IDispatcher) {
	if gm.keyFocus == disp {
		return
	}
	if gm.keyFocus != nil {
		gm.keyFocus.Dispatch(OnFocusLost, nil)
	}
	gm.keyFocus = disp
	if gm.keyFocus != nil {
		gm.keyFocus.Dispatch(OnFocus, nil)
	}
}

// SetCursorFocus sets the cursor-focused IDispatcher, which will exclusively receive OnCursor events.
func (gm *Manager) SetCursorFocus(disp core.IDispatcher) {
	if gm.cursorFocus == disp {
		return
	}
	gm.cursorFocus = disp
	if gm.cursorFocus == nil {
		gm.onCursor(OnCursor, gm.cev)
	}
}

// onKeyboard is called when char or key events are received.
// The events are dispatched to the focused IDispatcher or to non-GUI.
func (gm *Manager) onKeyboard(evname string, ev any) {
	if gm.keyFocus != nil {
		if gm.modal == nil {
			gm.keyFocus.Dispatch(evname, ev)
		} else if ipan, ok := gm.keyFocus.(IPanel); ok && gm.modal.IsAncestorOf(ipan) {
			gm.keyFocus.Dispatch(evname, ev)
		}
	} else {
		gm.Dispatch(evname, ev)
	}
}

// onMouse is called when mouse events are received.
// OnMouseDown/OnMouseUp are dispatched to gm.target or to non-GUI, while
// OnMouseDownOut/OnMouseUpOut are dispatched to all non-target panels.
func (gm *Manager) onMouse(evname string, ev any) {
	// To fix #299
	if gm.cev == nil {
		mev := ev.(*core.MouseEvent)
		gm.cev = &core.CursorEvent{
			Xpos: mev.Xpos,
			Ypos: mev.Ypos,
			Mods: mev.Mods,
		}
	}

	// Check if gm.scene is nil and if so then there are no IPanels to send events to
	if gm.scene == nil {
		gm.Dispatch(evname, ev) // Dispatch event to non-GUI since event was not filtered by any GUI component
		return
	}

	// Dispatch OnMouseDownOut/OnMouseUpOut to all panels except ancestors of target
	gm.forEachIPanel(func(ipan IPanel) {
		if gm.target == nil || !ipan.IsAncestorOf(gm.target) {
			switch evname {
			case OnMouseDown:
				ipan.Dispatch(OnMouseDownOut, ev)
			case OnMouseUp:
				ipan.Dispatch(OnMouseUpOut, ev)
			}
		}
	})

	// Appropriately dispatch the event to target panel's lowest subscribed ancestor or to non-GUI or not at all
	if gm.target != nil {
		if gm.modal == nil || gm.modal.IsAncestorOf(gm.target) {
			sendAncestry(gm.target, false, nil, gm.modal, evname, ev)
		}
	} else if gm.modal == nil {
		gm.Dispatch(evname, ev)
	}
}

// onScroll is called when scroll events are received.
// The events are dispatched to the target panel or to non-GUI.
func (gm *Manager) onScroll(evname string, ev any) {
	// Check if gm.scene is nil and if so then there are no IPanels to send events to
	if gm.scene == nil {
		gm.Dispatch(evname, ev) // Dispatch event to non-GUI since event was not filtered by any GUI component
		return
	}

	// Appropriately dispatch the event to target panel's lowest subscribed ancestor or to non-GUI or not at all
	if gm.target != nil {
		if gm.modal == nil || gm.modal.IsAncestorOf(gm.target) {
			sendAncestry(gm.target, false, nil, gm.modal, evname, ev)
		}
	} else if gm.modal == nil {
		gm.Dispatch(evname, ev)
	}
}

// onCursor is called when (mouse) cursor events are received.
// Updates the target/click panels and dispatches OnCursor, OnCursorEnter, OnCursorLeave events.
func (gm *Manager) onCursor(evname string, ev any) {
	// If an IDispatcher is capturing cursor events dispatch to it and return
	if gm.cursorFocus != nil {
		gm.cursorFocus.Dispatch(evname, ev)
		return
	}

	// If gm.scene is nil then there are no IPanels to send events to
	if gm.scene == nil {
		gm.Dispatch(evname, ev) // Dispatch event to non-GUI since event was not filtered by any GUI component
		return
	}

	// Get and store CursorEvent
	gm.cev = ev.(*core.CursorEvent)

	// Temporarily store last target and clear current one
	oldTarget := gm.target
	gm.target = nil

	// Find IPanel immediately under the cursor and store it in gm.target
	gm.forEachIPanel(func(ipan IPanel) {
		if ipan.InsideBorders(gm.cev.Xpos, gm.cev.Ypos) && (gm.target == nil || ipan.Position().Z < gm.target.GetPanel().Position().Z) {
			gm.target = ipan
		}
	})

	// If the cursor is now over a different panel, dispatch OnCursorLeave/OnCursorEnter
	if gm.target != oldTarget {
		// We are only interested in sending events up to the lowest common ancestor of target and oldTarget
		var commonAnc IPanel
		if gm.target != nil && oldTarget != nil {
			commonAnc, _ = gm.target.LowestCommonAncestor(oldTarget).(IPanel)
		}
		// If just left a panel and the new panel is not a descendant of the old panel
		if oldTarget != nil && !oldTarget.IsAncestorOf(gm.target) && (gm.modal == nil || gm.modal.IsAncestorOf(oldTarget)) {
			sendAncestry(oldTarget, true, commonAnc, gm.modal, OnCursorLeave, ev)
		}
		// If just entered a panel and it's not an ancestor of the old panel
		if gm.target != nil && !gm.target.IsAncestorOf(oldTarget) && (gm.modal == nil || gm.modal.IsAncestorOf(gm.target)) {
			sendAncestry(gm.target, true, commonAnc, gm.modal, OnCursorEnter, ev)
		}
	}

	// Appropriately dispatch the event to target panel's lowest subscribed ancestor or to non-GUI or not at all
	if gm.target != nil {
		if gm.modal == nil || gm.modal.IsAncestorOf(gm.target) {
			sendAncestry(gm.target, false, nil, gm.modal, evname, ev)
		}
	} else if gm.modal == nil {
		gm.Dispatch(evname, ev)
	}
}

// sendAncestry sends the specified event (evname/ev) to the specified target panel and its ancestors.
// If all is false, then the event is only sent to the lowest subscribed ancestor.
// If uptoEx (i.e. excluding) is not nil then the event will not be dispatched to that ancestor nor any higher ancestors.
// If uptoIn (i.e. including) is not nil then the event will be dispatched to that ancestor but not to any higher ancestors.
// uptoEx and uptoIn can both be defined.
func sendAncestry(ipan IPanel, all bool, uptoEx IPanel, uptoIn IPanel, evname string, ev any) {
	var ok bool
	for ipan != nil {
		if uptoEx != nil && ipan == uptoEx {
			break
		}
		count := ipan.Dispatch(evname, ev)
		if (uptoIn != nil && ipan == uptoIn) || (!all && count > 0) {
			break
		}
		ipan, ok = ipan.Parent().(IPanel)
		if !ok {
			break
		}
	}
}

// traverseIPanel traverses the descendants of the provided IPanel,
// executing the specified function for each IPanel.
func traverseIPanel(ipan IPanel, f func(ipan IPanel)) {
	// If panel not visible, ignore entire hierarchy below this point
	if !ipan.Visible() {
		return
	}
	if ipan.Enabled() {
		f(ipan) // Call specified function
	}
	// Check descendants (can assume they are IPanels)
	for _, child := range ipan.Children() {
		traverseIPanel(child.(IPanel), f)
	}
}

// traverseINode traverses the descendants of the specified INode,
// executing the specified function for each IPanel.
func traverseINode(inode core.INode, f func(ipan IPanel)) {
	if ipan, ok := inode.(IPanel); ok {
		traverseIPanel(ipan, f)
	} else {
		for _, child := range inode.Children() {
			traverseINode(child, f)
		}
	}
}

// forEachIPanel executes the specified function for each enabled and visible IPanel in gm.scene.
func (gm *Manager) forEachIPanel(f func(ipan IPanel)) {
	traverseINode(gm.scene, f)
}
