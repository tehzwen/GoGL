package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"./game"
	"./geometry"
	"./globals"
	"./mymath"
	"./shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var keys map[glfw.Key]bool
var buttons map[glfw.MouseButton]bool
var mouseMovement map[string]float64
var objectsToRender chan geometry.RenderObject

const (
	thread = false
)

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
	mouseMovement["sensitivity"] = 1.2

	state := geometry.State{
		Camera: geometry.Camera{
			Position: mgl32.Vec3{-1, 2.0, -3},
			Front:    mgl32.Vec3{0, 0, 1.0},
			Up:       mgl32.Vec3{0.0, 1.0, 0.0},
			Pitch:    0,
			Yaw:      90,
			Roll:     0,
		},
		PointLights:    []geometry.PointLight{},
		Objects:        []geometry.Geometry{},
		Keys:           make(map[glfw.Key]bool),
		LoadedObjects:  0,
		CurrentTexUnit: 0,
	}

	window := initGlfw()
	defer glfw.Terminate()

	window.SetKeyCallback(KeyHandler)
	window.SetMouseButtonCallback(MouseButtonHandler)
	window.SetCursorPosCallback(MouseMoveHandler)

	if len(argsWithoutProgram) <= 0 {
		geometry.ParseJSONFile("../Editor/statefiles/testsave.json", &state)
	} else {
		geometry.ParseJSONFile(argsWithoutProgram[0], &state)
	}

	then := 0.0

	game.Start(&state) //main logic start
	fmt.Println("PID: ", os.Getpid())
	gl.GenFramebuffers(1, &state.DepthFBO)

	//iterate through pointlights and create depth maps for each
	for l := 0; l < len(state.PointLights); l++ {
		state.PointLights[l].CreateLightSpaceTransforms(0.5, 25, 1024, 1024)
		state.PointLights[l].CreateCubeDepthMap(1024, 1024)
	}

	for l := 0; l < len(state.DirectionalLights); l++ {
		state.DirectionalLights[l].CreateLightSpaceTransforms(1.0, 7.5)
		state.DirectionalLights[l].CreateDirectionalDepthMap(1024, 1024)
	}

	//setup pointlightshadow shader program
	shadowShaderVals := make(map[string]bool)
	shadowShaderVals["uModelMatrix"] = true
	shadowShaderVals["aPosition"] = true
	shadowShaderVals["shadowMatrices"] = true
	shadowShaderVals["lightPos"] = true
	shadShader := &shader.OmniDirectionalShadow{}
	shadShader.Setup()
	pointLightShadowProgramInfo := geometry.ProgramInfo{}
	pointLightShadowProgramInfo.Program = geometry.InitOpenGL(shadShader.GetVertShader(), shadShader.GetFragShader(), shadShader.GetGeometryShader())
	shadowProgAttribs := geometry.Attributes{}
	shadowProgAttribs.SetPosition(0)
	pointLightShadowProgramInfo.SetAttributes(shadowProgAttribs)
	geometry.SetupAttributesMap(&pointLightShadowProgramInfo, shadowShaderVals)

	shadowShaderVals["shadowMatrices"] = false
	shadowShaderVals["lightPos"] = false
	shadowShaderVals["lightSpaceMatrix"] = true

	dirShadShader := &shader.DirectionalShadow{}
	dirShadShader.Setup()
	dirLightShadowProgramInfo := geometry.ProgramInfo{}
	dirLightShadowProgramInfo.Program = geometry.InitOpenGL(dirShadShader.GetVertShader(), dirShadShader.GetFragShader(), dirShadShader.GetGeometryShader())
	dirLightShadowProgramInfo.SetAttributes(shadowProgAttribs)
	geometry.SetupAttributesMap(&dirLightShadowProgramInfo, shadowShaderVals)

	for !window.ShouldClose() {
		if state.LoadedObjects == len(state.Objects) {

			now := glfw.GetTime()
			deltaTime := now - then
			then = now

			game.Update(&state, deltaTime) //main logic update

			state.Keys = keys

			if mouseMovement["move"] == 1 && buttons[glfw.MouseButton2] {
				front := mgl32.Vec3{0, 0, 0}
				state.Camera.Yaw += float32(mouseMovement["Xmove"] * mouseMovement["sensitivity"])
				state.Camera.Pitch += float32(mouseMovement["Ymove"] * mouseMovement["sensitivity"])

				if state.Camera.Pitch > 89 {
					state.Camera.Pitch = 89
				}
				if state.Camera.Pitch < -89 {
					state.Camera.Pitch = -89
				}

				front[0] = float32(math.Cos(geometry.ToRadians(state.Camera.Yaw)) * math.Cos(geometry.ToRadians(state.Camera.Pitch)))
				front[1] = float32(math.Sin(geometry.ToRadians(-state.Camera.Pitch)))
				front[2] = float32(math.Sin(geometry.ToRadians(state.Camera.Yaw)) * math.Cos(geometry.ToRadians(state.Camera.Pitch)))

				front = front.Normalize()

				state.Camera.Front = front

				//Rotation := geometry.RotateY(state.Camera.Center, state.Camera.Position, -(2 * deltaTime * mouseMovement["Xmove"]))
				//fmt.Println(mouseMovement["Ymove"])
				//state.Camera.Center = Rotation
			}
			mouseMovement["move"] = 0

			if thread {
				MultithreadRender(window, &state)
			} else {
				draw(window, &state, &pointLightShadowProgramInfo, &dirLightShadowProgramInfo)
			}
		}
	}
	fmt.Println("Program ended successfully!")
}

//TODO make cleaner pass of shadow programinfos
func draw(window *glfw.Window, state *geometry.State, pointLightShadowProgramInfo, dirLightShadowProgramInfo *geometry.ProgramInfo) {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// err := gl.GetError()

	// if err != gl.NO_ERROR {
	// 	fmt.Println(err)
	// 	//panic(err)
	// }

	glfw.PollEvents()

	//going to have to render depth for each pointlight here
	for l := 0; l < len(state.PointLights); l++ {
		state.PointLights[l].BindDepthMap(state)
		gl.Viewport(0, 0, 1024, 1024)
		gl.BindFramebuffer(gl.FRAMEBUFFER, state.DepthFBO)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		for x := 0; x < len(state.Objects); x++ {
			state.PointLights[l].ShadowRender(state, state.Objects[x], pointLightShadowProgramInfo)
		}
		gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, 0, 0)
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	}

	//Depth draw directional lights
	for l := 0; l < len(state.DirectionalLights); l++ {
		state.DirectionalLights[l].BindDepthMap(state)
		gl.Viewport(0, 0, 1024, 1024)
		gl.BindFramebuffer(gl.FRAMEBUFFER, state.DepthFBO)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		for x := 0; x < len(state.Objects); x++ {
			state.DirectionalLights[l].ShadowRender(state, state.Objects[x], dirLightShadowProgramInfo)
		}
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, 0, 0)
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	}

	//try the classical render method
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, int32(globals.Width), int32(globals.Height))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for i := 0; i < len(state.Objects); i++ {
		ClassicRender(state, state.Objects[i])
	}
	window.SwapBuffers()
}

//Classic non threaded render
func ClassicRender(state *geometry.State, object geometry.Geometry) {
	currentProgramInfo, err := object.GetProgramInfo()
	if err != nil {
		panic(err)
	}

	currentBuffers := object.GetBuffers()

	gl.UseProgram(currentProgramInfo.Program)

	currentForce := object.GetForce()
	object.Translate(currentForce)

	//now get the model
	currentModel, err := object.GetModel()
	if err != nil {
		//throw an error
		fmt.Printf("ERROR getting model!")
	}

	_, _, parent := object.GetDetails()
	currentCentroid := object.GetCentroid()
	currentMaterial := object.GetMaterial()
	currentVertices := object.GetVertices()

	state.RenderedObjects++

	var fovy = float32(60 * math.Pi / 180)
	var aspect = float32(globals.Width / globals.Height)
	var near = float32(0.1)
	var far = float32(1000.0)

	projection := mgl32.Perspective(fovy, aspect, near, far)
	//create a camfront value
	camFront := state.Camera.Position.Add(state.Camera.Front)
	viewMatrix := mgl32.LookAtV(state.Camera.Position, camFront, state.Camera.Up)
	camPosition := []float32{state.Camera.Position[0], state.Camera.Position[1], state.Camera.Position[2]}
	modelMatrix := mgl32.Ident4()

	//move to centroid
	centroidMat := mgl32.Translate3D(currentCentroid[0], currentCentroid[1], currentCentroid[2])
	modelMatrix = modelMatrix.Mul4(centroidMat)

	//rotation
	modelMatrix = modelMatrix.Mul4(currentModel.Rotation)

	//position
	positionMat := mgl32.Translate3D(currentModel.Position[0], currentModel.Position[1], currentModel.Position[2])
	modelMatrix = modelMatrix.Mul4(positionMat)

	//negative centroid
	negCent := mgl32.Translate3D(-currentCentroid[0], -currentCentroid[1], -currentCentroid[2])
	modelMatrix = modelMatrix.Mul4(negCent)

	//scale
	modelMatrix = geometry.ScaleM4(modelMatrix, currentModel.Scale)

	//camDist := geometry.VectorDistance(state.Camera.Position, currentModel.Position.Add(currentCentroid)) //calculate this for transparency

	// if currentMaterial.Alpha < 1.0 {
	// 	//name, _, _ := object.CurrentObject.GetDetails()
	// 	//fmt.Println("here: ", name)
	// 	gl.Enable(gl.BLEND)
	// 	gl.Disable(gl.DEPTH_TEST)
	// 	gl.BlendFunc(gl.ONE_MINUS_CONSTANT_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// 	gl.ClearDepth(float64(currentMaterial.Alpha))
	// } else {
	// 	gl.Enable(gl.DEPTH_TEST)
	// 	gl.DepthMask(true)
	// 	gl.Disable(gl.BLEND)
	// 	gl.ClearDepth(1.0)
	// 	gl.DepthFunc(gl.LEQUAL)
	// }

	if parent != "" {
		parentObj := geometry.GetSceneObject(parent, (*state))
		if parentObj != nil {
			parentMMatrix, err := parentObj.GetModelMatrix()
			if err == nil {
				modelMatrix = modelMatrix.Mul4(parentMMatrix)
			}
		} else {
			fmt.Println("ERROR GETTING PARENT OBJECT")
		}
	}

	state.ViewMatrix = viewMatrix

	object.SetModelMatrix(modelMatrix)
	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.Projection, 1, false, &projection[0])
	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.View, 1, false, &viewMatrix[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.CameraPosition, 1, &camPosition[0])

	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.Model, 1, false, &modelMatrix[0])

	model, err := object.GetModel()
	if err != nil {
		panic(err)
	}

	frustum := mymath.ConstructFrustrum(viewMatrix, projection)
	testLen := object.GetBoundingBox().Max.LenSqr()
	result := frustum.SphereIntersection(model.Position, testLen)

	if !result {
		return
	}

	gl.Uniform3fv(currentProgramInfo.UniformLocations.DiffuseVal, 1, &currentMaterial.Diffuse[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.AmbientVal, 1, &currentMaterial.Ambient[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.SpecularVal, 1, &currentMaterial.Specular[0])
	gl.Uniform1fv(currentProgramInfo.UniformLocations.NVal, 1, &currentMaterial.N)
	gl.Uniform1fv(currentProgramInfo.UniformLocations.Alpha, 1, &currentMaterial.Alpha)

	numPointLights := int32(len(state.PointLights))
	gl.Uniform1iv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("numPointLights\x00")), 1, &numPointLights)
	numDirLights := int32(len(state.DirectionalLights))
	gl.Uniform1iv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("numDirLights\x00")), 1, &numDirLights)

	diffuseTexture := object.GetDiffuseTexture()
	normalTexture := object.GetNormalTexture()

	if diffuseTexture != nil {

		diffuseTex := diffuseTexture.GetHandle()
		gl.ActiveTexture(gl.TEXTURE0 + diffuseTex)
		gl.BindTexture(gl.TEXTURE_2D, diffuseTex)
		gl.Uniform1i(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("uDiffuseTexture\x00")), int32(diffuseTex))
	}

	if normalTexture != nil {
		normTex := normalTexture.GetHandle()
		gl.ActiveTexture(gl.TEXTURE0 + normTex)
		gl.BindTexture(gl.TEXTURE_2D, normTex)
		gl.Uniform1i(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("uNormalTexture\x00")), int32(normTex))
	}

	for i := 0; i < len(state.PointLights); i++ {
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].position\x00")), 1, &state.PointLights[i].Position[0])
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].color\x00")), 1, &state.PointLights[i].Colour[0])
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].strength\x00")), 1, &state.PointLights[i].Strength)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].constant\x00")), 1, &state.PointLights[i].Constant)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].linear\x00")), 1, &state.PointLights[i].Linear)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].quadratic\x00")), 1, &state.PointLights[i].Quadratic)
		gl.ActiveTexture(gl.TEXTURE0 + state.PointLights[i].DepthMap)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, state.PointLights[i].DepthMap)
		gl.Uniform1i(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].depthMap\x00")), int32(state.PointLights[i].DepthMap))
	}

	// for i := 0; i < len(state.DirectionalLights); i++ {
	// 	gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"dirLights[", strconv.Itoa(i)}, "")+"].direction\x00")), 1, &state.DirectionalLights[i].Direction[0])
	// 	gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"dirLights[", strconv.Itoa(i)}, "")+"].color\x00")), 1, &state.DirectionalLights[i].Colour[0])
	// 	gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"dirLights[", strconv.Itoa(i)}, "")+"].strength\x00")), 1, &state.DirectionalLights[i].Strength)
	// 	gl.UniformMatrix4fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"dirLights[", strconv.Itoa(i)}, "")+"].lightSpaceMatrix\x00")), 1, false, &state.DirectionalLights[i].LightViewMatrix[0])
	// 	gl.ActiveTexture(gl.TEXTURE0 + state.DirectionalLights[i].DepthMap)
	// 	gl.BindTexture(gl.TEXTURE_2D, state.DirectionalLights[i].DepthMap)
	// 	gl.Uniform1i(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"dirLights[", strconv.Itoa(i)}, "")+"].depthMap\x00")), int32(state.DirectionalLights[i].DepthMap))
	// }

	if numDirLights > 0 {
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("dirLight.direction\x00")), 1, &state.DirectionalLights[0].Direction[0])
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("dirLight.color\x00")), 1, &state.DirectionalLights[0].Colour[0])
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("dirLight.position\x00")), 1, &state.DirectionalLights[0].Position[0])
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("dirLight.strength\x00")), 1, &state.DirectionalLights[0].Strength)
		gl.UniformMatrix4fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("dirLight.lightSpaceMatrix\x00")), 1, false, &state.DirectionalLights[0].LightViewMatrix[0])
		gl.ActiveTexture(gl.TEXTURE0 + state.DirectionalLights[0].DepthMap)
		gl.BindTexture(gl.TEXTURE_2D, state.DirectionalLights[0].DepthMap)
		gl.Uniform1i(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("dirLight.depthMap\x00")), int32(state.DirectionalLights[0].DepthMap))

	}

	gl.BindVertexArray(currentBuffers.Vao)
	if object.GetType() != "mesh" {
		gl.DrawElements(gl.TRIANGLES, int32(len(currentVertices.Vertices)), gl.UNSIGNED_INT, gl.Ptr(nil))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(currentVertices.Vertices)))
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindVertexArray(0)
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Samples, 3)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(globals.Width, globals.Height, "Go GL", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(frameBufferSizeCallback)

	return window
}

func collisionTest(state *geometry.State, object geometry.RenderObject) {
	currentName, _, _ := object.CurrentObject.GetDetails()
	currBox := object.CurrentObject.GetBoundingBox()

	//iterate through all collidable objects
	for x := 0; x < len(state.Objects); x++ {
		//need to check for itself
		name, _, _ := state.Objects[x].GetDetails()
		if name == currentName {
			continue
		}

		box := state.Objects[x].GetBoundingBox()
		if box.Collide {

			collide := geometry.Intersect(box, currBox)
			if collide && currBox.CollisionBody != name {
				object.CurrentObject.SetBoundingBox(geometry.BoundingBox{
					Max:            currBox.Max,
					Min:            currBox.Min,
					Collide:        currBox.Collide,
					CollisionCount: currBox.CollisionCount + 1,
					CollisionBody:  name,
				})
				object.CurrentObject.OnCollide(box)

			} else if currBox.CollisionBody == name {
				object.CurrentObject.SetBoundingBox(geometry.BoundingBox{
					Max:            currBox.Max,
					Min:            currBox.Min,
					Collide:        currBox.Collide,
					CollisionCount: 0,
					CollisionBody:  "",
				})
			}
		}
	}
}

func frameBufferSizeCallback(window *glfw.Window, width, height int) {
	globals.Width = width
	globals.Height = height
	gl.Viewport(0, 0, int32(width), int32(height))
}
