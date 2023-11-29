package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var n = flag.Int("n", 10, "line count")
var separate = flag.Bool("separate", false, "separate line clip paths")

func main() {
	flag.Parse()
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
	cloth := NewCloth(*n)

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
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			cloth.Update(0.015)
			cloth.Layout(gtx)

			op.InvalidateOp{}.Add(gtx.Ops)
			e.Frame(gtx.Ops)
		}
	}
}
