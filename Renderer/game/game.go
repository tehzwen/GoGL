package game

import (
	"fmt"

	"../geometry"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var angle = 0.0
var walkSpeed float64 = 1
var runSpeed float64 = 5
var lightToMove *geometry.Light

// Start : initialize our values for our game here
func Start(state *geometry.State) {
	fmt.Printf("Started!\n")
	lightToMove = &state.Lights[0]
}

// Update : runs each frame
func Update(state *geometry.State, deltaTime float64) {
	/*newRot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	state.Objects[0].SetRotation(newRot)
	angle += 0.5 * deltaTime */
	angle += 0.5 * deltaTime

	//lightToMove.Strength += float32(math.Sin(angle))
	//lightToMove.Position[2] += float32(math.Sin(angle) * 0.5)
	//angle += 0.5 * deltaTime
	//fmt.Printf("position: %+v\n", lightToMove.Position)
	speed := deltaTime * walkSpeed

	if state.Keys[glfw.KeyLeftShift] {
		speed *= runSpeed
	}

	if state.Keys[glfw.KeyQ] {
		lightToMove.Strength += 0.5
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

}
