// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"github.com/derekmu/g3n/gui"
	"image"
	_ "image/png"
	"os"
	"runtime"

	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/gui/assets"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var _ gui.IWindow = &window{}

// window encapsulates a GLFW window.
type window struct {
	*glfw.Window
	core.Dispatcher[core.WindowEvent]
	gls        *gls.GLS
	fullscreen bool
	lastX      int
	lastY      int
	lastWidth  int
	lastHeight int
	scaleX     float64
	scaleY     float64

	cursors       map[core.Cursor]*glfw.Cursor
	lastCursorKey core.Cursor
}

// newWindow creates a new window.
func newWindow(width, height int, title string) (*window, error) {
	// OpenGL functions must be executed in the same thread where the context was created by glfw.CreateWindow().
	runtime.LockOSThread()

	// Create wrapper window with dispatcher
	w := new(window)
	var err error

	// Initialize GLFW
	err = glfw.Init()
	if err != nil {
		return nil, err
	}

	// Set window hints
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Samples, 8)
	// Set OpenGL forward compatible context only for OSX because it is required for OSX.
	// When this is set, glLineWidth(width) only accepts width=1.0 and generates an error
	// for any other values although the spec says it should ignore unsupported widths
	// and generate an error only when width <= 0.
	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	// Create window and set it as the current context.
	// The window is created always as not full screen because if it is
	// created as full screen it not possible to revert it to windowed mode.
	// At the end of this function, the window will be set to full screen if requested.
	w.Window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	w.MakeContextCurrent()

	// Create OpenGL state
	w.gls, err = gls.NewGLS()
	if err != nil {
		return nil, err
	}

	// Compute and store scale
	fbw, fbh := w.GetFramebufferSize()
	w.scaleX = float64(fbw) / float64(width)
	w.scaleY = float64(fbh) / float64(height)

	// Create map for cursors
	w.cursors = make(map[core.Cursor]*glfw.Cursor)
	w.lastCursorKey = core.CursorLast

	// Preallocate GLFW standard cursors
	w.cursors[core.ArrowCursor] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	w.cursors[core.IBeamCursor] = glfw.CreateStandardCursor(glfw.IBeamCursor)
	w.cursors[core.CrosshairCursor] = glfw.CreateStandardCursor(glfw.CrosshairCursor)
	w.cursors[core.HandCursor] = glfw.CreateStandardCursor(glfw.HandCursor)
	w.cursors[core.HResizeCursor] = glfw.CreateStandardCursor(glfw.HResizeCursor)
	w.cursors[core.VResizeCursor] = glfw.CreateStandardCursor(glfw.VResizeCursor)

	// Preallocate extra G3N standard cursors (diagonal resize cursors)
	trblImage, _, err := assets.NewCursorTrblImage()
	if err != nil {
		return nil, err
	}
	tlbrImage, _, err := assets.NewCursorTlbrImage()
	if err != nil {
		return nil, err
	}
	w.cursors[core.DiagResizeTrblCursor] = glfw.CreateCursor(trblImage, 8, 8)
	w.cursors[core.DiagResizeTlbrCursor] = glfw.CreateCursor(tlbrImage, 8, 8)

	w.SetKeyCallback(func(x *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			w.Dispatch(core.KeyDownEvent{
				Key:  core.Key(key),
				Mods: core.ModifierKey(mods),
			})
		case glfw.Release:
			w.Dispatch(core.KeyUpEvent{
				Key:  core.Key(key),
				Mods: core.ModifierKey(mods),
			})
		case glfw.Repeat:
			w.Dispatch(core.KeyRepeatEvent{
				Key:  core.Key(key),
				Mods: core.ModifierKey(mods),
			})
		}
	})
	w.SetCharCallback(func(x *glfw.Window, char rune) {
		w.Dispatch(core.CharEvent{Char: char})
	})
	w.SetMouseButtonCallback(func(x *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		xpos, ypos := x.GetCursorPos()
		switch action {
		case glfw.Press:
			w.Dispatch(core.MouseDownEvent{
				X:      float32(xpos),
				Y:      float32(ypos),
				Button: core.MouseButton(button),
				Mods:   core.ModifierKey(mods),
			})
		case glfw.Release:
			w.Dispatch(core.MouseUpEvent{
				X:      float32(xpos),
				Y:      float32(ypos),
				Button: core.MouseButton(button),
				Mods:   core.ModifierKey(mods),
			})
		}
	})
	w.SetScrollCallback(func(x *glfw.Window, xoff float64, yoff float64) {
		w.Dispatch(core.ScrollEvent{
			X: float32(xoff),
			Y: float32(yoff),
		})
	})
	w.SetCursorPosCallback(func(x *glfw.Window, xpos float64, ypos float64) {
		w.Dispatch(core.CursorEvent{
			X: float32(xpos),
			Y: float32(ypos),
		})
	})
	w.SetCursorEnterCallback(func(x *glfw.Window, entered bool) {
		w.Dispatch(core.WindowCursorEnterEvent{
			Entered: entered,
		})
	})
	w.SetSizeCallback(func(x *glfw.Window, width int, height int) {
		fbw, fbh := x.GetFramebufferSize()
		w.scaleX = float64(fbw) / float64(width)
		w.scaleY = float64(fbh) / float64(height)
		w.Dispatch(core.WindowSizeEvent{
			Width:  width,
			Height: height,
		})
	})
	w.SetPosCallback(func(x *glfw.Window, xpos int, ypos int) {
		w.Dispatch(core.WindowPosEvent{
			X: xpos,
			Y: ypos,
		})
	})
	w.SetFocusCallback(func(x *glfw.Window, focused bool) {
		w.Dispatch(core.WindowFocusEvent{Focused: focused})
	})

	return w, nil
}

// Gls returns the associated OpenGL state.
func (w *window) Gls() *gls.GLS {
	return w.gls
}

// FullScreen returns whether this window is currently fullscreen.
func (w *window) FullScreen() bool {
	return w.fullscreen
}

// SetFullScreen sets this window as fullscreen on the primary monitor.
func (w *window) SetFullScreen(full bool) {
	// If already in the desired state, nothing to do
	if w.fullscreen == full {
		return
	}
	// Set window fullscreen on the primary monitor
	if full {
		// Save current position and size of the window
		w.lastX, w.lastY = w.GetPos()
		w.lastWidth, w.lastHeight = w.GetSize()
		// Get size of primary monitor
		mon := glfw.GetPrimaryMonitor()
		vmode := mon.GetVideoMode()
		width := vmode.Width
		height := vmode.Height
		// Set as fullscreen on the primary monitor
		w.SetMonitor(mon, 0, 0, width, height, vmode.RefreshRate)
		w.fullscreen = true
	} else {
		// Restore window to previous position and size
		w.SetMonitor(nil, w.lastX, w.lastY, w.lastWidth, w.lastHeight, glfw.DontCare)
		w.fullscreen = false
	}
}

// Destroy destroys this window and its context
func (w *window) Destroy() {
	w.Window.Destroy()
	glfw.Terminate()
	runtime.UnlockOSThread() // Important when using the execution tracer
}

// GetScale returns this window's DPI scale factor (FramebufferSize / Size)
func (w *window) GetScale() (x float64, y float64) {
	return w.scaleX, w.scaleY
}

// ScreenResolution returns the screen resolution
func (w *window) ScreenResolution() (width, height int) {
	mon := glfw.GetPrimaryMonitor()
	vmode := mon.GetVideoMode()
	return vmode.Width, vmode.Height
}

// PollEvents process events in the event queue
func (w *window) PollEvents() {
	glfw.PollEvents()
}

// SetSwapInterval sets the number of screen updates to wait from the time SwapBuffer() is called before swapping the buffers and returning.
func (w *window) SetSwapInterval(interval int) {
	glfw.SwapInterval(interval)
}

// SetCursor sets the window's cursor.
func (w *window) SetCursor(cursor core.Cursor) {
	cur, ok := w.cursors[cursor]
	if !ok {
		panic("Invalid cursor")
	}
	w.Window.SetCursor(cur)
}

// CreateCursor creates a new custom cursor and returns an int handle.
func (w *window) CreateCursor(imgFile string, xhot, yhot int) (core.Cursor, error) {
	file, err := os.Open(imgFile)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	w.lastCursorKey += 1
	w.cursors[w.lastCursorKey] = glfw.CreateCursor(img, xhot, yhot)

	return w.lastCursorKey, nil
}
