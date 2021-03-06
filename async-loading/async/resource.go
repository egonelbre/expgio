// SPDX-License-Identifier: Unlicense OR MIT

package async

import (
	"context"
	"sync"
	"sync/atomic"

	"gioui.org/layout"
)

// TODO:
//  * load multiple resources concurrently
//  * gradual loading
//  * loading progress
//  * resource limit based on size (e.g. max 10MiB of images)
//  * cancel loading when unloaded
//  * ensure purging doesn't block the rendering
//  * try to improve performance

type Loader struct {
	refresh sync.Cond
	mu      sync.Mutex

	maxLoaded int

	atomicActiveFrame   int64
	atomicFinishedFrame int64

	updated chan struct{}
	lookup  map[Tag]*resource
	queued  []*resource
}

func NewLoader(maxLoaded int) *Loader {
	loader := &Loader{
		updated:   make(chan struct{}, 1),
		lookup:    make(map[Tag]*resource),
		maxLoaded: maxLoaded,
	}
	loader.refresh.L = &loader.mu
	return loader
}

func (loader *Loader) Updated() <-chan struct{} { return loader.updated }

func (loader *Loader) update() {
	select {
	case loader.updated <- struct{}{}:
	default:
	}
}

type LoaderStats struct {
	Lookup int
	Queued int
}

func (loader *Loader) Stats() LoaderStats {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	return LoaderStats{
		Lookup: len(loader.lookup),
		Queued: len(loader.queued),
	}
}

func (loader *Loader) Frame(gtx layout.Context, w layout.Widget) layout.Dimensions {
	atomic.AddInt64(&loader.atomicActiveFrame, 1)
	dim := w(gtx)
	atomic.StoreInt64(&loader.atomicFinishedFrame, atomic.LoadInt64(&loader.atomicActiveFrame))

	// signal to maybe purge old entries
	loader.refresh.Signal()

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
		loader.refresh.Signal()
	}

	activeFrame := atomic.LoadInt64(&loader.atomicActiveFrame)
	atomic.StoreInt64(&r.atomicFrame, activeFrame)

	res := Resource{}
	res.State = State(atomic.LoadInt64(&r.atomicState))
	if res.State == Loaded {
		res.Value = r.value
	}
	return res
}

func (loader *Loader) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		loader.refresh.Signal()
	}()

	loader.mu.Lock()
	defer loader.mu.Unlock()

	for {
		loader.refresh.Wait()
		if ctx.Err() != nil {
			return
		}

		loader.purgeOld()

		for len(loader.queued) > 0 {
			active := loader.queued[0]
			loader.queued = loader.queued[1:]

			if atomic.LoadInt64(&active.atomicFrame) < atomic.LoadInt64(&loader.atomicFinishedFrame) {
				delete(loader.lookup, active.tag)
				continue
			}

			atomic.StoreInt64(&active.atomicState, int64(Loading))
			loader.update()

			loader.mu.Unlock()
			// TODO: implement concurrent loading
			value := active.load(ctx)
			loader.mu.Lock()

			active.value = value
			atomic.StoreInt64(&active.atomicState, int64(Loaded))

			loader.update()

			loader.purgeOld()
		}
	}
}

// TODO: this might end up blocking rendering
func (loader *Loader) purgeOld() {
	finishedFrame := atomic.LoadInt64(&loader.atomicFinishedFrame)
	for _, r := range loader.lookup {
		if len(loader.lookup) < loader.maxLoaded {
			break
		}
		if atomic.LoadInt64(&r.atomicFrame) < finishedFrame {
			delete(loader.lookup, r.tag)
		}
	}
}

type Tag interface{}

type Load func(ctx context.Context) interface{}

type Resource struct {
	State State
	Value interface{}
}

type resource struct {
	atomicFrame int64
	atomicState int64
	tag         Tag
	load        Load
	value       interface{}
}

type State byte

const (
	Queued State = iota
	Loading
	Loaded
)
