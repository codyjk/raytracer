package material

import (
	"raytracer/internal/color"
	"raytracer/internal/core"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type Metal struct {
	albedo color.Color
}

func NewMetal(albedo color.Color) Metal {
	return Metal{albedo}
}

func (m Metal) Scatter(rIn ray.Ray, rec *core.HitRecord, attenuation *color.Color, scattered *ray.Ray) bool {
	reflected := vector.Reflect(rIn.Direction(), rec.Normal())
	*scattered = ray.NewRay(rec.Point(), reflected)
	*attenuation = m.albedo
	return true
}
