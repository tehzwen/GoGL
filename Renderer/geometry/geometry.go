package geometry

import (
	"../shader"
	"../texture"
	"github.com/go-gl/mathgl/mgl32"
)

// Geometry : interface for all objects that will be rendered
type Geometry interface {
	Setup(Material, Model, string) error
	SetShader(string, string) error
	SetShaderVal(shader.Shader)
	GetProgramInfo() (ProgramInfo, error)
	GetModel() (Model, error)
	GetType() string
	GetCentroid() mgl32.Vec3
	GetShaderVal() shader.Shader
	GetDiffuseTexture() *texture.Texture
	GetNormalTexture() *texture.Texture
	GetMaterial() Material
	GetBuffers() ObjectBuffers
	GetVertices() VertexValues
	GetModelMatrix() (mgl32.Mat4, error)
	SetRotation(mgl32.Mat4)
	SetModelMatrix(mgl32.Mat4)
	SetParent(string)
	Scale(mgl32.Vec3)
	Translate(mgl32.Vec3)
	GetDetails() (string, string, string)
	GetBoundingBox() BoundingBox
}

// Attributes : struct for holding vertex attribute locations
type Attributes struct {
	position        uint32
	normal          uint32
	uv              uint32
	tangent         uint32
	bitangent       uint32
	vertexPosition  int32
	vertexNormal    int32
	vertexUV        int32
	vertexTangent   int32
	vertexBitangent int32
}

// ObjectBuffers : holds references to vertex buffers
type ObjectBuffers struct {
	Vao        uint32
	vbo        uint32
	attributes Attributes
}

// VertexValues : struct for holding vertex specific values
type VertexValues struct {
	Vertices []float32
	normals  []float32
	uvs      []float32
	faces    []uint32
}

// Uniforms : struct for holding all uniforms
type Uniforms struct {
	Projection     int32
	View           int32
	Model          int32
	NormalMatrix   int32
	DiffuseVal     int32
	AmbientVal     int32
	SpecularVal    int32
	NVal           int32
	Alpha          int32
	CameraPosition int32
	NumLights      int32
	LightPositions int32
	LightColours   int32
	LightStrengths int32
	DiffuseTexture int32
	PointLights    int32
}

// ProgramInfo : struct for holding program info (program, uniforms, attributes)
type ProgramInfo struct {
	Program          uint32
	UniformLocations Uniforms
	attributes       Attributes
	indexBuffer      uint32
}

// Material : struct for holding material info
type Material struct {
	Diffuse        []float32 `json:"diffuse"`
	Ambient        []float32 `json:"ambient"`
	Specular       []float32 `json:"specular"`
	N              float32   `json:"n"`
	ShaderType     int       `json:"shaderType"`
	Alpha          float32
	DiffuseTexture string
	NormalTexture  string
}

// Model : struct for holding model info
type Model struct {
	Position mgl32.Vec3
	Rotation mgl32.Mat4
	Scale    mgl32.Vec3
}

type BoundingBox struct {
	Min mgl32.Vec3
	Max mgl32.Vec3
}
