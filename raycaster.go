package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/oyberntzen/ecs"
	"github.com/oyberntzen/gogm"
	"github.com/oyberntzen/raycaster/renderer"

	_ "embed"
)

//go:embed vertex.glsl
var vertexShaderSource string

//go:embed fragment.glsl
var fragmentShaderSource string

//go:embed compute.glsl
var computeShaderSource string

var logger *log.Logger

const (
	width    = 700
	height   = 700
	rotSpeed = 0.01
)

var (
	cosPosRotSpeed float32 = float32(math.Cos(rotSpeed))
	sinPosRotSpeed float32 = float32(math.Sin(rotSpeed))
	cosNegRotSpeed float32 = float32(math.Cos(-rotSpeed))
	sinNegRotSpeed float32 = float32(math.Sin(-rotSpeed))
)

var worldMap []int32 = []int32{
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 7, 7, 7, 4,
	4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 4,
	4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4,
	4, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4,
	4, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 4,
	4, 0, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 7, 7, 0, 4,
	4, 0, 5, 0, 0, 0, 0, 5, 0, 5, 0, 5, 0, 5, 0, 5, 7, 0, 0, 4,
	4, 0, 6, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 4,
	4, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4,
	4, 0, 8, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 4,
	4, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 4,
	4, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 0, 5, 5, 5, 5, 7, 7, 7, 4,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 4,
	8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4,
	6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 4,
	4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 6, 0, 6, 2, 2, 2, 2, 2, 2, 4,
	4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 4,
	4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 2, 0, 0, 5, 0, 0, 4,
	4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 4,
	4, 4, 6, 4, 6, 4, 4, 4, 4, 4, 6, 4, 4, 4, 4, 4, 5, 4, 4, 4,
}

type raycaster struct {
	ecs.System
	window *renderer.Window

	va, ib, shader, ssbo, compute                             uint32
	uniformMap, uniformCamPos, uniformCamDir, uniformCamPlane int32

	vb *renderer.VertexBuffer

	camPos, camDir, camPlane gogm.Vec2[float32]
}

func (r *raycaster) Init() {
	logger = log.New(os.Stdout, "", log.Lshortfile)

	logger.Println("Init")

	// Initializing glfw and opengl
	var err error
	r.window, err = renderer.NewWindow(width, height, "Hello World")
	if err != nil {
		logger.Panic(err)
	}

	if err := gl.Init(); err != nil {
		logger.Panic(err)
	}

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(errorCallback, nil)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	logger.Println("OpenGL version", version)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	// Creating vertex buffer
	vbData := []float32{
		-1, -1,
		1, -1,
		-1, 1,
		1, 1,
	}
	r.vb = renderer.NewVertexBuffer(vbData)

	// Creating vertex array
	gl.GenVertexArrays(1, &r.va)
	gl.BindVertexArray(r.va)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 2*4, 0)

	// Creating index buffer
	gl.GenBuffers(1, &r.ib)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.ib)

	ibData := []uint32{
		0, 1, 2,
		1, 3, 2,
	}
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(ibData)*4, gl.Ptr(ibData), gl.STATIC_DRAW)

	// Creating vertex shader
	vs, err := compileShader(gl.VERTEX_SHADER, vertexShaderSource)
	if err != nil {
		logger.Fatal(err)
	}

	// Creating fragment shader
	fs, err := compileShader(gl.FRAGMENT_SHADER, fragmentShaderSource)
	if err != nil {
		logger.Fatal(err)
	}

	// Attach shaders
	r.shader = gl.CreateProgram()

	gl.AttachShader(r.shader, vs)
	gl.AttachShader(r.shader, fs)
	gl.LinkProgram(r.shader)
	gl.ValidateProgram(r.shader)

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	// Creating shader storage buffer object
	gl.GenBuffers(1, &r.ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, r.ssbo)

	rays := make([]int32, width)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(rays)*4, gl.Ptr(rays), gl.DYNAMIC_COPY)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, r.ssbo)

	// Creating compute shader
	cs, err := compileShader(gl.COMPUTE_SHADER, computeShaderSource)
	if err != nil {
		logger.Fatal(err)
	}

	r.compute = gl.CreateProgram()

	gl.AttachShader(r.compute, cs)
	gl.LinkProgram(r.compute)

	var result int32
	gl.GetProgramiv(r.compute, gl.LINK_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetProgramiv(r.compute, gl.INFO_LOG_LENGTH, &length)

		message := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(r.compute, length, &length, gl.Str(message))
		gl.DeleteProgram(r.compute)
		log.Fatalln(message)

	}

	gl.ValidateProgram(r.compute)

	gl.DeleteShader(cs)

	gl.UseProgram(r.compute)
	r.uniformMap = gl.GetUniformLocation(r.compute, gl.Str("map"+"\x00"))
	if r.uniformMap == -1 {
		logger.Fatalln("Uniform location not retrieved")
	}
	r.uniformCamPos = gl.GetUniformLocation(r.compute, gl.Str("camPos"+"\x00"))
	if r.uniformMap == -1 {
		logger.Fatalln("Uniform location not retrieved")
	}
	r.uniformCamDir = gl.GetUniformLocation(r.compute, gl.Str("camDir"+"\x00"))
	if r.uniformMap == -1 {
		logger.Fatalln("Uniform location not retrieved")
	}
	r.uniformCamPlane = gl.GetUniformLocation(r.compute, gl.Str("camPlane"+"\x00"))
	if r.uniformMap == -1 {
		logger.Fatalln("Uniform location not retrieved")
	}

	// Creating camera
	r.camPos = gogm.Vec2[float32]{11, 10}
	r.camDir = gogm.Vec2[float32]{-1, 1}
	r.camPlane = gogm.Vec2[float32]{0, 0.66}
}

func (r *raycaster) Update(dt float64) {
	if r.window.KeyPressed(glfw.KeyW) {
		var delta gogm.Vec2[float32]
		delta.MulS(&r.camDir, 0.01)
		r.camPos.Add(&r.camPos, &delta)
	}
	if r.window.KeyPressed(glfw.KeyS) {
		var delta gogm.Vec2[float32]
		delta.MulS(&r.camDir, -0.01)
		r.camPos.Add(&r.camPos, &delta)
	}
	if r.window.KeyPressed(glfw.KeyD) {
		oldDirX := r.camDir[0]
		r.camDir[0] = r.camDir[0]*cosNegRotSpeed - r.camDir[1]*sinNegRotSpeed
		r.camDir[1] = oldDirX*sinNegRotSpeed + r.camDir[1]*cosNegRotSpeed
		oldPlaneX := r.camPlane[0]
		r.camPlane[0] = r.camPlane[0]*cosNegRotSpeed - r.camPlane[1]*sinNegRotSpeed
		r.camPlane[1] = oldPlaneX*sinNegRotSpeed + r.camPlane[1]*cosNegRotSpeed
	}
	if r.window.KeyPressed(glfw.KeyA) {
		oldDirX := r.camDir[0]
		r.camDir[0] = r.camDir[0]*cosPosRotSpeed - r.camDir[1]*sinPosRotSpeed
		r.camDir[1] = oldDirX*sinPosRotSpeed + r.camDir[1]*cosPosRotSpeed
		oldPlaneX := r.camPlane[0]
		r.camPlane[0] = r.camPlane[0]*cosPosRotSpeed - r.camPlane[1]*sinPosRotSpeed
		r.camPlane[1] = oldPlaneX*sinPosRotSpeed + r.camPlane[1]*cosPosRotSpeed
	}

	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(r.compute)

	gl.Uniform1iv(r.uniformMap, int32(len(worldMap)), &worldMap[0])
	gl.Uniform2f(r.uniformCamPos, r.camPos[0], r.camPos[1])
	gl.Uniform2f(r.uniformCamDir, r.camDir[0], r.camDir[1])
	gl.Uniform2f(r.uniformCamPlane, r.camPlane[0], r.camPlane[1])

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, r.ssbo)
	gl.DispatchCompute(width, 1, 1)

	gl.BindVertexArray(r.va)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.ib)
	gl.UseProgram(r.shader)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

	r.window.Update()
}

func (r *raycaster) ShouldClose() bool {
	return r.window.ShouldClose()
}

func (r *raycaster) Delete() {
	r.vb.Delete()
	gl.DeleteVertexArrays(1, &r.va)
	gl.DeleteBuffers(1, &r.ib)
	gl.DeleteProgram(r.shader)

	gl.DeleteBuffers(1, &r.ssbo)
	gl.DeleteProgram(r.compute)

	r.window.Delete()
}

func init() {
	runtime.LockOSThread()
}

func errorCallback(source, gltype, id, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	if gltype == gl.DEBUG_TYPE_ERROR {
		logger.Printf("[OpengGL Error] type: %v, severity: %v, message: '%v'", gltype, severity, message)
	}
}

func compileShader(shaderType uint32, source string) (uint32, error) {
	id := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(id, 1, csources, nil)
	free()
	gl.CompileShader(id)

	var result int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)

		message := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(id, length, &length, gl.Str(message))
		gl.DeleteShader(id)
		if shaderType == gl.VERTEX_SHADER {
			return 0, fmt.Errorf("failed to compile vertex shader: %v", message)
		} else if shaderType == gl.FRAGMENT_SHADER {
			return 0, fmt.Errorf("failed to compile fragment shader: %v", message)
		} else {
			return 0, fmt.Errorf("failed to compile unknown shader: %v", message)
		}

	}

	return id, nil
}
