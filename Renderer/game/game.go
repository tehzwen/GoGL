package game

import (
	"fmt"
	"math"

	"../geometry"
	"../scene"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var angle = 0.0
var conveyerSpeed float32 = 0.008
var walkSpeed float64 = 5
var runSpeed float64 = 5
var lightToMove *geometry.PointLight
var warningLight1 *geometry.PointLight
var warningLight2 *geometry.PointLight

var transportCube1 geometry.Geometry
var chain1 geometry.Geometry
var transportHook1 geometry.Geometry

var transportCube2 geometry.Geometry
var chain2 geometry.Geometry
var transportHook2 geometry.Geometry

var viper geometry.Geometry
var err error
var warning1Reset bool

func TransportCollide(box geometry.BoundingBox, parent, hook, chain geometry.Geometry) {
	currentForce := parent.GetForce()
	parent.SetForce(mgl32.Vec3{-currentForce[0], -currentForce[1], -currentForce[2]})
	hook.SetForce(mgl32.Vec3{-currentForce[0], -currentForce[1], -currentForce[2]})
	chain.SetForce(mgl32.Vec3{-currentForce[0], -currentForce[1], -currentForce[2]})
}

// Start : initialize our values for our game here
func Start(state *geometry.State) {
	fmt.Printf("Started!\n")
	lightToMove, err = scene.GetLightFromScene(state, "pointLight1")
	warningLight1, err = scene.GetLightFromScene(state, "warningLight1")
	warningLight2, err = scene.GetLightFromScene(state, "warningLight2")

	transportCube1, err = scene.GetObjectFromScene(state, "transportCube1")
	chain1, err = scene.GetObjectFromScene(state, "chain1")
	transportHook1, err = scene.GetObjectFromScene(state, "transportHook1")

	transportCube2, err = scene.GetObjectFromScene(state, "transportCube2")
	chain2, err = scene.GetObjectFromScene(state, "chain2")
	transportHook2, err = scene.GetObjectFromScene(state, "transportHook2")

	warning1Reset = false

	if err != nil {
		panic(err)
	}

	transportCube1.SetForce(mgl32.Vec3{-conveyerSpeed, 0.0, 0.0})
	chain1.SetForce(mgl32.Vec3{-conveyerSpeed, 0.0, 0.0})
	transportHook1.SetForce(mgl32.Vec3{-conveyerSpeed, 0.0, 0.0})

	transportCube2.SetForce(mgl32.Vec3{-conveyerSpeed, 0.0, 0.0})
	chain2.SetForce(mgl32.Vec3{-conveyerSpeed, 0.0, 0.0})
	transportHook2.SetForce(mgl32.Vec3{-conveyerSpeed, 0.0, 0.0})

	transportCube2.SetOnCollide(func(box geometry.BoundingBox) {
		TransportCollide(box, transportCube2, chain2, transportHook2)
	})

	transportCube1.SetOnCollide(func(box geometry.BoundingBox) {
		TransportCollide(box, transportCube1, chain1, transportHook1)
	})

}

// Update : runs each frame
func Update(state *geometry.State, deltaTime float64) {
	speed := deltaTime * walkSpeed

	lightToMove.Move = false
	warningLight1.Move = false

	strength := float32(math.Abs(math.Sin(angle)) * 5)
	warningLight1.Strength = strength
	warningLight2.Strength = strength

	// 	rot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	// 	viper.SetRotation(rot)

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
