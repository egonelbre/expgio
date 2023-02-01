package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	perlin "github.com/aquilax/go-perlin"
)

type Cloth struct {
	N int

	Time     float32
	Perlin   *perlin.Perlin
	Vertices []Vertex
	Spring   []Spring
}

type Vertex struct {
	Static   bool
	Position Vector
	Velocity Vector
	Force    Vector
}

type Spring struct {
	A, B   int
	Length float32
}

func NewCloth(n int) *Cloth {
	cloth := &Cloth{}
	cloth.N = n
	cloth.Perlin = perlin.NewPerlin(2, 2, 3, 14)
	cloth.Vertices = make([]Vertex, n*n)
	cloth.Spring = make([]Spring, 0, 4*n*n+2*n)

	fn := float32(n-1) / float32(n)
	fninv := 1 / fn
	fndiag := fninv * sqrt(2)
	_ = fndiag

	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			v := &cloth.Vertices[y*n+x]
			fx, fy := float32(x), float32(y)
			v.Position = Vector{fx / fn, fy / fn}
		}
	}

	for y := 0; y < n-1; y++ {
		for x := 0; x < n-1; x++ {
			index := y*n + x
			if x == 0 {
				cloth.Spring = append(cloth.Spring, Spring{index, index + n, fninv})
			}
			if y == 0 {
				cloth.Spring = append(cloth.Spring, Spring{index, index + 1, fninv})
			}

			cloth.Spring = append(cloth.Spring, Spring{index + 1, index + 1 + n, fninv})
			cloth.Spring = append(cloth.Spring, Spring{index + n, index + n + 1, fninv})
			cloth.Spring = append(cloth.Spring, Spring{index, index + n + 1, fndiag})
			cloth.Spring = append(cloth.Spring, Spring{index + 1, index + n, fndiag})
		}
	}

	for i := range cloth.Vertices[:n] {
		cloth.Vertices[i].Static = true
	}

	return cloth
}

func (cloth *Cloth) Update(dt float32) {
	dt *= 3

	cloth.Time += dt

	windf := sin(cloth.Time*0.1) * 0.005
	for i := range cloth.Vertices {
		v := &cloth.Vertices[i]
		wind := float32(cloth.Perlin.Noise3D(float64(v.Position.X), float64(v.Position.Y), float64(cloth.Time)))

		v.Force = Vector{0.002 + windf + wind*0.005, 0.05}
	}

	dt *= 10

	// this is unstable calculation at the moment

	for _, spring := range cloth.Spring {
		a, b := &cloth.Vertices[spring.A], &cloth.Vertices[spring.B]
		delta := b.Position.Sub(a.Position)
		stretch := delta.Len() - spring.Length
		force := delta.NormalizeTo(stretch * 0.5)
		a.Force = a.Force.Add(force)
		b.Force = b.Force.Sub(force)
	}

	for i := range cloth.Vertices {
		v := &cloth.Vertices[i]
		if v.Static {
			continue
		}

		v.Velocity = v.Velocity.Add(v.Force.Scale(dt))
		v.Velocity = v.Velocity.Scale(0.97)
	}

	for i := range cloth.Vertices {
		v := &cloth.Vertices[i]
		if v.Static {
			continue
		}

		v.Position = v.Position.Add(v.Velocity.Scale(dt))
	}
}

func (cloth *Cloth) Layout(gtx layout.Context) {
	op.Offset(image.Point{
		X: gtx.Constraints.Max.X / 4,
		Y: gtx.Constraints.Max.Y / 6,
	}).Add(gtx.Ops)
	scale := float32(gtx.Constraints.Max.X) * 0.5 / float32(cloth.N)

	if !*separate {
		var p clip.Path
		p.Begin(gtx.Ops)
		for _, spring := range cloth.Spring {
			a, b := &cloth.Vertices[spring.A], &cloth.Vertices[spring.B]
			p.MoveTo(f32.Point(a.Position).Mul(scale))
			p.LineTo(f32.Point(b.Position).Mul(scale))
		}
		spec := p.End()
		paint.FillShape(gtx.Ops, color.NRGBA{A: 0xFF}, clip.Stroke{
			Path:  spec,
			Width: 1,
		}.Op())
	} else {
		for _, spring := range cloth.Spring {
			var p clip.Path
			p.Begin(gtx.Ops)
			a, b := &cloth.Vertices[spring.A], &cloth.Vertices[spring.B]
			p.MoveTo(f32.Point(a.Position).Mul(scale))
			p.LineTo(f32.Point(b.Position).Mul(scale))
			paint.FillShape(gtx.Ops, color.NRGBA{A: 0xFF}, clip.Stroke{
				Path:  p.End(),
				Width: 1,
			}.Op())
		}

	}

}
