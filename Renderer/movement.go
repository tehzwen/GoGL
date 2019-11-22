package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

//Callbacks for inputs

func KeyHandler(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		keys[key] = true
	} else if action == glfw.Release {
		keys[key] = false
	}
}

func MouseButtonHandler(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		buttons[button] = true
	} else if action == glfw.Release {
		buttons[button] = false
	}
}

func MouseMoveHandler(win *glfw.Window, xPos float64, yPos float64) {
	xDiff := xPos - mouseMovement["X"]
	yDiff := yPos - mouseMovement["Y"]

	mouseMovement["Xmove"] = xDiff
	mouseMovement["Ymove"] = yDiff
	mouseMovement["X"] = xPos
	mouseMovement["Y"] = yPos
	mouseMovement["move"] = 1
}
