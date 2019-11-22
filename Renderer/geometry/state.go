package geometry

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// State : struct for holding scene state
type State struct {
	Objects    []Geometry
	FragShader string
	VertShader string
	Camera     Camera
	Lights     []Light
	ViewMatrix mgl32.Mat4
	Keys       map[glfw.Key]bool
}

// Camera : struct for holding info about the camera
type Camera struct {
	Up       mgl32.Vec3
	Center   mgl32.Vec3
	Position mgl32.Vec3
}

// Light : struct for lights in the scene
type Light struct {
	Position []float32
	Colour   []float32
	Strength float32
}
