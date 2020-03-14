package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

const AngleSnap = Tau / 8

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
	start := time.Now()

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)

			fill(gtx, color.RGBA{R: 0x10, G: 0x14, B: 0x10, A: 0xFF})
			render(gtx, float32(time.Since(start).Seconds()))

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

func render(gtx *layout.Context, t float32) {
	t *= 0.3

	screenSize := f32.Point{
		X: float32(gtx.Constraints.Width.Max),
		Y: float32(gtx.Constraints.Height.Max),
	}

	radius := Min(screenSize.X, screenSize.Y) * 0.7 / 3

	const n = 256
	for i := 0; i < n; i++ {
		r := float32(i) / n
		p := curve(t+r*1.2+Sin(t+r)*3, radius+Sin(r*1.1)*30)

		var stack op.StackOp
		stack.Push(gtx.Ops)

		paint.ColorOp{
			Color: color.RGBA{R: 0xff, G: 0xd7, B: byte(i), A: 0xFF},
		}.Add(gtx.Ops)

		var builder clip.Path
		builder.Begin(gtx.Ops)
		builder.Move(p.Add(screenSize.Mul(0.5)))

		q := radius * 0.3 * pcurve(float32(i)/(n-1), 1.5, 0.6)

		builder.Cube(
			f32.Point{X: q, Y: 0},
			f32.Point{X: q, Y: 2 * q * 0.75},
			f32.Point{X: 0, Y: 2 * q * 0.75},
		)
		builder.Cube(
			f32.Point{X: -q, Y: 0},
			f32.Point{X: -q, Y: -2 * q * 0.75},
			f32.Point{X: 0, Y: -2 * q * 0.75},
		)
		builder.End().Add(gtx.Ops)

		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{},
				Max: f32.Point{
					X: float32(gtx.Constraints.Width.Max),
					Y: float32(gtx.Constraints.Height.Max),
				},
			}}.Add(gtx.Ops)

		stack.Pop()
	}
}

func curve(t, s float32) f32.Point {
	sn, cs := Sincos(-t * Tau)

	t = Mod(t, 2)
	if t < 1 {
		return f32.Point{X: -s + cs*s, Y: sn * s}
	} else {
		return f32.Point{X: +s - cs*s, Y: sn * s}
	}
}

func pcurve(p, a, b float32) float32 {
	k := Pow(a+b, a+b) / (Pow(a, a) * Pow(b, b))
	return k * Pow(p, a) * Pow(1-p, b)
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
