package scope

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/egonelbre/expgio/oscillator/generator"
)

type Display struct {
	Theme *material.Theme
	Data  generator.Data
}

func (display *Display) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}
