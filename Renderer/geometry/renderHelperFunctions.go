package geometry

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// InitOpenGL : initializes OpenGL and returns an intiialized Program.
func InitOpenGL(vertexShaderSource string, fragmentShaderSource string) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	/*
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)*/

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// SetupAttributes : helper function for setting up attribute/uniforms
func SetupAttributes(p *ProgramInfo) {
	(*p).attributes.vertexPosition = gl.GetAttribLocation((*p).Program, gl.Str("aPosition\x00"))
	(*p).attributes.vertexNormal = gl.GetAttribLocation((*p).Program, gl.Str("aNormal\x00"))
	(*p).UniformLocations.DiffuseVal = gl.GetUniformLocation((*p).Program, gl.Str("diffuseVal\x00"))
	(*p).UniformLocations.AmbientVal = gl.GetUniformLocation((*p).Program, gl.Str("ambientVal\x00"))
	(*p).UniformLocations.SpecularVal = gl.GetUniformLocation((*p).Program, gl.Str("specularVal\x00"))
	(*p).UniformLocations.NVal = gl.GetUniformLocation((*p).Program, gl.Str("nVal\x00"))
	(*p).UniformLocations.Projection = gl.GetUniformLocation((*p).Program, gl.Str("uProjectionMatrix\x00"))
	(*p).UniformLocations.View = gl.GetUniformLocation((*p).Program, gl.Str("uViewMatrix\x00"))
	(*p).UniformLocations.Model = gl.GetUniformLocation((*p).Program, gl.Str("uModelMatrix\x00"))
	(*p).UniformLocations.LightPositions = gl.GetUniformLocation((*p).Program, gl.Str("lightPositions\x00"))
	(*p).UniformLocations.LightColours = gl.GetUniformLocation((*p).Program, gl.Str("lightColours\x00"))
	(*p).UniformLocations.LightStrengths = gl.GetUniformLocation((*p).Program, gl.Str("lightStrengths\x00"))
	(*p).UniformLocations.NumLights = gl.GetUniformLocation((*p).Program, gl.Str("numLights\x00"))
	(*p).UniformLocations.CameraPosition = gl.GetUniformLocation((*p).Program, gl.Str("cameraPosition\x00"))

	if (*p).attributes.vertexPosition == -1 ||
		(*p).attributes.vertexNormal == -1 ||
		(*p).UniformLocations.Projection == -1 ||
		(*p).UniformLocations.View == -1 ||
		(*p).UniformLocations.Model == -1 ||
		(*p).UniformLocations.CameraPosition == -1 ||
		(*p).UniformLocations.LightPositions == -1 ||
		(*p).UniformLocations.LightColours == -1 ||
		(*p).UniformLocations.LightStrengths == -1 ||
		(*p).UniformLocations.NumLights == -1 ||
		(*p).UniformLocations.DiffuseVal == -1 ||
		(*p).UniformLocations.AmbientVal == -1 ||
		(*p).UniformLocations.NVal == -1 ||
		(*p).UniformLocations.SpecularVal == -1 {
		fmt.Printf("ERROR: One or more of the uniforms or attributes cannot be found in the shader\n")
	}
}

// CalculateCentroid : helper function for calculating the centroid of the geometry
func CalculateCentroid(vertices []float32, currScale mgl32.Vec3) mgl32.Vec3 {
	var xTotal = float32(0.0)
	var yTotal = float32(0.0)
	var zTotal = float32(0.0)

	for i := 0; i < len(vertices); i += 3 {
		xTotal += vertices[i] * currScale[0]
		yTotal += vertices[i+1] * currScale[1]
		zTotal += vertices[i+2] * currScale[2]
	}

	xTotal /= float32(len(vertices) / 3)
	yTotal /= float32(len(vertices) / 3)
	zTotal /= float32(len(vertices) / 3)

	return mgl32.Vec3{xTotal, yTotal, zTotal}

}

// CreateTriangleVAO : helper function for creating vertex buffer values
func CreateTriangleVAO(programInfo *ProgramInfo, vertices []float32, normals []float32, uvs []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)
	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	if indices != nil {
		var EBO uint32
		gl.GenBuffers(1, &EBO)
		// copy indices into element buffer
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	}

	// position
	gl.VertexAttribPointer((*programInfo).attributes.position, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray((*programInfo).attributes.position)

	//normals
	if normals != nil {
		var normalBuffer uint32
		gl.GenBuffers(1, &normalBuffer)

		gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(normals)*4, gl.Ptr(normals), gl.STATIC_DRAW)

		gl.VertexAttribPointer((*programInfo).attributes.normal, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray((*programInfo).attributes.normal)
	}

	if uvs != nil {
		var uvBuffer uint32
		gl.GenBuffers(1, &uvBuffer)

		gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(uvs)*4, gl.Ptr(uvs), gl.STATIC_DRAW)
		gl.VertexAttribPointer((*programInfo).attributes.uv, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray((*programInfo).attributes.uv)
	}

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray((*programInfo).attributes.position)

	return VAO
}
