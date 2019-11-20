package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// InitOpenGL : initializes OpenGL and returns an intiialized program.
func InitOpenGL(vertexShaderSource string, fragmentShaderSource string) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	/*
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)*/

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog
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

// SetupAttributes : helper function for setting up attribute/uniforms
func SetupAttributes(p *ProgramInfo) {
	(*p).attributes.vertexPosition = gl.GetAttribLocation((*p).program, gl.Str("aPosition\x00"))
	(*p).attributes.vertexNormal = gl.GetAttribLocation((*p).program, gl.Str("aNormal\x00"))
	(*p).uniformLocations.diffuseVal = gl.GetUniformLocation((*p).program, gl.Str("diffuseVal\x00"))
	(*p).uniformLocations.ambientVal = gl.GetUniformLocation((*p).program, gl.Str("ambientVal\x00"))
	(*p).uniformLocations.specularVal = gl.GetUniformLocation((*p).program, gl.Str("specularVal\x00"))
	(*p).uniformLocations.nVal = gl.GetUniformLocation((*p).program, gl.Str("nVal\x00"))
	(*p).uniformLocations.projection = gl.GetUniformLocation((*p).program, gl.Str("uProjectionMatrix\x00"))
	(*p).uniformLocations.view = gl.GetUniformLocation((*p).program, gl.Str("uViewMatrix\x00"))
	(*p).uniformLocations.model = gl.GetUniformLocation((*p).program, gl.Str("uModelMatrix\x00"))
	(*p).uniformLocations.lightPositions = gl.GetUniformLocation((*p).program, gl.Str("lightPositions\x00"))
	(*p).uniformLocations.lightColours = gl.GetUniformLocation((*p).program, gl.Str("lightColours\x00"))
	(*p).uniformLocations.lightStrengths = gl.GetUniformLocation((*p).program, gl.Str("lightStrengths\x00"))
	(*p).uniformLocations.numLights = gl.GetUniformLocation((*p).program, gl.Str("numLights\x00"))
	(*p).uniformLocations.cameraPosition = gl.GetUniformLocation((*p).program, gl.Str("cameraPosition\x00"))

	if (*p).attributes.vertexPosition == -1 ||
		(*p).attributes.vertexNormal == -1 ||
		(*p).uniformLocations.projection == -1 ||
		(*p).uniformLocations.view == -1 ||
		(*p).uniformLocations.model == -1 ||
		(*p).uniformLocations.cameraPosition == -1 ||
		(*p).uniformLocations.lightPositions == -1 ||
		(*p).uniformLocations.lightColours == -1 ||
		(*p).uniformLocations.lightStrengths == -1 ||
		(*p).uniformLocations.numLights == -1 ||
		(*p).uniformLocations.diffuseVal == -1 ||
		(*p).uniformLocations.ambientVal == -1 ||
		(*p).uniformLocations.nVal == -1 ||
		(*p).uniformLocations.specularVal == -1 {
		fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
	}
}

// CalculateCentroid : helper function for calculating the centroid of the geometry
func CalculateCentroid(vertices []float32) mgl32.Vec3 {
	var xTotal = float32(0.0)
	var yTotal = float32(0.0)
	var zTotal = float32(0.0)

	for i := 0; i < len(vertices); i += 3 {
		xTotal += vertices[i]
		yTotal += vertices[i+1]
		zTotal += vertices[i+2]
	}

	xTotal /= float32(len(vertices) / 3)
	yTotal /= float32(len(vertices) / 3)
	zTotal /= float32(len(vertices) / 3)

	return mgl32.Vec3{xTotal, yTotal, zTotal}

}

// CreateTriangleVAO : helper function for creating vertex buffer values
func CreateTriangleVAO(programInfo *ProgramInfo, vertices []float32, normals []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)
	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	if indices != nil {
		var EBO uint32
		gl.GenBuffers(1, &EBO)
		// copy indices into element buffer
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	}

	// position
	gl.VertexAttribPointer((*programInfo).attributes.position, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray((*programInfo).attributes.position)

	//normals
	if normals != nil {
		var normalBuffer uint32
		gl.GenBuffers(1, &normalBuffer)

		gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(normals)*4, gl.Ptr(normals), gl.STATIC_DRAW)

		gl.VertexAttribPointer((*programInfo).attributes.normal, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray((*programInfo).attributes.normal)
	}

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray((*programInfo).attributes.position)

	return VAO
}
