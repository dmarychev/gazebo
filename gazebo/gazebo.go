package main

import "log"
import "time"

//import "math"
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

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	rand.Seed(time.Now().UnixNano())

	initOpenGL()

	particlesSet := make([]particles.Particle, 0, 4096)

	for i := 0; i < 32; i++ {
		for j := 0; j < 256; j++ {
			particlesSet = append(particlesSet, particles.Particle{
				R: core.Vec2{X: 0.01 * float32(i), Y: 0.01 * float32(j)},
				//V: core.Vec2{Y: -2},
				M: 0.01,
			})
		}
	}

	for i := range particlesSet {
		particlesSet[i].R.X += -0.0005 + 0.001*rand.Float32()
		particlesSet[i].R.Y += -0.0005 + 0.001*rand.Float32()
	}

	/*	for i := 0; i < cap(particlesSet); i++ {
		rho := 0.5 + 0.5*rand.Float32()
		phi := float32(math.Pi * (0.25 + 0.5*rand.Float32()))
		particlesSet = append(particlesSet, particles.Particle{
			R: core.Vec2{X: rho * float32(math.Cos(phi)), Y: rho * float32(math.Sin(phi))},
			M: 0.1,
		})
	}*/

	renderThis, err := particles.NewRenderTechniqueFromFile("vfx/test.vs", "vfx/test.fs")
	if err != nil {
		panic(err)
	}

	indexClear, err := particles.NewComputeTechniqueFromFile("sph/index_clear.cs")
	if err != nil {
		panic(err)
	}

	indexMaxNeighbors := uint32(40)

	indexUpdate, err := particles.NewComputeTechniqueFromFile("sph/index_update.cs")
	if err != nil {
		panic(err)
	}
	indexUpdate.SetUniformUint("index_max_neighbors", indexMaxNeighbors)

	ps := particles.NewSystem(particlesSet, 0.5, renderThis, indexUpdate, indexClear, indexMaxNeighbors)

	// adding computation steps

	if err = ps.AddUpdateTechniqueFromFile("sph/density_and_pressure.cs"); err != nil {
		panic(err)
	}

	if err = ps.AddUpdateTechniqueFromFile("sph/accumulate_forces.cs"); err != nil {
		panic(err)
	}
	leapfrog, err := particles.NewComputeTechniqueFromFile("sph/leapfrog_integration.cs")
	if err != nil {
		panic(err)
	}
	leapfrog.SetUniformFloat32("dt", 0.01)
	log.Printf("dt=", leapfrog.GetUniformFloat32("dt"))
	ps.AddUpdateTechnique(leapfrog)

	if err = ps.AddUpdateTechniqueFromFile("sph/reflect_boundaries.cs"); err != nil {
		panic(err)
	}

	// calculations

	t0 := time.Now()
	fps := 0
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		ps.Update()
		ps.Render()
		fps++
		window.SwapBuffers()
		glfw.PollEvents()
		if t1 := time.Now(); t1.Sub(t0) >= 1000000000 {
			log.Printf("%v FPS", fps)
			fps, t0 = 0, t1
		}
		core.CheckError()
	}
}
