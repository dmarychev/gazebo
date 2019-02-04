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

uniform float g = 0.08; // gravity
uniform float mu = 5.0; // viscosity coefficient
uniform float h = 0.01; // smoothing parameter
uniform uint index_max_neighbors = 40; // maximum number of neighbors in the index

layout(std430, binding=0) buffer Particles {
    Particle current_particles[];
};

layout(std430, binding=1) buffer Index {
    uint index[];
};

void main()
{
    uint p_i = gl_GlobalInvocationID.x;
    Particle p = current_particles[p_i];

    uint index_base = p_i * index_max_neighbors;

    const float k_h6 = h * h * h * h * h * h;
    const float k_grad_coeff = -45.0f / (PI * k_h6);
    const float k_lap_coeff = 45.f / (PI * k_h6);

    vec2 f_press = vec2(.0f, .0f);
    vec2 f_vis = vec2(.0f, .0f);
    for (uint i = 0; i < index_max_neighbors; i++) {
        uint neighbor_idx = index[index_base + i];
        if (neighbor_idx == 0xdeadbeef) {
            break;
        }
        if (neighbor_idx != p_i) {
            Particle o = current_particles[neighbor_idx];

            vec2 dr = p.r - o.r;
            float ldr = length(dr);
            vec2 ndr = ldr > 0 ? normalize(dr) : vec2(0, 0);

            // pressure force
            f_press += -(o.m / o.d) * 0.5 * (o.p + p.p) * k_grad_coeff * (h - ldr) * (h - ldr) * ndr;

            // viscosity force
            f_vis += (o.m / o.d) * (o.v - p.v) * k_lap_coeff * (h - length(dr));
        }
    }

    vec2 f_grav = vec2(0, -p.d * g);

    p.f = f_press + f_grav + mu * f_vis;

    current_particles[p_i] = p;
}
