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

uniform float g = 0.098;

layout(binding=0) buffer Particles {
    Particle current_particles[];
};

void main()
{
    uint gid = gl_GlobalInvocationID.x;
    Particle p = current_particles[gid];

    p.f_grav = vec2(0, -p.m * g);
    p.f_total = p.f_vis + p.f_press + p.f_grav;

    current_particles[gid] = p;
}
