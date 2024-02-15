package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/egonelbre/expgio/f32color"
)

func main() {
	ui := NewUI()

	go func() {
		w := app.NewWindow(
			app.Title("Chat"),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

var (
	defaultMargin = unit.Dp(10)
)

type UI struct {
	Theme *material.Theme

	Groups *Groups
}

func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()

	ui.Groups = NewGroups(
		NewGroup("A", "Alpha"),
		NewGroup("B", "Bravo"),
		NewGroup("C", "Charlie"),
		NewGroup("D", "Delta"),
		NewGroup("E", "Echo"),
		NewGroup("F", "Foxtrot"),
		NewGroup("G", "Gopher"),
		NewGroup("H", "Hotel"),
		NewGroup("I", "India"),
		NewGroup("J", "Juliett"),
	)

	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.NextEvent().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			for {
				e, ok := gtx.Event(
					key.Filter{Name: key.NameEscape},
					key.Filter{Name: "1", Optional: key.ModCtrl},
					key.Filter{Name: "2", Optional: key.ModCtrl},
					key.Filter{Name: "3", Optional: key.ModCtrl},
					key.Filter{Name: "4", Optional: key.ModCtrl},
					key.Filter{Name: "5", Optional: key.ModCtrl},
					key.Filter{Name: "6", Optional: key.ModCtrl},
					key.Filter{Name: "7", Optional: key.ModCtrl},
					key.Filter{Name: "8", Optional: key.ModCtrl},
					key.Filter{Name: "9", Optional: key.ModCtrl},
				)
				if !ok {
					break
				}

				ev, ok := e.(key.Event)
				if !ok {
					continue
				}
				switch ev.Name {
				case key.NameEscape:
					return nil
				}
				if ev.Modifiers == key.ModCtrl {
					switch ev.Name {
					case "1":
						ui.activateGroup(w, 0)
					case "2":
						ui.activateGroup(w, 1)
					case "3":
						ui.activateGroup(w, 2)
					case "4":
						ui.activateGroup(w, 3)
					case "5":
						ui.activateGroup(w, 4)
					case "6":
						ui.activateGroup(w, 5)
					case "7":
						ui.activateGroup(w, 6)
					case "8":
						ui.activateGroup(w, 7)
					case "9":
						ui.activateGroup(w, 8)
					}
				}
			}

			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) activateGroup(w *app.Window, index int) {
	if index < 0 || index >= len(ui.Groups.Groups) {
		return
	}
	ui.Groups.Active = ui.Groups.Groups[index]
	w.Invalidate()
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.Groups.Layout(ui.Theme, gtx)
		}),
		layout.Rigid(ui.entries),
		layout.Flexed(1, Fill{Color: mainBackground}.Layout),
	)
}

func (ui *UI) entries(gtx layout.Context) layout.Dimensions {
	return entriesPanel.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(ui.entriesHeader),
		)
	})
}

func (ui *UI) entriesHeader(gtx layout.Context) layout.Dimensions {
	title := material.H4(ui.Theme, ui.Groups.Active.Name)
	title.Color = activeGroupTitle
	return entriesHeaderPanel.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, title.Layout)
		})
}

var (
	panelBackground  = f32color.HSL(0.5, 0.16, 0.26)
	panelBorder      = f32color.HSL(0.45, 0.16, 0.35)
	activeGroupTitle = f32color.HSL(0, 0, 0.97)

	mainBackground = f32color.HSL(0.0, 0.0, 0.97)

	borderWidth = unit.Dp(1)

	entriesPanel = Panel{
		Axis: layout.Vertical,
		Size: unit.Dp(270),

		Background:  panelBackground,
		Border:      panelBorder,
		BorderWidth: borderWidth,
	}

	entriesHeaderPanel = Panel{
		Axis: layout.Horizontal,
		Size: unit.Dp(80),

		Background:  panelBackground,
		Border:      panelBorder,
		BorderWidth: borderWidth,
	}
)
