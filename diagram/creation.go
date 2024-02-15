package main

import (
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type NodeCreationHud struct {
	drag gesture.Drag

	start   Vector
	end     Vector
	pointer pointer.ID
	drawing bool
}

func (hud *NodeCreationHud) Layout(gtx *Context) {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

	event.Op(gtx.Ops, hud)

	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: hud,
			Kinds:  pointer.Press | pointer.Drag | pointer.Release,
		})
		if !ok {
			break
		}

		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Kind {
			case pointer.Press:
				if hud.pointer == 0 {
					hud.start = gtx.FInv(ev.Position)
					hud.end = hud.start

					gtx.Diagram.Selection.Clear()
				}
			case pointer.Drag:
				if ev.PointerID == hud.pointer {
					hud.drawing = true
					hud.end = gtx.FInv(ev.Position).Add(Vector{X: 1, Y: 1})

					if ev.Priority < pointer.Grabbed {
						gtx.Execute(pointer.GrabCmd{
							Tag: hud,
							ID:  hud.pointer,
						})
					}
				}
			case pointer.Release:
				if hud.drawing && ev.PointerID == hud.pointer {
					hud.drawing = false
					hud.pointer = 0

					min := hud.start.Min(hud.end)
					max := hud.start.Max(hud.end)
					size := max.Sub(min)
					if size.X > 0 && size.Y > 0 {
						node := NewNode(min, size)
						gtx.Diagram.Nodes = append(gtx.Diagram.Nodes, node)
						gtx.Diagram.Selection.Select(node)
						gtx.Execute(op.InvalidateCmd{})
					}
				}
			case pointer.Cancel:
				if ev.PointerID == hud.pointer {
					hud.drawing = false
					hud.pointer = 0
				}
			}
		}
	}

	if hud.drawing {
		min := gtx.Pt(hud.start.Min(hud.end))
		max := gtx.Pt(hud.start.Max(hud.end))
		FillRect(gtx, image.Rectangle{
			Min: min,
			Max: max,
		}, WithAlpha(ActiveColor.Fill, 0xEE))
	}
}

type ConnectionCreationHud struct {
	source    *Port
	end       image.Point
	target    *Port
	newtarget *Port

	pointer pointer.ID
	drawing bool
}

func (hud *ConnectionCreationHud) Layout(gtx *Context) {
	hud.newtarget = nil
	for _, node := range gtx.Diagram.Nodes {
		for _, port := range node.Ports {
			hud.LayoutPort(gtx, port)
		}
	}
	hud.target = hud.newtarget

	if hud.drawing {
		from := gtx.Pt(hud.source.Position())
		var to image.Point
		var col color.NRGBA

		if hud.target == nil {
			to = hud.end
			col = WithAlpha(ActiveColor.Fill, 0xEE)
		} else {
			to = gtx.Pt(hud.target.Position())
			col = WithAlpha(DefaultConnection.Fill, 0xEE)
		}

		FillLine(gtx, from, to, gtx.PxPerUnit/4, col)
	}
}

type connectionCreationTag *Port

func (hud *ConnectionCreationHud) LayoutPort(gtx *Context, p *Port) {
	pos := gtx.Pt(p.Position())
	r := image.Rectangle{Min: pos, Max: pos}
	r = r.Inset(-gtx.PxPerUnit / 2)

	if hud.drawing && p != hud.source {
		if hud.end.In(r) {
			hud.newtarget = p
		}
	}

	defer clip.Rect(r).Push(gtx.Ops).Pop()

	tag := connectionCreationTag(p)
	event.Op(gtx.Ops, tag)
	pointer.CursorCrosshair.Add(gtx.Ops)

	for {
		ev, ok := gtx.Event(
			pointer.Filter{
				Target: tag,
				Kinds:  pointer.Press | pointer.Drag | pointer.Release,
			},
		)
		if !ok {
			break
		}
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Kind {
			case pointer.Press:
				if hud.pointer == 0 {
					hud.source = p
					hud.end = image.Point{
						X: int(ev.Position.X),
						Y: int(ev.Position.Y),
					}
					hud.pointer = ev.PointerID
					hud.drawing = true
				}
			case pointer.Drag:
				if ev.PointerID == hud.pointer {
					hud.drawing = true
					hud.end = image.Point{
						X: int(ev.Position.X),
						Y: int(ev.Position.Y),
					}

					if ev.Priority < pointer.Grabbed {
						gtx.Execute(pointer.GrabCmd{
							Tag: tag,
							ID:  hud.pointer,
						})
					}
				}
			case pointer.Release:
				if hud.drawing && ev.PointerID == hud.pointer {
					hud.drawing = false
					hud.pointer = 0

					if hud.target != nil {
						gtx.Diagram.Connections = append(gtx.Diagram.Connections,
							&Connection{From: hud.source, To: hud.target},
						)
					}
				}
			case pointer.Cancel:
				if ev.PointerID == hud.pointer {
					hud.drawing = false
					hud.pointer = 0
				}
			}
		}
	}
}
