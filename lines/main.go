package main

import (
	"flag"
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

	"github.com/loov/hrtime"
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
	state := NewState()

	now := hrtime.Now()
	lastRender := time.Duration(0)

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e.Queue, e.Config, e.Size)

			Stroke(gtx.Ops)

			op.InvalidateOp{}.Add(gtx.Ops)

			timeSinceStart := hrtime.Since(now)
			delta := timeSinceStart - lastRender
			lastRender = timeSinceStart

			dt := float32(delta.Seconds())
			if dt > 0.016 {
				dt = 0.016
			}
			state.Update(dt)
			state.Render(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

func Stroke(ops *op.Ops, points []f32.Point, thickness float32) op.CallOp {
	var path clip.Path
	path.Begin(ops)
	if len(points) < 2 {
		return path.End()
	}

	R := abs(thickness / 2)
	R2 := R * R
	s2R2 := math.Sqrt2 * R2

	// draw each segment, where
	// x1-------^--------a1-------^---------b1
	// |        | xn     |        | abn      |
	// | - - - - - - - - a - - - - - - - - - b
	// |                 |                   |
	// x2----------------a2-----------------b2
	// x1, x2, xn are the previous segments end corners and normal
	a, b := points[0], points[1]
	xn := segmentNormalScaled(a, b, R)
	x1, x2 := a.Add(xn), a.Sub(xn)

	for _, b := range points[1:] {
		abn := segmentNormalScaled(a, b, R)

		if dot(xn, abn) == 0 { // straight segment
			b1, b2 := b.Add(xn), b.Sub(xn)
			// QUAD(x1, b1, b2, x2)
			x1, x2 = b1, b2
		} else {
			b1, b2 := a.Add(xn), a.Sub(xn)
			// QUAD(x1, b1, b2, x2)

			x1, x2 = a.Add(abn), a.Sub(abn)
			if dot(rotate90(xn), abn) < 0 {
				// TRI(b1, x1, a)
			} else {
				// TRI(b2, a, x)
			}
		}
	}

	return path.End()
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func segmentNormal(a, b f32.Point) f32.Point {
	return rotate90(a.Sub(b))
}

func segmentNormalScaled(a, b f32.Point, size float32) f32.Point {
	return scaleTo(segmentNormal(a, b), size)
}

func dot(a, b f32.Point)         { return a.X*b.X + a.Y*b.Y }
func length(v f32.Point) float32 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

func scaleTo(v f32.Point, len f32.Point) f32.Point {
	ilen := length(v)
	if ilen > 0 {
		ilen = size / ilen
	}
	return scale(a, ilen)
}

func scale(v f32.Point, p float32) f32.Point {
	return f32.Point{X: v.X * p, Y: v.Y * p}
}

func rotate90(v f32.Point) f32.Point {
	return f32.Point{X: -v.Y, Y: v.X}
}

/*
func squashcircle(gtx layout.Context, p f32.Point, r float32, color color.RGBA) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	op.TransformOp{}.Offset(p).Add(gtx.Ops)

	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Pt(0, -r))
	path.Cube(f32.Pt(r, 0), f32.Pt(r, 2*r*0.75), f32.Pt(0, 2*r*0.75))
	path.Cube(f32.Pt(-r, 0), f32.Pt(-r, -2*r*0.75), f32.Pt(0, -2*r*0.75))
	path.End().Add(gtx.Ops)

	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(-r, -r, r, r)}.Add(gtx.Ops)
}

func fill(gtx layout.Context, col color.RGBA) layout.Dimensions {
	dr := f32.Rectangle{Max: layout.FPt(gtx.Constraints.Max)}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
*/
