package geometry

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"../shader"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Skybox struct {
	Path        string `json:"path"`
	Format      string `json:"format"`
	Vertices    []float32
	ProgramInfo ProgramInfo
	CubeMap     uint32
	VAO         uint32
}

func LoadCubeMap(path, extension string) uint32 {
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, textureID)

	for i := 0; i < 6; i++ {
		//fmt.Println(path + strconv.Itoa(i) + "." + extension)
		imgFile, err := os.Open(path + strconv.Itoa(i) + "." + extension)
		if err != nil {
			panic(err)
		}

		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			panic(err)
		}
		rgba := image.NewRGBA(img.Bounds())
		internalFmt := int32(gl.SRGB_ALPHA)
		format := uint32(gl.RGBA)
		width := int32(rgba.Rect.Size().X)
		height := int32(rgba.Rect.Size().Y)
		pixType := uint32(gl.UNSIGNED_BYTE)
		dataPtr := gl.Ptr(rgba.Pix)

		draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, internalFmt, width, height, 0, format, pixType, dataPtr)
	}

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)

	return textureID
}

func InitSkyBox(path, extension string, settingsSkyBox *Skybox) {
	skybox := LoadCubeMap(path, extension)

	skyboxVertices := []float32{
		//front face
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,

		//back face
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,

		//top face
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,

		//bottom face
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,

		//side face
		1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,

		//side face
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
	}

	skyboxIndices := []uint32{
		//front face
		2, 0, 1, 3, 0, 2,
		//backface
		5, 4, 6, 6, 4, 7,
		//top face
		10, 9, 8, 10, 8, 11,
		//bottom face
		13, 12, 14, 14, 12, 15,
		//
		18, 16, 17, 18, 17, 19,

		22, 21, 20, 23, 21, 22,
	}

	skyShaderVals := make(map[string]bool)
	skyShaderVals["uViewMatrix"] = true
	skyShaderVals["aPosition"] = true
	skyShaderVals["uProjectionMatrix"] = true
	skyShaderVals["skybox"] = true
	skyShader := &shader.SkyboxShader{}
	skyShader.Setup()
	skyShaderProgramInfo := ProgramInfo{}
	skyShaderProgramInfo.Program = InitOpenGL(skyShader.GetVertShader(), skyShader.GetFragShader(), skyShader.GetGeometryShader())
	skyShaderAttribs := Attributes{}
	skyShaderAttribs.SetPosition(0)
	skyShaderProgramInfo.SetAttributes(skyShaderAttribs)
	SetupAttributesMap(&skyShaderProgramInfo, skyShaderVals)

	//draw the skybox
	skyboxVAO := CreateTriangleVAO(&skyShaderProgramInfo, skyboxVertices, nil, nil, nil, nil, skyboxIndices)

	settingsSkyBox.CubeMap = skybox
	settingsSkyBox.VAO = skyboxVAO
	settingsSkyBox.ProgramInfo = skyShaderProgramInfo
	settingsSkyBox.Vertices = skyboxVertices
}
