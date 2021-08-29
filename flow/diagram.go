package main

import (
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type Diagram struct {
	Nodes []*Node
	Conns []*Conn

	Selection Set
	Focus     Set
}

func NewDiagram() *Diagram {
	diagram := &Diagram{}
	diagram.Selection = NewSet()
	diagram.Focus = NewSet()
	return diagram
}

type Node struct {
	Box

	Display Display
	In      []*Port
	Out     []*Port
	Ports   []*Port
}

type Display interface {
	Layout(gtx *Context)
}

type Port struct {
	Name   string
	Owner  *Node
	Offset Vector
}

func (p *Port) Position() Vector {
	return p.Owner.Pos.Add(p.Offset)
}

type Conn struct {
	From *Port
	To   *Port
}

func (diagram *Diagram) NewNode(display Display, pos, size Vector, in []*Port, out []*Port) *Node {
	node := &Node{
		Box: Box{
			Pos:  pos,
			Size: size,
		},
		Display: display,
	}

	for y, p := range in {
		p.Owner = node
		p.Offset = V(0, Unit(y)+0.5)
		node.In = append(node.In, p)
		node.Ports = append(node.Ports, p)
	}

	for y, p := range out {
		p.Owner = node
		p.Offset = V(size.X, Unit(y)+0.5)
		node.Out = append(node.Out, p)
		node.Ports = append(node.Ports, p)
	}

	diagram.Nodes = append(diagram.Nodes, node)

	return node
}

type Label string

func (label Label) Layout(gtx *Context) {
	w := material.Body1(gtx.Theme.Theme, string(label))
	w.Alignment = text.Middle
	w.Layout(gtx.Context)
}

type List []string

func (list List) Layout(gtx *Context) {
	defer op.Save(gtx.Ops).Load()

	for _, label := range list {
		w := material.Body1(gtx.Theme.Theme, string(label))
		w.Alignment = text.Middle
		w.Layout(gtx.Context)

		op.Offset(f32.Point{Y: float32(gtx.PxPerUnit)}).Add(gtx.Ops)
	}
}
