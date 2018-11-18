package particles

import "fmt"
import "strings"
import "github.com/go-gl/gl/v4.6-core/gl"

type FragmentShaderSource string
type VertexShaderSource string

// compile shader from GLSL text
func CompileShader(source interface{}) (uint32, error) {

	var shaderType uint32
	var shaderText string
	switch shader := source.(type) {
	case VertexShaderSource:
		shaderType = gl.VERTEX_SHADER
		shaderText = string(shader)
	case FragmentShaderSource:
		shaderType = gl.FRAGMENT_SHADER
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

type Technique struct {
	program uint32
}

func NewTechnique(vertexShader *VertexShaderSource, fragmentShader *FragmentShaderSource, varyings *[]string) (*Technique, error) {

	t := new(Technique)
	t.program = gl.CreateProgram()

	if vertexShader != nil {
		vshader, err := CompileShader(*vertexShader)
		if err != nil {
			return nil, fmt.Errorf("Failed to compile vertex shader: %v", err)
		}
		gl.AttachShader(t.program, vshader)
	}

	if fragmentShader != nil {
		fshader, err := CompileShader(*fragmentShader)
		if err != nil {
			return nil, fmt.Errorf("Failed to compile fragment shader: %v", err)
		}
		gl.AttachShader(t.program, fshader)
	}

	if varyings != nil && len(*varyings) > 0 {
		t.attachVaryings(*varyings)
	}

	gl.LinkProgram(t.program)
	checkError()
	return t, nil
}

func (t *Technique) attachVaryings(variables []string) {

	fmtVars := make([]string, len(variables))
	for i, v := range variables {
		fmtVars[i] = fmt.Sprintf("%v\x00", v)
	}

	fvaryings, free := gl.Strs(fmtVars...)
	defer free()

	gl.TransformFeedbackVaryings(t.program, int32(len(variables)), fvaryings, gl.INTERLEAVED_ATTRIBS)
	checkError()
}

func (t *Technique) SetUniform(name string, value interface{}) {
	// TODO: implement me
}

func (t *Technique) Enable() func() {
	gl.UseProgram(t.program)
	checkError()
	return func() {
		gl.UseProgram(0)
		checkError()
	}
}
