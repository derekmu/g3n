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

// LineStrip is a Graphic which is rendered as a collection of connected lines.
type LineStrip struct {
	Graphic
	uniMatrices gls.Uniform
}

// NewLineStrip creates LineStrip graphic with the specified geometry and material.
func NewLineStrip(igeom geometry.IGeometry, imat material.IMaterial) *LineStrip {
	l := new(LineStrip)
	l.InitGraphic(l, igeom, gls.LINE_STRIP)
	l.AddMaterial(l, imat, 0, 0)
	l.uniMatrices.Init("uMatrices")
	return l
}

// RenderSetup is called by the engine before drawing this geometry.
func (l *LineStrip) RenderSetup(gs *gls.GLS, _ *core.RenderInfo) {
	// Transfer model view projection matrix uniform
	gs.UniformMatrix4fv(l.uniMatrices.Location(gs), 3, false, &l.mdata.mvm[0])
}
