package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600
const floatSize = 4

var (
	projection       mgl32.Mat4
	view             mgl32.Mat4
	model            mgl32.Mat4
	selectedTriangle mgl32.Mat3
	mouseX           float32
	mouseY           float32
)

func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	mouseX = float32(xpos/(windowWidth*0.5) - 1.0)
	mouseY = float32(-(ypos/(windowHeight*0.5) - 1.0))
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Renders a textured spinning cube using GLFW 3.3 and OpenGL 2.1.
func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "CUBE", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Set mouse tracking callback
	window.SetCursorPosCallback(cursorPosCallback)

	// Initialize Gl
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	view = mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model = mgl32.Ident4()

	modelView := view.Mul4(model)
	modelViewUniform := gl.GetUniformLocation(program, gl.Str("modelView\x00"))
	gl.UniformMatrix4fv(modelViewUniform, 1, false, &modelView[0])

	normal := (modelView.Inv()).Transpose()
	normalUniform := gl.GetUniformLocation(program, gl.Str("normal\x00"))
	gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("texture\x00"))
	gl.Uniform1i(textureUniform, 0)

	selectedTriangle = mgl32.Ident3()
	selectedTriangle.SetCol(0, mgl32.Vec3{1, 1, 1})
	selectedTriangle.SetCol(1, mgl32.Vec3{-1, 1, 1})
	selectedTriangle.SetCol(2, mgl32.Vec3{-1, -1, 1})
	selectedTriangleUniform := gl.GetUniformLocation(program, gl.Str("selectedTriangle\x00"))
	gl.UniformMatrix3fv(selectedTriangleUniform, 1, false, &selectedTriangle[0])

	// Load the texture
	texture, err := newTexture("./res/square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*floatSize, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cubeIndices)*floatSize, gl.Ptr(cubeIndices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertPosition\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 11*floatSize, gl.PtrOffset(0))

	normalAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertNormal\x00")))
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointer(normalAttrib, 3, gl.FLOAT, false, 11*floatSize, gl.PtrOffset(3))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 11*floatSize, gl.PtrOffset(6*floatSize))

	colorAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 3, gl.FLOAT, false, 11*floatSize, gl.PtrOffset(8*floatSize))

	gl.BindVertexArray(0)
	gl.UseProgram(0)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.ClearColor(0.7, 0.7, 0.7, 1.0)
	gl.ClearStencil(0)
	gl.ClearDepth(1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle -= elapsed

		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		modelView = view.Mul4(model)
		normal = (modelView.Inv()).Transpose()

		// Ray-triangle intersection (with mouse coordinates)
		if (mouseX >= -1 && mouseX <= 1) && (mouseY >= -1 && mouseY <= 1) {
			// log.Printf("[Debug] Mouse position (x y): %v %v\n", mouseX, mouseY)

			invProjection := projection.Inv()
			invView := view.Inv()
			invModel := model.Inv()

			viewP1 := mgl32.TransformCoordinate(mgl32.Vec3{mouseX, mouseY, -1.0}, invProjection)

			R0 := mgl32.TransformCoordinate(mgl32.TransformCoordinate(mgl32.Vec3{0, 0, 0}, invView), invModel)
			R1 := mgl32.TransformCoordinate(mgl32.TransformCoordinate(viewP1, invView), invModel)
			D := mgl32.Vec3{R1[0] - R0[0], R1[1] - R0[1], R1[2] - R0[2]}.Normalize()

			triangleIsectIndex := -1
			minDist := float32(100000)
			for it := 0; it < len(cubeIndices); it += 3 {
				triangle := []int32{cubeIndices[it+0], cubeIndices[it+1], cubeIndices[it+2]}
				A := mgl32.Vec3{cubeVertices[triangle[0]*11+0], cubeVertices[triangle[0]*11+1], cubeVertices[triangle[0]*11+2]}
				B := mgl32.Vec3{cubeVertices[triangle[1]*11+0], cubeVertices[triangle[1]*11+1], cubeVertices[triangle[1]*11+2]}
				C := mgl32.Vec3{cubeVertices[triangle[2]*11+0], cubeVertices[triangle[2]*11+1], cubeVertices[triangle[2]*11+2]}

				P0 := A
				NV := mgl32.Vec3{B[0] - A[0], B[1] - A[1], B[2] - A[2]}.Cross(mgl32.Vec3{C[0] - A[0], C[1] - A[1], C[2] - A[2]}).Normalize()

				distIsect := mgl32.Vec3{P0[0] - R0[0], P0[1] - R0[1], P0[2] - R0[2]}.Dot(NV) / D.Dot(NV)
				if distIsect < 0.0 {
					continue
				}

				PIsect := mgl32.Vec3{R0[0] + D[0]*distIsect, R0[1] + D[1]*distIsect, R0[2] + D[2]*distIsect}

				if PointInOrOnTriangle(PIsect, A, B, C) {
					if distIsect < minDist {
						minDist = distIsect
						triangleIsectIndex = it / 3
					}
				}
			}

			if triangleIsectIndex >= 0 {
				// log.Printf("[Debug] Mouse is ON Triangle with Index: %v\n", triangleIsectIndex)

				triangle := []int32{cubeIndices[triangleIsectIndex*3+0], cubeIndices[triangleIsectIndex*3+1], cubeIndices[triangleIsectIndex*3+2]}
				selectedTriangle.SetCol(0, mgl32.Vec3{cubeVertices[triangle[0]*11+0], cubeVertices[triangle[0]*11+1], cubeVertices[triangle[0]*11+2]})
				selectedTriangle.SetCol(1, mgl32.Vec3{cubeVertices[triangle[1]*11+0], cubeVertices[triangle[1]*11+1], cubeVertices[triangle[1]*11+2]})
				selectedTriangle.SetCol(2, mgl32.Vec3{cubeVertices[triangle[2]*11+0], cubeVertices[triangle[2]*11+1], cubeVertices[triangle[2]*11+2]})
			} else {
				selectedTriangle = mgl32.Ident3()
			}
		}

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelViewUniform, 1, false, &modelView[0])
		gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])
		gl.UniformMatrix3fv(selectedTriangleUniform, 1, false, &selectedTriangle[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		// Draw cube
		gl.DrawElements(gl.TRIANGLES, int32(len(cubeIndices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Draw additional cube
		//
		// modelAdditional := model.Mul4(mgl32.Translate3D(2, 0, 0))
		// modelView = view.Mul4(modelAdditional)
		// normal = (modelView.Inv()).Transpose()
		// gl.UniformMatrix4fv(modelViewUniform, 1, false, &modelView[0])
		// gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])
		// gl.DrawElements(gl.TRIANGLES, int32(len(cubeIndices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

		gl.BindVertexArray(0)
		gl.UseProgram(0)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func PointInOrOn(P1 mgl32.Vec3, P2 mgl32.Vec3, A mgl32.Vec3, B mgl32.Vec3) bool {
	CP1 := mgl32.Vec3{B[0] - A[0], B[1] - A[1], B[2] - A[2]}.Cross(mgl32.Vec3{P1[0] - A[0], P1[1] - A[1], P1[2] - A[2]})
	CP2 := mgl32.Vec3{B[0] - A[0], B[1] - A[1], B[2] - A[2]}.Cross(mgl32.Vec3{P2[0] - A[0], P2[1] - A[1], P2[2] - A[2]})
	return CP1.Dot(CP2) >= 0
}

func PointInOrOnTriangle(P mgl32.Vec3, A mgl32.Vec3, B mgl32.Vec3, C mgl32.Vec3) bool {
	var isInA = PointInOrOn(P, A, B, C)
	var isInB = PointInOrOn(P, B, C, A)
	var isInC = PointInOrOn(P, C, A, B)
	return isInA && isInB && isInC
}
