package shaders

import _ "embed"

//go:embed lights.glsl
var lightsSource string

//go:embed pbr.glsl
var pbrSource string

//go:embed phong.glsl
var phongSource string

//go:embed bones.glsl
var bonesSource string

//go:embed morph.glsl
var morphSource string

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

func init() {
	AddInclude("lights", lightsSource)
	AddInclude("pbr", pbrSource)
	AddInclude("phong", phongSource)
	AddInclude("bones", bonesSource)
	AddInclude("morph", morphSource)
	AddShader("basic_vertex", basicVertexSource)
	AddShader("basic_fragment", basicFragmentSource)
	AddProgram("basic", "basic_vertex", "basic_fragment", "")
	AddShader("standard_vertex", standardVertexSource)
	AddShader("standard_fragment", standardFragmentSource)
	AddProgram("standard", "standard_vertex", "standard_fragment", "")
	AddShader("physical_vertex", physicalVertexSource)
	AddShader("physical_fragment", physicalFragmentSource)
	AddProgram("physical", "physical_vertex", "physical_fragment", "")
	AddShader("panel_vertex", panelVertexSource)
	AddShader("panel_fragment", panelFragmentSource)
	AddProgram("panel", "panel_vertex", "panel_fragment", "")
	AddShader("point_vertex", pointVertexSource)
	AddShader("point_fragment", pointFragmentSource)
	AddProgram("point", "point_vertex", "point_fragment", "")
}
