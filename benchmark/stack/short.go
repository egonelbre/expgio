// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Stack lays out child elements on top of each other,
// according to an alignment direction.
type Stack3 struct {
	// Alignment is the direction to align children
	// smaller than the available space.
	Alignment layout.Direction
}

// StackChild represents a child for a Stack layout.
type StackChild3 struct {
	expanded bool
	widget   layout.Widget
}

// Stacked returns a Stack child that is laid out with no minimum
// constraints and the maximum constraints passed to Stack.Layout.
func Stacked3(w layout.Widget) StackChild3 {
	return StackChild3{
		widget: w,
	}
}

// Expanded returns a Stack child with the minimum constraints set
// to the largest Stacked child. The maximum constraints are set to
// the same as passed to Stack.Layout.
func Expanded3(w layout.Widget) StackChild3 {
	return StackChild3{
		expanded: true,
		widget:   w,
	}
}

// Layout a stack of children. The position of the children are
// determined by the specified order, but Stacked children are laid out
// before Expanded children.
func (s Stack3) Layout3(gtx layout.Context, ax, bx, cx StackChild) layout.Dimensions {
	// Scratch space.
	var acall, bcall, ccall op.CallOp
	var adims, bdims, cdims layout.Dimensions

	var maxSZ image.Point
	// First lay out Stacked children.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}

	if !ax.expanded {
		macro := op.Record(gtx.Ops)
		adims = ax.widget(cgtx)
		acall = macro.Stop()
		maxSZ.X = max(maxSZ.X, adims.Size.X)
		maxSZ.Y = max(maxSZ.Y, adims.Size.Y)
	}
	if !bx.expanded {
		macro := op.Record(gtx.Ops)
		bdims = bx.widget(cgtx)
		bcall = macro.Stop()
		maxSZ.X = max(maxSZ.X, bdims.Size.X)
		maxSZ.Y = max(maxSZ.Y, bdims.Size.Y)
	}
	if !cx.expanded {
		macro := op.Record(gtx.Ops)
		cdims = cx.widget(cgtx)
		ccall = macro.Stop()
		maxSZ.X = max(maxSZ.X, cdims.Size.X)
		maxSZ.Y = max(maxSZ.Y, cdims.Size.Y)
	}

	if ax.expanded {
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		adims = ax.widget(cgtx)
		bcall = macro.Stop()
		maxSZ.X = max(maxSZ.X, adims.Size.X)
		maxSZ.Y = max(maxSZ.Y, adims.Size.Y)
	}
	if bx.expanded {
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		bdims = bx.widget(cgtx)
		ccall = macro.Stop()
		maxSZ.X = max(maxSZ.X, bdims.Size.X)
		maxSZ.Y = max(maxSZ.Y, bdims.Size.Y)
	}
	if cx.expanded {
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		cdims = cx.widget(cgtx)
		ccall = macro.Stop()
		maxSZ.X = max(maxSZ.X, cdims.Size.X)
		maxSZ.Y = max(maxSZ.Y, cdims.Size.Y)
	}

	maxSZ = gtx.Constraints.Constrain(maxSZ)
	var baseline int

	if ax.expanded {
		p := align(s.Alignment, adims.Size, maxSZ)
		trans := op.Offset(p).Push(gtx.Ops)
		acall.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := adims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - adims.Size.Y - p.Y
			}
		}
	}
	if bx.expanded {
		p := align(s.Alignment, bdims.Size, maxSZ)
		trans := op.Offset(p).Push(gtx.Ops)
		bcall.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := bdims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - bdims.Size.Y - p.Y
			}
		}
	}
	if cx.expanded {
		p := align(s.Alignment, cdims.Size, maxSZ)
		trans := op.Offset(p).Push(gtx.Ops)
		ccall.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := cdims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - cdims.Size.Y - p.Y
			}
		}
	}

	return layout.Dimensions{
		Size:     maxSZ,
		Baseline: baseline,
	}
}

func align(a layout.Direction, s, max image.Point) (p image.Point) {
	switch a {
	case layout.N, layout.S, layout.Center:
		p.X = (max.X - s.X) / 2
	case layout.NE, layout.SE, layout.E:
		p.X = max.X - s.X
	}
	switch a {
	case layout.W, layout.Center, layout.E:
		p.Y = (max.Y - s.Y) / 2
	case layout.SW, layout.S, layout.SE:
		p.Y = max.Y - s.Y
	}
	return p
}
