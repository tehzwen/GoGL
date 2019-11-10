package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 500
	height = 500

	vertexShaderSource = `
		#version 410
		in vec3 aPosition;
		
		uniform mat4 uProjectionMatrix;
		uniform mat4 uViewMatrix;
		uniform mat4 uModelMatrix;

		void main() {
			gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix * vec4(aPosition, 1.0);
			//gl_Position = vec4(aPosition, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		uniform vec3 diffuseVal;
		uniform vec3 ambientVal;
		uniform vec3 specularVal;

		out vec4 frag_colour;

		void main() {
			//frag_colour = vec4(diffuseVal * ambientVal * specularVal, 1.0);
			frag_colour = vec4(diffuseVal, 1.0);
		}
	` + "\x00"
)

var (
	triangle = []float32{
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.0, 0.0,

		0.0, 0.0, 0.5,
		0.0, 0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.0, 0.5,

		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.5,

		0.0, 0.0, 0.5,
		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.0, 0.0, 0.0,

		0.5, 0.0, 0.5,
		0.5, 0.0, 0.0,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.0,

		0.0, 0.0, 0.5,
		0.0, 0.0, 0.0,
		0.0, 0.5, 0.5,
		0.0, 0.5, 0.0,
	}
)

var triangleFaces = []uint32{
	0, 1, 2, 0, 2, 3,
	4, 5, 6, 4, 6, 7,
	8, 9, 10, 8, 10, 11,
	12, 13, 14, 12, 14, 15,
	16, 17, 18, 17, 18, 19,
	20, 21, 22, 21, 22, 23,
}

func main() {
	runtime.LockOSThread()

	var state = State{
		vertShader: vertexShaderSource,
		fragShader: fragmentShaderSource,
		camera: Camera{
			position: mgl32.Vec3{1.5, 0.5, -2.5},
			center:   mgl32.Vec3{0.5, 0.0, 0.0},
			up:       mgl32.Vec3{0.0, 1.0, 0.0},
		},
	}

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL(state.vertShader, state.fragShader)

	var testObject = Object{
		name:       "zach",
		fragShader: "1234",
		vertShader: "123456",
		programInfo: ProgramInfo{
			program: program,
			attributes: Attributes{
				position: 0,
			},
			uniformLocations: Uniforms{
				diffuseVal: -2,
			},
		},
		vertices: triangle,
		material: Material{
			diffuse:  []float32{0.6, 0.2, 0.6},
			ambient:  []float32{0.1, 0.1, 0.1},
			specular: []float32{0.8, 0.8, 0.8},
		},
		model: Model{
			position: mgl32.Vec3{1.0, 0.0, 0.0},
		},
	}

	SetupAttributes(&testObject.programInfo)

	var objectList = []Object{}

	//vao, vbo := initPositionBuffer(&testObject.programInfo, testObject.vertices)
	//indexBuffer := InitIndexBuffer(triangleFaces)

	//testObject.programInfo.indexBuffer = indexBuffer
	//testObject.buffers.vao = vao
	//testObject.buffers.vbo = vbo
	vao := CreateTriangleVAO(triangle, triangleFaces)

	testObject.buffers.vao = vao

	testObject.programInfo.uniformLocations.diffuseVal = gl.GetUniformLocation(testObject.programInfo.program, gl.Str("diffuseVal\x00"))

	objectList = append(objectList, testObject)

	state.objects = objectList

	for !window.ShouldClose() {
		//projection := glm.Mat4{}
		//fmt.Printf("%+v\n", projection)
		draw(window, state)
		//state.camera.center[0] -= 0.2
	}
}

func draw(window *glfw.Window, state State) {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	glfw.PollEvents()

	for i := 0; i < len(state.objects); i++ {

		gl.UseProgram(state.objects[i].programInfo.program)

		var fovy = float32(60 * math.Pi / 180)
		var aspect = float32(width / height)
		var near = float32(0.1)
		var far = float32(100.0)

		projection := mgl32.Perspective(fovy, aspect, near, far)

		gl.UniformMatrix4fv(state.objects[i].programInfo.uniformLocations.projection, 1, false, &projection[0])

		viewMatrix := mgl32.LookAtV(state.camera.position, state.camera.center, state.camera.up)

		gl.UniformMatrix4fv(state.objects[i].programInfo.uniformLocations.view, 1, false, &viewMatrix[0])

		modelMatrix := mgl32.Ident4()

		modelMatrix = mgl32.HomogRotate3D(float32(0), state.objects[i].model.position)

		gl.UniformMatrix4fv(state.objects[i].programInfo.uniformLocations.model, 1, false, &modelMatrix[0])

		gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.diffuseVal, 1, &state.objects[i].material.diffuse[0])
		gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.ambientVal, 1, &state.objects[i].material.ambient[0])
		gl.Uniform3fv(state.objects[i].programInfo.uniformLocations.specularVal, 1, &state.objects[i].material.specular[0])

		gl.BindVertexArray(state.objects[i].buffers.vao)
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
		gl.DrawElements(gl.TRIANGLES, int32(len(triangleFaces)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
		gl.BindVertexArray(0)

		window.SwapBuffers()
	}

}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Go GL", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL(vertexShaderSource string, fragmentShaderSource string) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func initPositionBuffer(programInfo *ProgramInfo, positionArray []float32) (uint32, uint32) {
	var vbo uint32
	var vao uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(positionArray), gl.Ptr(positionArray), gl.STATIC_DRAW)

	const numComponents = 3
	const dataType = gl.FLOAT
	const normalize = false
	const stride = 0

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(programInfo.attributes.position, numComponents, dataType, normalize, stride, nil)
	gl.EnableVertexAttribArray(programInfo.attributes.position)

	return vao, vbo
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
