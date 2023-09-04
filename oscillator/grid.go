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

	Gap    unit.Dp
	Margin unit.Dp
}

type Span struct {
	Min, Max image.Point
	Widget   layout.Widget
}

func CellAt(row, col int, w layout.Widget) Span {
	return Span{
		Min:    image.Pt(col, row),
		Max:    image.Pt(col, row),
		Widget: w,
	}
}

func CellRows(row0, row1, col int, w layout.Widget) Span {
	return Span{
		Min:    image.Pt(col, row0),
		Max:    image.Pt(col, row1),
		Widget: w,
	}
}

func CellCols(row, col0, col1 int, w layout.Widget) Span {
	return Span{
		Min:    image.Pt(col0, row),
		Max:    image.Pt(col1, row),
		Widget: w,
	}
}

func (g Grid) Layout(gtx layout.Context, spans ...Span) layout.Dimensions {
	var colBuffer, rowBuffer [16]float64

	// cumulative weights of row or column
	col := cumulativeProp(g.Col, colBuffer[:])
	row := cumulativeProp(g.Row, rowBuffer[:])

	// total size of the context
	size := gtx.Constraints.Max
	// gap between cells
	gap := gtx.Metric.Dp(g.Gap)
	// margin outside the grid
	margin := gtx.Metric.Dp(g.Margin)
	// total size of cells (excluding gaps and margins)
	display := image.Point{
		X: size.X - gap*(len(g.Col)-1) - margin*2,
		Y: size.Y - gap*(len(g.Row)-1) - margin*2,
	}

	// calculates the coordinates based on the cell coordinates.
	// bottomRight = 1, means it should calculate the bottom right corner of the cell.
	cellPosition := func(p image.Point, bottomRight int) image.Point {
		return image.Point{
			X: margin + int(col[p.X+bottomRight]*float64(display.X)+float64(gap*p.X)),
			Y: margin + int(row[p.Y+bottomRight]*float64(display.Y)+float64(gap*p.Y)),
		}
	}

	for _, span := range spans {
		func() {
			min := cellPosition(span.Min, 0)
			max := cellPosition(span.Max, 1)

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

func cumulativeProp(in []float64, out []float64) []float64 {
	if cap(out) < len(in)+1 {
		out = make([]float64, len(in)+1)
	} else {
		out = out[:len(in)+1]
	}
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
