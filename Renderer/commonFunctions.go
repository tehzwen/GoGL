package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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

func ParseJsonFile(filePath string, state *State) {
	fmt.Printf("here: %s\n", filePath)

	jsonFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var scene []Scene

	err = json.Unmarshal(byteValue, &scene)

	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(scene[0].Objects); i++ {
		if scene[0].Objects[i].ObjectType == "cube" {
			tempCube := Cube{}

			err := tempCube.SetShader(state.vertShader, state.fragShader)

			if err != nil {
				fmt.Println(err)
				panic(err)
			} else {
				//create model
				tempModel := Model{
					position: mgl32.Vec3{scene[0].Objects[i].Position[0], scene[0].Objects[i].Position[1], scene[0].Objects[i].Position[2]},
					scale:    mgl32.Vec3{scene[0].Objects[i].Scale[0], scene[0].Objects[i].Scale[1], scene[0].Objects[i].Scale[2]},
					rotation: mgl32.Ident4(),
				}
				tempCube.Setup(
					scene[0].Objects[i].Material,
					tempModel,
					scene[0].Objects[i].Name,
				)
				//tempArray = append(tempArray)
				fmt.Printf("%v+\n", tempModel)
				state.objects = append(state.objects, &tempCube)
			}

		}
	}

}
