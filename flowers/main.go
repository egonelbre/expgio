package main

import (
	"flag"
	"image/color"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"gioui.org/op"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
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
	state := NewState()

	gtx := layout.NewContext(w.Queue())

	now := time.Now()
	lastRender := time.Duration(0)

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)

			timeSinceStart := time.Since(now)
			delta := timeSinceStart - lastRender
			lastRender = timeSinceStart

			state.Update(float32(delta.Seconds() * 2))
			state.Render(gtx)

			e.Frame(gtx.Ops)
			w.Invalidate()
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
	state.Root.Path = []f32.Point{{X: 100, Y: 100}}
	return state
}

func (state *State) Update(delta float32) {
	state.Root.Update(delta)
	state.Camera = LerpPoint(0.7*delta, state.Camera, state.Root.Head())
}

func (state *State) Render(gtx *layout.Context) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	screenSize := f32.Point{
		X: float32(gtx.Constraints.Width.Max),
		Y: float32(gtx.Constraints.Height.Max),
	}

	camera := Neg(state.Root.Head()).Add(screenSize.Mul(0.5))
	op.TransformOp{}.Offset(camera).Add(gtx.Ops)

	func() {
		var stack op.StackOp
		stack.Push(gtx.Ops)
		defer stack.Pop()

		var builder clip.Path
		builder.Begin(gtx.Ops)
		//builder.Move(f32.Point{})
		builder.Move(camera)
		builder.Line(f32.Point{10, 10})
		builder.Line(f32.Point{-10, 10})
		builder.End().Add(gtx.Ops)

		paint.ColorOp{
			Color: color.RGBA{R: 0xff, G: 0, B: 0, A: 0xFF},
		}.Add(gtx.Ops)

		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{},
				Max: f32.Point{
					X: float32(gtx.Constraints.Width.Max),
					Y: float32(gtx.Constraints.Height.Max),
				},
			}}.Add(gtx.Ops)
	}()

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
	//for _, child := range branch.Branches {
	//	child.Render(gtx)
	//}

	LineClip(gtx, branch.Path, func(i int) float32 {
		return 10

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

		return radius
	})
}

func LineClip(gtx *layout.Context, path []f32.Point, thickness func(int) float32) {
	if len(path) < 2 {
		return
	}

	{
		var stack op.StackOp
		stack.Push(gtx.Ops)
		defer stack.Pop()

		var builder clip.Path
		builder.Begin(gtx.Ops)
		pred := path[0]
		builder.Move(pred)
		for _, p := range path[1:] {
			builder.Line(p.Sub(pred))
			pred = p
		}
		for i := len(path) - 1; i >= 0; i-- {
			p := path[i].Add(f32.Point{-10, -10})
			builder.Line(p.Sub(pred))
			pred = p
		}
		builder.Line(path[0].Sub(pred))
		builder.End().Add(gtx.Ops)

		paint.ColorOp{
			Color: color.RGBA{R: 0xff, G: 0, B: 0, A: 0xFF},
		}.Add(gtx.Ops)

		paint.PaintOp{
			Rect: f32.Rectangle{
				Min: f32.Point{},
				Max: f32.Point{
					X: float32(gtx.Constraints.Width.Max),
					Y: float32(gtx.Constraints.Height.Max),
				},
			}}.Add(gtx.Ops)
	}

	return

	var left []f32.Point
	var right []f32.Point

	left = append(left, path[0])
	right = append(right, path[0])

	a := path[0]
	//var x1, x2, xn f32.Point

	for i, b := range path[1:] {
		if Len(a.Sub(b)) < 1 {
			continue
		}
		radius := thickness(i)

		abn := ScaleTo(SegmentNormal(a, b), radius)
		// segment-corners
		a1, a2 := a.Add(abn), a.Sub(abn)
		b1, b2 := b.Add(abn), b.Sub(abn)

		/*
			if i > 0 && radius > 0.5 {
				d := Dot(Rotate(xn), abn)
				if d < 0 {
					left = append(left, a1, b1)
					m.Push(x1, a1, a)
				} else {
					m.Push(x2, a2, a)
					right = append(right, a2, a)
				}
				m.Polygon(0)
			}
		*/

		left = append(left, a1, b1)
		right = append(right, a2, b2)

		a = b
		//x1, x2, xn = b1, b2, abn
	}

	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	paint.ColorOp{
		Color: color.RGBA{R: 0xff, G: 0, B: 0, A: 0xFF},
	}.Add(gtx.Ops)

	var builder clip.Path
	builder.Begin(gtx.Ops)
	pred := left[0]
	builder.Move(pred)
	for _, p := range left[1:] {
		builder.Line(p.Sub(pred))
		pred = p
	}
	for i := len(right) - 1; i >= 0; i-- {
		p := right[i]
		builder.Line(p.Sub(pred))
		pred = p
	}
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
