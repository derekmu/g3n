precision highp float;

// Inputs from vertex shader
in vec4 Position;// Fragment position in camera coordinates
in vec3 Normal;// Fragment normal in camera coordinates
in vec2 FragTexcoord;// Fragment texture coordinates

// Lights uniforms
#if AMB_LIGHTS > 0
// Ambient lights color uniform
uniform vec3 AmbientLightColor[AMB_LIGHTS];
#endif
#if DIR_LIGHTS > 0
// Directional lights uniform array. Each directional light uses 2 elements
uniform vec3 DirLight[2 * DIR_LIGHTS];
// Macros to access elements inside the DirectionalLight uniform array
#define DirLightColor(a)    DirLight[2 * a]
#define DirLightPosition(a) DirLight[2 * a + 1]
#endif
#if POINT_LIGHTS > 0
// Point lights uniform array. Each point light uses 3 elements
uniform vec3 PointLight[3 * POINT_LIGHTS];
// Macros to access elements inside the PointLight uniform array
#define PointLightColor(a)          PointLight[3 * a]
#define PointLightPosition(a)       PointLight[3 * a + 1]
#define PointLightLinearDecay(a)    PointLight[3 * a + 2].x
#define PointLightQuadraticDecay(a) PointLight[3 * a + 2].y
#endif
#if SPOT_LIGHTS > 0
// Spot lights uniforms. Each spot light uses 5 elements
uniform vec3 SpotLight[5 * SPOT_LIGHTS];
// Macros to access elements inside the PointLight uniform array
#define SpotLightColor(a)           SpotLight[5 * a]
#define SpotLightPosition(a)        SpotLight[5 * a + 1]
#define SpotLightDirection(a)       SpotLight[5 * a + 2]
#define SpotLightAngularDecay(a)    SpotLight[5 * a + 3].x
#define SpotLightCutoffAngle(a)     SpotLight[5 * a + 3].y
#define SpotLightLinearDecay(a)     SpotLight[5 * a + 3].z
#define SpotLightQuadraticDecay(a)  SpotLight[5 * a + 4].x
#endif

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

/***
 phong lighting model
 Parameters:
    position:   input vertex position in camera coordinates
    normal:     input vertex normal in camera coordinates
    camDir:     input camera directions
    matAmbient: input material ambient color
    matDiffuse: input material diffuse color
    ambdiff:    output ambient+diffuse color
    spec:       output specular color
 Uniforms:
    AmbientLightColor[]
    DiffuseLightColor[]
    DiffuseLightPosition[]
    PointLightColor[]
    PointLightPosition[]
    PointLightLinearDecay[]
    PointLightQuadraticDecay[]
    MatSpecularColor
    MatShininess
*****/
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
        ambientTotal += AmbientLightColor[i] * matAmbient;
    }
    #endif

    #if DIR_LIGHTS > 0
    noLights = false;
    // Directional lights
    for (int i = 0; i < DIR_LIGHTS; ++i) {
        vec3 lightDirection = normalize(DirLightPosition(i)); // Vector from fragment to light source
        float dotNormal = dot(lightDirection, normal); // Dot product between light direction and fragment normal
        if (dotNormal > EPS) { // If the fragment is lit
                               diffuseTotal += DirLightColor(i) * matDiffuse * dotNormal;

                               #ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), MatShininess);
                               #else
            specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
                               #endif
            specularTotal += DirLightColor(i) * MatSpecularColor * specular;
        }
    }
    #endif

    #if POINT_LIGHTS > 0
    noLights = false;
    // Point lights
    for (int i = 0; i < POINT_LIGHTS; ++i) {
        vec3 lightDirection = PointLightPosition(i) - vec3(position); // Vector from fragment to light source
        float lightDistance = length(lightDirection); // Distance from fragment to light source
        lightDirection = lightDirection / lightDistance; // Normalize lightDirection
        float dotNormal = dot(lightDirection, normal);  // Dot product between light direction and fragment normal
        if (dotNormal > EPS) { // If the fragment is lit
                               float attenuation = 1.0 / (1.0 + lightDistance * (PointLightLinearDecay(i) + PointLightQuadraticDecay(i) * lightDistance));
                               vec3 attenuatedColor = PointLightColor(i) * attenuation;
                               diffuseTotal += attenuatedColor * matDiffuse * dotNormal;

                               #ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), MatShininess);
                               #else
            specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
                               #endif
            specularTotal += attenuatedColor * MatSpecularColor * specular;
        }
    }
    #endif

    #if SPOT_LIGHTS > 0
    noLights = false;
    for (int i = 0; i < SPOT_LIGHTS; ++i) {
        // Calculates the direction and distance from the current vertex to this spot light.
        vec3 lightDirection = SpotLightPosition(i) - vec3(position); // Vector from fragment to light source
        float lightDistance = length(lightDirection); // Distance from fragment to light source
        lightDirection = lightDirection / lightDistance; // Normalize lightDirection
        float angleDot = dot(-lightDirection, SpotLightDirection(i));
        float angle = acos(angleDot);
        float cutoff = radians(clamp(SpotLightCutoffAngle(i), 0.0, 90.0));
        if (angle < cutoff) { // Check if fragment is inside spotlight beam
                              float dotNormal = dot(lightDirection, normal); // Dot product between light direction and fragment normal
                              if (dotNormal > EPS) { // If the fragment is lit
                                                     float attenuation = 1.0 / (1.0 + lightDistance * (SpotLightLinearDecay(i) + SpotLightQuadraticDecay(i) * lightDistance));
                                                     float spotFactor = pow(angleDot, SpotLightAngularDecay(i));
                                                     vec3 attenuatedColor = SpotLightColor(i) * attenuation * spotFactor;
                                                     diffuseTotal += attenuatedColor * matDiffuse * dotNormal;

                                                     #ifdef BLINN
                specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), MatShininess);
                                                     #else
                specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
                                                     #endif
                specularTotal += attenuatedColor * MatSpecularColor * specular;
                              }
        }
    }
    #endif
    if (noLights) {
        diffuseTotal = matDiffuse;
    }
    // Sets output colors
    ambdiff = ambientTotal + MatEmissiveColor + diffuseTotal;
    spec = specularTotal;
}

// Final fragment color
out vec4 FragColor;

void main() {
    // Compute final texture color
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES > 0
    bool firstTex = true;
    if (MatTexVisible(0)) {
        vec4 texColor = texture(MatTexture[0], FragTexcoord * MatTexRepeat(0) + MatTexOffset(0));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 1
    if (MatTexVisible(1)) {
        vec4 texColor = texture(MatTexture[1], FragTexcoord * MatTexRepeat(1) + MatTexOffset(1));
        if (firstTex) {
            texMixed = texColor;
            firstTex = false;
        } else {
            texMixed = Blend(texMixed, texColor);
        }
    }
    #if MAT_TEXTURES > 2
    if (MatTexVisible(2)) {
        vec4 texColor = texture(MatTexture[2], FragTexcoord * MatTexRepeat(2) + MatTexOffset(2));
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
    vec4 matDiffuse = vec4(MatDiffuseColor, MatOpacity) * texMixed;
    vec4 matAmbient = vec4(MatAmbientColor, MatOpacity) * texMixed;

    // Normalize interpolated normal as it may have shrinked
    vec3 fragNormal = normalize(Normal);

    // Calculate the direction vector from the fragment to the camera (origin)
    vec3 camDir = normalize(-Position.xyz);

    // Workaround for gl_FrontFacing
    vec3 fdx = dFdx(Position.xyz);
    vec3 fdy = dFdy(Position.xyz);
    vec3 faceNormal = normalize(cross(fdx, fdy));
    if (dot(fragNormal, faceNormal) < 0.0) { // Back-facing
                                             fragNormal = -fragNormal;
    }

    // Calculates the Ambient+Diffuse and Specular colors for this fragment using the Phong model.
    vec3 Ambdiff, Spec;
    phongModel(Position, fragNormal, camDir, vec3(matAmbient), vec3(matDiffuse), Ambdiff, Spec);

    // Final fragment color
    FragColor = min(vec4(Ambdiff + Spec, matDiffuse.a), vec4(1.0));
}
