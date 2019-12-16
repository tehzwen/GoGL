package geometry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"../parser"

	"github.com/go-gl/mathgl/mgl32"
)

type SceneObject struct {
	Name           string    `json:"name"`
	Material       Material  `json:"material"`
	ObjectType     string    `json:"type"`
	Position       []float32 `json:"position"`
	Scale          []float32 `json:"scale"`
	DiffuseTexture string    `json:"diffuseTexture"`
	NormalTexture  string    `json:"normalTexture"`
	Parent         string    `json:"parent"`
	Model          string    `json:"model"`
}

type SceneLight struct {
	Name           string    `json:"name"`
	Material       Material  `json:"material"`
	ObjectType     string    `json:"type"`
	Position       []float32 `json:"position"`
	Scale          []float32 `json:"scale"`
	DiffuseTexture string    `json:"diffuseTexture"`
	NormalTexture  string    `json:"normalTexture"`
	Parent         string    `json:"parent"`
	Model          string    `json:"model"`
	Colour         []float32 `json:"colour"`
	Strength       float32   `json:"strength"`
}

type Settings struct {
}

type Scene struct {
	Objects  []SceneObject `json:"objects"`
	Lights   []SceneLight  `json:"lights"`
	Settings Settings      `json:"settings"`
}

// RotateY - rotates a vec3 by an angle around another vec3
func RotateY(a mgl32.Vec3, b mgl32.Vec3, c float64) mgl32.Vec3 {
	p := mgl32.Vec3{0, 0, 0}
	r := mgl32.Vec3{0, 0, 0}

	p[0] = a[0] - b[0]
	p[1] = a[1] - b[1]
	p[2] = a[2] - b[2]

	r[0] = p[2]*float32(math.Sin(c)) + p[0]*float32(math.Cos(c))
	r[1] = p[1]
	r[2] = p[2]*float32(math.Cos(c)) - p[0]*float32(math.Sin(c))

	return mgl32.Vec3{r[0] + b[0], r[1] + b[1], r[2] + b[2]}
}

// Normalize - normalizes a vec3
func Normalize(a mgl32.Vec3) mgl32.Vec3 {
	x := a[0]
	y := a[1]
	z := a[2]

	len := x*x + y*y + z*z

	if len > 0 {
		newLen := 1 / math.Sqrt(float64(len))
		len = float32(newLen)
	}

	return mgl32.Vec3{a[0] * len, a[1] * len, a[2] * len}
}

// ScaleM4 - function used for scaling matrix 4 by vec3
func ScaleM4(a mgl32.Mat4, v mgl32.Vec3) mgl32.Mat4 {
	x := v[0]
	y := v[1]
	z := v[2]

	out := mgl32.Mat4{
		a[0] * x,
		a[1] * x,
		a[2] * x,
		a[3] * x,
		a[4] * y,
		a[5] * y,
		a[6] * y,
		a[7] * y,
		a[8] * z,
		a[9] * z,
		a[10] * z,
		a[11] * z,
		a[12],
		a[13],
		a[14],
		a[15],
	}

	return out
}

func addObjectToState(object Geometry, state *State, sceneObj SceneObject) {
	err := object.SetShader(state.VertShader, state.FragShader)

	if err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		//create model
		tempModel := Model{
			Position: mgl32.Vec3{sceneObj.Position[0], sceneObj.Position[1], sceneObj.Position[2]},
			Scale:    mgl32.Vec3{sceneObj.Scale[0], sceneObj.Scale[1], sceneObj.Scale[2]},
			Rotation: mgl32.Ident4(),
		}
		object.Setup(
			sceneObj.Material,
			tempModel,
			sceneObj.Name,
		)
		state.Objects = append(state.Objects, object)
	}
}

func GetBoundingBox(vertices []float32) BoundingBox {
	min := mgl32.Vec3{0, 0, 0}
	max := mgl32.Vec3{0, 0, 0}

	for i := 0; i < len(vertices); i += 3 {
		if vertices[i] < min[0] {
			min[0] = vertices[i]
		}

		if vertices[i] > max[0] {
			max[0] = vertices[i]
		}

		if vertices[i+1] < min[1] {
			min[1] = vertices[i+1]
		}

		if vertices[i+1] > max[1] {
			max[1] = vertices[i+1]
		}

		if vertices[i+2] < min[2] {
			min[2] = vertices[i+2]
		}

		if vertices[i+2] > max[2] {
			max[2] = vertices[i+2]
		}
	}

	result := BoundingBox{}
	result.Max = max
	result.Min = min

	return result
}

func ScaleBoundingBox(box BoundingBox, scaleVec mgl32.Vec3) BoundingBox {
	result := BoundingBox{}
	result.Min = box.Min

	for i := 0; i < 3; i++ {
		if box.Max[i] == 0 && scaleVec[i] > 1 {
			result.Max[i] = box.Max[i] + scaleVec[i]
		} else {
			result.Max[i] = box.Max[i] * scaleVec[i]
		}
	}

	return result
}

func TranslateBoundingBox(box BoundingBox, translateVec mgl32.Vec3) BoundingBox {
	result := BoundingBox{}

	result.Min[0] = box.Min[0] + translateVec[0]
	result.Max[0] = box.Max[0] + translateVec[0]
	result.Min[1] = box.Min[1] + translateVec[1]
	result.Max[1] = box.Max[1] + translateVec[1]
	result.Min[2] = box.Min[2] + translateVec[2]
	result.Max[2] = box.Max[2] + translateVec[2]

	return result
}

func ParseJsonFile(filePath string, state *State) {
	fmt.Printf("Opening scene file: %s\n", filePath)

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)

	jsonFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var scene []Scene

	fmt.Println("Starting scene file read.....")
	err = json.Unmarshal(byteValue, &scene)
	fmt.Println("Reading scene file complete.....")

	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(scene[0].Objects); i++ {
		if scene[0].Objects[i].ObjectType == "cube" {
			fmt.Println(scene[0].Objects[i].Name, " loading....")
			tempCube := Cube{}
			addObjectToState(&tempCube, state, scene[0].Objects[i])
			state.LoadedObjects++
			fmt.Println(scene[0].Objects[i].Name, " loaded successfully!")

		} else if scene[0].Objects[i].ObjectType == "plane" {
			fmt.Println(scene[0].Objects[i].Name, " loading....")
			tempPlane := Plane{}
			addObjectToState(&tempPlane, state, scene[0].Objects[i])
			state.LoadedObjects++
			fmt.Println(scene[0].Objects[i].Name, " loaded successfully!")

		} else if scene[0].Objects[i].ObjectType == "mesh" {
			fmt.Println(scene[0].Objects[i].Name, " loading....")
			meshPath := exPath + "/../Editor/" + scene[0].Objects[i].Model
			tempMeshVals := parser.Parse(meshPath)

			tempModelObject := ModelObject{}

			err := tempModelObject.SetShader(state.VertShader, state.FragShader)

			if err != nil {
				fmt.Println(err)
				panic(err)
			} else {
				tempModelObject.SetVertexValues(tempMeshVals.Vertices, tempMeshVals.Normals, tempMeshVals.UVs)
				tempModel := Model{
					Position: mgl32.Vec3{scene[0].Objects[i].Position[0], scene[0].Objects[i].Position[1], scene[0].Objects[i].Position[2]},
					Scale:    mgl32.Vec3{scene[0].Objects[i].Scale[0], scene[0].Objects[i].Scale[1], scene[0].Objects[i].Scale[2]},
					Rotation: mgl32.Ident4(),
				}

				tempModelObject.Setup(
					scene[0].Objects[i].Material,
					tempModel,
					scene[0].Objects[i].Name,
				)

				state.Objects = append(state.Objects, &tempModelObject)
				state.LoadedObjects++
				fmt.Println(scene[0].Objects[i].Name, " loaded successfully!")
			}
		}
	}

	for j := 0; j < len(scene[0].Lights); j++ {
		tempLight := Light{
			Colour:   scene[0].Lights[j].Colour,
			Strength: scene[0].Lights[j].Strength,
			Position: scene[0].Lights[j].Position,
		}
		state.Lights = append(state.Lights, tempLight)
	}
}
