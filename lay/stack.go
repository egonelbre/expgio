package lay

import "gioui.org/layout"

type Padding struct {
	N, E, S, W Scale
}

func (p Padding) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Dimensions{}
}

type Scroll struct {
	Position int
}

type Stack struct {
	Gap Scale
}

func (stack Stack) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Dimensions{}
}

/*

Box{
	FlexH{
		Icon{},
		FlexV{
			Gap: Small,
			FlexH{
				Gap: Small,
				Text{Big, "Messages"},
				Color{Red, Text{Small, "| Global"}},
			},
			Text{Default, "lorem"},
		},
		DragIcon{},
	}
}

*/
