// compute shader
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
};

layout(binding=0) buffer Particles {
    Particle current_particles[];
};

void main()
{
    vec2 accel = vec2(0, -0.01);

    uint gid = gl_GlobalInvocationID.x;
    float dt = 0.01;

    Particle p = current_particles[gid];
    p.r += p.v * dt + accel * dt * dt / 2.0;
    p.r += accel * dt;

    if (p.r.y > 0.8) {
        vec2 n = vec2(0, -1);
        p.v = reflect(p.v, n);
        p.v *= 0.99;
        p.r = vec2(p.r.x, 0.8);
    }

    if (p.r.y < -0.8) {
        vec2 n = vec2(0, 1);
        p.v = reflect(p.v, n);
        p.v *= 0.99;
        p.r = vec2(p.r.x, -0.8);
    }

    if (p.r.x > 0.8) {
        vec2 n = vec2(-1, 0);
        p.v = reflect(p.v, n);
        p.v *= 0.99;
        p.r = vec2(0.8, p.r.y);
    }

    if (p.r.x < -0.8) {
        vec2 n = vec2(1, 0);
        p.v = reflect(p.v, n);
        p.v *= 0.99;
        p.r = vec2(-0.8, p.r.y);
    }

    current_particles[gid] = p;

    //next_particles[gid].velocity = p.velocity + accel * dt;
    //next_particles[gid].location = p.r + p.velocity * dt + accel * dt * dt / 2.0;
}
