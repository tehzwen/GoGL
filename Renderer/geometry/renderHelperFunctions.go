package geometry

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// InitOpenGL : initializes OpenGL and returns an intiialized Program.
func InitOpenGL(vertexShaderSource, fragmentShaderSource, geometryShaderSource string) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	var compiled int32
	/*
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)*/

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &compiled)
	if compiled == gl.FALSE {
		panic(compiled)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &compiled)
	if compiled == gl.FALSE {
		panic(compiled)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)

	//try and check for a geometry shader here
	if geometryShaderSource != "" {
		geoShader, err := compileShader(geometryShaderSource, gl.GEOMETRY_SHADER)
		if err != nil {
			panic(err)
		}

		gl.GetShaderiv(geoShader, gl.COMPILE_STATUS, &compiled)
		if compiled == gl.FALSE {
			panic(compiled)
		}

		gl.AttachShader(prog, geoShader)
		gl.LinkProgram(prog)
	} else {
		gl.LinkProgram(prog)
	}

	var status int32

	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		panic(status)
	}

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

func SetupAttributesMap(p *ProgramInfo, m map[string]bool) {
	//fmt.Println(p)
	//PrintActiveAttribs(p)

	//do checks to see what values need to be attributed
	if m["aPosition"] {
		(*p).attributes.vertexPosition = gl.GetAttribLocation((*p).Program, gl.Str("aPosition\x00"))
		if (*p).attributes.vertexPosition == -1 {
			fmt.Printf("ERROR: aPosition not found in shader\n")
		}
	}
	if m["aNormal"] {
		(*p).attributes.vertexNormal = gl.GetAttribLocation((*p).Program, gl.Str("aNormal\x00"))
		if (*p).attributes.vertexNormal == -1 {
			fmt.Printf("ERROR: aNormal not found in shader\n")
		}
	}
	if m["aUV"] {
		(*p).attributes.vertexUV = gl.GetAttribLocation((*p).Program, gl.Str("aUV\x00"))
		if (*p).attributes.vertexUV == -1 {
			fmt.Printf("ERROR: aUV not found in shader\n")
		}
	}
	if m["diffuseVal"] {
		(*p).UniformLocations.DiffuseVal = gl.GetUniformLocation((*p).Program, gl.Str("diffuseVal\x00"))
		if (*p).UniformLocations.DiffuseVal == -1 {
			fmt.Printf("ERROR: diffuseVal not found in shader\n")
		}
	}
	if m["ambientVal"] {
		(*p).UniformLocations.AmbientVal = gl.GetUniformLocation((*p).Program, gl.Str("ambientVal\x00"))
		if (*p).UniformLocations.AmbientVal == -1 {
			fmt.Printf("ERROR: ambientVal not found in shader\n")
		}

	}
	if m["specularVal"] {
		(*p).UniformLocations.SpecularVal = gl.GetUniformLocation((*p).Program, gl.Str("specularVal\x00"))
		if (*p).UniformLocations.SpecularVal == -1 {
			fmt.Printf("ERROR: specularVal not found in shader\n")
		}
	}
	if m["nVal"] {
		(*p).UniformLocations.NVal = gl.GetUniformLocation((*p).Program, gl.Str("nVal\x00"))
		if (*p).UniformLocations.NVal == -1 {
			fmt.Printf("ERROR: nVal not found in shader\n")
		}
	}
	if m["uProjectionMatrix"] {
		(*p).UniformLocations.Projection = gl.GetUniformLocation((*p).Program, gl.Str("uProjectionMatrix\x00"))
		if (*p).UniformLocations.Projection == -1 {
			fmt.Printf("ERROR: uProjectionMatrix not found in shader\n")
		}
	}
	if m["uViewMatrix"] {
		(*p).UniformLocations.View = gl.GetUniformLocation((*p).Program, gl.Str("uViewMatrix\x00"))
		if (*p).UniformLocations.DiffuseVal == -1 {
			fmt.Printf("ERROR: uViewMatrix not found in shader\n")
		}
	}
	if m["uModelMatrix"] {
		(*p).UniformLocations.Model = gl.GetUniformLocation((*p).Program, gl.Str("uModelMatrix\x00"))
		if (*p).UniformLocations.Model == -1 {
			fmt.Printf("ERROR: uModelMatrix not found in shader\n")
		}
	}
	if m["lightPositions"] {
		(*p).UniformLocations.LightPositions = gl.GetUniformLocation((*p).Program, gl.Str("lightPositions\x00"))
		if (*p).UniformLocations.LightPositions == -1 {
			fmt.Printf("ERROR: lightPositions not found in shader\n")
		}
	}
	if m["lightColours"] {
		(*p).UniformLocations.LightColours = gl.GetUniformLocation((*p).Program, gl.Str("lightColours\x00"))
		if (*p).UniformLocations.LightColours == -1 {
			fmt.Printf("ERROR: lightColours not found in shader\n")
		}
	}
	if m["lightStrengths"] {
		(*p).UniformLocations.LightStrengths = gl.GetUniformLocation((*p).Program, gl.Str("lightStrengths\x00"))
		if (*p).UniformLocations.LightStrengths == -1 {
			fmt.Printf("ERROR: lightStrengths not found in shader\n")
		}
	}
	if m["numLights"] {
		(*p).UniformLocations.NumLights = gl.GetUniformLocation((*p).Program, gl.Str("numLights\x00"))
		if (*p).UniformLocations.NumLights == -1 {
			fmt.Printf("ERROR: numLights not found in shader\n")
		}
	}
	if m["cameraPosition"] {
		(*p).UniformLocations.CameraPosition = gl.GetUniformLocation((*p).Program, gl.Str("cameraPosition\x00"))
		if (*p).UniformLocations.CameraPosition == -1 {
			fmt.Printf("ERROR: cameraPosition not found in shader\n")
		}
	}
	if m["uDiffuseTexture"] {
		(*p).UniformLocations.DiffuseTexture = gl.GetUniformLocation((*p).Program, gl.Str("uDiffuseTexture\x00"))
		if (*p).UniformLocations.DiffuseTexture == -1 {
			fmt.Printf("ERROR: uDiffuseTexture not found in shader\n")
		}
	}
	if m["depthMap"] {
		(*p).UniformLocations.DepthMap = gl.GetUniformLocation((*p).Program, gl.Str("depthMap\x00"))
		if (*p).UniformLocations.DepthMap == -1 {
			fmt.Printf("ERROR: Depth map not found in shader\n")
		}
	}
	if m["shadowMatrices"] {
		//fmt.Println("SHADOWS")
		(*p).UniformLocations.ShadowMatrices = gl.GetUniformLocation((*p).Program, gl.Str("shadowMatrices\x00"))
		if (*p).UniformLocations.ShadowMatrices == -1 {
			fmt.Printf("ERROR: Shadow matrices not found in shader\n")
		}
	}
	if m["lightPos"] {
		(*p).UniformLocations.LightPos = gl.GetUniformLocation((*p).Program, gl.Str("lightPos\x00"))
		if (*p).UniformLocations.LightPos == -1 {
			fmt.Printf("ERROR: lightPos not found in shader\n")
		}
	}

	// if m["pointLights"] {
	// 	(*p).UniformLocations.PointLights = gl.GetUniformLocation((*p).Program, gl.Str("pointLights\x00"))
	// 	if (*p).UniformLocations.PointLights == -1 {
	// 		fmt.Printf("ERROR: pointLights not found in shader\n")
	// 	}
	// }
}

func PrintActiveAttribs(p *ProgramInfo) {
	var count int32
	var length int32
	var bufsize int32 = 16
	var size int32
	var typeVal uint32
	var name uint8

	var i int32 = 0
	gl.GetProgramiv(p.Program, gl.ACTIVE_ATTRIBUTES, &count)

	for ; i < count; i++ {
		gl.GetActiveAttrib(p.Program, uint32(i), bufsize, &length, &size, &typeVal, &name)
		fmt.Printf("Attribute #%d Type: %d Name: %s\n", i, typeVal, string(name))
	}

	for i = 0; i < count; i++ {
		gl.GetActiveUniform(p.Program, uint32(i), bufsize, &length, &size, &typeVal, &name)
		fmt.Printf("Uniform #%d Type: %d Name: %s\n", i, typeVal, string(name))
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
func CreateTriangleVAO(programInfo *ProgramInfo, vertices []float32, normals []float32, uvs []float32, tangents []float32, bitangents []float32, indices []uint32) uint32 {

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

	if tangents != nil {
		var tangentBuffer uint32
		gl.GenBuffers(1, &tangentBuffer)

		gl.BindBuffer(gl.ARRAY_BUFFER, tangentBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(tangents)*4, gl.Ptr(tangents), gl.STATIC_DRAW)
		gl.VertexAttribPointer((*programInfo).attributes.tangent, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray((*programInfo).attributes.tangent)
	}

	if bitangents != nil {
		var bitangentBuffer uint32
		gl.GenBuffers(1, &bitangentBuffer)

		gl.BindBuffer(gl.ARRAY_BUFFER, bitangentBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(bitangents)*4, gl.Ptr(bitangents), gl.STATIC_DRAW)
		gl.VertexAttribPointer((*programInfo).attributes.bitangent, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray((*programInfo).attributes.bitangent)
	}

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray((*programInfo).attributes.position)

	return VAO
}
