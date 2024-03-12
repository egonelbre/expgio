package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"os"
	"runtime/pprof"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()
	go func() {
		w := &app.Window{}
		if err := loop(w); err != nil {
			log.Println(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
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

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			const size = 10
			for y := 0; y < e.Size.Y; y += size {
				for x := 0; x < e.Size.X; x += size {
					paint.ColorOp{
						Color: color.NRGBA{
							R: byte(x),
							G: byte(y),
							B: byte(x * y),
							A: 0xFF,
						},
					}.Add(gtx.Ops)
					stack := clip.RRect{
						Rect: image.Rect(x, y, x+size, y+size),
					}.Push(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					stack.Pop()
				}
			}

			e.Frame(gtx.Ops)
		}
	}
}
