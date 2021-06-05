package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"gonum.org/v1/gonum/mat"
)

const (
	Width, Height = 800, 800
)

func init() {
	// must be locked on main thread sadface
	runtime.LockOSThread()
}

func main() {

	identity := mat.NewDense(4, 4, IdentityMatrix4x4())

	array := make([]float32, 16)

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
	glfw.WindowHint(glfw.Resizable, glfw.False)

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
	fsrc, _ := os.Open("simple.glsl")

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

	buf := CreateBuffer(BufferType{
		Components: []VertexComponent{PositionComponent},
	})

	angle := float32(20.)
	angleStep := float32(360. / angle)

	for i := float32(0.); i <= angleStep; i += .1 {
		rads := i * angle / 180.0 * float32(math.Pi)
		x := float32(math.Cos(float64(rads))) * .85
		y := float32(math.Sin(float64(rads))) * .85

		buf.Vertex(NewVertex(float32(x), float32(y), 0))
	}

	buf.Vertex(NewVertex(0, 0, 0))

	buf.Bind(gl.ARRAY_BUFFER)
	buf.Upload(gl.ARRAY_BUFFER)
	buf.Unbind(gl.ARRAY_BUFFER)

	startTime := time.Now()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.ClearColor(1., 1, 1, 1.)

		shaderProgram.Bind()

		timeUniform := shaderProgram.GetUniform("iTime")

		timeUniform.Float(float32(time.Since(startTime).Seconds()))

		resolutionUniform := shaderProgram.GetUniform("iResolution")

		width, height := glfw.GetCurrentContext().GetFramebufferSize()

		resolutionUniform.Vecf(float32(width), float32(height), 3)

		viewMatrix := shaderProgram.GetUniform("viewMatrix")

		viewMatrix.Matrix(mat.NewDense(4, 4, IdentityMatrix4x4()), false)

		buf.Draw(gl.TRIANGLE_FAN)
		buf.Unbind(gl.ARRAY_BUFFER)

		shaderProgram.Unbind()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func equalArray(a, b []float32) (int, bool) {
	if len(b) != len(a) {
		return -1, false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return i, false
		}
	}

	return len(a), true
}
