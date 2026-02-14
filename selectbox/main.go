package main

import (
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Selectbox"))
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.FrameEvent:

				gtx := app.NewContext(&ops, e)
				Layout(gtx)
				e.Frame(gtx.Ops)

			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()

	app.Main()
}

type FocusBorderStyle struct {
	Focused     bool
	BorderWidth unit.Dp
	Color       color.NRGBA
}

func FocusBorder(th *material.Theme, focused bool) FocusBorderStyle {
	return FocusBorderStyle{
		Focused:     focused,
		BorderWidth: unit.Dp(2),
		Color:       th.ContrastBg,
	}
}

func (focus FocusBorderStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	inset := layout.UniformInset(focus.BorderWidth)
	if !focus.Focused {
		return inset.Layout(gtx, w)
	}

	return widget.Border{
		Color: focus.Color,
		Width: focus.BorderWidth,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return inset.Layout(gtx, w)
	})
}

type SelectList struct {
	widget.List

	Selected int
	Hovered  int

	ItemHeight unit.Dp

	focused bool
}

func (list *SelectList) Layout(th *material.Theme, gtx layout.Context, length int, element layout.ListElement) layout.Dimensions {
	return FocusBorder(th, list.focused).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := gtx.Constraints.Max
		gtx.Constraints = layout.Exact(size)
		defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()

		changed := false
		grabbed := false

		itemHeight := gtx.Metric.Dp(list.ItemHeight)

		pointerClicked := false
		pointerHovered := false
		pointerPosition := f32.Point{}

		event.Op(gtx.Ops, &list.List)

		for {
			ev, ok := gtx.Event(
				key.Filter{
					Focus: list,
					Name: key.NameUpArrow + "|" + key.NameDownArrow + "|" +
						key.NameHome + "|" + key.NameEnd + "|" +
						key.NamePageUp + "|" + key.NamePageDown,
				},
				pointer.Filter{
					Target: &list.List,
					Kinds:  pointer.Press | pointer.Move,
				},
			)
			if !ok {
				break
			}

			switch ev := ev.(type) {
			case key.Event:
				if ev.State == key.Press {
					offset := 0
					switch ev.Name {
					case key.NameHome:
						offset = -length
					case key.NameEnd:
						offset = length
					case key.NameUpArrow:
						offset = -1
					case key.NameDownArrow:
						offset = 1
					case key.NamePageUp:
						offset = -list.List.Position.Count
					case key.NamePageDown:
						offset = list.List.Position.Count
					}

					if offset != 0 {
						target := max(list.Selected+offset, 0)
						if target >= length {
							target = length - 1
						}
						if list.Selected != target {
							list.Selected = target
							changed = true
						}
					}
				}
			case key.FocusEvent:
				if list.focused != ev.Focus {
					list.focused = ev.Focus
					gtx.Execute(op.InvalidateCmd{})
				}
			case pointer.Event:
				switch ev.Kind {
				case pointer.Press:
					if !list.focused && !grabbed {
						grabbed = true
						gtx.Execute(key.FocusCmd{Tag: list})
					}
					// TODO: find the item
					pointerClicked = true
					pointerPosition = ev.Position
				case pointer.Move:
					pointerHovered = true
					pointerPosition = ev.Position
				case pointer.Cancel:
					list.Hovered = -1
				}
			}
		}

		if pointerClicked || pointerHovered {
			clientClickY := list.Position.First*itemHeight + list.Position.Offset + int(pointerPosition.Y)
			target := clientClickY / itemHeight
			if 0 <= target && target <= length {
				if pointerClicked && list.Selected != target {
					list.Selected = target
				}
				if pointerHovered && list.Hovered != target {
					list.Hovered = target
				}
			}
		}

		if changed {
			pos := &list.List.Position
			switch {
			case list.Selected < pos.First+1:
				list.List.Position = layout.Position{First: list.Selected - 1}
			case pos.First+pos.Count-1 <= list.Selected:
				list.List.Position = layout.Position{First: list.Selected - pos.Count + 2}
			}
		}

		return material.List(th, &list.List).Layout(gtx, length,
			func(gtx layout.Context, index int) layout.Dimensions {
				gtx.Constraints = layout.Exact(image.Point{
					X: gtx.Constraints.Max.X,
					Y: itemHeight,
				})
				return element(gtx, index)
			})
	})
}

var theme = material.NewTheme()

var state = SelectList{
	List: widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	},
	ItemHeight: unit.Dp(24),
}

var items = []string{
	"alpha",
	"beta",
	"gamma",
	"delta",
	"iota",
	"kappa",
	"lambda",
	"my",

	"alpha",
	"beta",
	"gamma",
	"delta",
	"iota",
	"kappa",
	"lambda",
	"my",

	"alpha",
	"beta",
	"gamma",
	"delta",
	"iota",
	"kappa",
	"lambda",
	"my",

	"alpha",
	"beta",
	"gamma",
	"delta",
	"iota",
	"kappa",
	"lambda",
	"my",

	"alpha",
	"beta",
	"gamma",
	"delta",
	"iota",
	"kappa",
	"lambda",
	"my",

	"alpha",
	"beta",
	"gamma",
	"delta",
	"iota",
	"kappa",
	"lambda",
	"my",
}

var editor widget.Editor

func Layout(gtx layout.Context) {
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			focused := gtx.Source.Focused(&editor)
			return FocusBorder(theme, focused).Layout(gtx,
				material.Editor(theme, &editor, "Hint").Layout)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return state.Layout(theme, gtx, len(items),
				func(gtx layout.Context, index int) layout.Dimensions {
					defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

					switch {
					case state.Selected == index:
						paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xF0, B: 0xF0, A: 0xFF})
					case state.Hovered == index:
						paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xFF, B: 0xF0, A: 0xFF})
					}

					inset := layout.Inset{Top: 1, Right: 4, Bottom: 1, Left: 4}
					return inset.Layout(gtx, material.Body1(theme, items[index]).Layout)
				})
		}),
	)
}
