// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
)

// Ambient represents an ambient light.
type Ambient struct {
	core.Node
	color math32.Color3
	uni   gls.Uniform
}

// NewAmbient returns a pointer to a new ambient color with the specified color.
func NewAmbient(color math32.Color) *Ambient {
	la := new(Ambient)
	la.InitNode(la)
	la.color = color.ToColor3()
	la.uni.Init("uAmbientLightColor")
	return la
}

// SetColor sets the color of this light.
func (la *Ambient) SetColor(color math32.Color) {
	la.color = color.ToColor3()
}

// Color returns the current color of this light.
func (la *Ambient) Color() math32.Color {
	return la.color
}

// RenderSetup is called by the engine before rendering the scene.
func (la *Ambient) RenderSetup(gs *gls.GLS, _ *core.RenderInfo, idx int) {
	location := la.uni.LocationIdx(gs, int32(idx))
	gs.Uniform3f(location, la.color.R, la.color.G, la.color.B)
}
