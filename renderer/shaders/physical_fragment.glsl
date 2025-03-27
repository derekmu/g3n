precision highp float;

// Input variables
in vec3 Position;
in vec3 Normal;
in vec3 CamDir;
in vec2 FragTexcoord;

// Output variables
out vec4 FragColor;

// Material parameters
uniform vec4 uMaterial[3];
#define uBaseColor          uMaterial[0]
#define uEmissiveColor      uMaterial[1].rgb
#define uMetallicFactor     uMaterial[2].x
#define uRoughnessFactor    uMaterial[2].y

#include <pbr>

void main() {
    vec4 color = pbr(uBaseColor, uEmissiveColor, uRoughnessFactor, uMetallicFactor);

    FragColor = vec4(pow(color.rgb, vec3(1.0 / 2.2)), color.a);
}
