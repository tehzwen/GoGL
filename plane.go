package main

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

// Plane : Primitive plane geometry struct
type Plane struct {
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
func (p *Plane) SetShader(vertShader string, fragShader string) error {

	if vertShader != "" && fragShader != "" {
		p.fragShader = fragShader
		p.vertShader = vertShader
		return nil
	} else {
		return errors.New("Error setting the shader, shader code must not be blank")
	}
}

// GetProgramInfo : getter for programinfo
func (p Plane) GetProgramInfo() (ProgramInfo, error) {
	if (p.programInfo != ProgramInfo{}) {
		return p.programInfo, nil
	}
	return ProgramInfo{}, errors.New("No program info")
}

// GetMaterial : getter for material
func (p Plane) GetMaterial() Material {
	return p.material
}

// GetVertices : getter for vertexValues
func (p Plane) GetVertices() VertexValues {
	return p.vertexValues
}

// GetCentroid : getter for centroid
func (p Plane) GetCentroid() mgl32.Vec3 {
	return p.centroid
}

// GetBuffers : getter for buffers
func (p Plane) GetBuffers() ObjectBuffers {
	return p.buffers
}

// SetRotation : helper function for setting rotation of cube to a mat4
func (p *Plane) SetRotation(rot mgl32.Mat4) {
	p.model.rotation = rot
}

// GetModel : getter for model values
func (p Plane) GetModel() (Model, error) {
	if (p.model != Model{}) {
		return p.model, nil
	}
	return Model{}, errors.New("no model info")
}

// Setup : function for initializing plane
func (p *Plane) Setup(mat Material, mod Model, name string) error {

	p.name = name
	p.programInfo = ProgramInfo{}

	p.programInfo.program = InitOpenGL(p.vertShader, p.fragShader)
	//fmt.Printf("%v+\n", p)
	//fmt.Printf("\nhere!\n")

	p.programInfo.attributes = Attributes{
		position: 0,
		normal:   1,
	}

	p.material = mat

	p.vertexValues.vertices = []float32{
		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.5,
	}

	p.vertexValues.faces = []uint32{
		0, 1, 2, 0, 2, 3,
	}

	p.vertexValues.normals = []float32{
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
	}
	SetupAttributes(&p.programInfo)
	//p.model.position = mgl32.Vec3{1, 0, 0}
	//p.model.scale = mgl32.Vec3{1, 1, 1}
	//p.model.rotation = mgl32.Ident4()
	p.model = mod
	p.centroid = CalculateCentroid(p.vertexValues.vertices)
	p.buffers.vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.vertices, p.vertexValues.normals, p.vertexValues.faces)

	// x := errors.New("Wrong")
	// return x
	return nil
}
