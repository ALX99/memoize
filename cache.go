package memoize

import (
	"context"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

// Cache is the interface needed to be satisfied by Cache implementations.
// Implementations needs be safe for concurrent use.
type Cache interface {
	// Get returns the cached result for the given args
	Get(args []reflect.Value) ([]reflect.Value, bool)
	// Set associates the given args with the given res
	Set(args []reflect.Value, res []reflect.Value)
	// Clear clears the cache
	Clear()
	// Len returns the number of cached argument calls
	Len() int
}

// MapCache is a cache that uses a map to store the results
type MapCache[T comparable] struct {
	cache map[T][]reflect.Value
	keyer KeyerFunc[T]

	cleanCtx context.Context
	cleanDur time.Duration

	sync.RWMutex
}

func NewMapCache[T comparable](keyer KeyerFunc[T], opts ...MapCacheOption) *MapCache[T] {
	cache := &MapCache[T]{cache: make(map[T][]reflect.Value), keyer: keyer}
	for _, opt := range opts {
		opt((*MapCache[any])(unsafe.Pointer(cache)))
	}

	if cache.cleanDur > 0 {
		go func() {
			t := time.NewTicker(cache.cleanDur)
			for {
				select {
				case <-t.C:
					cache.Clear()
				case <-cache.cleanCtx.Done():
					return
				}
			}
		}()
	}

	return cache
}

func (d *MapCache[T]) Get(args []reflect.Value) ([]reflect.Value, bool) {
	d.RLock()
	defer d.RUnlock()
	v, ok := d.cache[d.keyer(args)]
	return v, ok
}

func (d *MapCache[T]) Set(args []reflect.Value, res []reflect.Value) {
	d.Lock()
	d.cache[d.keyer(args)] = res
	d.Unlock()
}

func (d *MapCache[T]) Clear() {
	d.Lock()
	clear(d.cache)
	d.Unlock()
}

func (d *MapCache[T]) Len() int {
	d.RLock()
	defer d.RUnlock()
	return len(d.cache)
}

// MapCacheOption is an option that can be set to a MapCache
type MapCacheOption func(*MapCache[any])

// WithCleanDur sets the duration between when the cache is cleared
func WithCleanDur(ctx context.Context, dur time.Duration) MapCacheOption {
	return func(c *MapCache[any]) {
		c.cleanCtx = ctx
		c.cleanDur = dur
	}
}
