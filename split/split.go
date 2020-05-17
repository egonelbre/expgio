package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
)

type Split struct {
}

func (s *Split) Layout(gtx *layout.Context, left, right layout.Widget) {
	savedConstraints := gtx.Constraints
	defer func() {
		gtx.Constraints = savedConstraints
		gtx.Dimensions.Size = image.Point{
			X: savedConstraints.Width.Max,
			Y: savedConstraints.Height.Max,
		}
	}()
	gtx.Constraints.Height.Min = gtx.Constraints.Height.Max

	leftsize := gtx.Constraints.Width.Min / 2
	rightsize := gtx.Constraints.Width.Min - leftsize

	{
		var stack op.StackOp
		stack.Push(gtx.Ops)

		gtx.Constraints.Width.Min = leftsize
		gtx.Constraints.Width.Max = leftsize
		left()

		stack.Pop()
	}

	{
		var stack op.StackOp
		stack.Push(gtx.Ops)

		gtx.Constraints.Width.Min = rightsize
		gtx.Constraints.Width.Max = rightsize

		op.TransformOp{}.Offset(f32.Point{
			X: float32(leftsize),
		}).Add(gtx.Ops)
		right()

		stack.Pop()
	}
}
