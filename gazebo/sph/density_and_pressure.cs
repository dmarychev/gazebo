// calculate density and pressure
#version 460
#pragma optimize(off)

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
uniform float k = 0.01;

layout(std430, binding=0) buffer Particles {
    Particle current_particles[];
};

layout(std430, binding=1) buffer Index {
    uint index[];
};

uniform uint index_max_neighbors = 40;

void main()
{
    uint p_i = gl_GlobalInvocationID.x;
    Particle p = current_particles[p_i];

    uint index_base = p_i * index_max_neighbors;

    const float k_poly6_coeff = 315.0/(64.0 * PI * pow(h, 9));
    const float h2 = h * h;

    p.d = 0.0;
    for (uint i = 0; i < index_max_neighbors; i++) {
        uint neighbor_idx = index[index_base + i];
        if (neighbor_idx == 0xdeadbeef) {
            break;
        }
        Particle o = current_particles[neighbor_idx];

        vec2 dr = p.r - o.r;
        float dr2 = dot(dr, dr);
        float d_h2_dr2 = h2 - dr2;

        p.d += o.m * k_poly6_coeff * d_h2_dr2 * d_h2_dr2 * d_h2_dr2;
    }

    p.p = k * p.d;
    current_particles[p_i] = p;
}
