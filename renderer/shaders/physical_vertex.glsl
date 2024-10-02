// Physically Based Shading of a microfacet surface material - Vertex Shader
// Modified from reference implementation at https://github.com/KhronosGroup/glTF-WebGL-PBR

// Vertex attributes
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec3 VertexColor;
layout (location = 3) in vec2 VertexTexcoord;

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

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
out vec3 Position;
out vec3 Normal;
out vec3 CamDir;
out vec2 FragTexcoord;

void main() {
    vec3 vPosition = VertexPosition;

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

    mat4 finalWorld = mat4(1.0);

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

    // Transform this vertex position to camera coordinates.
    Position = vec3(ModelViewMatrix * finalWorld * vec4(vPosition, 1.0));

    // Transform this vertex normal to camera coordinates.
    Normal = normalize(NormalMatrix * finalWorld * VertexNormal);

    // Calculate the direction vector from the vertex to the camera
    // The camera is at 0,0,0
    CamDir = normalize(-Position.xyz);

    // Output texture coordinates to fragment shader
    FragTexcoord = VertexTexcoord;

    gl_Position = MVP * finalWorld * vec4(vPosition, 1.0);
}
