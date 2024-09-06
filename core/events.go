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

	WindowCursor
	WindowCursorEnter

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

type WindowCursorEnterEvent struct {
	Entered bool
}

func (e WindowCursorEnterEvent) WindowEventType() WindowEventType {
	return WindowCursorEnter
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

	GuiCursor
	GuiCursorEnter
	GuiCursorLeave

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

type KeyEvent interface {
	GetKey() Key
	GetMods() ModifierKey
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

func (e KeyUpEvent) GetKey() Key {
	return e.Key
}

func (e KeyUpEvent) GetMods() ModifierKey {
	return e.Mods
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

func (e KeyDownEvent) GetKey() Key {
	return e.Key
}

func (e KeyDownEvent) GetMods() ModifierKey {
	return e.Mods
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

func (e KeyRepeatEvent) GetKey() Key {
	return e.Key
}

func (e KeyRepeatEvent) GetMods() ModifierKey {
	return e.Mods
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

type CursorEvent struct {
	X float32
	Y float32
}

func (e CursorEvent) WindowEventType() WindowEventType {
	return WindowCursor
}

func (e CursorEvent) GuiEventType() GuiEventType {
	return GuiCursor
}

type GuiCursorEnterEvent struct{}

func (e GuiCursorEnterEvent) GuiEventType() GuiEventType {
	return GuiCursorEnter
}

type GuiCursorLeaveEvent struct{}

func (e GuiCursorLeaveEvent) GuiEventType() GuiEventType {
	return GuiCursorLeave
}

type MouseEvent interface {
	GetX() float32
	GetY() float32
	GetButton() MouseButton
	GetMods() ModifierKey
}

type MouseUpEvent struct {
	X      float32
	Y      float32
	Button MouseButton
	Mods   ModifierKey
}

func (e MouseUpEvent) WindowEventType() WindowEventType {
	return WindowMouseUp
}

func (e MouseUpEvent) GuiEventType() GuiEventType {
	return GuiMouseUp
}

func (e MouseUpEvent) GetX() float32 {
	return e.X
}

func (e MouseUpEvent) GetY() float32 {
	return e.Y
}

func (e MouseUpEvent) GetButton() MouseButton {
	return e.Button
}

func (e MouseUpEvent) GetMods() ModifierKey {
	return e.Mods
}

type MouseDownEvent struct {
	X      float32
	Y      float32
	Button MouseButton
	Mods   ModifierKey
}

func (e MouseDownEvent) WindowEventType() WindowEventType {
	return WindowMouseDown
}

func (e MouseDownEvent) GuiEventType() GuiEventType {
	return GuiMouseDown
}

func (e MouseDownEvent) GetX() float32 {
	return e.X
}

func (e MouseDownEvent) GetY() float32 {
	return e.Y
}

func (e MouseDownEvent) GetButton() MouseButton {
	return e.Button
}

func (e MouseDownEvent) GetMods() ModifierKey {
	return e.Mods
}

type ScrollEvent struct {
	X float32
	Y float32
}

func (e ScrollEvent) WindowEventType() WindowEventType {
	return WindowScroll
}

func (e ScrollEvent) GuiEventType() GuiEventType {
	return GuiScroll
}

type GuiClickEvent struct {
	X      float32
	Y      float32
	Button MouseButton
	Mods   ModifierKey
}

func (e GuiClickEvent) GuiEventType() GuiEventType {
	return GuiClick
}

func (e GuiClickEvent) GetX() float32 {
	return e.X
}

func (e GuiClickEvent) GetY() float32 {
	return e.Y
}

func (e GuiClickEvent) GetButton() MouseButton {
	return e.Button
}

func (e GuiClickEvent) GetMods() ModifierKey {
	return e.Mods
}
