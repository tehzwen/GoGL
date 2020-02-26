package geometry

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type DirectionalLight struct {
	Name            string    `json:"name"`
	Parent          string    `json:"parent"`
	Colour          []float32 `json:"colour"`
	Strength        float32   `json:"strength"`
	Direction       []float32 `json:"direction"`
	Position        []float32 `json:"position"`
	DepthMap        uint32
	LightViewMatrix mgl32.Mat4
}

func (light *DirectionalLight) CreateDirectionalDepthMap(width, height int32) {
	var depthMap uint32
	gl.GenTextures(1, &depthMap)
	gl.BindTexture(gl.TEXTURE_2D, depthMap)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, width, height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	light.DepthMap = depthMap
}

func (light *DirectionalLight) BindDepthMap(state *State) {
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, 0, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, state.DepthFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, light.DepthMap, 0)
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

func (light *DirectionalLight) CreateLightSpaceTransforms(near, far float32) {
	lightProj := mgl32.Ortho(-10.0, 10.0, -10.0, 10.0, near, far)
	lightView := mgl32.LookAtV(
		mgl32.Vec3{light.Position[0], light.Position[1], light.Position[2]},
		mgl32.Vec3{light.Direction[0], light.Direction[1], light.Direction[2]},
		mgl32.Vec3{0.0, 1.0, 0.0})
	lightSpace := lightProj.Mul4(lightView)
	light.LightViewMatrix = lightSpace
}

func (light *DirectionalLight) ShadowRender(state *State, object Geometry, shadowProgramInfo *ProgramInfo) {
	gl.UseProgram(shadowProgramInfo.Program)
	currentModel, err := object.GetModel()

	if err != nil {
		//throw an error
		fmt.Printf("ERROR getting model!")
	}
	_, _, parent := object.GetDetails()
	currentCentroid := object.GetCentroid()
	currentVertices := object.GetVertices()
	currentBuffers := object.GetBuffers()
	modelMatrix := mgl32.Ident4()

	//move to centroid
	centroidMat := mgl32.Translate3D(currentCentroid[0], currentCentroid[1], currentCentroid[2])
	modelMatrix = modelMatrix.Mul4(centroidMat)

	//rotation
	modelMatrix = modelMatrix.Mul4(currentModel.Rotation)

	//position
	positionMat := mgl32.Translate3D(currentModel.Position[0], currentModel.Position[1], currentModel.Position[2])
	modelMatrix = modelMatrix.Mul4(positionMat)

	//negative centroid
	negCent := mgl32.Translate3D(-currentCentroid[0], -currentCentroid[1], -currentCentroid[2])
	modelMatrix = modelMatrix.Mul4(negCent)

	//scale
	modelMatrix = ScaleM4(modelMatrix, currentModel.Scale)

	if parent != "" {
		parentObj := GetSceneObject(parent, (*state))
		if parentObj != nil {
			parentMMatrix, err := parentObj.GetModelMatrix()
			if err == nil {
				modelMatrix = modelMatrix.Mul4(parentMMatrix)
			}
		} else {
			fmt.Println("ERROR GETTING PARENT OBJECT")
		}
	}
	// if light.Move {
	// 	light.CreateLightSpaceTransforms(0.5, 25, 1024, 1024)
	// }

	light.CreateLightSpaceTransforms(0.5, 25.0)
	gl.UniformMatrix4fv(gl.GetUniformLocation(shadowProgramInfo.Program, gl.Str("uModelMatrix\x00")), 1, false, &modelMatrix[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(shadowProgramInfo.Program, gl.Str("lightSpaceMatrix\x00")), 1, false, &light.LightViewMatrix[0])
	gl.BindVertexArray(currentBuffers.Vao)

	if object.GetType() != "mesh" {
		gl.DrawElements(gl.TRIANGLES, int32(len(currentVertices.Vertices)), gl.UNSIGNED_INT, gl.Ptr(nil))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(currentVertices.Vertices)))
	}
	gl.BindVertexArray(0)
}
