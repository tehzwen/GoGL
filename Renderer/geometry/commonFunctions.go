package geometry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"../common"
	"../parser"

	"github.com/go-gl/mathgl/mgl32"
)

// RenderObject - Struct used for applying math in coroutine and used for rendering
type RenderObject struct {
	ViewMatrix               mgl32.Mat4
	ProjMatrix               mgl32.Mat4
	ModelMatrix              mgl32.Mat4
	CurrentBuffers           ObjectBuffers
	CurrentModel             Model
	CurrentMaterial          Material
	CurrentCentroid          mgl32.Vec3
	CurrentProgram           ProgramInfo
	CameraPosition           []float32
	CurrentVertices          VertexValues
	CurrentObject            Geometry
	DistanceToCamera         float32
	CurrentShadowProgramInfo ProgramInfo
	CurrentShadowBuffers     ObjectBuffers
}

// SceneObject - Object used for reading in from JSON scene file
type SceneObject struct {
	Name            string    `json:"name"`
	Material        Material  `json:"material"`
	ObjectType      string    `json:"type"`
	Position        []float32 `json:"position"`
	Scale           []float32 `json:"scale"`
	Rotation        []float32 `json:"rotation"`
	DiffuseTexture  string    `json:"diffuseTexture"`
	NormalTexture   string    `json:"normalTexture"`
	Parent          string    `json:"parent"`
	Model           string    `json:"model"`
	Collide         bool      `json:"collide"`
	Reflective      int       `json:"reflective"`
	RefractionIndex float32   `json:"refractionIndex"`
}

// Settings - WIP
type Settings struct {
	Cam             Camera    `json:"camera"`
	BackgroundColor []float32 `json:"backgroundColor"`
	Skybox          Skybox    `json:"skybox"`
}

// Scene - Struct for holding allthe info about the current scene
type Scene struct {
	Objects           []SceneObject      `json:"objects"`
	PointLights       []PointLight       `json:"pointLights"`
	DirectionalLights []DirectionalLight `json:"directionalLights"`
	Settings          Settings           `json:"settings"`
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

	//get rotation
	rot := CreateMat4FromArray(sceneObj.Rotation)

	//create model
	tempModel := Model{
		Scale:    mgl32.Vec3{sceneObj.Scale[0], sceneObj.Scale[1], sceneObj.Scale[2]},
		Rotation: rot,
	}

	if sceneObj.DiffuseTexture != "" {
		sceneObj.Material.DiffuseTexture = sceneObj.DiffuseTexture
	}

	if sceneObj.NormalTexture != "" {
		sceneObj.Material.NormalTexture = sceneObj.NormalTexture
	}

	object.Setup(
		sceneObj.Material,
		tempModel,
		sceneObj.Name,
		sceneObj.Collide,
		sceneObj.Reflective,
		sceneObj.RefractionIndex,
	)
	object.Translate(mgl32.Vec3{sceneObj.Position[0], sceneObj.Position[1], sceneObj.Position[2]})
	state.Objects = append(state.Objects, object)
}

// GetBoundingBox - Given a set of vertices, returns a bounding box object that contains the min & max of the box
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

// ScaleBoundingBox - Scales an existing bounding box by a vector, used when scaling an object to keep collision detection in scale
func ScaleBoundingBox(box BoundingBox, scaleVec mgl32.Vec3) BoundingBox {
	result := BoundingBox{
		Collide:        box.Collide,
		CollisionCount: box.CollisionCount,
		CollisionBody:  box.CollisionBody,
	}
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

// TranslateBoundingBox - Translates a bounding box by a vec3
func TranslateBoundingBox(box BoundingBox, translateVec mgl32.Vec3) BoundingBox {
	result := BoundingBox{
		Collide:        box.Collide,
		CollisionCount: box.CollisionCount,
		CollisionBody:  box.CollisionBody,
	}

	result.Min[0] = box.Min[0] + translateVec[0]
	result.Max[0] = box.Max[0] + translateVec[0]
	result.Min[1] = box.Min[1] + translateVec[1]
	result.Max[1] = box.Max[1] + translateVec[1]
	result.Min[2] = box.Min[2] + translateVec[2]
	result.Max[2] = box.Max[2] + translateVec[2]

	return result
}

// Intersect - Intersection function for checking if two bounding boxes are intersecting at all
func Intersect(a, b BoundingBox) bool {
	return (a.Min[0] <= b.Max[0] && a.Max[0] >= b.Min[0]) &&
		(a.Min[1] <= b.Max[1] && a.Max[1] >= b.Min[1]) &&
		(a.Min[2] <= b.Max[2] && a.Max[2] >= b.Min[2])
}

// ParseJSONFile - Given a json file that contains scene information, load it and put into global state
func ParseJSONFile(filePath string, state *State) {
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
		panic(err)
	}

	state.Settings = scene[0].Settings

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
			//check here if a cached version of the mesh already exists
			fmt.Println("CHECK THIS FILE ", scene[0].Objects[i].Model+".dat")

			var objects []parser.OBJObject

			//if there is no cached value
			if _, err := os.Stat("./game/.cache/" + scene[0].Objects[i].Model + ".dat"); os.IsNotExist(err) {
				meshPath := exPath + "/../Editor/models/" + scene[0].Objects[i].Model
				objects = parser.Parse(meshPath)
				b64Objects := parser.SerializeOBJ(objects)
				common.WriteB64("./game/.cache/"+scene[0].Objects[i].Model+".dat", b64Objects)
			} else {
				val := common.ReadB64("./game/.cache/" + scene[0].Objects[i].Model + ".dat")
				objects = parser.DeserializeOBJ(val)
			}
			for x := 0; x < len(objects); x++ {
				for j := 0; j < len(objects[x].Materials); j++ {
					tempModelObject := ModelObject{
						MTLPresent: true,
					}
					var parsedMaterial parser.ParsedMaterial

					//check for regular texture first
					if scene[0].Objects[i].DiffuseTexture != "" {
						tempModelObject.MTLPresent = false
						parsedMaterial.Kd = scene[0].Objects[i].Material.Diffuse
						parsedMaterial.Ka = scene[0].Objects[i].Material.Ambient
						parsedMaterial.Ks = scene[0].Objects[i].Material.Specular
						parsedMaterial.Ns = scene[0].Objects[i].Material.N
						parsedMaterial.D = scene[0].Objects[i].Material.Alpha
						parsedMaterial.MapKD = scene[0].Objects[i].DiffuseTexture

					} else {
						tempMaterial, err := parser.ParseMTLFile(objects[x].Materials[j].MTLLib, objects[x].Materials[j].Name)
						if err == nil {
							parsedMaterial = tempMaterial
						} else {
							parsedMaterial.Kd = scene[0].Objects[i].Material.Diffuse
							parsedMaterial.Ka = scene[0].Objects[i].Material.Ambient
							parsedMaterial.Ks = scene[0].Objects[i].Material.Specular
							parsedMaterial.Ns = scene[0].Objects[i].Material.N
							parsedMaterial.D = scene[0].Objects[i].Material.Alpha
							parsedMaterial.MapKD = scene[0].Objects[i].DiffuseTexture
						}
					}

					if len(objects[x].Geometry.UVs) > 0 {
						tempModelObject.SetVertexValues(objects[x].Geometry.Vertices[objects[x].Materials[j].Start*3:objects[x].Materials[j].End*3],
							objects[x].Geometry.Normals[objects[x].Materials[j].Start*3:objects[x].Materials[j].End*3],
							objects[x].Geometry.UVs[objects[x].Materials[j].Start*2:objects[x].Materials[j].End*2], nil)
					} else {
						tempModelObject.SetVertexValues(objects[x].Geometry.Vertices[objects[x].Materials[j].Start*3:objects[x].Materials[j].End*3],
							objects[x].Geometry.Normals[objects[x].Materials[j].Start*3:objects[x].Materials[j].End*3],
							nil, nil)
					}

					tempName := scene[0].Objects[i].Name
					rot := CreateMat4FromArray(scene[0].Objects[i].Rotation)
					tempModel := Model{
						Position: mgl32.Vec3{scene[0].Objects[i].Position[0], scene[0].Objects[i].Position[1], scene[0].Objects[i].Position[2]},
						Scale:    mgl32.Vec3{scene[0].Objects[i].Scale[0], scene[0].Objects[i].Scale[1], scene[0].Objects[i].Scale[2]},
						Rotation: mgl32.Ident4(),
					}

					if x > 1 {
						tempName = strings.Join([]string{tempName, strconv.Itoa(j)}, "")
						tempModelObject.SetParent(scene[0].Objects[j].Name)
						// newObject.name = object.name + i;
						//     newObject.parent = object.name;
						//     newObject.parentTransform = object.position;
					}

					if j > 0 {
						tempName = strings.Join([]string{tempName, strconv.Itoa(j)}, "")
						tempModelObject.SetParent(scene[0].Objects[i].Name)
						tempModel.Position = mgl32.Vec3{0, 0, 0}
						tempModel.Scale = mgl32.Vec3{1, 1, 1}
					} else {
						tempModel.Rotation = rot
					}

					tempMaterial := Material{
						Diffuse:  parsedMaterial.Kd,
						Ambient:  parsedMaterial.Ka,
						Specular: parsedMaterial.Ks,
						Alpha:    parsedMaterial.D,
						N:        parsedMaterial.Ns,
					}

					//create temp material, checking for values
					if parsedMaterial.MapKD != "" && parsedMaterial.MapBump != "" {
						tempMaterial.DiffuseTexture = parsedMaterial.MapKD
						tempMaterial.NormalTexture = parsedMaterial.MapBump
						tempMaterial.ShaderType = 4
					} else if parsedMaterial.MapKD != "" && parsedMaterial.MapBump == "" {
						tempMaterial.DiffuseTexture = parsedMaterial.MapKD
						tempMaterial.ShaderType = 3
					} else {
						tempMaterial.ShaderType = 1
					}

					tempModelObject.Setup(
						tempMaterial,
						tempModel,
						tempName,
						scene[0].Objects[i].Collide,
						scene[0].Objects[i].Reflective,
						scene[0].Objects[i].RefractionIndex,
					)

					if scene[0].Objects[i].Parent != "" {
						tempModelObject.SetParent(scene[0].Objects[i].Parent)
					}

					state.Objects = append(state.Objects, &tempModelObject)
					state.LoadedObjects++
					fmt.Println(tempModelObject.name, " loaded successfully!")
				}
			}
		}
	}

	for j := 0; j < len(scene[0].PointLights); j++ {
		state.PointLights = append(state.PointLights, scene[0].PointLights[j])
	}

	for l := 0; l < len(scene[0].DirectionalLights); l++ {
		tempLight := DirectionalLight{
			Colour:    scene[0].DirectionalLights[l].Colour,
			Strength:  scene[0].DirectionalLights[l].Strength,
			Direction: scene[0].DirectionalLights[l].Direction,
			Position:  scene[0].DirectionalLights[l].Position,
		}
		state.DirectionalLights = append(state.DirectionalLights, tempLight)
	}
}

// GetSceneObject - Helper function for getting an object by searching using name
func GetSceneObject(name string, state State) Geometry {
	//TODO make this a goroutine instead
	for i := 0; i < len(state.Objects); i++ {
		objName, _, _ := state.Objects[i].GetDetails()
		if objName == name {
			return state.Objects[i]
		}
	}
	return nil
}

func getVertexRowN(vertices []float32, n int) mgl32.Vec3 {
	return mgl32.Vec3{vertices[n*3], vertices[(n*3)+1], vertices[(n*3)+2]}
}

func getUVRowN(uvs []float32, n int) mgl32.Vec2 {
	return mgl32.Vec2{uvs[n*2], uvs[(n*2)+1]}
}

// CalculateBitangents - Function for calculating bitangent of a given object
func CalculateBitangents(vertices []float32, uvs []float32) ([]float32, []float32) {
	//TODO confirm that this function is correct with Dana
	var tangents []float32
	var bitangents []float32

	for i := 0; i < (len(vertices)/3)-2; i += 3 {
		v0 := getVertexRowN(vertices, i)
		v1 := getVertexRowN(vertices, i+1)
		v2 := getVertexRowN(vertices, i+2)

		uv0 := getUVRowN(uvs, i)
		uv1 := getUVRowN(uvs, i+1)
		uv2 := getUVRowN(uvs, i+2)

		deltaPos1 := v1.Sub(v0)
		deltaPos2 := v2.Sub(v0)

		deltaUV1 := uv1.Sub(uv0)
		deltaUV2 := uv2.Sub(uv0)

		r := 1.0 / (deltaUV1[0]*deltaUV2[1] - deltaUV1[1]*deltaUV2[0])

		tempTangent1 := deltaPos1.Mul(deltaUV2[1])
		tempTangent2 := deltaPos2.Mul(deltaUV1[1])

		tangent := tempTangent1.Sub(tempTangent2)
		tangent = tangent.Mul(r)

		tempBitangent1 := deltaPos2.Mul(deltaUV1[0])
		tempBitangent2 := deltaPos1.Mul(deltaUV2[0])

		bitangent := tempBitangent1.Sub(tempBitangent2)
		bitangent = bitangent.Mul(r)

		for j := 0; j < 3; j++ {
			bitangents = append(bitangents, bitangent[0])
			bitangents = append(bitangents, bitangent[1])
			bitangents = append(bitangents, bitangent[2])

			tangents = append(tangents, tangent[0])
			tangents = append(tangents, tangent[1])
			tangents = append(tangents, tangent[2])
		}
	}
	return tangents, bitangents
}

// ToRadians - Simple helper function to convert degrees to radians
func ToRadians(deg float32) float64 {
	return float64(deg * (math.Pi / 180))
}

// CreateMat4FromArray - Helper function for taking a JSON array and converting it into a mat4
func CreateMat4FromArray(arr []float32) mgl32.Mat4 {

	if len(arr) != 16 {
		return mgl32.Ident4()
	}

	return mgl32.Mat4{arr[0], arr[1], arr[2], arr[3],
		arr[4], arr[5], arr[6], arr[7],
		arr[8], arr[9], arr[10], arr[11],
		arr[12], arr[13], arr[14], arr[15]}
}

// VectorDistance - helper function for getting distance between 2 vectors
func VectorDistance(a, b mgl32.Vec3) float32 {
	xDiff := math.Pow(float64(a[0]-b[0]), 2)
	yDiff := math.Pow(float64(a[1]-b[1]), 2)
	zDiff := math.Pow(float64(a[2]-b[2]), 2)

	return float32(math.Sqrt(xDiff + yDiff + zDiff))
}
