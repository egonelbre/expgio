package main

type Unit int

type Box struct {
	Pos  Vector
	Size Vector
}

type Vector struct{ X, Y Unit }

func V(x, y Unit) Vector { return Vector{X: x, Y: y} }

func (v Vector) Add(b Vector) Vector {
	return Vector{
		X: v.X + b.X,
		Y: v.Y + b.Y,
	}
}

func (v Vector) Sub(b Vector) Vector {
	return Vector{
		X: v.X - b.X,
		Y: v.Y - b.Y,
	}
}
