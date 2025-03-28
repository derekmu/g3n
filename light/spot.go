// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
)

// Spot represents a spotlight.
type Spot struct {
	core.Node
	uni   gls.Uniform
	udata struct {
		color          math32.Color3
		position       math32.Vector3
		direction      math32.Vector3
		angularDecay   float32
		cutoffAngle    float32
		linearDecay    float32
		quadraticDecay float32
		_              float32
		_              float32
	}
}

// NewSpot returns a new spotlight.
func NewSpot(color math32.Color) *Spot {
	l := new(Spot)
	l.InitNode(l)
	l.uni.Init("uSpotLight")
	l.SetColor(color)
	l.SetAngularDecay(15.0)
	l.SetCutoffAngle(45.0)
	l.SetLinearDecay(1.0)
	l.SetQuadraticDecay(1.0)
	return l
}

// SetColor sets the color of this light.
func (l *Spot) SetColor(color math32.Color) {
	l.udata.color = color.ToColor3()
}

// Color returns the color of this light.
func (l *Spot) Color() math32.Color {
	return l.udata.color
}

// SetCutoffAngle sets the cutoff angle in degrees from 0 to 90.
func (l *Spot) SetCutoffAngle(angle float32) {
	l.udata.cutoffAngle = angle
}

// CutoffAngle returns the cutoff angle in degrees from 0 to 90.
func (l *Spot) CutoffAngle() float32 {
	return l.udata.cutoffAngle
}

// SetAngularDecay sets the angular decay exponent.
func (l *Spot) SetAngularDecay(decay float32) {
	l.udata.angularDecay = decay
}

// AngularDecay returns the angular decay exponent.
func (l *Spot) AngularDecay() float32 {
	return l.udata.angularDecay
}

// SetLinearDecay sets the linear decay factor.
func (l *Spot) SetLinearDecay(decay float32) {
	l.udata.linearDecay = decay
}

// LinearDecay returns the linear decay factor.
func (l *Spot) LinearDecay() float32 {
	return l.udata.linearDecay
}

// SetQuadraticDecay sets the quadratic decay factor.
func (l *Spot) SetQuadraticDecay(decay float32) {
	l.udata.quadraticDecay = decay
}

// QuadraticDecay returns the quadratic decay factor.
func (l *Spot) QuadraticDecay() float32 {
	return l.udata.quadraticDecay
}

// RenderSetup is called by the engine before rendering the scene.
func (l *Spot) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {
	pos := l.WorldPosition()
	pos4 := math32.Vector4{X: pos.X, Y: pos.Y, Z: pos.Z, W: 1.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	l.udata.position.X = pos4.X
	l.udata.position.Y = pos4.Y
	l.udata.position.Z = pos4.Z

	dir := l.WorldDirection()
	pos4.SetVector3(dir, 0.0)
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	l.udata.direction.X = pos4.X
	l.udata.direction.Y = pos4.Y
	l.udata.direction.Z = pos4.Z

	const vec3count = 5
	location := l.uni.LocationIdx(gs, vec3count*int32(idx))
	gs.Uniform3fv(location, vec3count, &l.udata.color.R)
}
