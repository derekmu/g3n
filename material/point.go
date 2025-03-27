// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/derekmu/g3n/math32"
)

// Point material is normally used for single point sprites.
type Point struct {
	Standard
}

// NewPoint creates a new Point material.
func NewPoint(color math32.Color) *Point {
	m := new(Point)
	m.InitPoint(color)
	m.SetShader("point")
	return m
}

// InitPoint initializes the material setting the specified color.
func (m *Point) InitPoint(color math32.Color) {
	m.InitStandard(color)
	m.udata.emissive = color.Color3()
	m.udata.psize = 1.0
	m.udata.protationZ = 0
}

// SetEmissiveColor sets the material's emissive color.
func (m *Point) SetEmissiveColor(color math32.Color) {
	m.udata.emissive = color.Color3()
}

// SetSize sets the point size.
func (m *Point) SetSize(size float32) {
	m.udata.psize = size
}

// SetRotationZ sets the point rotation around the Z axis.
func (m *Point) SetRotationZ(rot float32) {
	m.udata.protationZ = rot
}
