package particles

import "log"
import "io/ioutil"
import "github.com/dmarychev/gazebo/core"
import "github.com/dmarychev/gazebo/inspect"

// particle struct
type Particle struct {
	R     core.Vec2 // coordinates
	V     core.Vec2 // velocity
	F     core.Vec2 // F(total)
	prevF core.Vec2 // F(total) on previous step
	P     float32   // pressure
	D     float32   // density
	M     float32   // mass
	_     float32
}

type System struct {
	renderState *RenderState // objects related to rendering
}

func NewSystem(renderTechnique, indexUpdate, indexClear *core.Technique, indexMaxNeighbors uint32) *System {

	s := System{
		renderState: NewRenderState(renderTechnique, indexUpdate, indexClear, indexMaxNeighbors),
	}

	return &s
}

func (s *System) SetParticles(particles []Particle) {
	s.renderState.SetParticles(particles)
}

func (s *System) AddUpdateTechniqueFromFile(compShaderFile string) (err error) {
	technique, err := NewComputeTechniqueFromFile(compShaderFile)
	if err != nil {
		return
	}
	s.AddUpdateTechnique(technique)
	return
}

func (s *System) AddUpdateTechnique(t *core.Technique) {
	s.renderState.AddUpdateTechnique(t)
}

// shows current state on screen
func (s *System) Render() {
	s.renderState.Render()
	core.CheckError()
}

// updates particle system's state
func (s *System) Update() {
	s.renderState.Update()
	core.CheckError()
}

// synchronizes `s.particles` with current state in GPU memory
func (s *System) Sync() {

}

func NewComputeTechniqueFromFile(compShaderFile string) (*core.Technique, error) {
	log.Printf("Load compute technique: %v\n", compShaderFile)

	text, err := ioutil.ReadFile(compShaderFile)
	if err != nil {
		return nil, err
	}
	shaderSource := core.ComputeShaderSource(text)

	technique, err := core.NewComputeTechnique(&shaderSource)
	if err != nil {
		return nil, err
	}

	if err = LogTechniqueInfo(technique); err != nil {
		return nil, err
	}

	return technique, err
}

func NewRenderTechniqueFromFile(vertexShaderFile string, fragmentShaderFile string) (*core.Technique, error) {
	log.Printf("Load render technique: vs=%v fs=%v\n", vertexShaderFile, fragmentShaderFile)

	vsText, err := ioutil.ReadFile(vertexShaderFile)
	if err != nil {
		return nil, err
	}

	fsText, err := ioutil.ReadFile(fragmentShaderFile)
	if err != nil {
		return nil, err
	}

	vertexShaderSource := core.VertexShaderSource(vsText)
	fragmentShaderSource := core.FragmentShaderSource(fsText)

	technique, err := core.NewRenderTechnique(&vertexShaderSource, &fragmentShaderSource)
	if err != nil {
		return nil, err
	}

	if err = LogTechniqueInfo(technique); err != nil {
		return nil, err
	}

	return technique, err
}

func LogTechniqueInfo(t *core.Technique) error {
	log.Printf("Begin technique info\n")
	tinfo, err := inspect.InspectTechnique(t)
	if err != nil {
		return err
	}
	log.Printf("Uniform variables: \n")
	for _, ui := range tinfo.UniformVariables {
		log.Printf(" - %v\n", ui)
	}
	log.Printf("SSBO: \n")
	for _, ssbi := range tinfo.ShaderStorageBuffers {
		log.Printf("%v\n", ssbi)
		for _, variable := range ssbi.Variables {
			log.Printf(" - %v\n", variable)
		}
	}
	log.Printf("End technique info\n")
	return nil
}
