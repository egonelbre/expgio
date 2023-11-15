package main

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
)

func BenchmarkLayoutStd(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		layout.Stack{}.Layout(gtx,
			layout.Expanded(Example{
				Size: image.Point{X: 10, Y: 10},
			}.Layout),
			layout.Stacked(Example{
				Size: image.Point{X: 50, Y: 50},
			}.Layout),
			layout.Stacked(Example{
				Size: image.Point{X: 30, Y: 60},
			}.Layout),
		)
	}
}

func BenchmarkLayoutStd_ExpandedStacked(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		layout.Stack{}.Layout(gtx,
			layout.Expanded(Example{
				Size: image.Point{X: 10, Y: 10},
			}.Layout),
			layout.Stacked(Example{
				Size: image.Point{X: 50, Y: 50},
			}.Layout),
		)
	}
}

func BenchmarkLayout_ExpandedStacked(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		ExpandedStacked{}.Layout(gtx,
			Example{
				Size: image.Point{X: 10, Y: 10},
			}.Layout,
			Example{
				Size: image.Point{X: 50, Y: 50},
			}.Layout,
		)
	}
}

func BenchmarkLayout(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		Stack{}.Layout(gtx,
			Expanded(Example{
				Size: image.Point{X: 10, Y: 10},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 50, Y: 50},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 30, Y: 60},
			}.Layout),
		)
	}
}

func BenchmarkLayout3(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		Stack{}.Layout3(gtx,
			Expanded(Example{
				Size: image.Point{X: 10, Y: 10},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 50, Y: 50},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 30, Y: 60},
			}.Layout),
		)
	}
}

func BenchmarkShortLayout3(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		Stack3{}.Layout3(gtx,
			Expanded(Example{
				Size: image.Point{X: 10, Y: 10},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 50, Y: 50},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 30, Y: 60},
			}.Layout),
		)
	}
}

func BenchmarkShortLayout4(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		Stack4{}.Layout3(gtx,
			Expanded(Example{
				Size: image.Point{X: 10, Y: 10},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 50, Y: 50},
			}.Layout),
			Stacked(Example{
				Size: image.Point{X: 30, Y: 60},
			}.Layout),
		)
	}
}

func BenchmarkShortLayout5(b *testing.B) {
	gtx := layout.Context{
		Ops: new(op.Ops),
		Constraints: layout.Constraints{
			Max: image.Point{X: 100, Y: 100},
		},
	}

	for i := 0; i < b.N; i++ {
		gtx.Ops.Reset()

		Stack5{
			Children: []StackChild5{
				Expanded5(Example{
					Size: image.Point{X: 10, Y: 10},
				}.Layout),
				Stacked5(Example{
					Size: image.Point{X: 50, Y: 50},
				}.Layout),
				Stacked5(Example{
					Size: image.Point{X: 30, Y: 60},
				}.Layout),
			},
		}.Layout(gtx)
	}
}

type Example struct {
	Size image.Point
}

func (e Example) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: e.Size.X,
			Y: e.Size.Y,
		},
	}
}
