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
	materialLeft := material.NewDielectric(1.50)
	// The air bubble's refractive index is relative to the containing sphere (glass)
	materialBubble := material.NewDielectric(1.00 / 1.50)
	materialRight := material.NewMetal(color.NewColor(0.8, 0.6, 0.2), 1.0)

	world.Add(hittable.NewSphere(vector.NewPoint3(0, -100.5, -1), 100, materialGround))
	world.Add(hittable.NewSphere(vector.NewPoint3(0, 0, -1.2), 0.5, materialCenter))
	world.Add(hittable.NewSphere(vector.NewPoint3(-1.0, 0, -1.0), 0.5, materialLeft))
	world.Add(hittable.NewSphere(vector.NewPoint3(-1.0, 0, -1.0), 0.4, materialBubble))
	world.Add(hittable.NewSphere(vector.NewPoint3(1.0, 0, -1.0), 0.5, materialRight))

	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	samplesPerPixel := 100
	maxDepth := 50
	vFov := 20.0
	lookFrom := vector.NewPoint3(-2, 2, 1)
	lookAt := vector.NewPoint3(0, 0, -1)
	vUp := vector.NewVec3(0, 1, 0)
	cam := camera.NewCamera(aspectRatio, imageWidth, samplesPerPixel, maxDepth, vFov, lookFrom, lookAt, vUp)

	cam.Render(os.Stdout, os.Stderr, world)
}
