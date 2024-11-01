package material

import (
	"raytracer/internal/color"
	"raytracer/internal/core"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type Lambertian struct {
	albedo color.Color
}

func NewLambertian(albedo color.Color) Lambertian {
	return Lambertian{albedo}
}

func (l Lambertian) Scatter(rIn ray.Ray, rec *core.HitRecord, attenuation *color.Color, scattered *ray.Ray) bool {
	scatterDirection := rec.Normal().Add(vector.RandomUnitVector())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal()
	}

	*scattered = ray.NewRay(rec.Point(), scatterDirection)
	*attenuation = l.albedo
	return true
}
