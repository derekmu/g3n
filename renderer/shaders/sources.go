package shaders

import _ "embed"

//go:embed basic_vertex.glsl
var basicVertexSource string

//go:embed basic_fragment.glsl
var basicFragmentSource string

//go:embed standard_vertex.glsl
var standardVertexSource string

//go:embed standard_fragment.glsl
var standardFragmentSource string

//go:embed physical_vertex.glsl
var physicalVertexSource string

//go:embed physical_fragment.glsl
var physicalFragmentSource string

//go:embed panel_vertex.glsl
var panelVertexSource string

//go:embed panel_fragment.glsl
var panelFragmentSource string

//go:embed point_vertex.glsl
var pointVertexSource string

//go:embed point_fragment.glsl
var pointFragmentSource string

// Maps shader name with its source code
var shaderMap = map[string]string{
	"panel_fragment":    panelFragmentSource,
	"physical_vertex":   physicalVertexSource,
	"basic_vertex":      basicVertexSource,
	"standard_vertex":   standardVertexSource,
	"point_vertex":      pointVertexSource,
	"standard_fragment": standardFragmentSource,
	"point_fragment":    pointFragmentSource,
	"physical_fragment": physicalFragmentSource,
	"basic_fragment":    basicFragmentSource,
	"panel_vertex":      panelVertexSource,
}

// Maps program name with ProgramInfo struct with shaders names
var programMap = map[string]ProgramInfo{
	"basic":    {"basic_vertex", "basic_fragment", ""},
	"panel":    {"panel_vertex", "panel_fragment", ""},
	"physical": {"physical_vertex", "physical_fragment", ""},
	"point":    {"point_vertex", "point_fragment", ""},
	"standard": {"standard_vertex", "standard_fragment", ""},
}
