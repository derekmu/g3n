// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"strconv"

	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/renderer/shaders"
)

// ShaderSpecs describes the specification of a compiled shader program.
type ShaderSpecs struct {
	Name             string              // Shader name
	UseLights        material.UseLights  // Bitmask indicating which lights to consider
	AmbientLightsMax int                 // Number of ambient lights
	DirLightsMax     int                 // Number of directional lights
	PointLightsMax   int                 // Number of point lights
	SpotLightsMax    int                 // Number of spot lights
	MatTexturesMax   int                 // Number of material textures
	MaterialDefines  gls.MaterialDefines // Additional shader defines
	GeometryDefines  gls.GeometryDefines // Additional shader defines
	GraphicDefines   gls.GraphicDefines  // Additional shader defines
}

// Shaman is the shader manager.
type Shaman struct {
	gs           *gls.GLS
	shaderSource map[string]string              // maps shader name to its template
	programInfo  map[string]shaders.ProgramInfo // maps name of the program to ProgramInfo
	programs     map[ShaderSpecs]*gls.Program   // compiled programs
	specs        ShaderSpecs                    // current shader specs
	program      *gls.Program                   // current program
	frameNumber  int
}

// Init initializes the shader manager
func (sm *Shaman) Init(gs *gls.GLS) {
	sm.gs = gs
	sm.shaderSource = make(map[string]string)
	sm.programInfo = make(map[string]shaders.ProgramInfo)
	sm.programs = make(map[ShaderSpecs]*gls.Program)
}

// AddShaders adds registered shaders and programs to the Shaman.
func (sm *Shaman) AddShaders() {
	for name, source := range shaders.Shaders() {
		sm.shaderSource[name] = source
	}
	for name, programInfo := range shaders.Programs() {
		sm.programInfo[name] = programInfo
	}
}

// SetProgram sets the shader program to satisfy the specs.
//
// Returns whether this is the first time this program was activated this frame and an error if one occurred.
func (sm *Shaman) SetProgram(specs ShaderSpecs) (bool, error) {
	// overwrite light specs based on if the material uses lights
	if (specs.UseLights & material.UseLightAmbient) == 0 {
		specs.AmbientLightsMax = 0
	}
	if (specs.UseLights & material.UseLightDirectional) == 0 {
		specs.DirLightsMax = 0
	}
	if (specs.UseLights & material.UseLightPoint) == 0 {
		specs.PointLightsMax = 0
	}
	if (specs.UseLights & material.UseLightSpot) == 0 {
		specs.SpotLightsMax = 0
	}
	// switch programs if the specs are different from the active specs
	if sm.specs != specs {
		// use a prebuilt program if possible
		var program *gls.Program
		var ok bool
		if program, ok = sm.programs[specs]; !ok {
			// build a new program
			var err error
			program, err = sm.genProgram(&specs)
			if err != nil {
				return false, err
			}
			sm.programs[specs] = program
		}
		// activate the program
		sm.specs = specs
		sm.program = program
		sm.gs.UseProgram(program)
	}
	return sm.program.SetFrameNumber(sm.frameNumber), nil
}

// genProgram generates shader program from the specified specs
func (sm *Shaman) genProgram(specs *ShaderSpecs) (*gls.Program, error) {
	progInfo, ok := sm.programInfo[specs.Name]
	if !ok {
		return nil, fmt.Errorf("program %s not found", specs.Name)
	}
	// prepare the defines map
	defines := map[string]string{
		"AMB_LIGHTS":   strconv.Itoa(specs.AmbientLightsMax),
		"DIR_LIGHTS":   strconv.Itoa(specs.DirLightsMax),
		"POINT_LIGHTS": strconv.Itoa(specs.PointLightsMax),
		"SPOT_LIGHTS":  strconv.Itoa(specs.SpotLightsMax),
		"MAT_TEXTURES": strconv.Itoa(specs.MatTexturesMax),
	}
	specs.MaterialDefines.AddToMap(defines)
	specs.GeometryDefines.AddToMap(defines)
	specs.GraphicDefines.AddToMap(defines)
	// prepare the vertex shader
	vertexSource, ok := sm.shaderSource[progInfo.Vertex]
	if !ok {
		return nil, fmt.Errorf("vertex shader %s not found", progInfo.Vertex)
	}
	vertexSource = sm.preprocess(vertexSource, defines)
	// prepare the fragment shader
	fragSource, ok := sm.shaderSource[progInfo.Fragment]
	if !ok {
		return nil, fmt.Errorf("fragment shader %s not found", progInfo.Fragment)
	}
	fragSource = sm.preprocess(fragSource, defines)
	// prepare the geometry shader
	var geomSource = ""
	if progInfo.Geometry != "" {
		geomSource, ok = sm.shaderSource[progInfo.Geometry]
		if !ok {
			return nil, fmt.Errorf("geometry shader %s not found", progInfo.Geometry)
		}
		geomSource = sm.preprocess(geomSource, defines)
	}
	// build the program
	prog := sm.gs.NewProgram(specs.Name)
	prog.AddShader(gls.VERTEX_SHADER, vertexSource)
	prog.AddShader(gls.FRAGMENT_SHADER, fragSource)
	if progInfo.Geometry != "" {
		prog.AddShader(gls.GEOMETRY_SHADER, geomSource)
	}
	err := prog.Build()
	if err != nil {
		return nil, err
	}
	return prog, nil
}

func (sm *Shaman) preprocess(source string, defines map[string]string) string {
	var prefix = fmt.Sprintf("#version %s\n", GLSL_VERSION)
	for name, value := range defines {
		prefix = prefix + fmt.Sprintf("#define %s %s\n", name, value)
	}
	return prefix + source
}

func (sm *Shaman) NewFrame() {
	sm.frameNumber++
}
