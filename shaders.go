package shaders

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// CompileBasic compiles two OpenGL vertex and fragment shaders
// and returns a program with the current OpenGL context.
// It returns an error if compilation, linking or validation fails.
func CompileBasic(vertexSrcCode, fragmentSrcCode string) (program uint32, err error) {
	program = gl.CreateProgram()
	vid, err := compile(gl.VERTEX_SHADER, vertexSrcCode)
	if err != nil {
		return 0, fmt.Errorf("vertex shader compile: %w", err)
	}
	fid, err := compile(gl.FRAGMENT_SHADER, fragmentSrcCode)
	if err != nil {
		return 0, fmt.Errorf("fragment shader compile: %w", err)
	}
	gl.AttachShader(program, vid)
	gl.AttachShader(program, fid)
	gl.LinkProgram(program)
	log := ivLog(program, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog)
	if len(log) > 0 {
		return 0, fmt.Errorf("link failed: %v", log)
	}
	// We should technically call DetachShader after linking... https://www.youtube.com/watch?v=71BLZwRGUJE&list=PLlrATfBNZ98foTJPJ_Ev03o2oq3-GGOS2&index=7&ab_channel=TheCherno
	gl.ValidateProgram(program)
	log = ivLog(program, gl.VALIDATE_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog)
	if len(log) > 0 {
		return 0, fmt.Errorf("validation failed: %v", log)
	}

	// We can clean up.
	gl.DeleteShader(vid)
	gl.DeleteShader(fid)
	return program, nil
}

func compile(shaderType uint32, sourceCode string) (uint32, error) {
	id := gl.CreateShader(shaderType)
	csources, free := gl.Strs(sourceCode)
	gl.ShaderSource(id, 1, csources, nil)
	free()
	gl.CompileShader(id)

	// We now check the errors during compile, if there were any.
	log := ivLog(id, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog)
	if len(log) > 0 {
		return 0, errors.New(log)
	}
	return id, nil
}

// ivLog is a helper function for extracting log data
// from a Shader compilation step or program linking.
//
//	log := ivLog(id, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog)
//	if len(log) > 0 {
//		return 0, errors.New(log)
//	}
func ivLog(id, plName uint32, getIV func(program uint32, pname uint32, params *int32), getInfo func(program uint32, bufSize int32, length *int32, infoLog *uint8)) string {
	var iv int32
	getIV(id, plName, &iv)
	if iv == gl.FALSE {
		var logLength int32
		getIV(id, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength)
		getInfo(id, logLength, &logLength, &log[0])
		return string(log[:len(log)-1]) // we exclude the last null character.
	}
	return ""
}
