package memoize

import (
	"encoding"
	"math"
	"reflect"
	"sync"
	"unique"
	"unsafe"

	fnv "github.com/alx99/memoize/internal"
)

type KeyerFunc[T comparable] func(args []reflect.Value) T

func Fn64aKeyer() KeyerFunc[uint64] {
	var (
		keepAlive []unique.Handle[any]
		mu        sync.Mutex
		zeroArg   reflect.Value
	)

	return func(args []reflect.Value) uint64 {
		fnv := fnv.New64a()
		for _, arg := range args {
			fnv.WriteByte(1)
			switch arg.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fnv.WriteUint64(uint64(arg.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				fnv.WriteUint64(arg.Uint())
			case reflect.Float32, reflect.Float64:
				fnv.WriteUint64(math.Float64bits(arg.Float()))
			case reflect.Complex64, reflect.Complex128:
				v := arg.Complex()
				fnv.WriteUint64(math.Float64bits(real(v)))
				fnv.WriteUint64(math.Float64bits(imag(v)))
			case reflect.String:
				fnv.Write([]byte(arg.String()))
			case reflect.Bool:
				if arg.Bool() {
					fnv.WriteByte(1)
				} else {
					fnv.WriteByte(0)
				}
			case reflect.Pointer, reflect.UnsafePointer:
				fnv.WriteUint64(uint64(arg.Pointer()))
			default:
				if arg != zeroArg && arg.CanInterface() {
					v, ok := arg.Interface().(encoding.BinaryMarshaler)
					if ok {
						b, err := v.MarshalBinary()
						if err != nil {
							continue
						}
						fnv.Write(b)
						continue
					}
				}

				if arg.Comparable() {
					// The unique package is documented as being able to create a globally unique handle for any comparable type.
					// This handle is should always be the same for values that are equal.
					uniq := unique.Make(arg.Interface())
					fnv.Write(unsafe.Slice((*byte)(unsafe.Pointer(&uniq)), unsafe.Sizeof(uniq)))

					mu.Lock()
					keepAlive = append(keepAlive, uniq) // keep the handle alive
					mu.Unlock()
				}
			}
		}

		return fnv.Sum()
	}
}
