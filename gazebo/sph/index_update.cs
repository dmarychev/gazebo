// update index
#version 460

layout(local_size_x = 16, local_size_y = 16, local_size_z = 1) in;

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

layout(std430, binding=0) buffer Particles {
    Particle current_particles[];
};

layout(std430, binding=1) buffer Index {
    uint index[];
};

uniform uint index_max_neighbors = 40;
uniform float h = 0.01;

void main()
{
    uint p_i = gl_GlobalInvocationID.x;
    uint candidate_i = gl_GlobalInvocationID.y;

    Particle p = current_particles[p_i];
    Particle candidate = current_particles[candidate_i];

    vec2 d = p.r - candidate.r;
    if (length(d) < h) {
        uint index_base = p_i * index_max_neighbors;
        for (uint i = 0; i < index_max_neighbors; i++) {
            if (atomicCompSwap(index[index_base + i], 0xdeadbeef, candidate_i) == 0xdeadbeef) {
                break;
            }
        }
    }
}
