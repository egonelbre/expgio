package main

import (
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
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
		w := app.NewWindow(app.Title("Selectbox"))
		var ops op.Ops
		for e := range w.Events() {
			switch e := e.(type) {
			case system.FrameEvent:

				gtx := layout.NewContext(&ops, e)
				Layout(gtx)
				e.Frame(gtx.Ops)

			case system.DestroyEvent:
				os.Exit(0)
			}
		}
	}()

	app.Main()
}

type SelectList struct {
	widget.List

	Selected int
	Hovered  int

	ItemHeight unit.Dp

	focused bool
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

func (list *SelectList) Layout(th *material.Theme, gtx layout.Context, length int, element layout.ListElement) layout.Dimensions {
	return FocusBorder(th, list.focused).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := gtx.Constraints.Max
		gtx.Constraints = layout.Exact(size)
		defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()

		key.InputOp{
			Tag:  list,
			Keys: key.NameHome + "|" + key.NameUpArrow + "|" + key.NameDownArrow + "|" + key.NameEnd,
		}.Add(gtx.Ops)

		pointer.InputOp{
			Tag:          list,
			Types:        pointer.Press | pointer.Move,
			ScrollBounds: image.Rectangle{},
		}.Add(gtx.Ops)

		changed := false
		grabbed := false

		itemHeight := gtx.Metric.Dp(list.ItemHeight)

		pointerClicked := false
		pointerHovered := false
		pointerPosition := f32.Point{}
		for _, ev := range gtx.Events(list) {
			switch ev := ev.(type) {
			case key.Event:
				if ev.State == key.Press {
					switch ev.Name {
					case key.NameHome:
						if list.Selected != 0 {
							list.Selected = 0
							changed = true
						}
					case key.NameEnd:
						if list.Selected != length-1 {
							list.Selected = length - 1
							changed = true
						}
					case key.NameUpArrow:
						if list.Selected > 0 {
							list.Selected--
							changed = true
						}
					case key.NameDownArrow:
						if list.Selected < length-1 {
							list.Selected++
							changed = true
						}
					}
				}
			case key.FocusEvent:
				if list.focused != ev.Focus {
					list.focused = ev.Focus
					op.InvalidateOp{}.Add(gtx.Ops)
				}
			case pointer.Event:
				switch ev.Type {
				case pointer.Press:
					if !list.focused && !grabbed {
						grabbed = true
						key.FocusOp{Tag: list}.Add(gtx.Ops)
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

var theme = material.NewTheme(gofont.Collection())
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
			return FocusBorder(theme, editor.Focused()).Layout(gtx,
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
