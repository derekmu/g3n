// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graphic implements scene objects which have a graphic representation.
package graphic

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/geometry"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
)

// IGraphic is the interface for all Graphic objects.
type IGraphic interface {
	core.INode
	GetGraphic() *Graphic
	GetGeometry() *geometry.Geometry
	IGeometry() geometry.IGeometry
	SetRenderable(bool)
	Renderable() bool
	SetCullable(bool)
	Cullable() bool
	RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo)
}

// Graphic is a Node which has a visible representation in the scene.
// It has an associated geometry and one or more materials.
// It is the base type used by other graphics such as lines, line strips, points, and meshes.
type Graphic struct {
	core.Node                        // Embedded Node
	igeom         geometry.IGeometry // Associated IGeometry
	materials     []GraphicMaterial  // Materials
	mode          uint32             // OpenGL primitive
	cullable      bool               // Cullable flag
	renderOrder   int                // Render order
	ShaderDefines gls.GraphicDefines // Graphic-specific shader defines
	mdata         struct {
		mvm  math32.Matrix4 // ModelViewMatrix
		mvpm math32.Matrix4 // ModelViewProjectionMatrix
		nm   math32.Matrix4 // NormalMatrix
	}
}

// InitGraphic initializes a Graphic type embedded in another type
// with the specified geometry and OpenGL mode.
func (gr *Graphic) InitGraphic(igr IGraphic, igeom geometry.IGeometry, mode uint32) *Graphic {
	gr.InitNode(igr)
	igeom.Incref()
	gr.igeom = igeom
	gr.mode = mode
	gr.materials = nil
	gr.cullable = true
	return gr
}

// GetGraphic satisfies the IGraphic interface and
// returns pointer to the base Graphic.
func (gr *Graphic) GetGraphic() *Graphic {
	return gr
}

// GetGeometry satisfies the IGraphic interface and returns
// a pointer to the geometry associated with this graphic.
func (gr *Graphic) GetGeometry() *geometry.Geometry {
	return gr.igeom.GetGeometry()
}

// IGeometry satisfies the IGraphic interface and returns
// a pointer to the IGeometry associated with this graphic.
func (gr *Graphic) IGeometry() geometry.IGeometry {
	return gr.igeom
}

// Dispose overrides the embedded Node Dispose method.
func (gr *Graphic) Dispose() {
	if gr.igeom.Decref() {
		gr.igeom.Dispose()
	}
	gr.igeom = nil
	gr.ClearMaterials()
}

// ClearMaterials removes all the materials from this Graphic.
func (gr *Graphic) ClearMaterials() {
	for _, mat := range gr.materials {
		if mat.imat.Decref() {
			mat.imat.Dispose()
		}
	}
	clear(gr.materials)
	gr.materials = nil
}

// SetCullable satisfies the IGraphic interface and
// sets the cullable state of this Graphic (default = true).
func (gr *Graphic) SetCullable(state bool) {
	gr.cullable = state
}

// Cullable satisfies the IGraphic interface and
// returns the cullable state of this graphic.
func (gr *Graphic) Cullable() bool {
	return gr.cullable
}

// SetRenderOrder sets the render order of the object.
// All objects have renderOrder of 0 by default.
// To render before renderOrder 0 set a lower renderOrder e.g. -1.
// To render after renderOrder 0 set a higher renderOrder e.g. 1
func (gr *Graphic) SetRenderOrder(order int) {
	gr.renderOrder = order
}

// RenderOrder returns the render order of the object.
func (gr *Graphic) RenderOrder() int {
	return gr.renderOrder
}

// AddMaterial adds a material for the specified subset of vertices.
// If the material applies to all vertices, start and count must be 0.
func (gr *Graphic) AddMaterial(igr IGraphic, imat material.IMaterial, start, count int) {
	gmat := GraphicMaterial{
		imat:     imat,
		start:    start,
		count:    count,
		igraphic: igr,
	}
	imat.Incref()
	gr.materials = append(gr.materials, gmat)
}

// AddGroupMaterial adds a material for the specified geometry group.
func (gr *Graphic) AddGroupMaterial(igr IGraphic, imat material.IMaterial, gindex int) {
	geom := gr.igeom.GetGeometry()
	if gindex < 0 || gindex >= geom.GroupCount() {
		panic("Invalid group index")
	}
	group := geom.GroupAt(gindex)
	gr.AddMaterial(igr, imat, group.Start, group.Count)
}

// Materials returns slice with this graphic materials.
func (gr *Graphic) Materials() []GraphicMaterial {
	return gr.materials
}

// GetMaterial returns the material associated with the specified vertex position.
func (gr *Graphic) GetMaterial(vpos int) material.IMaterial {
	for _, gmat := range gr.materials {
		// One material
		if gmat.count == 0 {
			return gmat.imat
		}
		if gmat.start <= vpos && gmat.start+gmat.count >= vpos {
			return gmat.imat
		}
	}
	return nil
}

// SetIGraphic sets the IGraphic on all this Graphic's GraphicMaterials.
func (gr *Graphic) SetIGraphic(igr IGraphic) {
	for i := range gr.materials {
		gr.materials[i].igraphic = igr
	}
}

// BoundingBox recursively calculates and returns the bounding box
// containing this node and all its children.
func (gr *Graphic) BoundingBox() math32.Box3 {
	geom := gr.igeom.GetGeometry()
	bbox := geom.BoundingBox()
	m := gr.MatrixWorld()
	bbox.ApplyMatrix4(&m)
	for _, inode := range gr.Children() {
		childGraphic, ok := inode.(*Graphic)
		if ok {
			childBbox := childGraphic.BoundingBox()
			bbox.Union(&childBbox)
		}
	}
	return bbox
}

// CalculateMatrices calculates the model view and model view projection matrices.
func (gr *Graphic) CalculateMatrices(rinfo *core.RenderInfo) {
	mm := gr.MatrixWorld()
	gr.mdata.mvm.MultiplyMatrices(&rinfo.ViewMatrix, &mm)
	gr.mdata.mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &gr.mdata.mvm)
}

// ModelViewMatrix returns the last cached model view matrix for this graphic.
func (gr *Graphic) ModelViewMatrix() *math32.Matrix4 {
	return &gr.mdata.mvm
}

// ModelViewProjectionMatrix returns the last cached model view projection matrix for this graphic.
func (gr *Graphic) ModelViewProjectionMatrix() *math32.Matrix4 {
	return &gr.mdata.mvpm
}

// GraphicMaterial specifies the material to be used for a subset of vertices from a Graphic's geometry.
type GraphicMaterial struct {
	imat     material.IMaterial // Associated material
	start    int                // Index of first element in the geometry
	count    int                // Number of elements
	igraphic IGraphic           // Graphic which contains this GraphicMaterial
}

// IMaterial returns the material associated with the GraphicMaterial.
func (grmat *GraphicMaterial) IMaterial() material.IMaterial {
	return grmat.imat
}

// IGraphic returns the graphic associated with the GraphicMaterial.
func (grmat *GraphicMaterial) IGraphic() IGraphic {
	return grmat.igraphic
}

// Render is called by the renderer to render this graphic material.
func (grmat *GraphicMaterial) Render(gs *gls.GLS, rinfo *core.RenderInfo) {
	grmat.imat.RenderSetup(gs)
	gr := grmat.igraphic.GetGraphic()
	gr.igeom.RenderSetup(gs)
	grmat.igraphic.RenderSetup(gs, rinfo)

	count := grmat.count
	geom := gr.igeom.GetGeometry()
	indices := geom.Indices()
	if indices.Len() > 0 {
		if count == 0 {
			count = indices.Len()
		}
		gs.DrawElements(gr.mode, int32(count), gls.UNSIGNED_INT, 4*uint32(grmat.start))
	} else {
		if count == 0 {
			count = geom.Items()
		}
		gs.DrawArrays(gr.mode, int32(grmat.start), int32(count))
	}
}
