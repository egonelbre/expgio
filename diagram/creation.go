package main

import (
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/op"
)

type NodeCreationDisplay struct {
	drag gesture.Drag

	start   Vector
	end     Vector
	pointer pointer.ID
	drawing bool
}

func (d *NodeCreationDisplay) Layout(gtx *Context) {
	defer op.Save(gtx.Ops).Load()

	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   d,
		Grab:  d.drawing,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}.Add(gtx.Ops)

	for _, ev := range gtx.Events(d) {
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Type {
			case pointer.Press:
				if d.pointer == 0 {
					d.start = gtx.FInv(ev.Position)
					d.pointer = ev.PointerID
				}
			case pointer.Drag:
				if ev.PointerID == d.pointer {
					d.drawing = true
					d.end = gtx.FInv(ev.Position).Add(Vector{X: 1, Y: 1})
				}
			case pointer.Release:
				if d.drawing && ev.PointerID == d.pointer {
					d.drawing = false

					min := d.start.Min(d.end)
					max := d.start.Max(d.end)
					size := max.Sub(min)
					if size.X > 0 && size.Y > 0 {
						gtx.Diagram.Nodes = append(gtx.Diagram.Nodes, NewNode(min, size))
					}
				}
			case pointer.Cancel:
				if ev.PointerID == d.pointer {
					d.drawing = false
					d.pointer = 0
				}
			}
		}
	}

	if d.drawing {
		min := gtx.Pt(d.start.Min(d.end))
		max := gtx.Pt(d.start.Max(d.end))
		FillRect(gtx, image.Rectangle{
			Min: min,
			Max: max,
		}, color.NRGBA{G: 0x88, A: 0x88})
	}
}
