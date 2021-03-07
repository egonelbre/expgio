package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type NodeDisplay struct{}

func (d *NodeDisplay) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		d.LayoutNode(gtx, node)
	}
}

func (d *NodeDisplay) LayoutNode(gtx *Context, n *Node) {
	b := gtx.Bounds(n.Box)
	FillRect(gtx, b, color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF})
}

type ConnectionDisplay struct{}

func (d *ConnectionDisplay) Layout(gtx *Context) {
	for _, conn := range gtx.Diagram.Connections {
		d.LayoutConnection(gtx, conn)
	}
}

func (d *ConnectionDisplay) LayoutConnection(gtx *Context, c *Connection) {
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

type PortDisplay struct{}

func (d *PortDisplay) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		for _, port := range node.Ports {
			d.LayoutPort(gtx, port)
		}
	}
}

func (d *PortDisplay) LayoutPort(gtx *Context, p *Port) {
	pos := gtx.Pt(p.Position())
	r := image.Rectangle{Min: pos, Max: pos}
	r = r.Inset(-gtx.PxPerUnit / 4)
	FillRect(gtx, r, color.NRGBA{R: 0xA0, B: 0xA0, A: 0xFF})
}
