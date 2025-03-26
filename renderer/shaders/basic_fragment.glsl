precision highp float;

// Input variables
in vec3 Color;

// Output variables
out vec4 FragColor;

void main() {
    FragColor = vec4(Color, 1.0);
}
