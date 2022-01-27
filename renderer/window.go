package renderer

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	if err := glfw.Init(); err != nil {
		log.Panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
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
