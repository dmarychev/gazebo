package main

import "runtime"
import "github.com/go-gl/glfw/v3.2/glfw"
import "github.com/go-gl/gl/v4.6-core/gl"
import "github.com/dmarychev/gazebo/particles"

import "log"
import "time"
import "math"

import "math/rand"

var vertexShader particles.VertexShaderSource = `
	#version 460

	in vec3 p_location;
	in vec3 p_velocity;
	in float p_size;
	in float p_weight;
	in float p_time;

	out OutputBlock {
		vec3 o_location;
		vec3 o_velocity;
		float o_size;
		float o_weight;
		float o_time;
	};

	void main() {
		float dt = 0.001;
		o_time = p_time;
		o_weight = p_weight * 1;
		o_size = p_size * 1;

		if (p_location.y < 0) {
			o_location = vec3(0, 0, 0);
			o_velocity = -1.0 * p_velocity;
		} else {
			vec3 accel = vec3(0, -0.5, 0);

			o_velocity = p_velocity + accel * dt;
			o_location = p_location + p_velocity * dt + accel * dt * dt / 2.0;
		}
		gl_Position = vec4(o_location.xyz, 1);
		gl_PointSize = 4.0;
	}
`

var fragmentShader particles.FragmentShaderSource = `
	#version 460
	out vec4 frag_color;

	void main() {
		frag_color = vec4(0, 1, 0, 1.0);
	}
`

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

	varyings := []string{"o_location", "o_velocity", "o_size", "o_weight", "o_time"}
	update, err := particles.NewTechnique(&vertexShader, nil, &varyings)
	if err != nil {
		panic(err)
	}

	render, err := particles.NewTechnique(nil, &fragmentShader, nil)
	if err != nil {
		panic(err)
	}

	particlesSet := make([]particles.Particle, 0, 10000)
	for i := 0; i < cap(particlesSet); i++ {
		angle := float64(math.Pi * rand.Float32())
		particlesSet = append(particlesSet, particles.Particle{
			Vx: 0.5 * (1.0 - 2.0*float32(math.Cos(angle))),
			Vy: 0.5 + 0.5*float32(math.Sin(angle)),
			//Size:   0.1 * (2.0 * (rand.Float32() - 0.5)),
			//Weight: 0.1 * (2.0 * (rand.Float32() - 0.5)),
			//Time:   100.0 * rand.Float32(),
		})
	}

	ps := particles.NewSystem(particlesSet, 0.5, update, render)

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
