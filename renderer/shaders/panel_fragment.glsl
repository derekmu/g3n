precision highp float;

// Texture uniforms
uniform sampler2D uMatTexture;
uniform vec2 uMatTexInfo[3];
#define uMatTexOffset    uMatTexInfo[0]
#define uMatTexRepeat    uMatTexInfo[1]

// Inputs from vertex shader
in vec2 FragTexcoord;

// Input uniform
uniform vec4 uPanel[8];
#define uBounds          uPanel[0]// panel bounds in texture coordinates
#define uBorder          uPanel[1]// panel border in texture coordinates
#define uPadding         uPanel[2]// panel padding in texture coordinates
#define uContent         uPanel[3]// panel content area in texture coordinates
#define uBorderColor     uPanel[4]// panel border color
#define uPaddingColor    uPanel[5]// panel padding color
#define uContentColor    uPanel[6]// panel content color
#define uTextureValid    bool(uPanel[7].x)// texture valid flag

// Output
out vec4 FragColor;

/***
* Checks if current fragment texture coordinate is inside the
* supplied rectangle in texture coordinates:
* rect[0] - position x [0,1]
* rect[1] - position y [0,1]
* rect[2] - width [0,1]
* rect[3] - height [0,1]
*/
bool checkRect(vec4 rect) {
    if (FragTexcoord.x < rect[0]) {
        return false;
    }
    if (FragTexcoord.x > rect[0] + rect[2]) {
        return false;
    }
    if (FragTexcoord.y < rect[1]) {
        return false;
    }
    if (FragTexcoord.y > rect[1] + rect[3]) {
        return false;
    }
    return true;
}


void main() {
    // Discard fragment outside of received bounds
    // Bounds[0] - xmin
    // Bounds[1] - ymin
    // Bounds[2] - xmax
    // Bounds[3] - ymax
    if (FragTexcoord.x <= uBounds[0] || FragTexcoord.x >= uBounds[2]) {
        discard;
    }
    if (FragTexcoord.y <= uBounds[1] || FragTexcoord.y >= uBounds[3]) {
        discard;
    }

    // Check if fragment is inside content area
    if (checkRect(uContent)) {
        // If no texture, the color will be the material color.
        vec4 color = uContentColor;

        if (uTextureValid) {
            // Adjust texture coordinates to fit texture inside the content area
            vec2 offset = vec2(-uContent[0], -uContent[1]);
            vec2 factor = vec2(1.0 / uContent[2], 1.0 / uContent[3]);
            vec2 texcoord = (FragTexcoord + offset) * factor;
            vec4 texColor = texture(uMatTexture, texcoord * uMatTexRepeat + uMatTexOffset);

            // Mix content color with texture color.
            // Note that doing a simple linear interpolation (e.g. using mix()) is not correct!
            // The right formula can be found here: https://en.wikipedia.org/wiki/Alpha_compositing#Alpha_blending
            // For a more in-depth discussion: http://apoorvaj.io/alpha-compositing-opengl-blending-and-premultiplied-alpha.html#toc4
            // Another great discussion here: https://ciechanow.ski/alpha-compositing/

            // Alpha premultiply the content color
            vec4 contentPre = uContentColor;
            contentPre.rgb *= contentPre.a;

            // Alpha premultiply the content color
            vec4 texPre = texColor;
            texPre.rgb *= texPre.a;

            // Combine colors to obtain the alpha premultiplied final color
            color = texPre + contentPre * (1.0 - texPre.a);

            // Un-alpha-premultiply
            color.rgb /= color.a;
        }

        FragColor = color;
        return;
    }

    // Checks if fragment is inside paddings area
    if (checkRect(uPadding)) {
        FragColor = uPaddingColor;
        return;
    }

    // Checks if fragment is inside borders area
    if (checkRect(uBorder)) {
        FragColor = uBorderColor;
        return;
    }

    // Fragment is in margins area (always transparent)
    FragColor = vec4(1, 1, 1, 0);
}
