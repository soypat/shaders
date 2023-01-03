// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	_ "embed"
	_ "image/png"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/soypat/shaders"
	"golang.org/x/exp/slog"
)

const windowWidth = 800
const windowHeight = 600

// shader contains our source code by embedding.
//
//go:embed triangle.glsl
var shader string

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type colorVertex struct {
	pos   [2]float32
	color [4]float32
}

var vertices = []colorVertex{
	{pos: [2]float32{-0.5, -0.5}, color: [4]float32{1, 0, 0, 1}},
	{pos: [2]float32{0, 0.5}, color: [4]float32{0, 1, 0, 1}},
	{pos: [2]float32{0.5, -0.5}, color: [4]float32{0, 0, 1, 1}},
}

// use this buffer for time being.
var colors [][4]float32

func init() {
	for i := range vertices {
		colors = append(colors, vertices[i].color)
	}
}

func main() {
	if err := glfw.Init(); err != nil {
		slog.Error("glfw.Init failed", err)
		os.Exit(1)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Hello colorful triangle", nil, nil)
	if err != nil {
		slog.Error("create window failed", err)
		os.Exit(1)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		slog.Error("initializing openGL", err)
		os.Exit(1)
	}

	slog.Info("openGL init via Glow", slog.String("version", gl.GoStr(gl.GetString(gl.VERSION))))

	vertexSource, fragSource, err := shaders.ParseCombinedBasic(strings.NewReader(shader))
	if err != nil {
		slog.Error("parsing combined shaders", err)
		os.Exit(1)
	}
	// Configure the vertex and fragment shaders
	program, err := shaders.CompileBasic(vertexSource, fragSource)
	if err != nil {
		slog.Error("compiling program", err)
		os.Exit(1)
	}
	gl.UseProgram(program)
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// float32 is 4 bytes wide.
	const attrSize = 4
	const vertexSize = attrSize * (2 + 4) // 2 positions, 4 colors
	var vbo uint32
	vertPtr := unsafe.Pointer(&vertices[0])
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, vertexSize*len(vertices), vertPtr, gl.STATIC_DRAW)

	posAttr := gl.GetAttribLocation(program, gl.Str("vert\x00"))
	if posAttr < 0 {
		slog.Error("vert attr not found", nil, slog.Int("index", int(posAttr)))
		os.Exit(1)
	}
	gl.EnableVertexAttribArray(uint32(posAttr))
	gl.VertexAttribPointerWithOffset(uint32(posAttr), 2, gl.FLOAT, false, vertexSize, 0)

	// setting up colors.
	// colorPtr := unsafe.Pointer(&colors[0])
	// var vboColor uint32
	// gl.GenBuffers(1, &vboColor)
	// gl.BindBuffer(gl.ARRAY_BUFFER, vboColor)
	// gl.BufferData(gl.ARRAY_BUFFER, len(colors[0])*attrSize*len(colors), colorPtr, gl.STATIC_DRAW)
	vcolorAttr := gl.GetAttribLocation(program, gl.Str("vert_color\x00"))
	if vcolorAttr < 0 {
		slog.Error("color not found", nil, slog.Int("index", int(vcolorAttr)))
		os.Exit(1)
	}
	colorAttr := gl.GetAttribLocation(program, gl.Str("color\x00"))
	gl.EnableVertexAttribArray(uint32(vcolorAttr))
	gl.VertexAttribPointerWithOffset(uint32(vcolorAttr), 4, gl.FLOAT, true, vertexSize, 0)

	slog.Info("program attr",
		slog.Int("poslayout", int(posAttr)),
		slog.Int("vcolorlayout", int(vcolorAttr)),
		slog.Int("colorlayout", int(colorAttr)),
		slog.Any("colors", colors),
	)

	for !window.ShouldClose() {
		gl.ClearColor(0.4, 0.4, 0.6, 1)
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
