package material

import (
	"raytracer/internal/color"
	"raytracer/internal/core"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type Dielectric struct {
	refractionIndex float64
}

func NewDielectric(ri float64) Dielectric {
	return Dielectric{ri}
}

func (d Dielectric) Scatter(rIn ray.Ray, rec *core.HitRecord, attenuation *color.Color, scattered *ray.Ray) bool {
	*attenuation = color.NewColor(1, 1, 1)

	ri := d.refractionIndex
	if rec.FrontFace() {
		ri = 1.0 / d.refractionIndex
	}

	unitDirection := rIn.Direction().Unit()
	refracted := vector.Refract(unitDirection, rec.Normal(), ri)

	*scattered = ray.NewRay(rec.Point(), refracted)

	return true
}
