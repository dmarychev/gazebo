// cleanup index
#version 460

layout(local_size_x = 16, local_size_y = 1, local_size_z = 1) in;

layout(std430, binding=1) buffer Index {
    uint index[];
};

uniform uint index_max_neighbors = 40;

void main()
{
    uint index_base = gl_GlobalInvocationID.x * index_max_neighbors;
    for (uint i = 0; i < index_max_neighbors; i++) {
        index[index_base + i] = 0xdeadbeef;
    }
}
