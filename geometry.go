package main

import (
	//"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Geometry : interface for all objects that will be rendered
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

// Attributes : struct for holding vertex attribute locations
type Attributes struct {
	position       uint32
	normal         uint32
	uv             uint32
	vertexPosition int32
	vertexNormal   int32
}

// ObjectBuffers : holds references to vertex buffers
type ObjectBuffers struct {
	vao        uint32
	vbo        uint32
	attributes Attributes
}

// VertexValues : struct for holding vertex specific values
type VertexValues struct {
	vertices []float32
	normals  []float32
	faces    []uint32
}

// Uniforms : struct for holding all uniforms
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

// ProgramInfo : struct for holding program info (program, uniforms, attributes)
type ProgramInfo struct {
	program          uint32
	uniformLocations Uniforms
	attributes       Attributes
	indexBuffer      uint32
}

// Material : struct for holding material info
type Material struct {
	diffuse  []float32
	ambient  []float32
	specular []float32
	n        float32
	texture  string
	alpha    float32
}

// Model : struct for holding model info
type Model struct {
	position mgl32.Vec3
	rotation mgl32.Mat4
	scale    mgl32.Vec3
}
