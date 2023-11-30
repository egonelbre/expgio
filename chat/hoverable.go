package main

import (
	"image"
	"image/color"
	"time"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Hoverable struct {
	hovered bool
}

func (h *Hoverable) Layout(gtx layout.Context) layout.Dimensions {
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()

	event.Op(gtx.Ops, h)

	for {
		e, ok := gtx.Event(pointer.Filter{
			Target: h,
			Kinds:  pointer.Enter | pointer.Leave,
		})
		if !ok {
			break
		}

		ev, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		switch ev.Kind {
		case pointer.Enter:
			h.hovered = true
		case pointer.Leave:
			h.hovered = false
		}
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Min,
	}
}

func (h *Hoverable) Active() bool { return h.hovered }

type AnimationTimer struct {
	Duration time.Duration

	progress time.Duration
	last     time.Time
}

func (anim *AnimationTimer) Progress() float32 {
	return float32(anim.progress) / float32(anim.Duration)
}

func (anim *AnimationTimer) Update(gtx layout.Context, active bool) float32 {
	delta := gtx.Now.Sub(anim.last)
	anim.last = gtx.Now

	if delta > 15*time.Millisecond {
		delta = 15 * time.Millisecond
	}

	if active {
		if anim.progress < anim.Duration {
			anim.progress += delta
			gtx.Execute(op.InvalidateCmd{})
			if anim.progress > anim.Duration {
				anim.progress = anim.Duration
			}
		}
	} else {
		if anim.progress > 0 {
			anim.progress -= delta
			gtx.Execute(op.InvalidateCmd{})
			if anim.progress < 0 {
				anim.progress = 0
			}
		}
	}

	return anim.Progress()
}

// BorderSmooth lays out a widget and draws a border inside it, with non-pixel-perfect border.
type BorderSmooth struct {
	Color        color.NRGBA
	CornerRadius unit.Dp
	Width        float32
}

func (b BorderSmooth) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)

	paint.FillShape(gtx.Ops, b.Color, clip.Stroke{
		Path:  clip.UniformRRect(image.Rectangle{Max: dims.Size}, gtx.Dp(b.CornerRadius)).Path(gtx.Ops),
		Width: b.Width,
	}.Op())

	return dims
}
