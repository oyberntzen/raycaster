package renderer

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

}

func errorCallback(source, gltype, id, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	if gltype == gl.DEBUG_TYPE_ERROR {
		log.Printf("[OpengGL Error] type: %v, severity: %v, message: '%v'", gltype, severity, message)
	}
}

type Window struct {
	window *glfw.Window
}

func NewWindow(width, height int, title string) (*Window, error) {
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		log.Panic(err)
	}

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(errorCallback, nil)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	return &Window{
		window: window,
	}, nil
}

func (window *Window) Update() {
	window.window.SwapBuffers()
	glfw.PollEvents()
}

func (window *Window) KeyPressed(key glfw.Key) bool {
	return window.window.GetKey(key) == glfw.Press
}

func (window *Window) ShouldClose() bool {
	return window.window.ShouldClose()
}

func (window *Window) Delete() {
	glfw.Terminate()
}

func (window *Window) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (window *Window) Draw(vb *VertexBuffer, ib *IndexBuffer, shader *Shader) {
	vb.Bind()
	ib.Bind()
	shader.Bind()
	gl.DrawElements(gl.TRIANGLES, int32(ib.count), gl.UNSIGNED_INT, nil)
}
