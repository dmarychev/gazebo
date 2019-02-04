// calculate forces
#version 460

layout(local_size_x = 16, local_size_y = 1, local_size_z = 1) in;

struct Particle {
    vec2 r;
    vec2 v;
    vec2 f;
    vec2 prev_f;
    float p; // pressure
    float d; // density
    float m; // mass
    float _;
};

const float PI = 3.1415926535897932384626433832795;

uniform float h = 0.01;
uniform float g = 0.08;
uniform float mu = 5.0;

layout(std430, binding=0) buffer Particles {
    Particle current_particles[];
};

layout(std430, binding=1) buffer Index {
    uint index[];
};

uniform uint index_max_neighbors = 40;

vec2 gradW(vec2 r, float h) {
    float nr = length(r);
    return (nr < h && nr > 0) ? (-45.0f / (PI * pow(h, 6))) * pow(h - nr, 2) * normalize(r): vec2(0, 0);
}

float laplacianW(float r, float h) {
    return (45.f / (PI * pow(h, 6))) * (h - r);
}

void main()
{
    uint p_i = gl_GlobalInvocationID.x;
    Particle p = current_particles[p_i];

    uint index_base = p_i * index_max_neighbors;

    vec2 f_press = vec2(.0f, .0f);
    vec2 f_vis = vec2(.0f, .0f);
    for (uint i = 0; i < index_max_neighbors; i++) {
        uint neighbor_idx = index[index_base + i];
        if (neighbor_idx != 0xdeadbeef && neighbor_idx != p_i) {
            Particle o = current_particles[neighbor_idx];
            f_press += -(o.m / o.d) * 0.5 * (o.p + p.p) * gradW(p.r - o.r, h);
            f_vis += (o.m / o.d) * (o.v - p.v) * laplacianW(length(p.r - o.r), h);
        }
    }

    vec2 f_grav = vec2(0, -p.d * g);

    p.f = f_press + f_grav + mu * f_vis;

    current_particles[p_i] = p;
}
