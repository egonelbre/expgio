package main

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type NodeLayer struct{}

func (layer *NodeLayer) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		layer.LayoutNode(gtx, node)
	}
}

func (layer *NodeLayer) LayoutNode(gtx *Context, n *Node) {
	b := gtx.Bounds(n.Box)
	FillRect(gtx, b, gtx.Theme.Node.Fill)
	FillRectBorder(gtx, b, float32(gtx.Dp), gtx.Theme.Node.Border)

	defer op.Save(gtx.Ops).Load()
	clip.Rect(b).Op().Add(gtx.Ops)
	op.Offset(layout.FPt(b.Min)).Add(gtx.Ops)

	before := gtx.Constraints
	defer func() { gtx.Constraints = before }()

	gtx.Constraints.Min = b.Size()
	gtx.Constraints.Max = b.Size()

	n.Display.Layout(gtx)
}

type ConnLayer struct{}

func (layer *ConnLayer) Layout(gtx *Context) {
	for _, conn := range gtx.Diagram.Conns {
		layer.LayoutConn(gtx, conn)
	}
}

func (layer *ConnLayer) LayoutConn(gtx *Context, conn *Conn) {
	connectionWidth := gtx.PxPerUnit / 4

	from := gtx.FPt(conn.From.Position())
	to := gtx.FPt(conn.To.Position()).Sub(from)

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

	defer op.Save(gtx.Ops).Load()
	clip.Stroke{
		Path: path(gtx.Ops),
		Style: clip.StrokeStyle{
			Cap:   clip.RoundCap,
			Width: float32(connectionWidth),
		},
	}.Op().Add(gtx.Ops)
	paint.ColorOp{Color: gtx.Theme.Conn.Border}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
