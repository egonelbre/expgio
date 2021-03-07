package main

type Selecter struct{}

func (d *Selecter) Enabled() bool { return true }

func (d *Selecter) Layout(gtx *Context) {
	for _, node := range gtx.Diagram.Nodes {
		_ = node
	}
}
