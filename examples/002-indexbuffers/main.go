package main

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/soypat/shaders"
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

//go:embed triangle.glsl
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
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, projectName, nil, nil)
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

	// Separate vertex and fragment shaders from source code.
	vertexSource, fragSource, err := shaders.ParseCombinedBasic(strings.NewReader(shader))
	if err != nil {
		panic(err)
	}

	// Configure the vertex and fragment shaders
	program, err := shaders.CompileBasic(vertexSource, fragSource)
	if err != nil {
		panic(err)
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
	// Size is 2 since each vertex contains 2 gl.FLOATs
	// Stride is 2 since our data is 2D.
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 2*attrSize, 0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// BEWARE WITH gl.UNSIGNED_INT, all buffer indices take unsigned integers.
		// If gl.INT is used nothing will be drawn.
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}
