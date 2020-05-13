package main

import (
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
)

const defaultDuration = 1000 * time.Millisecond

type Slider struct {
	Duration time.Duration

	last op.Ops
	next op.Ops

	t0     time.Time
	offset float32
}

func (slider *Slider) PushLeft(gtx *layout.Context) {
	slider.last = slider.next
	slider.offset = 1.0
	slider.t0 = gtx.Now()
}

func (slider *Slider) PushRight(gtx *layout.Context) {
	slider.last = slider.next
	slider.offset = -1.0
	slider.t0 = gtx.Now()
}

func (s *Slider) Layout(gtx *layout.Context, w layout.Widget) {
	var delta time.Duration
	if !s.t0.IsZero() {
		now := gtx.Now()
		delta = now.Sub(slider.t0)
		slider.t0 = now
	}

	if s.offset != 0 {
		duration := s.Duration
		if duration == 0 {
			duration = defaultDuration
		}
		movement := float32(delta.Seconds()) / float32(duration.Seconds())
		if s.offset < 0 {
			s.offset += movement
			if s.offset >= 0 {
				s.offset = 0
			}
		} else {
			s.offset -= movement
			if s.offset <= 0 {
				s.offset = 0
			}
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

	if slider.offset == 0 {
		op.CallOp{Ops: &s.next}.Add(gtx.Ops)
		return
	}

	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	if slider.offset > 0 {
		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X) * (s.offset - 1),
		}).Add(gtx.Ops)

		op.CallOp{Ops: &s.last}.Add(gtx.Ops)

		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X),
		}).Add(gtx.Ops)

		op.CallOp{Ops: &s.next}.Add(gtx.Ops)
	} else {
		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X) * (s.offset + 1),
		}).Add(gtx.Ops)

		op.CallOp{Ops: &s.last}.Add(gtx.Ops)

		op.TransformOp{}.Offset(f32.Point{
			X: float32(-gtx.Dimensions.Size.X),
		}).Add(gtx.Ops)

		op.CallOp{Ops: &s.next}.Add(gtx.Ops)
	}
}
