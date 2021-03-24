package main

import (
	"image"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/op"
)

type ManipulationHud struct {
	start        image.Point
	current      image.Point
	appliedDelta Vector

	pointer  pointer.ID
	dragging bool
}

func (hud *ManipulationHud) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		hud.HandleNode(gtx, node)
	}

	for selected := range gtx.Diagram.Selection.Selected {
		switch sel := selected.(type) {
		case *Node:
			hud.LayoutNode(gtx, sel)
		}
	}
}

type manipulationTag *Node

func (hud *ManipulationHud) HandleNode(gtx *Context, node *Node) {
	defer op.Save(gtx.Ops).Load()
	tag := manipulationTag(node)

	pointer.Rect(gtx.Bounds(node.Box)).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   tag,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}.Add(gtx.Ops)

	for _, ev := range gtx.Events(tag) {
		if ev, ok := ev.(pointer.Event); ok {
			switch ev.Type {
			case pointer.Press:
				if !hud.dragging {
					if ev.Modifiers.Contain(key.ModCtrl) {
						gtx.Diagram.Selection.Toggle(node)
					} else {
						gtx.Diagram.Selection.Select(node)
					}
				}
				if !hud.dragging {
					hud.pointer = ev.PointerID
					hud.start = image.Point{
						X: int(ev.Position.X),
						Y: int(ev.Position.Y),
					}
					hud.current = hud.start
				}
			case pointer.Drag:
				if ev.PointerID == hud.pointer {
					hud.dragging = true
					hud.current = image.Point{
						X: int(ev.Position.X),
						Y: int(ev.Position.Y),
					}
					hud.updateDelta(gtx)
				}
			case pointer.Release:
				if ev.PointerID == hud.pointer && hud.dragging {
					hud.applyDelta(gtx)
					hud.dragging = false
					hud.pointer = 0
					hud.start = image.Point{}
					hud.current = image.Point{}

				}
			case pointer.Cancel:
				if ev.PointerID == hud.pointer && hud.dragging {
					hud.resetDelta(gtx)
					hud.dragging = false
					hud.pointer = 0
					hud.start = image.Point{}
					hud.current = image.Point{}
				}
			}
		}
	}
}

func (hud *ManipulationHud) updateDelta(gtx *Context) {
	lastDelta := hud.appliedDelta
	newDelta := gtx.Inv(hud.current.Sub(hud.start))
	if lastDelta == newDelta {
		return
	}
	hud.appliedDelta = newDelta

	for selected := range gtx.Diagram.Selection.Selected {
		switch sel := selected.(type) {
		case *Node:
			sel.Pos = sel.Pos.Sub(lastDelta).Add(newDelta)
		}
	}
}

func (hud *ManipulationHud) applyDelta(gtx *Context) {
	hud.appliedDelta = Vector{}
}

func (hud *ManipulationHud) resetDelta(gtx *Context) {
	hud.current = hud.start
	hud.updateDelta(gtx)
}

func (hud *ManipulationHud) LayoutNode(gtx *Context, node *Node) {
	b := gtx.Bounds(node.Box)
	FillRectBorder(gtx, b, 4, FocusColor.Fill)
}
