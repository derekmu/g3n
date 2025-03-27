// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Output variables
out vec4 Position;
out vec3 Normal;
out vec3 CamDir;
out vec2 FragTexcoord;

// Model uniforms
uniform mat4 uMatrices[3];
#define uModelViewMatrix           uMatrices[0]
#define uModelViewProjectionMatrix uMatrices[1]
#define uNormalMatrix              mat3(uMatrices[2])

// Texture parameters
#if MAT_TEXTURES > 0
uniform vec2 uMatTexInfo[3 * MAT_TEXTURES];
#define uMatTexFlipY(a)      bool(uMatTexInfo[(3 * a) + 2].x)
#endif

#include <morph>
#include <bones>

void main() {
    vec3 vPosition = VertexPosition;
    morph(vPosition);

    mat4 finalWorld = mat4(1.0);
    mat3 finalNormal = mat3(1.0);
    bones(finalWorld, finalNormal);

    Position = uModelViewMatrix * finalWorld * vec4(vPosition, 1.0);
    Normal = normalize(uNormalMatrix * finalNormal * VertexNormal);
    CamDir = normalize(-Position.xyz);
    vec2 texcoord = VertexTexcoord;
    #if MAT_TEXTURES > 0
    if (uMatTexFlipY(0)) {
        texcoord.y = 1.0 - texcoord.y;
    }
    #endif
    FragTexcoord = texcoord;
    gl_Position = uModelViewProjectionMatrix * finalWorld * vec4(vPosition, 1.0);
}
