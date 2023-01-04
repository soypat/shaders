package main

import (
	"errors"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/soypat/shaders"
	"golang.org/x/exp/slog"
)

var (
	ErrStringNotNullTerminated = errors.New("string not null terminated")
)

type Renderer struct {
}

// VertexArray ties data layout with vertex buffer(s).
// Is aware of data layout via VertexAttribPointer* calls.
type VertexArray struct {
	rid uint32
}

// AttribLayout is a low level configuration struct
// for adding vertex buffers attribute layouts to a vertex array object.
type AttribLayout struct {
	// The OpenGL program identifier.
	Program Program
	// Type is a OpenGL enum representing the underlying type. Valid types include
	// gl.FLOAT, gl.UNSIGNED_INT, gl.UNSIGNED_BYTE, gl.BYTE etc.
	Type uint32
	// Name is the identifier of the attribute in the
	// vertex shader source code finished with a null terminator.
	Name string
	// Packing is a value between 1 and 4 and represents how many
	// of the type are present at the attribute location.
	//
	// Example:
	// When w orking with a vec3 attribute in the shader source code
	// with a gl.Float type, then the Packing is 3 since there are
	// 3 floats packed at each attribute location.
	Packing int
	// Stride is the distance in bytes between attributes in the buffer.
	Stride int
	// Offset is the starting offset with which to start
	// traversing the vertex buffer.
	Offset int
	// specifies whether fixed-point data values should be normalized (when true)
	// or converted directly as fixed-point values (when false) when they are accessed.
	// Usually left as false?
	Normalize bool
}

func NewVAO() VertexArray {
	// Configure the Vertex Array Object.
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	return VertexArray{rid: vao}
}

func (vao VertexArray) AddAttribute(vbo VertexBuffer, layout AttribLayout) error {
	if !strings.HasSuffix(layout.Name, "\x00") {
		return ErrStringNotNullTerminated
	}
	vbo.Bind()
	vertAttrib := uint32(gl.GetAttribLocation(layout.Program.rid, gl.Str(layout.Name)))
	gl.EnableVertexAttribArray(vertAttrib)
	// VAO: Vertex Array Object is bound to the vertex buffer on this call.
	// What this line is saying is that `vertAttrib`` index is going to be bound
	// to the current gl.ARRAY_BUFFER (vbo).
	// It also stores size, type, normalized, stride and pointer as vertex array
	// state, in addition to the current vertex array buffer object binding. https://registry.khronos.org/OpenGL-Refpages/gl4/html/glVertexAttribPointer.xhtml
	gl.VertexAttribPointerWithOffset(vertAttrib, int32(layout.Packing), gl.FLOAT,
		layout.Normalize, int32(layout.Stride), 0)
	return glCheckError()
}

// VertexBuffer contains bytes, no information on the layout or type.
type VertexBuffer struct {
	// Renderer ID. If using OpenGL is the id set on buffer creation.
	rid uint32
}

func NewVertexBuffer[T any](data []T) (VertexBuffer, error) {
	return newVertexBuffer(gl.STATIC_DRAW, data)
}

func newVertexBuffer[T any](usage uint32, data []T) (VertexBuffer, error) {
	var vbo VertexBuffer
	vertexSize := unsafe.Sizeof(data[0])
	vertPtr := unsafe.Pointer(&data[0])
	gl.GenBuffers(1, &vbo.rid)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.rid)
	gl.BufferData(gl.ARRAY_BUFFER, int(vertexSize)*len(data), vertPtr, usage)
	return vbo, glCheckError()
}

func (vbo VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.rid)
}
func (vbo VertexBuffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (vbo VertexBuffer) Delete() {
	gl.DeleteBuffers(1, &vbo.rid)
}

type IndexBuffer struct {
	// Renderer ID. If using OpenGL is the id set on buffer creation.
	rid uint32
}

func NewIndexBuffer(data []uint32) (IndexBuffer, error) {
	return newIndexBuffer(gl.STATIC_DRAW, data)
}

func newIndexBuffer(usage uint32, data []uint32) (IndexBuffer, error) {
	var ibo IndexBuffer
	const IndexSize = unsafe.Sizeof(data[0])
	vertPtr := unsafe.Pointer(&data[0])
	gl.GenBuffers(1, &ibo.rid)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo.rid)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(IndexSize)*len(data), vertPtr, usage)
	return ibo, glCheckError()
}

func (vbo IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbo.rid)
}
func (vbo IndexBuffer) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}
func (vbo IndexBuffer) Delete() {
	gl.DeleteBuffers(1, &vbo.rid)
}

// Vertex and Fragment are null terminated strings with source code.
type ShaderSource struct {
	// Vertex and Fragment are null terminated strings with source code.
	Vertex   string
	Fragment string
}

func NewProgram(ss ShaderSource) (prog Program, err error) {
	prog.rid, err = shaders.CompileBasic(ss.Vertex, ss.Fragment)
	return prog, err
}

type Program struct {
	rid uint32
}

func (p Program) Bind() {
	gl.UseProgram(p.rid)
}

func (p Program) BindFrag(name string) error {
	if !strings.HasSuffix(name, "\x00") {
		return ErrStringNotNullTerminated
	}
	gl.BindFragDataLocation(p.rid, 0, gl.Str(name))
	return nil
}

func (p Program) Unbind() {
	gl.UseProgram(0)
}
func (p Program) Delete() { gl.DeleteProgram(p.rid) }

func (p Program) SetUniformName4f(name string, v0, v1, v2, v3 float32) error {
	if !strings.HasSuffix(name, "\x00") {
		return ErrStringNotNullTerminated
	}
	loc := gl.GetUniformLocation(p.rid, gl.Str(name))
	if loc < 0 {
		return errors.New("unable to find uniform in program- did you use the identifier so it was not stripped from program?")
	}
	gl.Uniform4f(loc, v0, v1, v2, v3)
	return nil
}

func glClearError() {
	for gl.GetError() != gl.NO_ERROR {
	}
}

func glCheckError() error {
	code := gl.GetError()
	// We have a nil context check since it makes no sense to
	// loop forever if there is no current context.
	if code == gl.NO_ERROR {
		return nil
	}
	errs := GLErrors{code}
	if glfw.GetCurrentContext() == nil {
		slog.Error("glfw context nil", errs)
		return nil
	}
	for {
		code = gl.GetError()
		if code == gl.NO_ERROR {
			return errs
		}
		errs = append(errs, code)
	}
}

type GLErrors []uint32

func (ge GLErrors) Error() (errstr string) {
	if len(ge) == 0 {
		return "no gl errors"
	}
	for i := range ge {
		var s string
		switch ge[i] {
		case gl.INVALID_ENUM:
			s = "invalid enum"
		case gl.INVALID_FRAMEBUFFER_OPERATION:
			s = "invalid framebuffer operation"
		case gl.INVALID_INDEX:
			s = "invalid index"
		case gl.INVALID_OPERATION:
			s = "invalid operation"
		case gl.INVALID_VALUE:
			s = "invalid value"
		default:
			s = "unknown error enum"
		}
		errstr += s
		if i != len(ge)-1 {
			errstr += "; "
		}
	}
	return errstr
}
