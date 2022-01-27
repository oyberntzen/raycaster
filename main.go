package main

import "github.com/oyberntzen/ecs"

func main() {
	scene := ecs.Scene{}
	ren := &raycaster{}
	scene.AddSystem(ren)

	scene.Init()

	for !ren.window.ShouldClose() {
		scene.Update(0)
	}
}
