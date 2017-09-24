package vector

import "math"

type Vec struct {
	X float64
	Y float64
}

const degToRad = math.Pi / 180.

func (vec Vec) Norm() float64 {
	return math.Sqrt(vec.Dot(vec))
}

func (vec Vec) Dot(b Vec) float64 {
	return vec.X*b.X + vec.Y*b.Y
}

func (vec Vec) Dist(b Vec) float64 {
	return vec.Minus(b).Norm()
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

func (vec Vec) Mult(dist float64) Vec {
	return Vec{
		X: vec.X * dist,
		Y: vec.Y * dist,
	}
}

func (vec Vec) MoveDegree(degree float64, distance float64) Vec{
	direction := Vec{
		X: distance * math.Sin(degree * degToRad),
		Y: distance * math.Cos(degree * degToRad),
	}

	return vec.Plus(direction)
}
