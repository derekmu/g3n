#include <lights>

void phong(vec3 matAmbient, vec3 matDiffuse, float matShininess, vec3 matSpecularColor, vec3 matEmissiveColor, out vec3 ambientColor, out vec3 specularColor) {
    const float EPS = 0.00001;

    vec3 normal = normalize(Normal);
    vec3 fdx = dFdx(Position.xyz);
    vec3 fdy = dFdy(Position.xyz);
    vec3 faceNormal = normalize(cross(fdx, fdy));
    if (dot(normal, faceNormal) < 0.0) {
        normal = -normal;
    }

    vec3 ambientTotal = vec3(0.0);
    vec3 diffuseTotal = vec3(0.0);
    vec3 specularTotal = vec3(0.0);
    float specular;
    bool noLights = true;

    #if AMB_LIGHTS > 0
    noLights = false;
    for (int i = 0; i < AMB_LIGHTS; ++i) {
        ambientTotal += uAmbientLightColor[i] * matAmbient;
    }
    #endif

    #if DIR_LIGHTS > 0
    noLights = false;
    for (int i = 0; i < DIR_LIGHTS; ++i) {
        vec3 lightDirection = normalize(uDirLightPosition(i));
        float dotNormal = dot(lightDirection, normal);
        if (dotNormal > EPS) {
            diffuseTotal += uDirLightColor(i) * matDiffuse * dotNormal;
            #ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + CamDir)), 0.0), matShininess);
            #else
            specular = pow(max(dot(reflect(-lightDirection, normal), CamDir), 0.0), matShininess);
            #endif
            specularTotal += uDirLightColor(i) * matSpecularColor * specular;
        }
    }
    #endif

    #if POINT_LIGHTS > 0
    noLights = false;
    for (int i = 0; i < POINT_LIGHTS; ++i) {
        vec3 lightDirection = uPointLightPosition(i) - vec3(Position);
        float lightDistance = length(lightDirection);
        lightDirection = lightDirection / lightDistance;
        float dotNormal = dot(lightDirection, normal);
        if (dotNormal > EPS) {
            float attenuation = 1.0 / (1.0 + lightDistance * (uPointLightLinearDecay(i) + uPointLightQuadraticDecay(i) * lightDistance));
            vec3 attenuatedColor = uPointLightColor(i) * attenuation;
            diffuseTotal += attenuatedColor * matDiffuse * dotNormal;
            #ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + CamDir)), 0.0), matShininess);
            #else
            specular = pow(max(dot(reflect(-lightDirection, normal), CamDir), 0.0), matShininess);
            #endif
            specularTotal += attenuatedColor * matSpecularColor * specular;
        }
    }
    #endif

    #if SPOT_LIGHTS > 0
    noLights = false;
    for (int i = 0; i < SPOT_LIGHTS; ++i) {
        vec3 lightDirection = uSpotLightPosition(i) - vec3(Position);
        float lightDistance = length(lightDirection);
        lightDirection = lightDirection / lightDistance;
        float angleDot = dot(-lightDirection, uSpotLightDirection(i));
        float angle = acos(angleDot);
        float cutoff = radians(clamp(uSpotLightCutoffAngle(i), 0.0, 90.0));
        if (angle < cutoff) {
            float dotNormal = dot(lightDirection, normal);
            if (dotNormal > EPS) {
                float attenuation = 1.0 / (1.0 + lightDistance * (uSpotLightLinearDecay(i) + uSpotLightQuadraticDecay(i) * lightDistance));
                float spotFactor = pow(angleDot, uSpotLightAngularDecay(i));
                vec3 attenuatedColor = uSpotLightColor(i) * attenuation * spotFactor;
                diffuseTotal += attenuatedColor * matDiffuse * dotNormal;
                #ifdef BLINN
                specular = pow(max(dot(normal, normalize(lightDirection + CamDir)), 0.0), matShininess);
                #else
                specular = pow(max(dot(reflect(-lightDirection, normal), CamDir), 0.0), matShininess);
                #endif
                specularTotal += attenuatedColor * matSpecularColor * specular;
            }
        }
    }
    #endif

    if (noLights) {
        diffuseTotal = matDiffuse;
    }

    ambientColor = ambientTotal + matEmissiveColor + diffuseTotal;
    specularColor = specularTotal;
}