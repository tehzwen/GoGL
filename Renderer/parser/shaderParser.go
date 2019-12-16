package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Shader struct {
	VertShaderSource []string `json:"vertShader"`
	FragShaderSource []string `json:"fragShader"`
	AttribsSource    []string `json:"attribs"`
	UniformsSource   []string `json:"uniforms"`
	VertShaderText   string
	FragShaderText   string
	ShaderName       string
}

func ParseShaderFiles(shaderFiles []string) (error, []Shader) {

	var shaders []Shader

	for i := 0; i < len(shaderFiles); i++ {

		var tempShader Shader

		jsonFile, err := os.Open(shaderFiles[i])

		fmt.Printf("%v+\n", jsonFile)
		if err != nil {
			return errors.New("Error reading the shader file"), nil
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, &tempShader)

		if err != nil {
			return errors.New("Error reading the shader file"), nil
		}

		tempShader.ShaderName = filepath.Base(shaderFiles[i])
		//join the array of strings into a single string for the shader
		tempShader.VertShaderText = strings.Join(tempShader.VertShaderSource[:], "\n")
		tempShader.FragShaderText = strings.Join(tempShader.FragShaderSource[:], "\n")

		shaders = append(shaders, tempShader)

	}
	return nil, shaders

}
