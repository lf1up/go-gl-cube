package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

var vertexShader = `
	#version 120

	uniform mat4 projection;
	uniform mat4 modelView;
	uniform mat4 normal;

	uniform mat3 selectedTriangle;

	attribute vec3 vertPosition;
	attribute vec3 vertNormal;
	attribute vec2 vertTexCoord;
	attribute vec3 vertColor;

	varying vec3 fragPosition;
	varying vec3 fragNormal;
	varying vec2 fragTexCoord;
	varying vec3 fragColor;
	varying float fragSelected;

	void main() {
			vec4 vertPosition4 = modelView * vec4(vertPosition, 1.0);

			fragPosition = vec3(vertPosition4) / vertPosition4.w;
			fragNormal = vec3(normal * vec4(vertNormal, 0.0));

			fragTexCoord = vertTexCoord;
			fragColor = vertColor;

			if (vertPosition == selectedTriangle[0] || vertPosition == selectedTriangle[1] || vertPosition == selectedTriangle[2])
				fragSelected = 1.0;
			else
			  fragSelected = 0.0;

			gl_Position = projection * modelView * vec4(vertPosition, 1.0);
	}
` + "\x00"

// Blinn–Phong shading model with gamma correction
var fragmentShader = `
	#version 120

	uniform sampler2D texture;

	varying vec3 fragPosition;
	varying vec3 fragNormal;
	varying vec2 fragTexCoord;
	varying vec3 fragColor;
	varying float fragSelected;

	const vec3 lightPosition = vec3(3.0, 3.0, 3.0);
	const vec3 lightColor = vec3(1.0, 1.0, 1.0);
	const vec3 ambientColor = vec3(0.3, 0.3, 0.3);
	const vec3 diffuseColor = vec3(0.5, 0.5, 0.5);
	const vec3 specularColor = vec3(1.0, 1.0, 1.0);

	const float lightPower = 40.0;
	const float shininess = 16.0;
	const float screenGamma = 2.2; // Assume the monitor is calibrated to the sRGB color space

	const int mode = 1;

	void main() {
		vec3 normal = normalize(fragNormal);
		vec3 lightDir = lightPosition - fragPosition;
		float distance = length(lightDir);
		distance = distance * distance;
		lightDir = normalize(lightDir);
	
		float lambertian = max(dot(lightDir, normal), 0.0);
		float specular = 0.0;
	
		if (lambertian > 0.0) {
	
			vec3 viewDir = normalize(-fragPosition);
	
			// this is blinn phong
			vec3 halfDir = normalize(lightDir + viewDir);
			float specAngle = max(dot(halfDir, normal), 0.0);
			specular = pow(specAngle, shininess);
				 
			// this is phong (for comparison)
			if (mode == 2) {
				vec3 reflectDir = reflect(-lightDir, normal);
				specAngle = max(dot(reflectDir, viewDir), 0.0);
				// note that the exponent is different here
				specular = pow(specAngle, shininess/4.0);
			}
		}
		vec3 colorLinear = ambientColor +
											 diffuseColor * lambertian * lightColor * lightPower / distance +
											 specularColor * specular * lightColor * lightPower / distance;
		// apply gamma correction (assume ambientColor, diffuseColor and specularColor
		// have been linearized, i.e. have no gamma correction in them)
		vec3 colorGammaCorrected = pow(colorLinear, vec3(1.0 / screenGamma));
		// use the gamma corrected color in the fragment
		if (fragSelected == 1.0)
			gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
		else
			gl_FragColor = texture2D(texture, fragTexCoord) * vec4(fragColor, 1.0) * vec4(colorGammaCorrected, 1.0);
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
	defer free()

	gl.ShaderSource(shader, 1, csources, nil)
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
