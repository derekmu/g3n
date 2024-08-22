// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

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

#ifdef MORPHTARGETS
uniform float morphTargetInfluences[8];
#if MORPHTARGETS > 0
in vec3 MorphPosition0;
#if MORPHTARGETS > 1
in vec3 MorphPosition1;
#if MORPHTARGETS > 2
in vec3 MorphPosition2;
#if MORPHTARGETS > 3
in vec3 MorphPosition3;
#if MORPHTARGETS > 4
in vec3 MorphPosition4;
#if MORPHTARGETS > 5
in vec3 MorphPosition5;
#if MORPHTARGETS > 6
in vec3 MorphPosition6;
#if MORPHTARGETS > 7
in vec3 MorphPosition7;
#endif
#endif
#endif
#endif
#endif
#endif
#endif
#endif
#endif

#ifdef BONE_INFLUENCERS
#if BONE_INFLUENCERS > 0
uniform mat4 mBones[TOTAL_BONES];
in vec4 matricesIndices;
in vec4 matricesWeights;
#endif
#endif

// Output variables for Fragment shader
out vec4 Position;
out vec3 Normal;
out vec2 FragTexcoord;

void main() {
    // Transform vertex position to camera coordinates
    Position = ModelViewMatrix * vec4(VertexPosition, 1.0);

    // Transform vertex normal to camera coordinates
    Normal = normalize(NormalMatrix * VertexNormal);

    vec2 texcoord = VertexTexcoord;
    #if MAT_TEXTURES > 0
    // Flip texture coordinate Y if requested.
    if (MatTexFlipY(0)) {
        texcoord.y = 1.0 - texcoord.y;
    }
    #endif
    FragTexcoord = texcoord;
    vec3 vPosition = VertexPosition;
    mat4 finalWorld = mat4(1.0);

    #ifdef MORPHTARGETS
    #if MORPHTARGETS > 0
    vPosition += MorphPosition0 * morphTargetInfluences[0];
    #if MORPHTARGETS > 1
    vPosition += MorphPosition1 * morphTargetInfluences[1];
    #if MORPHTARGETS > 2
    vPosition += MorphPosition2 * morphTargetInfluences[2];
    #if MORPHTARGETS > 3
    vPosition += MorphPosition3 * morphTargetInfluences[3];
    #if MORPHTARGETS > 4
    vPosition += MorphPosition4 * morphTargetInfluences[4];
    #if MORPHTARGETS > 5
    vPosition += MorphPosition5 * morphTargetInfluences[5];
    #if MORPHTARGETS > 6
    vPosition += MorphPosition6 * morphTargetInfluences[6];
    #if MORPHTARGETS > 7
    vPosition += MorphPosition7 * morphTargetInfluences[7];
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif

    #ifdef BONE_INFLUENCERS
    #if BONE_INFLUENCERS > 0
    mat4 influence = mBones[int(matricesIndices[0])] * matricesWeights[0];
    #if BONE_INFLUENCERS > 1
    influence += mBones[int(matricesIndices[1])] * matricesWeights[1];
    #if BONE_INFLUENCERS > 2
    influence += mBones[int(matricesIndices[2])] * matricesWeights[2];
    #if BONE_INFLUENCERS > 3
    influence += mBones[int(matricesIndices[3])] * matricesWeights[3];
    #endif
    #endif
    #endif
    finalWorld = finalWorld * influence;
    #endif
    #endif

    // Output projected and transformed vertex position
    gl_Position = MVP * finalWorld * vec4(vPosition, 1.0);
}
