package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

var vertexShader = `
	#version 120

	uniform mat4 projection;
	uniform mat4 view;
	uniform mat4 model;

	attribute vec3 vertPosition;
	attribute vec3 vertNormal;
	attribute vec2 vertTexCoord;
	attribute vec3 vertColor;

	varying vec3 fragNormal;
	varying vec2 fragTexCoord;
	varying vec3 fragColor;

	void main() {
			fragNormal = vertNormal;
			fragTexCoord = vertTexCoord;
			fragColor = vertColor;
			gl_Position = projection * view * model * vec4(vertPosition, 1.0);
	}
` + "\x00"

var fragmentShader = `
	#version 120

	uniform sampler2D tex;

	varying vec3 fragNormal;
	varying vec2 fragTexCoord;
	varying vec3 fragColor;

	void main() {
			gl_FragColor = texture2D(tex, fragTexCoord) * vec4(fragColor, 1.0);
	}
` + "\x00"

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
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
