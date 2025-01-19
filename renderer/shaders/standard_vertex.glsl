// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Model uniforms
uniform mat4 uMatrices[3];
#define uModelViewMatrix           uMatrices[0]
#define uModelViewProjectionMatrix uMatrices[1]
#define uNormalMatrix              mat3(uMatrices[2])

// Material parameters uniform array
#if MAT_TEXTURES > 0
uniform vec2 uMatTexInfo[3 * MAT_TEXTURES];
#define uMatTexFlipY(a)      bool(uMatTexInfo[(3 * a) + 2].x)
#endif

#ifdef MORPHTARGETS
uniform float uMorphWeights[8];
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

#ifdef TOTAL_BONES
uniform mat4 Bones[TOTAL_BONES];
in vec4 matricesIndices;
in vec4 matricesWeights;
#endif

// Output variables for Fragment shader
out vec4 Position;
out vec3 Normal;
out vec2 FragTexcoord;

void main() {
    vec3 vPosition = VertexPosition;

    #ifdef MORPHTARGETS
    #if MORPHTARGETS > 0
    vPosition += MorphPosition0 * uMorphWeights[0];
    #if MORPHTARGETS > 1
    vPosition += MorphPosition1 * uMorphWeights[1];
    #if MORPHTARGETS > 2
    vPosition += MorphPosition2 * uMorphWeights[2];
    #if MORPHTARGETS > 3
    vPosition += MorphPosition3 * uMorphWeights[3];
    #if MORPHTARGETS > 4
    vPosition += MorphPosition4 * uMorphWeights[4];
    #if MORPHTARGETS > 5
    vPosition += MorphPosition5 * uMorphWeights[5];
    #if MORPHTARGETS > 6
    vPosition += MorphPosition6 * uMorphWeights[6];
    #if MORPHTARGETS > 7
    vPosition += MorphPosition7 * uMorphWeights[7];
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif
    #endif

    mat4 finalWorld = mat4(1.0);
    mat3 finalNormal = mat3(1.0);

    #ifdef TOTAL_BONES
    mat4 influence = mat4(0.0);
    mat3 normalInfluence = mat3(0.0);
    for (int i = 0; i < 4; i++) {
        mat4 bone = Bones[int(matricesIndices[i])];
        float weight = matricesWeights[i];
        influence += bone * weight;
        mat3 boneNormal = mat3(transpose(inverse(bone)));
        normalInfluence += boneNormal * weight;
    }
    finalWorld = finalWorld * influence;
    finalNormal = finalNormal * normalInfluence;
    #endif

    // Transform this vertex position to camera coordinates.
    Position = uModelViewMatrix * finalWorld * vec4(vPosition, 1.0);

    // Transform this vertex normal to camera coordinates.
    Normal = normalize(uNormalMatrix * finalNormal * VertexNormal);

    vec2 texcoord = VertexTexcoord;
    #if MAT_TEXTURES > 0
    // Flip texture coordinate Y if requested.
    if (uMatTexFlipY(0)) {
        texcoord.y = 1.0 - texcoord.y;
    }
    #endif
    FragTexcoord = texcoord;

    // Output projected and transformed vertex position
    gl_Position = uModelViewProjectionMatrix * finalWorld * vec4(vPosition, 1.0);
}
