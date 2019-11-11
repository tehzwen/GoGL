package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 500
	height = 500

	vertexShaderSource = `
		#version 410
		in vec3 aPosition;
		in vec3 aNormal;

		out vec3 oNormal;
		out vec3 normalInterp;
		out vec3 oFragPosition;
		
		uniform mat4 uProjectionMatrix;
		uniform mat4 uViewMatrix;
		uniform mat4 uModelMatrix;

		void main() {
			mat4 normalMatrix = transpose(inverse(uModelMatrix));
			oNormal = normalize((uModelMatrix * vec4(aNormal, 1.0)).xyz);
			normalInterp = vec3(normalMatrix * vec4(aNormal, 0.0));
			oFragPosition = (uModelMatrix * vec4(aPosition, 1.0)).xyz;
			gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix * vec4(aPosition, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		#define MAX_LIGHTS 128

		in vec3 oFragPosition;
		in vec3 normalInterp;
		in vec3 oNormal;

		uniform vec3 cameraPosition;
		uniform vec3 diffuseVal;
		uniform vec3 ambientVal;
		uniform vec3 specularVal;
		uniform float nVal;
		uniform int numLights;
		uniform vec3 lightPositions[MAX_LIGHTS];
		uniform vec3 lightColours[MAX_LIGHTS];
		uniform float lightStrengths[MAX_LIGHTS];

		out vec4 frag_colour;

		void main() {
			vec3 diffuse;
			vec3 ambient;
			vec3 specular;
			vec3 normal = normalize(normalInterp);

			for (int i = 0; i < numLights; i++) {
				vec3 lightDirection = normalize(lightPositions[i] - oFragPosition);

				//ambient
				ambient += (ambientVal * lightColours[i]) * lightStrengths[i];

				//diffuse
				float NdotL = max(dot(lightDirection, normal), 0.0);
				diffuse += ((diffuseVal * lightColours[i]) * NdotL * lightStrengths[i]);

				//specular
				vec3 nCameraPosition = normalize(cameraPosition); // Normalize the camera position
                vec3 V = normalize(nCameraPosition - oFragPosition);
				vec3 H = normalize(V + lightDirection); // H = V + L normalized
				
				if (NdotL > 0.0f)
				{
					float NDotH = max(dot(normal, H), 0.0);
                    float NHPow = pow(NDotH, nVal); // (N dot H)^n
                    specular += ((specularVal * lightColours[i]) * NHPow);
				}
			}
			//frag_colour = vec4(testNormal, 1.0);
			frag_colour = vec4(diffuse + ambient + specular, 1.0);
		}
	` + "\x00"
)

var angle = 0.0

var (
	triangle = []float32{
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.0, 0.0,

		0.0, 0.0, 0.5,
		0.0, 0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.0, 0.5,

		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.5,

		0.0, 0.0, 0.5,
		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.0, 0.0, 0.0,

		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.0,

		0.0, 0.0, 0.5,
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
	}
)

var triangleNormals = []float32{
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

var triangleFaces = []uint32{
	0, 1, 2, 0, 2, 3,
	4, 5, 6, 4, 6, 7,
	8, 9, 10, 8, 10, 11,
	12, 13, 14, 12, 14, 15,
	16, 17, 18, 17, 18, 19,
	20, 21, 22, 21, 22, 23,
}

func main() {
	runtime.LockOSThread()

	var state = State{
		vertShader: vertexShaderSource,
		fragShader: fragmentShaderSource,
		camera: Camera{
			position: mgl32.Vec3{1.5, 0.5, -2.5},
			center:   mgl32.Vec3{0.5, 0.0, 0.0},
			up:       mgl32.Vec3{0.0, 1.0, 0.0},
		},
		lights: []Light{
			Light{
				colour:   []float32{1.0, 1.0, 1.0},
				strength: 0.50,
				position: []float32{1.0, 0.0, -2.0},
			},
		},
	}

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL(state.vertShader, state.fragShader)

	var testObject = Object{
		name:       "zach",
		fragShader: "1234",
		vertShader: "123456",
		programInfo: ProgramInfo{
			program: program,
			attributes: Attributes{
				position: 0,
				normal:   1,
			},
		},
		vertices: triangle,
		material: Material{
			diffuse:  []float32{0.6, 0.2, 0.6},
			ambient:  []float32{0.1, 0.1, 0.1},
			specular: []float32{0.8, 0.8, 0.8},
			n:        10,
		},
		model: Model{
			position: mgl32.Vec3{1.0, 0.0, 0.0},
			rotation: mgl32.Mat4{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0},
		},
	}

	SetupAttributes(&testObject.programInfo)

	var objectList = []Object{}
	vao := CreateTriangleVAO(&testObject.programInfo, triangle, triangleNormals, triangleFaces)
	//normalBuffer := InitNormalAttribute(&testObject.programInfo, triangleNormals)
	centroid := CalculateCentroid(testObject.vertices)

	testObject.centroid = centroid
	testObject.buffers.vao = vao
	objectList = append(objectList, testObject)

	state.objects = objectList

	then := 0.0

	for !window.ShouldClose() {
		now := glfw.GetTime()
		deltaTime := now - then
		then = now
		angle += 0.5 * deltaTime

		state.objects[0].model.rotation = mgl32.HomogRotate3D(float32(angle), state.objects[0].model.position)
		draw(window, state)
	}
}

func draw(window *glfw.Window, state State) {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Disable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	glfw.PollEvents()

	//create light arrays
	lightPositionArray := []float32{}
	lightColorArray := []float32{}
	lightStrengthArray := []float32{}

	for t := 0; t < len(state.lights); t++ {
		//get all values in the arrays
		for u := 0; u < len(state.lights[t].position); u++ {
			lightPositionArray = append(lightPositionArray, state.lights[t].position[u])
			lightColorArray = append(lightColorArray, state.lights[t].colour[u])
		}
		lightStrengthArray = append(lightStrengthArray, state.lights[t].strength)
	}

	for i := 0; i < len(state.objects); i++ {

		gl.UseProgram(state.objects[i].programInfo.program)

		var fovy = float32(60 * math.Pi / 180)
		var aspect = float32(width / height)
		var near = float32(0.1)
		var far = float32(100.0)

		//projection matrix
		projection := mgl32.Perspective(fovy, aspect, near, far)
		gl.UniformMatrix4fv(state.objects[i].programInfo.uniformLocations.projection, 1, false, &projection[0])

		//view matrix
		viewMatrix := mgl32.LookAtV(state.camera.position, state.camera.center, state.camera.up)
		gl.UniformMatrix4fv(state.objects[i].programInfo.uniformLocations.view, 1, false, &viewMatrix[0])

		//model matrix
		modelMatrix := mgl32.Ident4()
		positionMat := mgl32.Translate3D(state.objects[i].model.position[0], state.objects[i].model.position[1], state.objects[i].model.position[2])
		modelMatrix = modelMatrix.Mul4(positionMat)
		centroidMat := mgl32.Translate3D(state.objects[i].centroid[0], state.objects[i].centroid[1], state.objects[i].centroid[2])
		modelMatrix = modelMatrix.Mul4(centroidMat)
		modelMatrix = modelMatrix.Mul4(state.objects[i].model.rotation)
		negCent := mgl32.Translate3D(-state.objects[i].centroid[0], -state.objects[i].centroid[1], -state.objects[i].centroid[2])
		modelMatrix = modelMatrix.Mul4(negCent)

		gl.UniformMatrix4fv(state.objects[i].programInfo.uniformLocations.model, 1, false, &modelMatrix[0])

		gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.diffuseVal, 1, &state.objects[i].material.diffuse[0])
		gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.ambientVal, 1, &state.objects[i].material.ambient[0])
		gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.specularVal, 1, &state.objects[i].material.specular[0])
		gl.Uniform1f(state.objects[i].programInfo.uniformLocations.nVal, state.objects[i].material.n)

		gl.Uniform1i(state.objects[i].programInfo.uniformLocations.numLights, int32(len(state.lights)))

		//update camera
		camPosition := []float32{state.camera.position[0], state.camera.position[1], state.camera.position[2]}
		gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.cameraPosition, 1, &camPosition[0])

		//update lights
		if len(lightPositionArray) > 0 && len(lightColorArray) > 0 && len(lightStrengthArray) > 0 {
			gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.lightPositions, int32(len(lightPositionArray)/3), &lightPositionArray[0])
			gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.lightColours, int32(len(lightColorArray)/3), &lightColorArray[0])
			gl.Uniform1fv(state.objects[i].programInfo.uniformLocations.lightStrengths, int32(len(lightStrengthArray)), &lightStrengthArray[0])
		}

		gl.BindVertexArray(state.objects[i].buffers.vao)
		gl.DrawElements(gl.TRIANGLES, int32(len(triangleFaces)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
		gl.BindVertexArray(0)

		window.SwapBuffers()
	}

}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Go GL", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL(vertexShaderSource string, fragmentShaderSource string) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

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
