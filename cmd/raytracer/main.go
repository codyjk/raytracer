package main

import (
	"os"
	"raytracer/internal/camera"
	"raytracer/internal/hittable"
	"raytracer/internal/vector"
)

func main() {
	world := hittable.NewHittableList()
	world.Add(hittable.NewSphere(vector.NewPoint3(0, 0, -1), 0.5))
	world.Add(hittable.NewSphere(vector.NewPoint3(0, -100.5, -1), 100))

	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	cam := camera.NewCamera(aspectRatio, imageWidth)

	cam.Render(os.Stdout, os.Stderr, world)
}
