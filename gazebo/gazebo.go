package main

import "log"
import "time"
import "math"
import "math/rand"
import "fmt"
import "io/ioutil"
import "runtime"
import "github.com/go-gl/glfw/v3.2/glfw"
import "github.com/go-gl/gl/v4.6-core/gl"
import "github.com/dmarychev/gazebo/particles"
import "github.com/dmarychev/gazebo/core"
import "github.com/dmarychev/gazebo/inspect"

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

	window, err := glfw.CreateWindow(800, 600, "Test", nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetAspectRatio(4, 3)
	window.MakeContextCurrent()

	return window
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	rand.Seed(time.Now().UnixNano())

	initOpenGL()

	vertexShader, err := ioutil.ReadFile("vfx/test.vs")
	if err != nil {
		panic(err)
	}

	fragmentShader, err := ioutil.ReadFile("vfx/test.fs")
	if err != nil {
		panic(err)
	}

	computeShader, err := ioutil.ReadFile("vfx/test.cs")
	if err != nil {
		panic(err)
	}

	vertexShaderSource := core.VertexShaderSource(vertexShader)
	fragmentShaderSource := core.FragmentShaderSource(fragmentShader)
	computeShaderSource := core.ComputeShaderSource(computeShader)

	updateTechnique, err := core.NewComputeTechnique(&computeShaderSource)
	if err != nil {
		panic(err)
	}

	renderTechnique, err := core.NewRenderTechnique(&vertexShaderSource, &fragmentShaderSource)
	if err != nil {
		panic(err)
	}

	tinfo, err := inspect.InspectTechnique(updateTechnique)
	if err != nil {
		panic(err)
	}
	for _, ssbi := range tinfo.ShaderStorageBuffers {
		fmt.Printf("%v\n", ssbi)
		for _, variable := range ssbi.Variables {
			fmt.Printf("%v\n", variable)
		}
	}

	particlesSet := make([]particles.Particle, 0, 100000)

	for i := 0; i < cap(particlesSet); i++ {
		rho := 0.05 + 0.2*rand.Float32()
		phi := float64(math.Pi * (0.25 + 0.5*rand.Float32()))
		particlesSet = append(particlesSet, particles.Particle{
			Vx: rho * float32(math.Cos(phi)),
			Vy: rho * float32(math.Sin(phi)),
		})
	}

	ps := particles.NewSystem(particlesSet, 0.5, updateTechnique, renderTechnique)

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
			log.Printf("%v FPS\r", fps)
			fps, t0 = 0, t1
		}
	}
}
