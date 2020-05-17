package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

type Split struct {
	// Ratio keeps the current layout.
	// 0 is center, -1 completely to the left, 1 completely to the right.
	Ratio float32
	// Bar is the width for resizing the layout
	Bar int

	drag   bool
	dragID pointer.ID
	dragX  float32
}

const defaultBarWidth = 10

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

	bar := s.Bar
	if bar <= 0 {
		bar = defaultBarWidth
	}

	proportion := (s.Ratio + 1) / 2
	leftsize := int(proportion*float32(gtx.Constraints.Width.Max) - float32(bar))

	rightoffset := leftsize + bar
	rightsize := gtx.Constraints.Width.Max - rightoffset

	{ // handle input
		for _, ev := range gtx.Events(s) {
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			switch e.Type {
			case pointer.Press:
				if s.drag {
					break
				}

				s.drag = true
				s.dragID = e.PointerID
				s.dragX = e.Position.X

			case pointer.Move:
				if !s.drag || s.dragID != e.PointerID {
					break
				}

				deltaX := e.Position.X - s.dragX
				s.dragX = e.Position.X

				deltaRatio := deltaX * 2 / float32(gtx.Constraints.Width.Max)
				s.Ratio += deltaRatio

			case pointer.Release:
				fallthrough
			case pointer.Cancel:
				if !s.drag || s.dragID != e.PointerID {
					break
				}
				s.drag = false
			}
		}

		// register for input
		barRect := image.Rect(leftsize, 0, rightoffset, gtx.Constraints.Width.Max)
		pointer.Rect(barRect).Add(gtx.Ops)
		pointer.InputOp{Key: s, Grab: s.drag}.Add(gtx.Ops)
	}

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
			X: float32(rightoffset),
		}).Add(gtx.Ops)
		right()

		stack.Pop()
	}
}
