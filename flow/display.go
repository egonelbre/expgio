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

	const p = 0.6
	const c = 0.2

	tabR2 := pxPerUnit * p
	tabP := pxPerUnit * (1.0 - p) / 2.0
	portR := pxPerUnit * 0.15

	path := func(ops *op.Ops) clip.PathSpec {
		min := layout.FPt(b.Min)
		max := layout.FPt(b.Max)

		var p clip.Path
		p.Begin(ops)
		p.MoveTo(min)

		p.LineTo(f32.Point{X: max.X, Y: min.Y})
		for range n.Out {
			p.Line(f32.Pt(0, tabP))
			p.Cube(
				f32.Pt(1.4*tabR2/2, 0),
				f32.Pt(1.4*tabR2/2, tabR2),
				f32.Pt(0, tabR2),
			)
			p.Line(f32.Pt(0, tabP))
		}
		p.LineTo(max)
		p.LineTo(f32.Point{X: min.X, Y: max.Y})

		p.LineTo(f32.Point{X: min.X, Y: min.Y + pxPerUnit*float32(len(n.In))})
		for range n.In {
			p.Line(f32.Pt(0, -tabP))
			p.Cube(
				f32.Pt(-1.4*tabR2/2, 0),
				f32.Pt(-1.4*tabR2/2, -tabR2),
				f32.Pt(0, -tabR2),
			)
			p.Line(f32.Pt(0, -tabP))
		}
		p.Close()

		return p.End()
	}

	pixelAlignLine := f32.Pt(-gtx.Metric.PxPerDp/2, -gtx.Metric.PxPerDp/2)
	pixelAlign := f32.Affine2D{}.Offset(pixelAlignLine)

	//  background
	func() {
		defer op.Affine(pixelAlign).Push(gtx.Ops).Pop()
		paint.FillShape(gtx.Ops, gtx.Theme.Node.Fill, clip.Outline{
			Path: path(gtx.Ops),
		}.Op())
	}()

	// border
	func() {
		defer op.Affine(pixelAlign).Push(gtx.Ops).Pop()
		paint.FillShape(gtx.Ops, gtx.Theme.Node.Border, clip.Stroke{
			Path:  path(gtx.Ops),
			Width: gtx.Metric.PxPerDp,
		}.Op())
	}()

	// ports
	for _, p := range n.Ports {
		center := gtx.FPt(p.Position())
		center = center.Add(pixelAlignLine)
		paint.FillShape(gtx.Ops, gtx.Theme.Node.Border, Circle{
			Center: center,
			Radius: portR,
		}.Op(gtx.Ops))
	}

	defer clip.Rect(b).Op().Push(gtx.Ops).Pop()
	defer op.Offset(b.Min).Push(gtx.Ops).Pop()

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

	paint.FillShape(gtx.Ops, gtx.Theme.Conn.Border, clip.Stroke{
		Path:  path(gtx.Ops),
		Width: float32(connectionWidth),
	}.Op())
}
