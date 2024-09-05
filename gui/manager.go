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
	if gm != nil {
		return gm
	}
	panic(fmt.Errorf("gui.InitManager not called"))
}

// IWindow is the interface that the manager uses to interact with the window.
type IWindow interface {
	core.IDispatcher[core.WindowEvent]
	GetScale() (x float64, y float64)
	SetCursor(cursor core.Cursor)
}

// Manager routes events to the appropriate GUI components or outside the GUI if not applicable.
type Manager struct {
	core.Dispatcher[core.GuiEvent]
	core.TimerManager
	window      IWindow
	scene       core.INode
	target      IPanel
	keyFocus    core.IDispatcher[core.GuiEvent]
	cursorFocus core.IDispatcher[core.GuiEvent]
}

// InitManager creates the Manager singleton or panics if it's already been called.
func InitManager(window IWindow) {
	if gm != nil {
		panic(fmt.Errorf("gui.InitManager already called"))
	}
	gm = new(Manager)
	gm.TimerManager.Initialize()
	gm.window = window
	window.Subscribe(gm.onWindowEvent)
}

// SetScene sets the INode to watch for events.
func (gm *Manager) SetScene(scene core.INode) {
	gm.scene = scene
}

// SetKeyFocus sets the key-focused IDispatcher, which will exclusively receive key and char events.
func (gm *Manager) SetKeyFocus(disp core.IDispatcher[core.GuiEvent]) {
	if gm.keyFocus == disp {
		return
	}
	if gm.keyFocus != nil {
		gm.keyFocus.Dispatch(core.GuiFocusLostEvent{})
	}
	gm.keyFocus = disp
	if gm.keyFocus != nil {
		gm.keyFocus.Dispatch(core.GuiFocusEvent{})
	}
}

// SetCursorFocus sets the cursor-focused IDispatcher, which will exclusively receive cursor events.
func (gm *Manager) SetCursorFocus(disp core.IDispatcher[core.GuiEvent]) {
	if gm.cursorFocus == disp {
		return
	}
	gm.cursorFocus = disp
}

// onKeyEvent is called when char or key events are received.
func (gm *Manager) onKeyEvent(ev core.GuiEvent) {
	if gm.keyFocus != nil {
		gm.keyFocus.Dispatch(ev)
	} else {
		gm.Dispatch(ev)
	}
}

// onMouse is called when mouse events are received.
func (gm *Manager) onMouse(ev core.GuiEvent) {
	if gm.scene != nil && gm.target != nil {
		sendAncestry(gm.target, false, nil, ev)
	} else {
		gm.Dispatch(ev)
	}
}

// onScroll is called when scroll events are received.
func (gm *Manager) onScroll(ev core.ScrollEvent) {
	if gm.scene != nil && gm.target != nil {
		sendAncestry(gm.target, false, nil, ev)
	} else {
		gm.Dispatch(ev)
	}
}

func (gm *Manager) updateMouseTarget(x, y float32) {
	oldTarget := gm.target
	gm.target = nil
	// Find IPanel immediately under the cursor and store it in gm.target
	gm.forEachIPanel(func(ipan IPanel) {
		if ipan.InsideBorders(x, y) && (gm.target == nil || ipan.Position().Z < gm.target.GetPanel().Position().Z) {
			gm.target = ipan
		}
	})
	if gm.target != oldTarget {
		// Only send events up to the lowest common ancestor of target and oldTarget
		var commonAnc IPanel
		if gm.target != nil && oldTarget != nil {
			commonAnc, _ = gm.target.LowestCommonAncestor(oldTarget).(IPanel)
		}
		if oldTarget != nil && !oldTarget.IsAncestorOf(gm.target) {
			sendAncestry(oldTarget, true, commonAnc, core.GuiCursorLeaveEvent{})
		}
		if gm.target != nil && !gm.target.IsAncestorOf(oldTarget) {
			sendAncestry(gm.target, true, commonAnc, core.GuiCursorEnterEvent{})
		}
	}
}

// onCursor is called when cursor events are received.
func (gm *Manager) onCursor(ev core.CursorEvent) {
	if gm.cursorFocus != nil {
		gm.cursorFocus.Dispatch(ev)
		return
	}
	if gm.scene == nil {
		gm.Dispatch(ev)
		return
	}
	gm.updateMouseTarget(ev.X, ev.Y)
	if gm.target != nil {
		sendAncestry(gm.target, false, nil, ev)
	} else {
		gm.Dispatch(ev)
	}
}

// sendAncestry sends the specified event to the specified target panel and its ancestors.
// If all is false, only send to the lowest subscribed ancestor.
// If upToExclude is not nil then the event will not be dispatched to that panel or it's ancestors.
// If upToInclude is not nil then the event will be dispatched to that panel but not it's ancestors.
func sendAncestry(ipan IPanel, all bool, upToExclude IPanel, ev core.GuiEvent) {
	var ok bool
	for ipan != nil {
		if upToExclude != nil && ipan == upToExclude {
			break
		}
		count := ipan.Dispatch(ev)
		if !all && count > 0 {
			break
		}
		ipan, ok = ipan.Parent().(IPanel)
		if !ok {
			break
		}
	}
}

// traverseIPanel traverses the descendants of the provided IPanel, executing the specified function for each IPanel.
func traverseIPanel(ipan IPanel, f func(ipan IPanel)) {
	if !ipan.Visible() {
		return
	}
	if ipan.Enabled() {
		f(ipan)
	}
	for _, child := range ipan.Children() {
		traverseIPanel(child.(IPanel), f)
	}
}

// traverseINode traverses the descendants of the specified INode, executing the specified function for each IPanel.
func traverseINode(inode core.INode, f func(ipan IPanel)) {
	if ipan, ok := inode.(IPanel); ok {
		traverseIPanel(ipan, f)
	} else {
		for _, child := range inode.Children() {
			traverseINode(child, f)
		}
	}
}

// forEachIPanel executes the specified function for each enabled and visible IPanel in the scene.
func (gm *Manager) forEachIPanel(f func(ipan IPanel)) {
	traverseINode(gm.scene, f)
}

func (gm *Manager) onWindowEvent(event core.WindowEvent) bool {
	switch ev := event.(type) {
	case core.KeyUpEvent:
		gm.onKeyEvent(ev)
	case core.KeyDownEvent:
		gm.onKeyEvent(ev)
	case core.KeyRepeatEvent:
		gm.onKeyEvent(ev)
	case core.CharEvent:
		gm.onKeyEvent(ev)
	case core.CursorEvent:
		gm.onCursor(ev)
	case core.MouseUpEvent:
		gm.onMouse(ev)
	case core.MouseDownEvent:
		gm.onMouse(ev)
	case core.ScrollEvent:
		gm.onScroll(ev)
	default:
		return false
	}
	return true
}
