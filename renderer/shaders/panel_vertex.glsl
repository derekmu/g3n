// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Model uniforms
uniform mat4 ModelMatrix;

// Outputs for fragment shader
out vec2 FragTexcoord;

// Texture uniforms
uniform vec2 MatTexinfo[3];

#define MatTexFlipY            bool(MatTexinfo[2].x)

void main() {
    vec2 texcoord = VertexTexcoord;
    // Flip texture coordinate Y if requested.
    if (MatTexFlipY) {
        texcoord.y = 1.0 - texcoord.y;
    }
    FragTexcoord = texcoord;

    // Set position
    vec4 pos = vec4(VertexPosition.xyz, 1);
    gl_Position = ModelMatrix * pos;
}
