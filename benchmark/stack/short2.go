// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Stack lays out child elements on top of each other,
// according to an alignment direction.
type Stack4 struct {
	// Alignment is the direction to align children
	// smaller than the available space.
	Alignment layout.Direction
}

// StackChild represents a child for a Stack layout.
type StackChild4 struct {
	expanded bool
	widget   layout.Widget
	call     op.CallOp
	dims     layout.Dimensions
}

// Stacked returns a Stack child that is laid out with no minimum
// constraints and the maximum constraints passed to Stack.Layout.
func Stacked4(w layout.Widget) StackChild4 {
	return StackChild4{
		widget: w,
	}
}

// Expanded returns a Stack child with the minimum constraints set
// to the largest Stacked child. The maximum constraints are set to
// the same as passed to Stack.Layout.
func Expanded4(w layout.Widget) StackChild4 {
	return StackChild4{
		expanded: true,
		widget:   w,
	}
}

// Layout a stack of children. The position of the children are
// determined by the specified order, but Stacked children are laid out
// before Expanded children.
func (s Stack4) Layout3(gtx layout.Context, ax, bx, cx StackChild) layout.Dimensions {
	var maxSZ image.Point
	// First lay out Stacked children.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}

	stacked := func(ch *StackChild) {
		if ch.expanded {
			return
		}
		macro := op.Record(gtx.Ops)
		ch.dims = ch.widget(cgtx)
		ch.call = macro.Stop()
		updateMax(&maxSZ, ch.dims.Size)
	}

	stacked(&ax)
	stacked(&bx)
	stacked(&cx)

	expanded := func(ch *StackChild) {
		if !ch.expanded {
			return
		}
		macro := op.Record(gtx.Ops)
		cgtx.Constraints.Min = maxSZ
		ch.dims = ch.widget(cgtx)
		ch.call = macro.Stop()
		updateMax(&maxSZ, ch.dims.Size)
	}

	expanded(&ax)
	expanded(&bx)
	expanded(&cx)

	maxSZ = gtx.Constraints.Constrain(maxSZ)
	var baseline int

	if ax.expanded {
		p := align(s.Alignment, ax.dims.Size, maxSZ)
		trans := op.Offset(p).Push(gtx.Ops)
		ax.call.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := ax.dims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - ax.dims.Size.Y - p.Y
			}
		}
	}
	if bx.expanded {
		p := align(s.Alignment, bx.dims.Size, maxSZ)
		trans := op.Offset(p).Push(gtx.Ops)
		bx.call.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := bx.dims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - bx.dims.Size.Y - p.Y
			}
		}
	}
	if cx.expanded {
		p := align(s.Alignment, cx.dims.Size, maxSZ)
		trans := op.Offset(p).Push(gtx.Ops)
		cx.call.Add(gtx.Ops)
		trans.Pop()
		if baseline == 0 {
			if b := cx.dims.Baseline; b != 0 {
				baseline = b + maxSZ.Y - cx.dims.Size.Y - p.Y
			}
		}
	}

	return layout.Dimensions{
		Size:     maxSZ,
		Baseline: baseline,
	}
}

func updateMax(a *image.Point, b image.Point) {
	a.X = max(a.X, b.X)
	a.Y = max(a.Y, b.Y)
}
