package main

import (
	"fmt"
	"os"
	"raytracer/internal/camera"
	"raytracer/internal/color"
	"raytracer/internal/core"
	"raytracer/internal/hittable"
	"raytracer/internal/material"
	"raytracer/internal/util"
	"raytracer/internal/vector"
)

func main() {
	world := hittable.NewHittableList()

	groundMaterial := material.NewLambertian(color.NewColor(0.5, 0.5, 0.5))
	world.Add(hittable.NewSphere(vector.NewPoint3(0, -1000, 0), 1000, groundMaterial))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := util.RandomFloat()
			center := vector.NewPoint3(float64(a)+0.9*util.RandomFloat(), 0.2, float64(b)+0.9*util.RandomFloat())

			var sphereMaterial core.Material

			if chooseMat < 0.8 {
				// Diffuse
				albedo := color.Random().Mul(color.Random())
				sphereMaterial = material.NewLambertian(albedo)
			} else if chooseMat < 0.95 {
				// Metal
				albedo := color.RandomFromRange(0.5, 1)
				fuzz := util.RandomFloatFromRange(0, 0.5)
				sphereMaterial = material.NewMetal(albedo, fuzz)
			} else {
				// Glass
				sphereMaterial = material.NewDielectric(1.5)
			}

			world.Add(hittable.NewSphere(center, 0.2, sphereMaterial))
		}
	}

	mat1 := material.NewDielectric(1.5)
	world.Add(hittable.NewSphere(vector.NewPoint3(0, 1, 0), 1.0, mat1))

	mat2 := material.NewLambertian(color.NewColor(0.4, 0.2, 0.1))
	world.Add(hittable.NewSphere(vector.NewPoint3(-4, 1, 0), 1.0, mat2))

	mat3 := material.NewMetal(color.NewColor(0.7, 0.6, 0.5), 0.0)
	world.Add(hittable.NewSphere(vector.NewPoint3(4, 1, 0), 1.0, mat3))

	camConfig := camera.DefaultConfig()

	camConfig.AspectRatio = 16.0 / 9.0
	camConfig.ImageWidth = 1200
	camConfig.SamplesPerPixel = 50
	camConfig.MaxDepth = 50

	camConfig.VerticalFOV = 20
	camConfig.LookFrom = vector.NewPoint3(13, 2, 3)
	camConfig.LookAt = vector.NewPoint3(0, 0, 0)
	camConfig.VUp = vector.NewVec3(0, 1, 0)

	camConfig.DefocusAngle = 0.6
	camConfig.FocusDist = 10.0

	cam, err := camera.New(camConfig)
	if err != nil {
		fmt.Print(fmt.Errorf("failed to create camera: %w", err))
		return
	}

	cam.Render(os.Stdout, os.Stderr, world)
}
