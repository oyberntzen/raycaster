#version 430 core

#define MAP_WIDTH 20
#define MAP_HEIGHT 20

#define SCREEN_WIDTH 700
#define SCREEN_HEIGHT 700

layout(local_size_x = 1) in;
layout(std430, binding = 0) buffer allRays
{
    int rays[];
};

uniform int map[MAP_WIDTH*MAP_HEIGHT];
uniform vec2 camPos;
uniform vec2 camDir;
uniform vec2 camPlane;

void main() {
    int x = int(gl_GlobalInvocationID.x);
    
    float cameraX = 2 * x / float(SCREEN_WIDTH) - 1;
    vec2 rayDir = camDir + (camPlane * cameraX);

    ivec2 mapPos = ivec2(camPos);

    vec2 sideDist;
    
    vec2 deltaDist = vec2(1e30, 1e30);
    if (rayDir.x != 0) {
        deltaDist.x = abs(1 / rayDir.x);
    }
    if (rayDir.y != 0) {
        deltaDist.y = abs(1 / rayDir.y);
    }
    float perpWallDist;

    ivec2 stepDir;

    int hit = 0;
    int side;

    if (rayDir.x < 0) {
        stepDir.x = -1;
        sideDist.x = (camPos.x - mapPos.x) * deltaDist.x;
    } else {
        stepDir.x = 1;
        sideDist.x = (mapPos.x + 1 - camPos.x) * deltaDist.x;        
    }
    if (rayDir.y < 0) {
        stepDir.y = -1;
        sideDist.y = (camPos.y - mapPos.y) * deltaDist.y;
    } else {
        stepDir.y = 1;
        sideDist.y = (mapPos.y + 1 - camPos.y) * deltaDist.y;        
    }

    while (hit == 0) {
        if (sideDist.x < sideDist.y) {
            sideDist.x += deltaDist.x;
            mapPos.x += stepDir.x;
            side = 0;
        } else {
            sideDist.y += deltaDist.y;
            mapPos.y += stepDir.y;
            side = 1;
        }

        if (map[mapPos.y*MAP_WIDTH+mapPos.x] > 0) {
            hit = 1;
        }
    }

    if (side == 0) {
        perpWallDist = (sideDist.x - deltaDist.x);
    } else {
        perpWallDist = (sideDist.y - deltaDist.y);
    }

    int lineHeight = int(SCREEN_HEIGHT / perpWallDist);
    rays[x] = lineHeight;
}