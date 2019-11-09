package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type State struct {
	objects    []Object
	fragShader string
	vertShader string
	camera     Camera
}

type Camera struct {
	up       mgl32.Vec3
	center   mgl32.Vec3
	position mgl32.Vec3
}
