package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/loov/hrtime"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var n = flag.Int("n", 10, "line count")

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

	now := hrtime.Now()

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			next := hrtime.Now()
			dt := next - now
			if dt > 60*time.Millisecond {
				dt = 60 * time.Millisecond
			}
			now = next

			cloth.Update(float32(dt.Seconds()))
			cloth.Layout(gtx)

			op.InvalidateOp{}.Add(gtx.Ops)
			e.Frame(gtx.Ops)
		}
	}
}
