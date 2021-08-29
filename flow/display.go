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

	pxPerUnit := float32(gtx.PxPerUnit)
	const (
		tabX  = 0.3
		tabY  = 0.25
		portR = 0.15
	)

	path := func(ops *op.Ops) clip.PathSpec {

		b := layout.FRect(b)

		var p clip.Path
		p.Begin(ops)
		p.MoveTo(b.Min)

		p.LineTo(f32.Point{X: b.Max.X, Y: b.Min.Y})
		for range n.Out {
			p.Line(f32.Pt(pxPerUnit*tabX, pxPerUnit*tabY))
			p.Line(f32.Pt(0, pxPerUnit*(1-tabY*2)))
			p.Line(f32.Pt(-pxPerUnit*tabX, pxPerUnit*tabY))
		}
		p.LineTo(b.Max)
		p.LineTo(f32.Point{X: b.Min.X, Y: b.Max.Y})

		p.LineTo(f32.Point{X: b.Min.X, Y: b.Min.Y + pxPerUnit*float32(len(n.In))})
		for range n.In {
			p.Line(f32.Pt(-pxPerUnit*tabX, -pxPerUnit*tabY))
			p.Line(f32.Pt(0, -pxPerUnit*(1-tabY*2)))
			p.Line(f32.Pt(pxPerUnit*tabX, -pxPerUnit*tabY))
		}
		p.Close()

		return p.End()
	}

	//  background
	func() {
		defer op.Save(gtx.Ops).Load()
		clip.Outline{Path: path(gtx.Ops)}.Op().Add(gtx.Ops)
		paint.ColorOp{Color: gtx.Theme.Node.Fill}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}()

	// border
	func() {
		defer op.Save(gtx.Ops).Load()
		clip.Stroke{
			Path: path(gtx.Ops),
			Style: clip.StrokeStyle{
				Width: gtx.Metric.PxPerDp,
			},
		}.Op().Add(gtx.Ops)
		paint.ColorOp{Color: gtx.Theme.Node.Border}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}()

	// ports
	for _, p := range n.Ports {
		paint.FillShape(gtx.Ops, gtx.Theme.Node.Border, clip.Circle{
			Center: gtx.FPt(p.Position()),
			Radius: pxPerUnit * portR,
		}.Op(gtx.Ops))
	}

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
	connectionWidth := gtx.PxPerUnit / 8

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
