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

// ShaderSpecs describes the specification of a compiled shader program
type ShaderSpecs struct {
	Name             string              // Shader name
	Version          string              // GLSL version
	UseLights        material.UseLights  // Bitmask indicating which lights to consider
	AmbientLightsMax int                 // Current number of ambient lights
	DirLightsMax     int                 // Current Number of directional lights
	PointLightsMax   int                 // Current Number of point lights
	SpotLightsMax    int                 // Current Number of spot lights
	MatTexturesMax   int                 // Current Number of material textures
	MaterialDefines  gls.MaterialDefines // Additional shader defines
	GeometryDefines  gls.GeometryDefines // Additional shader defines
	GraphicDefines   gls.GraphicDefines  // Additional shader defines
}

// Shaman is the shader manager
type Shaman struct {
	gs       *gls.GLS                       // Reference to OpenGL state
	includes map[string]string              // include files sources
	shadersm map[string]string              // maps shader name to its template
	proginfo map[string]shaders.ProgramInfo // maps name of the program to ProgramInfo
	programs map[ShaderSpecs]*gls.Program   // compiled programs with specs
	specs    ShaderSpecs                    // Current shader specs
}

// NewShaman creates and returns a pointer to a new shader manager
func NewShaman(gs *gls.GLS) *Shaman {
	sm := new(Shaman)
	sm.Init(gs)
	return sm
}

// Init initializes the shader manager
func (sm *Shaman) Init(gs *gls.GLS) {
	sm.gs = gs
	sm.includes = make(map[string]string)
	sm.shadersm = make(map[string]string)
	sm.proginfo = make(map[string]shaders.ProgramInfo)
}

// AddDefaultShaders adds to this shader manager all default
// include chunks, shaders and programs statically registered.
func (sm *Shaman) AddDefaultShaders() error {
	for _, name := range shaders.Shaders() {
		sm.AddShader(name, shaders.ShaderSource(name))
	}

	for _, name := range shaders.Programs() {
		sm.proginfo[name] = shaders.GetProgramInfo(name)
	}
	return nil
}

// AddChunk adds a shader chunk with the specified name and source code
func (sm *Shaman) AddChunk(name, source string) {
	sm.includes[name] = source
}

// AddShader adds a shader program with the specified name and source code
func (sm *Shaman) AddShader(name, source string) {
	sm.shadersm[name] = source
}

// AddProgram adds a program with the specified name and associated vertex
// and fragment shaders names (previously registered)
func (sm *Shaman) AddProgram(name, vertexName, fragName string, others ...string) {
	geomName := ""
	if len(others) > 0 {
		geomName = others[0]
	}
	sm.proginfo[name] = shaders.ProgramInfo{
		Vertex:   vertexName,
		Fragment: fragName,
		Geometry: geomName,
	}
}

// SetProgram sets the shader program to satisfy the specified specs.
// Returns an indication if the current shader has changed and a possible error
// when creating a new shader program.
// Receives a copy of the specs because it changes the fields which specify the
// number of lights depending on the UseLights flags.
func (sm *Shaman) SetProgram(s *ShaderSpecs) (bool, error) {
	// Checks material use lights bit mask
	// copy so we don't change light settings provided
	specs := *s
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

	// If current shader specs are the same as the specified specs, nothing to do.
	if sm.specs == specs {
		return false, nil
	}

	// Search for compiled program with the specified specs
	if program, ok := sm.programs[specs]; ok {
		sm.gs.UseProgram(program)
		sm.specs = specs
		return true, nil
	}

	// Generates new program with the specified specs
	program, err := sm.GenProgram(&specs)
	if err != nil {
		return false, err
	}

	// Save specs as current specs, adds new program to the list and activates the program
	sm.specs = specs
	sm.programs[specs] = program
	sm.gs.UseProgram(program)
	return true, nil
}

// GenProgram generates shader program from the specified specs
func (sm *Shaman) GenProgram(specs *ShaderSpecs) (*gls.Program, error) {
	// Get info for the specified shader program
	progInfo, ok := sm.proginfo[specs.Name]
	if !ok {
		return nil, fmt.Errorf("Program:%s not found", specs.Name)
	}

	// Sets the defines map
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

	vertexSource, ok := sm.shadersm[progInfo.Vertex]
	if !ok {
		return nil, fmt.Errorf("Vertex shader:%s not found", progInfo.Vertex)
	}
	vertexSource = sm.preprocess(vertexSource, defines)

	fragSource, ok := sm.shadersm[progInfo.Fragment]
	if !ok {
		return nil, fmt.Errorf("Fragment shader:%s not found", progInfo.Fragment)
	}
	fragSource = sm.preprocess(fragSource, defines)

	// Checks for optional geometry shader compiled template
	var geomSource = ""
	if progInfo.Geometry != "" {
		geomSource, ok = sm.shadersm[progInfo.Geometry]
		if !ok {
			return nil, fmt.Errorf("Geometry shader:%s not found", progInfo.Geometry)
		}
		geomSource = sm.preprocess(geomSource, defines)
	}

	// Creates shader program
	prog := sm.gs.NewProgram()
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
	// If defines map supplied, generate prefix with glsl version directive first,
	// followed by "#define" directives
	var prefix = ""
	if defines != nil { // This is only true for the outer call
		prefix = fmt.Sprintf("#version %s\n", GLSL_VERSION)
		for name, value := range defines {
			prefix = prefix + fmt.Sprintf("#define %s %s\n", name, value)
		}
	}
	return prefix + source
}
