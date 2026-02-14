package rig

import (
	"image"
	"image/color"
	"math/rand/v2"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/egonelbre/expgio/f32color"
)

const (
	BorderWidth = unit.Dp(4)
	AreaRadius  = unit.Dp(3)
	HandleSize  = unit.Dp(8)
	SplitSize   = unit.Dp(12)
	MinAreaSize = unit.Dp(50)
)

var (
	BackgroundColor = color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xFF}
	DragColor       = color.NRGBA{R: 0xFF, G: 0x88, B: 0x88, A: 0xFF}
	HoverColor      = color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	PassiveColor    = color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0x30}
	SplitColor      = color.NRGBA{R: 0x88, G: 0xCC, B: 0x88, A: 0x60}
	SplitHoverColor = color.NRGBA{R: 0x88, G: 0xFF, B: 0x88, A: 0xFF}
)

type Axis int

const (
	Horizontal Axis = iota // edge runs horizontally (Y boundary)
	Vertical               // edge runs vertically (X boundary)
)

type EditorID string

type EditorDef struct {
	Name string
	New  func() layout.Widget
}

// Screen holds a flat list of areas and dynamically computes edges.
type Screen struct {
	Registry      *Registry
	Bounds        image.Rectangle
	Areas         []*Area
	Edges         []*Edge
	Intersections []*Intersection
}

// Edge represents a shared boundary between areas.
type Edge struct {
	Axis   Axis
	Pos    int // X for vertical, Y for horizontal
	Start  int // start of the edge span
	End    int // end of the edge span
	Before []*Area
	After  []*Area

	Sizer        Sizer
	dragStartPos int
}

// Intersection represents a point where a vertical and horizontal edge cross.
type Intersection struct {
	Pos   image.Point
	VEdge *Edge
	HEdge *Edge

	Sizer        Sizer
	dragStartPos image.Point
}

type Area struct {
	Bounds image.Rectangle
	Editor *Editor
}

type Editor struct {
	Widget layout.Widget
}

func NewScreen() *Screen {
	s := &Screen{
		Registry: NewRegistry(),
		Bounds:   image.Rect(0, 0, 1024, 1024),
		Areas: []*Area{
			{
				Bounds: image.Rect(0, 0, 512, 512),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0, 0.6, 0.6)}.Layout,
				},
			},
			{
				Bounds: image.Rect(512, 0, 1024, 512),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0.2, 0.6, 0.6)}.Layout,
				},
			},
			{
				Bounds: image.Rect(0, 512, 512, 1024),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0.4, 0.6, 0.6)}.Layout,
				},
			},
			{
				Bounds: image.Rect(512, 512, 1024, 1024),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0.6, 0.6, 0.6)}.Layout,
				},
			},
		},
	}
	s.recomputeEdges()
	return s
}

func dprect(gtx layout.Context, r image.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: gtx.Metric.Dp(unit.Dp(r.Min.X)),
			Y: gtx.Metric.Dp(unit.Dp(r.Min.Y)),
		},
		Max: image.Point{
			X: gtx.Metric.Dp(unit.Dp(r.Max.X)),
			Y: gtx.Metric.Dp(unit.Dp(r.Max.Y)),
		},
	}
}

func (screen *Screen) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(gtx.Constraints.Max)

	defer clip.Rect(dprect(gtx, screen.Bounds)).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: BackgroundColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// Render areas.
	for _, area := range screen.Areas {
		area.Layout(gtx)
	}

	// Render edges.
	for _, edge := range screen.Edges {
		edge.Layout(screen, gtx)
	}

	// Render intersections.
	for _, inter := range screen.Intersections {
		inter.Layout(screen, gtx)
	}

	// Render split handles.
	screen.layoutSplitHandles(gtx)

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

// Area.Layout renders the editor content.
func (area *Area) Layout(gtx layout.Context) {
	inside := dprect(gtx, area.Bounds).Inset(gtx.Dp(BorderWidth) / 2)
	defer clip.UniformRRect(inside, gtx.Dp(AreaRadius)).Push(gtx.Ops).Pop()
	defer op.Offset(inside.Min).Push(gtx.Ops).Pop()

	egtx := gtx
	egtx.Constraints = layout.Exact(inside.Size())
	area.Editor.Widget(egtx)
}

// Edge.Layout renders the edge handle and processes drag events.
func (edge *Edge) Layout(screen *Screen, gtx layout.Context) {
	handleSize := gtx.Dp(HandleSize)

	// Use dragStartPos during drag to keep clip area stable.
	// Gio delivers drag events relative to the clip origin,
	// so a moving origin causes jittery feedback loops.
	pos := edge.Pos
	if edge.Sizer.Dragging() {
		pos = edge.dragStartPos
	}

	var hitRect image.Rectangle
	switch edge.Axis {
	case Vertical:
		px := gtx.Metric.Dp(unit.Dp(pos))
		startPx := gtx.Metric.Dp(unit.Dp(edge.Start))
		endPx := gtx.Metric.Dp(unit.Dp(edge.End))
		hitRect = image.Rect(px-handleSize/2, startPx, px+handleSize/2, endPx)
	case Horizontal:
		py := gtx.Metric.Dp(unit.Dp(pos))
		startPx := gtx.Metric.Dp(unit.Dp(edge.Start))
		endPx := gtx.Metric.Dp(unit.Dp(edge.End))
		hitRect = image.Rect(startPx, py-handleSize/2, endPx, py+handleSize/2)
	}

	defer op.Offset(hitRect.Min).Push(gtx.Ops).Pop()
	areaStack := clip.Rect(image.Rectangle{Max: hitRect.Size()}).Push(gtx.Ops)
	defer areaStack.Pop()

	var axis gesture.Axis
	var cursor pointer.Cursor
	switch edge.Axis {
	case Vertical:
		axis = gesture.Horizontal
		cursor = pointer.CursorColResize
	case Horizontal:
		axis = gesture.Vertical
		cursor = pointer.CursorRowResize
	}

	for _, ev := range edge.Sizer.Events(gtx, axis) {
		switch ev.Kind {
		case pointer.Press:
			edge.dragStartPos = edge.Pos
		case pointer.Drag:
			var deltaDp int
			switch edge.Axis {
			case Vertical:
				deltaDp = int(gtx.Metric.PxToDp(int(ev.Position.X - edge.Sizer.start.X)))
			case Horizontal:
				deltaDp = int(gtx.Metric.PxToDp(int(ev.Position.Y - edge.Sizer.start.Y)))
			}
			edge.setPosition(screen, edge.dragStartPos+deltaDp)
		case pointer.Release, pointer.Cancel:
			// Request one more frame so the handle redraws at the final position.
			// On this frame the clip was still at dragStartPos.
			gtx.Execute(op.InvalidateCmd{})
		}
	}
	cursor.Add(gtx.Ops)
	event.Op(gtx.Ops, &edge.Sizer)

	var c color.NRGBA
	switch {
	case edge.Sizer.Dragging():
		c = DragColor
	case edge.Sizer.Hovered():
		c = HoverColor
	default:
		c = PassiveColor
	}
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// setPosition moves the edge to newPos, clamping to min area size constraints.
func (edge *Edge) setPosition(screen *Screen, newPos int) {
	minSize := int(MinAreaSize)

	// Compute clamping bounds from before/after areas.
	minPos := screen.Bounds.Min.X
	maxPos := screen.Bounds.Max.X
	if edge.Axis == Horizontal {
		minPos = screen.Bounds.Min.Y
		maxPos = screen.Bounds.Max.Y
	}

	for _, a := range edge.Before {
		var areaMin int
		if edge.Axis == Vertical {
			areaMin = a.Bounds.Min.X + minSize
		} else {
			areaMin = a.Bounds.Min.Y + minSize
		}
		if areaMin > minPos {
			minPos = areaMin
		}
	}
	for _, a := range edge.After {
		var areaMax int
		if edge.Axis == Vertical {
			areaMax = a.Bounds.Max.X - minSize
		} else {
			areaMax = a.Bounds.Max.Y - minSize
		}
		if areaMax < maxPos {
			maxPos = areaMax
		}
	}

	if newPos < minPos {
		newPos = minPos
	}
	if newPos > maxPos {
		newPos = maxPos
	}

	delta := newPos - edge.Pos
	if delta == 0 {
		return
	}

	for _, a := range edge.Before {
		if edge.Axis == Vertical {
			a.Bounds.Max.X += delta
		} else {
			a.Bounds.Max.Y += delta
		}
	}
	for _, a := range edge.After {
		if edge.Axis == Vertical {
			a.Bounds.Min.X += delta
		} else {
			a.Bounds.Min.Y += delta
		}
	}
	edge.Pos = newPos
}

// Intersection.Layout handles both-axis resize at edge crossings.
func (inter *Intersection) Layout(screen *Screen, gtx layout.Context) {
	handleSize := gtx.Dp(HandleSize)

	// Use dragStartPos during drag to keep clip area stable.
	pos := inter.Pos
	if inter.Sizer.Dragging() {
		pos = inter.dragStartPos
	}

	px := gtx.Metric.Dp(unit.Dp(pos.X))
	py := gtx.Metric.Dp(unit.Dp(pos.Y))

	hitRect := image.Rect(px-handleSize, py-handleSize, px+handleSize, py+handleSize)

	defer op.Offset(hitRect.Min).Push(gtx.Ops).Pop()
	areaStack := clip.Rect(image.Rectangle{Max: hitRect.Size()}).Push(gtx.Ops)
	defer areaStack.Pop()

	for _, ev := range inter.Sizer.Events(gtx, gesture.Both) {
		switch ev.Kind {
		case pointer.Press:
			inter.dragStartPos = inter.Pos
		case pointer.Drag:
			dx := int(gtx.Metric.PxToDp(int(ev.Position.X - inter.Sizer.start.X)))
			dy := int(gtx.Metric.PxToDp(int(ev.Position.Y - inter.Sizer.start.Y)))
			inter.VEdge.setPosition(screen, inter.dragStartPos.X+dx)
			inter.HEdge.setPosition(screen, inter.dragStartPos.Y+dy)
			inter.Pos.X = inter.VEdge.Pos
			inter.Pos.Y = inter.HEdge.Pos
		case pointer.Release, pointer.Cancel:
			gtx.Execute(op.InvalidateCmd{})
		}
	}
	pointer.CursorCrosshair.Add(gtx.Ops)
	event.Op(gtx.Ops, &inter.Sizer)

	var c color.NRGBA
	switch {
	case inter.Sizer.Dragging():
		c = DragColor
	case inter.Sizer.Hovered():
		c = HoverColor
	default:
		c = PassiveColor
	}
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

// recomputeEdges rebuilds edges and intersections from current area bounds.
func (screen *Screen) recomputeEdges() {
	screen.Edges = screen.findEdges()
	screen.Intersections = screen.findIntersections()
}

func (screen *Screen) findEdges() []*Edge {
	var edges []*Edge

	// Find vertical edges (shared X boundaries).
	xVals := map[int]bool{}
	for _, a := range screen.Areas {
		xVals[a.Bounds.Min.X] = true
		xVals[a.Bounds.Max.X] = true
	}
	for x := range xVals {
		if x == screen.Bounds.Min.X || x == screen.Bounds.Max.X {
			continue
		}
		var before, after []*Area
		for _, a := range screen.Areas {
			if a.Bounds.Max.X == x {
				before = append(before, a)
			}
			if a.Bounds.Min.X == x {
				after = append(after, a)
			}
		}
		if len(before) > 0 && len(after) > 0 {
			// Find the span of the edge.
			start, end := screen.Bounds.Max.Y, screen.Bounds.Min.Y
			for _, a := range before {
				if a.Bounds.Min.Y < start {
					start = a.Bounds.Min.Y
				}
				if a.Bounds.Max.Y > end {
					end = a.Bounds.Max.Y
				}
			}
			edges = append(edges, &Edge{
				Axis:   Vertical,
				Pos:    x,
				Start:  start,
				End:    end,
				Before: before,
				After:  after,
			})
		}
	}

	// Find horizontal edges (shared Y boundaries).
	yVals := map[int]bool{}
	for _, a := range screen.Areas {
		yVals[a.Bounds.Min.Y] = true
		yVals[a.Bounds.Max.Y] = true
	}
	for y := range yVals {
		if y == screen.Bounds.Min.Y || y == screen.Bounds.Max.Y {
			continue
		}
		var before, after []*Area
		for _, a := range screen.Areas {
			if a.Bounds.Max.Y == y {
				before = append(before, a)
			}
			if a.Bounds.Min.Y == y {
				after = append(after, a)
			}
		}
		if len(before) > 0 && len(after) > 0 {
			start, end := screen.Bounds.Max.X, screen.Bounds.Min.X
			for _, a := range before {
				if a.Bounds.Min.X < start {
					start = a.Bounds.Min.X
				}
				if a.Bounds.Max.X > end {
					end = a.Bounds.Max.X
				}
			}
			edges = append(edges, &Edge{
				Axis:   Horizontal,
				Pos:    y,
				Start:  start,
				End:    end,
				Before: before,
				After:  after,
			})
		}
	}

	return edges
}

func (screen *Screen) findIntersections() []*Intersection {
	var intersections []*Intersection

	var vedges, hedges []*Edge
	for _, e := range screen.Edges {
		if e.Axis == Vertical {
			vedges = append(vedges, e)
		} else {
			hedges = append(hedges, e)
		}
	}

	for _, ve := range vedges {
		for _, he := range hedges {
			// Vertical edge spans Y from ve.Start to ve.End, at X = ve.Pos.
			// Horizontal edge spans X from he.Start to he.End, at Y = he.Pos.
			if ve.Start < he.Pos && he.Pos < ve.End &&
				he.Start < ve.Pos && ve.Pos < he.End {
				intersections = append(intersections, &Intersection{
					Pos:   image.Pt(ve.Pos, he.Pos),
					VEdge: ve,
					HEdge: he,
				})
			}
		}
	}

	return intersections
}

// Split handle types.
type splitHandle struct {
	Area  *Area
	Pos   image.Point // corner position in Dp
	Axis  Axis        // split direction
	Sizer Sizer
}

// Screen keeps split handles persistent across frames.
var splitHandles []*splitHandle

func (screen *Screen) layoutSplitHandles(gtx layout.Context) {
	// Rebuild split handle list if needed.
	if splitHandles == nil {
		splitHandles = screen.buildSplitHandles()
	}

	for _, sh := range splitHandles {
		sh.layout(screen, gtx)
	}
}

func (screen *Screen) buildSplitHandles() []*splitHandle {
	var handles []*splitHandle
	for _, area := range screen.Areas {
		// Add split handles at interior corners.
		bounds := area.Bounds

		type cornerDef struct {
			pos  image.Point
			axis Axis
		}

		corners := []cornerDef{
			// Top-left: split horizontally (create panel to the left) or vertically (above)
			{image.Pt(bounds.Min.X, bounds.Min.Y), Horizontal},
			{image.Pt(bounds.Min.X, bounds.Min.Y), Vertical},
			// Top-right
			{image.Pt(bounds.Max.X, bounds.Min.Y), Horizontal},
			{image.Pt(bounds.Max.X, bounds.Min.Y), Vertical},
			// Bottom-left
			{image.Pt(bounds.Min.X, bounds.Max.Y), Horizontal},
			{image.Pt(bounds.Min.X, bounds.Max.Y), Vertical},
			// Bottom-right
			{image.Pt(bounds.Max.X, bounds.Max.Y), Horizontal},
			{image.Pt(bounds.Max.X, bounds.Max.Y), Vertical},
		}

		for _, c := range corners {
			// Only add handles at interior corners (not on screen boundary).
			if c.pos.X == screen.Bounds.Min.X || c.pos.X == screen.Bounds.Max.X {
				if c.axis == Vertical {
					continue
				}
			}
			if c.pos.Y == screen.Bounds.Min.Y || c.pos.Y == screen.Bounds.Max.Y {
				if c.axis == Horizontal {
					continue
				}
			}
			// Skip if both coordinates are on screen boundary.
			if (c.pos.X == screen.Bounds.Min.X || c.pos.X == screen.Bounds.Max.X) &&
				(c.pos.Y == screen.Bounds.Min.Y || c.pos.Y == screen.Bounds.Max.Y) {
				continue
			}

			handles = append(handles, &splitHandle{
				Area: area,
				Pos:  c.pos,
				Axis: c.axis,
			})
		}
	}

	// Deduplicate handles at same position+axis.
	type key struct {
		pos  image.Point
		axis Axis
	}
	seen := map[key]bool{}
	var deduped []*splitHandle
	for _, h := range handles {
		k := key{h.Pos, h.Axis}
		if !seen[k] {
			seen[k] = true
			deduped = append(deduped, h)
		}
	}

	return deduped
}

func (sh *splitHandle) layout(screen *Screen, gtx layout.Context) {
	splitSize := gtx.Dp(SplitSize)
	px := gtx.Metric.Dp(unit.Dp(sh.Pos.X))
	py := gtx.Metric.Dp(unit.Dp(sh.Pos.Y))

	var hitRect image.Rectangle
	switch sh.Axis {
	case Vertical:
		hitRect = image.Rect(px-splitSize/2, py-splitSize/2, px+splitSize/2, py+splitSize/2)
	case Horizontal:
		hitRect = image.Rect(px-splitSize/2, py-splitSize/2, px+splitSize/2, py+splitSize/2)
	}

	defer op.Offset(hitRect.Min).Push(gtx.Ops).Pop()
	areaStack := clip.Rect(image.Rectangle{Max: hitRect.Size()}).Push(gtx.Ops)
	defer areaStack.Pop()

	for _, ev := range sh.Sizer.Events(gtx, gesture.Both) {
		if ev.Kind == pointer.Release {
			screen.splitArea(sh.Area, sh.Axis)
			splitHandles = nil // rebuild
		}
	}

	pointer.CursorCrosshair.Add(gtx.Ops)
	event.Op(gtx.Ops, &sh.Sizer)

	var c color.NRGBA
	if sh.Sizer.Hovered() || sh.Sizer.Dragging() {
		c = SplitHoverColor
	} else {
		c = SplitColor
	}
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func (screen *Screen) splitArea(area *Area, axis Axis) {
	minSize := int(MinAreaSize)

	var newArea *Area
	switch axis {
	case Vertical:
		mid := (area.Bounds.Min.Y + area.Bounds.Max.Y) / 2
		if area.Bounds.Max.Y-area.Bounds.Min.Y < minSize*2 {
			return
		}
		newArea = &Area{
			Bounds: image.Rect(area.Bounds.Min.X, mid, area.Bounds.Max.X, area.Bounds.Max.Y),
			Editor: &Editor{
				Widget: Color{Color: f32color.HSL(float32(rand.Float64()), 0.6, 0.6)}.Layout,
			},
		}
		area.Bounds.Max.Y = mid
	case Horizontal:
		mid := (area.Bounds.Min.X + area.Bounds.Max.X) / 2
		if area.Bounds.Max.X-area.Bounds.Min.X < minSize*2 {
			return
		}
		newArea = &Area{
			Bounds: image.Rect(mid, area.Bounds.Min.Y, area.Bounds.Max.X, area.Bounds.Max.Y),
			Editor: &Editor{
				Widget: Color{Color: f32color.HSL(float32(rand.Float64()), 0.6, 0.6)}.Layout,
			},
		}
		area.Bounds.Max.X = mid
	}

	if newArea != nil {
		screen.Areas = append(screen.Areas, newArea)
		screen.recomputeEdges()
	}
}

