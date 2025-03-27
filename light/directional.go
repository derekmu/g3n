// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
)

// Directional represents a directional light.
type Directional struct {
	core.Node
	uni   gls.Uniform
	udata struct {
		color    math32.Color3
		position math32.Vector3
	}
}

// NewDirectional returns a new directional light.
func NewDirectional(color math32.Color) *Directional {
	ld := new(Directional)
	ld.InitNode(ld)

	ld.uni.Init("uDirLight")
	ld.SetColor(color)
	return ld
}

// SetColor sets the color of this light.
func (ld *Directional) SetColor(color math32.Color) {
	ld.udata.color = color.Color3()
}

// Color returns the color of this light.
func (ld *Directional) Color() math32.Color {
	return ld.udata.color
}

// RenderSetup is called by the engine before rendering the scene
func (ld *Directional) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {
	pos := ld.WorldPosition()
	pos4 := math32.Vector4{X: pos.X, Y: pos.Y, Z: pos.Z}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	ld.udata.position.X = pos4.X
	ld.udata.position.Y = pos4.Y
	ld.udata.position.Z = pos4.Z

	const vec3count = 2
	location := ld.uni.LocationIdx(gs, vec3count*int32(idx))
	gs.Uniform3fv(location, vec3count, &ld.udata.color.R)
}
