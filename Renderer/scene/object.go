package scene

import (
	"errors"

	"../geometry"
)

func GetObjectFromScene(state *geometry.State, name string) (geometry.Geometry, error) {
	for i := 0; i < len(state.Objects); i++ {
		oName, _, _ := state.Objects[i].GetDetails()
		if oName == name {
			return state.Objects[i], nil
		}
	}

	return nil, errors.New("No object found of name: " + name)
}

func GetLightFromScene(state *geometry.State, name string) (*geometry.PointLight, error) {
	for i := 0; i < len(state.PointLights); i++ {
		lName := state.PointLights[i].Name
		if lName == name {
			return &state.PointLights[i], nil
		}
	}

	return &geometry.PointLight{}, errors.New("No object found of name: " + name)
}
