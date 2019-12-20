package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ParsedShader struct {
	VertShaderSource []string `json:"vertShader"`
	FragShaderSource []string `json:"fragShader"`
	AttribsSource    []string `json:"attribs"`
	UniformsSource   []string `json:"uniforms"`
	VertShaderText   string
	FragShaderText   string
	ShaderName       string
}

func ParseShaderFiles(shaderFiles []string) (error, []ParsedShader) {

	var shaders []ParsedShader

	for i := 0; i < len(shaderFiles); i++ {

		var tempShader ParsedShader

		jsonFile, err := os.Open(shaderFiles[i])
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

func GetShaderByName(name string, shaders []ParsedShader) ParsedShader {
	for i := 0; i < len(shaders); i++ {
		if name == shaders[i].ShaderName {
			return shaders[i]
		}
	}

	return ParsedShader{}
}
