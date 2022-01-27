#version 430 core

#define SCREEN_WIDTH 700
#define SCREEN_HEIGHT 700

layout(std430, binding = 0) buffer allRays
{
    int rays[];
};

out vec4 color;

void main()
{
    int x = int(gl_FragCoord.x);
    int y = int(gl_FragCoord.y);

    color = vec4(0, 0, 0, 1.0);

    if (y >= (SCREEN_HEIGHT - rays[x]) / 2 && y <= (SCREEN_HEIGHT + rays[x]) / 2) {
        color = vec4(1.0, 1.0, 1.0, 1.0);
    }
}