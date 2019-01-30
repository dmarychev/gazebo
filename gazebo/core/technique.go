package core

import "fmt"
import "strings"
import "github.com/go-gl/gl/v4.6-core/gl"

type FragmentShaderSource string
type VertexShaderSource string
type ComputeShaderSource string

// compile shader from GLSL text
func compileShader(source interface{}) (uint32, error) {

	var shaderType uint32
	var shaderText string
	switch shader := source.(type) {
	case VertexShaderSource:
		shaderType = gl.VERTEX_SHADER
		shaderText = string(shader)
	case FragmentShaderSource:
		shaderType = gl.FRAGMENT_SHADER
		shaderText = string(shader)
	case ComputeShaderSource:
		shaderType = gl.COMPUTE_SHADER
		shaderText = string(shader)
	default:
		return 0, fmt.Errorf("Shader type is not supported %T", source)
	}

	shader := gl.CreateShader(shaderType)

	csource, free := gl.Strs(shaderText)
	defer free()
	gl.ShaderSource(shader, 1, csource, nil)

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

type Technique uint32

func NewRenderTechnique(vertexShader *VertexShaderSource, fragmentShader *FragmentShaderSource) (*Technique, error) {
	return newTechnique(vertexShader, fragmentShader, nil)
}

func NewComputeTechnique(computeShader *ComputeShaderSource) (*Technique, error) {
	return newTechnique(nil, nil, computeShader)
}

func newTechnique(vertexShader *VertexShaderSource, fragmentShader *FragmentShaderSource, computeShader *ComputeShaderSource) (*Technique, error) {

	t := Technique(gl.CreateProgram())

	if vertexShader != nil {
		vshader, err := compileShader(*vertexShader)
		if err != nil {
			return nil, fmt.Errorf("Failed to compile vertex shader: %v", err)
		}
		gl.AttachShader(uint32(t), vshader)
	}

	if fragmentShader != nil {
		fshader, err := compileShader(*fragmentShader)
		if err != nil {
			return nil, fmt.Errorf("Failed to compile fragment shader: %v", err)
		}
		gl.AttachShader(uint32(t), fshader)
	}

	if computeShader != nil {
		cshader, err := compileShader(*computeShader)
		if err != nil {
			return nil, fmt.Errorf("Failed to compile compute shader: %v", err)
		}
		gl.AttachShader(uint32(t), cshader)
	}

	if err := t.linkAndValidate(); err != nil {
		return nil, err
	}
	CheckError()

	return &t, nil
}

func (t *Technique) linkAndValidate() error {
	p := uint32(*t)

	var status int32
	var logLength int32

	gl.LinkProgram(p)
	CheckError()
	gl.GetProgramiv(p, gl.LINK_STATUS, &status)
	gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)
	if logLength > 0 {
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p, logLength, nil, gl.Str(log))
		fmt.Printf("Link log:\n%v\n", log)
	}
	if status == gl.FALSE {
		return fmt.Errorf("failed to link program: %v", status)
	}

	gl.ValidateProgram(p)
	CheckError()
	gl.GetProgramiv(p, gl.VALIDATE_STATUS, &status)
	gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)
	if logLength > 0 {
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p, logLength, nil, gl.Str(log))
		fmt.Printf("Validate log:\n%v\n", log)
	}
	if status == gl.FALSE {
		return fmt.Errorf("failed to validate program: %v", status)
	}

	return nil
}

func (t *Technique) GetUniformFloat32(name string) (value float32) {
	uLocation := gl.GetUniformLocation(uint32(*t), gl.Str(fmt.Sprintf("%v\x00", name)))
	gl.GetUniformfv(uint32(*t), uLocation, &value)
	return
}

func (t *Technique) SetUniformFloat32(name string, value float32) {
	uLocation := gl.GetUniformLocation(uint32(*t), gl.Str(fmt.Sprintf("%v\x00", name)))
	disable := t.Enable()
	defer disable()
	gl.Uniform1f(uLocation, value)
}

func (t *Technique) Enable() func() {
	gl.UseProgram(uint32(*t))
	CheckError()
	return func() {
		gl.UseProgram(0)
		CheckError()
	}
}
