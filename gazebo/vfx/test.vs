#version 460

#pragma optimize(off)

in vec3 p_location;

void main() {
    gl_Position = vec4(p_location.xyz, 1);
    gl_PointSize = 4.0;
}
