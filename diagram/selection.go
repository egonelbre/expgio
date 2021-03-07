package main

import (
	"image/color"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type SelectionDisplay struct{}

func (d *SelectionDisplay) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		d.HandleNode(gtx, node)
	}

	for selected := range gtx.Diagram.Selection.Selected {
		switch sel := selected.(type) {
		case *Node:
			d.LayoutNode(gtx, sel)
		}
	}
}

type selectionClickTag *Node

func (d *SelectionDisplay) HandleNode(gtx *Context, node *Node) {
	defer op.Save(gtx.Ops).Load()
	tag := selectionClickTag(node)

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
					gtx.Diagram.Selection.Set(node)
				}
			}
		}
	}
}

func (d *SelectionDisplay) LayoutNode(gtx *Context, node *Node) {
	defer op.Save(gtx.Ops).Load()

	b := gtx.Bounds(node.Box)

	clip.Border{
		Rect:  layout.FRect(b),
		Width: 4,
	}.Add(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{B: 0xFF, A: 0xA0}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
