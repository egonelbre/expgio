package main

import (
	"flag"
	"image/color"
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
)

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
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			const size = 10
			for y := 0; y < e.Size.Y; y += size {
				for x := 0; x < e.Size.X; x += size {
					stack := op.Push(gtx.Ops)
					paint.ColorOp{
						Color: color.RGBA{
							R: byte(x),
							G: byte(y),
							B: byte(x * y),
							A: 0xFF,
						},
					}.Add(gtx.Ops)
					clip.RRect{Rect: f32.Rectangle{
						Min: f32.Point{
							X: float32(x),
							Y: float32(y),
						},
						Max: f32.Point{
							X: float32(x + size),
							Y: float32(y + size),
						},
					}}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					stack.Pop()
				}
			}

			e.Frame(gtx.Ops)
		}
	}
}
