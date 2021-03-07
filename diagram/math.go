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

func (v Vector) Min(b Vector) Vector {
	return Vector{
		X: minUnit(v.X, b.X),
		Y: minUnit(v.Y, b.Y),
	}
}

func (v Vector) Max(b Vector) Vector {
	return Vector{
		X: maxUnit(v.X, b.X),
		Y: maxUnit(v.Y, b.Y),
	}
}

func minUnit(a, b Unit) Unit {
	if a < b {
		return a
	}
	return b
}

func maxUnit(a, b Unit) Unit {
	if a > b {
		return a
	}
	return b
}
