package vector

import (
	"fmt"
	"math"
	"raytracer/internal/util"
)

// Represents a 3d vector.
type Vec3 struct {
	e [3]float64
}

func NewVec3(e0, e1, e2 float64) Vec3 {
	return Vec3{e: [3]float64{e0, e1, e2}}
}

// An alias for Vec3, used for geometric clarity.
type Point3 = Vec3

func NewPoint3(e0, e1, e2 float64) Point3 {
	return Point3{e: [3]float64{e0, e1, e2}}
}

func ZeroVec3() Vec3 {
	return Vec3{e: [3]float64{0, 0, 0}}
}

func (v Vec3) X() float64                { return v.e[0] }
func (v Vec3) Y() float64                { return v.e[1] }
func (v Vec3) Z() float64                { return v.e[2] }
func (v Vec3) At(i int) float64          { return v.e[i] }
func (v *Vec3) Set(i int, value float64) { v.e[i] = value }

func (v Vec3) Neg() Vec3 {
	return NewVec3(-v.e[0], -v.e[1], -v.e[2])
}

func (v Vec3) Scale(t float64) Vec3 {
	return NewVec3(v.e[0]*t, v.e[1]*t, v.e[2]*t)
}

func (v Vec3) Div(t float64) Vec3 {
	return v.Scale(1 / t)
}

func (v Vec3) LengthSquared() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v Vec3) Add(other Vec3) Vec3 {
	return NewVec3(v.e[0]+other.e[0], v.e[1]+other.e[1], v.e[2]+other.e[2])
}

func (v Vec3) AddScaled(other Vec3, t float64) Vec3 {
	return NewVec3(v.e[0]+t*other.e[0], v.e[1]+t*other.e[1], v.e[2]+t*other.e[2])
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return NewVec3(v.e[0]-other.e[0], v.e[1]-other.e[1], v.e[2]-other.e[2])
}

func (v Vec3) Mul(other Vec3) Vec3 {
	return NewVec3(v.e[0]*other.e[0], v.e[1]*other.e[1], v.e[2]*other.e[2])
}

// Unit returns a unit vector in the same direction
func (v Vec3) Unit() Vec3 {
	return v.Div(v.Length())
}

func (v Vec3) String() string {
	return fmt.Sprintf("%v %v %v", v.e[0], v.e[1], v.e[2])
}

// Utility functions

func Dot(u, v Vec3) float64 {
	return u.e[0]*v.e[0] + u.e[1]*v.e[1] + u.e[2]*v.e[2]
}

func Cross(u, v Vec3) Vec3 {
	return NewVec3(
		u.e[1]*v.e[2]-u.e[2]*v.e[1],
		u.e[2]*v.e[0]-u.e[0]*v.e[2],
		u.e[0]*v.e[1]-u.e[1]*v.e[0],
	)
}

func Random() Vec3 {
	return NewVec3(util.RandomFloat(), util.RandomFloat(), util.RandomFloat())
}

func RandomFromRange(min, max float64) Vec3 {
	return NewVec3(util.RandomFloatFromRange(min, max), util.RandomFloatFromRange(min, max), util.RandomFloatFromRange(min, max))
}

func RandomUnitVector() Vec3 {
	for {
		p := RandomFromRange(-1.0, 1.0)
		lenSq := p.LengthSquared()
		// Very small values can underflow to 0 when squared, so add a lower bound.
		if 1e-160 < lenSq && lenSq <= 1 {
			return p.Div(math.Sqrt(lenSq))
		}
	}
}

func RandomOnHemisphere(normal Vec3) Vec3 {
	onUnitSphere := RandomUnitVector()
	if Dot(onUnitSphere, normal) > 0.0 { // In the same hemisphere as the normal
		return onUnitSphere
	} else {
		return onUnitSphere.Neg()
	}
}
