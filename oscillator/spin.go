package main

import (
	"slices"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Spinnable[T comparable] interface {
	comparable
	Options() []T
	String() string
}

type Spin[T Spinnable[T]] struct {
	Current *T // we use a pointer here to change values on the pending Config

	Next widget.Clickable
	Prev widget.Clickable
}

func (spin *Spin[T]) Spin(offset int) {
	values := (*spin.Current).Options()
	i := slices.Index(values, *spin.Current)
	if i < 0 { // didn't find the current value in the options list
		*spin.Current = values[0]
		return
	}

	// implement wrap around
	i = i + offset
	if i < 0 {
		i += len(values)
	}
	i = i % len(values)

	// update the value
	*spin.Current = values[i]
}

func (spin *Spin[T]) update(gtx layout.Context) {
	if spin.Prev.Clicked(gtx) {
		spin.Spin(-1)
	}
	if spin.Next.Clicked(gtx) {
		spin.Spin(1)
	}
}

func (spin *Spin[T]) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	spin.update(gtx)
	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceBetween,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(material.Button(th, &spin.Next, "<").Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx,
				material.Body1(th, (*spin.Current).String()).Layout)
		}),
		layout.Rigid(material.Button(th, &spin.Prev, ">").Layout),
	)
}
