package main

type Selection struct{}

func (d *Selection) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		_ = node
	}
}
