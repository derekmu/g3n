// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/texture"
)

// Image is a Panel which resizes to match the texture size by default.
type Image struct {
	Panel
	resize bool
}

// NewImage creates a new Image.
func NewImage() *Image {
	i := new(Image)
	i.InitImage()
	return i
}

// InitImage initializes an image.
func (i *Image) InitImage() {
	i.Panel.InitPanel(i, 0, 0)
	i.resize = true
}

// SetResize sets whether this image resizes automatically to match its texture size.
func (i *Image) SetResize(resize bool) {
	i.resize = resize
}

// GetResize returns whether this image resizes automatically to match its texture size.
func (i *Image) GetResize(resize bool) bool {
	return resize
}

// SetTexture changes the image's texture and resizes the panel.
// It returns a pointer to the previous texture.
func (i *Image) SetTexture(tex *texture.Texture2D) *texture.Texture2D {
	if i.resize {
		if tex != nil {
			i.Panel.SetContentSize(float32(tex.Width()), float32(tex.Height()))
		} else {
			i.Panel.SetContentSize(0, 0)
		}
	}
	return i.Panel.SetTexture(tex)
}
