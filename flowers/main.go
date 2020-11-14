package main

import (
	"flag"
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

	"github.com/egonelbre/expgio/f32color"
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
			gtx := layout.NewContext(&ops, e)
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
	Camera f32.Point

	Root *Branch
}

func NewState() *State {
	state := &State{
		Root: NewRoot(),
	}
	return state
}

func (state *State) Update(delta float32) {
	state.Root.Update(delta)
	state.Camera = LerpPoint(0.8*delta, state.Camera, state.Root.Head())
}

func (state *State) Render(gtx layout.Context) {
	// fill(gtx, color.RGBA{R: 0x10, G: 0x14, B: 0x10, A: 0xFF})
	fill(gtx, color.RGBA{R: 0xFF, G: 0xFF, B: 0xEE, A: 0xFF})

	defer op.Push(gtx.Ops).Pop()

	screenSize := layout.FPt(gtx.Constraints.Min)
	offset := Neg(state.Camera).Add(screenSize.Mul(0.5))
	op.Offset(offset).Add(gtx.Ops)

	state.Root.Render(gtx)
}

type Branch struct {
	Time float32

	PathLimit int
	Path      []f32.Point

	Fill   color.RGBA
	Stroke color.RGBA

	Thickness  float32
	Lightness  float32
	Accelerate float32
	Speed      float32
	Direction  float32
	Turn       float32
	Length     float32
	Travel     float32
	Bounce     float32

	IsRoot         bool
	Life           int
	SpawnCountdown float32
	SpawnInterval  float32
	SpawnSide      bool
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
	branch.PathLimit = 200
	branch.Speed = 150.0
	branch.Direction = random(0, Tau)
	branch.IsRoot = true
	branch.Thickness = 12.0
	branch.Lightness = 0.8

	branch.Fill = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	branch.Stroke = f32color.HSL(0, branch.Lightness*0.2, 0.4)

	branch.SpawnInterval = 0.3
	branch.SpawnCountdown = branch.SpawnInterval

	branch.NextSegment()

	return branch
}

func NewBranch(root *Branch) *Branch {
	child := NewRoot()
	child.Thickness = root.Thickness * 0.3
	child.IsRoot = false
	child.Lightness = root.Lightness * 0.8
	child.Life = len(root.Path)
	child.Path = []f32.Point{root.Head()}
	child.Direction = root.Direction
	if root.SpawnSide {
		child.Turn = randomSnapped(-Tau*4/4, Tau*1/4, AngleSnap)
	} else {
		child.Turn = randomSnapped(-Tau*1/4, Tau*4/4, AngleSnap)
	}
	child.Turn += root.Turn
	root.SpawnSide = !root.SpawnSide
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

	branch.Bounce = 0
}

func (branch *Branch) IsDying() bool {
	return !branch.IsRoot && (branch.Life < 0)
}
func (branch *Branch) IsDead() bool {
	return !branch.IsRoot && (branch.Life < 0) && (len(branch.Path) == 0)
}

func (branch *Branch) Update(dt float32) {
	branch.Life--
	if branch.Bounce < 1 {
		branch.Bounce += dt / 2.5
		if branch.Bounce > 1 {
			branch.Bounce = 1
		}
	}

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

func (branch *Branch) Render(gtx layout.Context) {
	stroke := 8 - Bounce(branch.Bounce, 0, 4, 1)
	branch.renderPath(gtx, stroke, branch.Stroke)

	for _, child := range branch.Branches {
		child.Render(gtx)
	}

	branch.renderPath(gtx, 0, branch.Fill)
}

func (branch *Branch) renderPath(gtx layout.Context, radiusAdd float32, color color.RGBA) {
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

		squashcircle(gtx, pt, radius+radiusAdd, color)
	}
}

func squashcircle(gtx layout.Context, p f32.Point, r float32, color color.RGBA) {
	defer op.Push(gtx.Ops).Pop()

	op.Offset(p).Add(gtx.Ops)

	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Pt(0, -r))
	path.Cube(f32.Pt(r, 0), f32.Pt(r, 2*r*0.75), f32.Pt(0, 2*r*0.75))
	path.Cube(f32.Pt(-r, 0), f32.Pt(-r, -2*r*0.75), f32.Pt(0, -2*r*0.75))
	path.Outline().Add(gtx.Ops)

	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func fill(gtx layout.Context, col color.RGBA) layout.Dimensions {
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
