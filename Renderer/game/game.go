package game

import (
	"fmt"

	"../geometry"
	"../scene"
	"github.com/go-gl/glfw/v3.1/glfw"
	_ "github.com/go-gl/mathgl/mgl32"
)

var angle = 0.0
var walkSpeed float64 = 5
var runSpeed float64 = 5
var lightToMove *geometry.PointLight


var testcube geometry.Geometry
var err error

// Start : initialize our values for our game here
func Start(state *geometry.State) {
	fmt.Printf("Started!\n")
	lightToMove = &state.PointLights[0]
	testcube, err = scene.GetObjectFromScene(state, "container1")
	
	if err != nil {
	    panic(err)
	}

}

// Update : runs each frame
func Update(state *geometry.State, deltaTime float64) {
	speed := deltaTime * walkSpeed

	lightToMove.Move = false
	
	//rot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	//testcube.SetRotation(rot)

	if state.Keys[glfw.KeyQ] {
		lightToMove.Position[2] += 0.1
		lightToMove.Move = true
	}

	if state.Keys[glfw.KeyR] {
		lightToMove.Position[2] -= 0.1
		lightToMove.Move = true
	}

	if state.Keys[glfw.Key3] {
		lightToMove.Position[0] += 0.1
		lightToMove.Move = true
	}

	if state.Keys[glfw.Key1] {
		lightToMove.Position[0] -= 0.1
		lightToMove.Move = true
	}

	// if state.Keys[glfw.KeyT] {
	// 	if lightToMove.Constant > 0 {
	// 		lightToMove.Constant -= 0.025
	// 	}
	// }

	if state.Keys[glfw.KeyLeftShift] {
		speed *= runSpeed
	}

	if state.Keys[glfw.KeyW] {
		MoveForward(state, speed)
	}
	if state.Keys[glfw.KeyS] {
		MoveBackward(state, speed)
	}
	if state.Keys[glfw.KeyA] {
		MoveLeft(state, speed)
	}
	if state.Keys[glfw.KeyD] {
		MoveRight(state, speed)
	}

	angle += 0.5 * deltaTime

}
