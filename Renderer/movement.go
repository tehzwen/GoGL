package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//Callbacks for inputs

func KeyHandler(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		keys[key] = true
	} else if action == glfw.Release {
		keys[key] = false
	}
}

func MouseButtonHandler(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		buttons[button] = true
	} else if action == glfw.Release {
		buttons[button] = false
	}
}

func MouseMoveHandler(win *glfw.Window, xPos float64, yPos float64) {
	xDiff := xPos - mouseMovement["X"]
	yDiff := yPos - mouseMovement["Y"]

	mouseMovement["Xmove"] = xDiff
	mouseMovement["Ymove"] = yDiff
	mouseMovement["X"] = xPos
	mouseMovement["Y"] = yPos
	mouseMovement["move"] = 1
}

func MoveForward(state *State, deltaTime float64) {
	forwardVector := state.viewMatrix.Row(2)
	forwardVector = forwardVector.Mul(float32(deltaTime))

	state.camera.position = state.camera.position.Add(mgl32.Vec3{-forwardVector[0], -forwardVector[1], -forwardVector[2]})
	state.camera.center = state.camera.center.Add(mgl32.Vec3{-forwardVector[0], -forwardVector[1], -forwardVector[2]})
}

func MoveBackward(state *State, deltaTime float64) {
	forwardVector := state.viewMatrix.Row(2)
	forwardVector = forwardVector.Mul(float32(deltaTime))

	state.camera.position = state.camera.position.Add(mgl32.Vec3{forwardVector[0], forwardVector[1], forwardVector[2]})
	state.camera.center = state.camera.center.Add(mgl32.Vec3{forwardVector[0], forwardVector[1], forwardVector[2]})
}

func MoveLeft(state *State, deltaTime float64) {
	forwardVector := state.viewMatrix.Row(2)
	newForwardVector := mgl32.Vec3{forwardVector[0], forwardVector[1], forwardVector[2]}
	sideWaysVector := newForwardVector.Cross(state.camera.up)
	sideWaysVector = Normalize(sideWaysVector)

	sideWaysVector = sideWaysVector.Mul(float32(deltaTime))

	state.camera.position = state.camera.position.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
	state.camera.center = state.camera.center.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
}

func MoveRight(state *State, deltaTime float64) {
	forwardVector := state.viewMatrix.Row(2)
	newForwardVector := mgl32.Vec3{-forwardVector[0], -forwardVector[1], -forwardVector[2]}
	sideWaysVector := newForwardVector.Cross(state.camera.up)
	sideWaysVector = Normalize(sideWaysVector)

	sideWaysVector = sideWaysVector.Mul(float32(deltaTime))

	state.camera.position = state.camera.position.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
	state.camera.center = state.camera.center.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
}
