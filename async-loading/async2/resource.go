package async2

import (
	"context"
	"sync/atomic"

	"gioui.org/layout"
)

type Tag interface{}

type Frame int64

type Cache struct {
	atomicActiveFrame   int64
	atomicFinishedFrame int64

	maxLoaders int

	queue   *Queue
	loaders []*Loader
}

func (cache *Cache) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	newFrame := atomic.AddInt64(&cache.atomicActiveFrame, 1)
	defer atomic.StoreInt64(&cache.atomicFinishedFrame, newFrame)

	return w(gtx)
}

type Loader struct {
}

type Queue struct {
	queued []*Resource
}

type Resource struct {
	AtomicRequestedAt Frame
	AtomicState       Frame

	Tag  Tag
	Load Load

	Priority int64
}

type Load func(ctx context.Context) interface{}

func (cache *Cache) Fetch(tag Tag, load Load) Resource {

}
