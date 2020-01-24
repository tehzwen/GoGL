package geometry

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func CreateCubeDepthMap(width, height int32) uint32 {

	var depthMapFBO uint32
	gl.GenFramebuffers(1, &depthMapFBO)

	var cubeMap uint32
	gl.GenTextures(1, &cubeMap)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeMap)
	for i := 0; i < 6; i++ {
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.DEPTH_COMPONENT, width, height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_COMPONENT, gl.TEXTURE_2D, cubeMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return depthMapFBO
}

func CreateLightSpaceTransforms(light PointLight, near, far float32, width, height int32) {
	aspect := float32(width) / float32(height)
	shadowProj := mgl32.Perspective(float32(60*math.Pi/180), aspect, near, far)

	shadowTransforms := []mgl32.Mat4{}
	lightPosition := mgl32.Vec3{light.Position[0], light.Position[1], light.Position[2]}

	shadowProjVal := shadowProj.Mul4x1(mgl32.Vec4{light.Position[0], light.Position[1], light.Position[2], 1})

	shadowTransforms = append(shadowTransforms, mgl32.LookAtV(mgl32.Vec3{shadowProjVal[0], shadowProjVal[1], shadowProjVal[2]}, lightPosition.Add(mgl32.Vec3{1.0, 0.0, 0.0}), mgl32.Vec3{0.0, -1.0, 0.0}))
	shadowTransforms = append(shadowTransforms, mgl32.LookAtV(mgl32.Vec3{shadowProjVal[0], shadowProjVal[1], shadowProjVal[2]}, lightPosition.Add(mgl32.Vec3{-1.0, 0.0, 0.0}), mgl32.Vec3{0.0, -1.0, 0.0}))
	shadowTransforms = append(shadowTransforms, mgl32.LookAtV(mgl32.Vec3{shadowProjVal[0], shadowProjVal[1], shadowProjVal[2]}, lightPosition.Add(mgl32.Vec3{0.0, -1.0, 0.0}), mgl32.Vec3{0.0, -1.0, 0.0}))
	shadowTransforms = append(shadowTransforms, mgl32.LookAtV(mgl32.Vec3{shadowProjVal[0], shadowProjVal[1], shadowProjVal[2]}, lightPosition.Add(mgl32.Vec3{0.0, 1.0, 0.0}), mgl32.Vec3{0.0, -1.0, 0.0}))
	shadowTransforms = append(shadowTransforms, mgl32.LookAtV(mgl32.Vec3{shadowProjVal[0], shadowProjVal[1], shadowProjVal[2]}, lightPosition.Add(mgl32.Vec3{0.0, 0.0, 1.0}), mgl32.Vec3{0.0, -1.0, 0.0}))
	shadowTransforms = append(shadowTransforms, mgl32.LookAtV(mgl32.Vec3{shadowProjVal[0], shadowProjVal[1], shadowProjVal[2]}, lightPosition.Add(mgl32.Vec3{1.0, 0.0, -1.0}), mgl32.Vec3{0.0, -1.0, 0.0}))
}
