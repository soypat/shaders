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

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Separate vertex and fragment shaders from source code.
	vertexSource, fragSource, err := shaders.ParseCombinedBasic(strings.NewReader(shader))
	if err != nil {
		slog.Error("parse combined source fail", err)
		return
	}

	// Configure the vertex and fragment shaders
	program, err := NewProgram(ShaderSource{Vertex: vertexSource, Fragment: fragSource})
	if err != nil {
		slog.Error("compile fail", err)
		return
	}
	defer program.Delete()
	program.Bind()

	err = program.BindFrag("outputColor\x00")
	if err != nil {
		slog.Error("program bind frag fail", err)
		return
	}
	// Configure the Vertex Array Object.
	vao := NewVAO()

	// Create the Position Buffer Object.
	vbo, err := NewVertexBuffer(positions)
	if err != nil {
		slog.Error("creating positions vertex buffer", err)
		return
	}
	err = vao.AddAttribute(vbo, AttribLayout{
		Program: program,
		Type:    gl.FLOAT,
		Name:    "vert\x00",
		Packing: 2,
		Stride:  2 * 4, // 2 floats, each 4 bytes wide.
	})
	if err != nil {
		slog.Error("adding attribute vert", err)
		return
	}

	// Create Index Buffer Object.
	_, err = NewIndexBuffer(indices)
	if err != nil {
		slog.Error("creating index buffer", err)
		return
	}

	// Set uniform variable `u_color` in source code.
	err = program.SetUniformName4f("u_color\x00", 0.2, 0.3, 0.8, 1)
	if err != nil {
		slog.Error("creating index buffer", err)
		return
	}
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))

		program.SetUniformName4f("u_color\x00", float32(time.Now().UnixMilli()%1000)/1000, .5, .3, 1)
		// Maintenance
		glfw.SwapInterval(1) // Can prevent epilepsy for high frequency
		window.SwapBuffers()
		glfw.PollEvents()
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}
