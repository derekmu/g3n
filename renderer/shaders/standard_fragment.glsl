precision highp float;

// Input variables
in vec4 Position;
in vec3 Normal;
in vec2 FragTexcoord;

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

#include <lights>

void phongModel(vec4 position, vec3 normal, vec3 camDir, vec3 matAmbient, vec3 matDiffuse, out vec3 ambdiff, out vec3 spec) {
    vec3 ambientTotal = vec3(0.0);
    vec3 diffuseTotal = vec3(0.0);
    vec3 specularTotal = vec3(0.0);

    bool noLights = true;
    const float EPS = 0.00001;

    float specular;

    #if AMB_LIGHTS > 0
    noLights = false;
    // Ambient lights
    for (int i = 0; i < AMB_LIGHTS; ++i) {
        ambientTotal += uAmbientLightColor[i] * matAmbient;
    }
    #endif

    #if DIR_LIGHTS > 0
    noLights = false;
    // Directional lights
    for (int i = 0; i < DIR_LIGHTS; ++i) {
        vec3 lightDirection = normalize(uDirLightPosition(i));// Vector from fragment to light source
        float dotNormal = dot(lightDirection, normal);// Dot product between light direction and fragment normal
        if (dotNormal > EPS) {
            // If the fragment is lit
            diffuseTotal += uDirLightColor(i) * matDiffuse * dotNormal;
            #ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), uMatShininess);
            #else
            specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), uMatShininess);
            #endif
            specularTotal += uDirLightColor(i) * uMatSpecularColor * specular;
        }
    }
    #endif

    #if POINT_LIGHTS > 0
    noLights = false;
    // Point lights
    for (int i = 0; i < POINT_LIGHTS; ++i) {
        vec3 lightDirection = uPointLightPosition(i) - vec3(position);// Vector from fragment to light source
        float lightDistance = length(lightDirection);// Distance from fragment to light source
        lightDirection = lightDirection / lightDistance;// Normalize lightDirection
        float dotNormal = dot(lightDirection, normal);// Dot product between light direction and fragment normal
        if (dotNormal > EPS) {
            // If the fragment is lit
            float attenuation = 1.0 / (1.0 + lightDistance * (uPointLightLinearDecay(i) + uPointLightQuadraticDecay(i) * lightDistance));
            vec3 attenuatedColor = uPointLightColor(i) * attenuation;
            diffuseTotal += attenuatedColor * matDiffuse * dotNormal;

            #ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), uMatShininess);
            #else
            specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), uMatShininess);
            #endif
            specularTotal += attenuatedColor * uMatSpecularColor * specular;
        }
    }
    #endif

    #if SPOT_LIGHTS > 0
    noLights = false;
    for (int i = 0; i < SPOT_LIGHTS; ++i) {
        // Calculates the direction and distance from the current vertex to this spot light.
        vec3 lightDirection = uSpotLightPosition(i) - vec3(position);// Vector from fragment to light source
        float lightDistance = length(lightDirection);// Distance from fragment to light source
        lightDirection = lightDirection / lightDistance;// Normalize lightDirection
        float angleDot = dot(-lightDirection, uSpotLightDirection(i));
        float angle = acos(angleDot);
        float cutoff = radians(clamp(uSpotLightCutoffAngle(i), 0.0, 90.0));
        if (angle < cutoff) {
            // Check if fragment is inside spotlight beam
            float dotNormal = dot(lightDirection, normal);// Dot product between light direction and fragment normal
            if (dotNormal > EPS) {
                // If the fragment is lit
                float attenuation = 1.0 / (1.0 + lightDistance * (uSpotLightLinearDecay(i) + uSpotLightQuadraticDecay(i) * lightDistance));
                float spotFactor = pow(angleDot, uSpotLightAngularDecay(i));
                vec3 attenuatedColor = uSpotLightColor(i) * attenuation * spotFactor;
                diffuseTotal += attenuatedColor * matDiffuse * dotNormal;

                #ifdef BLINN
                specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), uMatShininess);
                #else
                specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), uMatShininess);
                #endif
                specularTotal += attenuatedColor * uMatSpecularColor * specular;
            }
        }
    }
    #endif
    if (noLights) {
        diffuseTotal = matDiffuse;
    }
    // Sets output colors
    ambdiff = ambientTotal + uMatEmissiveColor + diffuseTotal;
    spec = specularTotal;
}

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

    // Combine material with texture colors
    vec4 matDiffuse = vec4(uMatDiffuseColor, uMatOpacity) * texMixed;
    vec4 matAmbient = vec4(uMatAmbientColor, uMatOpacity) * texMixed;

    // Normalize interpolated normal as it may have shrinked
    vec3 fragNormal = normalize(Normal);

    // Calculate the direction vector from the fragment to the camera (origin)
    vec3 camDir = normalize(-Position.xyz);

    // Workaround for gl_FrontFacing
    vec3 fdx = dFdx(Position.xyz);
    vec3 fdy = dFdy(Position.xyz);
    vec3 faceNormal = normalize(cross(fdx, fdy));
    if (dot(fragNormal, faceNormal) < 0.0) {
        // Back-facing
        fragNormal = -fragNormal;
    }

    // Calculates the Ambient+Diffuse and Specular colors for this fragment using the Phong model.
    vec3 Ambdiff, Spec;
    phongModel(Position, fragNormal, camDir, vec3(matAmbient), vec3(matDiffuse), Ambdiff, Spec);

    // Final fragment color
    FragColor = min(vec4(Ambdiff + Spec, matDiffuse.a), vec4(1.0));
}
