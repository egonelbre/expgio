package main

type Diagram struct {
	Selection   Set
	Nodes       []*Node
	Connections []*Connection
}

func NewDemoDiagram() *Diagram {
	diagram := &Diagram{}
	diagram.Nodes = []*Node{
		NewNode(V(1, 1), V(6, 3)),
		NewNode(V(1, 10), V(6, 3)),
		NewNode(V(10, 1), V(6, 3)),
		NewNode(V(10, 10), V(6, 3)),
	}
	ns := diagram.Nodes
	diagram.Connections = []*Connection{
		{From: ns[0].Ports[3], To: ns[2].Ports[4]},
		{From: ns[0].Ports[5], To: ns[3].Ports[4]},
		{From: ns[1].Ports[1], To: ns[3].Ports[6]},
	}
	return diagram
}

type Node struct {
	Box
	Ports []*Port
}

func NewNode(pos, size Vector) *Node {
	node := &Node{
		Box: Box{
			Pos:  pos,
			Size: size,
		},
	}
	for y := Unit(0); y <= size.Y; y++ {
		node.Ports = append(node.Ports,
			&Port{
				Owner:  node,
				Offset: V(0, y),
			},
			&Port{
				Owner:  node,
				Offset: V(size.X, y),
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
	return p.Owner.Pos.Add(p.Offset)
}
