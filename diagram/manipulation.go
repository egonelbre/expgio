package main

import (
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ManipulationHud struct{}

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
		Types: pointer.Press,
	}.Add(gtx.Ops)

	for _, ev := range gtx.Events(tag) {
		if x, ok := ev.(pointer.Event); ok {
			switch x.Type {
			case pointer.Press:
				if x.Modifiers.Contain(key.ModCtrl) {
					gtx.Diagram.Selection.Toggle(node)
				} else {
					gtx.Diagram.Selection.Select(node)
				}
			}
		}
	}
}

func (hud *ManipulationHud) LayoutNode(gtx *Context, node *Node) {
	defer op.Save(gtx.Ops).Load()

	b := gtx.Bounds(node.Box)

	clip.Border{
		Rect:  layout.FRect(b),
		Width: 4,
	}.Add(gtx.Ops)
	// TODO: use dashed border

	paint.ColorOp{Color: FocusColor.Fill}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
