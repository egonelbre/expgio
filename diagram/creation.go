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
				}
			case pointer.Drag:
				if ev.PointerID == hud.pointer {
					hud.drawing = true
					hud.end = gtx.FInv(ev.Position).Add(Vector{X: 1, Y: 1})
				}
			case pointer.Release:
				if hud.drawing && ev.PointerID == hud.pointer {
					hud.drawing = false

					min := hud.start.Min(hud.end)
					max := hud.start.Max(hud.end)
					size := max.Sub(min)
					if size.X > 0 && size.Y > 0 {
						gtx.Diagram.Nodes = append(gtx.Diagram.Nodes, NewNode(min, size))
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
		}, color.NRGBA{G: 0x88, A: 0x88})
	}
}
