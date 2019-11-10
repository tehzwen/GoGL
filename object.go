package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Object struct {
	name        string
	fragShader  string
	vertShader  string
	buffers     ObjectBuffers
	programInfo ProgramInfo
	vertices    []float32
	faces       []float32
	material    Material
	model       Model
}

type Attributes struct {
	position       uint32
	normal         uint32
	uv             uint32
	vertexPosition int32
}

type ObjectBuffers struct {
	vao        uint32
	vbo        uint32
	attributes Attributes
}

type Uniforms struct {
	projection     int32
	view           int32
	model          int32
	normalMatrix   int32
	diffuseVal     int32
	ambientVal     int32
	specularVal    int32
	nVal           int32
	cameraPosition int32
	numLights      int32
	lightPositions int32
	lightColours   int32
	lightStrengths int32
}

type ProgramInfo struct {
	program          uint32
	uniformLocations Uniforms
	attributes       Attributes
	indexBuffer      uint32
}

type Material struct {
	diffuse  []float32
	ambient  []float32
	specular []float32
	n        float32
	texture  string
	alpha    float32
}

type Model struct {
	position mgl32.Vec3
	rotation mgl32.Mat4
	scale    mgl32.Vec3
}

func InitIndexBuffer(indices []uint32) uint32 {
	var indexBuffer uint32
	gl.GenBuffers(1, &indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	return indexBuffer
}

func SetupAttributes(p *ProgramInfo) {
	//fmt.Print("%d\n", (*p).uniformLocations.diffuseVal)
	(*p).attributes.vertexPosition = gl.GetAttribLocation((*p).program, gl.Str("aPosition\x00"))
	(*p).uniformLocations.diffuseVal = gl.GetUniformLocation((*p).program, gl.Str("diffuseVal\x00"))
	(*p).uniformLocations.ambientVal = gl.GetUniformLocation((*p).program, gl.Str("ambientVal\x00"))
	(*p).uniformLocations.specularVal = gl.GetUniformLocation((*p).program, gl.Str("specularVal\x00"))
	(*p).uniformLocations.projection = gl.GetUniformLocation((*p).program, gl.Str("uProjectionMatrix\x00"))
	(*p).uniformLocations.view = gl.GetUniformLocation((*p).program, gl.Str("uViewMatrix\x00"))
	(*p).uniformLocations.model = gl.GetUniformLocation((*p).program, gl.Str("uModelMatrix\x00"))

	if (*p).attributes.vertexPosition == -1 ||
		(*p).uniformLocations.projection == -1 ||
		(*p).uniformLocations.view == -1 ||
		(*p).uniformLocations.model == -1 ||
		(*p).uniformLocations.diffuseVal == -1 ||
		(*p).uniformLocations.ambientVal == -1 ||
		(*p).uniformLocations.specularVal == -1 {
		fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
	}
}

func CreateTriangleVAO(vertices []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}
