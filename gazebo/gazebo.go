package main

import "log"
import "time"

import "math/rand"
import "runtime"
import "github.com/go-gl/glfw/v3.2/glfw"
import "github.com/go-gl/gl/v4.6-core/gl"
import "github.com/dmarychev/gazebo/particles"
import "github.com/dmarychev/gazebo/core"

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.PROGRAM_POINT_SIZE)
}

func initGlfw() *glfw.Window {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(1920, 1080, "Test", nil, nil)
	if err != nil {
		panic(err)
	}

	//window.SetAspectRatio(4, 3)
	window.MakeContextCurrent()

	return window
}

// initialize simple SPH particles system
func simpleSPHSystem() *particles.System {

	smoothingRadius := float32(0.01)
	maxNeighborParticles := uint32(40)
	viscosity := float32(5.0)
	gravity := float32(0.08)
	//gravity := float32(0.0)
	pressureCoefficient := float32(0.015)
	modellingTimeStep := float32(0.01)
	damping := float32(-0.99)

	renderThis, err := particles.NewRenderTechniqueFromFile("vfx/test.vs", "vfx/test.fs")
	if err != nil {
		panic(err)
	}

	indexClear, err := particles.NewComputeTechniqueFromFile("sph/index_clear.cs")
	if err != nil {
		panic(err)
	}

	indexUpdate, err := particles.NewComputeTechniqueFromFile("sph/index_update.cs")
	if err != nil {
		panic(err)
	}
	indexUpdate.SetUniformFloat32("h", smoothingRadius)
	indexUpdate.SetUniformUint("index_max_neighbors", maxNeighborParticles)

	densityAndPressure, err := particles.NewComputeTechniqueFromFile("sph/density_and_pressure.cs")
	if err != nil {
		panic(err)
	}
	densityAndPressure.SetUniformFloat32("h", smoothingRadius)
	densityAndPressure.SetUniformFloat32("k", pressureCoefficient)
	densityAndPressure.SetUniformUint("index_max_neighbors", maxNeighborParticles)

	accumulateForces, err := particles.NewComputeTechniqueFromFile("sph/accumulate_forces.cs")
	if err != nil {
		panic(err)
	}
	accumulateForces.SetUniformFloat32("h", smoothingRadius)
	accumulateForces.SetUniformFloat32("mu", viscosity)
	accumulateForces.SetUniformFloat32("g", gravity)
	accumulateForces.SetUniformUint("index_max_neighbors", maxNeighborParticles)

	leapfrog, err := particles.NewComputeTechniqueFromFile("sph/leapfrog_integration.cs")
	if err != nil {
		panic(err)
	}
	leapfrog.SetUniformFloat32("dt", modellingTimeStep)

	reflectBoundaries, err := particles.NewComputeTechniqueFromFile("sph/reflect_boundaries.cs")
	if err != nil {
		panic(err)
	}
	reflectBoundaries.SetUniformFloat32("damping_coeff", damping)

	ps := particles.NewSystem(renderThis, indexUpdate, indexClear, maxNeighborParticles)

	ps.AddUpdateTechnique(densityAndPressure)
	ps.AddUpdateTechnique(accumulateForces)
	ps.AddUpdateTechnique(leapfrog)
	ps.AddUpdateTechnique(reflectBoundaries)

	return ps
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	rand.Seed(time.Now().UnixNano())

	initOpenGL()

	particlesSet := make([]particles.Particle, 0, 4096)

	for i := 0; i < 16; i++ {
		for j := 0; j < 256; j++ {
			particlesSet = append(particlesSet, particles.Particle{
				R: core.Vec2{X: -0.8 + 0.01*float32(i), Y: -0.8 + 0.01*float32(j)},
				//V: core.Vec2{Y: -2},
				M: 0.01,
			})
		}
	}

	// drop
	/*	for i := 0; i < 16; i++ {
			for j := 0; j < 16; j++ {
				particlesSet = append(particlesSet, particles.Particle{
					R: core.Vec2{X: 0.01 * float32(i), Y: 20.0 + 0.01*float32(j)},
					//V: core.Vec2{Y: -2},
					M: 0.01,
				})
			}
		}
	*/
	for i := range particlesSet {
		particlesSet[i].R.X += -0.0005 + 0.001*rand.Float32()
		particlesSet[i].R.Y += -0.0005 + 0.001*rand.Float32()
	}

	ps := simpleSPHSystem()
	ps.SetParticles(particlesSet)

	log.Printf("Press Enter to toggle simulation")

	simulationOn := false
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEnter && action == glfw.Press {
			simulationOn = !simulationOn
			if simulationOn {
				log.Printf("Simulation Off")
			} else {
				log.Printf("Simulation On")
			}
		}
	})

	t0 := time.Now()
	fps := 0
	gl.ClearColor(0.8, 0.8, 0.8, 1.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		if simulationOn {
			ps.Update()
		}
		ps.Render()
		fps++
		window.SwapBuffers()
		glfw.PollEvents()
		if t1 := time.Now(); t1.Sub(t0) >= 1E9 {
			log.Printf("%v FPS", fps)
			fps, t0 = 0, t1
		}
		core.CheckError()
	}
}
