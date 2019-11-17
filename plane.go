package main

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

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

func (p *Plane) SetShader(vertShader string, fragShader string) error {

	if vertShader != "" && fragShader != "" {
		p.fragShader = fragShader
		p.vertShader = vertShader
		return nil
	} else {
		return errors.New("Error setting the shader, shader code must not be blank")
	}
}

func (p Plane) GetProgramInfo() (ProgramInfo, error) {
	if (p.programInfo != ProgramInfo{}) {
		return p.programInfo, nil
	}
	return ProgramInfo{}, errors.New("No program info!")
}

func (p Plane) GetMaterial() Material {
	return p.material
}

func (p Plane) GetVertices() VertexValues {
	return p.vertexValues
}

func (p Plane) GetCentroid() mgl32.Vec3 {
	return p.centroid
}

func (p Plane) GetBuffers() ObjectBuffers {
	return p.buffers
}

func (p *Plane) SetRotation(rot mgl32.Mat4) {
	p.model.rotation = rot
}

func (p Plane) GetModel() (Model, error) {
	if (p.model != Model{}) {
		return p.model, nil
	}
	return Model{}, errors.New("No model info!")
}

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
