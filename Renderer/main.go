package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sort"
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
			Position: mgl32.Vec3{6, 0.0, 0},
			Front:    mgl32.Vec3{0, 0, 1.0},
			Up:       mgl32.Vec3{0.0, 1.0, 0.0},
			Pitch:    0,
			Yaw:      90,
			Roll:     0,
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

			draw(window, &state)
		}
	}
	fmt.Println("Program ended successfully!")
}

func draw(window *glfw.Window, state *geometry.State) {

	glfw.PollEvents()

	//getting the math values for each object using go routines
	for i := 0; i < len(state.Objects); i++ {
		go doObjectMath(state.Objects[i], (*state), objectsToRender)
	}

	// fmt.Println(state.Objects)

	// fmt.Println("SORTED: ", state.Objects)

	//render sequentially using channel values
	state.RenderedObjects = 0
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.MULTISAMPLE)

	tempList := []geometry.RenderObject{}

	for x := 0; x < len(state.Objects); x++ {
		tempList = append(tempList, <-objectsToRender)
	}

	//sort for transparency
	sort.Slice(tempList, func(a, b int) bool {
		//get model info
		nameA, _, _ := tempList[a].CurrentObject.GetDetails()
		nameB, _, _ := tempList[b].CurrentObject.GetDetails()
		distA := tempList[a].DistanceToCamera
		distB := tempList[b].DistanceToCamera

		if distA > distB {
			return true
		} else if distB > distA {
			return false
		} else {
			return nameA > nameB
		}
	})

	for x := 0; x < len(tempList); x++ {
		renderObject(state, tempList[x])
	}

	window.SwapBuffers()
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
				//fmt.Println("RESET")
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

func renderObject(state *geometry.State, object geometry.RenderObject) {
	currentProgramInfo := object.CurrentProgram
	projection := object.ProjMatrix
	viewMatrix := object.ViewMatrix
	camPosition := object.CameraPosition
	modelMatrix := object.ModelMatrix
	currentMaterial := object.CurrentMaterial
	currentBuffers := object.CurrentBuffers
	currentVertices := object.CurrentVertices

	//do movement physics and collision testing
	currentForce := object.CurrentObject.GetForce()
	object.CurrentObject.Translate(currentForce)

	if object.CurrentObject.GetBoundingBox().Collide {
		collisionTest(state, object)
	}

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
	testLen := object.CurrentObject.GetBoundingBox().Max.LenSqr()
	result := frustum.SphereIntersection(model.Position, testLen)

	if !result {
		return
	}
	diffuseTexture := object.CurrentObject.GetDiffuseTexture()
	normalTexture := object.CurrentObject.GetNormalTexture()

	if diffuseTexture != nil {
		diffuseTexture.Bind(gl.TEXTURE0)
		err := diffuseTexture.SetUniform(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("uDiffuseTexture\x00")))
		if err != nil {
			panic(err)
		}
	}

	if normalTexture != nil {
		normalTexture.Bind(gl.TEXTURE1)
		err := normalTexture.SetUniform(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("uNormalTexture\x00")))
		if err != nil {
			panic(err)
		}
	}

	state.RenderedObjects++

	gl.Uniform3fv(currentProgramInfo.UniformLocations.DiffuseVal, 1, &currentMaterial.Diffuse[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.AmbientVal, 1, &currentMaterial.Ambient[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.SpecularVal, 1, &currentMaterial.Specular[0])
	gl.Uniform1fv(currentProgramInfo.UniformLocations.NVal, 1, &currentMaterial.N)
	gl.Uniform1fv(currentProgramInfo.UniformLocations.Alpha, 1, &currentMaterial.Alpha)

	if currentMaterial.Alpha < 1.0 {
		//name, _, _ := object.CurrentObject.GetDetails()
		//fmt.Println("here: ", name)
		gl.Enable(gl.BLEND)
		gl.Disable(gl.DEPTH_TEST)
		gl.BlendFunc(gl.ONE_MINUS_CONSTANT_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.ClearDepth(float64(currentMaterial.Alpha))
	} else {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)
		gl.ClearDepth(1.0)
		gl.DepthFunc(gl.LEQUAL)
	}

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

	if normalTexture != nil {
		normalTexture.UnBind()
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
	//create a camfront value
	camFront := state.Camera.Position.Add(state.Camera.Front)
	viewMatrix := mgl32.LookAtV(state.Camera.Position, camFront, state.Camera.Up)
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
	camDist := geometry.VectorDistance(state.Camera.Position, currentModel.Position.Add(currentCentroid)) //calculate this for transparency

	if parent != "" {
		parentObj := geometry.GetSceneObject(parent, state)
		if parentObj != nil {
			parentMMatrix, err := parentObj.GetModelMatrix()
			if err == nil {
				modelMatrix = modelMatrix.Mul4(parentMMatrix)
			}
		} else {
			fmt.Println("ERROR GETTING PARENT OBJECT")
		}
	}

	object.SetModelMatrix(modelMatrix)
	result := geometry.RenderObject{
		ModelMatrix:      modelMatrix,
		ViewMatrix:       viewMatrix,
		ProjMatrix:       projection,
		CameraPosition:   camPosition,
		CurrentBuffers:   currentBuffers,
		CurrentCentroid:  currentCentroid,
		CurrentMaterial:  currentMaterial,
		CurrentModel:     currentModel,
		CurrentProgram:   currentProgramInfo,
		CurrentVertices:  currentVertices,
		CurrentObject:    object,
		DistanceToCamera: camDist,
	}
	objects <- result
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Samples, 4)
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
