package core

import (
	"raytracer/internal/color"
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type HitRecord struct {
	p         vector.Point3
	normal    vector.Vec3
	t         float64
	frontFace bool
	mat       Material
}

func (hr HitRecord) Point() vector.Point3 {
	return hr.p
}

func (hr *HitRecord) SetPoint(p vector.Point3) {
	hr.p = p
}

func (hr HitRecord) Normal() vector.Vec3 {
	return hr.normal
}

func (hr HitRecord) T() float64 {
	return hr.t
}

func (hr HitRecord) Material() Material {
	return hr.mat
}

func (hr HitRecord) FrontFace() bool {
	return hr.frontFace
}

func (hr *HitRecord) SetT(t float64) {
	hr.t = t
}

func (hr *HitRecord) SetFaceNormal(r ray.Ray, outwardNormal vector.Vec3) {
	hr.frontFace = vector.Dot(r.Direction(), outwardNormal) < 0
	if hr.frontFace {
		hr.normal = outwardNormal
	} else {
		hr.normal = outwardNormal.Scale(-1.0)
	}
}

func (hr *HitRecord) SetMaterial(mat Material) {
	hr.mat = mat
}

type Material interface {
	Scatter(rIn ray.Ray, rec *HitRecord, attenuation *color.Color, scattered *ray.Ray) bool
}
