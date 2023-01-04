package main

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"os"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/soypat/shaders"
	"golang.org/x/exp/slog"
)

// Very basic index buffer example.
const (
	projectName  = "Index Buffers"
	windowWidth  = 800
	windowHeight = 800
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

//go:embed uniformtriangle.glsl
var shader string

// Square with indices:
// 3----2
// |    |
// 0----1
var positions = []float32{
	-0.5, -0.5, // 0
	0.5, -0.5, // 1
	0.5, 0.5, // 2
	-0.5, 0.5, //3
}
var indices = []uint32{
	0, 1, 2, // Lower right triangle.
	0, 2, 3, // Upper left triangle.
}

func main() {
	if err := glfw.Init(); err != nil {
		slog.Error("failed to initialize glfw", err)
		os.Exit(1)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, projectName, nil, nil)
	if err != nil {
		slog.Error("create glfw window failed", err)
		return
	}
	window.MakeContextCurrent()
	// Initialize Glow
	if err := gl.Init(); err != nil {
		slog.Error("init glow fail", err)
		return
	}
	glClearError()
	// gl.DebugMessageCallback(func(source, gltype, id, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	// 	slog.Warn(message,
	// 		slog.Uint64("source", uint64(source)),
	// 		slog.Uint64("gltype", uint64(gltype)),
	// 		slog.Uint64("id", uint64(id)),
	// 		slog.Uint64("severity", uint64(severity)),
	// 		slog.Int("length", int(length)),
	// 	)
	// }, unsafe.Pointer(nil))
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Separate vertex and fragment shaders from source code.
	vertexSource, fragSource, err := shaders.ParseCombinedBasic(strings.NewReader(shader))
	if err != nil {
		slog.Error("parse combined source fail", err)
		return
	}

	// Configure the vertex and fragment shaders
	program, err := shaders.CompileBasic(vertexSource, fragSource)
	if err != nil {
		slog.Error("compile fail", err)
		return
	}
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// float32 is 4 bytes wide.
	const attrSize = 4

	// Configure the Vertex Array Object.
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Create the Position Buffer Object.
	var pbo uint32
	vertPtr := unsafe.Pointer(&positions[0])
	gl.GenBuffers(1, &pbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, pbo)
	gl.BufferData(gl.ARRAY_BUFFER, attrSize*len(positions), vertPtr, gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)

	// Create Index Buffer Object.
	var ibo uint32
	indPtr := unsafe.Pointer(&indices[0])
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, attrSize*len(indices), indPtr, gl.STATIC_DRAW)

	// Get Uniform variable location in program.
	colorUniform := gl.GetUniformLocation(program, gl.Str("u_color\x00"))
	if colorUniform < 0 {
		slog.Error("could not find uniform attribute u_color", nil)
	}
	// Set the uniform variable to a certain color.
	gl.Uniform4f(colorUniform, 0.2, 0.3, 0.8, 1)

	// VAO: Vertex Array Object is bound to the vertex buffer on this call.
	// What this line is saying is that `vertAttrib`` index is going to be bound
	// to the current gl.ARRAY_BUFFER (vbo).
	// It also stores size, type, normalized, stride and pointer as vertex array state, in addition to the current vertex array buffer object binding. https://registry.khronos.org/OpenGL-Refpages/gl4/html/glVertexAttribPointer.xhtml
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 2*attrSize, 0)
	if err := glCheckError(); err != nil {
		slog.Error("after setup", err)
	}
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// OpenGL is telling us 1280 == 0x0500.
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
		if err := glCheckError(); err != nil {
			// It is telling us the gl.INT enum is incorrect...
			// not very descriptive error if you ask me.
			slog.Error("rendering", err)
		}
		// We can also set the color in between renders.
		gl.Uniform4f(colorUniform, float32(time.Now().UnixMilli()%1000)/1000, .5, .3, 1)
		// Maintenance
		glfw.SwapInterval(1) // Can prevent epilepsy for high frequency
		window.SwapBuffers()
		glfw.PollEvents()
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}

func glClearError() {
	for gl.GetError() != gl.NO_ERROR {
	}
}

func glCheckError() error {
	code := gl.GetError()
	if code == gl.NO_ERROR {
		return nil
	}
	errs := GLErrors{code}
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
