package main

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/soypat/shaders"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Hello triangle", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := shaders.CompileBasic(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// float32 is 4 bytes wide.
	const attrSize = 4
	var vbo uint32
	vertPtr := unsafe.Pointer(&triangleVertices[0])
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, attrSize*len(triangleVertices), vertPtr, gl.STATIC_DRAW)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)

	// Size is 2 since each vertex contains 2 gl.FLOATs
	// Stride is 2 since our data is 2D.
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 2*attrSize, 0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// NOTE: If nothing is visible maybe add a gl.BindVertexArray(vao) call in here and file a bug!
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}

//go:embed triangle_vert.glsl
var vertexShader string

//go:embed triangle_frag.glsl
var fragmentShader string

var triangleVertices = []float32{
	-0.5, -0.5,
	0.0, 0.5,
	0.5, -0.5,
}
