# Golang Rendering Engine
## Requirements

### GLFW 3.1
The following instructions are taken from the GLFW github page used in this project (I know I am using GLFW 3.1, and the readme below says 3.3, either should work but for now I am sticking with 3.1).


### GLFW 3.3 for Go [![Build Status](https://travis-ci.org/go-gl/glfw.svg?branch=master)](https://travis-ci.org/go-gl/glfw) [![GoDoc](https://godoc.org/github.com/go-gl/glfw/v3.3/glfw?status.svg)](https://godoc.org/github.com/go-gl/glfw/v3.3/glfw)

#### Installation

* GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
	* On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
	* On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
	* On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.
	* See [here](http://www.glfw.org/docs/latest/compile.html#compile_deps) for full details.
* Go 1.4+ is required on Windows (otherwise you must use MinGW v4.8.1 exactly, see [Go issue 8811](https://github.com/golang/go/issues/8811)).

```
go get -u github.com/go-gl/glfw/v3.3/glfw
```

### MGL32

This is a math library used by the engine for matrices, vectors and so on.

```
go get -u github.com/go-gl/mathgl/mgl32
```

### OpenGL for Go

This repository holds Go bindings to various OpenGL versions. 

```
go get -u github.com/go-gl/gl/v{3.2,3.3,4.1,4.2,4.3,4.4,4.5,4.6}-{core,compatibility}/gl
```
for this project I am using ```github.com/go-gl/gl/v4.1-core/gl```

### GOMBZ

This is a Go library that provides a serializable data structure for 3d models and animations.

```
go get -u github.com/tbogdala/gombz
```

### Assimp

This library provides the reading of 3D object files to be used in the engine. For ubuntu, I simply use ```sudo apt install assimp-utils```


