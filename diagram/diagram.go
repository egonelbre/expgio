package main

type Diagram struct {
	Selection   *Selection
	Nodes       []*Node
	Connections []*Connection
}

func NewDiagram() *Diagram {
	diagram := &Diagram{}
	diagram.Selection = NewSelection()
	return diagram
}

func NewDemoDiagram() *Diagram {
	diagram := NewDiagram()
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
	Style *Style
	Ports []*Port
}

func NewNode(pos, size Vector) *Node {
	node := &Node{
		Box: Box{
			Pos:  pos,
			Size: size,
		},
		Style: Default,
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

type Selection struct {
	Selected Set
}

func NewSelection() *Selection {
	return &Selection{
		Selected: make(Set),
	}
}

func (sel *Selection) Clear() {
	sel.Selected = Set{}
}

func (sel *Selection) Contains(v any) bool {
	_, ok := sel.Selected[v]
	return ok
}

func (sel *Selection) Toggle(v any) {
	sel.Selected.Toggle(v)
}

func (sel *Selection) Set(v any) {
	sel.Clear()
	sel.Selected.Include(v)
}

func (sel *Selection) Select(v any) {
	if !sel.Selected.Contains(v) {
		sel.Set(v)
	}
}
