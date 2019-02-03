#version 460

#pragma optimize(off)

in vec2 p_location;

void main() {
    gl_Position = vec4(p_location.xy, 0, 1);
    gl_PointSize = 5.0;
}
