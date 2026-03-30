package gui

import (
	"github.com/derekmu/g3n/core"
)

// MouseTracker is a utility for panel mouse tracking.
//
// # Tracks whether the mouse is over the Panel and
//
// Dispatches GuiClickEvents based on MouseUpEvent & MouseDownEvents.
type MouseTracker struct {
	Over    bool            // Whether the mouse is currently over the panel.
	Pressed core.MouseState // Which mouse buttons are currently pressed over the panel.
	panel   IPanel
}

// Init sets the panel for to mouse events.
func (m *MouseTracker) Init(panel IPanel) {
	m.panel = panel
	panel.Subscribe(m.onGuiEvent)
}

// onGuiEvent is the callback for events on the panel.
func (m *MouseTracker) onGuiEvent(event core.GuiEvent) bool {
	switch ev := event.(type) {
	case core.MouseUpEvent:
		if m.panel.Enabled() {
			clicked := m.Pressed.IsSet(ev.Button)
			m.Pressed = m.Pressed.Unset(ev.Button)
			if clicked {
				m.panel.Dispatch(core.GuiClickEvent{
					X:      ev.X,
					Y:      ev.Y,
					Button: ev.Button,
					Mods:   ev.Mods,
				})
			}
		}
	case core.MouseDownEvent:
		if m.panel.Enabled() {
			m.Pressed = m.Pressed.Set(ev.Button)
		}
	case core.GuiMouseEnterEvent:
		m.Over = true
	case core.GuiMouseLeaveEvent:
		m.Over = false
		// Dragging out cancels clicks
		m.Pressed = 0
	case core.GuiEnableEvent:
		// Enabling or disabling a button cancels clicks
		m.Pressed = 0
	default:
		return false
	}
	return true
}
