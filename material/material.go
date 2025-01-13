// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package material contains virtual materials which
// specify the appearance of graphic objects.
package material

import (
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/texture"
	"slices"
)

// Side represents the material's visible side(s).
type Side int

// The face side(s) to be rendered. The non-rendered side will be culled to improve performance.
const (
	SideFront = Side(iota)
	SideBack
	SideDouble
)

// Blending specifies the blending mode.
type Blending int

// The various blending types.
const (
	BlendNone = Blending(iota)
	BlendNormal
	BlendAdditive
	BlendSubtractive
	BlendMultiply
	BlendCustom
)

// UseLights is a bitmask that specifies which types of lights affect the material.
type UseLights int

// The possible UseLights values.
const (
	UseLightNone        UseLights = 0x00
	UseLightAmbient     UseLights = 0x01
	UseLightDirectional UseLights = 0x02
	UseLightPoint       UseLights = 0x04
	UseLightSpot        UseLights = 0x08
	UseLightAll         UseLights = 0xFF
)

// IMaterial is the interface for all materials.
type IMaterial interface {
	GetMaterial() *Material
	RenderSetup(gs *gls.GLS)
	Dispose()
}

// Material is the base material.
type Material struct {
	refcount int // Current number of references

	// Shader specification
	shader        string              // Shader name
	ShaderDefines gls.MaterialDefines // shader defines

	side          Side                 // Face side(s) visibility
	blending      Blending             // Blending mode
	useLights     UseLights            // Which light types to consider
	transparent   bool                 // Whether at all transparent
	wireframe     bool                 // Whether to render only the wireframe
	lineWidth     float32              // Line width for lines and wireframe
	textures      []*texture.Texture2D // List of textures
	samplerCounts map[string]int

	polyOffsetFactor float32 // polygon offset factor
	polyOffsetUnits  float32 // polygon offset units

	depthMask bool   // Enable writing into the depth buffer
	depthTest bool   // Enable depth buffer test
	depthFunc uint32 // Active depth test function

	// Equations used for custom blending (when blending=BlendCustom)
	blendRGB      uint32 // separate blending equation for RGB
	blendAlpha    uint32 // separate blending equation for Alpha
	blendSrcRGB   uint32 // separate blending func source RGB
	blendDstRGB   uint32 // separate blending func dest RGB
	blendSrcAlpha uint32 // separate blending func source Alpha
	blendDstAlpha uint32 // separate blending func dest Alpha
}

// NewMaterial creates a new Material.
func NewMaterial() *Material {
	m := new(Material)
	return m.InitMaterial()
}

// InitMaterial initializes the material.
func (m *Material) InitMaterial() *Material {
	m.refcount = 1
	m.useLights = UseLightAll
	m.side = SideFront
	m.transparent = false
	m.wireframe = false
	m.depthMask = true
	m.depthFunc = gls.LEQUAL
	m.depthTest = true
	m.blending = BlendNormal
	m.lineWidth = 1.0
	m.polyOffsetFactor = 0
	m.polyOffsetUnits = 0
	m.textures = make([]*texture.Texture2D, 0)
	m.samplerCounts = make(map[string]int)
	return m
}

// GetMaterial satisfies the IMaterial interface.
func (m *Material) GetMaterial() *Material {
	return m
}

// Incref increments the reference count for this material
// and returns a pointer to the material.
// It should be used when this material is shared by another
// Graphic object.
func (m *Material) Incref() *Material {
	m.refcount++
	return m
}

// Dispose decrements this material reference count and
// if necessary releases OpenGL resources, C memory
// and textures associated with this material.
func (m *Material) Dispose() {
	// Only dispose if last
	if m.refcount > 1 {
		m.refcount--
		return
	}
	// Delete textures
	for i := 0; i < len(m.textures); i++ {
		m.textures[i].Dispose()
	}
	m.InitMaterial()
}

// SetShader sets the name of the shader program for this material
func (m *Material) SetShader(sname string) {
	m.shader = sname
}

// Shader returns the current name of the shader program for this material
func (m *Material) Shader() string {
	return m.shader
}

// SetUseLights sets the material use lights bit mask specifying which
// light types will be used when rendering the material
// By default the material will use all lights
func (m *Material) SetUseLights(lights UseLights) {
	m.useLights = lights
}

// UseLights returns the current use lights bitmask
func (m *Material) UseLights() UseLights {
	return m.useLights
}

// SetSide sets the visible side(s) (SideFront | SideBack | SideDouble)
func (m *Material) SetSide(side Side) {
	m.side = side
}

// Side returns the current side visibility for this material
func (m *Material) Side() Side {
	return m.side
}

// SetTransparent sets whether this material is transparent.
func (m *Material) SetTransparent(state bool) {
	m.transparent = state
}

// Transparent returns whether this material has been set as transparent.
func (m *Material) Transparent() bool {
	return m.transparent
}

// SetWireframe sets whether only the wireframe is rendered.
func (m *Material) SetWireframe(state bool) {
	m.wireframe = state
}

// Wireframe returns whether only the wireframe is rendered.
func (m *Material) Wireframe() bool {
	return m.wireframe
}

func (m *Material) SetDepthMask(state bool) {
	m.depthMask = state
}

func (m *Material) SetDepthTest(state bool) {
	m.depthTest = state
}

func (m *Material) SetDepthFunc(state uint32) {
	m.depthFunc = state
}

func (m *Material) SetBlending(blending Blending) {
	m.blending = blending
}

func (m *Material) SetLineWidth(width float32) {
	m.lineWidth = width
}

func (m *Material) SetPolygonOffset(factor, units float32) {
	m.polyOffsetFactor = factor
	m.polyOffsetUnits = units
}

// RenderSetup is called by the renderer before drawing objects with this material.
func (m *Material) RenderSetup(gs *gls.GLS) {
	// Sets triangle side view mode
	switch m.side {
	case SideFront:
		gs.Enable(gls.CULL_FACE)
		gs.FrontFace(gls.CCW)
	case SideBack:
		gs.Enable(gls.CULL_FACE)
		gs.FrontFace(gls.CW)
	case SideDouble:
		gs.Disable(gls.CULL_FACE)
		gs.FrontFace(gls.CCW)
	}

	if m.depthTest {
		gs.Enable(gls.DEPTH_TEST)
	} else {
		gs.Disable(gls.DEPTH_TEST)
	}
	gs.DepthMask(m.depthMask)
	gs.DepthFunc(m.depthFunc)

	if m.wireframe {
		gs.PolygonMode(gls.FRONT_AND_BACK, gls.LINE)
	} else {
		gs.PolygonMode(gls.FRONT_AND_BACK, gls.FILL)
	}

	// Set polygon offset if requested
	gs.PolygonOffset(m.polyOffsetFactor, m.polyOffsetUnits)

	// Sets line width
	gs.LineWidth(m.lineWidth)

	// Sets blending
	switch m.blending {
	case BlendNone:
		gs.Disable(gls.BLEND)
	case BlendNormal:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.SRC_ALPHA, gls.ONE_MINUS_SRC_ALPHA)
	case BlendAdditive:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.SRC_ALPHA, gls.ONE)
	case BlendSubtractive:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.ZERO, gls.ONE_MINUS_SRC_COLOR)
		break
	case BlendMultiply:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.ZERO, gls.SRC_COLOR)
		break
	case BlendCustom:
		gs.BlendEquationSeparate(m.blendRGB, m.blendAlpha)
		gs.BlendFuncSeparate(m.blendSrcRGB, m.blendDstRGB, m.blendSrcAlpha, m.blendDstAlpha)
		break
	default:
		panic("Invalid blending")
	}

	// Render textures
	// Keep track of counts of unique sampler names to correctly index sampler arrays
	clear(m.samplerCounts)
	for slotIdx, tex := range m.textures {
		samplerName, _ := tex.GetUniformNames()
		uniIdx, _ := m.samplerCounts[samplerName]
		tex.RenderSetup(gs, slotIdx, uniIdx)
		m.samplerCounts[samplerName] = uniIdx + 1
	}
}

// AddTexture adds the specified texture to the material
func (m *Material) AddTexture(tex *texture.Texture2D) {
	m.textures = append(m.textures, tex)
}

// RemoveTexture removes the specified texture from the material
func (m *Material) RemoveTexture(tex *texture.Texture2D) {
	for pos, curr := range m.textures {
		if curr == tex {
			copy(m.textures[pos:], m.textures[pos+1:])
			m.textures[len(m.textures)-1] = nil
			m.textures = m.textures[:len(m.textures)-1]
			break
		}
	}
}

// HasTexture checks if the material contains the specified texture
func (m *Material) HasTexture(tex *texture.Texture2D) bool {
	return slices.Index(m.textures, tex) >= 0
}

// TextureCount returns the current number of textures
func (m *Material) TextureCount() int {
	return len(m.textures)
}

// Textures returns a slice with this material's textures
func (m *Material) Textures() []*texture.Texture2D {
	return m.textures
}
