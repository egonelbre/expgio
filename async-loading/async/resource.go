// SPDX-License-Identifier: Unlicense OR MIT

package async

import (
	"context"
	"math"
	"sync"
	"sync/atomic"

	"gioui.org/layout"
)

type Loader struct {
	mu sync.Mutex

	config LoaderConfig

	atomicActiveFrame   int64
	atomicFinishedFrame int64

	signal  chan struct{}
	updated chan struct{}
	workCh  chan *resource
	wg      sync.WaitGroup

	lookup map[Tag]*resource
	queued []*resource

	// LRU doubly-linked list: head = most recent, tail = least recent
	lruHead   *resource
	lruTail   *resource
	lruLen    int
	totalBytes int64
}

type LoaderConfig struct {
	MaxCount    int   // 0 = unlimited
	MaxBytes    int64 // 0 = unlimited
	Concurrency int   // default 4
}

func NewLoader(config LoaderConfig) *Loader {
	if config.Concurrency <= 0 {
		config.Concurrency = 4
	}
	return &Loader{
		config:  config,
		signal:  make(chan struct{}, 1),
		updated: make(chan struct{}, 1),
		workCh:  make(chan *resource),
		lookup:  make(map[Tag]*resource),
	}
}

func (loader *Loader) Updated() <-chan struct{} { return loader.updated }

func (loader *Loader) update() {
	select {
	case loader.updated <- struct{}{}:
	default:
	}
}

func (loader *Loader) notify() {
	select {
	case loader.signal <- struct{}{}:
	default:
	}
}

type LoaderStats struct {
	Lookup     int
	Queued     int
	LRULen     int
	TotalBytes int64
}

func (loader *Loader) Stats() LoaderStats {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	return LoaderStats{
		Lookup:     len(loader.lookup),
		Queued:     len(loader.queued),
		LRULen:     loader.lruLen,
		TotalBytes: loader.totalBytes,
	}
}

func (loader *Loader) Frame(gtx layout.Context, w layout.Widget) layout.Dimensions {
	atomic.AddInt64(&loader.atomicActiveFrame, 1)
	dim := w(gtx)
	atomic.StoreInt64(&loader.atomicFinishedFrame, atomic.LoadInt64(&loader.atomicActiveFrame))
	loader.notify()
	return dim
}

func (loader *Loader) Schedule(tag Tag, load Load) Resource {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	r, ok := loader.lookup[tag]
	if !ok {
		r = &resource{
			tag:  tag,
			load: load,
		}
		loader.lookup[tag] = r
		loader.queued = append(loader.queued, r)
		loader.notify()
	}

	activeFrame := atomic.LoadInt64(&loader.atomicActiveFrame)
	atomic.StoreInt64(&r.atomicFrame, activeFrame)

	// Promote to LRU head if already in the list.
	if r.inLRU {
		loader.lruPromote(r)
	}

	res := Resource{}
	res.State = State(atomic.LoadInt64(&r.atomicState))
	switch res.State {
	case Loaded:
		res.Value = r.value
	case Loading:
		res.Progress = math.Float32frombits(uint32(atomic.LoadInt64(&r.atomicProgress)))
	}
	return res
}

func (loader *Loader) Run(ctx context.Context) {
	// Start worker goroutines.
	for range loader.config.Concurrency {
		loader.wg.Add(1)
		go loader.worker(ctx)
	}

	// Dispatcher loop.
	for {
		select {
		case <-ctx.Done():
			close(loader.workCh)
			loader.wg.Wait()
			return
		case <-loader.signal:
		}

		loader.dispatch(ctx)
	}
}

func (loader *Loader) dispatch(ctx context.Context) {
	loader.mu.Lock()
	queued := loader.queued
	loader.queued = nil
	loader.mu.Unlock()

	finishedFrame := atomic.LoadInt64(&loader.atomicFinishedFrame)

	for _, r := range queued {
		// Skip stale items that weren't visible in the last finished frame.
		if atomic.LoadInt64(&r.atomicFrame) < finishedFrame {
			loader.mu.Lock()
			delete(loader.lookup, r.tag)
			loader.mu.Unlock()
			continue
		}

		atomic.StoreInt64(&r.atomicState, int64(Loading))
		loader.update()

		// Create per-resource cancellation context.
		rctx, cancel := context.WithCancel(ctx)
		r.ctx = rctx
		r.cancel = cancel

		select {
		case loader.workCh <- r:
		case <-ctx.Done():
			cancel()
			return
		}
	}
}

func (loader *Loader) worker(_ context.Context) {
	defer loader.wg.Done()

	for r := range loader.workCh {
		if r.ctx.Err() != nil {
			continue
		}

		progressFn := func(p float32) {
			atomic.StoreInt64(&r.atomicProgress, int64(math.Float32bits(p)))
			loader.update()
		}

		value := r.load(r.ctx, progressFn)

		if r.ctx.Err() != nil {
			// Load was cancelled during execution; discard result.
			continue
		}

		var bytes int64
		if sv, ok := value.(SizedValue); ok {
			value = sv.Value
			bytes = sv.Bytes
		}

		r.value = value
		r.bytes = bytes
		atomic.StoreInt64(&r.atomicState, int64(Loaded))

		loader.mu.Lock()
		loader.lruAddHead(r)
		loader.evictExcess(16)
		loader.mu.Unlock()

		loader.notify()
		loader.update()
	}
}

// LRU list operations (must hold loader.mu).

func (loader *Loader) lruRemove(r *resource) {
	if !r.inLRU {
		return
	}
	if r.lruPrev != nil {
		r.lruPrev.lruNext = r.lruNext
	} else {
		loader.lruHead = r.lruNext
	}
	if r.lruNext != nil {
		r.lruNext.lruPrev = r.lruPrev
	} else {
		loader.lruTail = r.lruPrev
	}
	r.lruPrev = nil
	r.lruNext = nil
	r.inLRU = false
	loader.lruLen--
	loader.totalBytes -= r.bytes
}

func (loader *Loader) lruAddHead(r *resource) {
	loader.lruRemove(r) // remove first if already present
	r.lruNext = loader.lruHead
	r.lruPrev = nil
	if loader.lruHead != nil {
		loader.lruHead.lruPrev = r
	}
	loader.lruHead = r
	if loader.lruTail == nil {
		loader.lruTail = r
	}
	r.inLRU = true
	loader.lruLen++
	loader.totalBytes += r.bytes
}

func (loader *Loader) lruPromote(r *resource) {
	if loader.lruHead == r {
		return
	}
	loader.lruRemove(r)
	loader.lruAddHead(r)
}

func (loader *Loader) evictExcess(maxIter int) {
	finishedFrame := atomic.LoadInt64(&loader.atomicFinishedFrame)

	for i := 0; i < maxIter && loader.lruTail != nil; i++ {
		overCount := loader.config.MaxCount > 0 && loader.lruLen > loader.config.MaxCount
		overBytes := loader.config.MaxBytes > 0 && loader.totalBytes > loader.config.MaxBytes
		if !overCount && !overBytes {
			return
		}

		victim := loader.lruTail
		// Don't evict resources seen in the current frame.
		if atomic.LoadInt64(&victim.atomicFrame) >= finishedFrame {
			return
		}

		loader.lruRemove(victim)
		delete(loader.lookup, victim.tag)
		if victim.cancel != nil {
			victim.cancel()
		}
	}
}

type Tag any

type Load func(ctx context.Context, progress func(float32)) any

type SizedValue struct {
	Value any
	Bytes int64
}

type Resource struct {
	State    State
	Value    any
	Progress float32
}

type resource struct {
	atomicFrame    int64
	atomicState    int64
	atomicProgress int64

	tag    Tag
	load   Load
	value  any
	bytes  int64

	ctx    context.Context
	cancel context.CancelFunc

	// LRU doubly-linked list pointers.
	lruPrev *resource
	lruNext *resource
	inLRU   bool
}

type State byte

const (
	Queued State = iota
	Loading
	Loaded
)
