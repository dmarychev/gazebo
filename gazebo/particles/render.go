package particles

import "github.com/go-gl/gl/v4.6-core/gl"
import "unsafe"

//import "log"

const (
	ATTRIB_COORDINATES = iota // index of coordinates attribute buffer
)

func AttachVertexAttributes() func() {
	gl.EnableVertexAttribArray(ATTRIB_COORDINATES)
	gl.VertexAttribPointer(ATTRIB_COORDINATES, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), unsafe.Pointer(unsafe.Offsetof(Particle{}.X)))
	return func() {
		gl.DisableVertexAttribArray(ATTRIB_COORDINATES)
	}
}

type RenderState struct {
	updateTechnique  *Technique            // a technique used to update system
	renderTechnique  *Technique            // a technique used to render system
	vao              VertexArrayObject     // array buffer associated with the state
	vbo              [2]VertexBufferObject // one VBO contains particles to be drawn now, another is a receiver for transformed particles. VBOs are swapped after render.
	indexSource      uint32                // VBO with this index is being read
	indexDestination uint32                // VBO with this index is being written
	countParticles   uint32                // number of particles in process
}

func NewRenderState(update *Technique, render *Technique) *RenderState {
	rs := RenderState{
		updateTechnique:  update,
		renderTechnique:  render,
		indexSource:      0,
		indexDestination: 1,
	}

	// initialize buffer; use single buffer, pointing to attributes with offset and stride
	for index := range rs.vbo {
		rs.vbo[index] = MakeVertexBufferObject(0, nil)
	}

	return &rs
}

func (rs *RenderState) SourceVbo() VertexBufferObject {
	return rs.vbo[rs.indexSource]
}

func (rs *RenderState) DestinationVbo() VertexBufferObject {
	return rs.vbo[rs.indexDestination]
}

func (rs *RenderState) Swap() {
	rs.indexSource, rs.indexDestination = rs.indexDestination, rs.indexSource
}

func (rs *RenderState) SetParticles(particles []Particle) {
	if len(particles) > 0 {
		sizeBytes := uint32(len(particles)) * uint32(unsafe.Sizeof(Particle{}))
		rs.SourceVbo().SetData(gl.Ptr(particles), sizeBytes)
		rs.DestinationVbo().SetData(nil, sizeBytes)
		rs.countParticles = uint32(len(particles))
	}
}

func (rs *RenderState) Update() {

	disable := rs.updateTechnique.Enable()
	defer disable()

	unbind := rs.SourceVbo().BindBase(gl.SHADER_STORAGE_BUFFER, 0)
	defer unbind()

	unbind = rs.DestinationVbo().BindBase(gl.SHADER_STORAGE_BUFFER, 1)
	defer unbind()

	gl.DispatchCompute(rs.countParticles, 1, 1)
	checkError()

	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
	checkError()
}

func (rs *RenderState) Render() {

	defer func() {
		rs.Swap()
	}()

	unbind := rs.vao.Bind()
	defer unbind()

	unbind = rs.DestinationVbo().Bind(gl.ARRAY_BUFFER)
	defer unbind()

	detach := AttachVertexAttributes()
	defer detach()

	disable := rs.renderTechnique.Enable()
	defer disable()

	gl.DrawArrays(gl.POINTS, 0, int32(rs.countParticles))
}
