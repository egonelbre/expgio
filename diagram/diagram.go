package main

import (
	"image"
)

type Diagram struct {
	Selection   Set
	Nodes       []*Node
	Connections []*Connection
}

func NewDemoDiagram() *Diagram {
	diagram := &Diagram{}
	diagram.Nodes = []*Node{
		NewNode(image.Pt(1, 1), image.Pt(6, 3)),
		NewNode(image.Pt(1, 10), image.Pt(6, 3)),
		NewNode(image.Pt(10, 1), image.Pt(6, 3)),
		NewNode(image.Pt(10, 10), image.Pt(6, 3)),
	}
	ns := diagram.Nodes
	diagram.Connections = []*Connection{
		{From: ns[0].Ports[3], To: ns[2].Ports[2]},
		{From: ns[0].Ports[5], To: ns[3].Ports[4]},
		{From: ns[1].Ports[1], To: ns[3].Ports[6]},
	}
	return diagram
}

type Node struct {
	Position Vector
	Size     Vector

	Ports []*Port
}

func NewNode(pos, size image.Point) *Node {
	node := &Node{
		Position: pos,
		Size:     size,
	}
	for y := 0; y <= size.Y; y++ {
		node.Ports = append(node.Ports,
			&Port{
				Owner:  node,
				Offset: image.Point{X: 0, Y: y},
			},
			&Port{
				Owner:  node,
				Offset: image.Point{X: size.X, Y: y},
			},
		)
	}
	return node
}

type Port struct {
	Owner  *Node
	Offset Vector
}

type Connection struct {
	From *Port
	To   *Port
}

func (p *Port) Position() Vector {
	return p.Offset.Add(p.Owner.Position)
}
