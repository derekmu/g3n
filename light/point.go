// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
)

// Point is an omnidirectional light source.
type Point struct {
	core.Node
	uni   gls.Uniform
	udata struct {
		color          math32.Color3
		position       math32.Vector3
		linearDecay    float32
		quadraticDecay float32
		_              float32
	}
}

// NewPoint returns a point light.
func NewPoint(color math32.Color) *Point {
	lp := new(Point)
	lp.InitNode(lp)

	lp.uni.Init("uPointLight")
	lp.SetColor(color)
	lp.SetLinearDecay(1.0)
	lp.SetQuadraticDecay(1.0)
	return lp
}

// SetColor sets the color of this light.
func (lp *Point) SetColor(color math32.Color) {
	lp.udata.color = color.ToColor3()
}

// Color returns the color of this light.
func (lp *Point) Color() math32.Color {
	return lp.udata.color
}

// SetLinearDecay sets the linear decay factor.
func (lp *Point) SetLinearDecay(decay float32) {
	lp.udata.linearDecay = decay
}

// LinearDecay returns the current linear decay factor.
func (lp *Point) LinearDecay() float32 {
	return lp.udata.linearDecay
}

// SetQuadraticDecay sets the quadratic decay factor.
func (lp *Point) SetQuadraticDecay(decay float32) {
	lp.udata.quadraticDecay = decay
}

// QuadraticDecay returns the current quadratic decay factor.
func (lp *Point) QuadraticDecay() float32 {
	return lp.udata.quadraticDecay
}

// RenderSetup is called by the engine before rendering the scene.
func (lp *Point) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {
	pos := lp.WorldPosition()
	pos4 := math32.Vector4{X: pos.X, Y: pos.Y, Z: pos.Z, W: 1.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	lp.udata.position.X = pos4.X
	lp.udata.position.Y = pos4.Y
	lp.udata.position.Z = pos4.Z

	const vec3count = 3
	location := lp.uni.LocationIdx(gs, vec3count*int32(idx))
	gs.Uniform3fv(location, vec3count, &lp.udata.color.R)
}
