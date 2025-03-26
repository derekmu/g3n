#ifdef TOTAL_BONES
uniform mat4 uBones[TOTAL_BONES];
in vec4 matricesIndices;
in vec4 matricesWeights;
#endif

void bones(inout mat4 finalWorld, inout mat3 finalNormal) {
    #ifdef TOTAL_BONES
    mat4 influence = mat4(0.0);
    mat3 normalInfluence = mat3(0.0);
    for (int i = 0; i < 4; i++) {
        mat4 bone = uBones[int(matricesIndices[i])];
        float weight = matricesWeights[i];
        influence += bone * weight;
        mat3 boneNormal = mat3(transpose(inverse(bone)));
        normalInfluence += boneNormal * weight;
    }
    finalWorld = finalWorld * influence;
    finalNormal = finalNormal * normalInfluence;
    #endif
}