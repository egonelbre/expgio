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
	FillRect(gtx, b, Default.Fill)
	FillRectBorder(gtx, b, 1, Default.Border)
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

	func() {
		defer op.Save(gtx.Ops).Load()

		var p clip.Path
		p.Begin(gtx.Ops)
		p.MoveTo(from)
		if to.X-from.X != 0 {
			p.Cube(curveOffset, to.Sub(curveOffset), to)
		} else {
			p.Line(to)
		}
		pathOp := p.End()

		clip.Stroke{
			Path: pathOp,
			Style: clip.StrokeStyle{
				Cap:   clip.RoundCap,
				Width: float32(connectionWidth + 2),
			},
		}.Op().Add(gtx.Ops)
		paint.ColorOp{Color: DefaultConnection.Border}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}()

	func() {
		defer op.Save(gtx.Ops).Load()

		var p clip.Path
		p.Begin(gtx.Ops)
		p.MoveTo(from)
		p.Cube(curveOffset, to.Sub(curveOffset), to)
		pathOp := p.End()

		clip.Stroke{
			Path: pathOp,
			Style: clip.StrokeStyle{
				Cap:   clip.RoundCap,
				Width: float32(connectionWidth),
			},
		}.Op().Add(gtx.Ops)
		paint.ColorOp{Color: DefaultConnection.Fill}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}()
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
	FillRectBorder(gtx, b, 1, DefaultPort.Border)
}
