package main

import (
	"fmt"
	"math"
	"os"
	"raytracer/internal/color"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

// This math solves for the ray-sphere intersection. The math is explained here:
// https://raytracing.github.io/books/RayTracingInOneWeekend.html#addingasphere/ray-sphereintersection
func hitSphere(center vector.Point3, radius float64, r ray.Ray) float64 {
	oc := center.Sub(r.Origin())
	a := vector.Dot(r.Direction(), r.Direction())
	b := -2.0 * vector.Dot(r.Direction(), oc)
	c := vector.Dot(oc, oc) - radius*radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return -1.0
	} else {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a)
	}
}

// We are rendering a horizontal greadient.
// blendedValue = (1 - a) * startValue + a * endValue, where a is the linear
// scale of the ray direction.
func rayColor(r ray.Ray) color.Color {
	sphereCenter := vector.NewPoint3(0, 0, -1)
	t := hitSphere(sphereCenter, 0.5, r)
	if t > 0.0 {
		N := r.At(t).Sub(sphereCenter).Unit()
		return color.NewColor(N.X()+1, N.Y()+1, N.Z()+1).Scale(0.5)
	}

	unitDirection := r.Direction().Unit()
	a := 0.5 * (unitDirection.Y() + 1.0)
	return color.NewColor(1.0, 1.0, 1.0).Scale(1.0 - a).Add(color.NewColor(0.5, 0.7, 1.0).Scale(a))
}

func main() {
	// Image
	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)
	if imageHeight < 1 {
		imageHeight = 1
	}

	// Camera
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	cameraCenter := vector.NewPoint3(0, 0, 0)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	// In our image, (0,0) represents the top left pixel. In the coordinate space
	// of this program, it represents the bottom left. So, the directions in which
	// the code iterates through the pictures needs to be inverted to match our
	// mental model of the image starting in the top left.
	viewportU := vector.NewVec3(viewportWidth, 0, 0)
	viewportV := vector.NewVec3(0, -viewportHeight, 0)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel
	pixelDeltaU := viewportU.Div(float64(imageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))

	// Calculate the locatino of the upper left pixel.
	viewportUpperLeft := cameraCenter.Sub(vector.NewVec3(0, 0, focalLength)).Sub(viewportU.Div(2.0)).Sub(viewportV.Div(2.0))
	pixel00Loc := viewportUpperLeft.Add(pixelDeltaU.Add(pixelDeltaV).Scale(0.5))

	// Render
	fmt.Printf("P3\n%d %d\n255\n", imageWidth, imageHeight)

	for j := 0; j < imageHeight; j++ {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d", imageHeight-j)
		for i := 0; i < imageWidth; i++ {
			pixelCenter := pixel00Loc.Add(pixelDeltaU.Scale(float64(i))).Add(pixelDeltaV.Scale(float64(j)))
			rayDirection := pixelCenter.Sub(cameraCenter)
			ray := ray.NewRay(cameraCenter, rayDirection)

			pixelColor := rayColor(ray)
			color.WriteColor(os.Stdout, pixelColor)
		}
	}

	fmt.Fprintln(os.Stderr, "\nDone.")
}
