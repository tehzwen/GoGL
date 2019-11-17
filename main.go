package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 500
	height = 500

	vertexShaderSource = `
		#version 330
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
		#version 330
		precision highp float;
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
			vec3 diffuse = vec3(0, 0, 0);
			vec3 ambient = vec3(0, 0, 0);
			vec3 specular = vec3(0, 0, 0);
			vec3 normal = normalize(normalInterp);

			for (int i = 0; i < numLights; i++) {
				vec3 nCameraPosition = normalize(cameraPosition); // Normalize the camera position
				vec3 V = normalize(nCameraPosition - oFragPosition);

				vec3 lightDirection = normalize(lightPositions[i] - oFragPosition);
				float diff = max(dot(normal, lightDirection), 0.0);
				vec3 reflectDir = reflect(-lightDirection, normal);
				float spec = pow(max(dot(V, reflectDir), 0.0), nVal);
				float distance = length(lightPositions[i] - oFragPosition);
				float attenuation = 1.0 / (distance * distance);
				attenuation *= lightStrengths[i];

				ambient += ambientVal * lightColours[i];
				diffuse += diffuseVal * lightColours[i] * diff;
				specular += specularVal * lightColours[i] * spec;

				ambient *= attenuation;
				diffuse *= attenuation;
				specular *= attenuation;

			}
			frag_colour = vec4(diffuse + ambient + specular, 1.0);
		}
	` + "\x00"
)

var angle = 0.0

func main() {
	runtime.LockOSThread()

	var state = State{
		vertShader: vertexShaderSource,
		fragShader: fragmentShaderSource,
		camera: Camera{
			position: mgl32.Vec3{0.0, 0.5, -2.5},
			center:   mgl32.Vec3{0.5, 0.0, 0.0},
			up:       mgl32.Vec3{0.0, 1.0, 0.0},
		},
		lights: []Light{
			Light{
				colour:   []float32{1.0, 1.0, 1.0},
				strength: 10,
				position: []float32{1.0, 0.0, -2.0},
			},
			Light{
				colour:   []float32{1.0, 1.0, 1.0},
				strength: 0.1,
				position: []float32{2.0, 0.0, -0.5},
			},
		},
	}

	window := initGlfw()
	defer glfw.Terminate()
	var objectList = []Geometry{}

	testCube := Cube{}
	testCube2 := Cube{}

	err := testCube.SetShader(vertexShaderSource, fragmentShaderSource)
	err = testCube2.SetShader(vertexShaderSource, fragmentShaderSource)

	if err != nil {
		panic(err)
	} else {
		testCube.Setup(
			Material{
				diffuse:  []float32{0.6, 0.8, 0.6},
				ambient:  []float32{0.1, 0.1, 0.1},
				specular: []float32{0.8, 0.8, 0.8},
				n:        100,
			},
			Model{
				position: mgl32.Vec3{1, 0, 0},
				scale:    mgl32.Vec3{1, 1, 1},
				rotation: mgl32.Ident4(),
			}, "testcube1")

		testCube2.Setup(
			Material{
				diffuse:  []float32{0.6, 0.2, 0.6},
				ambient:  []float32{0.1, 0.1, 0.1},
				specular: []float32{0.8, 0.8, 0.8},
				n:        100,
			}, Model{
				position: mgl32.Vec3{0, 0, 0},
				scale:    mgl32.Vec3{1, 1, 1},
				rotation: mgl32.Ident4(),
			}, "testcube2")

		objectList = append(objectList, &testCube)
		objectList = append(objectList, &testCube2)
	}

	state.objects = objectList

	then := 0.0

	for !window.ShouldClose() {
		now := glfw.GetTime()
		deltaTime := now - then
		then = now
		angle += 0.5 * deltaTime

		if err == nil {
			newRot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
			state.objects[0].SetRotation(newRot)
		}
		draw(window, state)
	}
}

func draw(window *glfw.Window, state State) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
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

		currentProgramInfo, err := state.objects[i].GetProgramInfo()

		if err != nil {
			//else lets throw an error here
			fmt.Printf("ERROR getting program info!")
		}
		gl.UseProgram(currentProgramInfo.program)
		//now get the model
		currentModel, err := state.objects[i].GetModel()
		if err != nil {
			//throw an error
			fmt.Printf("ERROR getting model!")
		}

		currentCentroid := state.objects[i].GetCentroid()
		currentMaterial := state.objects[i].GetMaterial()
		currentBuffers := state.objects[i].GetBuffers()
		currentVertices := state.objects[i].GetVertices()

		var fovy = float32(60 * math.Pi / 180)
		var aspect = float32(width / height)
		var near = float32(0.1)
		var far = float32(100.0)

		//projection matrix
		projection := mgl32.Perspective(fovy, aspect, near, far)
		gl.UniformMatrix4fv(currentProgramInfo.uniformLocations.projection, 1, false, &projection[0])

		//view matrix
		viewMatrix := mgl32.LookAtV(state.camera.position, state.camera.center, state.camera.up)
		gl.UniformMatrix4fv(currentProgramInfo.uniformLocations.view, 1, false, &viewMatrix[0])

		//model matrix
		modelMatrix := mgl32.Ident4()
		positionMat := mgl32.Translate3D(currentModel.position[0], currentModel.position[1], currentModel.position[2])
		modelMatrix = modelMatrix.Mul4(positionMat)
		centroidMat := mgl32.Translate3D(currentCentroid[0], currentCentroid[1], currentCentroid[2])
		modelMatrix = modelMatrix.Mul4(centroidMat)
		modelMatrix = modelMatrix.Mul4(currentModel.rotation)
		negCent := mgl32.Translate3D(-currentCentroid[0], -currentCentroid[1], -currentCentroid[2])
		modelMatrix = modelMatrix.Mul4(negCent)

		gl.UniformMatrix4fv(currentProgramInfo.uniformLocations.model, 1, false, &modelMatrix[0])

		gl.Uniform3fv(currentProgramInfo.uniformLocations.diffuseVal, 1, &currentMaterial.diffuse[0])
		gl.Uniform3fv(currentProgramInfo.uniformLocations.ambientVal, 1, &currentMaterial.ambient[0])
		gl.Uniform3fv(currentProgramInfo.uniformLocations.specularVal, 1, &currentMaterial.specular[0])
		gl.Uniform1f(currentProgramInfo.uniformLocations.nVal, currentMaterial.n)

		gl.Uniform1i(currentProgramInfo.uniformLocations.numLights, int32(len(state.lights)))

		//update camera
		camPosition := []float32{state.camera.position[0], state.camera.position[1], state.camera.position[2]}
		gl.Uniform3fv(currentProgramInfo.uniformLocations.cameraPosition, 1, &camPosition[0])

		//update lights
		if len(lightPositionArray) > 0 && len(lightColorArray) > 0 && len(lightStrengthArray) > 0 {
			gl.Uniform3fv(currentProgramInfo.uniformLocations.lightPositions, int32(len(lightPositionArray)/3), &lightPositionArray[0])
			gl.Uniform3fv(currentProgramInfo.uniformLocations.lightColours, int32(len(lightColorArray)/3), &lightColorArray[0])
			gl.Uniform1fv(currentProgramInfo.uniformLocations.lightStrengths, int32(len(lightStrengthArray)), &lightStrengthArray[0])
		}

		gl.BindVertexArray(currentBuffers.vao)
		gl.DrawElements(gl.TRIANGLES, int32(len(currentVertices.vertices)), gl.UNSIGNED_INT, gl.Ptr(nil))
		gl.BindVertexArray(0)

		//window.SwapBuffers()
	}
	window.SwapBuffers()

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
