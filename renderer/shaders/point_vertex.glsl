// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Model uniforms
uniform mat4 MVP;
uniform mat4 MV;

// Material parameters uniform array
uniform vec3 Material[6];
// Macros to access elements inside the Material array
#define MatAmbientColor     Material[0]
#define MatDiffuseColor     Material[1]
#define MatSpecularColor    Material[2]
#define MatEmissiveColor    Material[3]
#define MatShininess        Material[4].x
#define MatOpacity          Material[4].y
#define MatPointSize        Material[4].z
#define MatPointRotationZ   Material[5].x
#if MAT_TEXTURES > 0
// Texture unit sampler array
uniform sampler2D MatTexture[MAT_TEXTURES];
// Texture parameters (3*vec2 per texture)
uniform vec2 MatTexinfo[3 * MAT_TEXTURES];
// Macros to access elements inside the MatTexinfo array
#define MatTexOffset(a)     MatTexinfo[(3 * a)]
#define MatTexRepeat(a)     MatTexinfo[(3 * a) + 1]
#define MatTexFlipY(a)      bool(MatTexinfo[(3 * a) + 2].x)
#define MatTexVisible(a)    bool(MatTexinfo[(3 * a) + 2].y)
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

// Outputs for fragment shader
out vec3 Color;
flat out mat2 Rotation;

void main() {
    // Rotation matrix for fragment shader
    float rotSin = sin(MatPointRotationZ);
    float rotCos = cos(MatPointRotationZ);
    Rotation = mat2(rotCos, rotSin, -rotSin, rotCos);

    // Sets the vertex position
    vec4 pos = MVP * vec4(VertexPosition, 1.0);
    gl_Position = pos;

    // Sets the size of the rasterized point decreasing with distance
    vec4 posMV = MV * vec4(VertexPosition, 1.0);
    gl_PointSize = MatPointSize / -posMV.z;

    // Outputs color
    Color = MatEmissiveColor;
}

