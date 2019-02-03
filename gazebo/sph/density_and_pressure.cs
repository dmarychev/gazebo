// calculate density and pressure
#version 460
//#pragma optimize(off)

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct Particle {
    vec2 r;
    vec2 v;
    vec2 f;
    vec2 prev_f;
    float p; // pressure
    float d; // density
    float m; // mass
};

const float PI = 3.1415926535897932384626433832795;

uniform float h = 0.05;
uniform float k = 0.1;

layout(std430, binding=0) buffer Particles {
    Particle current_particles[];
};

float W(float r, float h) {
    return r < h ? (315.0/(64.0 * PI * pow(h, 9))) * pow(pow(h, 2) - pow(r, 2), 3) : 0.0;
}

void main()
{
    uint gid = gl_GlobalInvocationID.x;
    Particle p = current_particles[gid];

    uint num_particles = (gl_NumWorkGroups * gl_WorkGroupSize).x;

    p.d = .0f;
    for (uint p_i = 0; p_i < num_particles; ++p_i) {
        Particle o = current_particles[p_i];
        p.d += o.m * W(length(p.r - o.r), h);
    }

    p.p = k * p.d;

    current_particles[gid] = p;
}
