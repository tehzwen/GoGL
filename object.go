package main

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	//"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Geometry interface {
	Setup(Material, Model, string) error
	SetShader(string, string) error
	GetProgramInfo() (ProgramInfo, error)
	GetModel() (Model, error)
	GetCentroid() mgl32.Vec3
	GetMaterial() Material
	GetBuffers() ObjectBuffers
	GetVertices() VertexValues
	SetRotation(mgl32.Mat4)
}

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
	centroid    mgl32.Vec3
}

func (o *Object) SetShader(vertShader string, fragShader string) error {

	if vertShader != "" && fragShader != "" {
		o.fragShader = fragShader
		o.vertShader = vertShader
		return nil
	} else {
		return errors.New("Error setting the shader, shader code must not be blank")
	}
}

//WIP
func (o Object) Setup(mat Material, name string) error {
	return nil
}

func (o Object) GetBuffers() ObjectBuffers {
	return o.buffers
}

func (o Object) GetProgramInfo() (ProgramInfo, error) {
	if (o.programInfo != ProgramInfo{}) {
		return o.programInfo, nil
	}
	return ProgramInfo{}, errors.New("No program info!")
}

func (o Object) GetMaterial() Material {
	return o.material
}

func (o Object) GetModel() (Model, error) {
	if (o.model != Model{}) {
		return o.model, nil
	}
	return Model{}, errors.New("No model info!")
}

func (o Object) GetCentroid() mgl32.Vec3 {
	return o.centroid
}

type Attributes struct {
	position       uint32
	normal         uint32
	uv             uint32
	vertexPosition int32
	vertexNormal   int32
}

type ObjectBuffers struct {
	vao        uint32
	vbo        uint32
	attributes Attributes
}

type VertexValues struct {
	vertices []float32
	normals  []float32
	faces    []uint32
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
