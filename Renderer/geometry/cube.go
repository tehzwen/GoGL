package geometry

import (
	"errors"
	"fmt"

	"../shader"
	"../texture"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Cube : Primitive cube geometry struct
type Cube struct {
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

func (c Cube) GetModelMatrix() mgl32.Mat4 {
	return c.modelMatrix
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
}

func (c *Cube) Translate(translateVec mgl32.Vec3) {
	c.model.Position = c.model.Position.Add(translateVec)
	c.centroid = c.centroid.Add(translateVec)
}

// Setup : function for initializing cube
func (c *Cube) Setup(mat Material, mod Model, name string) error {
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

	c.vertexValues.faces = []uint32{
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
	c.vertexValues.uvs = []float32{
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

	if mat.ShaderType == 0 {
		bS := &shader.BasicShader{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader())
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}
		c.programInfo.attributes.vertexPosition = gl.GetAttribLocation(c.programInfo.Program, gl.Str("vertexPosition\x00"))

		if c.programInfo.attributes.vertexPosition == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}

		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, nil, nil, c.vertexValues.faces)

	} else if mat.ShaderType == 1 {
		bS := &shader.BlinnNoTexture{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader())
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}

		c.programInfo.attributes.vertexPosition = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aPosition\x00"))
		c.programInfo.attributes.vertexNormal = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aNormal\x00"))
		c.programInfo.UniformLocations.DiffuseVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("diffuseVal\x00"))
		c.programInfo.UniformLocations.AmbientVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("ambientVal\x00"))
		c.programInfo.UniformLocations.SpecularVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("specularVal\x00"))
		c.programInfo.UniformLocations.NVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("nVal\x00"))
		c.programInfo.UniformLocations.Projection = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uProjectionMatrix\x00"))
		c.programInfo.UniformLocations.View = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uViewMatrix\x00"))
		c.programInfo.UniformLocations.Model = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uModelMatrix\x00"))
		c.programInfo.UniformLocations.LightPositions = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightPositions\x00"))
		c.programInfo.UniformLocations.LightColours = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightColours\x00"))
		c.programInfo.UniformLocations.LightStrengths = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightStrengths\x00"))
		c.programInfo.UniformLocations.NumLights = gl.GetUniformLocation(c.programInfo.Program, gl.Str("numLights\x00"))
		c.programInfo.UniformLocations.CameraPosition = gl.GetUniformLocation(c.programInfo.Program, gl.Str("cameraPosition\x00"))

		if c.programInfo.attributes.vertexPosition == -1 ||
			c.programInfo.attributes.vertexNormal == -1 ||
			c.programInfo.UniformLocations.Projection == -1 ||
			c.programInfo.UniformLocations.View == -1 ||
			c.programInfo.UniformLocations.Model == -1 ||
			c.programInfo.UniformLocations.CameraPosition == -1 ||
			c.programInfo.UniformLocations.LightPositions == -1 ||
			c.programInfo.UniformLocations.LightColours == -1 ||
			c.programInfo.UniformLocations.LightStrengths == -1 ||
			c.programInfo.UniformLocations.NumLights == -1 ||
			c.programInfo.UniformLocations.DiffuseVal == -1 ||
			c.programInfo.UniformLocations.AmbientVal == -1 ||
			c.programInfo.UniformLocations.NVal == -1 ||
			c.programInfo.UniformLocations.SpecularVal == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}

		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, c.vertexValues.normals, nil, c.vertexValues.faces)

	} else if mat.ShaderType == 2 {
		c.programInfo.Program = InitOpenGL(c.vertShader, c.fragShader)
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
			uv:       2,
		}

		SetupAttributes(&c.programInfo)

	} else if mat.ShaderType == 3 {
		bS := &shader.BlinnDiffuseTexture{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader())
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
			uv:       2,
		}
		texture0, err := texture.NewTextureFromFile("../Editor/"+c.material.DiffuseTexture,
			gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

		if err != nil {
			panic(err)
		}
		c.diffuseTexture = texture0

		c.programInfo.attributes.vertexPosition = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aPosition\x00"))
		c.programInfo.attributes.vertexNormal = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aNormal\x00"))
		c.programInfo.attributes.vertexUV = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aUV\x00"))
		c.programInfo.UniformLocations.DiffuseVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("diffuseVal\x00"))
		c.programInfo.UniformLocations.AmbientVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("ambientVal\x00"))
		c.programInfo.UniformLocations.SpecularVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("specularVal\x00"))
		c.programInfo.UniformLocations.NVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("nVal\x00"))
		c.programInfo.UniformLocations.Projection = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uProjectionMatrix\x00"))
		c.programInfo.UniformLocations.View = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uViewMatrix\x00"))
		c.programInfo.UniformLocations.Model = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uModelMatrix\x00"))
		c.programInfo.UniformLocations.LightPositions = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightPositions\x00"))
		c.programInfo.UniformLocations.LightColours = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightColours\x00"))
		c.programInfo.UniformLocations.LightStrengths = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightStrengths\x00"))
		c.programInfo.UniformLocations.NumLights = gl.GetUniformLocation(c.programInfo.Program, gl.Str("numLights\x00"))
		c.programInfo.UniformLocations.CameraPosition = gl.GetUniformLocation(c.programInfo.Program, gl.Str("cameraPosition\x00"))

		if c.programInfo.attributes.vertexPosition == -1 ||
			c.programInfo.attributes.vertexNormal == -1 ||
			c.programInfo.attributes.vertexUV == -1 ||
			c.programInfo.UniformLocations.Projection == -1 ||
			c.programInfo.UniformLocations.View == -1 ||
			c.programInfo.UniformLocations.Model == -1 ||
			c.programInfo.UniformLocations.CameraPosition == -1 ||
			c.programInfo.UniformLocations.LightPositions == -1 ||
			c.programInfo.UniformLocations.LightColours == -1 ||
			c.programInfo.UniformLocations.LightStrengths == -1 ||
			c.programInfo.UniformLocations.NumLights == -1 ||
			c.programInfo.UniformLocations.DiffuseVal == -1 ||
			c.programInfo.UniformLocations.AmbientVal == -1 ||
			c.programInfo.UniformLocations.NVal == -1 ||
			c.programInfo.UniformLocations.SpecularVal == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}
		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, c.vertexValues.normals, c.vertexValues.uvs, c.vertexValues.faces)

	} else if mat.ShaderType == 4 {
		bS := &shader.BlinnDiffuseTexture{}
		bS.Setup()
		c.shaderVal = bS
		c.programInfo.Program = InitOpenGL(c.shaderVal.GetVertShader(), c.shaderVal.GetFragShader())
		c.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
			uv:       2,
		}
		texture0, err := texture.NewTextureFromFile("../Editor/"+c.material.DiffuseTexture,
			gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

		if err != nil {
			panic(err)
		}
		c.diffuseTexture = texture0

		c.programInfo.attributes.vertexPosition = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aPosition\x00"))
		c.programInfo.attributes.vertexNormal = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aNormal\x00"))
		c.programInfo.attributes.vertexUV = gl.GetAttribLocation(c.programInfo.Program, gl.Str("aUV\x00"))
		c.programInfo.UniformLocations.DiffuseVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("diffuseVal\x00"))
		c.programInfo.UniformLocations.AmbientVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("ambientVal\x00"))
		c.programInfo.UniformLocations.SpecularVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("specularVal\x00"))
		c.programInfo.UniformLocations.NVal = gl.GetUniformLocation(c.programInfo.Program, gl.Str("nVal\x00"))
		c.programInfo.UniformLocations.Projection = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uProjectionMatrix\x00"))
		c.programInfo.UniformLocations.View = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uViewMatrix\x00"))
		c.programInfo.UniformLocations.Model = gl.GetUniformLocation(c.programInfo.Program, gl.Str("uModelMatrix\x00"))
		c.programInfo.UniformLocations.LightPositions = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightPositions\x00"))
		c.programInfo.UniformLocations.LightColours = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightColours\x00"))
		c.programInfo.UniformLocations.LightStrengths = gl.GetUniformLocation(c.programInfo.Program, gl.Str("lightStrengths\x00"))
		c.programInfo.UniformLocations.NumLights = gl.GetUniformLocation(c.programInfo.Program, gl.Str("numLights\x00"))
		c.programInfo.UniformLocations.CameraPosition = gl.GetUniformLocation(c.programInfo.Program, gl.Str("cameraPosition\x00"))

		if c.programInfo.attributes.vertexPosition == -1 ||
			c.programInfo.attributes.vertexNormal == -1 ||
			c.programInfo.attributes.vertexUV == -1 ||
			c.programInfo.UniformLocations.Projection == -1 ||
			c.programInfo.UniformLocations.View == -1 ||
			c.programInfo.UniformLocations.Model == -1 ||
			c.programInfo.UniformLocations.CameraPosition == -1 ||
			c.programInfo.UniformLocations.LightPositions == -1 ||
			c.programInfo.UniformLocations.LightColours == -1 ||
			c.programInfo.UniformLocations.LightStrengths == -1 ||
			c.programInfo.UniformLocations.NumLights == -1 ||
			c.programInfo.UniformLocations.DiffuseVal == -1 ||
			c.programInfo.UniformLocations.AmbientVal == -1 ||
			c.programInfo.UniformLocations.NVal == -1 ||
			c.programInfo.UniformLocations.SpecularVal == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}
		c.buffers.Vao = CreateTriangleVAO(&c.programInfo, c.vertexValues.Vertices, c.vertexValues.normals, c.vertexValues.uvs, c.vertexValues.faces)

	}

	c.boundingBox = GetBoundingBox(c.vertexValues.Vertices)
	c.Scale(mod.Scale)
	c.boundingBox = ScaleBoundingBox(c.boundingBox, mod.Scale)
	c.model.Position = mod.Position
	c.boundingBox = TranslateBoundingBox(c.boundingBox, mod.Position)
	c.model.Rotation = mod.Rotation
	c.centroid = CalculateCentroid(c.vertexValues.Vertices, c.model.Scale)

	return nil
}
