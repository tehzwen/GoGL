package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"./geometry"
	"./mymath"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func MultithreadRender(window *glfw.Window, state *geometry.State) {
	glfw.PollEvents()
	// err := gl.GetError()

	// if err != gl.NO_ERROR {
	// 	fmt.Println(err)
	// }

	//getting the math values for each object using go routines
	for i := 0; i < len(state.Objects); i++ {
		go doObjectMath(state.Objects[i], (*state), objectsToRender)
	}

	//fmt.Println(state.Objects)

	//fmt.Println("SORTED: ", state.Objects)

	//render sequentially using channel values
	state.RenderedObjects = 0
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.MULTISAMPLE)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	tempList := []geometry.RenderObject{}

	for x := 0; x < len(state.Objects); x++ {
		tempList = append(tempList, <-objectsToRender)
	}

	//sort for transparency
	sort.Slice(tempList, func(a, b int) bool {
		//get model info
		nameA, _, _ := tempList[a].CurrentObject.GetDetails()
		nameB, _, _ := tempList[b].CurrentObject.GetDetails()
		distA := tempList[a].DistanceToCamera
		distB := tempList[b].DistanceToCamera

		if distA > distB {
			return true
		} else if distB > distA {
			return false
		} else {
			return nameA > nameB
		}
	})

	// state.ShadowMatrices = geometry.CreateLightSpaceTransforms(state.Lights[0], 1.0, 25.0, 1024, 1024)
	// //fmt.Println(state.ShadowMatrices[0])
	// gl.Viewport(0, 0, 1024, 1024)
	// gl.BindFramebuffer(gl.FRAMEBUFFER, state.DepthMapFBO)
	// gl.Clear(gl.DEPTH_BUFFER_BIT)
	// for x := 0; x < len(tempList); x++ {
	// 	renderObject(state, tempList[x], true)
	// }

	// gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	// gl.Viewport(0, 0, width, height)
	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//normal render here
	for x := 0; x < len(tempList); x++ {
		renderObject(state, tempList[x], false)
	}

	window.SwapBuffers()

}

func doObjectMath(object geometry.Geometry, state geometry.State, objects chan<- geometry.RenderObject) {
	currentProgramInfo, err := object.GetProgramInfo()

	if err != nil {
		//else lets throw an error here
		fmt.Printf("ERROR getting program info!")
	}

	currentShaderProgramInfo, err := object.GetShadowProgramInfo()

	if err != nil {
		//else lets throw an error here
		fmt.Printf("ERROR getting shadow program info!")
	}

	//now get the model
	currentModel, err := object.GetModel()
	if err != nil {
		//throw an error
		fmt.Printf("ERROR getting model!")
	}

	_, _, parent := object.GetDetails()
	currentCentroid := object.GetCentroid()
	currentMaterial := object.GetMaterial()
	currentBuffers := object.GetBuffers()
	shadowBuffers := object.GetShadowBuffers()
	currentVertices := object.GetVertices()

	var fovy = float32(60 * math.Pi / 180)
	var aspect = float32(width / height)
	var near = float32(0.1)
	var far = float32(100.0)

	projection := mgl32.Perspective(fovy, aspect, near, far)
	//create a camfront value
	camFront := state.Camera.Position.Add(state.Camera.Front)
	viewMatrix := mgl32.LookAtV(state.Camera.Position, camFront, state.Camera.Up)
	camPosition := []float32{state.Camera.Position[0], state.Camera.Position[1], state.Camera.Position[2]}
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
	modelMatrix = geometry.ScaleM4(modelMatrix, currentModel.Scale)

	camDist := geometry.VectorDistance(state.Camera.Position, currentModel.Position.Add(currentCentroid)) //calculate this for transparency

	if parent != "" {
		parentObj := geometry.GetSceneObject(parent, state)
		if parentObj != nil {
			parentMMatrix, err := parentObj.GetModelMatrix()
			if err == nil {
				modelMatrix = modelMatrix.Mul4(parentMMatrix)
			}
		} else {
			fmt.Println("ERROR GETTING PARENT OBJECT")
		}
	}

	object.SetModelMatrix(modelMatrix)
	result := geometry.RenderObject{
		ModelMatrix:              modelMatrix,
		ViewMatrix:               viewMatrix,
		ProjMatrix:               projection,
		CameraPosition:           camPosition,
		CurrentBuffers:           currentBuffers,
		CurrentCentroid:          currentCentroid,
		CurrentMaterial:          currentMaterial,
		CurrentModel:             currentModel,
		CurrentProgram:           currentProgramInfo,
		CurrentVertices:          currentVertices,
		CurrentObject:            object,
		CurrentShadowProgramInfo: currentShaderProgramInfo,
		CurrentShadowBuffers:     shadowBuffers,
		DistanceToCamera:         camDist,
	}
	objects <- result
}

func renderObject(state *geometry.State, object geometry.RenderObject, shadowPass bool) {

	var currentProgramInfo geometry.ProgramInfo
	var currentBuffers geometry.ObjectBuffers
	projection := object.ProjMatrix
	viewMatrix := object.ViewMatrix
	camPosition := object.CameraPosition
	modelMatrix := object.ModelMatrix
	currentMaterial := object.CurrentMaterial
	currentVertices := object.CurrentVertices

	//do movement physics and collision testing
	currentForce := object.CurrentObject.GetForce()
	object.CurrentObject.Translate(currentForce)

	if object.CurrentObject.GetBoundingBox().Collide {
		collisionTest(state, object)
	}

	if !shadowPass {
		currentProgramInfo = object.CurrentProgram
		currentBuffers = object.CurrentBuffers
		gl.UseProgram(currentProgramInfo.Program)

	} else {
		currentProgramInfo = object.CurrentShadowProgramInfo
		currentBuffers = object.CurrentShadowBuffers
		gl.UseProgram(currentProgramInfo.Program)

		//fmt.Println(currentProgramInfo)
		//bind the matrices to the shader
		for i := 0; i < 6; i++ {
			gl.UniformMatrix4fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"shadowMatrices[", strconv.Itoa(i)}, "")+"]\x00")), 1, false, &state.ShadowMatrices[i][0])
		}

	}

	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.Projection, 1, false, &projection[0])
	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.View, 1, false, &viewMatrix[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.CameraPosition, 1, &camPosition[0])
	gl.UniformMatrix4fv(currentProgramInfo.UniformLocations.Model, 1, false, &modelMatrix[0])

	model, err := object.CurrentObject.GetModel()
	if err != nil {
		panic(err)
	}

	frustum := mymath.ConstructFrustrum(viewMatrix, projection)
	testLen := object.CurrentObject.GetBoundingBox().Max.LenSqr()
	result := frustum.SphereIntersection(model.Position, testLen)

	if !result {
		return
	}
	diffuseTexture := object.CurrentObject.GetDiffuseTexture()
	normalTexture := object.CurrentObject.GetNormalTexture()

	if diffuseTexture != nil {
		diffuseTexture.Bind(gl.TEXTURE1)
		err := diffuseTexture.SetUniform(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("uDiffuseTexture\x00")))
		if err != nil {
			panic(err)
		}
	}

	if normalTexture != nil {
		normalTexture.Bind(gl.TEXTURE2)
		err := normalTexture.SetUniform(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str("uNormalTexture\x00")))
		if err != nil {
			panic(err)
		}
	}

	// if !shadowPass {
	// 	//bind the cubemap

	// }

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, state.DepthCubeMap)
	gl.Uniform1ui(currentProgramInfo.UniformLocations.DepthMap, state.DepthCubeMap)

	state.RenderedObjects++

	gl.Uniform3fv(currentProgramInfo.UniformLocations.DiffuseVal, 1, &currentMaterial.Diffuse[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.AmbientVal, 1, &currentMaterial.Ambient[0])
	gl.Uniform3fv(currentProgramInfo.UniformLocations.SpecularVal, 1, &currentMaterial.Specular[0])
	gl.Uniform1fv(currentProgramInfo.UniformLocations.NVal, 1, &currentMaterial.N)
	gl.Uniform1fv(currentProgramInfo.UniformLocations.Alpha, 1, &currentMaterial.Alpha)

	if currentMaterial.Alpha < 1.0 {
		//name, _, _ := object.CurrentObject.GetDetails()
		//fmt.Println("here: ", name)
		gl.Enable(gl.BLEND)
		gl.Disable(gl.DEPTH_TEST)
		gl.BlendFunc(gl.ONE_MINUS_CONSTANT_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.ClearDepth(float64(currentMaterial.Alpha))
	} else {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthMask(true)
		gl.Disable(gl.BLEND)
		gl.ClearDepth(1.0)
		gl.DepthFunc(gl.LEQUAL)
	}

	num := int32(len(state.Lights))
	gl.Uniform1iv(currentProgramInfo.UniformLocations.NumLights, 1, &num)

	for i := 0; i < len(state.Lights); i++ {
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].position\x00")), 1, &state.Lights[i].Position[0])
		gl.Uniform3fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].color\x00")), 1, &state.Lights[i].Colour[0])
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].strength\x00")), 1, &state.Lights[i].Strength)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].constant\x00")), 1, &state.Lights[i].Constant)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].linear\x00")), 1, &state.Lights[i].Linear)
		gl.Uniform1fv(gl.GetUniformLocation(currentProgramInfo.Program, gl.Str(strings.Join([]string{"pointLights[", strconv.Itoa(i)}, "")+"].quadratic\x00")), 1, &state.Lights[i].Quadratic)
	}

	state.ViewMatrix = viewMatrix
	gl.BindVertexArray(currentBuffers.Vao)

	if object.CurrentObject.GetType() != "mesh" {
		gl.DrawElements(gl.TRIANGLES, int32(len(currentVertices.Vertices)), gl.UNSIGNED_INT, gl.Ptr(nil))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(currentVertices.Vertices)))
	}

	gl.BindVertexArray(0)

	if diffuseTexture != nil {
		diffuseTexture.UnBind()
	}

	if normalTexture != nil {
		normalTexture.UnBind()
	}

}
