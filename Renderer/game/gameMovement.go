package game

import (
	"../geometry"
	"github.com/go-gl/mathgl/mgl32"
)

func translate(object geometry.Geometry, translateVector mgl32.Vec3) {
	object.Translate(translateVector)
}

func MoveForward(state *geometry.State, deltaTime float64) {
	forwardVector := state.ViewMatrix.Row(2)
	forwardVector = forwardVector.Mul(float32(deltaTime))

	state.Camera.Position = state.Camera.Position.Add(mgl32.Vec3{-forwardVector[0], -forwardVector[1], -forwardVector[2]})
}

func MoveBackward(state *geometry.State, deltaTime float64) {
	forwardVector := state.ViewMatrix.Row(2)
	forwardVector = forwardVector.Mul(float32(deltaTime))

	state.Camera.Position = state.Camera.Position.Add(mgl32.Vec3{forwardVector[0], forwardVector[1], forwardVector[2]})
}

func MoveLeft(state *geometry.State, deltaTime float64) {
	forwardVector := state.ViewMatrix.Row(2)
	newForwardVector := mgl32.Vec3{forwardVector[0], forwardVector[1], forwardVector[2]}
	sideWaysVector := newForwardVector.Cross(state.Camera.Up)
	sideWaysVector = geometry.Normalize(sideWaysVector)

	sideWaysVector = sideWaysVector.Mul(float32(deltaTime))

	state.Camera.Position = state.Camera.Position.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
	state.Camera.Center = state.Camera.Center.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
}

func MoveRight(state *geometry.State, deltaTime float64) {
	forwardVector := state.ViewMatrix.Row(2)
	newForwardVector := mgl32.Vec3{-forwardVector[0], -forwardVector[1], -forwardVector[2]}
	sideWaysVector := newForwardVector.Cross(state.Camera.Up)
	sideWaysVector = geometry.Normalize(sideWaysVector)

	sideWaysVector = sideWaysVector.Mul(float32(deltaTime))

	state.Camera.Position = state.Camera.Position.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
	state.Camera.Center = state.Camera.Center.Add(mgl32.Vec3{sideWaysVector[0], sideWaysVector[1], sideWaysVector[2]})
}
