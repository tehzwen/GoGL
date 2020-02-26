package geometry

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// State : struct for holding scene state
type State struct {
	Objects           []Geometry
	FragShader        string
	VertShader        string
	Camera            Camera
	PointLights       []PointLight
	DirectionalLights []DirectionalLight
	ViewMatrix        mgl32.Mat4
	Keys              map[glfw.Key]bool
	LoadedObjects     int
	RenderedObjects   int
	ShadowMatrices    []mgl32.Mat4
	CurrentTexUnit    uint32
	DepthFBO          uint32
}

// Camera : struct for holding info about the camera
type Camera struct {
	Up       mgl32.Vec3
	Center   mgl32.Vec3
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Pitch    float32
	Yaw      float32
	Roll     float32
}
