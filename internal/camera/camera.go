package camera

import (
	"fmt"
	"io"
	"math"
	"raytracer/internal/color"
	"raytracer/internal/hittable"
	"raytracer/internal/interval"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type Camera struct {
	aspectRatio float64
	imageWidth  int
	imageHeight int
	center      vector.Point3
	pixel00Loc  vector.Point3
	pixelDeltaU vector.Vec3
	pixelDeltaV vector.Vec3
}

func NewCamera(aspectRatio float64, imageWidth int) Camera {
	imageHeight := int(float64(imageWidth) / aspectRatio)
	if imageHeight < 1 {
		imageHeight = 1
	}

	center := vector.NewPoint3(0, 0, 0)

	// Determine viewport dimensions.
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))

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
	viewportUpperLeft := center.Sub(vector.NewVec3(0, 0, focalLength)).Sub(viewportU.Div(2.0)).Sub(viewportV.Div(2.0))
	pixel00Loc := viewportUpperLeft.Add(pixelDeltaU.Add(pixelDeltaV).Scale(0.5))

	return Camera{
		aspectRatio,
		imageWidth,
		imageHeight,
		center,
		pixel00Loc,
		pixelDeltaU,
		pixelDeltaV,
	}
}

func (c Camera) rayColor(r ray.Ray, world hittable.Hittable) color.Color {
	var rec hittable.HitRecord

	if world.Hit(r, interval.NewInterval(0, math.Inf(1)), &rec) {
		return rec.Normal().Add(color.NewColor(1, 1, 1)).Scale(0.5)
	}

	unitDirection := r.Direction().Unit()
	a := 0.5 * (unitDirection.Y() + 1.0)
	return color.NewColor(1.0, 1.0, 1.0).Scale(1.0 - a).Add(color.NewColor(0.5, 0.7, 1.0).Scale(a))
}

func (c Camera) Render(out io.Writer, log io.Writer, world hittable.Hittable) {
	fmt.Fprintf(out, "P3\n%d %d\n255\n", c.imageWidth, c.imageHeight)

	for j := 0; j < c.imageHeight; j++ {
		fmt.Fprintf(log, "\rScanlines remaining: %d", c.imageHeight-j)
		for i := 0; i < c.imageWidth; i++ {
			pixelCenter := c.pixel00Loc.Add(c.pixelDeltaU.Scale(float64(i))).Add(c.pixelDeltaV.Scale(float64(j)))
			rayDirection := pixelCenter.Sub(c.center)
			ray := ray.NewRay(c.center, rayDirection)

			pixelColor := c.rayColor(ray, world)
			color.WriteColor(out, pixelColor)
		}
	}

	fmt.Fprintln(log, "\nDone.")
}
