package geometry

import (
	"errors"

	"../shader"
	"../texture"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Plane : Primitive plane geometry struct
type Plane struct {
	name           string
	fragShader     string
	vertShader     string
	shaderType     string
	parent         string
	boundingBox    BoundingBox
	buffers        ObjectBuffers
	programInfo    ProgramInfo
	material       Material
	model          Model
	centroid       mgl32.Vec3
	vertexValues   VertexValues
	modelMatrix    mgl32.Mat4
	shaderVal      shader.Shader
	diffuseTexture *texture.Texture
	normalTexture  *texture.Texture
}

// SetShader : helper function for applying frag/vert shader
func (p *Plane) SetShader(vertShader string, fragShader string) error {

	if vertShader != "" && fragShader != "" {
		p.fragShader = fragShader
		p.vertShader = vertShader
		return nil
	}
	return errors.New("Error setting the shader, shader code must not be blank")
}

func (p *Plane) SetModelMatrix(mm mgl32.Mat4) {
	p.modelMatrix = mm
}

func (p Plane) GetModelMatrix() mgl32.Mat4 {
	return p.modelMatrix
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

func (p Plane) GetDetails() (string, string, string) {
	return p.name, p.shaderType, p.parent
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

func (p Plane) GetType() string {
	return "plane"
}

func (p Plane) GetBoundingBox() BoundingBox {
	return p.boundingBox
}

func (p Plane) GetShaderVal() shader.Shader {
	return p.shaderVal
}

func (p *Plane) SetShaderVal(s shader.Shader) {
	p.shaderVal = s
}

func (p Plane) GetDiffuseTexture() *texture.Texture {
	return p.diffuseTexture
}

func (p Plane) GetNormalTexture() *texture.Texture {
	return p.normalTexture
}

// Scale : function used to scale the cube and recalculate the centroid
func (p *Plane) Scale(scaleVec mgl32.Vec3) {
	p.model.Scale = scaleVec
	p.centroid = CalculateCentroid(p.vertexValues.Vertices, p.model.Scale)
}

func (p *Plane) Translate(translateVec mgl32.Vec3) {
	p.model.Position = p.model.Position.Add(translateVec)
	p.centroid = p.centroid.Add(translateVec)
}

// SetRotation : helper function for setting rotation of cube to a mat4
func (p *Plane) SetRotation(rot mgl32.Mat4) {
	p.model.Rotation = rot
}

func (p *Plane) SetParent(parent string) {
	p.parent = parent
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

	p.vertexValues.Vertices = []float32{
		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.5,
	}

	p.vertexValues.faces = []uint32{
		0, 2, 1, 2, 0, 3,
	}

	p.vertexValues.normals = []float32{
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
	}
	p.vertexValues.uvs = []float32{
		0.0, 0.0,
		5.0, 0.0,
		5.0, 5.0,
		0.0, 5.0,
	}

	p.name = name
	p.programInfo = ProgramInfo{}
	p.material = mat

	var shaderVals map[string]bool
	shaderVals = make(map[string]bool)

	if mat.ShaderType == 0 {
		shaderVals["aPosition"] = true
		bS := &shader.BasicShader{}
		bS.Setup()
		p.shaderVal = bS
		p.programInfo.Program = InitOpenGL(p.shaderVal.GetVertShader(), p.shaderVal.GetFragShader())
		p.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}

		SetupAttributesMap(&p.programInfo, shaderVals)

		p.buffers.Vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.Vertices, nil, nil, nil, nil, p.vertexValues.faces)

	} else if mat.ShaderType == 1 {
		shaderVals["aPosition"] = true
		shaderVals["aNormal"] = true
		shaderVals["diffuseVal"] = true
		shaderVals["ambientVal"] = true
		shaderVals["specularVal"] = true
		shaderVals["nVal"] = true
		shaderVals["uProjectionMatrix"] = true
		shaderVals["uViewMatrix"] = true
		shaderVals["uModelMatrix"] = true
		shaderVals["pointLights"] = true
		shaderVals["numLights"] = true
		shaderVals["cameraPosition"] = true

		bS := &shader.BlinnNoTexture{}
		bS.Setup()
		p.shaderVal = bS
		p.programInfo.Program = InitOpenGL(p.shaderVal.GetVertShader(), p.shaderVal.GetFragShader())
		p.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}

		SetupAttributesMap(&p.programInfo, shaderVals)

		p.buffers.Vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.Vertices, p.vertexValues.normals, nil, nil, nil, p.vertexValues.faces)

	} else if mat.ShaderType == 2 {
		p.programInfo.Program = InitOpenGL(p.vertShader, p.fragShader)
		p.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
			uv:       2,
		}

	} else if mat.ShaderType == 3 {
		shaderVals["aPosition"] = true
		shaderVals["aNormal"] = true
		shaderVals["aUV"] = true
		shaderVals["diffuseVal"] = true
		shaderVals["ambientVal"] = true
		shaderVals["specularVal"] = true
		shaderVals["nVal"] = true
		shaderVals["uProjectionMatrix"] = true
		shaderVals["uViewMatrix"] = true
		shaderVals["uModelMatrix"] = true
		shaderVals["pointLights"] = true
		shaderVals["numLights"] = true
		shaderVals["cameraPosition"] = true
		shaderVals["uDiffuseTexture"] = true

		bS := &shader.BlinnDiffuseTexture{}
		bS.Setup()
		p.shaderVal = bS
		p.programInfo.Program = InitOpenGL(p.shaderVal.GetVertShader(), p.shaderVal.GetFragShader())
		p.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
			uv:       2,
		}
		texture0, err := texture.NewTextureFromFile("../Editor/"+p.material.DiffuseTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}
		p.diffuseTexture = texture0

		SetupAttributesMap(&p.programInfo, shaderVals)
		p.buffers.Vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.Vertices, p.vertexValues.normals, p.vertexValues.uvs, nil, nil, p.vertexValues.faces)

	} else if mat.ShaderType == 4 {
		shaderVals["aPosition"] = true
		shaderVals["aNormal"] = true
		shaderVals["aUV"] = true
		shaderVals["diffuseVal"] = true
		shaderVals["ambientVal"] = true
		shaderVals["specularVal"] = true
		shaderVals["nVal"] = true
		shaderVals["uProjectionMatrix"] = true
		shaderVals["uViewMatrix"] = true
		shaderVals["uModelMatrix"] = true
		shaderVals["pointLights"] = true
		shaderVals["numLights"] = true
		shaderVals["cameraPosition"] = true
		shaderVals["uDiffuseTexture"] = true

		//calculate tangents and bitangents
		tangents, bitangents := CalculateBitangents(p.vertexValues.Vertices, p.vertexValues.uvs)

		bS := &shader.BlinnDiffuseAndNormal{}
		bS.Setup()
		p.shaderVal = bS
		p.programInfo.Program = InitOpenGL(p.shaderVal.GetVertShader(), p.shaderVal.GetFragShader())
		p.programInfo.attributes = Attributes{
			position:  0,
			normal:    1,
			uv:        2,
			tangent:   3,
			bitangent: 4,
		}
		//load diffuse texture
		texture0, err := texture.NewTextureFromFile("../Editor/"+p.material.DiffuseTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}
		//load normal texture
		texture1, err := texture.NewTextureFromFile("../Editor/"+p.material.NormalTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}

		p.diffuseTexture = texture0
		p.normalTexture = texture1

		SetupAttributesMap(&p.programInfo, shaderVals)
		p.buffers.Vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.Vertices, p.vertexValues.normals, p.vertexValues.uvs, tangents, bitangents, p.vertexValues.faces)

	}

	p.boundingBox = GetBoundingBox(p.vertexValues.Vertices)
	p.Scale(mod.Scale)
	p.boundingBox = ScaleBoundingBox(p.boundingBox, mod.Scale)
	p.model.Position = mod.Position
	p.boundingBox = TranslateBoundingBox(p.boundingBox, mod.Position)
	p.model.Rotation = mod.Rotation
	p.centroid = CalculateCentroid(p.vertexValues.Vertices, p.model.Scale)

	return nil
}
