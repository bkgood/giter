# giter

![did i break the build/tests?](https://github.com/bkgood/giter/actions/workflows/ci.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/bkgood/giter.svg)](https://pkg.go.dev/github.com/bkgood/giter)

i wanted to pass a function an iterator over a map's keys and now I have this.

there is a non-zero chance that:

1. someone else has already done this, or
1. this is very non-idiomatic Go.

I'm ok with both.

this whole exercise is somewhat poisoned by how verbose and without-inference go's function literal
syntax is ╮ (. ❛ ᴗ ❛.) ╭

## usage

```go
package main

import (
	"fmt"

	i "github.com/bkgood/giter"
)

func main() {
	// sum the even numbers in [0, 100)
	x := i.Sum(i.Filter(func(x int) bool { return x%2 == 0 }, i.Range[int](0, 100)))

	// 2450
	fmt.Println(x)

	// build an index of some structs by an id
	type foo struct {
		id   int
		name string
	}

	index := i.Collect( // or just ToMap
		i.MapCollector[int, foo](),
		i.Map(
			func(f foo) i.KVPair[int, foo] {
				return i.KVPair[int, foo]{f.id, f}
			},
			i.Slice([]foo{foo{1, "willy"}})))

	// map[1:{1 willy}]
	fmt.Println(index)
}
```

## things i vaguely want to do

### add docs

i'm still getting used to writing godoc stuff; this would be decent practice.

### sort out if I can have sized iterators next to unsized ones

if i can see a size, I can preallocate in collection.

### other functions

#### reverse ranges

ranges only ascend now. they should be able to descend.

#### reverse slices

90% of the time i want a reverse range i really just want to traverse a slice backwards. so just
offer that.

## isn't this slow?

probably? i wouldn't try to write a blas implementation with it. but for typical webby
line-of-business type stuff, it is predictably slow, and behaves the same for large and small n.

but if someone did write a blas implementation with this i would be amused.
