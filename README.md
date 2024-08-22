# memoize

> **WARNING**: This project is in early development and is not yet ready for use.

Memoize is a library for memorizing function return values based on function arguments.
The main use case is to avoid making duplicate expensive IO operations or computations multiple times.

The most basic usage looks like this.

```go
package main

import (
	"fmt"

	"github.com/alx99/memoize"
)

func expensiveFunc(val int) int {
	fmt.Println("Expensive function executed")
	return val
}

func main() {
	memoized := memoize.Auto(expensiveFunc)
	memoized(1) // expensiveFunc is called
	memoized(2) // expensiveFunc is called
	memoized(1) // expensiveFunc is not called

	manualMemoized := memoize.Manual[func(int) int, string](expensiveFunc)
	manualMemoized("key1")(1) // expensiveFunc is called
	manualMemoized("key2")(2) // expensiveFunc is called
	manualMemoized("key1")(2) // expensiveFunc is not called and 1 is returned
}
```

Please read the [documentation](https://pkg.go.dev/github.com/alx99/memoize) for more information.
