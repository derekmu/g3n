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

// Points represents a geometry containing only points
type Points struct {
	Graphic
	uniMatrices gls.Uniform
}

// NewPoints creates Points object with the specified geometry and material.
func NewPoints(igeom geometry.IGeometry, imat material.IMaterial) *Points {
	p := new(Points)
	p.InitGraphic(p, igeom, gls.POINTS)
	if imat != nil {
		p.AddMaterial(p, imat, 0, 0)
	}
	p.uniMatrices.Init("uMatrices")
	return p
}

// RenderSetup is called by the engine before rendering this graphic.
func (p *Points) RenderSetup(gs *gls.GLS, _ *core.RenderInfo) {
	// Transfer model view projection matrix uniform
	gs.UniformMatrix4fv(p.uniMatrices.Location(gs), 3, false, &p.mdata.mvm[0])
}
