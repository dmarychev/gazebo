package core

import "unsafe"
import "github.com/go-gl/gl/v4.6-core/gl"

type TransformFeedbackObject uint32
type VertexArrayObject uint32
type VertexBufferObject uint32

// Methods of Transform Feedback Object

func MakeTransformFeedbackObject() TransformFeedbackObject {
	var tfb uint32
	gl.GenTransformFeedbacks(1, &tfb)
	return TransformFeedbackObject(tfb)
}

func (tfb TransformFeedbackObject) Bind() func() {
	gl.BindTransformFeedback(gl.TRANSFORM_FEEDBACK, uint32(tfb))
	return func() {
		gl.BindTransformFeedback(gl.TRANSFORM_FEEDBACK, 0)
	}
}

func (tfb TransformFeedbackObject) AttachVertexBuffer(vbo VertexBufferObject) func() {
	unbind_tfb := tfb.Bind()
	unbind_vbo := vbo.BindBase(gl.TRANSFORM_FEEDBACK_BUFFER, 0)
	CheckError()
	return func() {
		unbind_vbo()
		unbind_tfb()
	}
}

// Methods of Vertex Array Object

func MakeVertexArrayObject() VertexArrayObject {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	return VertexArrayObject(vao)
}

func (vao VertexArrayObject) Bind() func() {
	gl.BindVertexArray(uint32(vao))
	return func() {
		gl.BindVertexArray(0)
	}
}

// Methods of Vertex Buffer Object

func MakeVertexBufferObject(sizeBytes int, data unsafe.Pointer) VertexBufferObject {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	if sizeBytes > 0 {
		gl.NamedBufferData(vbo, sizeBytes, data, gl.DYNAMIC_DRAW)
	}
	return VertexBufferObject(vbo)
}

func (vbo VertexBufferObject) SetData(data unsafe.Pointer, size uint32) uint32 {
	// TODO: consider keeping buffer if it's size is enough
	CheckError()

	unbind := vbo.Bind(gl.SHADER_STORAGE_BUFFER)
	defer unbind()

	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(size), data, gl.DYNAMIC_DRAW)

	CheckError()
	return size
}

func (vbo VertexBufferObject) GetData(data unsafe.Pointer, size uint32) {
	CheckError()

	unbind := vbo.Bind(gl.SHADER_STORAGE_BUFFER)
	defer unbind()

	gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, int(size), data)

	CheckError()
}

func (vbo VertexBufferObject) Bind(target uint32) func() {
	gl.BindBuffer(target, uint32(vbo))
	return func() {
		gl.BindBuffer(target, 0)
	}
}

func (vbo VertexBufferObject) BindBase(target, index uint32) func() {
	gl.BindBufferBase(target, index, uint32(vbo))
	return func() {
		gl.BindBufferBase(target, index, 0)
	}
}

func (vbo VertexBufferObject) Size() (sizeBytes int32) {
	unbind := vbo.Bind(gl.ARRAY_BUFFER)
	defer unbind()

	gl.GetBufferParameteriv(gl.ARRAY_BUFFER, gl.BUFFER_SIZE, &sizeBytes)
	CheckError()
	return
}
