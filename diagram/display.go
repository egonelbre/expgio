package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type NodeHud struct{}

func (hud *NodeHud) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		hud.LayoutNode(gtx, node)
	}
}

func (hud *NodeHud) LayoutNode(gtx *Context, n *Node) {
	b := gtx.Bounds(n.Box)
	FillRect(gtx, b, n.Style.Fill)
	FillRectBorder(gtx, b, float32(gtx.Transform.Dp), n.Style.Border)
}

type ConnectionHud struct{}

func (hud *ConnectionHud) Layout(gtx *Context) {
	for _, conn := range gtx.Diagram.Connections {
		hud.LayoutConnection(gtx, conn)
	}
}

func (hud *ConnectionHud) LayoutConnection(gtx *Context, c *Connection) {
	connectionWidth := gtx.PxPerUnit / 4

	from := gtx.FPt(c.From.Position())
	to := gtx.FPt(c.To.Position()).Sub(from)

	curveOffset := f32.Point{X: float32(to.X) / 2}

	path := func(ops *op.Ops) clip.PathSpec {
		var p clip.Path
		p.Begin(ops)
		p.MoveTo(from)
		if to.X != 0 {
			p.Cube(curveOffset, to.Sub(curveOffset), to)
		} else {
			p.Line(to)
		}
		return p.End()
	}

	paint.FillShape(gtx.Ops, DefaultConnection.Border, clip.Stroke{
		Path:  path(gtx.Ops),
		Width: float32(connectionWidth + gtx.Transform.Dp*2),
	}.Op())

	paint.FillShape(gtx.Ops, DefaultConnection.Fill, clip.Stroke{
		Path:  path(gtx.Ops),
		Width: float32(connectionWidth),
	}.Op())
}

type PortHud struct {
	ShowAll *bool
}

func (hud *PortHud) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		if !(hud.ShowAll != nil && *hud.ShowAll) && !gtx.Diagram.Selection.Contains(node) {
			continue
		}
		for _, port := range node.Ports {
			hud.LayoutPort(gtx, port)
		}
	}
}

func (hud *PortHud) LayoutPort(gtx *Context, p *Port) {
	pos := gtx.Pt(p.Position())
	r := image.Rectangle{Min: pos, Max: pos}
	b := r.Inset(-gtx.PxPerUnit / 4)
	FillRect(gtx, b, DefaultPort.Fill)
	FillRectBorder(gtx, b, float32(gtx.Transform.Dp), DefaultPort.Border)
}
