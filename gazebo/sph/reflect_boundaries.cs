// reflect boundaries
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

uniform float damping_coeff = 0.99;

layout(binding=0) buffer Particles {
    Particle current_particles[];
};

void main()
{
    uint gid = gl_GlobalInvocationID.x;
    Particle p = current_particles[gid];

    if (p.r.y > 0.8) {
        vec2 n = vec2(0, -1);
        p.v = reflect(p.v, n);
        p.v *= damping_coeff;
        p.r = vec2(p.r.x, 0.8);
    }

    if (p.r.y < -0.8) {
        vec2 n = vec2(0, 1);
        p.v = reflect(p.v, n);
        p.v *= damping_coeff;
        p.r = vec2(p.r.x, -0.8);
    }

    if (p.r.x > 0.8) {
        vec2 n = vec2(-1, 0);
        p.v = reflect(p.v, n);
        p.v *= damping_coeff;
        p.r = vec2(0.8, p.r.y);
    }

    if (p.r.x < -0.8) {
        vec2 n = vec2(1, 0);
        p.v = reflect(p.v, n);
        p.v *= damping_coeff;
        p.r = vec2(-0.8, p.r.y);
    }

    current_particles[gid] = p;
}
