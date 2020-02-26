package game

import (
	"fmt"

	"../geometry"
	"../scene"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var angle = 0.0
var walkSpeed float64 = 1
var runSpeed float64 = 5
var lightToMove *geometry.PointLight

var cube3 geometry.Geometry

// var wall0 geometry.Geometry
var err error

// Start : initialize our values for our game here
func Start(state *geometry.State) {
	fmt.Printf("Started!\n")
	lightToMove = &state.PointLights[0]
	cube3, err = scene.GetObjectFromScene(state, "testCube3")

	if err != nil {
		panic(err)
	}

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

	lightToMove.Move = false

	// 	if state.Keys[glfw.Key5] {
	// 		rot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	// 		wall0.SetRotation(rot)
	// 	}

	rot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	cube3.SetRotation(rot)

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

	if state.Keys[glfw.KeyLeft] {
		cube3.AddForce(mgl32.Vec3{-0.5 * float32(deltaTime), 0.0, 0.0})
	}
	if state.Keys[glfw.KeyRight] {
		cube3.AddForce(mgl32.Vec3{0.5 * float32(deltaTime), 0.0, 0.0})
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
