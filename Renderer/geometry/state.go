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
	Settings          Settings
}

// Camera : struct for holding info about the camera
type Camera struct {
	Name     string     `json:"name"`
	Up       mgl32.Vec3 `json:"up"`
	Position mgl32.Vec3 `json:"position"`
	Front    mgl32.Vec3 `json:"front"`
	Pitch    float32    `json:"pitch"`
	Yaw      float32    `json:"yaw"`
	Roll     float32    `json:"roll"`
}
