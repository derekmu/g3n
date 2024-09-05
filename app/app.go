// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package app implements a cross-platform G3N app.
package app

import (
	"fmt"
	"github.com/derekmu/g3n/audio/al"
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gui"
	"github.com/derekmu/g3n/renderer"
	"log"
	"time"
)

// Application is the main overall G3N container.
type Application struct {
	*window
	renderer   *renderer.Renderer // Renderer object
	audioDev   *al.Device         // Default audio device
	startTime  time.Time          // Application start time
	frameStart time.Time          // Frame start time
	frameDelta time.Duration      // Duration of last frame
}

// NewApplication creates a new Application.
func NewApplication(width, height int, title string) *Application {
	app := new(Application)
	win, err := newWindow(width, height, title)
	if err != nil {
		panic(err)
	}
	app.window = win
	app.renderer = renderer.NewRenderer(win.Gls())
	err = app.renderer.AddDefaultShaders()
	if err != nil {
		panic(fmt.Errorf("AddDefaultShaders:%v", err))
	}
	gui.InitManager(win)
	return app
}

// Run starts the update loop and calls the user-provided update function every frame.
func (a *Application) Run(update func(rend *renderer.Renderer, deltaTime time.Duration)) {
	a.startTime = time.Now()
	a.frameStart = time.Now()
	for {
		if a.ShouldClose() {
			a.Dispatch(core.AppExitEvent{})
			break
		}
		now := time.Now()
		a.frameDelta = now.Sub(a.frameStart)
		a.frameStart = now
		update(a.renderer, a.frameDelta)
		a.window.SwapBuffers()
		a.window.PollEvents()
	}
	if a.audioDev != nil {
		err := al.CloseDevice(a.audioDev)
		if err != nil {
			log.Print("failed to close audio device", err)
		}
	}
	a.Destroy()
}

// Exit requests to terminate the application.
func (a *Application) Exit() {
	a.window.SetShouldClose(true)
}

// Renderer returns the application's renderer.
func (a *Application) Renderer() *renderer.Renderer {
	return a.renderer
}

// RunTime returns the elapsed duration since the call to Run().
func (a *Application) RunTime() time.Duration {
	return time.Since(a.startTime)
}

// OpenDefaultAudioDevice opens the default audio device setting it to the current context.
func (a *Application) OpenDefaultAudioDevice() error {
	// Opens default audio device
	var err error
	a.audioDev, err = al.OpenDevice("")
	if err != nil {
		return fmt.Errorf("opening OpenAL default device: %s", err)
	}
	// Check for OpenAL effects extension support
	var attribs []int
	if al.IsExtensionPresent("ALC_EXT_EFX") {
		attribs = []int{al.MAX_AUXILIARY_SENDS, 4}
	}
	// Create audio context
	acx, err := al.CreateContext(a.audioDev, attribs)
	if err != nil {
		return fmt.Errorf("creating OpenAL context: %s", err)
	}
	// Makes the context the current one
	err = al.MakeContextCurrent(acx)
	if err != nil {
		return fmt.Errorf("setting OpenAL context current: %s", err)
	}
	return nil
}
