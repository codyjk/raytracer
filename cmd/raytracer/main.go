package main

import (
	"math"
	"os"
	"raytracer/internal/camera"
	"raytracer/internal/color"
	"raytracer/internal/hittable"
	"raytracer/internal/material"
	"raytracer/internal/vector"
)

func main() {
	world := hittable.NewHittableList()

	R := math.Cos(math.Pi / 4)
	materialLeft := material.NewLambertian(color.NewColor(0, 0, 1))
	materialRight := material.NewLambertian(color.NewColor(1, 0, 0))
	world.Add(hittable.NewSphere(vector.NewPoint3(-R, 0, -1), R, materialLeft))
	world.Add(hittable.NewSphere(vector.NewPoint3(R, 0, -1), R, materialRight))

	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	samplesPerPixel := 100
	maxDepth := 50
	vFov := 90.0
	cam := camera.NewCamera(aspectRatio, imageWidth, samplesPerPixel, maxDepth, vFov)

	cam.Render(os.Stdout, os.Stderr, world)
}
