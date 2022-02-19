package main

import (
	"flag"
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
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	start := time.Now()
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			fill(gtx, color.NRGBA{R: 0x10, G: 0x14, B: 0x10, A: 0xFF})
			render(gtx, float32(time.Since(start).Seconds()))

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

func render(gtx layout.Context, t float32) {
	t *= 0.3

	screenSize := layout.FPt(gtx.Constraints.Max)
	radius := Min(screenSize.X, screenSize.Y) * 0.7 / 3

	const n = 256
	for i := 0; i < n; i++ {
		r := float32(i) / n
		p := curve(t+r*1.2+Sin(t+r)*3, radius+Sin(r*1.1)*30)
		q := radius * 0.3 * pcurve(float32(i)/(n-1), 1.5, 0.6)
		squashCircle(gtx, p.Add(screenSize.Mul(0.5)), q, color.NRGBA{R: 0xff, G: 0xd7, B: byte(i), A: 0xFF})
	}
}

func squashCircle(gtx layout.Context, p f32.Point, r float32, color color.NRGBA) {
	defer op.Offset(p).Push(gtx.Ops).Pop()

	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Pt(0, -r))
	path.Cube(f32.Pt(r, 0), f32.Pt(r, 2*r*0.75), f32.Pt(0, 2*r*0.75))
	path.Cube(f32.Pt(-r, 0), f32.Pt(-r, -2*r*0.75), f32.Pt(0, -2*r*0.75))

	paint.FillShape(gtx.Ops, color, clip.Outline{Path: path.End()}.Op())
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

func fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Min}
}
