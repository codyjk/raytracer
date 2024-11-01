package main

import (
	"os"
	"raytracer/internal/camera"
	"raytracer/internal/color"
	"raytracer/internal/hittable"
	"raytracer/internal/material"
	"raytracer/internal/vector"
)

func main() {
	world := hittable.NewHittableList()

	materialGround := material.NewLambertian(color.NewColor(0.8, 0.8, 0.0))
	materialCenter := material.NewLambertian(color.NewColor(0.1, 0.2, 0.5))
	materialLeft := material.NewMetal(color.NewColor(0.8, 0.8, 0.8))
	materialRight := material.NewMetal(color.NewColor(0.8, 0.6, 0.2))

	world.Add(hittable.NewSphere(vector.NewPoint3(0, -100.5, -1), 100, materialGround))
	world.Add(hittable.NewSphere(vector.NewPoint3(0, 0, -1.2), 0.5, materialCenter))
	world.Add(hittable.NewSphere(vector.NewPoint3(-1.0, 0, -1.0), 0.5, materialLeft))
	world.Add(hittable.NewSphere(vector.NewPoint3(1.0, 0, -1.0), 0.5, materialRight))

	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	samplesPerPixel := 100
	maxDepth := 50
	cam := camera.NewCamera(aspectRatio, imageWidth, samplesPerPixel, maxDepth)

	cam.Render(os.Stdout, os.Stderr, world)
}
