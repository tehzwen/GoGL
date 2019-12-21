package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"./game"
	"./geometry"
	"./mymath"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var keys map[glfw.Key]bool
var buttons map[glfw.MouseButton]bool
var mouseMovement map[string]float64
var mu sync.Mutex
var objectsToRender chan geometry.RenderObject

const (
	width  = 1280
	height = 960

	vertexShaderSource = `
		#version 300 es
		//needed to add layout location for mac to work properly
		layout (location = 0) in vec3 aPosition;
		layout (location = 1) in vec3 aNormal;

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
		#version 300 es
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
				vec3 nCameraPosition = normalize(cameraPosition); // Normalize the camera Position
				vec3 V = normalize(nCameraPosition - oFragPosition);

				vec3 lightDirection = normalize(lightPositions[i] - oFragPosition);
				float diff = max(dot(normal, lightDirection), 0.0);
				vec3 reflectDir = reflect(-lightDirection, normal);
				float spec = pow(max(dot(V, reflectDir), 0.0), nVal);
				float distance = length(lightPositions[i] - oFragPosition);
				float attenuation = 1.0 / (distance * distance);
				attenuation *= lightStrengths[i];

				ambient += ambientVal * lightColours[i] * diffuseVal;
				diffuse += diffuseVal * lightColours[i] * diff;

				if (diff > 0.0f) {
					specular += specularVal * lightColours[i] * spec;
					specular *= attenuation;
				}
				//ambient *= attenuation; causes much darker scene
				diffuse *= attenuation;
			}
			frag_colour = vec4(diffuse + ambient + specular, 1.0); 
		}
	` + "\x00"
)

var angle = 0.0

func main() {
	runtime.LockOSThread()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Terminated successfully!")
		os.Exit(1)
	}()

	//get arguments
	argsWithoutProgram := os.Args[1:]

	objectsToRender = make(chan geometry.RenderObject, 10)
	keys = make(map[glfw.Key]bool)
	buttons = make(map[glfw.MouseButton]bool)
	mouseMovement = make(map[string]float64)

	state := geometry.State{
		VertShader: vertexShaderSource,
		FragShader: fragmentShaderSource,
		Camera: geometry.Camera{
			Position: mgl32.Vec3{0.5, 0.0, -2.5},
			Center:   mgl32.Vec3{0.5, 0.0, 0.0},
			Up:       mgl32.Vec3{0.0, 1.0, 0.0},
		},
		Lights:        []geometry.Light{},
		Objects:       []geometry.Geometry{},
		Keys:          make(map[glfw.Key]bool),
		LoadedObjects: 0,
	}

	window := initGlfw()
	defer glfw.Terminate()

	window.SetKeyCallback(KeyHandler)
	window.SetMouseButtonCallback(MouseButtonHandler)
	window.SetCursorPosCallback(MouseMoveHandler)

	geometry.ParseJsonFile(argsWithoutProgram[0], &state)

	fmt.Println(len(state.Objects))

	then := 0.0

	game.Start(&state) //main logic start
	fmt.Println("PID: ", os.Getpid())

	for !window.ShouldClose() {
		if state.LoadedObjects == len(state.Objects) {

			now := glfw.GetTime()
			deltaTime := now - then
			then = now
			angle += 0.5 * deltaTime

			game.Update(&state, deltaTime) //main logic update

			state.Keys = keys

			if mouseMovement["move"] == 1 && buttons[glfw.MouseButton2] {
				Rotation := geometry.RotateY(state.Camera.Center, state.Camera.Position, -(2 * deltaTime * mouseMovement["Xmove"]))
				state.Camera.Center = Rotation
			}
			mouseMovement["move"] = 0
			draw(window, &state)
		}
	}
	fmt.Println("Program ended succesfully!")
}

func draw(window *glfw.Window, state *geometry.State) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearDepth(1.0)
	glfw.PollEvents()

	//getting the math values for each object using go routines
	for i := 0; i < len(state.Objects); i++ {
		go doObjectMath(state.Objects[i], (*state), objectsToRender)
	}

	//render sequentially using channel values
	state.RenderedObjects = 0
	for x := 0; x < len(state.Objects); x++ {
		tempObject := <-objectsToRender
		renderObject(state, tempObject)
	}
	//fmt.Println("Rendered ", state.RenderedObjects, " vs ", len(state.Objects))

	window.SwapBuffers()
}

func renderObject(state *geometry.State, object geometry.RenderObject) {

	currentProgramInfo := object.CurrentProgram
	projection := object.ProjMatrix
	viewMatrix := object.ViewMatrix
	camPosition := object.CameraPosition
	modelMatrix := object.ModelMatrix
	currentMaterial := object.CurrentMaterial
	currentBuffers := object.CurrentBuffers
	currentVertices := object.CurrentVertices

	gl.UseProgram(currentProgramInfo.Program)

	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.Projection, 1, false, &projection[0])
	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.View, 1, false, &viewMatrix[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.CameraPosition, 1, &camPosition[0])
	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.Model, 1, false, &modelMatrix[0])

	model, err := object.CurrentObject.GetModel()
	if err != nil {
		panic(err)
	}

	frustum := mymath.ConstructFrustrum(viewMatrix, projection)
	testLen := object.CurrentObject.GetBoundingBox().Max.Len()
	result := frustum.SphereIntersection(model.Position, testLen)

	if !result {
		return
	}
	diffuseTexture := object.CurrentObject.GetDiffuseTexture()

	if diffuseTexture != nil {
		diffuseTexture.Bind(gl.TEXTURE0)
		err := diffuseTexture.SetUniform(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("uDiffuseTexture\x00")))
		if err != nil {
			panic(err)
		}
	}

	state.RenderedObjects++

	gl.Uniform3fv(currentProgramInfo.UniformLocations.DiffuseVal, 1, &currentMaterial.Diffuse[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.AmbientVal, 1, &currentMaterial.Ambient[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.SpecularVal, 1, &currentMaterial.Specular[0])
	gl.Uniform1fv(currentProgramInfo.UniformLocations.NVal, 1, &currentMaterial.N)

	num := int32(len(state.Lights))
	gl.Uniform1iv(currentProgramInfo.UniformLocations.NumLights, 1, &num)

	for i := 0; i < len(state.Lights); i++ {
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].position\x00")), 1, &state.Lights[i].Position[0])
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].color\x00")), 1, &state.Lights[i].Colour[0])
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].strength\x00")), 1, &state.Lights[i].Strength)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].constant\x00")), 1, &state.Lights[i].Constant)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].linear\x00")), 1, &state.Lights[i].Linear)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].quadratic\x00")), 1, &state.Lights[i].Quadratic)
	}

	state.ViewMatrix = viewMatrix
	gl.BindVertexArray(currentBuffers.Vao)
	if object.CurrentObject.GetType() != "mesh" {
		gl.DrawElements(gl.TRIANGLES, int32(len(currentVertices.Vertices)), gl.UNSIGNED_INT, gl.Ptr(nil))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(currentVertices.Vertices)))
	}

	gl.BindVertexArray(0)

	if diffuseTexture != nil {
		diffuseTexture.UnBind()
	}

}

func doObjectMath(object geometry.Geometry, state geometry.State, objects chan<- geometry.RenderObject) {
	currentProgramInfo, err := object.GetProgramInfo()

	if err != nil {
		//else lets throw an error here
		fmt.Printf("ERROR getting program info!")
	}

	//now get the model
	currentModel, err := object.GetModel()
	if err != nil {
		//throw an error
		fmt.Printf("ERROR getting model!")
	}

	_, _, parent := object.GetDetails()
	currentCentroid := object.GetCentroid()
	currentMaterial := object.GetMaterial()
	currentBuffers := object.GetBuffers()
	currentVertices := object.GetVertices()

	var fovy = float32(60 * math.Pi / 180)
	var aspect = float32(width / height)
	var near = float32(0.1)
	var far = float32(100.0)

	projection := mgl32.Perspective(fovy, aspect, near, far)
	viewMatrix := mgl32.LookAtV(state.Camera.Position, state.Camera.Center, state.Camera.Up)
	camPosition := []float32{state.Camera.Position[0], state.Camera.Position[1], state.Camera.Position[2]}
	modelMatrix := mgl32.Ident4()
	positionMat := mgl32.Translate3D(currentModel.Position[0], currentModel.Position[1], currentModel.Position[2])
	modelMatrix = modelMatrix.Mul4(positionMat)
	centroidMat := mgl32.Translate3D(currentCentroid[0], currentCentroid[1], currentCentroid[2])
	modelMatrix = modelMatrix.Mul4(centroidMat)
	modelMatrix = modelMatrix.Mul4(currentModel.Rotation)
	negCent := mgl32.Translate3D(-currentCentroid[0], -currentCentroid[1], -currentCentroid[2])
	modelMatrix = modelMatrix.Mul4(negCent)
	modelMatrix = geometry.ScaleM4(modelMatrix, currentModel.Scale)

	if parent != "" {
		parentObj := geometry.GetSceneObject(parent, state)
		modelMatrix = modelMatrix.Mul4(parentObj.GetModelMatrix())
	}

	object.SetModelMatrix(modelMatrix)
	result := geometry.RenderObject{
		ModelMatrix:     modelMatrix,
		ViewMatrix:      viewMatrix,
		ProjMatrix:      projection,
		CameraPosition:  camPosition,
		CurrentBuffers:  currentBuffers,
		CurrentCentroid: currentCentroid,
		CurrentMaterial: currentMaterial,
		CurrentModel:    currentModel,
		CurrentProgram:  currentProgramInfo,
		CurrentVertices: currentVertices,
		CurrentObject:   object,
	}
	objects <- result
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
