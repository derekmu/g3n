// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/geometry"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/material"
)

// Lines is a Graphic which is rendered as a collection of independent lines.
type Lines struct {
	Graphic
	uniMatrices gls.Uniform
}

// NewLines returns a pointer to a new Lines object.
func NewLines(igeom geometry.IGeometry, imat material.IMaterial) *Lines {
	l := new(Lines)
	l.Init(igeom, imat)
	return l
}

// Init initializes the Lines object and adds the specified material.
func (l *Lines) Init(igeom geometry.IGeometry, imat material.IMaterial) {
	l.InitGraphic(l, igeom, gls.LINES)
	l.AddMaterial(l, imat, 0, 0)
	l.uniMatrices.Init("uMatrices")
}

// RenderSetup is called by the engine before drawing this geometry.
func (l *Lines) RenderSetup(gs *gls.GLS, _ *core.RenderInfo) {
	// Transfer model view projection matrix uniform
	gs.UniformMatrix4fv(l.uniMatrices.Location(gs), 3, false, &l.mdata.mvm[0])
}
