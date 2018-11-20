package particles

// particle struct
type Particle struct {
	X, Y, Z, _    float32 // coordinates
	Vx, Vy, Vz, _ float32 // velocities
	//	Weight     float32 // weight
	//	Time       float32 // particle's current time, seconds from spawn
}

type System struct {
	particles     []Particle   // particles
	renderState   *RenderState // objects related to rendering
	modellingStep float32      // increase for particle's time
}

func NewSystem(particles []Particle, modellingStep float32, updateTechnique, renderTechnique *Technique) *System {

	s := System{
		modellingStep: modellingStep,
		particles:     make([]Particle, len(particles)),
		renderState:   NewRenderState(updateTechnique, renderTechnique),
	}

	// work with copy of particles for safety
	copy(s.particles, particles)

	s.renderState.SetParticles(s.particles)

	return &s
}

func (s *System) Size() uint32 {
	return uint32(len(s.particles))
}

// shows current state on screen
func (s *System) Render() {
	s.renderState.Render()
	checkError()
}

// updates particle system's state
func (s *System) Update() {
	s.renderState.Update()
	checkError()
}

// synchronizes `s.particles` with current state in GPU memory
func (s *System) Sync() {

}
