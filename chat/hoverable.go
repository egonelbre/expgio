package main

import (
	"image"
	"time"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

type Hoverable struct {
	hovered bool
}

func (h *Hoverable) Layout(gtx layout.Context) layout.Dimensions {
	h.update(gtx)

	stack := op.Push(gtx.Ops)
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   h,
		Types: pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)
	stack.Pop()

	return layout.Dimensions{
		Size: gtx.Constraints.Min,
	}
}

func (h *Hoverable) Active() bool { return h.hovered }

func (h *Hoverable) update(gtx layout.Context) {
	for _, e := range gtx.Events(h) {
		ev, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		switch ev.Type {
		case pointer.Enter:
			h.hovered = true
		case pointer.Leave:
			h.hovered = false
		}
	}
}

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

	if active {
		if anim.progress < anim.Duration {
			anim.progress += delta
			op.InvalidateOp{}.Add(gtx.Ops)
			if anim.progress > anim.Duration {
				anim.progress = anim.Duration
			}
		}
	} else {
		if anim.progress > 0 {
			anim.progress -= delta
			op.InvalidateOp{}.Add(gtx.Ops)
			if anim.progress < 0 {
				anim.progress = 0
			}
		}
	}

	return anim.Progress()
}
