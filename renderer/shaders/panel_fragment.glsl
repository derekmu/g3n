precision highp float;

// Input variables
in vec2 FragTexcoord;

// Output variables
out vec4 FragColor;

// Texture uniforms
uniform sampler2D uMatTexture;
uniform vec2 uMatTexInfo[3];
#define uMatTexOffset    uMatTexInfo[0]
#define uMatTexRepeat    uMatTexInfo[1]

// Panel uniform
uniform vec4 uPanel[8];
#define uBounds          uPanel[0]// bounds in texture coordinates
#define uColor           uPanel[1]// panel color
#define uTextureValid    bool(uPanel[2].x)// texture valid flag

void main() {
    // Discard fragment outside of received bounds
    if (FragTexcoord.x <= uBounds[0] || FragTexcoord.x >= uBounds[2]) {
        discard;
    }
    if (FragTexcoord.y <= uBounds[1] || FragTexcoord.y >= uBounds[3]) {
        discard;
    }
    if (uTextureValid) {
        // Use texture
        FragColor = texture(uMatTexture, FragTexcoord * uMatTexRepeat + uMatTexOffset);
    } else {
        // Use material
        FragColor = uColor;
    }
}
