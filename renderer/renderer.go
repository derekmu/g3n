// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package renderer implements the scene renderer.
package renderer

import (
	"cmp"
	"slices"
	"sort"

	"github.com/derekmu/g3n/camera"
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/graphic"
	"github.com/derekmu/g3n/gui"
	"github.com/derekmu/g3n/light"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
)

// Renderer renders a scene containing 3D objects and/or 2D GUI elements.
type Renderer struct {
	Shaman                      // Embedded shader manager
	gs          *gls.GLS        // Reference to OpenGL state
	rinfo       core.RenderInfo // Preallocated Render info
	specs       ShaderSpecs     // Preallocated Shader specs
	sortObjects bool            // Flag indicating whether objects should be sorted before rendering
	stats       Stats           // Renderer statistics

	// Populated each frame
	ambLights    []*light.Ambient           // Ambient lights in the scene
	dirLights    []*light.Directional       // Directional lights in the scene
	pointLights  []*light.Point             // Point lights in the scene
	spotLights   []*light.Spot              // Spot lights in the scene
	others       []core.INode               // Other nodes (audio, players, etc)
	graphics     []*graphic.Graphic         // Graphics to be rendered
	grmatsOpaque []*graphic.GraphicMaterial // Opaque graphic materials to be rendered
	grmatsTransp []*graphic.GraphicMaterial // Transparent graphic materials to be rendered
	zLayers      map[int][]gui.IPanel       // All IPanels to be rendered organized by Z-layer
	zLayerKeys   []int                      // Z-layers being used (initially in no particular order, sorted later)
}

// Stats describes how many objects of each type are being rendered.
// It is cleared at the start of each render.
type Stats struct {
	GraphicMats int // Number of graphic materials rendered
	Lights      int // Number of lights rendered
	Panels      int // Number of GUI panels rendered
	Others      int // Number of other objects rendered
}

// NewRenderer creates and returns a pointer to a new Renderer.
func NewRenderer(gs *gls.GLS) *Renderer {
	r := new(Renderer)
	r.gs = gs
	r.Shaman.Init(gs)
	r.sortObjects = true

	r.ambLights = make([]*light.Ambient, 0)
	r.dirLights = make([]*light.Directional, 0)
	r.pointLights = make([]*light.Point, 0)
	r.spotLights = make([]*light.Spot, 0)
	r.others = make([]core.INode, 0)
	r.graphics = make([]*graphic.Graphic, 0)
	r.grmatsOpaque = make([]*graphic.GraphicMaterial, 0)
	r.grmatsTransp = make([]*graphic.GraphicMaterial, 0)
	r.zLayers = make(map[int][]gui.IPanel)
	r.zLayers[0] = nil
	r.zLayerKeys = append(r.zLayerKeys, 0)

	return r
}

// Stats returns a copy of the statistics for the last frame.
// Should be called after the frame was rendered.
func (r *Renderer) Stats() Stats {
	return r.stats
}

// SetObjectSorting sets whether objects will be sorted before rendering.
func (r *Renderer) SetObjectSorting(sort bool) {
	r.sortObjects = sort
}

// ObjectSorting returns whether objects will be sorted before rendering.
func (r *Renderer) ObjectSorting() bool {
	return r.sortObjects
}

// Render renders the specified scene using the specified camera. Returns an error.
func (r *Renderer) Render(scene core.INode, cam camera.ICamera) error {
	// Updates world matrices of all scene nodes
	scene.UpdateMatrixWorld()

	// Build RenderInfo
	cam.ViewMatrix(&r.rinfo.ViewMatrix)
	cam.ProjMatrix(&r.rinfo.ProjMatrix)

	// Clear stats and scene arrays
	r.stats = Stats{}
	r.ambLights = r.ambLights[0:0]
	r.dirLights = r.dirLights[0:0]
	r.pointLights = r.pointLights[0:0]
	r.spotLights = r.spotLights[0:0]
	r.others = r.others[0:0]
	r.graphics = r.graphics[0:0]
	r.grmatsOpaque = r.grmatsOpaque[0:0]
	r.grmatsTransp = r.grmatsTransp[0:0]
	clear(r.zLayers)
	r.zLayers[0] = nil
	r.zLayerKeys = r.zLayerKeys[0:1]
	r.zLayerKeys[0] = 0

	// Advance the frame, for shaders to do once-per-frame processing
	r.Shaman.NewFrame()

	// Prepare for frustum culling
	var proj math32.Matrix4
	proj.MultiplyMatrices(&r.rinfo.ProjMatrix, &r.rinfo.ViewMatrix)
	frustum := math32.NewFrustumFromMatrix(&proj)

	// Classify scene and all scene nodes, culling IGraphics which are fully outside the camera frustum
	r.classifyAndCull(scene, frustum, 0)

	// Set light counts in shader specs
	r.specs.AmbientLightsMax = len(r.ambLights)
	r.specs.DirLightsMax = len(r.dirLights)
	r.specs.PointLightsMax = len(r.pointLights)
	r.specs.SpotLightsMax = len(r.spotLights)

	// Pre-calculate MV and MVP matrices and compile initial lists of opaque and transparent graphic materials
	for _, gr := range r.graphics {
		// Calculate MV and MVP matrices for all non-GUI graphics to be rendered
		gr.CalculateMatrices(&r.rinfo)
		// Append all graphic materials of this graphic to lists of graphic materials to be rendered
		materials := gr.Materials()
		for i := range materials {
			r.stats.GraphicMats++
			if materials[i].IMaterial().GetMaterial().Transparent() {
				r.grmatsTransp = append(r.grmatsTransp, &materials[i])
			} else {
				r.grmatsOpaque = append(r.grmatsOpaque, &materials[i])
			}
		}
	}

	// Z-sort graphic materials back to front
	if r.sortObjects {
		zSort(r.grmatsOpaque)
		zSort(r.grmatsTransp)
	}

	// Sort zLayers back to front
	sort.Ints(r.zLayerKeys)

	// Iterate over all panels from back to front, setting Z and adding graphic materials to grmatsTransp/grmatsOpaque
	const deltaZ = 0.00001
	panZ := -1 + float32(r.stats.Panels)*deltaZ
	for _, k := range r.zLayerKeys {
		for _, ipan := range r.zLayers[k] {
			// Set panel Z
			ipan.SetPositionZ(panZ)
			panZ -= deltaZ
			if ipan.Renderable() {
				// Append the panel's graphic material to lists of graphic materials to be rendered
				mat := ipan.GetGraphic().Materials()[0]
				if mat.IMaterial().GetMaterial().Transparent() {
					r.grmatsTransp = append(r.grmatsTransp, &mat)
				} else {
					r.grmatsOpaque = append(r.grmatsOpaque, &mat)
				}
			}
		}
	}

	// Render opaque objects front to back
	for i := len(r.grmatsOpaque) - 1; i >= 0; i-- {
		err := r.renderGraphicMaterial(r.grmatsOpaque[i])
		if err != nil {
			return err
		}
	}

	// Render transparent objects back to front
	for _, grmat := range r.grmatsTransp {
		err := r.renderGraphicMaterial(grmat)
		if err != nil {
			return err
		}
	}

	// Render other nodes (audio players, etc.)
	for _, inode := range r.others {
		inode.Render(r.gs)
	}

	// Enable depth mask so that clearing the depth buffer works
	r.gs.DepthMask(true)

	return nil
}

// classifyAndCull classifies the provided INode and all of its descendants.
// It ignores (culls) renderable IGraphics which are fully outside the specified frustum.
func (r *Renderer) classifyAndCull(inode core.INode, frustum *math32.Frustum, zLayer int) {
	// Ignore invisible nodes and their descendants
	if !inode.Visible() {
		return
	}
	switch node := inode.(type) {
	case gui.IPanel:
		// non-renderable panels still affect z layers and need to have their z position updated for mouse events
		zLayer += node.ZLayerDelta()
		layer, ok := r.zLayers[zLayer]
		if !ok {
			r.zLayerKeys = append(r.zLayerKeys, zLayer)
		}
		r.zLayers[zLayer] = append(layer, node)
		r.stats.Panels++
	case graphic.IGraphic:
		if node.Renderable() {
			gr := node.GetGraphic()
			if node.Cullable() {
				mw := gr.MatrixWorld()
				bb := node.GetGeometry().BoundingBox()
				bb.ApplyMatrix4(&mw)
				if frustum.IntersectsBox(&bb) {
					r.graphics = append(r.graphics, gr)
				}
			} else {
				r.graphics = append(r.graphics, gr)
			}
		}
	case *light.Ambient:
		if node.Renderable() {
			r.ambLights = append(r.ambLights, node)
		}
	case *light.Directional:
		if node.Renderable() {
			r.dirLights = append(r.dirLights, node)
		}
	case *light.Point:
		if node.Renderable() {
			r.pointLights = append(r.pointLights, node)
		}
	case *light.Spot:
		if node.Renderable() {
			r.spotLights = append(r.spotLights, node)
		}
	default:
		if node.Renderable() {
			r.others = append(r.others, inode)
			r.stats.Others++
		}
	}
	// Process children
	for _, ichild := range inode.Children() {
		r.classifyAndCull(ichild, frustum, zLayer)
	}
}

// zSort sorts a list of graphic materials based on the user-specified render order
// then based on their Z position relative to the camera, back to front.
func zSort(grmats []*graphic.GraphicMaterial) {
	slices.SortFunc(grmats, func(a, b *graphic.GraphicMaterial) int {
		gr1 := a.IGraphic().GetGraphic()
		gr2 := b.IGraphic().GetGraphic()
		// Check for user-supplied render order
		rO1 := gr1.RenderOrder()
		rO2 := gr2.RenderOrder()
		if rO1 != rO2 {
			return cmp.Compare(rO1, rO2)
		}
		mvm1 := gr1.ModelViewMatrix()
		mvm2 := gr2.ModelViewMatrix()
		g1pos := gr1.Position()
		g2pos := gr2.Position()
		g1pos.ApplyMatrix4(mvm1)
		g2pos.ApplyMatrix4(mvm2)
		return cmp.Compare(g1pos.Z, g2pos.Z)
	})
}

// renderGraphicMaterial renders the specified graphic material.
func (r *Renderer) renderGraphicMaterial(grmat *graphic.GraphicMaterial) error {
	mat := grmat.IMaterial().GetMaterial()
	geom := grmat.IGraphic().GetGeometry()
	gr := grmat.IGraphic().GetGraphic()

	// Set the shader specs for this material and set shader program
	r.specs.Name = mat.Shader()
	r.specs.UseLights = mat.UseLights()
	r.specs.MatTexturesMax = mat.TextureCount()
	r.specs.MaterialDefines = mat.ShaderDefines
	r.specs.GeometryDefines = geom.ShaderDefines
	r.specs.GraphicDefines = gr.ShaderDefines

	// Set active program and apply shader specs
	update, err := r.Shaman.SetProgram(r.specs)
	if err != nil {
		return err
	}

	// only update lights the first time each frame for each program
	if update {
		// Set up lights (transfer lights' uniforms)
		if r.specs.UseLights != material.UseLightNone {
			if r.specs.UseLights&material.UseLightAmbient != 0 {
				for idx, l := range r.ambLights {
					l.RenderSetup(r.gs, &r.rinfo, idx)
					r.stats.Lights++
				}
			}
			if r.specs.UseLights&material.UseLightDirectional != 0 {
				for idx, l := range r.dirLights {
					l.RenderSetup(r.gs, &r.rinfo, idx)
					r.stats.Lights++
				}
			}
			if r.specs.UseLights&material.UseLightPoint != 0 {
				for idx, l := range r.pointLights {
					l.RenderSetup(r.gs, &r.rinfo, idx)
					r.stats.Lights++
				}
			}
			if r.specs.UseLights&material.UseLightSpot != 0 {
				for idx, l := range r.spotLights {
					l.RenderSetup(r.gs, &r.rinfo, idx)
					r.stats.Lights++
				}
			}
		}
	}

	// Render this graphic material
	grmat.Render(r.gs, &r.rinfo)

	return nil
}
