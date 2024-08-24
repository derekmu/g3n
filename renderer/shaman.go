// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"maps"
	"strconv"

	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/renderer/shaders"
)

// ShaderSpecs describes the specification of a compiled shader program
type ShaderSpecs struct {
	Name             string             // Shader name
	Version          string             // GLSL version
	ShaderUnique     bool               // indicates if shader is independent of lights and textures
	UseLights        material.UseLights // Bitmask indicating which lights to consider
	AmbientLightsMax int                // Current number of ambient lights
	DirLightsMax     int                // Current Number of directional lights
	PointLightsMax   int                // Current Number of point lights
	SpotLightsMax    int                // Current Number of spot lights
	MatTexturesMax   int                // Current Number of material textures
	MatDefines       gls.ShaderDefines  // Additional shader defines
	GeomDefines      gls.ShaderDefines  // Additional shader defines
	GrDefines        gls.ShaderDefines  // Additional shader defines
}

// ProgSpecs represents a compiled shader program along with its specs
type ProgSpecs struct {
	program *gls.Program // program object
	specs   ShaderSpecs  // associated specs
}

// Shaman is the shader manager
type Shaman struct {
	gs       *gls.GLS                       // Reference to OpenGL state
	includes map[string]string              // include files sources
	shadersm map[string]string              // maps shader name to its template
	proginfo map[string]shaders.ProgramInfo // maps name of the program to ProgramInfo
	programs []ProgSpecs                    // list of compiled programs with specs
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
	if sm.specs.equals(&specs) {
		return false, nil
	}

	// Search for compiled program with the specified specs
	for _, pinfo := range sm.programs {
		if pinfo.specs.equals(&specs) {
			sm.gs.UseProgram(pinfo.program)
			sm.specs = specs
			return true, nil
		}
	}

	// Generates new program with the specified specs
	prog, err := sm.GenProgram(&specs)
	if err != nil {
		return false, err
	}

	// create copy of defines so cached specs don't get modified
	// do this after looking at cached specs so we don't create copies of maps every frame
	if specs.MatDefines != nil {
		specs.MatDefines = maps.Clone(s.MatDefines)
	}
	if specs.GeomDefines != nil {
		specs.GeomDefines = maps.Clone(s.GeomDefines)
	}
	if specs.GrDefines != nil {
		specs.GrDefines = maps.Clone(s.GrDefines)
	}

	// Save specs as current specs, adds new program to the list and activates the program
	sm.specs = specs
	sm.programs = append(sm.programs, ProgSpecs{prog, specs})
	sm.gs.UseProgram(prog)
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
	for name, value := range specs.MatDefines {
		defines[name] = value
	}
	for name, value := range specs.GeomDefines {
		defines[name] = value
	}
	for name, value := range specs.GrDefines {
		defines[name] = value
	}

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

// equals compares two ShaderSpecs and returns true if they are effectively equal.
func (ss *ShaderSpecs) equals(other *ShaderSpecs) bool {
	if ss.Name != other.Name {
		return false
	}
	if other.ShaderUnique {
		return true
	}
	return ss.AmbientLightsMax == other.AmbientLightsMax &&
		ss.DirLightsMax == other.DirLightsMax &&
		ss.PointLightsMax == other.PointLightsMax &&
		ss.SpotLightsMax == other.SpotLightsMax &&
		ss.MatTexturesMax == other.MatTexturesMax &&
		maps.Equal(ss.MatDefines, other.MatDefines) &&
		maps.Equal(ss.GeomDefines, other.GeomDefines) &&
		maps.Equal(ss.GrDefines, other.GrDefines)
}
