package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime/pprof"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var patch = (func() paint.ImageOp {
	f, err := os.Open("9patch.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	m, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	return paint.NewImageOp(m)
})()

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Println(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	gtx := layout.NewContext(w.Queue())

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)

			layout.UniformInset(unit.Dp(30)).Layout(gtx, func() {
				img := Patch9{
					Src:    patch,
					Left:   30,
					Top:    30,
					Right:  30,
					Bottom: 30,
				}
				img.Layout(gtx)
			})

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

type Patch9 struct {
	Src paint.ImageOp

	Left   int
	Top    int
	Right  int
	Bottom int
}

func (im Patch9) Layout(gtx *layout.Context) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Max, Y: cs.Height.Max}

	imageScale := float32(72.0 / 160.0)

	wf, hf := float32(d.X), float32(d.Y)
	_ = hf
	var s op.StackOp
	s.Push(gtx.Ops)
	clip.Rect{Rect: f32.Rectangle{Max: toPointF(d)}}.Op(gtx.Ops).Add(gtx.Ops)

	orig := im.Src.Rect
	defer func() { im.Src.Rect = orig }()

	leftf := float32(im.Left) * imageScale
	rightf := float32(im.Right) * imageScale
	topf := float32(im.Top) * imageScale
	bottomf := float32(im.Bottom) * imageScale
	_ = bottomf

	{ // top-left
		im.Src.Rect = orig
		im.Src.Rect.Max.X = orig.Min.X + im.Left
		im.Src.Rect.Max.Y = orig.Min.Y + im.Top
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: 0, Y: 0},
				Max: f32.Point{X: leftf, Y: topf},
			},
		}.Add(gtx.Ops)
	}

	{ // top-center
		im.Src.Rect = orig
		im.Src.Rect.Min.X = orig.Min.X + im.Left
		im.Src.Rect.Max.X = orig.Max.X - im.Right
		im.Src.Rect.Max.Y = orig.Min.Y + im.Top
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: leftf, Y: 0},
				Max: f32.Point{X: wf - rightf, Y: topf},
			},
		}.Add(gtx.Ops)
	}

	{ // top-right
		im.Src.Rect = orig
		im.Src.Rect.Min.X = orig.Max.X - im.Right
		im.Src.Rect.Max.Y = orig.Min.Y + im.Top
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: wf - rightf, Y: 0},
				Max: f32.Point{X: wf, Y: topf},
			},
		}.Add(gtx.Ops)
	}

	{ // center-left
		im.Src.Rect = orig
		im.Src.Rect.Max.X = orig.Min.X + im.Left
		im.Src.Rect.Min.Y = orig.Min.Y + im.Top
		im.Src.Rect.Max.Y = orig.Max.Y - im.Bottom
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: 0, Y: topf},
				Max: f32.Point{X: leftf, Y: hf - bottomf},
			},
		}.Add(gtx.Ops)
	}

	{ // center-center
		im.Src.Rect = orig
		im.Src.Rect.Min.X = orig.Min.X + im.Left
		im.Src.Rect.Max.X = orig.Max.X - im.Right
		im.Src.Rect.Min.Y = orig.Min.Y + im.Top
		im.Src.Rect.Max.Y = orig.Max.Y - im.Bottom
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: leftf, Y: topf},
				Max: f32.Point{X: wf - rightf, Y: hf - bottomf},
			},
		}.Add(gtx.Ops)
	}

	{ // center-right
		im.Src.Rect = orig
		im.Src.Rect.Min.X = orig.Max.X - im.Right
		im.Src.Rect.Min.Y = orig.Min.Y + im.Top
		im.Src.Rect.Max.Y = orig.Max.Y - im.Bottom
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: wf - rightf, Y: topf},
				Max: f32.Point{X: wf, Y: hf - bottomf},
			},
		}.Add(gtx.Ops)
	}

	{ // bottom-left
		im.Src.Rect = orig
		im.Src.Rect.Max.X = orig.Min.X + im.Left
		im.Src.Rect.Min.Y = orig.Max.Y - im.Bottom
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: 0, Y: hf - bottomf},
				Max: f32.Point{X: leftf, Y: hf},
			},
		}.Add(gtx.Ops)
	}

	{ // bottom-center
		im.Src.Rect = orig
		im.Src.Rect.Min.X = orig.Min.X + im.Left
		im.Src.Rect.Max.X = orig.Max.X - im.Right
		im.Src.Rect.Min.Y = orig.Max.Y - im.Bottom
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: leftf, Y: hf - bottomf},
				Max: f32.Point{X: wf - rightf, Y: hf},
			},
		}.Add(gtx.Ops)
	}

	{ // bottom-right
		im.Src.Rect = orig
		im.Src.Rect.Min.X = orig.Max.X - im.Right
		im.Src.Rect.Min.Y = orig.Max.Y - im.Bottom
		im.Src.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{X: wf - rightf, Y: hf - bottomf},
				Max: f32.Point{X: wf, Y: hf},
			},
		}.Add(gtx.Ops)
	}

	s.Pop()
	gtx.Dimensions = layout.Dimensions{Size: d}
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func bounds(gtx *layout.Context) f32.Rectangle {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	return f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
}

func fill(gtx *layout.Context, col color.RGBA) {
	dr := bounds(gtx)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
}
