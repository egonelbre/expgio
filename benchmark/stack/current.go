// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Stack lays out child elements on top of each other,
// according to an alignment direction.
type Stack struct {
	// Alignment is the direction to align children
	// smaller than the available space.
	Alignment layout.Direction
}

// StackChild represents a child for a Stack layout.
type StackChild struct {
	expanded bool
	widget   layout.Widget

	// Scratch space.
	call op.CallOp
	dims layout.Dimensions
}

// Stacked returns a Stack child that is laid out with no minimum
// constraints and the maximum constraints passed to Stack.Layout.
func Stacked(w layout.Widget) StackChild {
	return StackChild{
		widget: w,
	}
}

// Expanded returns a Stack child with the minimum constraints set
// to the largest Stacked child. The maximum constraints are set to
// the same as passed to Stack.Layout.
func Expanded(w layout.Widget) StackChild {
	return StackChild{
		expanded: true,
		widget:   w,
	}
}

// Layout a stack of children. The position of the children are
// determined by the specified order, but Stacked children are laid out
// before Expanded children.
func (s Stack) Layout(gtx layout.Context, children ...StackChild) layout.Dimensions {
	var maxSZ image.Point
	// First lay out Stacked children.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}
	for i := range children {
		w := &children[i]
		if w.expanded {
			continue
		}
		macro := op.Record(gtx.Ops)
		w.dims = w.widget(cgtx)
		w.call = macro.Stop()
		if w := w.dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := w.dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
	}
	// Then lay out Expanded children.
	for i := range children {
		w := &children[i]
		if !w.expanded {
			continue
		}
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		w.dims = w.widget(cgtx)
		w.call = macro.Stop()
		if w := w.dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := w.dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
	}

	maxSZ = gtx.Constraints.Constrain(maxSZ)
	var baseline int
	for i := range children {
		w := &children[i]
		sz := w.dims.Size
		var p image.Point
		switch s.Alignment {
		case layout.N, layout.S, layout.Center:
			p.X = (maxSZ.X - sz.X) / 2
		case layout.NE, layout.SE, layout.E:
			p.X = maxSZ.X - sz.X
		}
		switch s.Alignment {
		case layout.W, layout.Center, layout.E:
			p.Y = (maxSZ.Y - sz.Y) / 2
		case layout.SW, layout.S, layout.SE:
			p.Y = maxSZ.Y - sz.Y
		}
		trans := op.Offset(p).Push(gtx.Ops)
		w.call.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := w.dims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - sz.Y - p.Y
			}
		}
	}
	return layout.Dimensions{
		Size:     maxSZ,
		Baseline: baseline,
	}
}

// Layout a stack of children. The position of the children are
// determined by the specified order, but Stacked children are laid out
// before Expanded children.
func (s Stack) Layout3(gtx layout.Context, ax, bx, cx StackChild) layout.Dimensions {
	var maxSZ image.Point
	// First lay out Stacked children.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}

	if !ax.expanded {
		macro := op.Record(gtx.Ops)
		dims := ax.widget(cgtx)
		call := macro.Stop()
		if w := dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
		ax.call = call
		ax.dims = dims
	}
	if !bx.expanded {
		macro := op.Record(gtx.Ops)
		dims := bx.widget(cgtx)
		call := macro.Stop()
		if w := dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
		bx.call = call
		bx.dims = dims
	}
	if !cx.expanded {
		macro := op.Record(gtx.Ops)
		dims := cx.widget(cgtx)
		call := macro.Stop()
		if w := dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
		cx.call = call
		cx.dims = dims
	}

	if ax.expanded {
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		dims := ax.widget(cgtx)
		call := macro.Stop()
		if w := dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
		ax.call = call
		ax.dims = dims
	}
	if bx.expanded {
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		dims := bx.widget(cgtx)
		call := macro.Stop()
		if w := dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
		bx.call = call
		bx.dims = dims
	}
	if cx.expanded {
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		dims := cx.widget(cgtx)
		call := macro.Stop()
		if w := dims.Size.X; w > maxSZ.X {
			maxSZ.X = w
		}
		if h := dims.Size.Y; h > maxSZ.Y {
			maxSZ.Y = h
		}
		cx.call = call
		cx.dims = dims
	}

	maxSZ = gtx.Constraints.Constrain(maxSZ)
	var baseline int

	if ax.expanded {
		sz := ax.dims.Size
		var p image.Point
		switch s.Alignment {
		case layout.N, layout.S, layout.Center:
			p.X = (maxSZ.X - sz.X) / 2
		case layout.NE, layout.SE, layout.E:
			p.X = maxSZ.X - sz.X
		}
		switch s.Alignment {
		case layout.W, layout.Center, layout.E:
			p.Y = (maxSZ.Y - sz.Y) / 2
		case layout.SW, layout.S, layout.SE:
			p.Y = maxSZ.Y - sz.Y
		}
		trans := op.Offset(p).Push(gtx.Ops)
		ax.call.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := ax.dims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - sz.Y - p.Y
			}
		}
	}
	if bx.expanded {
		sz := bx.dims.Size
		var p image.Point
		switch s.Alignment {
		case layout.N, layout.S, layout.Center:
			p.X = (maxSZ.X - sz.X) / 2
		case layout.NE, layout.SE, layout.E:
			p.X = maxSZ.X - sz.X
		}
		switch s.Alignment {
		case layout.W, layout.Center, layout.E:
			p.Y = (maxSZ.Y - sz.Y) / 2
		case layout.SW, layout.S, layout.SE:
			p.Y = maxSZ.Y - sz.Y
		}
		trans := op.Offset(p).Push(gtx.Ops)
		bx.call.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := bx.dims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - sz.Y - p.Y
			}
		}
	}
	if cx.expanded {
		sz := cx.dims.Size
		var p image.Point
		switch s.Alignment {
		case layout.N, layout.S, layout.Center:
			p.X = (maxSZ.X - sz.X) / 2
		case layout.NE, layout.SE, layout.E:
			p.X = maxSZ.X - sz.X
		}
		switch s.Alignment {
		case layout.W, layout.Center, layout.E:
			p.Y = (maxSZ.Y - sz.Y) / 2
		case layout.SW, layout.S, layout.SE:
			p.Y = maxSZ.Y - sz.Y
		}
		trans := op.Offset(p).Push(gtx.Ops)
		cx.call.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := cx.dims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - sz.Y - p.Y
			}
		}
	}

	return layout.Dimensions{
		Size:     maxSZ,
		Baseline: baseline,
	}
}
