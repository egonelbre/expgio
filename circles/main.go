package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/gesture"
	"gioui.org/io/key"
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
		w := new(app.Window)
		w.Option(
			app.Title("Circles"),
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
	Theme    *material.Theme
	Overlays []*Overlay
	Change   gesture.Click
}

func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Overlays = []*Overlay{
		NewOverlay("Hello", f32.Pt(0.5, 0.5)),
	}
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			_, ok := gtx.Event(key.Filter{Name: key.NameEscape})
			if ok {
				return nil
			}

			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(gtx.Constraints.Max)

	ui.Change.Add(gtx.Ops)
	for {
		click, ok := ui.Change.Update(gtx.Source)
		if !ok {
			break
		}
		if click.Kind != gesture.KindClick {
			continue
		}
		last := ui.Overlays[len(ui.Overlays)-1]

		pos := layout.FPt(click.Position)
		max := layout.FPt(gtx.Constraints.Max)
		pos.X /= max.X
		pos.Y /= max.Y

		next := NewOverlay(last.Text, pos)
		next.Fg = last.Bg
		next.Dark = !last.Dark
		if next.Dark {
			next.Bg = f32color.HSL(rand.Float32(), 0.5, 0.15)
		} else {
			next.Bg = f32color.HSL(rand.Float32(), 0.5, 0.85)
		}
		ui.Overlays = append(ui.Overlays, next)
	}

	for len(ui.Overlays) >= 2 && ui.Overlays[1].Show.Done() {
		ui.Overlays = ui.Overlays[1:]
	}
	for _, overlay := range ui.Overlays {
		_ = overlay.Layout(ui.Theme, gtx)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

type Overlay struct {
	Dark bool
	Fg   color.NRGBA
	Bg   color.NRGBA

	Flood f32.Point
	Text  string
	Show  Animation
}

func NewOverlay(t string, flood f32.Point) *Overlay {
	return &Overlay{
		Dark:  true,
		Fg:    f32color.HSL(0, 0.5, 0.85),
		Bg:    f32color.HSL(0, 0.5, 0.15),
		Flood: flood,
		Text:  t,
		Show:  NewAnimation(1500 * time.Millisecond),
	}
}

func (overlay *Overlay) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	progress := overlay.Show.Update(gtx)
	if progress < 1 {
		p := layout.FPt(gtx.Constraints.Max)
		p.X *= overlay.Flood.X
		p.Y *= overlay.Flood.Y

		r := magnitude(gtx.Constraints.Max) * progress
		var rect image.Rectangle
		rect.Min = p.Sub(f32.Pt(r, r)).Round()
		rect.Max = p.Add(f32.Pt(r, r)).Round()

		defer clip.UniformRRect(rect, int(r)).Push(gtx.Ops).Pop()
	}

	paint.ColorOp{Color: overlay.Bg}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	selectColorMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: overlay.Fg}.Add(gtx.Ops)
	selectColor := selectColorMacro.Stop()

	macro := op.Record(gtx.Ops)
	textgtx := gtx
	textgtx.Constraints.Min = image.Point{}
	dims := widget.Label{}.Layout(textgtx, th.Shaper, font.Font{Weight: font.Bold}, 128, overlay.Text, selectColor)
	text := macro.Stop()

	center := gtx.Constraints.Max.Div(2)
	defer op.Offset(image.Point{
		X: center.X - dims.Size.X/2,
		Y: center.Y - dims.Size.Y/2,
	}).Push(gtx.Ops).Pop()

	text.Add(gtx.Ops)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func magnitude(p image.Point) float32 {
	return float32(math.Sqrt(float64(p.X*p.X + p.Y*p.Y)))
}

type Animation struct {
	now      time.Time
	progress time.Duration
	duration time.Duration
}

func NewAnimation(duration time.Duration) Animation {
	return Animation{duration: duration}
}

func (anim *Animation) Update(gtx layout.Context) float32 {
	if anim.now.IsZero() {
		anim.now = gtx.Now
	}

	delta := gtx.Now.Sub(anim.now)
	anim.now = gtx.Now

	if delta > 15*time.Millisecond {
		delta = 15 * time.Millisecond
	}

	if anim.progress < anim.duration {
		anim.progress += delta
		if anim.progress > anim.duration {
			anim.progress = anim.duration
		}
		gtx.Execute(op.InvalidateCmd{})
	}

	return float32(float64(anim.progress) / float64(anim.duration))
}

func (anim *Animation) Progress() float32 {
	return float32(float64(anim.progress) / float64(anim.duration))
}

func (anim *Animation) Done() bool { return anim.Progress() >= 1 }
