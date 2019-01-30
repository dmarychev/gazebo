// leapfrog integration of movement equations
#version 460
#pragma optimize(off)

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct Particle {
    vec2 r;
    vec2 v;
    vec2 f_vis;
    vec2 f_press;
    vec2 f_grav;
    vec2 f_total;
    float p; // pressure
    float d; // density
    float m; // mass
    float t; // time
};

uniform float dt = 0.01;

layout(binding=0) buffer Particles {
    Particle current_particles[];
};

void main()
{
    uint gid = gl_GlobalInvocationID.x;
    Particle p = current_particles[gid];

    vec2 a = p.f_total / p.m;
//    vec2 v_half = p.v + a * 0.5 * dt;
//    p.r += v_half * dt;
//    p.v = v_half + a * 0.5 * dt;
    p.r = p.r + p.v * dt + a * dt * dt / 2.0;
    p.v = p.v + a * dt;

    current_particles[gid] = p;
}
