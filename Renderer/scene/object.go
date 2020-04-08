package scene

import (
	"../geometry"
)

func GetObjectFromScene(state *geometry.State, name string) geometry.Geometry {
	for i := 0; i < len(state.Objects); i++ {
		oName, _, _ := state.Objects[i].GetDetails()
		if oName == name {
			return state.Objects[i]
		}
	}

	panic("Cannot find object " + name)
}

func GetLightFromScene(state *geometry.State, name string) *geometry.PointLight {
	for i := 0; i < len(state.PointLights); i++ {
		lName := state.PointLights[i].Name
		if lName == name {
			return &state.PointLights[i]
		}
	}

	panic("No object found of name: " + name)
}
