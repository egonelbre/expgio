package main

import (
	"flag"
	"image/color"
	"log"
	"math"
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

func main() {
	lines := flag.Int("line-count", 1000, "line count")
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
		if err := loop(*lines, w); err != nil {
			log.Println(err)
		}
	}()
	app.Main()
}

func loop(linecount int, w *app.Window) error {
	start := time.Now()

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			now := time.Since(start)
			_ = now

			paint.ColorOp{
				Color: color.NRGBA{
					R: byte(0xff),
					G: byte(0),
					B: byte(0),
					A: 0xFF,
				},
			}.Add(gtx.Ops)

			step := float32(math.Pi) / 90
			radius := float32(e.Size.Y) / 8.0
			center := f32.Point{
				X: float32(e.Size.X) / 2,
				Y: float32(e.Size.Y) / 2,
			}

			ripple := float32(float32(math.Sin(now.Seconds()*0.79))*2 + 5)
			wobbler := float32(float32(math.Sin(now.Seconds()*3.91))*radius/9 + radius/4)

			var builder clip.Path
			builder.Begin(gtx.Ops)

			pt := func(phi, radius, wobble float32) f32.Point {
				w := float32(math.Sin(float64(phi*ripple)))*wobble + radius
				return f32.Point{
					X: float32(math.Cos(float64(phi))) * w,
					Y: float32(math.Sin(float64(phi))) * w,
				}
			}

			start := pt(0, radius, 10.0)
			prev := start
			builder.Move(prev.Add(center))
			for phi := step; phi < 2*math.Pi; phi += step {
				next := pt(phi, radius, wobbler)
				builder.Line(next.Sub(prev))
				prev = next
			}
			builder.Line(start.Sub(prev))
			builder.Outline().Add(gtx.Ops)

			paint.PaintOp{}.Add(gtx.Ops)

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}
