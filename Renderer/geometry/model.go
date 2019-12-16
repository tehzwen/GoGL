package geometry

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

// ModelObject : Primitive ModelObject geometry struct
type ModelObject struct {
	name         string
	fragShader   string
	vertShader   string
	shaderType   string
	boundingBox  BoundingBox
	buffers      ObjectBuffers
	programInfo  ProgramInfo
	material     Material
	Model        Model
	centroid     mgl32.Vec3
	vertexValues VertexValues
}

// SetShader : helper function for applying frag/vert shader
func (m *ModelObject) SetShader(vertShader string, fragShader string) error {

	if vertShader != "" && fragShader != "" {
		m.fragShader = fragShader
		m.vertShader = vertShader
		return nil
	}

	return errors.New("Error setting the shader, shader code must not be blank")

}

// GetProgramInfo : getter for programinfo
func (m ModelObject) GetProgramInfo() (ProgramInfo, error) {
	if (m.programInfo != ProgramInfo{}) {
		return m.programInfo, nil
	}
	return ProgramInfo{}, errors.New("No program info")
}

// GetMaterial : getter for material
func (m ModelObject) GetMaterial() Material {
	return m.material
}

func (m ModelObject) GetDetails() (string, string) {
	return m.name, m.shaderType
}

// GetVertices : getter for vertexValues
func (m ModelObject) GetVertices() VertexValues {
	return m.vertexValues
}

func (m ModelObject) GetBoundingBox() BoundingBox {
	return m.boundingBox
}

// GetCentroid : getter for centroid
func (m ModelObject) GetCentroid() mgl32.Vec3 {
	return m.centroid
}

// GetBuffers : getter for buffers
func (m ModelObject) GetBuffers() ObjectBuffers {
	return m.buffers
}

func (m ModelObject) GetType() string {
	return "mesh"
}

// SetRotation : helper function for setting rotation of ModelObject to a mat4
func (m *ModelObject) SetRotation(rot mgl32.Mat4) {
	m.Model.Rotation = rot
}

// GetModel : getter for ModelObject values
func (m ModelObject) GetModel() (Model, error) {
	if (m.Model != Model{}) {
		return m.Model, nil
	}
	return Model{}, errors.New("no ModelObject info")
}

// Scale : function used to scale the ModelObject and recalculate the centroid
func (m *ModelObject) Scale(scaleVec mgl32.Vec3) {
	m.Model.Scale = scaleVec
	m.centroid = CalculateCentroid(m.vertexValues.Vertices, m.Model.Scale)
}

func (m *ModelObject) Translate(translateVec mgl32.Vec3) {
	m.Model.Position = m.Model.Position.Add(translateVec)
	m.centroid = m.centroid.Add(translateVec)
}

func (m *ModelObject) SetVertexValues(vertices []float32, normals []float32, uvs []float32) {
	m.vertexValues.Vertices = vertices
	m.vertexValues.normals = normals
}

// Setup : function for initializing ModelObject
func (m *ModelObject) Setup(mat Material, mod Model, name string) error {
	m.name = name
	m.programInfo = ProgramInfo{}
	m.programInfo.Program = InitOpenGL(m.vertShader, m.fragShader)
	m.programInfo.attributes = Attributes{
		position: 0,
		normal:   1,
	}

	m.boundingBox = GetBoundingBox(m.vertexValues.Vertices)
	m.material = mat
	SetupAttributes(&m.programInfo)
	m.Scale(mod.Scale)
	m.boundingBox = ScaleBoundingBox(m.boundingBox, mod.Scale)
	m.Model.Position = mod.Position
	m.boundingBox = TranslateBoundingBox(m.boundingBox, mod.Position)
	m.Model.Rotation = mod.Rotation
	m.centroid = CalculateCentroid(m.vertexValues.Vertices, m.Model.Scale)
	m.buffers.Vao = CreateTriangleVAO(&m.programInfo, m.vertexValues.Vertices, m.vertexValues.normals, nil)

	return nil
}
