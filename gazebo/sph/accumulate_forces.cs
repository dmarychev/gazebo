// calculate forces
#version 460

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

uniform float g = 2.0;
uniform float h = 0.05;
uniform float mu = 5.0;

layout(std430, binding=0) buffer Particles {
    Particle current_particles[];
};

vec2 gradW(vec2 r, float h) {
    float nr = length(r);
    return nr < h ? -45.0f / (PI * pow(h, 6)) * pow(h - nr, 2) * normalize(r) : vec2(.0f, .0f);
}

float laplacianW(float r, float h) {
    return r < h ? 45.f / (PI * pow(h, 6)) * (h - r) : .0f;
}

void main()
{
    uint gid = gl_GlobalInvocationID.x;
    Particle p = current_particles[gid];

    uint num_particles = (gl_NumWorkGroups * gl_WorkGroupSize).x;

    vec2 f_press = vec2(.0f, .0f);
    vec2 f_vis = vec2(.0f, .0f);
    for (uint p_i = 0; p_i < num_particles; ++p_i) {
        if (p_i != gid) {
            Particle o = current_particles[p_i];
            f_press += -(o.m / o.d) * 0.5 * (o.p + p.p) * gradW(p.r - o.r, h);
            f_vis += (o.m / o.d) * (o.v - p.v) * laplacianW(length(p.r - o.r), h);
        }
    }

    vec2 f_grav = vec2(0, -p.d * g);

    p.f = vec2(0, 0) + f_press + f_grav + mu * f_vis;

    current_particles[gid] = p;
}
