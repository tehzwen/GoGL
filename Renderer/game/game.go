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

// var cube3 geometry.Geometry
// var wall0 geometry.Geometry
// var err error

// Start : initialize our values for our game here
func Start(state *geometry.State) {
	fmt.Printf("Started!\n")
	lightToMove = &state.Lights[0]
	// 	cube3, err = scene.GetObjectFromScene(state, "testCube3")

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	wall0, err = scene.GetObjectFromScene(state, "wall1")

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	//we can set an event listener for when this object collides
	// 	cube3.SetOnCollide(func(box geometry.BoundingBox) {
	// 		fmt.Println("Cube collision!")
	// 		//cube3.SetForce(mgl32.Vec3{0, 0, 0})
	// 		currentForce := cube3.GetForce()
	// 		//reduce the force by a small margin due to collision
	// 		currentForce = currentForce.Mul(0.9)
	// 		cube3.SetForce(mgl32.Vec3{-currentForce[0], -currentForce[1], -currentForce[2]})
	// 	})

}

// Update : runs each frame
func Update(state *geometry.State, deltaTime float64) {
	speed := deltaTime * walkSpeed

	if state.Keys[glfw.KeyLeftShift] {
		speed *= runSpeed
	}

	// 	if state.Keys[glfw.Key5] {
	// 		rot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	// 		wall0.SetRotation(rot)
	// 	}

	// 	if state.Keys[glfw.Key1] {
	// 		cube3.AddForce(mgl32.Vec3{0.5 * float32(deltaTime), 0.0, 0.0})
	// 	}

	// 	if state.Keys[glfw.Key2] {
	// 		cube3.AddForce(mgl32.Vec3{-0.5 * float32(deltaTime), 0.0, 0.0})
	// 	}

	if state.Keys[glfw.KeyQ] {
		lightToMove.Position[2] += 1.0
	}

	if state.Keys[glfw.KeyR] {
		lightToMove.Position[2] -= 1.0
	}

	if state.Keys[glfw.KeyT] {
		if lightToMove.Constant > 0 {
			lightToMove.Constant -= 0.025
		}

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
