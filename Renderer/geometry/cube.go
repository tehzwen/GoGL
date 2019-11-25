package geometry

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
	}

	return errors.New("Error setting the shader, shader code must not be blank")

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

func (c Cube) GetType() string {
	return "cube"
}

// SetRotation : helper function for setting rotation of cube to a mat4
func (c *Cube) SetRotation(rot mgl32.Mat4) {
	c.model.Rotation = rot
}

// GetModel : getter for model values
func (c Cube) GetModel() (Model, error) {
	if (c.model != Model{}) {
		return c.model, nil
	}
	return Model{}, errors.New("no model info")
}

// Scale : function used to scale the cube and recalculate the centroid
func (c *Cube) Scale(scaleVec mgl32.Vec3) {
	c.model.Scale = scaleVec
	c.centroid = CalculateCentroid(c.vertexValues.Vertices, c.model.Scale)
}

func (c *Cube) Translate(translateVec mgl32.Vec3) {
	c.model.Position = c.model.Position.Add(translateVec)
	c.centroid = c.centroid.Add(translateVec)
}

// Setup : function for initializing cube
func (c *Cube) Setup(mat Material, mod Model, name string) error {
	c.name = name
	c.programInfo = ProgramInfo{}
	c.programInfo.Program = InitOpenGL(c.vertShader, c.fragShader)
	c.programInfo.attributes = Attributes{
		position: 0,
		normal:   1,
	}

	c.material = mat
	c.vertexValues.Vertices = []float32{
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
	c.Scale(mod.Scale)
	c.model.Position = mod.Position
	c.model.Rotation = mod.Rotation
	c.centroid = CalculateCentroid(c.vertexValues.Vertices, c.model.Scale)
	c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, c.vertexValues.normals, c.vertexValues.faces)

	return nil
}
