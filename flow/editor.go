package main

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type Editor struct {
	Zoom    Zoom
	Diagram *Diagram

	Layers    []Layer
	Exclusive Layer
}

type Layer interface {
	Layout(*Context)
}

func NewEditor(diagram *Diagram) *Editor {
	editor := &Editor{
		Diagram: diagram,
	}

	editor.Zoom.Level = defaultZoom

	editor.AddLayer(&BackgroundLayer{})
	editor.AddLayer(&GridLayer{})
	editor.AddLayer(&ConnLayer{})
	editor.AddLayer(&NodeLayer{})

	return editor
}

func (editor *Editor) AddLayer(layer Layer) {
	editor.Layers = append(editor.Layers, layer)
}

func (editor *Editor) Layout(th *Theme, gtx layout.Context) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	clip.Rect{Max: gtx.Constraints.Max}.Add(gtx.Ops)

	for _, layer := range editor.Layers {
		var lgtx layout.Context
		if editor.Exclusive == nil || editor.Exclusive == layer {
			lgtx = gtx
		} else {
			lgtx = gtx.Disabled()
		}
		dgtx := NewContext(lgtx, th, &editor.Zoom, editor.Diagram)
		layer.Layout(dgtx)
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}
