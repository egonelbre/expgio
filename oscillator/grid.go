package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

// Example how to do custom layouting when you have a need.

type Grid struct {
	// row and column proportional size
	Row []float64
	Col []float64

	Gap unit.Dp
}

type Span struct {
	Min, Max image.Point
	Widget   layout.Widget
}

func CellAt(row, col int, w layout.Widget) Span {
	return Span{
		Min:    image.Pt(col, row),
		Max:    image.Pt(col+1, row+1),
		Widget: w,
	}
}

func CellRows(row0, row1, col int, w layout.Widget) Span {
	return Span{
		Min:    image.Pt(col, row0),
		Max:    image.Pt(col+1, row1+1),
		Widget: w,
	}
}

func CellCols(row, col0, col1 int, w layout.Widget) Span {
	return Span{
		Min:    image.Pt(col0, row),
		Max:    image.Pt(col1+1, row+1),
		Widget: w,
	}
}

func (g Grid) Layout(gtx layout.Context, spans ...Span) layout.Dimensions {
	row := proportions(g.Row)
	col := proportions(g.Col)

	size := gtx.Constraints.Max
	// TODO: fix gap calculation for edges -- it actually should place straight to the edge.
	gap := gtx.Metric.Dp(g.Gap)

	coordToDisplay := func(p image.Point) image.Point {
		return image.Point{
			X: int(col[p.X] * float64(size.X)),
			Y: int(row[p.Y] * float64(size.Y)),
		}
	}

	for _, span := range spans {
		func() {
			min := coordToDisplay(span.Min).Add(image.Point{X: gap / 2, Y: gap / 2})
			max := coordToDisplay(span.Max).Sub(image.Point{X: gap / 2, Y: gap / 2})

			defer op.Offset(min).Push(gtx.Ops).Pop()
			size := max.Sub(min)
			defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()

			gtx := gtx
			gtx.Constraints = layout.Exact(size)
			_ = span.Widget(gtx)
		}()
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func proportions(in []float64) []float64 {
	out := make([]float64, len(in)+1)
	cum := 0.0
	for i, r := range in {
		out[i] = cum
		cum += r
	}
	out[len(out)-1] = cum
	for i := range out {
		out[i] /= cum
	}
	return out
}
