package main

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type NodeDisplay struct{}

func (d *NodeDisplay) Enabled() bool { return true }

func (d *NodeDisplay) Layout(diagram *Diagram, gtx layout.Context) {
	for _, node := range diagram.Nodes {
		d.LayoutNode(node, gtx)
	}
}

func (d *NodeDisplay) LayoutNode(n *Node, gtx layout.Context) {
	defer op.Save(gtx.Ops).Load()

	pos := VectorPx(n.Position, gtx)
	size := VectorPx(n.Size, gtx)

	paint.ColorOp{Color: color.NRGBA{R: 0xA0, G: 0xA0, A: 0xFF}}.Add(gtx.Ops)
	clip.Rect{Min: pos, Max: pos.Add(size)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

type ConnectionDisplay struct{}

func (d *ConnectionDisplay) Enabled() bool { return true }

func (d *ConnectionDisplay) Layout(diagram *Diagram, gtx layout.Context) {
	for _, conn := range diagram.Connections {
		d.LayoutConnection(conn, gtx)
	}
}

func (d *ConnectionDisplay) LayoutConnection(c *Connection, gtx layout.Context) {
	defer op.Save(gtx.Ops).Load()

	cell := VectorPx(image.Pt(1, 1), gtx)

	from := VectorPx(c.From.Position(), gtx)
	to := VectorPx(c.To.Position(), gtx)

	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(layout.FPt(from))
	p.LineTo(layout.FPt(to))

	clip.Stroke{
		Path: p.End(),
		Style: clip.StrokeStyle{
			Width: float32(cell.X) / 4,
		},
	}.Op().Add(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{G: 0xA0, B: 0xA0, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

type PortDisplay struct{}

func (d *PortDisplay) Enabled() bool { return true }

func (d *PortDisplay) Layout(diagram *Diagram, gtx layout.Context) {
	for _, node := range diagram.Nodes {
		for _, port := range node.Ports {
			d.LayoutPort(port, gtx)
		}
	}
}

func (d *PortDisplay) LayoutPort(p *Port, gtx layout.Context) {
	defer op.Save(gtx.Ops).Load()

	pos := VectorPx(p.Position(), gtx)
	size := VectorPx(image.Point{X: 1, Y: 1}, gtx)

	var r clip.Rect
	r.Min = pos
	r.Min.X -= size.X / 4
	r.Min.Y -= size.Y / 4
	r.Max = pos
	r.Max.X += size.X / 4
	r.Max.Y += size.Y / 4
	r.Add(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{R: 0xA0, B: 0xA0, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
