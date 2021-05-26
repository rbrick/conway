package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"gonum.org/v1/gonum/mat"
)

const (
	Width, Height = 800, 600
)

func init() {
	// must be locked on main thread sadface
	runtime.LockOSThread()
}

func main() {

	identity := mat.NewDense(4, 4, IdentityMatrix4x4())

	array := make([]float64, 16)

	MatrixToArray(identity, array)

	fmt.Println(array)

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize:", err) // failed to initialize GLFW
	}

	// terminate when we exit the game loop
	defer glfw.Terminate()

	// set the version to OpenGL 3.2
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)

	// set it to the core profile
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// attempt to create a window
	window, err := glfw.CreateWindow(Width, Height, "3D Conway", nil, nil)

	if err != nil {
		log.Fatalln("failed to create window:", err)
	}

	window.MakeContextCurrent()

	// initialize OpenGL
	if err = gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vsrc, _ := os.Open("vertexShader.glsl")
	fsrc, _ := os.Open("fragmentShader.glsl")

	vertexShader, err := ReadShader(vsrc, gl.VERTEX_SHADER)

	if err != nil {
		panic(err)
	}

	fragmentShader, err := ReadShader(fsrc, gl.FRAGMENT_SHADER)

	if err != nil {
		panic(err)
	}

	shaderProgram := NewProgram()

	shaderProgram.Attach(vertexShader)
	shaderProgram.Attach(fragmentShader)
	shaderProgram.Link()

	buffer := CreateBuffer()

	buffer.Vertex(NewVertex(-1, -1, 0))
	buffer.Vertex(NewVertex(1, -1, 0))
	buffer.Vertex(NewVertex(1, 1, 0))
	buffer.Vertex(NewVertex(1, -1, 0))

	// create a vao

	var vao uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao) // bind it

	buffer.Bind(gl.ARRAY_BUFFER) // bind as an array for this buffer
	buffer.Upload(gl.ARRAY_BUFFER)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 1, gl.PtrOffset(0))
	// gl.EnableVertexArrayAttrib(vao, 0)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shaderProgram.Bind()

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.QUADS, 0, 4)
		gl.BindVertexArray(0)

		shaderProgram.Unbind()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
