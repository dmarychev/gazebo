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
	updateTechniques  []*core.Technique       // a techniques used to update system
	renderTechnique   *core.Technique         // a technique used to render system
	indexUpdate       *core.Technique         // a technique used to update index system
	indexClear        *core.Technique         // a technique used to clear index system
	indexMaxNeighbors uint32                  // maximum neighbors in index
	vao               core.VertexArrayObject  // array buffer associated with the state
	vbo               core.VertexBufferObject // a VBO containing particles' state.
	indexVbo          core.VertexBufferObject // a VBO containing index data
	countParticles    uint32                  // number of particles in process
}

func NewRenderState(render, indexUpdate, indexClear *core.Technique, indexMaxNeighbors uint32) *RenderState {
	rs := RenderState{
		updateTechniques:  make([]*core.Technique, 0, 10),
		renderTechnique:   render,
		indexUpdate:       indexUpdate,
		indexClear:        indexClear,
		indexMaxNeighbors: indexMaxNeighbors,
	}
	rs.vbo = core.MakeVertexBufferObject(0, nil)
	rs.indexVbo = core.MakeVertexBufferObject(0, nil)
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
		rs.indexVbo.SetData(nil, rs.countParticles*rs.indexMaxNeighbors*uint32(unsafe.Sizeof(uint32(0))))
	}
}

func (rs *RenderState) Update() {
	/*
		p_data := make([]float32, rs.countParticles*uint32(unsafe.Sizeof(Particle{}))/uint32(unsafe.Sizeof(float32(0))))
		rs.vbo.GetData(gl.Ptr(p_data), rs.countParticles*uint32(unsafe.Sizeof(Particle{})))
		log.Printf("Updated Particles: ")
		log.Printf("%v", p_data)
		log.Printf("End of particles")
	*/
	unbindParticles := rs.vbo.BindBase(gl.SHADER_STORAGE_BUFFER, 0)
	unbindIndex := rs.indexVbo.BindBase(gl.SHADER_STORAGE_BUFFER, 1)

	if rs.indexClear != nil {
		disable := rs.indexClear.Enable()

		gl.DispatchCompute(rs.countParticles/16, 1, 1)
		core.CheckError()

		gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
		disable()
		/*
			data := make([]uint32, rs.countParticles*rs.indexMaxNeighbors, rs.countParticles*rs.indexMaxNeighbors)
			rs.indexVbo.GetData(gl.Ptr(data), uint32(len(data))*uint32(unsafe.Sizeof(uint32(0))))
			log.Printf("Cleared Index: ")
			log.Printf("%v", data)
			log.Printf("End of index")*/
	}

	if rs.indexUpdate != nil {
		disableUpdate := rs.indexUpdate.Enable()

		gl.DispatchCompute(rs.countParticles/16, rs.countParticles/16, 1)
		core.CheckError()

		gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
		disableUpdate()

	}

	for _, technique := range rs.updateTechniques {
		disable := technique.Enable()

		gl.DispatchCompute(rs.countParticles/16, 1, 1)
		core.CheckError()

		gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
		core.CheckError()
		/*
			p_data := make([]float32, rs.countParticles*uint32(unsafe.Sizeof(Particle{}))/uint32(unsafe.Sizeof(float32(0))))
			rs.vbo.GetData(gl.Ptr(p_data), rs.countParticles*uint32(unsafe.Sizeof(Particle{})))
			log.Printf("Updated Particles #: ")
			log.Printf("%v", p_data)
			log.Printf("End of particles #")*/

		disable()
	}

	unbindIndex()
	unbindParticles()

	/*i_data := make([]uint32, rs.countParticles*rs.indexMaxNeighbors, rs.countParticles*rs.indexMaxNeighbors)
	rs.indexVbo.GetData(gl.Ptr(i_data), uint32(len(i_data))*uint32(unsafe.Sizeof(uint32(0))))
	log.Printf("Updated Index: ")
	log.Printf("%v", i_data)
	log.Printf("End of index")*/
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
