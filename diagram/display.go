package main

import (
	"image"
	"image/color"

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
	FillRect(gtx, b, color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF})
}

type ConnectionHud struct{}

func (hud *ConnectionHud) Layout(gtx *Context) {
	for _, conn := range gtx.Diagram.Connections {
		hud.LayoutConnection(gtx, conn)
	}
}

func (hud *ConnectionHud) LayoutConnection(gtx *Context, c *Connection) {
	defer op.Save(gtx.Ops).Load()

	connectionWidth := gtx.PxPerUnit / 4

	from := gtx.FPt(c.From.Position())
	to := gtx.FPt(c.To.Position()).Sub(from)

	curveOffset := f32.Point{X: float32(gtx.PxPerUnit)}

	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(from)
	p.Cube(curveOffset, to.Sub(curveOffset), to)
	clip.Stroke{
		Path: p.End(),
		Style: clip.StrokeStyle{
			Width: float32(connectionWidth),
		},
	}.Op().Add(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{G: 0xA0, B: 0xA0, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

type PortHud struct{}

func (hud *PortHud) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		for _, port := range node.Ports {
			hud.LayoutPort(gtx, port)
		}
	}
}

func (hud *PortHud) LayoutPort(gtx *Context, p *Port) {
	pos := gtx.Pt(p.Position())
	r := image.Rectangle{Min: pos, Max: pos}
	r = r.Inset(-gtx.PxPerUnit / 4)
	FillRect(gtx, r, color.NRGBA{R: 0xA0, B: 0xA0, A: 0xFF})
}
