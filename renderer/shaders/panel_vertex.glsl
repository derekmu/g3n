// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Output variables
out vec2 FragTexcoord;

// Model uniforms
uniform mat4 uModelMatrix;

void main() {
    FragTexcoord = VertexTexcoord;
    vec4 pos = vec4(VertexPosition.xyz, 1);
    gl_Position = uModelMatrix * pos;
}
