package interval

import "math"

type Interval struct {
	min, max float64
}

func NewInterval(min, max float64) Interval {
	return Interval{min, max}
}

func Empty() Interval {
	return NewInterval(math.Inf(1), math.Inf(-1))
}

func Universe() Interval {
	return NewInterval(math.Inf(-1), math.Inf(1))
}

func (i Interval) Min() float64 {
	return i.min
}

func (i Interval) Max() float64 {
	return i.max
}

func (i Interval) Surrounds(x float64) bool {
	return i.min < x && x < i.max
}

func (i Interval) Clamp(x float64) float64 {
	if x < i.min {
		return i.min
	} else if x > i.max {
		return i.max
	}
	return x
}
