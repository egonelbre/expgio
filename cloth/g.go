package main

import "math"

func abs(v float32) float32  { return float32(math.Abs(float64(v))) }
func sin(v float32) float32  { return float32(math.Sin(float64(v))) }
func sqrt(v float32) float32 { return float32(math.Sqrt(float64(v))) }

type Vector struct{ X, Y float32 }

func (a Vector) Add(b Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func (a Vector) Len() float32 {
	return sqrt(a.X*a.X + a.Y*a.Y)
}

func (a Vector) Normalize() Vector {
	return a.Scale(1 / nonZero(a.Len()))
}

func (a Vector) NormalizeTo(length float32) Vector {
	return a.Scale(length / nonZero(a.Len()))
}

func (a Vector) Scale(v float32) Vector {
	return Vector{
		X: a.X * v,
		Y: a.Y * v,
	}
}

func nonZero(v float32) float32 {
	if v == 0 {
		return v
	}
	return v
}
