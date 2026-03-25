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

void morph(inout vec3 vPosition) {
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
}