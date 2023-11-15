// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Stack lays out child elements on top of each other,
// according to an alignment direction.
type Stack5 struct {
	// Alignment is the direction to align children
	// smaller than the available space.
	Alignment layout.Direction

	Children []StackChild5
}

// StackChild represents a child for a Stack layout.
type StackChild5 struct {
	expanded bool
	widget   layout.Widget

	// Scratch space.
	call op.CallOp
	dims layout.Dimensions
}

// Stacked returns a Stack child that is laid out with no minimum
// constraints and the maximum constraints passed to Stack.Layout.
func Stacked5(w layout.Widget) StackChild5 {
	return StackChild5{
		widget: w,
	}
}

// Expanded returns a Stack child with the minimum constraints set
// to the largest Stacked child. The maximum constraints are set to
// the same as passed to Stack.Layout.
func Expanded5(w layout.Widget) StackChild5 {
	return StackChild5{
		expanded: true,
		widget:   w,
	}
}

// Layout a stack of children. The position of the children are
// determined by the specified order, but Stacked children are laid out
// before Expanded children.
func (s Stack5) Layout(gtx layout.Context) layout.Dimensions {
	var maxSZ image.Point
	// First lay out Stacked children.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}
	for i := range s.Children {
		w := &s.Children[i]
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
	for i := range s.Children {
		w := &s.Children[i]
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
	for i := range s.Children {
		w := &s.Children[i]
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
