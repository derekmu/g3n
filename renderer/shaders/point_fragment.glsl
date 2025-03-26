precision highp float;

// Input variables
in vec3 Color;
flat in mat2 Rotation;

// Output variables
out vec4 FragColor;

// Material parameters uniform array
uniform vec3 uMaterial[6];
#define uMatAmbientColor     uMaterial[0]
#define uMatDiffuseColor     uMaterial[1]
#define uMatSpecularColor    uMaterial[2]
#define uMatEmissiveColor    uMaterial[3]
#define uMatShininess        uMaterial[4].x
#define uMatOpacity          uMaterial[4].y
#define uMatPointSize        uMaterial[4].z
#define uMatPointRotationZ   uMaterial[5].x

#if MAT_TEXTURES > 0
// Texture unit sampler array
uniform sampler2D uMatTexture[MAT_TEXTURES];
// Texture parameters (3*vec2 per texture)
uniform vec2 uMatTexInfo[3 * MAT_TEXTURES];
#define uMatTexOffset(a)     uMatTexInfo[(3 * a)]
#define uMatTexRepeat(a)     uMatTexInfo[(3 * a) + 1]
#define uMatTexFlipY(a)      bool(uMatTexInfo[(3 * a) + 2].x)
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
    // Compute final texture color
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES > 0
    vec2 pointCoord = Rotation * gl_PointCoord - vec2(0.5) + vec2(0.5);
    bool firstTex = true;
    if (uMatTexVisible(0)) {
        vec4 texColor = texture(uMatTexture[0], pointCoord * uMatTexRepeat(0) + uMatTexOffset(0));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 1
    if (uMatTexVisible(1)) {
        vec4 texColor = texture(uMatTexture[1], pointCoord * uMatTexRepeat(1) + uMatTexOffset(1));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 2
    if (uMatTexVisible(2)) {
        vec4 texColor = texture(uMatTexture[2], pointCoord * uMatTexRepeat(2) + uMatTexOffset(2));
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
    FragColor = min(vec4(Color, uMatOpacity) * texMixed, vec4(1));
}
