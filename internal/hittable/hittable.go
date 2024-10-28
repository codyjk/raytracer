package hittable

import (
	"raytracer/internal/ray"
	"raytracer/internal/vector"
)

type HitRecord struct {
	p         vector.Point3
	normal    vector.Vec3
	t         float64
	frontFace bool
}

func (hr HitRecord) Normal() vector.Vec3 {
	return hr.normal
}

func (hr *HitRecord) setFaceNormal(r ray.Ray, outwardNormal vector.Vec3) {
	hr.frontFace = vector.Dot(r.Direction(), outwardNormal) < 0
	if hr.frontFace {
		hr.normal = outwardNormal
	} else {
		hr.normal = outwardNormal.Scale(-1.0)
	}
}

func (hr HitRecord) T() float64 {
	return hr.t
}

type Hittable interface {
	Hit(r ray.Ray, rayTmin float64, rayTmax float64, rec *HitRecord) bool
}

type HittableList struct {
	objects []Hittable
}

func NewHittableList() HittableList {
	return HittableList{[]Hittable{}}
}

func (hl *HittableList) Add(object Hittable) {
	hl.objects = append(hl.objects, object)
}

func (hl HittableList) Hit(r ray.Ray, rayTmin float64, rayTmax float64, rec *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := rayTmax

	for _, object := range hl.objects {
		if object.Hit(r, rayTmin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t
			*rec = tempRec
		}
	}

	return hitAnything
}
