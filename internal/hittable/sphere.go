package hittable

import (
	"math"
	"raytracer/internal/core"
	"raytracer/internal/interval"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type Sphere struct {
	center vector.Point3
	radius float64
	mat    core.Material
}

func NewSphere(center vector.Point3, radius float64, mat core.Material) Sphere {
	return Sphere{center, math.Max(0.0, radius), mat}
}

func (s Sphere) Center() vector.Point3 {
	return s.center
}

func (s Sphere) Radius() float64 {
	return s.radius
}

// This math solves for the ray-sphere intersection. The math is explained here:
// https://raytracing.github.io/books/RayTracingInOneWeekend.html#addingasphere/ray-sphereintersection
// https://raytracing.github.io/books/RayTracingInOneWeekend.html#surfacenormalsandmultipleobjects/simplifyingtheray-sphereintersectioncode
func (s Sphere) Hit(r ray.Ray, rayT interval.Interval, rec *core.HitRecord) bool {
	oc := s.center.Sub(r.Origin())
	a := r.Direction().LengthSquared()
	h := vector.Dot(r.Direction(), oc)
	c := oc.LengthSquared() - s.radius*s.radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range.
	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}

	rec.SetT(root)
	rec.SetPoint(r.At(rec.T()))
	outwardNormal := rec.Point().Sub(s.center).Div(s.radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.SetMaterial(s.mat)

	return true
}
