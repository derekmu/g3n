// Physically Based Shading of a microfacet surface material - Fragment Shader
// Modified from reference implementation at https://github.com/KhronosGroup/glTF-WebGL-PBR
//
// References:
// [1] Real Shading in Unreal Engine 4
//     http://blog.selfshadow.com/publications/s2013-shading-course/karis/s2013_pbs_epic_notes_v2.pdf
// [2] Physically Based Shading at Disney
//     http://blog.selfshadow.com/publications/s2012-shading-course/burley/s2012_pbs_disney_brdf_notes_v3.pdf
// [3] README.md - Environment Maps
//     https://github.com/KhronosGroup/glTF-WebGL-PBR/#environment-maps
// [4] "An Inexpensive BRDF Model for Physically based Rendering" by Christophe Schlick
//     https://www.cs.virginia.edu/~jdl/bib/appearance/analytic%20models/schlick94b.pdf

#include <lights>

#ifdef HAS_BASECOLORMAP
uniform sampler2D uBaseColorSampler;
#endif
#ifdef HAS_METALROUGHNESSMAP
uniform sampler2D uMetallicRoughnessSampler;
#endif
#ifdef HAS_NORMALMAP
uniform sampler2D uNormalSampler;
#endif
#ifdef HAS_EMISSIVEMAP
uniform sampler2D uEmissiveSampler;
#endif
#ifdef HAS_OCCLUSIONMAP
uniform sampler2D uOcclusionSampler;
uniform float uOcclusionStrength;
#endif

// Encapsulate the various inputs used by the various functions in the shading equation
// We store values in this struct to simplify the integration of alternative implementations
// of the shading terms, outlined in the Readme.MD Appendix.F
struct PBRLightInfo {
    float NdotL;// cos angle between normal and light direction
    float NdotV;// cos angle between normal and view direction
    float NdotH;// cos angle between normal and half vector
    float LdotH;// cos angle between light direction and half vector
    float VdotH;// cos angle between view direction and half vector
};

struct PBRInfo {
    float perceptualRoughness;// roughness value, as authored by the model creator (input to shader)
    float metalness;// metallic value at the surface
    vec3 reflectance0;// full reflectance color (normal incidence angle)
    vec3 reflectance90;// reflectance color at grazing angle
    float alphaRoughness;// roughness mapped to a more linear change in the roughness (proposed by [2])
    vec3 diffuseColor;// color contribution from diffuse lighting
    vec3 specularColor;// color contribution from specular lighting
};

const float M_PI = 3.141592653589793;
const float c_MinRoughness = 0.04;

vec4 SRGBtoLINEAR(vec4 srgbIn) {
    vec3 bLess = step(vec3(0.04045), srgbIn.xyz);
    vec3 linOut = mix(srgbIn.xyz / vec3(12.92), pow((srgbIn.xyz + vec3(0.055)) / vec3(1.055), vec3(2.4)), bLess);
    return vec4(linOut, srgbIn.w);
}

// Find the normal for this fragment, pulling either from a predefined normal map
// or from the interpolated mesh normal and tangent attributes.
vec3 getNormal() {
    // Retrieve the tangent space matrix
    vec3 pos_dx = dFdx(Position);
    vec3 pos_dy = dFdy(Position);
    vec3 tex_dx = dFdx(vec3(FragTexcoord, 0.0));
    vec3 tex_dy = dFdy(vec3(FragTexcoord, 0.0));
    vec3 t = (tex_dy.t * pos_dx - tex_dx.t * pos_dy) / (tex_dx.s * tex_dy.t - tex_dy.s * tex_dx.t);
    vec3 ng = normalize(Normal);
    t = normalize(t - ng * dot(ng, t));
    vec3 b = normalize(cross(ng, t));
    mat3 tbn = mat3(t, b, ng);

    #ifdef HAS_NORMALMAP
    vec3 n = texture(uNormalSampler, FragTexcoord).rgb;
    n = normalize(tbn * ((2.0 * n - 1.0) * vec3(1.0, 1.0, 1.0)));
    #else
    // The tbn matrix is linearly interpolated, so we need to re-normalize
    vec3 n = normalize(tbn[2].xyz);
    #endif

    return n;
}

// Basic Lambertian diffuse
// Implementation from Lambert's Photometria https://archive.org/details/lambertsphotome00lambgoog
// See also [1], Equation 1
vec3 diffuse(PBRInfo pbrInputs) {
    return pbrInputs.diffuseColor / M_PI;
}

// The following equation models the Fresnel reflectance term of the spec equation (aka F())
// Implementation of fresnel from [4], Equation 15
vec3 specularReflection(PBRInfo pbrInputs, PBRLightInfo pbrLight) {
    return pbrInputs.reflectance0 + (pbrInputs.reflectance90 - pbrInputs.reflectance0) * pow(clamp(1.0 - pbrLight.VdotH, 0.0, 1.0), 5.0);
}

// This calculates the specular geometric attenuation (aka G()),
// where rougher material will reflect less light back to the viewer.
// This implementation is based on [1] Equation 4, and we adopt their modifications to
// alphaRoughness as input as originally proposed in [2].
float geometricOcclusion(PBRInfo pbrInputs, PBRLightInfo pbrLight) {
    float NdotL = pbrLight.NdotL;
    float NdotV = pbrLight.NdotV;
    float r = pbrInputs.alphaRoughness;
    float attenuationL = 2.0 * NdotL / (NdotL + sqrt(r * r + (1.0 - r * r) * (NdotL * NdotL)));
    float attenuationV = 2.0 * NdotV / (NdotV + sqrt(r * r + (1.0 - r * r) * (NdotV * NdotV)));
    return attenuationL * attenuationV;
}

// The following equation(s) model the distribution of microfacet normals across the area being drawn (aka D())
// Implementation from "Average Irregularity Representation of a Roughened Surface for Ray Reflection" by T. S. Trowbridge, and K. P. Reitz
// Follows the distribution function recommended in the SIGGRAPH 2013 course notes from EPIC Games [1], Equation 3.
float microfacetDistribution(PBRInfo pbrInputs, PBRLightInfo pbrLight) {
    float roughnessSq = pbrInputs.alphaRoughness * pbrInputs.alphaRoughness;
    float f = (pbrLight.NdotH * roughnessSq - pbrLight.NdotH) * pbrLight.NdotH + 1.0;
    return roughnessSq / (M_PI * f * f);
}

vec3 pbrModel(PBRInfo pbrInputs, vec3 lightColor, vec3 lightDir) {
    vec3 n = getNormal();// normal at surface point
    vec3 v = normalize(CamDir);// Vector from surface point to camera
    vec3 l = normalize(lightDir);// Vector from surface point to light
    vec3 h = normalize(l + v);// Half vector between both l and v
    vec3 reflection = -normalize(reflect(v, n));

    float NdotL = clamp(dot(n, l), 0.001, 1.0);
    float NdotV = abs(dot(n, v)) + 0.001;
    float NdotH = clamp(dot(n, h), 0.0, 1.0);
    float LdotH = clamp(dot(l, h), 0.0, 1.0);
    float VdotH = clamp(dot(v, h), 0.0, 1.0);

    PBRLightInfo pbrLight = PBRLightInfo(
    NdotL,
    NdotV,
    NdotH,
    LdotH,
    VdotH
    );

    // Calculate the shading terms for the microfacet specular shading model
    vec3 F = specularReflection(pbrInputs, pbrLight);
    float G = geometricOcclusion(pbrInputs, pbrLight);
    float D = microfacetDistribution(pbrInputs, pbrLight);

    // Calculation of analytical lighting contribution
    vec3 diffuseContrib = (1.0 - F) * diffuse(pbrInputs);
    vec3 specContrib = F * G * D / (4.0 * NdotL * NdotV);
    // Obtain final intensity as reflectance (BRDF) scaled by the energy of the light (cosine law)
    vec3 color = NdotL * lightColor * (diffuseContrib + specContrib);

    return color;
}

vec4 pbr(vec4 baseColor, vec3 emissiveColor, float roughnessFactor, float metallicFactor) {
    float perceptualRoughness = roughnessFactor;
    float metallic = metallicFactor;

    #ifdef HAS_METALROUGHNESSMAP
    // Roughness is stored in the 'g' channel, metallic is stored in the 'b' channel.
    // This layout intentionally reserves the 'r' channel for (optional) occlusion map data
    vec4 mrSample = texture(uMetallicRoughnessSampler, FragTexcoord);
    perceptualRoughness = mrSample.g * perceptualRoughness;
    metallic = mrSample.b * metallic;
    #endif

    perceptualRoughness = clamp(perceptualRoughness, c_MinRoughness, 1.0);
    metallic = clamp(metallic, 0.0, 1.0);
    // Roughness is authored as perceptual roughness.
    // Convert to material roughness by squaring the perceptual roughness [2].
    float alphaRoughness = perceptualRoughness * perceptualRoughness;

    // The albedo may be defined from a base texture or a flat color
    #ifdef HAS_BASECOLORMAP
    baseColor = SRGBtoLINEAR(texture(uBaseColorSampler, FragTexcoord)) * baseColor;
    #endif

    vec3 f0 = vec3(0.04);
    vec3 diffuseColor = baseColor.rgb * (vec3(1.0) - f0);
    diffuseColor *= 1.0 - metallic;

    vec3 specularColor = mix(f0, baseColor.rgb, metallicFactor);

    float reflectance = max(max(specularColor.r, specularColor.g), specularColor.b);

    // For typical incident reflectance range (between 4% to 100%) set the grazing reflectance to 100% for typical fresnel effect.
    // For very low reflectance range on highly diffuse objects (below 4%), incrementally reduce grazing reflectance to 0%.
    float reflectance90 = clamp(reflectance * 25.0, 0.0, 1.0);
    vec3 specularEnvironmentR0 = specularColor.rgb;
    vec3 specularEnvironmentR90 = vec3(1.0, 1.0, 1.0) * reflectance90;

    PBRInfo pbrInputs = PBRInfo(
    perceptualRoughness,
    metallic,
    specularEnvironmentR0,
    specularEnvironmentR90,
    alphaRoughness,
    diffuseColor,
    specularColor
    );

    vec3 color = vec3(0.0);

    #if AMB_LIGHTS > 0
    for (int i = 0; i < AMB_LIGHTS; i++) {
        color += uAmbientLightColor[i] * pbrInputs.diffuseColor;
    }
    #endif

    #if DIR_LIGHTS > 0
    for (int i = 0; i < DIR_LIGHTS; i++) {
        // Direction of the current light
        vec3 lightDirection = normalize(uDirLightPosition(i));
        color += pbrModel(pbrInputs, uDirLightColor(i), lightDirection);
    }
    #endif

    #if POINT_LIGHTS > 0
    for (int i = 0; i < POINT_LIGHTS; i++) {
        // Direction and distance from the current vertex to the light
        vec3 lightDirection = uPointLightPosition(i) - vec3(Position);
        float lightDistance = length(lightDirection);
        lightDirection = lightDirection / lightDistance;
        // Attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + uPointLightLinearDecay(i) * lightDistance +
        uPointLightQuadraticDecay(i) * lightDistance * lightDistance);
        vec3 attenuatedColor = uPointLightColor(i) * attenuation;
        color += pbrModel(pbrInputs, attenuatedColor, lightDirection);
    }
    #endif

    #if SPOT_LIGHTS > 0
    for (int i = 0; i < SPOT_LIGHTS; i++) {
        // Direction and distance from the current vertex to the light
        vec3 lightDirection = uSpotLightPosition(i) - vec3(Position);
        float lightDistance = length(lightDirection);
        lightDirection = lightDirection / lightDistance;

        // Attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + uSpotLightLinearDecay(i) * lightDistance +
        uSpotLightQuadraticDecay(i) * lightDistance * lightDistance);

        // Angle between the vertex direction and light direction
        float angle = acos(dot(-lightDirection, uSpotLightDirection(i)));
        float cutoff = radians(clamp(uSpotLightCutoffAngle(i), 0.0, 90.0));
        // If the angle is greater than the cutoff the will not contribute to the final color
        if (angle < cutoff) {
            float spotFactor = pow(dot(-lightDirection, uSpotLightDirection(i)), uSpotLightAngularDecay(i));
            vec3 attenuatedColor = uSpotLightColor(i) * attenuation * spotFactor;
            color += pbrModel(pbrInputs, attenuatedColor, lightDirection);
        }
    }
    #endif

    #ifdef HAS_OCCLUSIONMAP
    float ao = texture(uOcclusionSampler, FragTexcoord).r;
    color = mix(color, color * ao, 1.0);
    #endif

    #ifdef HAS_EMISSIVEMAP
    emissiveColor = SRGBtoLINEAR(texture(uEmissiveSampler, FragTexcoord)).rgb * emissiveColor.rgb;
    #endif
    color += emissiveColor.rgb;

    // Alternative colors for testing:
    // Base Color
    //    FragColor = baseColor;
    // Normal
    //    FragColor = vec4(getNormal(), 1.0);
    // Emissive Color
    //    FragColor = vec4(emissiveColor, 1.0);
    // F
    //    color = F;
    // G
    //    color = vec3(G);
    // D
    //    color = vec3(D);
    // Specular
    //    color = specContrib;
    // Diffuse
    //    color = diffuseContrib;
    // Roughness
    //    color = vec3(perceptualRoughness);
    // Metallic
    //    color = vec3(metallic);

    return vec4(color, baseColor.a);
}