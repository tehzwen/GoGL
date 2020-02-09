package geometry

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"../shader"
	"../texture"
	"github.com/go-gl/mathgl/mgl32"
)

type collisionFunction func(BoundingBox)

// Geometry : interface for all objects that will be rendered
type Geometry interface {
	Setup(Material, Model, string, bool) error
	SetShader(string, string) error
	SetShaderVal(shader.Shader)
	GetProgramInfo() (ProgramInfo, error)
	GetShadowProgramInfo() (ProgramInfo, error)
	GetModel() (Model, error)
	GetType() string
	GetCentroid() mgl32.Vec3
	GetShaderVal() shader.Shader
	GetDiffuseTexture() *texture.Texture
	GetNormalTexture() *texture.Texture
	GetMaterial() Material
	GetBuffers() ObjectBuffers
	GetShadowBuffers() ObjectBuffers
	GetVertices() VertexValues
	GetModelMatrix() (mgl32.Mat4, error)
	SetRotation(mgl32.Mat4)
	SetModelMatrix(mgl32.Mat4)
	SetParent(string)
	Scale(mgl32.Vec3)
	Translate(mgl32.Vec3)
	GetDetails() (string, string, string)
	GetBoundingBox() BoundingBox
	SetBoundingBox(BoundingBox)
	SetOnCollide(collisionFunction)
	OnCollide(BoundingBox)
	AddForce(mgl32.Vec3)
	GetForce() mgl32.Vec3
	SetForce(mgl32.Vec3)
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

func (a *Attributes) SetPosition(pos uint32) {
	a.position = pos
}

func (a *Attributes) SetNormal(norm uint32) {
	a.normal = norm
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
	Normals  []float32
	Uvs      []float32
	Faces    []uint32
}

func (v *VertexValues) Serialize() string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)

	err := e.Encode(v)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (v *VertexValues) Deserialize(value string) error {
	by, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}

	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&v)
	return err
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
	DepthMap       int32
	ShadowMatrices int32
	LightPos       int32
}

// ProgramInfo : struct for holding program info (program, uniforms, attributes)
type ProgramInfo struct {
	Program          uint32
	UniformLocations Uniforms
	attributes       Attributes
	indexBuffer      uint32
}

func (p *ProgramInfo) SetAttributes(a Attributes) {
	p.attributes = a
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
	Min            mgl32.Vec3
	Max            mgl32.Vec3
	Collide        bool
	CollisionCount int
	CollisionBody  string
}
