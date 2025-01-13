// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/texture"
)

// Physical is a physically based rendered material which uses the metallic-roughness model.
type Physical struct {
	Material                                // Embedded material
	baseColorTex         *texture.Texture2D // Optional base color texture
	metallicRoughnessTex *texture.Texture2D // Optional metallic-roughness
	normalTex            *texture.Texture2D // Optional normal texture
	occlusionTex         *texture.Texture2D // Optional occlusion texture
	emissiveTex          *texture.Texture2D // Optional emissive texture
	uni                  gls.Uniform        // Uniform location cache
	udata                struct {           // Combined uniform data
		baseColorFactor math32.Color4
		emissiveFactor  math32.Color4
		metallicFactor  float32
		roughnessFactor float32
	}
}

// Number of glsl shader vec4 elements used by uniform data.
const physicalVec4Count = 3

// NewPhysical creates a new Physical material.
func NewPhysical() *Physical {
	m := new(Physical)
	m.InitPhysical()
	m.SetShader("physical")
	return m
}

// InitPhysical initializes the material.
func (m *Physical) InitPhysical() {
	m.InitMaterial()
	// Creates uniform and set default values
	m.uni.Init("Material")
	m.udata.baseColorFactor = math32.Color4{1, 1, 1, 1}
	m.udata.emissiveFactor = math32.Color4{0, 0, 0, 1}
	m.udata.metallicFactor = 1
	m.udata.roughnessFactor = 1
}

// SetBaseColorFactor sets this material's base color.
// Its default value is {1,1,1,1}.
func (m *Physical) SetBaseColorFactor(c math32.Color4) {
	m.udata.baseColorFactor = c
}

// SetMetallicFactor sets this material's metallic factor.
// Its default value is 1.
func (m *Physical) SetMetallicFactor(v float32) {
	m.udata.metallicFactor = v
}

// SetRoughnessFactor sets this material's roughness factor.
// Its default value is 1.
func (m *Physical) SetRoughnessFactor(v float32) {
	m.udata.roughnessFactor = v
}

// SetEmissiveFactor sets this material's emissive.
// Its default is {1, 1, 1}.
func (m *Physical) SetEmissiveFactor(c math32.Color) {
	m.udata.emissiveFactor.R = c.R
	m.udata.emissiveFactor.G = c.G
	m.udata.emissiveFactor.B = c.B
}

// SetBaseColorMap sets this material's optional texture base color.
func (m *Physical) SetBaseColorMap(tex *texture.Texture2D) {
	m.baseColorTex = tex
	if m.baseColorTex != nil {
		m.baseColorTex.SetUniformNames("uBaseColorSampler", "")
		m.ShaderDefines.HAS_BASECOLORMAP = true
		m.AddTexture(m.baseColorTex)
	} else {
		m.ShaderDefines.HAS_BASECOLORMAP = false
		m.RemoveTexture(m.baseColorTex)
	}
}

// SetMetallicRoughnessMap sets this material's optional metallic-roughness texture.
func (m *Physical) SetMetallicRoughnessMap(tex *texture.Texture2D) {
	m.metallicRoughnessTex = tex
	if m.metallicRoughnessTex != nil {
		m.metallicRoughnessTex.SetUniformNames("uMetallicRoughnessSampler", "")
		m.ShaderDefines.HAS_METALROUGHNESSMAP = true
		m.AddTexture(m.metallicRoughnessTex)
	} else {
		m.ShaderDefines.HAS_METALROUGHNESSMAP = false
		m.RemoveTexture(m.metallicRoughnessTex)
	}
}

// SetNormalMap sets this material's optional normal texture.
func (m *Physical) SetNormalMap(tex *texture.Texture2D) {
	m.normalTex = tex
	if m.normalTex != nil {
		m.normalTex.SetUniformNames("uNormalSampler", "")
		m.ShaderDefines.HAS_NORMALMAP = true
		m.AddTexture(m.normalTex)
	} else {
		m.ShaderDefines.HAS_NORMALMAP = false
		m.RemoveTexture(m.normalTex)
	}
}

// SetOcclusionMap sets this material's optional occlusion texture.
func (m *Physical) SetOcclusionMap(tex *texture.Texture2D) {
	m.occlusionTex = tex
	if m.occlusionTex != nil {
		m.occlusionTex.SetUniformNames("uOcclusionSampler", "")
		m.ShaderDefines.HAS_OCCLUSIONMAP = true
		m.AddTexture(m.occlusionTex)
	} else {
		m.ShaderDefines.HAS_OCCLUSIONMAP = false
		m.RemoveTexture(m.occlusionTex)
	}
}

// SetEmissiveMap sets this material's optional emissive texture.
func (m *Physical) SetEmissiveMap(tex *texture.Texture2D) {
	m.emissiveTex = tex
	if m.emissiveTex != nil {
		m.emissiveTex.SetUniformNames("uEmissiveSampler", "")
		m.ShaderDefines.HAS_EMISSIVEMAP = true
		m.AddTexture(m.emissiveTex)
	} else {
		m.ShaderDefines.HAS_EMISSIVEMAP = false
		m.RemoveTexture(m.emissiveTex)
	}
}

// RenderSetup transfers this material's uniforms and textures to the shader.
func (m *Physical) RenderSetup(gl *gls.GLS) {
	m.Material.RenderSetup(gl)
	location := m.uni.Location(gl)
	gl.Uniform4fv(location, physicalVec4Count, &m.udata.baseColorFactor.R)
}
