package hittable

import (
	"raytracer/internal/core"
	"raytracer/internal/interval"
	"raytracer/internal/ray"
)

type Hittable interface {
	Hit(r ray.Ray, rayT interval.Interval, rec *core.HitRecord) bool
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

func (hl HittableList) Hit(r ray.Ray, rayT interval.Interval, rec *core.HitRecord) bool {
	var tempRec core.HitRecord
	hitAnything := false
	closestSoFar := rayT.Max()

	for _, object := range hl.objects {
		if object.Hit(r, interval.NewInterval(rayT.Min(), closestSoFar), &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T()
			*rec = tempRec
		}
	}

	return hitAnything
}
