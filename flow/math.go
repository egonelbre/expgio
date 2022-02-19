package main

import (
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type Unit float32

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

// Circle represents the clip area of a circle.
type Circle struct {
	Center f32.Point
	Radius float32
}

// Op returns the op for the circle.
func (c Circle) Op(ops *op.Ops) clip.Op {
	return clip.Outline{Path: c.Path(ops)}.Op()
}

// Push the circle clip on the clip stack.
func (c Circle) Push(ops *op.Ops) clip.Stack {
	return c.Op(ops).Push(ops)
}

// Path returns the PathSpec for the circle.
func (c Circle) Path(ops *op.Ops) clip.PathSpec {
	var p clip.Path
	p.Begin(ops)

	center := c.Center
	r := c.Radius

	// https://pomax.github.io/bezierinfo/#circles_cubic.
	const q = 4 * (math.Sqrt2 - 1) / 3

	curve := r * q
	top := f32.Point{X: center.X, Y: center.Y - r}

	p.MoveTo(top)
	p.CubeTo(
		f32.Point{X: center.X + curve, Y: center.Y - r},
		f32.Point{X: center.X + r, Y: center.Y - curve},
		f32.Point{X: center.X + r, Y: center.Y},
	)
	p.CubeTo(
		f32.Point{X: center.X + r, Y: center.Y + curve},
		f32.Point{X: center.X + curve, Y: center.Y + r},
		f32.Point{X: center.X, Y: center.Y + r},
	)
	p.CubeTo(
		f32.Point{X: center.X - curve, Y: center.Y + r},
		f32.Point{X: center.X - r, Y: center.Y + curve},
		f32.Point{X: center.X - r, Y: center.Y},
	)
	p.CubeTo(
		f32.Point{X: center.X - r, Y: center.Y - curve},
		f32.Point{X: center.X - curve, Y: center.Y - r},
		top,
	)
	return p.End()
}
