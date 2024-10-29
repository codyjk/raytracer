package color

import (
	"fmt"
	"io"
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

	// Translate [0,1] to [0,255]
	intensity := interval.NewInterval(0.000, 0.999)
	rbyte := int(256 * intensity.Clamp(r))
	gbyte := int(256 * intensity.Clamp(g))
	bbyte := int(256 * intensity.Clamp(b))

	// Write the color components
	_, err := fmt.Fprintf(w, "%d %d %d\n", rbyte, gbyte, bbyte)
	return err
}
