package color

import (
	"fmt"
	"io"
	"math"
	"raytracer/internal/interval"
	"raytracer/internal/vector"
)

type Color = vector.Vec3

func NewColor(r, g, b float64) Color {
	return vector.NewVec3(r, g, b)
}

func WriteColor(w io.Writer, pixelColor Color) error {
	r := pixelColor.X()
	g := pixelColor.Y()
	b := pixelColor.Z()

	// Apply a linear to gamma transform for gamma 2
	r = linearToGamma(r)
	g = linearToGamma(g)
	b = linearToGamma(b)

	// Translate [0,1] to [0,255]
	intensity := interval.NewInterval(0.000, 0.999)
	rbyte := int(256 * intensity.Clamp(r))
	gbyte := int(256 * intensity.Clamp(g))
	bbyte := int(256 * intensity.Clamp(b))

	// Write the color components
	_, err := fmt.Fprintf(w, "%d %d %d\n", rbyte, gbyte, bbyte)
	return err
}

func Random() Color {
	return vector.Random()
}

func RandomFromRange(min, max float64) Color {
	return vector.RandomFromRange(min, max)
}

func linearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		// We're applying "gamma 2" - 2 being the power when going from gamma to linear space.
		// From linear to gamma, therefore, is power 1/2.
		return math.Sqrt(linearComponent)
	}

	return 0
}
