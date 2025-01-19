// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Model uniforms
uniform mat4 uModelMatrix;

// Outputs for fragment shader
out vec2 FragTexcoord;

// Texture uniforms
uniform vec2 uMatTexInfo[3];
#define uMatTexFlipY            bool(uMatTexInfo[2].x)

void main() {
    vec2 texcoord = VertexTexcoord;
    // Flip texture coordinate Y if requested.
    if (uMatTexFlipY) {
        texcoord.y = 1.0 - texcoord.y;
    }
    FragTexcoord = texcoord;

    // Set position
    vec4 pos = vec4(VertexPosition.xyz, 1);
    gl_Position = uModelMatrix * pos;
}
