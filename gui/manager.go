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
	mouseTarget IPanel
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
func (m *Manager) SetScene(scene core.INode) {
	m.scene = scene
}

// SetKeyFocus sets the key-focused IDispatcher, which will exclusively receive key and char events.
func (m *Manager) SetKeyFocus(disp core.IDispatcher[core.GuiEvent]) {
	if m.keyFocus == disp {
		return
	}
	if m.keyFocus != nil {
		m.keyFocus.Dispatch(core.GuiFocusLostEvent{})
	}
	m.keyFocus = disp
	if m.keyFocus != nil {
		m.keyFocus.Dispatch(core.GuiFocusEvent{})
	}
}

// SetCursorFocus sets the cursor-focused IDispatcher, which will exclusively receive cursor events.
func (m *Manager) SetCursorFocus(disp core.IDispatcher[core.GuiEvent]) {
	if m.cursorFocus == disp {
		return
	}
	m.cursorFocus = disp
}

// onKeyEvent is called when char or key events are received.
func (m *Manager) onKeyEvent(ev core.GuiEvent) {
	if m.keyFocus != nil {
		m.keyFocus.Dispatch(ev)
	} else {
		m.Dispatch(ev)
	}
}

// onMouse is called when mouse events are received.
func (m *Manager) onMouse(ev core.GuiEvent) {
	if m.scene != nil && m.mouseTarget != nil {
		sendAncestry(m.mouseTarget, false, nil, ev)
	} else {
		m.Dispatch(ev)
	}
}

// onScroll is called when scroll events are received.
func (m *Manager) onScroll(ev core.ScrollEvent) {
	if m.scene != nil && m.mouseTarget != nil {
		sendAncestry(m.mouseTarget, false, nil, ev)
	} else {
		m.Dispatch(ev)
	}
}

func (m *Manager) updateMouseTarget(x, y float64) {
	oldTarget := m.mouseTarget
	m.mouseTarget = nil
	// Find IPanel immediately under the cursor and store it in gm.target
	m.forEachIPanel(func(ipan IPanel) bool {
		if ipan.ContainsMouse(x, y) {
			if m.mouseTarget == nil || ipan.Position().Z < m.mouseTarget.Position().Z {
				m.mouseTarget = ipan
			}
			return true
		} else {
			// don't process children if the parent doesn't contain the mouse location
			return false
		}
	})
	if m.mouseTarget != oldTarget {
		// Only send events up to the lowest common ancestor of target and oldTarget
		var commonAnc IPanel
		if m.mouseTarget != nil && oldTarget != nil {
			commonAnc, _ = m.mouseTarget.LowestCommonAncestor(oldTarget).(IPanel)
		}
		if oldTarget != nil && !oldTarget.IsAncestorOf(m.mouseTarget) {
			sendAncestry(oldTarget, true, commonAnc, core.GuiCursorLeaveEvent{})
		}
		if m.mouseTarget != nil && !m.mouseTarget.IsAncestorOf(oldTarget) {
			sendAncestry(m.mouseTarget, true, commonAnc, core.GuiCursorEnterEvent{})
		}
	}
}

// onCursor is called when cursor events are received.
func (m *Manager) onCursor(ev core.CursorEvent) {
	if m.cursorFocus != nil {
		m.cursorFocus.Dispatch(ev)
		return
	}
	if m.scene == nil {
		m.Dispatch(ev)
		return
	}
	m.updateMouseTarget(ev.X, ev.Y)
	if m.mouseTarget != nil {
		sendAncestry(m.mouseTarget, false, nil, ev)
	} else {
		m.Dispatch(ev)
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
func traverseIPanel(ipan IPanel, f func(ipan IPanel) bool) {
	if !ipan.Visible() {
		return
	}
	if ipan.Enabled() && !f(ipan) {
		return
	}
	for _, child := range ipan.Children() {
		traverseIPanel(child.(IPanel), f)
	}
}

// traverseINode traverses the descendants of the specified INode, executing the specified function for each IPanel.
func traverseINode(inode core.INode, f func(ipan IPanel) bool) {
	if ipan, ok := inode.(IPanel); ok {
		traverseIPanel(ipan, f)
	} else {
		for _, child := range inode.Children() {
			traverseINode(child, f)
		}
	}
}

// forEachIPanel executes the specified function for each enabled and visible IPanel in the scene.
func (m *Manager) forEachIPanel(f func(ipan IPanel) bool) {
	traverseINode(m.scene, f)
}

func (m *Manager) onWindowEvent(event core.WindowEvent) bool {
	switch ev := event.(type) {
	case core.KeyUpEvent:
		m.onKeyEvent(ev)
	case core.KeyDownEvent:
		m.onKeyEvent(ev)
	case core.KeyRepeatEvent:
		m.onKeyEvent(ev)
	case core.CharEvent:
		m.onKeyEvent(ev)
	case core.CursorEvent:
		m.onCursor(ev)
	case core.MouseUpEvent:
		m.onMouse(ev)
	case core.MouseDownEvent:
		m.onMouse(ev)
	case core.ScrollEvent:
		m.onScroll(ev)
	default:
		return false
	}
	return true
}
