package main

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

// Cube : Primitive cube geometry struct
type Cube struct {
	name         string
	fragShader   string
	vertShader   string
	buffers      ObjectBuffers
	programInfo  ProgramInfo
	material     Material
	model        Model
	centroid     mgl32.Vec3
	vertexValues VertexValues
}

// SetShader : helper function for applying frag/vert shader
func (c *Cube) SetShader(vertShader string, fragShader string) error {

	if vertShader != "" && fragShader != "" {
		c.fragShader = fragShader
		c.vertShader = vertShader
		return nil
	} else {
		return errors.New("Error setting the shader, shader code must not be blank")
	}
}

// GetProgramInfo : getter for programinfo
func (c Cube) GetProgramInfo() (ProgramInfo, error) {
	if (c.programInfo != ProgramInfo{}) {
		return c.programInfo, nil
	}
	return ProgramInfo{}, errors.New("No program info")
}

// GetMaterial : getter for material
func (c Cube) GetMaterial() Material {
	return c.material
}

// GetVertices : getter for vertexValues
func (c Cube) GetVertices() VertexValues {
	return c.vertexValues
}

// GetCentroid : getter for centroid
func (c Cube) GetCentroid() mgl32.Vec3 {
	return c.centroid
}

// GetBuffers : getter for buffers
func (c Cube) GetBuffers() ObjectBuffers {
	return c.buffers
}

// SetRotation : helper function for setting rotation of cube to a mat4
func (c *Cube) SetRotation(rot mgl32.Mat4) {
	c.model.rotation = rot
}

// GetModel : getter for model values
func (c Cube) GetModel() (Model, error) {
	if (c.model != Model{}) {
		return c.model, nil
	}
	return Model{}, errors.New("no model info")
}

// Setup : function for initializing cube
func (c *Cube) Setup(mat Material, mod Model, name string) error {
	c.name = name
	c.programInfo = ProgramInfo{}

	c.programInfo.program = InitOpenGL(c.vertShader, c.fragShader)
	//fmt.Printf("%v+\n", c)
	//fmt.Printf("\nhere!\n")

	c.programInfo.attributes = Attributes{
		position: 0,
		normal:   1,
	}

	c.material = mat

	c.vertexValues.vertices = []float32{
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.0, 0.0,

		0.0, 0.0, 0.5,
		0.0, 0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.0, 0.5,

		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.5,

		0.0, 0.0, 0.5,
		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.0, 0.0, 0.0,

		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.0,

		0.0, 0.0, 0.5,
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
	}

	c.vertexValues.faces = []uint32{
		0, 1, 2, 0, 2, 3,
		4, 5, 6, 4, 6, 7,
		8, 9, 10, 8, 10, 11,
		12, 13, 14, 12, 14, 15,
		16, 17, 18, 17, 18, 19,
		20, 21, 22, 21, 22, 23,
	}

	c.vertexValues.normals = []float32{
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,

		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,

		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,

		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,

		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,

		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
	}
	SetupAttributes(&c.programInfo)
	//c.model.position = mgl32.Vec3{1, 0, 0}
	//c.model.scale = mgl32.Vec3{1, 1, 1}
	//c.model.rotation = mgl32.Ident4()
	c.model = mod
	c.centroid = CalculateCentroid(c.vertexValues.vertices)
	c.buffers.vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.vertices, c.vertexValues.normals, c.vertexValues.faces)

	// x := errors.New("Wrong")
	// return x
	return nil
}
