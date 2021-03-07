package main

import (
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/op"
)

type NodeCreationHud struct {
	drag gesture.Drag

	start   Vector
	end     Vector
	pointer pointer.ID
	drawing bool
}

func (hud *NodeCreationHud) Layout(gtx *Context) {
	defer op.Save(gtx.Ops).Load()

	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   hud,
		Grab:  hud.drawing,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}.Add(gtx.Ops)

	for _, ev := range gtx.Events(hud) {
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Type {
			case pointer.Press:
				if hud.pointer == 0 {
					hud.start = gtx.FInv(ev.Position)
					hud.pointer = ev.PointerID

					gtx.Diagram.Selection.Clear()
				}
			case pointer.Drag:
				if ev.PointerID == hud.pointer {
					hud.drawing = true
					hud.end = gtx.FInv(ev.Position).Add(Vector{X: 1, Y: 1})
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
	defer op.Save(gtx.Ops).Load()
	pos := gtx.Pt(p.Position())
	r := image.Rectangle{Min: pos, Max: pos}
	r = r.Inset(-gtx.PxPerUnit / 2)

	if hud.drawing && p != hud.source {
		if hud.end.In(r) {
			hud.newtarget = p
		}
	}

	tag := connectionCreationTag(p)

	pointer.Rect(r).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   tag,
		Grab:  hud.drawing,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}.Add(gtx.Ops)

	for _, ev := range gtx.Events(tag) {
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Type {
			case pointer.Press:
				if hud.pointer == 0 {
					hud.source = p
					hud.pointer = ev.PointerID
				}
			case pointer.Drag:
				if ev.PointerID == hud.pointer {
					hud.drawing = true
					hud.end = image.Point{
						X: int(ev.Position.X),
						Y: int(ev.Position.Y),
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
