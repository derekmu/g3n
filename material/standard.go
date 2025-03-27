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
	Material
	uni   gls.Uniform
	udata struct {
		ambient    math32.Color3
		diffuse    math32.Color3
		specular   math32.Color3
		emissive   math32.Color3
		shininess  float32
		opacity    float32
		psize      float32
		protationZ float32
		_          float32
		_          float32
	}
}

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
	m.uni.Init("uMaterial")
	m.SetColor(color)
	m.SetSpecularColor(math32.Color3{R: 0.5, G: 0.5, B: 0.5})
	m.SetEmissiveColor(math32.Color3{})
	m.SetShininess(30.0)
	m.SetOpacity(1.0)
}

// AmbientColor returns the material ambient color reflectivity.
func (m *Standard) AmbientColor() math32.Color {
	return m.udata.ambient
}

// SetAmbientColor sets the material ambient color reflectivity.
func (m *Standard) SetAmbientColor(color math32.Color) {
	m.udata.ambient = color.Color3()
}

// SetColor sets the material diffuse color and also the material ambient color reflectivity.
func (m *Standard) SetColor(color math32.Color) {
	m.udata.diffuse = color.Color3()
	m.udata.ambient = color.Color3()
}

// SetEmissiveColor sets the material emissive color.
func (m *Standard) SetEmissiveColor(color math32.Color) {
	m.udata.emissive = color.Color3()
}

// EmissiveColor returns the material current emissive color.
func (m *Standard) EmissiveColor() math32.Color {
	return m.udata.emissive
}

// SetSpecularColor sets the material specular color reflectivity.
func (m *Standard) SetSpecularColor(color math32.Color) {
	m.udata.specular = color.Color3()
}

// SetShininess sets the specular highlight factor.
func (m *Standard) SetShininess(shininess float32) {
	m.udata.shininess = shininess
}

// SetOpacity sets the material opacity (alpha).
func (m *Standard) SetOpacity(opacity float32) {
	m.udata.opacity = opacity
}

// RenderSetup is called before drawing the object which uses this material.
func (m *Standard) RenderSetup(gs *gls.GLS) {
	m.Material.RenderSetup(gs)
	location := m.uni.Location(gs)
	gs.Uniform3fv(location, 6, &m.udata.ambient.R)
}
