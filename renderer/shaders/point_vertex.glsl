// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Output variables
out vec3 Color;
flat out mat2 Rotation;

// Model uniforms
uniform mat4 uMatrices[3];
#define uModelViewMatrix           uMatrices[0]
#define uModelViewProjectionMatrix uMatrices[1]

// Material parameters
uniform vec3 uMaterial[6];
#define uMatAmbientColor     uMaterial[0]
#define uMatDiffuseColor     uMaterial[1]
#define uMatSpecularColor    uMaterial[2]
#define uMatEmissiveColor    uMaterial[3]
#define uMatShininess        uMaterial[4].x
#define uMatOpacity          uMaterial[4].y
#define uMatPointSize        uMaterial[4].z
#define uMatPointRotationZ   uMaterial[5].x

// Texture parameters
#if MAT_TEXTURES > 0
uniform sampler2D uMatTexture[MAT_TEXTURES];
uniform vec2 uMatTexInfo[3 * MAT_TEXTURES];
#define uMatTexOffset(a)     uMatTexInfo[(3 * a)]
#define uMatTexRepeat(a)     uMatTexInfo[(3 * a) + 1]
#define uMatTexVisible(a)    bool(uMatTexInfo[(3 * a) + 2].y)
// Alpha compositing (see here: https://ciechanow.ski/alpha-compositing/)
vec4 Blend(vec4 texMixed, vec4 texColor) {
    texMixed.rgb *= texMixed.a;
    texColor.rgb *= texColor.a;
    texMixed = texColor + texMixed * (1 - texColor.a);
    if (texMixed.a > 0.0) {
        texMixed.rgb /= texMixed.a;
    }
    return texMixed;
}
#endif

void main() {
    // Rotation matrix for fragment shader
    float rotSin = sin(uMatPointRotationZ);
    float rotCos = cos(uMatPointRotationZ);
    Rotation = mat2(rotCos, rotSin, -rotSin, rotCos);

    // Sets the vertex position
    vec4 pos = uModelViewProjectionMatrix * vec4(VertexPosition, 1.0);
    gl_Position = pos;

    // Sets the size of the rasterized point decreasing with distance
    vec4 posMV = uModelViewMatrix * vec4(VertexPosition, 1.0);
    gl_PointSize = uMatPointSize / -posMV.z;

    // Outputs color
    Color = uMatEmissiveColor;
}

