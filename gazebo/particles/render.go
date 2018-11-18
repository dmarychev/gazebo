package particles

import "github.com/go-gl/gl/v4.6-core/gl"
import "unsafe"

//import "log"

const (
	ATTRIB_COORDINATES = iota // index of coordinates attribute buffer
	ATTRIB_VELOCITIES  = iota // index of velocities attribute buffer
	ATTRIB_SIZES       = iota // index of sizes attribute buffer
	ATTRIB_WEIGHTS     = iota // index of weights attribute buffer
	ATTRIB_TIMES       = iota // index of times attribute buffer
)

func AttachVertexAttributes() func() {
	gl.EnableVertexAttribArray(ATTRIB_COORDINATES)
	gl.EnableVertexAttribArray(ATTRIB_VELOCITIES)
	gl.EnableVertexAttribArray(ATTRIB_SIZES)
	gl.EnableVertexAttribArray(ATTRIB_WEIGHTS)
	gl.EnableVertexAttribArray(ATTRIB_TIMES)

	gl.VertexAttribPointer(ATTRIB_COORDINATES, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), unsafe.Pointer(unsafe.Offsetof(Particle{}.X)))
	gl.VertexAttribPointer(ATTRIB_VELOCITIES, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), unsafe.Pointer(unsafe.Offsetof(Particle{}.Vx)))
	gl.VertexAttribPointer(ATTRIB_SIZES, 1, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), unsafe.Pointer(unsafe.Offsetof(Particle{}.Size)))
	gl.VertexAttribPointer(ATTRIB_WEIGHTS, 1, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), unsafe.Pointer(unsafe.Offsetof(Particle{}.Weight)))
	gl.VertexAttribPointer(ATTRIB_TIMES, 1, gl.FLOAT, false, int32(unsafe.Sizeof(Particle{})), unsafe.Pointer(unsafe.Offsetof(Particle{}.Time)))

	return func() {
		gl.DisableVertexAttribArray(ATTRIB_COORDINATES)
		gl.DisableVertexAttribArray(ATTRIB_VELOCITIES)
		gl.DisableVertexAttribArray(ATTRIB_SIZES)
		gl.DisableVertexAttribArray(ATTRIB_WEIGHTS)
		gl.DisableVertexAttribArray(ATTRIB_TIMES)
	}
}

type RenderState struct {
	updateTechnique  *Technique              // a technique used to update system
	renderTechnique  *Technique              // a technique used to render system
	vao              VertexArrayObject       // array buffer associated with the state
	vbo              [2]VertexBufferObject   // one VBO contains particles to be drawn now, another is a receiver for transformed particles. VBOs are swapped after render.
	tfb              TransformFeedbackObject // transform feedback object, used to intercept data after vertex processing stage and moving it to destination VBO
	indexSource      uint32                  // VBO with this index is being read
	indexDestination uint32                  // VBO with this index is being written
	countParticles   uint32                  // number of particles in process
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

	rs.tfb = MakeTransformFeedbackObject()

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

func (rs *RenderState) Update() {
	gl.Enable(gl.RASTERIZER_DISCARD)

	unbind := rs.vao.Bind()
	defer unbind()

	detach := rs.tfb.AttachVertexBuffer(rs.DestinationVbo())
	defer detach()

	disable := rs.updateTechnique.Enable()
	defer disable()

	unbind = rs.SourceVbo().Bind(gl.ARRAY_BUFFER)
	defer unbind()

	detach = AttachVertexAttributes()
	defer detach()

	//var query uint32
	//gl.GenQueries(1, &query)
	//gl.BeginQuery(gl.TRANSFORM_FEEDBACK_PRIMITIVES_WRITTEN, query)

	gl.BeginTransformFeedback(gl.POINTS)
	gl.DrawArrays(gl.POINTS, 0, int32(rs.countParticles))
	gl.EndTransformFeedback()

	//gl.EndQuery(gl.TRANSFORM_FEEDBACK_PRIMITIVES_WRITTEN)
	//var primWritten uint32
	//gl.GetQueryObjectuiv(query, gl.QUERY_RESULT, &primWritten)
	//	log.Printf("Primitives written: %v\n", primWritten)

	gl.Disable(gl.RASTERIZER_DISCARD)

}

func (rs *RenderState) SetParticles(particles []Particle) {
	if len(particles) > 0 {
		sizeBytes := uint32(len(particles)) * uint32(unsafe.Sizeof(Particle{}))
		rs.SourceVbo().SetData(unsafe.Pointer(&particles[0]), sizeBytes)
		rs.DestinationVbo().SetData(nil, sizeBytes)
		rs.countParticles = uint32(len(particles))
	}
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

	unbind = rs.DestinationVbo().Bind(gl.ARRAY_BUFFER)
	defer unbind()

	disable := rs.renderTechnique.Enable()
	defer disable()

	gl.DrawArrays(gl.POINTS, 0, int32(rs.countParticles))
}
