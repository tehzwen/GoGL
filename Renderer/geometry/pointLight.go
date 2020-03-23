package geometry

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// PointLight - struct for a pointlight in the scene
type PointLight struct {
	Name              string    `json:"name"`
	Position          []float32 `json:"position"`
	Parent            string    `json:"parent"`
	Colour            []float32 `json:"colour"`
	Strength          float32   `json:"strength"`
	Quadratic         float32   `json:"quadratic"`
	Linear            float32   `json:"linear"`
	Constant          float32   `json:"constant"`
	FarPlane          float32   `json:"farPlane"`
	NearPlane         float32   `json:"nearPlane"`
	Shadow            int32     `json:"shadow"`
	DepthMap          uint32
	LightViewMatrices []mgl32.Mat4
	Move              bool
}

func (light *PointLight) CreateCubeDepthMap(width, height int32) {
	var depthCubeMap uint32
	gl.GenTextures(1, &depthCubeMap)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, depthCubeMap)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, gl.DEPTH_COMPONENT, 1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_X, 0, gl.DEPTH_COMPONENT, 1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Y, 0, gl.DEPTH_COMPONENT, 1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, 0, gl.DEPTH_COMPONENT, 1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_Z, 0, gl.DEPTH_COMPONENT, 1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexImage2D(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, 0, gl.DEPTH_COMPONENT, 1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	light.DepthMap = depthCubeMap
}

func (light *PointLight) BindDepthMap(state *State) {
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, 0, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, state.DepthFBO)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, light.DepthMap, 0)
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

func (light *PointLight) CreateLightSpaceTransforms(width, height float32) {
	aspect := float32(width) / float32(height)
	shadowProj := mgl32.Perspective(90*math.Pi/180, aspect, light.NearPlane, light.FarPlane)

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

	light.LightViewMatrices = shadowTransforms
}

func (light *PointLight) ShadowRender(state *State, object Geometry, shadowProgramInfo *ProgramInfo) {
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

	if object.GetType() == "mesh" {
		//move to centroid
		centroidMat := mgl32.Translate3D(currentCentroid[0], currentCentroid[1], currentCentroid[2])
		modelMatrix = modelMatrix.Mul4(centroidMat)
		//position
		positionMat := mgl32.Translate3D(currentModel.Position[0], currentModel.Position[1], currentModel.Position[2])
		modelMatrix = modelMatrix.Mul4(positionMat)
		//negative centroid
		negCent := mgl32.Translate3D(-currentCentroid[0], -currentCentroid[1], -currentCentroid[2])
		modelMatrix = modelMatrix.Mul4(negCent)
		//rotation
		modelMatrix = modelMatrix.Mul4(currentModel.Rotation)
		//scale
		modelMatrix = ScaleM4(modelMatrix, currentModel.Scale)
	} else {
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
	}

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
	if light.Move {
		light.CreateLightSpaceTransforms(1024, 1024)
	}

	for i := 0; i < 6; i++ {
		gl.UniformMatrix4fv(gl.GetUniformLocation(shadowProgramInfo.Program, gl.Str(strings.Join([]string{"shadowMatrices[", strconv.Itoa(i)}, "")+"]\x00")), 1, false, &light.LightViewMatrices[i][0])
	}

	gl.UniformMatrix4fv(gl.GetUniformLocation(shadowProgramInfo.Program, gl.Str("uModelMatrix\x00")), 1, false, &modelMatrix[0])
	gl.Uniform3fv(gl.GetUniformLocation(shadowProgramInfo.Program, gl.Str("lightPos\x00")), 1, &light.Position[0])
	gl.Uniform1fv(gl.GetUniformLocation(shadowProgramInfo.Program, gl.Str("farPlane\x00")), 1, &light.FarPlane)
	gl.BindVertexArray(currentBuffers.Vao)

	if object.GetType() != "mesh" {
		gl.DrawElements(gl.TRIANGLES, int32(len(currentVertices.Vertices)), gl.UNSIGNED_INT, gl.Ptr(nil))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(currentVertices.Vertices)))
	}
	gl.BindVertexArray(0)
}
