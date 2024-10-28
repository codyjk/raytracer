package color

import (
	"fmt"
	"io"
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

	// Translate [0,1] to [0,255]
	rbyte := int(255.999 * r)
	gbyte := int(255.999 * g)
	bbyte := int(255.999 * b)

	// Write the color components
	_, err := fmt.Fprintf(w, "%d %d %d\n", rbyte, gbyte, bbyte)
	return err
}
