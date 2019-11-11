package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type State struct {
	objects    []Object
	fragShader string
	vertShader string
	camera     Camera
	lights     []Light
}

type Camera struct {
	up       mgl32.Vec3
	center   mgl32.Vec3
	position mgl32.Vec3
}

type Light struct {
	position []float32
	colour   []float32
	strength float32
}
