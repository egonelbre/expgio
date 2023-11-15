package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/egonelbre/expgio/f32color"
)

func main() {
	ui := NewUI()

	go func() {
		w := app.NewWindow(app.Title("Panels"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

var defaultMargin = unit.Dp(10)

type UI struct {
	Theme  *material.Theme
	Panels []*Panel
}

func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	for i := 0; i < 4; i++ {
		ui.Panels = append(ui.Panels, NewPanel(ui))
	}
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.NextEvent().(type) {
		case system.FrameEvent:

			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case key.Event:
			switch e.Name {
			case key.NameEscape:
				return nil
			}

		case system.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	unchanged := append([]*Panel{}, ui.Panels...)
	var left *Panel
	for _, panel := range unchanged {
		panel.Update(left, gtx)
		_ = panel.Layout(ui.Theme, gtx)
		left = panel
	}
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (ui *UI) AddPanel(before *Panel, add *Panel) {
	for k, p := range ui.Panels {
		if p == before {
			k++
			ui.Panels = append(ui.Panels, add)
			copy(ui.Panels[k+1:], ui.Panels[k:])
			ui.Panels[k] = add
			return
		}
	}
	ui.Panels = append(ui.Panels, add)
}

func (ui *UI) ClosePanel(toClose *Panel) {
	for k, p := range ui.Panels {
		if p == toClose {
			ui.Panels = append(ui.Panels[:k], ui.Panels[k+1:]...)
			return
		}
	}
}

type Panel struct {
	UI *UI

	DeltaTime DeltaTime

	LeftPx  float32
	WidthPx float32

	TargetLeftPx float32

	Color color.NRGBA
	Title string

	Insert widget.Clickable
	Close  widget.Clickable
}

func NewPanel(ui *UI) *Panel {
	return &Panel{
		UI:      ui,
		WidthPx: 200,
		Color:   f32color.HSL(rand.Float32(), 0.2+0.2*rand.Float32(), 0.7+0.1*rand.Float32()),
		Title:   fmt.Sprintf("%04x", rand.Int()&0xFFFF),
	}
}

func (panel *Panel) Update(left *Panel, gtx layout.Context) {
	target := float32(0)
	if left != nil {
		target = left.TargetLeftPx + left.WidthPx
	}
	panel.TargetLeftPx = target

	dt := float32(panel.DeltaTime.Update(gtx).Seconds())
	if panel.LeftPx == target {
		return
	}
	dt *= 6
	if dt > 0.5 {
		dt = 0.5
	}

	panel.LeftPx = panel.LeftPx*(1-dt) + target*dt
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (panel *Panel) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(image.Pt(int(panel.WidthPx), int(gtx.Constraints.Max.Y)))

	defer op.Offset(image.Pt(int(panel.LeftPx), 0)).Push(gtx.Ops).Pop()
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: panel.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	for range panel.Insert.Update(gtx) {
		panel.UI.AddPanel(panel, NewPanel(panel.UI))
	}

	_ = panel.Insert.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{
			Size: gtx.Constraints.Max,
		}
	})

	if panel.Close.Clicked(gtx) {
		panel.UI.ClosePanel(panel)
	}

	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(material.Body1(th, panel.Title).Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(th.TextSize)}.Layout), // TODO: fix add gtx.DpToSp and gtx.SpToDp
					layout.Rigid(material.Button(th, &panel.Close, "Close").Layout),
				)
			})
	})
}

type DeltaTime struct{ Last time.Time }

func (dt *DeltaTime) Update(gtx layout.Context) time.Duration {
	if dt.Last.IsZero() {
		dt.Last = gtx.Now
	}
	delta := gtx.Now.Sub(dt.Last)
	if delta == 0 {
		delta = 8 * time.Millisecond
	}
	if delta > 16*time.Millisecond {
		delta = 16 * time.Millisecond
	}
	dt.Last = gtx.Now
	return delta
}
