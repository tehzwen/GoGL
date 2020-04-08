package game

import (
	"fmt"

	"../geometry"
	"github.com/go-gl/glfw/v3.1/glfw"
	_ "github.com/go-gl/mathgl/mgl32"
)

var angle = 0.0
var walkSpeed float64 = 5
var runSpeed float64 = 5
var lightSpeed float32 = 0.3

//var lightToMove *geometry.PointLight

// var lights []*geometry.PointLight
// var lightNames = [...]string{"pointLight1"}

// var dragon geometry.Geometry
// var err error

// var objects []geometry.Geometry
// //var objectNames = [...]string{"dragonReflect", "dragonRefract", "bunnyDiffuse", "reflectCube", "refractCube", "normalCube", "doomguy", "bunnyFlat", "flatCube", "diffuseCube"}
// var objectNames = [...]string{"dragonReflect", "dragonRefract"}

// Start : initialize our values for our game here
func Start(state *geometry.State) {
	fmt.Printf("Started!\n")

	// 	for i := 0; i < len(objectNames); i++ {
	// 		objects = append(objects, scene.GetObjectFromScene(state, objectNames[i]))
	// 	}

	//     for i := 0; i < len(lightNames); i++ {
	//         lights = append(lights, scene.GetLightFromScene(state, lightNames[i]))
	//     }

	// // 	lightToMove = scene.GetLightFromScene(state, "pointLight1")
	// // 	lightToMove2 = scene.GetLightFromScene(state, "pointLight2")
	// // 	lightToMove3 = scene.GetLightFromScene(state, "pointLight3")
	// 	//dragon = scene.GetObjectFromScene(state, "dragon")

}

// Update : runs each frame
func Update(state *geometry.State, deltaTime float64) {
	speed := deltaTime * walkSpeed

	// rot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	// //rotate the objects
	// for i := 0; i < len(objects); i++ {
	//     objects[i].SetRotation(rot)
	// }

	//rot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	//dragon.SetRotation(rot)

	// 	if state.Keys[glfw.KeyQ] {

	// 	    for i := 0; i < len(lights); i++ {
	// 	        lights[i].Position[2] += lightSpeed
	// 	        lights[i].Move = true
	// 	    }

	// // 		lightToMove.Position[2] += 0.1
	// // 		lightToMove.Move = true
	// 	}

	// 	if state.Keys[glfw.KeyR] {
	// // 		lightToMove.Position[2] -= 0.1
	// // 		lightToMove.Move = true
	// 		for i := 0; i < len(lights); i++ {
	// 	        lights[i].Position[2] -= lightSpeed
	// 	        lights[i].Move = true
	// 	    }
	// 	}

	// 	if state.Keys[glfw.Key3] {
	// // 		lightToMove.Position[0] += 0.1
	// // 		lightToMove.Move = true
	// 		for i := 0; i < len(lights); i++ {
	// 	        lights[i].Position[0] += lightSpeed
	// 	        lights[i].Move = true
	// 	    }
	// 	}

	// 	if state.Keys[glfw.Key1] {
	// // 		lightToMove.Position[0] -= 0.1
	// // 		lightToMove.Move = true
	// 		for i := 0; i < len(lights); i++ {
	// 	        lights[i].Position[0] -= lightSpeed
	// 	        lights[i].Move = true
	// 	    }
	// 	}

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
