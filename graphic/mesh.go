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
	"log"
)

// Mesh is a Graphic with uniforms for the model, view, projection, and normal matrices.
type Mesh struct {
	Graphic
	uniMatrices gls.Uniform
	skeleton    *Skeleton
	uniBones    gls.Uniform
}

// NewMesh creates and returns a pointer to a mesh with the specified geometry and material.
// If the mesh has multi materials, the material specified here must be nil and the
// individual materials must be added using "AddMaterial" or AddGroupMaterial".
func NewMesh(igeom geometry.IGeometry, imat material.IMaterial) *Mesh {
	m := new(Mesh)
	m.Init(igeom, imat)
	m.ShaderDefines.TOTAL_BONES = 0
	return m
}

// Init initializes the Mesh and its uniforms.
func (m *Mesh) Init(igeom geometry.IGeometry, imat material.IMaterial) {
	m.InitGraphic(m, igeom, gls.TRIANGLES)
	m.uniMatrices.Init("uMatrices")
	m.uniBones.Init("uBones")
	if imat != nil {
		m.AddMaterial(imat, 0, 0)
	}
}

// SetMaterial clears all materials and adds the specified material for all vertices.
func (m *Mesh) SetMaterial(imat material.IMaterial) {
	m.Graphic.ClearMaterials()
	m.Graphic.AddMaterial(m, imat, 0, 0)
}

// AddMaterial adds a material for the specified subset of vertices.
func (m *Mesh) AddMaterial(imat material.IMaterial, start, count int) {
	m.Graphic.AddMaterial(m, imat, start, count)
}

// AddGroupMaterial adds a material for the specified geometry group.
func (m *Mesh) AddGroupMaterial(imat material.IMaterial, gindex int) {
	m.Graphic.AddGroupMaterial(m, imat, gindex)
}

// RenderSetup is called by the engine before drawing the mesh geometry.
// It is responsible for updating the current shader uniforms with the model matrices.
func (m *Mesh) RenderSetup(gs *gls.GLS, _ *core.RenderInfo) {
	// Calculates normal matrix and transfer uniform
	var nm math32.Matrix3
	_ = nm.GetNormalMatrix(&m.mdata.mvm)
	m.mdata.nm.SetFromMatrix3(&nm)
	gs.UniformMatrix4fv(m.uniMatrices.Location(gs), 3, false, &m.mdata.mvm[0])

	if m.skeleton != nil {
		// Get inverse matrix world
		var invMat math32.Matrix4
		node := m.GetNode()
		nMW := node.MatrixWorld()
		err := invMat.GetInverse(&nMW)
		if err != nil {
			log.Print("Skeleton.BoneMatrices: inverting matrix failed")
		}

		// Transfer bone matrices
		boneMatrices := m.skeleton.BoneMatrices(&invMat)
		gs.UniformMatrix4fv(m.uniBones.Location(gs), int32(len(boneMatrices)), false, &boneMatrices[0][0])
	}
}

// SetSkeleton sets the skeleton used by the mesh.
func (m *Mesh) SetSkeleton(skeleton *Skeleton) {
	m.skeleton = skeleton
	if skeleton != nil {
		m.ShaderDefines.TOTAL_BONES = len(m.skeleton.Bones())
	} else {
		m.ShaderDefines.TOTAL_BONES = 0
	}
}

// Skeleton returns the skeleton used by the mesh.
func (m *Mesh) Skeleton() *Skeleton {
	return m.skeleton
}
