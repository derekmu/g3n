#include <attributes>

// Model uniforms
uniform mat4 ModelMatrix;

// Outputs for fragment shader
out vec2 FragTexcoord;

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
