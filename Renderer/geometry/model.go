package geometry

import (
	"errors"

	"../shader"
	"../texture"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// ModelObject : Primitive ModelObject geometry struct
type ModelObject struct {
	name           string
	fragShader     string
	vertShader     string
	shaderType     string
	parent         string
	boundingBox    BoundingBox
	buffers        ObjectBuffers
	programInfo    ProgramInfo
	material       Material
	Model          Model
	centroid       mgl32.Vec3
	vertexValues   VertexValues
	modelMatrix    mgl32.Mat4
	shaderVal      shader.Shader
	diffuseTexture *texture.Texture
	normalTexture  *texture.Texture
	onCollide      collisionFunction
	velocity       mgl32.Vec3
}

func (m *ModelObject) SetBoundingBox(b BoundingBox) {
	m.boundingBox = b
}

func (m *ModelObject) SetForce(v mgl32.Vec3) {
	m.velocity = v
}

func (m *ModelObject) GetForce() mgl32.Vec3 {
	return m.velocity
}

func (m *ModelObject) AddForce(v mgl32.Vec3) {
	m.velocity.Add(v)
}

func (m *ModelObject) OnCollide(box BoundingBox) {
	m.onCollide(box)
}

func (m *ModelObject) SetOnCollide(colFunc collisionFunction) {
	m.onCollide = colFunc
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

func (m *ModelObject) SetModelMatrix(mm mgl32.Mat4) {
	m.modelMatrix = mm
}

func (m ModelObject) GetModelMatrix() (mgl32.Mat4, error) {
	if (m.modelMatrix != mgl32.Mat4{}) {
		return m.modelMatrix, nil
	}
	return mgl32.Mat4{}, errors.New("No matrix yet")
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

func (m ModelObject) GetDetails() (string, string, string) {
	return m.name, m.shaderType, m.parent
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

func (m ModelObject) GetDiffuseTexture() *texture.Texture {
	return m.diffuseTexture
}

func (m ModelObject) GetNormalTexture() *texture.Texture {
	return m.normalTexture
}

func (m ModelObject) GetShaderVal() shader.Shader {
	return m.shaderVal
}

func (m *ModelObject) SetShaderVal(s shader.Shader) {
	m.shaderVal = s
}

// SetRotation : helper function for setting rotation of ModelObject to a mat4
func (m *ModelObject) SetRotation(rot mgl32.Mat4) {
	m.Model.Rotation = rot
}

func (m *ModelObject) SetParent(parent string) {
	m.parent = parent
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
	m.boundingBox = TranslateBoundingBox(m.boundingBox, translateVec)
}

func (m *ModelObject) SetVertexValues(vertices []float32, normals []float32, uvs []float32, faces []uint32) {
	m.vertexValues.Vertices = vertices
	m.vertexValues.normals = normals
	m.vertexValues.uvs = uvs
	m.vertexValues.faces = faces
}

// Setup : function for initializing ModelObject
func (m *ModelObject) Setup(mat Material, mod Model, name string, collide bool) error {
	m.name = name
	m.material = mat
	m.programInfo = ProgramInfo{}

	var shaderVals map[string]bool
	shaderVals = make(map[string]bool)

	if mat.ShaderType == 0 {
		shaderVals["aPosition"] = true
		bS := &shader.BasicShader{}
		bS.Setup()
		m.shaderVal = bS
		m.programInfo.Program = InitOpenGL(m.shaderVal.GetVertShader(), m.shaderVal.GetFragShader())
		m.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}
		SetupAttributesMap(&m.programInfo, shaderVals)
		m.buffers.Vao = CreateTriangleVAO(&m.programInfo, m.vertexValues.Vertices, nil, nil, nil, nil, m.vertexValues.faces)
	} else if mat.ShaderType == 1 {
		shaderVals["aPosition"] = true
		shaderVals["aNormal"] = true
		shaderVals["diffuseVal"] = true
		shaderVals["ambientVal"] = true
		shaderVals["specularVal"] = true
		shaderVals["nVal"] = true
		shaderVals["Alpha"] = true
		shaderVals["uProjectionMatrix"] = true
		shaderVals["uViewMatrix"] = true
		shaderVals["uModelMatrix"] = true
		shaderVals["pointLights"] = true
		shaderVals["numLights"] = true
		shaderVals["cameraPosition"] = true

		bS := &shader.BlinnNoTexture{}
		bS.Setup()
		m.shaderVal = bS
		m.programInfo.Program = InitOpenGL(m.shaderVal.GetVertShader(), m.shaderVal.GetFragShader())
		m.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}

		SetupAttributesMap(&m.programInfo, shaderVals)
		m.buffers.Vao = CreateTriangleVAO(&m.programInfo, m.vertexValues.Vertices, m.vertexValues.normals, nil, nil, nil, m.vertexValues.faces)

	} else if mat.ShaderType == 2 {
		//not sure yet what TODO here
	} else if mat.ShaderType == 3 {
		shaderVals["aPosition"] = true
		shaderVals["aNormal"] = true
		shaderVals["aUV"] = true
		shaderVals["Alpha"] = true
		shaderVals["diffuseVal"] = true
		shaderVals["ambientVal"] = true
		shaderVals["specularVal"] = true
		shaderVals["nVal"] = true
		shaderVals["uProjectionMatrix"] = true
		shaderVals["uViewMatrix"] = true
		shaderVals["uModelMatrix"] = true
		shaderVals["numLights"] = true
		shaderVals["cameraPosition"] = true
		shaderVals["uDiffuseTexture"] = true
		shaderVals["pointLights"] = true

		if m.material.DiffuseTexture != "" {
			bS := &shader.BlinnDiffuseTexture{}
			bS.Setup()
			m.shaderVal = bS
			m.programInfo.Program = InitOpenGL(m.shaderVal.GetVertShader(), m.shaderVal.GetFragShader())
			m.programInfo.attributes = Attributes{
				position: 0,
				normal:   1,
				uv:       2,
			}
			texture0, err := texture.NewTextureFromFile("../Editor/models/"+m.material.DiffuseTexture,
				gl.REPEAT, gl.REPEAT)

			if err != nil {
				panic(err)
			}
			m.diffuseTexture = texture0
		} else {
			bS := &shader.BlinnNoTexture{}
			bS.Setup()
			m.shaderVal = bS
			m.programInfo.Program = InitOpenGL(m.shaderVal.GetVertShader(), m.shaderVal.GetFragShader())
			m.programInfo.attributes = Attributes{
				position: 0,
				normal:   1,
				uv:       2,
			}
		}

		SetupAttributesMap(&m.programInfo, shaderVals)

		//check if UVS or not
		if len(m.vertexValues.uvs) > 0 {
			m.buffers.Vao = CreateTriangleVAO(&m.programInfo, m.vertexValues.Vertices, m.vertexValues.normals, m.vertexValues.uvs, nil, nil, m.vertexValues.faces)
		} else {
			m.buffers.Vao = CreateTriangleVAO(&m.programInfo, m.vertexValues.Vertices, m.vertexValues.normals, nil, nil, nil, m.vertexValues.faces)
		}

	} else if mat.ShaderType == 4 {
		shaderVals["aPosition"] = true
		shaderVals["aNormal"] = true
		shaderVals["aUV"] = true
		shaderVals["diffuseVal"] = true
		shaderVals["ambientVal"] = true
		shaderVals["specularVal"] = true
		shaderVals["Alpha"] = true
		shaderVals["nVal"] = true
		shaderVals["uProjectionMatrix"] = true
		shaderVals["uViewMatrix"] = true
		shaderVals["uModelMatrix"] = true
		shaderVals["numLights"] = true
		shaderVals["pointLights"] = true
		shaderVals["cameraPosition"] = true
		shaderVals["uDiffuseTexture"] = true
		shaderVals["uNormalTexture"] = true

		tangents, bitangents := CalculateBitangents(m.vertexValues.Vertices, m.vertexValues.uvs)

		bS := &shader.BlinnDiffuseAndNormal{}
		bS.Setup()
		m.shaderVal = bS
		m.programInfo.Program = InitOpenGL(m.shaderVal.GetVertShader(), m.shaderVal.GetFragShader())
		m.programInfo.attributes = Attributes{
			position:  0,
			normal:    1,
			uv:        2,
			tangent:   3,
			bitangent: 4,
		}
		//load diffuse texture
		texture0, err := texture.NewTextureFromFile("../Editor/models/"+m.material.DiffuseTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}
		//load normal texture
		texture1, err := texture.NewTextureFromFile("../Editor/models/"+m.material.NormalTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}

		m.diffuseTexture = texture0
		m.normalTexture = texture1
		SetupAttributesMap(&m.programInfo, shaderVals)
		m.buffers.Vao = CreateTriangleVAO(&m.programInfo, m.vertexValues.Vertices, m.vertexValues.normals, m.vertexValues.uvs, tangents, bitangents, m.vertexValues.faces)

	}

	m.boundingBox = GetBoundingBox(m.vertexValues.Vertices)

	if collide {
		m.boundingBox.Collide = true
	} else {
		m.boundingBox.Collide = false
	}

	m.Scale(mod.Scale)
	m.boundingBox = ScaleBoundingBox(m.boundingBox, mod.Scale)
	m.Model.Position = mod.Position
	m.boundingBox = TranslateBoundingBox(m.boundingBox, mod.Position)
	m.Model.Rotation = mod.Rotation
	m.centroid = CalculateCentroid(m.vertexValues.Vertices, m.Model.Scale)
	m.onCollide = func(box BoundingBox) {}

	return nil
}
