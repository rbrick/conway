package main

import (
	"fmt"
	"log"
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

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
