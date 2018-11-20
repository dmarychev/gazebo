// compute shader
#version 460
#pragma optimize(off)

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct Particle {
    vec4 location;
    vec4 velocity;
};

layout(binding=0) buffer Current {
    Particle current_particles[];
};

layout(binding=1) buffer Next {
    Particle next_particles[];
};


void main()
{
    vec4 accel = vec4(0, -0.01, 0, 0);

    uint gid = gl_GlobalInvocationID.x;
    float dt = 0.01;

    Particle p = current_particles[gid];

    p.location += p.velocity * dt + accel * dt * dt / 2.0;
    p.velocity += accel * dt;

    if (p.location.y > 0.8) {
        vec4 n = vec4(0, -1, 0, 0);
        p.velocity = reflect(p.velocity, n);
        p.velocity *= 0.99;
        p.location = vec4(p.location.x, 0.8, 0, 0);
    }

    if (p.location.y < -0.8) {
        vec4 n = vec4(0, 1, 0, 0);
        p.velocity = reflect(p.velocity, n);
        p.velocity *= 0.99;
        p.location = vec4(p.location.x, -0.8, 0, 0);
    }

    if (p.location.x > 0.8) {
        vec4 n = vec4(-1, 0, 0, 0);
        p.velocity = reflect(p.velocity, n);
        p.velocity *= 0.99;
        p.location = vec4(0.8, p.location.y, 0, 0);
    }

    if (p.location.x < -0.8) {
        vec4 n = vec4(1, 0, 0, 0);
        p.velocity = reflect(p.velocity, n);
        p.velocity *= 0.99;
        p.location = vec4(-0.8, p.location.y, 0, 0);
    }

    next_particles[gid].location = p.location;
    next_particles[gid].velocity = p.velocity;

    //next_particles[gid].velocity = p.velocity + accel * dt;
    //next_particles[gid].location = p.location + p.velocity * dt + accel * dt * dt / 2.0;
}
