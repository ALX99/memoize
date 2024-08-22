package memoize

import (
	"reflect"
	"strings"
	"unsafe"
)

type memoizer[T any] struct {
	fv    reflect.Value
	ft    reflect.Type
	sb    strings.Builder
	cache Cache
}

// Auto returns a memoized version of the function f
// By default it uses [MapCache] as the cache and [Fn64aKeyer] as the [KeyerFunc].
func Auto[T any](f T, opts ...MemoizerOption) T {
	return newMemoizer(f, opts...).runFunc()
}

// Manual returns a function F(V) T that returns the memoized function f.
// The key needs to be explicitly be passed to F(V).
func Manual[T any, V comparable](f T, opts ...MemoizerOption) func(V) T {
	m := newMemoizer(f, opts...)

	return func(v V) T {
		// res, ok := m.cache.Get(key)
		// if ok {
		// 	returnVal := func() []reflect.Value {
		// 		return res
		// 	}
		// 	return reflect.MakeFunc(reflect.TypeOf(returnVal), func(args []reflect.Value) []reflect.Value { return returnVal() }).Interface().(T)
		// }

		return reflect.MakeFunc(m.ft, func(args []reflect.Value) []reflect.Value {
			key := []reflect.Value{reflect.ValueOf(v)}
			if res, ok := m.cache.Get(key); ok {
				return res
			}

			results := m.fv.Call(args)
			m.cache.Set(key, results)
			return results
		}).Interface().(T)
	}
}

// newMemoizer creates a new memoizer
func newMemoizer[T any](f T, opts ...MemoizerOption) *memoizer[T] {
	m := &memoizer[T]{
		fv:    reflect.ValueOf(f),
		ft:    reflect.TypeOf(f),
		cache: NewMapCache(Fn64aKeyer()),
	}
	for _, opt := range opts {
		opt((*memoizer[any])(unsafe.Pointer(m)))
	}

	return m
}

func (m *memoizer[T]) runFunc() T {
	return reflect.MakeFunc(m.ft, func(args []reflect.Value) []reflect.Value {
		if res, ok := m.cache.Get(args); ok {
			return res
		}

		results := m.fv.Call(args)
		m.cache.Set(args, results)
		return results
	}).Interface().(T)
}
