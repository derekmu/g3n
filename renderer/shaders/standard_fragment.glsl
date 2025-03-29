precision highp float;

// Input variables
in vec3 Position;
in vec3 Normal;
in vec3 CamDir;
in vec2 FragTexcoord;
in vec3 VPosition;

// Output variables
out vec4 FragColor;

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
#define uMatTexVisible(a)    bool(uMatTexInfo[(3 * a) + 2].x)
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

#include <phong>

void main() {
    // Compute final texture color
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES > 0
    bool firstTex = true;
    if (uMatTexVisible(0)) {
        vec4 texColor = texture(uMatTexture[0], FragTexcoord * uMatTexRepeat(0) + uMatTexOffset(0));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 1
    if (uMatTexVisible(1)) {
        vec4 texColor = texture(uMatTexture[1], FragTexcoord * uMatTexRepeat(1) + uMatTexOffset(1));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 2
    if (uMatTexVisible(2)) {
        vec4 texColor = texture(uMatTexture[2], FragTexcoord * uMatTexRepeat(2) + uMatTexOffset(2));
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

    vec4 matDiffuse = vec4(uMatDiffuseColor, uMatOpacity) * texMixed;
    vec4 matAmbient = vec4(uMatAmbientColor, uMatOpacity) * texMixed;

    vec3 ambientColor, specularColor;
    phong(vec3(matAmbient), vec3(matDiffuse), uMatShininess, uMatSpecularColor, uMatEmissiveColor, ambientColor, specularColor);

    FragColor = min(vec4(ambientColor + specularColor, matDiffuse.a), vec4(1.0));
}
