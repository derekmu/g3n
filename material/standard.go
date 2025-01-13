// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/math32"
)

// Standard is a material that supports the classic lighting model with ambient, diffuse, specular and emissive lights.
// The lighting calculation is implemented in the vertex shader.
type Standard struct {
	Material             // Embedded material
	uni      gls.Uniform // Uniform location cache
	udata    struct {    // Combined uniform data in 6 vec3:
		ambient    math32.Color // Ambient color reflectivity
		diffuse    math32.Color // Diffuse color reflectivity
		specular   math32.Color // Specular color reflectivity
		emissive   math32.Color // Emissive color
		shininess  float32      // Specular shininess factor
		opacity    float32      // Opacity
		psize      float32      // Point size
		protationZ float32      // Point rotation around Z axis
	}
}

// Number of glsl shader vec3 elements used by uniform data.
const standardVec3Count = 6

// NewStandard creates and returns a pointer to a new standard material.
func NewStandard(color math32.Color) *Standard {
	m := new(Standard)
	m.InitStandard(color)
	m.SetShader("standard")
	return m
}

// NewBlinnPhong creates and returns a pointer to a new Standard material using Blinn-Phong model.
// It is very close to Standard (Phong) model so we need only pass a parameter.
func NewBlinnPhong(color math32.Color) *Standard {
	m := new(Standard)
	m.InitStandard(color)
	m.SetShader("standard")
	m.ShaderDefines.BLINN = true
	return m
}

// InitStandard initializes the material setting the specified color.
func (m *Standard) InitStandard(color math32.Color) {
	m.InitMaterial()
	// Creates uniforms and set initial values
	m.uni.Init("Material")
	m.SetColor(color)
	m.SetSpecularColor(math32.Color{0.5, 0.5, 0.5})
	m.SetEmissiveColor(math32.Color{0, 0, 0})
	m.SetShininess(30.0)
	m.SetOpacity(1.0)
}

// AmbientColor returns the material ambient color reflectivity.
func (m *Standard) AmbientColor() math32.Color {
	return m.udata.ambient
}

// SetAmbientColor sets the material ambient color reflectivity.
// The default is the same as the diffuse color.
func (m *Standard) SetAmbientColor(color math32.Color) {
	m.udata.ambient = color
}

// SetColor sets the material diffuse color and also the material ambient color reflectivity.
func (m *Standard) SetColor(color math32.Color) {
	m.udata.diffuse = color
	m.udata.ambient = color
}

// SetEmissiveColor sets the material emissive color.
// The default is {0,0,0}.
func (m *Standard) SetEmissiveColor(color math32.Color) {
	m.udata.emissive = color
}

// EmissiveColor returns the material current emissive color.
func (m *Standard) EmissiveColor() math32.Color {
	return m.udata.emissive
}

// SetSpecularColor sets the material specular color reflectivity.
// The default is {0.5, 0.5, 0.5}.
func (m *Standard) SetSpecularColor(color math32.Color) {
	m.udata.specular = color
}

// SetShininess sets the specular highlight factor.
// The default is 30.
func (m *Standard) SetShininess(shininess float32) {
	m.udata.shininess = shininess
}

// SetOpacity sets the material opacity (alpha).
// The default is 1.0.
func (m *Standard) SetOpacity(opacity float32) {
	m.udata.opacity = opacity
}

// RenderSetup is called before drawing the object which uses this material.
func (m *Standard) RenderSetup(gs *gls.GLS) {
	m.Material.RenderSetup(gs)
	location := m.uni.Location(gs)
	gs.Uniform3fv(location, standardVec3Count, &m.udata.ambient.R)
}
