package material

import (
	"math"
	"raytracer/internal/color"
	"raytracer/internal/core"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type Metal struct {
	albedo color.Color
	fuzz   float64
}

func NewMetal(albedo color.Color, fuzz float64) Metal {
	fuzz = math.Min(fuzz, 1.0)
	return Metal{albedo, fuzz}
}

func (m Metal) Scatter(rIn ray.Ray, rec *core.HitRecord, attenuation *color.Color, scattered *ray.Ray) bool {
	reflected := vector.Reflect(rIn.Direction(), rec.Normal())
	reflected = reflected.Unit().Add(vector.RandomUnitVector().Scale(m.fuzz))
	*scattered = ray.NewRay(rec.Point(), reflected)
	*attenuation = m.albedo
	return vector.Dot(scattered.Direction(), rec.Normal()) > 0
}
