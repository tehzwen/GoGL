package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

// State : struct for holding scene state
type State struct {
	objects    []Geometry
	fragShader string
	vertShader string
	camera     Camera
	lights     []Light
	viewMatrix mgl32.Mat4
}

// Camera : struct for holding info about the camera
type Camera struct {
	up       mgl32.Vec3
	center   mgl32.Vec3
	position mgl32.Vec3
}

// Light : struct for lights in the scene
type Light struct {
	position []float32
	colour   []float32
	strength float32
}
