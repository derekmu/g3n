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

	WindowMouseUp
	WindowMouseDown
	WindowScroll
)

type WindowEvent interface {
	WindowEventType() WindowEventType
}

type AppExitEvent struct{}

func (i AppExitEvent) WindowEventType() WindowEventType {
	return AppExit
}

type WindowFocusEvent struct {
	Focused bool
}

func (i WindowFocusEvent) WindowEventType() WindowEventType {
	return WindowFocus
}

type WindowPosEvent struct {
	X int
	Y int
}

func (i WindowPosEvent) WindowEventType() WindowEventType {
	return WindowPos
}

type WindowSizeEvent struct {
	Width  int
	Height int
}

func (i WindowSizeEvent) WindowEventType() WindowEventType {
	return WindowSize
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

func (i GuiResizeEvent) GuiEventType() GuiEventType {
	return GuiResize
}

type GuiEnableEvent struct {
	Enabled bool
}

func (i GuiEnableEvent) GuiEventType() GuiEventType {
	return GuiEnable
}

type GuiFocusEvent struct{}

func (i GuiFocusEvent) GuiEventType() GuiEventType {
	return GuiFocus
}

type GuiFocusLostEvent struct{}

func (i GuiFocusLostEvent) GuiEventType() GuiEventType {
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

func (i KeyUpEvent) WindowEventType() WindowEventType {
	return WindowKeyUp
}

func (i KeyUpEvent) GuiEventType() GuiEventType {
	return GuiKeyUp
}

func (i KeyUpEvent) GetKey() Key {
	return i.Key
}

func (i KeyUpEvent) GetMods() ModifierKey {
	return i.Mods
}

type KeyDownEvent struct {
	Key  Key
	Mods ModifierKey
}

func (i KeyDownEvent) WindowEventType() WindowEventType {
	return WindowKeyDown
}

func (i KeyDownEvent) GuiEventType() GuiEventType {
	return GuiKeyDown
}

func (i KeyDownEvent) GetKey() Key {
	return i.Key
}

func (i KeyDownEvent) GetMods() ModifierKey {
	return i.Mods
}

type KeyRepeatEvent struct {
	Key  Key
	Mods ModifierKey
}

func (i KeyRepeatEvent) WindowEventType() WindowEventType {
	return WindowKeyRepeat
}

func (i KeyRepeatEvent) GuiEventType() GuiEventType {
	return GuiKeyRepeat
}

func (i KeyRepeatEvent) GetKey() Key {
	return i.Key
}

func (i KeyRepeatEvent) GetMods() ModifierKey {
	return i.Mods
}

type CharEvent struct {
	Char rune
}

func (i CharEvent) WindowEventType() WindowEventType {
	return WindowChar
}

func (i CharEvent) GuiEventType() GuiEventType {
	return GuiChar
}

type CursorEvent struct {
	X float32
	Y float32
}

func (i CursorEvent) WindowEventType() WindowEventType {
	return WindowCursor
}

func (i CursorEvent) GuiEventType() GuiEventType {
	return GuiCursor
}

type GuiCursorEnterEvent struct{}

func (i GuiCursorEnterEvent) GuiEventType() GuiEventType {
	return GuiCursorEnter
}

type GuiCursorLeaveEvent struct{}

func (i GuiCursorLeaveEvent) GuiEventType() GuiEventType {
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

func (i MouseUpEvent) WindowEventType() WindowEventType {
	return WindowMouseUp
}

func (i MouseUpEvent) GuiEventType() GuiEventType {
	return GuiMouseUp
}

func (i MouseUpEvent) GetX() float32 {
	return i.X
}

func (i MouseUpEvent) GetY() float32 {
	return i.Y
}

func (i MouseUpEvent) GetButton() MouseButton {
	return i.Button
}

func (i MouseUpEvent) GetMods() ModifierKey {
	return i.Mods
}

type MouseDownEvent struct {
	X      float32
	Y      float32
	Button MouseButton
	Mods   ModifierKey
}

func (i MouseDownEvent) WindowEventType() WindowEventType {
	return WindowMouseDown
}

func (i MouseDownEvent) GuiEventType() GuiEventType {
	return GuiMouseDown
}

func (i MouseDownEvent) GetX() float32 {
	return i.X
}

func (i MouseDownEvent) GetY() float32 {
	return i.Y
}

func (i MouseDownEvent) GetButton() MouseButton {
	return i.Button
}

func (i MouseDownEvent) GetMods() ModifierKey {
	return i.Mods
}

type ScrollEvent struct {
	X float32
	Y float32
}

func (i ScrollEvent) WindowEventType() WindowEventType {
	return WindowScroll
}

func (i ScrollEvent) GuiEventType() GuiEventType {
	return GuiScroll
}

type GuiClickEvent struct {
	X      float32
	Y      float32
	Button MouseButton
	Mods   ModifierKey
}

func (i GuiClickEvent) GuiEventType() GuiEventType {
	return GuiClick
}

func (i GuiClickEvent) GetX() float32 {
	return i.X
}

func (i GuiClickEvent) GetY() float32 {
	return i.Y
}

func (i GuiClickEvent) GetButton() MouseButton {
	return i.Button
}

func (i GuiClickEvent) GetMods() ModifierKey {
	return i.Mods
}
