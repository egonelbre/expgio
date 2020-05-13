package main

import (
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
)

const defaultDuration = 300 * time.Millisecond

type Slider struct {
	Duration time.Duration

	push int

	last *op.Ops
	next *op.Ops

	t0     time.Time
	offset float32
}

func (s *Slider) PushLeft(gtx *layout.Context) {
	s.push = 1
}

func (s *Slider) PushRight(gtx *layout.Context) {
	s.push = -1
}

func (s *Slider) Layout(gtx *layout.Context, w layout.Widget) {
	if s.push != 0 {
		s.last, s.next = s.next, new(op.Ops)
		s.offset = float32(s.push)
		s.t0 = gtx.Now()
		s.push = 0
	}

	var delta time.Duration
	if !s.t0.IsZero() {
		now := gtx.Now()
		delta = now.Sub(s.t0)
		s.t0 = now
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
		if s.next == nil {
			s.next = new(op.Ops)
		}
		s.next.Reset()
		gtx.Ops = s.next
		w()
		gtx.Ops = prev
	}

	if s.offset == 0 {
		op.CallOp{Ops: s.next}.Add(gtx.Ops)
		return
	}

	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	offset := absfn(s.offset, easeInOutCubic)

	if s.offset > 0 {
		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X) * (offset - 1),
		}).Add(gtx.Ops)
		op.CallOp{Ops: s.last}.Add(gtx.Ops)

		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X),
		}).Add(gtx.Ops)
		op.CallOp{Ops: s.next}.Add(gtx.Ops)
	} else {
		op.TransformOp{}.Offset(f32.Point{
			X: float32(gtx.Dimensions.Size.X) * (offset + 1),
		}).Add(gtx.Ops)
		op.CallOp{Ops: s.last}.Add(gtx.Ops)

		op.TransformOp{}.Offset(f32.Point{
			X: float32(-gtx.Dimensions.Size.X),
		}).Add(gtx.Ops)
		op.CallOp{Ops: s.next}.Add(gtx.Ops)
	}
}

func absfn(t float32, fn func(float32) float32) float32 {
	if t < 0 {
		return -fn(-t)
	}
	return fn(t)
}

func easeInOutCubic(t float32) float32 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return (t-1)*(2*t-2)*(2*t-2) + 1
}

func hesitant(t float32) float32 {
	s := float32(math.Sin(float64(t * math.Pi * 2)))
	return t + s*s*s/2
}

// return 0.0005 * (t - 1) * (t - 1) * t * (44851 - 224256 * t + 224256 * t * t)
