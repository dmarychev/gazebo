package particles

import "github.com/go-gl/gl/v4.6-core/gl"
import "github.com/dmarychev/gazebo/core"
import "unsafe"

//import "log"

const (
	ATTRIB_COORDINATES = iota // index of coordinates attribute buffer
)

func AttachVertexAttributes() func() {
	gl.EnableVertexAttribArray(ATTRIB_COORDINATES)
	gl.VertexAttribPointer(ATTRIB_COORDINATES, 2, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), unsafe.Pointer(unsafe.Offsetof(Particle{}.R)))
	return func() {
		gl.DisableVertexAttribArray(ATTRIB_COORDINATES)
	}
}

type RenderState struct {
	updateTechniques []*core.Technique       // a techniques used to update system
	renderTechnique  *core.Technique         // a technique used to render system
	vao              core.VertexArrayObject  // array buffer associated with the state
	vbo              core.VertexBufferObject // a VBO containing particles' state.
	countParticles   uint32                  // number of particles in process
}

func NewRenderState(render *core.Technique) *RenderState {
	rs := RenderState{
		updateTechniques: make([]*core.Technique, 0, 10),
		renderTechnique:  render,
	}
	rs.vbo = core.MakeVertexBufferObject(0, nil)
	return &rs
}

func (rs *RenderState) AddUpdateTechnique(t *core.Technique) {
	rs.updateTechniques = append(rs.updateTechniques, t)
}

func (rs *RenderState) SetParticles(particles []Particle) {
	if len(particles) > 0 {
		sizeBytes := uint32(len(particles)) * uint32(unsafe.Sizeof(Particle{}))
		rs.vbo.SetData(gl.Ptr(particles), sizeBytes)
		rs.countParticles = uint32(len(particles))
	}
}

func (rs *RenderState) Update() {

	for _, technique := range rs.updateTechniques {
		disable := technique.Enable()
		defer disable()

		unbind := rs.vbo.BindBase(gl.SHADER_STORAGE_BUFFER, 0)
		defer unbind()

		gl.DispatchCompute(rs.countParticles, 1, 1)
		core.CheckError()

		gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
		core.CheckError()
	}
}

func (rs *RenderState) Render() {

	unbind := rs.vao.Bind()
	defer unbind()

	unbind = rs.vbo.Bind(gl.ARRAY_BUFFER)
	defer unbind()

	detach := AttachVertexAttributes()
	defer detach()

	disable := rs.renderTechnique.Enable()
	defer disable()

	gl.DrawArrays(gl.POINTS, 0, int32(rs.countParticles))
}
