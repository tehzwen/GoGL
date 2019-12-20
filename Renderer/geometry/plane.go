package geometry

import (
	"errors"
	"fmt"

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

	if mat.ShaderType == 0 {
		bS := &shader.BasicShader{}
		bS.Setup()
		p.shaderVal = bS
		p.programInfo.Program = InitOpenGL(p.shaderVal.GetVertShader(), p.shaderVal.GetFragShader())
		p.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}
		p.programInfo.attributes.vertexPosition = gl.GetAttribLocation(p.programInfo.Program, gl.Str("vertexPosition\x00"))

		if p.programInfo.attributes.vertexPosition == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}

	} else if mat.ShaderType == 1 {
		bS := &shader.BlinnNoTexture{}
		bS.Setup()
		p.shaderVal = bS
		p.programInfo.Program = InitOpenGL(p.shaderVal.GetVertShader(), p.shaderVal.GetFragShader())
		p.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
		}

		p.programInfo.attributes.vertexPosition = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aPosition\x00"))
		p.programInfo.attributes.vertexNormal = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aNormal\x00"))
		p.programInfo.UniformLocations.DiffuseVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("diffuseVal\x00"))
		p.programInfo.UniformLocations.AmbientVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("ambientVal\x00"))
		p.programInfo.UniformLocations.SpecularVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("specularVal\x00"))
		p.programInfo.UniformLocations.NVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("nVal\x00"))
		p.programInfo.UniformLocations.Projection = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uProjectionMatrix\x00"))
		p.programInfo.UniformLocations.View = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uViewMatrix\x00"))
		p.programInfo.UniformLocations.Model = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uModelMatrix\x00"))
		p.programInfo.UniformLocations.LightPositions = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightPositions\x00"))
		p.programInfo.UniformLocations.LightColours = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightColours\x00"))
		p.programInfo.UniformLocations.LightStrengths = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightStrengths\x00"))
		p.programInfo.UniformLocations.NumLights = gl.GetUniformLocation(p.programInfo.Program, gl.Str("numLights\x00"))
		p.programInfo.UniformLocations.CameraPosition = gl.GetUniformLocation(p.programInfo.Program, gl.Str("cameraPosition\x00"))

		if p.programInfo.attributes.vertexPosition == -1 ||
			p.programInfo.attributes.vertexNormal == -1 ||
			p.programInfo.UniformLocations.Projection == -1 ||
			p.programInfo.UniformLocations.View == -1 ||
			p.programInfo.UniformLocations.Model == -1 ||
			p.programInfo.UniformLocations.CameraPosition == -1 ||
			p.programInfo.UniformLocations.LightPositions == -1 ||
			p.programInfo.UniformLocations.LightColours == -1 ||
			p.programInfo.UniformLocations.LightStrengths == -1 ||
			p.programInfo.UniformLocations.NumLights == -1 ||
			p.programInfo.UniformLocations.DiffuseVal == -1 ||
			p.programInfo.UniformLocations.AmbientVal == -1 ||
			p.programInfo.UniformLocations.NVal == -1 ||
			p.programInfo.UniformLocations.SpecularVal == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}

	} else if mat.ShaderType == 2 {
		p.programInfo.Program = InitOpenGL(p.vertShader, p.fragShader)
		p.programInfo.attributes = Attributes{
			position: 0,
			normal:   1,
			uv:       2,
		}

		SetupAttributes(&p.programInfo)

	} else if mat.ShaderType == 3 {
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
			gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

		if err != nil {
			panic(err)
		}
		p.diffuseTexture = texture0

		p.programInfo.attributes.vertexPosition = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aPosition\x00"))
		p.programInfo.attributes.vertexNormal = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aNormal\x00"))
		p.programInfo.attributes.vertexUV = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aUV\x00"))
		p.programInfo.UniformLocations.DiffuseVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("diffuseVal\x00"))
		p.programInfo.UniformLocations.AmbientVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("ambientVal\x00"))
		p.programInfo.UniformLocations.SpecularVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("specularVal\x00"))
		p.programInfo.UniformLocations.NVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("nVal\x00"))
		p.programInfo.UniformLocations.Projection = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uProjectionMatrix\x00"))
		p.programInfo.UniformLocations.View = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uViewMatrix\x00"))
		p.programInfo.UniformLocations.Model = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uModelMatrix\x00"))
		p.programInfo.UniformLocations.LightPositions = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightPositions\x00"))
		p.programInfo.UniformLocations.LightColours = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightColours\x00"))
		p.programInfo.UniformLocations.LightStrengths = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightStrengths\x00"))
		p.programInfo.UniformLocations.NumLights = gl.GetUniformLocation(p.programInfo.Program, gl.Str("numLights\x00"))
		p.programInfo.UniformLocations.CameraPosition = gl.GetUniformLocation(p.programInfo.Program, gl.Str("cameraPosition\x00"))

		if p.programInfo.attributes.vertexPosition == -1 ||
			p.programInfo.attributes.vertexNormal == -1 ||
			p.programInfo.attributes.vertexUV == -1 ||
			p.programInfo.UniformLocations.Projection == -1 ||
			p.programInfo.UniformLocations.View == -1 ||
			p.programInfo.UniformLocations.Model == -1 ||
			p.programInfo.UniformLocations.CameraPosition == -1 ||
			p.programInfo.UniformLocations.LightPositions == -1 ||
			p.programInfo.UniformLocations.LightColours == -1 ||
			p.programInfo.UniformLocations.LightStrengths == -1 ||
			p.programInfo.UniformLocations.NumLights == -1 ||
			p.programInfo.UniformLocations.DiffuseVal == -1 ||
			p.programInfo.UniformLocations.AmbientVal == -1 ||
			p.programInfo.UniformLocations.NVal == -1 ||
			p.programInfo.UniformLocations.SpecularVal == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}
		p.buffers.Vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.Vertices, p.vertexValues.normals, p.vertexValues.uvs, p.vertexValues.faces)

	} else if mat.ShaderType == 4 {
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
			gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

		if err != nil {
			panic(err)
		}
		p.diffuseTexture = texture0

		p.programInfo.attributes.vertexPosition = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aPosition\x00"))
		p.programInfo.attributes.vertexNormal = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aNormal\x00"))
		p.programInfo.attributes.vertexUV = gl.GetAttribLocation(p.programInfo.Program, gl.Str("aUV\x00"))
		p.programInfo.UniformLocations.DiffuseVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("diffuseVal\x00"))
		p.programInfo.UniformLocations.AmbientVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("ambientVal\x00"))
		p.programInfo.UniformLocations.SpecularVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("specularVal\x00"))
		p.programInfo.UniformLocations.NVal = gl.GetUniformLocation(p.programInfo.Program, gl.Str("nVal\x00"))
		p.programInfo.UniformLocations.Projection = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uProjectionMatrix\x00"))
		p.programInfo.UniformLocations.View = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uViewMatrix\x00"))
		p.programInfo.UniformLocations.Model = gl.GetUniformLocation(p.programInfo.Program, gl.Str("uModelMatrix\x00"))
		p.programInfo.UniformLocations.LightPositions = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightPositions\x00"))
		p.programInfo.UniformLocations.LightColours = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightColours\x00"))
		p.programInfo.UniformLocations.LightStrengths = gl.GetUniformLocation(p.programInfo.Program, gl.Str("lightStrengths\x00"))
		p.programInfo.UniformLocations.NumLights = gl.GetUniformLocation(p.programInfo.Program, gl.Str("numLights\x00"))
		p.programInfo.UniformLocations.CameraPosition = gl.GetUniformLocation(p.programInfo.Program, gl.Str("cameraPosition\x00"))

		if p.programInfo.attributes.vertexPosition == -1 ||
			p.programInfo.attributes.vertexNormal == -1 ||
			p.programInfo.attributes.vertexUV == -1 ||
			p.programInfo.UniformLocations.Projection == -1 ||
			p.programInfo.UniformLocations.View == -1 ||
			p.programInfo.UniformLocations.Model == -1 ||
			p.programInfo.UniformLocations.CameraPosition == -1 ||
			p.programInfo.UniformLocations.LightPositions == -1 ||
			p.programInfo.UniformLocations.LightColours == -1 ||
			p.programInfo.UniformLocations.LightStrengths == -1 ||
			p.programInfo.UniformLocations.NumLights == -1 ||
			p.programInfo.UniformLocations.DiffuseVal == -1 ||
			p.programInfo.UniformLocations.AmbientVal == -1 ||
			p.programInfo.UniformLocations.NVal == -1 ||
			p.programInfo.UniformLocations.SpecularVal == -1 {
			fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
		}
		p.buffers.Vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.Vertices, p.vertexValues.normals, p.vertexValues.uvs, p.vertexValues.faces)

	}

	p.boundingBox = GetBoundingBox(p.vertexValues.Vertices)
	p.Scale(mod.Scale)
	p.boundingBox = ScaleBoundingBox(p.boundingBox, mod.Scale)
	p.model.Position = mod.Position
	p.boundingBox = TranslateBoundingBox(p.boundingBox, mod.Position)
	p.model.Rotation = mod.Rotation
	p.centroid = CalculateCentroid(p.vertexValues.Vertices, p.model.Scale)
	p.buffers.Vao = CreateTriangleVAO(&p.programInfo, p.vertexValues.Vertices, p.vertexValues.normals, p.vertexValues.uvs, p.vertexValues.faces)

	// x := errors.New("Wrong")
	// return x
	return nil
}
