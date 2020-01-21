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

	"../parser"

	"github.com/go-gl/mathgl/mgl32"
)

type RenderObject struct {
	ViewMatrix       mgl32.Mat4
	ProjMatrix       mgl32.Mat4
	ModelMatrix      mgl32.Mat4
	CurrentBuffers   ObjectBuffers
	CurrentModel     Model
	CurrentMaterial  Material
	CurrentCentroid  mgl32.Vec3
	CurrentProgram   ProgramInfo
	CameraPosition   []float32
	CurrentVertices  VertexValues
	CurrentObject    Geometry
	DistanceToCamera float32
}

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

type PointLight struct {
	Name      string    `json:"name"`
	Position  []float32 `json:"position"`
	Parent    string    `json:"parent"`
	Colour    []float32 `json:"colour"`
	Strength  float32   `json:"strength"`
	Quadratic float32   `json:"quadratic"`
	Linear    float32   `json:"linear"`
	Constant  float32   `json:"constant"`
}

type Settings struct {
}

type Scene struct {
	Objects  []SceneObject `json:"objects"`
	Lights   []PointLight  `json:"pointLights"`
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

	//create model
	tempModel := Model{
		Scale:    mgl32.Vec3{sceneObj.Scale[0], sceneObj.Scale[1], sceneObj.Scale[2]},
		Rotation: mgl32.Ident4(),
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
	)
	object.Translate(mgl32.Vec3{sceneObj.Position[0], sceneObj.Position[1], sceneObj.Position[2]})
	state.Objects = append(state.Objects, object)
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
		panic(err)
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

			objects := parser.Parse(meshPath)

			//tempModelObject := ModelObject{}

			for x := 0; x < len(objects); x++ {
				for j := 0; j < len(objects[x].Materials); j++ {

					fmt.Println("verts: ", len(objects[x].Geometry.Vertices), " normals: ", len(objects[x].Geometry.Normals), " uvs: ", len(objects[x].Geometry.UVs))
					fmt.Println("Start: ", objects[x].Materials[j].Start, " End: ", objects[x].Materials[j].End)

					parsedMaterial := parser.ParseMTLFile(objects[x].Materials[j].MTLLib, objects[x].Materials[j].Name)
					tempModelObject := ModelObject{}

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

					if x > 1 {
						tempModelObject.SetParent(scene[0].Objects[i].Name)
						tempName = strings.Join([]string{tempName, strconv.Itoa(j)}, "")
					} else {

					}

					tempModel := Model{
						Position: mgl32.Vec3{scene[0].Objects[i].Position[0], scene[0].Objects[i].Position[1], scene[0].Objects[i].Position[2]},
						Scale:    mgl32.Vec3{scene[0].Objects[i].Scale[0], scene[0].Objects[i].Scale[1], scene[0].Objects[i].Scale[2]},
						Rotation: mgl32.Ident4(),
					}

					if j > 1 {
						tempName = strings.Join([]string{tempName, strconv.Itoa(j)}, "")
						tempModelObject.SetParent(scene[0].Objects[i].Name)
						tempModel.Position = mgl32.Vec3{0, 0, 0}
						tempModel.Scale = mgl32.Vec3{1, 1, 1}
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
						tempName)

					state.Objects = append(state.Objects, &tempModelObject)
					state.LoadedObjects++
					fmt.Println(tempModelObject.name, " loaded successfully!")

				}
			}
			/*meshes, err := assimp.ParseFile(meshPath)

			if err != nil {
				panic(err)
			}

			for j := 0; j < len(meshes); j++ {
				// fmt.Println("Material: ", meshes[j].Material, "Num verts: ", len(meshes[j].Mesh.Vertices), "Num uvs: ", len(meshes[j].Mesh.UVChannels[0]))
				// fmt.Println(len(meshes[j].Mesh.Faces))
				//flatten the verts, normals and uvs
				verts := []float32{}
				norms := []float32{}
				uvs := []float32{}
				faces := []uint32{}

				tempMaterial := Material{
					DiffuseTexture: meshes[j].Material.DiffuseMap,
					Diffuse:        meshes[j].Material.Diffuse,
					Ambient:        meshes[j].Material.Ambient,
					Specular:       meshes[j].Material.Specular,
					Alpha:          meshes[j].Material.Alpha,
					N:              meshes[j].Material.Shininess,
					ShaderType:     scene[0].Objects[i].Material.ShaderType,
				}

				tempModelObject := ModelObject{}
				tempName := scene[0].Objects[i].Name

				for x := 0; x < len(meshes[j].Mesh.Vertices); x++ {
					verts = append(verts, meshes[j].Mesh.Vertices[x][0])
					verts = append(verts, meshes[j].Mesh.Vertices[x][1])
					verts = append(verts, meshes[j].Mesh.Vertices[x][2])

					norms = append(norms, meshes[j].Mesh.Normals[x][0])
					norms = append(norms, meshes[j].Mesh.Normals[x][1])
					norms = append(norms, meshes[j].Mesh.Normals[x][2])

				}
				for v := 0; v < len(meshes[j].Mesh.UVChannels[0]); v++ {
					if len(meshes[j].Mesh.UVChannels[0]) > 0 {
						uvs = append(uvs, meshes[j].Mesh.UVChannels[0][v][0])
						uvs = append(uvs, meshes[j].Mesh.UVChannels[0][v][1])
					}
				}

				for z := 0; z < len(meshes[j].Mesh.Faces); z++ {
					faces = append(faces, meshes[j].Mesh.Faces[z][0])
					faces = append(faces, meshes[j].Mesh.Faces[z][1])
					faces = append(faces, meshes[j].Mesh.Faces[z][2])
				}

				tempModel := Model{
					Scale:    mgl32.Vec3{scene[0].Objects[i].Scale[0], scene[0].Objects[i].Scale[1], scene[0].Objects[i].Scale[2]},
					Rotation: mgl32.Ident4(),
				}

				if j > 0 {
					tempModelObject.SetParent(scene[0].Objects[i].Name)
					tempName = strings.Join([]string{tempName, strconv.Itoa(j)}, "")
					tempModelObject.SetParent(scene[0].Objects[i].Name)
					tempModel.Position = mgl32.Vec3{0, 0, 0}
					tempModel.Scale = mgl32.Vec3{1, 1, 1}
				}

				tempModelObject.SetVertexValues(verts, norms, uvs, faces)

				tempModelObject.Setup(
					tempMaterial,
					tempModel,
					tempName)

				if j == 0 {
					tempModelObject.Translate(mgl32.Vec3{scene[0].Objects[i].Position[0], scene[0].Objects[i].Position[1], scene[0].Objects[i].Position[2]})
				} else {
					tempModelObject.boundingBox = TranslateBoundingBox(tempModelObject.boundingBox, mgl32.Vec3{scene[0].Objects[i].Position[0], scene[0].Objects[i].Position[1], scene[0].Objects[i].Position[2]})
				}

				state.Objects = append(state.Objects, &tempModelObject)
				state.LoadedObjects++
				fmt.Println(tempModelObject.name, " loaded successfully!")

			} */
		}
	}

	for j := 0; j < len(scene[0].Lights); j++ {
		tempLight := Light{
			Colour:    scene[0].Lights[j].Colour,
			Strength:  scene[0].Lights[j].Strength,
			Position:  scene[0].Lights[j].Position,
			Quadratic: scene[0].Lights[j].Quadratic,
			Linear:    scene[0].Lights[j].Linear,
			Constant:  scene[0].Lights[j].Constant,
		}
		state.Lights = append(state.Lights, tempLight)
	}

}

func GetSceneObject(name string, state State) Geometry {
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

func CalculateBitangents(vertices []float32, uvs []float32) ([]float32, []float32) {
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

func ToRadians(deg float32) float64 {
	return float64(deg * (math.Pi / 180))
}

func VectorDistance(a, b mgl32.Vec3) float32 {
	//fmt.Println(a)
	xDiff := math.Pow(float64(a[0]-b[0]), 2)
	yDiff := math.Pow(float64(a[1]-b[1]), 2)
	zDiff := math.Pow(float64(a[2]-b[2]), 2)

	return float32(math.Sqrt(xDiff + yDiff + zDiff))
}
