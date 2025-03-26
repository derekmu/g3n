// Lights uniforms
#if AMB_LIGHTS > 0
// Ambient lights color uniform
uniform vec3 uAmbientLightColor[AMB_LIGHTS];
#endif

#if DIR_LIGHTS > 0
// Directional lights uniform array. Each directional light uses 2 elements
uniform vec3 uDirLight[2 * DIR_LIGHTS];
#define uDirLightColor(a)    uDirLight[2 * a]
#define uDirLightPosition(a) uDirLight[2 * a + 1]
#endif

#if POINT_LIGHTS > 0
// Point lights uniform array. Each point light uses 3 elements
uniform vec3 uPointLight[3 * POINT_LIGHTS];
#define uPointLightColor(a)          uPointLight[3 * a]
#define uPointLightPosition(a)       uPointLight[3 * a + 1]
#define uPointLightLinearDecay(a)    uPointLight[3 * a + 2].x
#define uPointLightQuadraticDecay(a) uPointLight[3 * a + 2].y
#endif

#if SPOT_LIGHTS > 0
// Spot lights uniforms. Each spot light uses 5 elements
uniform vec3 uSpotLight[5 * SPOT_LIGHTS];
#define uSpotLightColor(a)           uSpotLight[5 * a]
#define uSpotLightPosition(a)        uSpotLight[5 * a + 1]
#define uSpotLightDirection(a)       uSpotLight[5 * a + 2]
#define uSpotLightAngularDecay(a)    uSpotLight[5 * a + 3].x
#define uSpotLightCutoffAngle(a)     uSpotLight[5 * a + 3].y
#define uSpotLightLinearDecay(a)     uSpotLight[5 * a + 3].z
#define uSpotLightQuadraticDecay(a)  uSpotLight[5 * a + 4].x
#endif