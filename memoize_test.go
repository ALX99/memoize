package memoize

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleAuto() {
	i := 0
	expensiveFunc := func(x int) {
		i += x
	}

	memoized := Auto(expensiveFunc) // memoize.Auto

	memoized(100)
	memoized(1)
	memoized(1)

	fmt.Println(i)
	// Output: 101
}

func ExampleManual() {
	i := 0
	expensiveFunc := func(x int) {
		i += x
	}

	memoized := Manual[func(int), string](expensiveFunc) // memoize.Manual

	memoized("key1")(100)
	memoized("key2")(1)
	memoized("key2")(1)

	fmt.Println(i)
	// Output: 101
}

func BenchmarkMemoize(b *testing.B) {
	type ss struct{ int }
	f := func(x, y, z int, a ss) int {
		return x + y
	}
	memoizedFunc := Auto(f)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memoizedFunc(1, 2, 3, ss{1})
	}
}

func TestMemoieManual(t *testing.T) {
	i := 0
	incr := Manual[func(int), int](func(x int) { i += x })
	incr(100)(100)
	incr(1)(1)
	incr(1)(1)
	assert.Equal(t, 101, i)
}

func TestCacheGrowsCorrectly(t *testing.T) {
	f := func(x int) {}
	m := newMemoizer[func(int)](f)
	m.runFunc()(1)
	assert.Equal(t, m.cache.Len(), 1)
	m.runFunc()(1)
	assert.Equal(t, m.cache.Len(), 1)
	m.runFunc()(2)
	assert.Equal(t, m.cache.Len(), 2)
}

func TestMemoizeNoArgs(t *testing.T) {
	i := 0
	noArgs := func() { i++ }
	memoizedFunc := Auto(noArgs)
	memoizedFunc()
	memoizedFunc()
	assert.Equal(t, 1, i)
}

func TestMemoizeStringArgs(t *testing.T) {
	i := 0
	stringArgs := func(s string) { i++ }
	memoizedFunc := Auto(stringArgs)
	memoizedFunc("hello")
	memoizedFunc("hello")
	assert.Equal(t, 1, i)
}

func TestMemoizeUintArgs(t *testing.T) {
	i := 0
	uintArgs := func(x uint) { i++ }
	memoizedFunc := Auto(uintArgs)
	memoizedFunc(1)
	memoizedFunc(1)
	assert.Equal(t, 1, i)
}

func TestMemoizeFloatArgs(t *testing.T) {
	i := 0
	floatArgs := func(x float64) { i++ }
	memoizedFunc := Auto(floatArgs)
	memoizedFunc(1.0)
	memoizedFunc(1.0)
	assert.Equal(t, 1, i)
}

func TestMemoizeIntArgs(t *testing.T) {
	i := 0
	intArgs := func(x int) { i++ }
	memoizedFunc := Auto(intArgs)
	memoizedFunc(1)
	memoizedFunc(1)
	assert.Equal(t, 1, i)
}

func TestMemoizeMultipleArgs(t *testing.T) {
	i := 0
	multipleArgs := func(x, y int) { i++ }
	memoizedFunc := Auto(multipleArgs)
	memoizedFunc(1, 2)
	memoizedFunc(1, 2)
	assert.Equal(t, 1, i)
}

func TestMemoizeMapArgs(t *testing.T) {
	i := 0
	mapArgs := func(m map[string]int) { i++ }
	memoizedFunc := Auto(mapArgs)
	memoizedFunc(map[string]int{"a": 1})
	memoizedFunc(map[string]int{"a": 1})
	assert.Equal(t, 1, i)
}

func TestMemoizeSliceArgs(t *testing.T) {
	i := 0
	sliceArgs := func(s []int) { i++ }
	memoizedFunc := Auto(sliceArgs)
	memoizedFunc([]int{1, 2})
	memoizedFunc([]int{1, 2})
	assert.Equal(t, 1, i)
}

// func TestMemoizeStructArgs(t *testing.T) {
// 	i := 0
// 	type s struct {
// 		x int
// 	}
// 	structArgs := func(s) { i++ }
// 	memoizedFunc := Memoize(structArgs)
// 	memoizedFunc(s{1})
// 	memoizedFunc(s{1})
// 	assert.Equal(t, 1, i)
// }

func TestMemoizeIntPointerArgs(t *testing.T) {
	i := 0
	intPointerArgs := func(x *int) { i++ }
	memoizedFunc := Auto(intPointerArgs)
	x := 1
	memoizedFunc(&x)
	memoizedFunc(&x)
	assert.Equal(t, 1, i)
}

func TestMemoizeDifferentIntPointerArgs(t *testing.T) {
	i := 0
	intPointerArgs := func(x *int) { i++ }
	memoizedFunc := Auto(intPointerArgs)
	x, y := 1, 1
	memoizedFunc(&x)
	memoizedFunc(&y)
	assert.Equal(t, 2, i)
}

// func TestMemoizeManyDifferentArgs(t *testing.T) {
// 	i := 0
// 	s := struct{ suen string }{suen: "suen"}
// 	manyDifferentArgs := func(x, y int, z string, m map[string]int, sl []int, st struct{ suen string }) { i++ }
// 	memoizedFunc := Memoize(manyDifferentArgs)
// 	memoizedFunc(1, 2, "hello", map[string]int{"a": 1}, []int{1, 2}, struct{ suen string }{suen: "suen"})
// 	memoizedFunc(1, 2, "hello", map[string]int{"a": 1}, []int{1, 2}, s)
// 	assert.Equal(t, 1, i)
// }

func TestCreateKeyer(t *testing.T) {
	type dummyStruct struct {
		x int
		y int
	}
	type dummyStructHalfExported struct {
		X int
		y [2]int
	}
	type dummyStructExported struct {
		X int
		Y [2]int
	}

	type dummyInterface interface{}

	dummyStructI1 := dummyStruct{x: 1, y: 2}
	dummyStructI2 := dummyStruct{x: 1, y: 23}

	type args struct{}
	tests := []struct {
		name          string
		args          args
		callWith      [][]reflect.Value
		wantKeysCount int
	}{
		{
			name: "no args",
			callWith: [][]reflect.Value{
				{},
				{},
			},
			wantKeysCount: 1,
		},
		{
			name: "all whole number args are compared by value",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(1)},
				{reflect.ValueOf(uint64(1))},
				{reflect.ValueOf(int8(1))},
				{reflect.ValueOf(int16(1))},
				{reflect.ValueOf(int32(1))},
				{reflect.ValueOf(int64(1))},
				{reflect.ValueOf(uint(1))},
				{reflect.ValueOf(uint8(1))},
				{reflect.ValueOf(uint16(1))},
				{reflect.ValueOf(uint32(1))},
				{reflect.ValueOf(uint64(1))},
			},
			wantKeysCount: 1,
		},
		{
			name: "all float args are compared by value",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(1.0)},
				{reflect.ValueOf(float32(1.0))},
				{reflect.ValueOf(float64(1.0))},
			},
			wantKeysCount: 1,
		},
		{
			name: "all string args are compared by value",
			callWith: [][]reflect.Value{
				{reflect.ValueOf("hello")},
				{reflect.ValueOf("hello")},
			},
			wantKeysCount: 1,
		},
		{
			name: "all complex args are compared by value",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(complex(1, 1))},
				{reflect.ValueOf(complex64(1 + 1i))},
				{reflect.ValueOf(complex128(1 + 1i))},
			},
			wantKeysCount: 1,
		},
		{
			name: "all slice args are compared by value",
			callWith: [][]reflect.Value{
				{reflect.ValueOf([]int{1, 2})},
				{reflect.ValueOf([]int{1, 2})},
			},
			wantKeysCount: 1,
		},
		{
			name: "all map args are compared by value",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(map[string]int{"a": 1})},
				{reflect.ValueOf(map[string]int{"a": 1})},
			},
			wantKeysCount: 1,
		},
		{
			name: "test single int arg",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(1)},
				{reflect.ValueOf(2)},
			},
			wantKeysCount: 2,
		},
		{
			name: "test single string arg",
			callWith: [][]reflect.Value{
				{reflect.ValueOf("hello")},
				{reflect.ValueOf("world")},
			},
			wantKeysCount: 2,
		},
		{
			name: "test pointers",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(ptr(0))},
				{reflect.ValueOf(ptr(0))},
				{reflect.ValueOf(uintptr(0))},
				{reflect.ValueOf(ptr(2))},
			},
			wantKeysCount: 4,
		},
		{
			name: "test nil",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(nil)},
				{reflect.ValueOf(nil)},
			},
			wantKeysCount: 1,
		},
		{
			name: "test variadic args 1",
			callWith: [][]reflect.Value{
				{reflect.ValueOf("a"), reflect.ValueOf("b")},
				{reflect.ValueOf("a")},
			},
			wantKeysCount: 2,
		},
		{
			name: "test variadic args 2",
			callWith: [][]reflect.Value{
				{reflect.ValueOf("hi"), reflect.ValueOf("there")},
				{reflect.ValueOf("hi"), reflect.ValueOf("there")},
				{reflect.ValueOf("hi")},
				{reflect.ValueOf("hi")},
			},
			wantKeysCount: 2,
		},
		// {
		// 	name: "test anonymous structs (completely ignored)",
		// 	args: args{compareSlices: false, compareMaps: false},
		// 	callWith: [][]reflect.Value{
		// 		{reflect.ValueOf(struct{ x, y int }{x: 1, y: 2})},
		// 		{reflect.ValueOf(struct{ x, y int }{y: 1, x: 2})},
		// 		{reflect.ValueOf(struct{ x, y int }{y: 1, x: 3})},
		// 	},
		// 	wantKeysCount: 1,
		// },
		{
			name: "test structs (no exported fields)",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(dummyStruct{x: 1, y: 2})},
				{reflect.ValueOf(dummyStruct{x: 1, y: 2})},
				{reflect.ValueOf(dummyStruct{x: 1, y: 3})},
			},
			wantKeysCount: 2,
		},
		{
			name: "test structs (half exported fields)",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(dummyStructHalfExported{X: 1})},
				{reflect.ValueOf(dummyStructHalfExported{X: 1})},
				{reflect.ValueOf(dummyStructHalfExported{X: 2})},
			},
			wantKeysCount: 2,
		},
		{
			name: "test structs (all exported fields)",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(dummyStructExported{X: 1, Y: [2]int{1, 2}})},
				{reflect.ValueOf(dummyStructExported{X: 1, Y: [2]int{1, 2}})},
				{reflect.ValueOf(dummyStructExported{X: 1, Y: [2]int{2, 1}})},
				{reflect.ValueOf(dummyStructExported{X: 1, Y: [2]int{2}})},
			},
			wantKeysCount: 3,
		},
		{
			name: "test context",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(context.Background())},
				{reflect.ValueOf(context.Background())},
				{reflect.ValueOf(context.WithValue(context.Background(), "a", 1))},
				{reflect.ValueOf(context.WithValue(context.Background(), "a", 1))},
			},
			wantKeysCount: 3,
		},
		{
			name: "test interface",
			callWith: [][]reflect.Value{
				{reflect.ValueOf(dummyStructI1)},
				{reflect.ValueOf(dummyStructI1)},
				{reflect.ValueOf(dummyStructI2)},
			},
			wantKeysCount: 2,
		},
	}
	for _, tt := range tests {
		cache := NewMapCache(Fn64aKeyer())
		t.Run(tt.name, func(t *testing.T) {
			for _, args := range tt.callWith {
				cache.Set(args, []reflect.Value{})
			}

			assert.Len(t, cache.cache, tt.wantKeysCount, fmt.Sprintf("%q failed", tt.name))
		})
	}
}

func ptr[T any](x T) *T {
	return &x
}
