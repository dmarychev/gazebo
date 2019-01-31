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
    vec2 prev_f_total;
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

    vec2 v_half = p.v + 0.5 * dt * p.prev_f_total / p.m;
    p.r = p.r + v_half * dt;
    p.v = v_half + 0.5 * dt * p.f_total / p.m;
    p.prev_f_total = vec2(0,0) + p.f_total; // save from previous step

    current_particles[gid] = p;
}
