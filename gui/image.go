// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/texture"
)

// Image is a Panel which resizes to match the texture size.
type Image struct {
	Panel
}

// NewImage creates a new Image.
func NewImage() *Image {
	i := new(Image)
	i.InitImage()
	return i
}

func (i *Image) InitImage() {
	i.Panel.InitPanel(i, 0, 0)
}

// SetTexture changes the image's texture and resizes the panel.
// It returns a pointer to the previous texture.
func (i *Image) SetTexture(tex *texture.Texture2D) *texture.Texture2D {
	prevtex := i.Panel.SetTexture(tex)
	if tex != nil {
		i.Panel.SetContentSize(float32(i.tex.Width()), float32(i.tex.Height()))
	} else {
		i.Panel.SetContentSize(0, 0)
	}
	return prevtex
}
