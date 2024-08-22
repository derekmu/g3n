precision highp float;

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

// Inputs from vertex shader
in vec3 Color;
flat in mat2 Rotation;

// Output
out vec4 FragColor;

void main() {
    // Compute final texture color
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES > 0
    vec2 pointCoord = Rotation * gl_PointCoord - vec2(0.5) + vec2(0.5);
    bool firstTex = true;
    if (MatTexVisible(0)) {
        vec4 texColor = texture(MatTexture[0], pointCoord * MatTexRepeat(0) + MatTexOffset(0));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 1
    if (MatTexVisible(1)) {
        vec4 texColor = texture(MatTexture[1], pointCoord * MatTexRepeat(1) + MatTexOffset(1));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 2
    if (MatTexVisible(2)) {
        vec4 texColor = texture(MatTexture[2], pointCoord * MatTexRepeat(2) + MatTexOffset(2));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #endif
    #endif
    #endif

    // Generates final color
    FragColor = min(vec4(Color, MatOpacity) * texMixed, vec4(1));
}
