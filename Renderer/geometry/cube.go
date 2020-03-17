package geometry

import (
	"errors"

	"../shader"
	"../texture"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Cube : Primitive cube geometry struct
type Cube struct {
	name              string
	fragShader        string
	vertShader        string
	shaderType        string
	parent            string
	boundingBox       BoundingBox
	buffers           ObjectBuffers
	programInfo       ProgramInfo
	material          Material
	model             Model
	centroid          mgl32.Vec3
	vertexValues      VertexValues
	modelMatrix       mgl32.Mat4
	shaderVal         shader.Shader
	diffuseTexture    *texture.Texture
	normalTexture     *texture.Texture
	onCollide         collisionFunction
	velocity          mgl32.Vec3
	shadowProgramInfo ProgramInfo
	shadowShaderVal   shader.Shader
	shadowBuffers     ObjectBuffers
}

func (c *Cube) SetBoundingBox(b BoundingBox) {
	c.boundingBox = b
}

func (c *Cube) SetForce(v mgl32.Vec3) {
	c.velocity = v
}

func (c *Cube) GetForce() mgl32.Vec3 {
	return c.velocity
}

func (c *Cube) AddForce(v mgl32.Vec3) {
	c.velocity = c.velocity.Add(v)
}

func (c *Cube) OnCollide(box BoundingBox) {
	c.onCollide(box)
}

func (c *Cube) SetOnCollide(colFunc collisionFunction) {
	c.onCollide = colFunc
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

func (c *Cube) SetModelMatrix(mm mgl32.Mat4) {
	c.modelMatrix = mm
}

func (c Cube) GetModelMatrix() (mgl32.Mat4, error) {
	if (c.modelMatrix != mgl32.Mat4{}) {
		return c.modelMatrix, nil
	}
	return mgl32.Mat4{}, errors.New("No matrix yet")
}

// GetProgramInfo : getter for programinfo
func (c Cube) GetProgramInfo() (ProgramInfo, error) {
	if (c.programInfo != ProgramInfo{}) {
		return c.programInfo, nil
	}
	return ProgramInfo{}, errors.New("No program info")
}

func (c Cube) GetShadowProgramInfo() (ProgramInfo, error) {
	if (c.shadowProgramInfo != ProgramInfo{}) {
		return c.shadowProgramInfo, nil
	}
	return ProgramInfo{}, errors.New("No shadow program info")
}

func (c Cube) GetShadowBuffers() ObjectBuffers {
	return c.shadowBuffers
}

// GetMaterial : getter for material
func (c Cube) GetMaterial() Material {
	return c.material
}

func (c Cube) GetDetails() (string, string, string) {
	return c.name, c.shaderType, c.parent
}

// GetVertices : getter for vertexValues
func (c Cube) GetVertices() VertexValues {
	return c.vertexValues
}

func (c Cube) GetShaderVal() shader.Shader {
	return c.shaderVal
}

func (c *Cube) SetShaderVal(s shader.Shader) {
	c.shaderVal = s
}

// GetCentroid : getter for centroid
func (c Cube) GetCentroid() mgl32.Vec3 {
	return c.centroid
}

func (c Cube) GetDiffuseTexture() *texture.Texture {
	return c.diffuseTexture
}

func (c Cube) GetNormalTexture() *texture.Texture {
	return c.normalTexture
}

// GetBuffers : getter for buffers
func (c Cube) GetBuffers() ObjectBuffers {
	return c.buffers
}

func (c Cube) GetBoundingBox() BoundingBox {
	return c.boundingBox
}

func (c Cube) GetType() string {
	return "cube"
}

// SetRotation : helper function for setting rotation of cube to a mat4
func (c *Cube) SetRotation(rot mgl32.Mat4) {
	c.model.Rotation = rot
}

func (c *Cube) SetParent(parent string) {
	c.parent = parent
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
	c.boundingBox = ScaleBoundingBox(c.boundingBox, scaleVec)
}

func (c *Cube) Translate(translateVec mgl32.Vec3) {
	c.model.Position = c.model.Position.Add(translateVec)
	c.centroid = c.centroid.Add(translateVec)
	c.boundingBox = TranslateBoundingBox(c.boundingBox, translateVec)
}

// Setup : function for initializing cube
func (c *Cube) Setup(mat Material, mod Model, name string, collide bool) error {
	c.vertexValues.Vertices = []float32{
		//front face
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.0, 0.0,

		//back face
		0.0, 0.0, 0.5,
		0.0, 0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.0, 0.5,

		//top face
		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.5,

		//bottom face
		0.0, 0.0, 0.5,
		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.0, 0.0, 0.0,

		//side face
		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.0,

		//side face
		0.0, 0.0, 0.5,
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
	}

	c.vertexValues.Faces = []uint32{
		//front face
		2, 0, 1, 3, 0, 2,
		//backface
		5, 4, 6, 6, 4, 7,
		//top face
		10, 9, 8, 10, 8, 11,
		//bottom face
		13, 12, 14, 14, 12, 15,
		//
		18, 16, 17, 18, 17, 19,

		22, 21, 20, 23, 21, 22,
	}

	c.vertexValues.Normals = []float32{
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
	c.vertexValues.Uvs = []float32{
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,

		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,

		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,

		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,

		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,

		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
	}

	c.name = name
	c.programInfo = ProgramInfo{}
	c.material = mat

	var shaderVals map[string]bool
	shaderVals = make(map[string]bool)

	if mat.ShaderType == 0 {
		shaderVals["aPosition"] = true
		bS := &shader.BasicShader{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader(), c.shaderVal.GetGeometryShader())
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}

		SetupAttributesMap(&c.programInfo, shaderVals)

		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, nil, nil, nil, nil, c.vertexValues.Faces)

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
		shaderVals["cameraPosition"] = true

		bS := &shader.BlinnNoTexture{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader(), c.shaderVal.GetGeometryShader())
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}

		SetupAttributesMap(&c.programInfo, shaderVals)

		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, c.vertexValues.Normals, nil, nil, nil, c.vertexValues.Faces)

	} else if mat.ShaderType == 2 {
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader(), c.shaderVal.GetGeometryShader())
		c.programInfo.attributes = Attributes{
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
		shaderVals["cameraPosition"] = true
		shaderVals["uDiffuseTexture"] = true
		shaderVals["pointLights"] = true

		bS := &shader.BlinnDiffuseTexture{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader(), c.shaderVal.GetGeometryShader())
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
			uv:       2,
		}
		texture0, err := texture.NewTextureFromFile("../Editor/materials/"+c.material.DiffuseTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}
		c.diffuseTexture = texture0

		SetupAttributesMap(&c.programInfo, shaderVals)
		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, c.vertexValues.Normals, c.vertexValues.Uvs, nil, nil, c.vertexValues.Faces)

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
		shaderVals["cameraPosition"] = true
		shaderVals["uDiffuseTexture"] = true
		shaderVals["uNormalTexture"] = true

		//calculate tangents and bitangents
		tangents, bitangents := CalculateBitangents(c.vertexValues.Vertices, c.vertexValues.Uvs)

		bS := &shader.BlinnDiffuseAndNormal{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader(), c.shaderVal.GetGeometryShader())
		c.programInfo.attributes = Attributes{
			position:  0,
			normal:    1,
			uv:        2,
			tangent:   3,
			bitangent: 4,
		}

		//load diffuse texture
		texture0, err := texture.NewTextureFromFile("../Editor/materials/"+c.material.DiffuseTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}
		//load normal texture
		texture1, err := texture.NewTextureFromFile("../Editor/materials/"+c.material.NormalTexture,
			gl.REPEAT, gl.REPEAT)

		if err != nil {
			panic(err)
		}

		c.diffuseTexture = texture0
		c.normalTexture = texture1

		SetupAttributesMap(&c.programInfo, shaderVals)
		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, c.vertexValues.Normals, c.vertexValues.Uvs, tangents, bitangents, c.vertexValues.Faces)

	}

	c.boundingBox = GetBoundingBox(c.vertexValues.Vertices)

	if collide {
		c.boundingBox.Collide = true
	} else {
		c.boundingBox.Collide = false
	}

	c.Scale(mod.Scale)
	c.boundingBox = ScaleBoundingBox(c.boundingBox, mod.Scale)
	c.model.Position = mod.Position
	c.boundingBox = TranslateBoundingBox(c.boundingBox, mod.Position)
	c.model.Rotation = mod.Rotation
	c.centroid = CalculateCentroid(c.vertexValues.Vertices, c.model.Scale)
	c.onCollide = func(box BoundingBox) {}

	return nil
}
