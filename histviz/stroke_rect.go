package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type StrokeRect struct {
	Rect  image.Rectangle
	Inset int
	Color color.NRGBA
}

func (stroke StrokeRect) Add(ops *op.Ops) {
	var path clip.Path

	// TODO: draw using boxes

	path.Begin(ops)

	inner := stroke.Rect

	path.MoveTo(f32.Point{
		X: float32(inner.Min.X),
		Y: float32(inner.Min.Y),
	})
	path.LineTo(f32.Point{
		X: float32(inner.Max.X),
		Y: float32(inner.Min.Y),
	})
	path.LineTo(f32.Point{
		X: float32(inner.Max.X),
		Y: float32(inner.Max.Y),
	})
	path.LineTo(f32.Point{
		X: float32(inner.Min.X),
		Y: float32(inner.Max.Y),
	})
	path.Close()

	outer := stroke.Rect.Inset(stroke.Inset)

	path.MoveTo(f32.Point{
		X: float32(outer.Min.X),
		Y: float32(outer.Min.Y),
	})
	path.LineTo(f32.Point{
		X: float32(outer.Min.X),
		Y: float32(outer.Max.Y),
	})
	path.LineTo(f32.Point{
		X: float32(outer.Max.X),
		Y: float32(outer.Max.Y),
	})
	path.LineTo(f32.Point{
		X: float32(outer.Max.X),
		Y: float32(outer.Min.Y),
	})
	path.Close()

	paint.FillShape(ops, stroke.Color, clip.Outline{
		Path: path.End(),
	}.Op())
}
