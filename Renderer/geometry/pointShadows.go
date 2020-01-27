package geometry

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func CreateCubeDepthMap(state *State, width, height int32) {

	//var depthMapFBO uint32
	gl.GenFramebuffers(1, &state.DepthMapFBO)
	gl.GenTextures(1, &state.DepthCubeMap)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, state.DepthCubeMap)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	var i uint32 = 0
	for i = 0; i < 6; i++ {
		side := gl.TEXTURE_CUBE_MAP_POSITIVE_X + i
		gl.TexImage2D(side, 0, gl.DEPTH_COMPONENT, width, height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)

	}
	//gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, state.DepthMapFBO)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, state.DepthCubeMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)

	//error check the framebuffer
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	if status != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println("ERROR WITH FRAMEBUFFER ", status)
		panic(status)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func CreateLightSpaceTransforms(light Light, near, far float32, width, height float32) []mgl32.Mat4 {
	aspect := float32(width) / float32(height)
	shadowProj := mgl32.Perspective(90*math.Pi/180, aspect, near, far)

	shadowTransforms := []mgl32.Mat4{}
	lightPosition := mgl32.Vec3{light.Position[0], light.Position[1], light.Position[2]}

	//create lookat mat4
	look1Add := lightPosition.Add(mgl32.Vec3{1.0, 0.0, 0.0})
	lookAt1 := mgl32.LookAtV(lightPosition, look1Add, mgl32.Vec3{0.0, -1.0, 0.0})
	final1 := shadowProj.Mul4(lookAt1)
	shadowTransforms = append(shadowTransforms, final1)

	look2Add := lightPosition.Add(mgl32.Vec3{-1.0, 0.0, 0.0})
	lookAt2 := mgl32.LookAtV(lightPosition, look2Add, mgl32.Vec3{0.0, -1.0, 0.0})
	final2 := shadowProj.Mul4(lookAt2)
	shadowTransforms = append(shadowTransforms, final2)

	look3Add := lightPosition.Add(mgl32.Vec3{0.0, 1.0, 0.0})
	lookAt3 := mgl32.LookAtV(lightPosition, look3Add, mgl32.Vec3{0.0, 0.0, 1.0})
	final3 := shadowProj.Mul4(lookAt3)
	shadowTransforms = append(shadowTransforms, final3)

	lookAt4 := mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{0.0, -1.0, 0.0}), mgl32.Vec3{0.0, 0.0, -1.0})
	final4 := shadowProj.Mul4(lookAt4)
	shadowTransforms = append(shadowTransforms, final4)

	lookAt5 := mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{0.0, 0.0, 1.0}), mgl32.Vec3{0.0, -1.0, 0.0})
	final5 := shadowProj.Mul4(lookAt5)
	shadowTransforms = append(shadowTransforms, final5)

	lookAt6 := mgl32.LookAtV(lightPosition, lightPosition.Add(mgl32.Vec3{0.0, 0.0, -1.0}), mgl32.Vec3{0.0, -1.0, 0.0})
	final6 := shadowProj.Mul4(lookAt6)
	shadowTransforms = append(shadowTransforms, final6)

	return shadowTransforms
}
