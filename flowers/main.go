package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"math/rand"
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

	gtx := layout.NewContext(w.Queue())

	now := hrtime.Now()
	lastRender := time.Duration(0)

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)
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

type State struct {
	Root *Branch

	Camera f32.Point
}

func NewState() *State {
	state := &State{
		Root: NewRoot(),
	}
	state.Root.Path = []f32.Point{{X: 400, Y: 400}}
	return state
}

func (state *State) Update(delta float32) {
	state.Root.Update(delta)
	state.Camera = LerpPoint(0.7*delta, state.Camera, state.Root.Head())
}

func (state *State) Render(gtx *layout.Context) {
	fill(gtx, color.RGBA{R: 0x10, G: 0x14, B: 0x10, A: 0xFF})

	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	state.Root.Render(gtx)
}

type Branch struct {
	Time float32

	PathLimit int
	Path      []f32.Point

	Thickness  float32
	Lightness  float32
	Accelerate float32
	Speed      float32
	Direction  float32
	Turn       float32
	Length     float32
	Travel     float32

	IsRoot         bool
	Life           int
	SpawnCountdown float32
	SpawnInterval  float32
	Branches       []*Branch
}

func randomSnapped(min, max float32, snap float32) float32 {
	return min + Floor(rand.Float32()*(max-min)/snap)*snap
}

func random(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func NewRoot() *Branch {
	branch := &Branch{}
	branch.PathLimit = 300
	branch.Speed = 100.0
	branch.Direction = random(0, Tau)
	branch.IsRoot = true
	branch.Thickness = 8.0
	branch.Lightness = 0.7

	branch.SpawnInterval = 0.3
	branch.SpawnCountdown = branch.SpawnInterval

	branch.NextSegment()

	return branch
}

func NewBranch(root *Branch) *Branch {
	child := NewRoot()
	child.Thickness = root.Thickness * 0.3
	child.IsRoot = false
	child.Lightness = root.Lightness * 0.5
	child.Life = len(root.Path)
	child.Path = []f32.Point{root.Head()}
	child.Direction = root.Direction
	child.Turn = randomSnapped(-Tau*3/4, Tau*3/4, AngleSnap)
	return child
}

func (branch *Branch) Head() f32.Point {
	if len(branch.Path) == 0 {
		return f32.Point{}
	}
	return branch.Path[len(branch.Path)-1]
}

func (branch *Branch) NextSegment() {
	branch.Accelerate += 1.0
	branch.Turn = randomSnapped(-Tau/2, Tau/2, AngleSnap)
	branch.Length = randomSnapped(100, 200, 25)
	branch.Travel = branch.Length
}

func (branch *Branch) IsDying() bool {
	return !branch.IsRoot && (branch.Life < 0)
}
func (branch *Branch) IsDead() bool {
	return !branch.IsRoot && (branch.Life < 0) && (len(branch.Path) == 0)
}

func (branch *Branch) Update(dt float32) {
	branch.Life--

	alive := branch.Branches[:0]
	for _, child := range branch.Branches {
		child.Update(dt)
		if !child.IsDead() {
			alive = append(alive, child)
		}
	}
	branch.Branches = alive

	if branch.Accelerate > 0 {
		branch.Accelerate -= dt * 5
		dt += Bounce(branch.Accelerate, 0, 1, 1) * dt
	}
	branch.Time += dt

	if branch.IsRoot {
		branch.SpawnCountdown -= dt
		if branch.SpawnCountdown < 0 {
			branch.SpawnCountdown = branch.SpawnInterval
			child := NewBranch(branch)
			branch.Branches = append(branch.Branches, child)
		}
	}

	if !branch.IsDying() {
		distance := branch.Speed * dt

		if branch.IsRoot || branch.Travel > 0 {
			sn, cs := Sincos(branch.Direction)
			branch.Direction += branch.Turn * distance / branch.Length
			dir := f32.Point{X: cs, Y: sn}.Mul(distance)
			branch.Travel -= distance

			branch.Path = append(branch.Path, branch.Head().Add(dir))
			if len(branch.Path) > branch.PathLimit {
				copy(branch.Path, branch.Path[1:])
				branch.Path = branch.Path[:len(branch.Path)-1]
			}
		}

		if branch.IsRoot && branch.Travel <= 0 {
			branch.NextSegment()
		}
	} else {
		branch.Path = branch.Path[1:]
	}
}

func (branch *Branch) Render(gtx *layout.Context) {
	for _, child := range branch.Branches {
		child.Render(gtx)
	}

	if len(branch.Path) == 0 {
		return
	}

	pred := branch.Path[0]
	for i, pt := range branch.Path {
		if i > 0 && Len(pred.Sub(pt)) < 1 {
			continue
		}
		pred = pt

		p := float32(i) / float32(branch.PathLimit)

		radius := branch.Thickness * (Sin(p*4*Tau+branch.Time*6) + 5) / 5
		pp := float32(i) / float32(len(branch.Path))
		if pp > 0.85 {
			radius *= 1 - (pp-0.85)/0.3
		}
		if !branch.IsRoot {
			if pp < 0.05 {
				radius *= pp / 0.05
			}
		}

		squashcircle(gtx, pt, radius, color.RGBA{R: 0xff, G: 0, B: 0, A: 0xFF})
	}
}

func squashcircle(gtx *layout.Context, p f32.Point, r float32, color color.RGBA) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	paint.ColorOp{Color: color}.Add(gtx.Ops)

	var builder clip.Path
	builder.Begin(gtx.Ops)
	builder.Move(p.Add(f32.Point{X: 0, Y: -r}))

	builder.Cube(
		f32.Point{X: r, Y: 0},
		f32.Point{X: r, Y: 2 * r * 0.75},
		f32.Point{X: 0, Y: 2 * r * 0.75},
	)
	builder.Cube(
		f32.Point{X: -r, Y: 0},
		f32.Point{X: -r, Y: -2 * r * 0.75},
		f32.Point{X: 0, Y: -2 * r * 0.75},
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
