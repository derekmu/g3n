// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

type WindowEventType int32

const (
	AppExit WindowEventType = iota
	WindowFocus
	WindowPos
	WindowSize

	WindowKeyUp
	WindowKeyDown
	WindowKeyRepeat
	WindowChar

	WindowMouse
	WindowMouseEnter
	WindowMouseUp
	WindowMouseDown
	WindowScroll
)

type WindowEvent interface {
	WindowEventType() WindowEventType
}

type AppExitEvent struct{}

func (e AppExitEvent) WindowEventType() WindowEventType {
	return AppExit
}

type WindowFocusEvent struct {
	Focused bool
}

func (e WindowFocusEvent) WindowEventType() WindowEventType {
	return WindowFocus
}

type WindowPosEvent struct {
	X int
	Y int
}

func (e WindowPosEvent) WindowEventType() WindowEventType {
	return WindowPos
}

type WindowSizeEvent struct {
	Width  int
	Height int
}

func (e WindowSizeEvent) WindowEventType() WindowEventType {
	return WindowSize
}

type WindowMouseEnterEvent struct {
	Entered bool
}

func (e WindowMouseEnterEvent) WindowEventType() WindowEventType {
	return WindowMouseEnter
}

type GuiEventType int32

const (
	GuiResize GuiEventType = iota
	GuiEnable
	GuiFocus
	GuiFocusLost

	GuiKeyUp
	GuiKeyDown
	GuiKeyRepeat
	GuiChar

	GuiMouse
	GuiMouseEnter
	GuiMouseLeave

	GuiMouseDown
	GuiMouseUp
	GuiScroll
	GuiClick
)

type GuiEvent interface {
	GuiEventType() GuiEventType
}

type GuiResizeEvent struct{}

func (e GuiResizeEvent) GuiEventType() GuiEventType {
	return GuiResize
}

type GuiEnableEvent struct {
	Enabled bool
}

func (e GuiEnableEvent) GuiEventType() GuiEventType {
	return GuiEnable
}

type GuiFocusEvent struct{}

func (e GuiFocusEvent) GuiEventType() GuiEventType {
	return GuiFocus
}

type GuiFocusLostEvent struct{}

func (e GuiFocusLostEvent) GuiEventType() GuiEventType {
	return GuiFocusLost
}

type KeyUpEvent struct {
	Key  Key
	Mods ModifierKey
}

func (e KeyUpEvent) WindowEventType() WindowEventType {
	return WindowKeyUp
}

func (e KeyUpEvent) GuiEventType() GuiEventType {
	return GuiKeyUp
}

type KeyDownEvent struct {
	Key  Key
	Mods ModifierKey
}

func (e KeyDownEvent) WindowEventType() WindowEventType {
	return WindowKeyDown
}

func (e KeyDownEvent) GuiEventType() GuiEventType {
	return GuiKeyDown
}

type KeyRepeatEvent struct {
	Key  Key
	Mods ModifierKey
}

func (e KeyRepeatEvent) WindowEventType() WindowEventType {
	return WindowKeyRepeat
}

func (e KeyRepeatEvent) GuiEventType() GuiEventType {
	return GuiKeyRepeat
}

type CharEvent struct {
	Char rune
}

func (e CharEvent) WindowEventType() WindowEventType {
	return WindowChar
}

func (e CharEvent) GuiEventType() GuiEventType {
	return GuiChar
}

type MouseEvent struct {
	X float64
	Y float64
}

func (e MouseEvent) WindowEventType() WindowEventType {
	return WindowMouse
}

func (e MouseEvent) GuiEventType() GuiEventType {
	return GuiMouse
}

type GuiMouseEnterEvent struct{}

func (e GuiMouseEnterEvent) GuiEventType() GuiEventType {
	return GuiMouseEnter
}

type GuiMouseLeaveEvent struct{}

func (e GuiMouseLeaveEvent) GuiEventType() GuiEventType {
	return GuiMouseLeave
}

type MouseUpEvent struct {
	X      float64
	Y      float64
	Button MouseButton
	Mods   ModifierKey
}

func (e MouseUpEvent) WindowEventType() WindowEventType {
	return WindowMouseUp
}

func (e MouseUpEvent) GuiEventType() GuiEventType {
	return GuiMouseUp
}

type MouseDownEvent struct {
	X      float64
	Y      float64
	Button MouseButton
	Mods   ModifierKey
}

func (e MouseDownEvent) WindowEventType() WindowEventType {
	return WindowMouseDown
}

func (e MouseDownEvent) GuiEventType() GuiEventType {
	return GuiMouseDown
}

type ScrollEvent struct {
	X float64
	Y float64
}

func (e ScrollEvent) WindowEventType() WindowEventType {
	return WindowScroll
}

func (e ScrollEvent) GuiEventType() GuiEventType {
	return GuiScroll
}

type GuiClickEvent struct {
	X      float64
	Y      float64
	Button MouseButton
	Mods   ModifierKey
}

func (e GuiClickEvent) GuiEventType() GuiEventType {
	return GuiClick
}
