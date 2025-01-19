// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/geometry"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/texture"
)

// Skybox is the Graphic that represents a skybox.
type Skybox struct {
	Graphic
	uniMatrices gls.Uniform
}

// SkyboxData contains the data necessary to locate the textures for a Skybox in a concise manner.
type SkyboxData struct {
	DirAndPrefix string
	Extension    string
	Suffixes     [6]string
}

// NewSkybox creates and returns a pointer to a Skybox with the specified textures.
func NewSkybox(data SkyboxData) (*Skybox, error) {
	s := new(Skybox)

	geom := geometry.NewCube(1)
	s.Graphic.Init(s, geom, gls.TRIANGLES)
	s.Graphic.SetCullable(false)

	for i := 0; i < 6; i++ {
		tex, err := texture.NewTexture2DFromImage(data.DirAndPrefix + data.Suffixes[i] + "." + data.Extension)
		if err != nil {
			return nil, err
		}
		matFace := material.NewStandard(math32.Color{R: 1, G: 1, B: 1})
		matFace.AddTexture(tex)
		matFace.SetSide(material.SideBack)
		matFace.SetUseLights(material.UseLightNone)

		// Disable writes to the depth buffer (call glDepthMask(GL_FALSE)).
		// This will cause every other object to draw over the skybox, making it always appear behind everything else.
		// It doesn't matter how small/big the skybox is as long as it's visible by the camera (within near/far planes).
		matFace.SetDepthMask(false)

		s.AddGroupMaterial(s, matFace, i)
	}

	// Creates uniforms
	s.uniMatrices.Init("uMatrices")

	// The skybox should always be rendered last among the opaque objects
	s.SetRenderOrder(100)

	return s, nil
}

// RenderSetup is called by the engine before drawing the skybox geometry.
// It is responsible for updating the current shader uniforms with the model matrices.
func (s *Skybox) RenderSetup(gs *gls.GLS, _ *core.RenderInfo) {
	// Clear translation
	s.mdata.mvm[12] = 0
	s.mdata.mvm[13] = 0
	s.mdata.mvm[14] = 0

	// Calculates normal matrix and transfer uniform
	var nm math32.Matrix3
	_ = nm.GetNormalMatrix(&s.mdata.mvm)
	s.mdata.nm.SetFromMatrix3(&nm)
	gs.UniformMatrix4fv(s.uniMatrices.Location(gs), 3, false, &s.mdata.mvm[0])
}
