package vector

import "math"

type Vec struct {
	X float64
	Y float64
}

func (vec Vec) Norm() float64 {
	return math.Sqrt(vec.Dot(vec))
}

func (a Vec) Dot(b Vec) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vec) Dist(b Vec) float64 {
	return a.Minus(b).Norm()
}

func (vec Vec) Plus(vec2 Vec) Vec {
	return Vec{
		X: vec.X + vec2.X,
		Y: vec.Y + vec2.Y,
	}
}

func (vec Vec) Minus(vec2 Vec) Vec {
	return Vec{
		X: vec.X - vec2.X,
		Y: vec.Y - vec2.Y,
	}
}
