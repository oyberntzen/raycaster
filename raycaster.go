package main

import (
	"log"
	"math"
	"os"

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

	vb      *renderer.VertexBuffer
	ib      *renderer.IndexBuffer
	shader  *renderer.Shader
	compute *renderer.ComputeShader
	ssbo    *renderer.ShaderStorageBufferObject

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

	// Creating vertex buffer
	vbData := []float32{
		-1, -1,
		1, -1,
		-1, 1,
		1, 1,
	}
	r.vb = renderer.NewVertexBuffer(vbData, []renderer.LayoutElement{
		{ShaderDataType: renderer.ShaderDataTypeFloat32, Count: 2},
	})

	// Creating index buffer

	ibData := []uint32{
		0, 1, 2,
		1, 3, 2,
	}
	r.ib = renderer.NewIndexBuffer(ibData)

	// Attach shaders
	r.shader = renderer.NewShader(vertexShaderSource, fragmentShaderSource)

	// Creating shader storage buffer object
	rays := make([]int32, width)
	r.ssbo = renderer.NewShaderStorageBufferObject(rays, 0)

	// Creating compute shader
	r.compute = renderer.NewComputeShader(computeShaderSource)

	// Creating camera
	r.camPos = gogm.Vec2[float32]{11, 10}
	r.camDir = gogm.Vec2[float32]{-1, 0}
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

	r.window.Clear()

	r.compute.Bind()
	r.compute.UploadUniformArrayInt32("map", worldMap)
	r.compute.UploadUniformVec2Float32("camPos", &r.camPos)
	r.compute.UploadUniformVec2Float32("camDir", &r.camDir)
	r.compute.UploadUniformVec2Float32("camPlane", &r.camPlane)

	r.ssbo.Bind()
	r.compute.Run(width, 1, 1)

	r.window.Draw(r.vb, r.ib, r.shader)

	r.window.Update()
}

func (r *raycaster) ShouldClose() bool {
	return r.window.ShouldClose()
}

func (r *raycaster) Delete() {
	r.vb.Delete()
	r.ib.Delete()
	r.shader.Delete()

	r.ssbo.Delete()
	r.compute.Delete()

	r.window.Delete()
}
