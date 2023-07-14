package rig

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
)

// Sizer represents a resizer.
type Sizer struct {
	entered  bool
	dragging bool
	start    f32.Point
	grab     bool
	pid      pointer.ID
}

// Add the handler to the operation list to receive drag events.
func (s *Sizer) Add(ops *op.Ops) {
	pointer.InputOp{
		Tag:  s,
		Grab: s.grab,
		Types: pointer.Press | pointer.Drag | pointer.Release |
			pointer.Enter | pointer.Leave,
	}.Add(ops)
}

func (s *Sizer) Hovered() bool  { return s.entered }
func (s *Sizer) Dragging() bool { return s.dragging }

const touchSlop = unit.Dp(3)

// Events returns the next drag events, if any.
func (s *Sizer) Events(cfg unit.Metric, q event.Queue, axis gesture.Axis) []pointer.Event {
	var events []pointer.Event
	for _, e := range q.Events(s) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Type {
		case pointer.Press:
			if !(e.Buttons == pointer.ButtonPrimary || e.Source == pointer.Touch) {
				continue
			}
			if s.dragging {
				continue
			}
			s.dragging = true
			s.pid = e.PointerID
			s.start = e.Position
		case pointer.Drag:
			if !s.dragging || e.PointerID != s.pid {
				continue
			}
			switch axis {
			case gesture.Horizontal:
				e.Position.Y = s.start.Y
			case gesture.Vertical:
				e.Position.X = s.start.X
			case gesture.Both:
				// Do nothing
			}
			if e.Priority < pointer.Grabbed {
				diff := e.Position.Sub(s.start)
				slop := cfg.Dp(touchSlop)
				if diff.X*diff.X+diff.Y*diff.Y > float32(slop*slop) {
					s.grab = true
				}
			}

		case pointer.Enter:
			if !s.entered {
				s.pid = e.PointerID
			}
			if s.pid == e.PointerID {
				s.entered = true
			}
		case pointer.Leave:
			if s.entered && s.pid == e.PointerID {
				s.entered = false
			}

		case pointer.Release, pointer.Cancel:
			if !s.dragging || e.PointerID != s.pid {
				continue
			}
			s.dragging = false
			s.grab = false

			if e.Type == pointer.Cancel && s.entered && s.pid == e.PointerID {
				s.entered = false
			}

		}

		events = append(events, e)
	}

	return events
}
