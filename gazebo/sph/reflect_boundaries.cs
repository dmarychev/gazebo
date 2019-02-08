// reflect boundaries
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

uniform float damping_coeff = -0.5;

layout(std430, binding=0) buffer Particles {
    Particle current_particles[];
};

const float half_h_size = 0.8;
const float half_w_size = 0.8;
const float eps = 0.001;

void main()
{
    uint gid = gl_GlobalInvocationID.x;
    Particle p = current_particles[gid];

/*    if (p.r.y >= half_h_size) {
        p.v *= damping_coeff;
        p.r.y = half_h_size - eps;
    } else*/ if (p.r.y <= -half_h_size) {
        p.v *= damping_coeff;
        p.r.y = -half_h_size + eps;
    } else if (p.r.x >= half_w_size) {
        p.v *= damping_coeff;
        p.r.x = half_w_size - eps;
    } else if (p.r.x <= -half_w_size) {
        p.v *= damping_coeff;
        p.r.x = -half_w_size + eps;
    }

    current_particles[gid] = p;
}
