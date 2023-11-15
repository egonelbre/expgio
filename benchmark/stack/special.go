package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Stack lays out child elements on top of each other,
// according to an alignment direction.
type ExpandedStacked struct {
}

func (s ExpandedStacked) Layout(gtx layout.Context, ex, st layout.Widget) layout.Dimensions {
	// First lay out Stacked children.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}

	macro := op.Record(gtx.Ops)
	stdims := st(cgtx)
	stcall := macro.Stop()

	cgtx.Constraints.Min = stdims.Size
	exdims := ex(cgtx)

	maxSZ := gtx.Constraints.Constrain(maxPt(stdims.Size, exdims.Size))
	var baseline int

	p := image.Point{
		X: (maxSZ.X - stdims.Size.X) / 2,
		Y: (maxSZ.Y - stdims.Size.Y) / 2,
	}

	trans := op.Offset(p).Push(gtx.Ops)
	stcall.Add(gtx.Ops)
	trans.Pop()
	if baseline == 0 {
		if b := stdims.Baseline; b != 0 {
			baseline = b + maxSZ.Y - stdims.Size.Y - p.Y
		}
	}

	return layout.Dimensions{
		Size:     maxSZ,
		Baseline: baseline,
	}

}

func maxPt(a, b image.Point) image.Point {
	return image.Point{
		X: max(a.X, b.X),
		Y: max(a.Y, b.Y),
	}
}
