// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Model uniforms
uniform mat4 uMatrices[3];
#define uModelViewProjectionMatrix uMatrices[1]

// Final output color for fragment shader
out vec3 Color;

void main() {
    Color = VertexColor;
    gl_Position = uModelViewProjectionMatrix * vec4(VertexPosition, 1.0);
}
