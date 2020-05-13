package main

import (
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
)

const defaultDuration = 250 * time.Millisecond

type Slider struct {
	Duration time.Duration

	last op.Ops
	next op.Ops

	t0       time.Time
	progress float32
}

func (slider *Slider) Push(gtx *layout.Context) {
	slider.last = slider.next
	slider.progress = 1.0
	slider.t0 = gtx.Now()
}

func (s *Slider) Layout(gtx *layout.Context, w layout.Widget) {
	var delta time.Duration
	if !s.t0.IsZero() {
		now := gtx.Now()
		delta = now.Sub(slider.t0)
		slider.t0 = now
	}

	if s.progress > 0 {
		duration := s.Duration
		if duration == 0 {
			duration = defaultDuration
		}
		s.progress -= float32(delta.Seconds()) / float32(duration.Seconds())
		if s.progress < 0 {
			s.progress = 0
		}
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	{
		prev := gtx.Ops
		s.next.Reset()
		gtx.Ops = &s.next
		w()
		gtx.Ops = prev
	}

	var stack op.StackOp
	stack.Push(gtx.Ops)

	if slider.progress > 0 {
		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X) * (s.progress - 1),
		}).Add(gtx.Ops)
		op.CallOp{Ops: &s.last}.Add(gtx.Ops)

		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X),
		}).Add(gtx.Ops)

		op.CallOp{Ops: &s.next}.Add(gtx.Ops)
	} else {
		op.CallOp{Ops: &s.next}.Add(gtx.Ops)
	}

	stack.Pop()
}
