package ray

import "raytracer/internal/vector"

type Ray struct {
	origin    vector.Point3
	direction vector.Vec3
}

func (r Ray) Origin() vector.Point3 {
	return r.origin
}

func (r Ray) Direction() vector.Vec3 {
	return r.direction
}

func (r Ray) At(t float64) vector.Point3 {
	return r.origin.Add(r.direction.Scale(t))
}
