# giter
i wanted to pass a function an iterator over a map's keys and now I have this.

there is a non-zero chance that:

1. someone else has already done this, or
1. this is very non-idiomatic Go.

I'm ok with both.

this whole exercise is somewhat poisoned by how verbose and without-inference go's function literal
syntax is ╮ (. ❛ ᴗ ❛.) ╭

## usage

```go
import (
	i "github.com/bkgood/giter"
)

	// sum the even numbers in [0, 100)
	_ = i.Sum(i.Filter(func(x int) bool { return x%2 == 0 }, i.Range[int](0, 100)))

	// build an index of some structs by an id
	type foo struct {
		id   int
		name string
	}

	_ = i.Collect(
		i.MapCollector[int, foo](),
		i.Map(
			func(f foo) i.KVPair[int, foo] {
				return i.KVPair[int, foo]{f.id, f}
			}, i.Slice([]foo{foo{1, "hi"}})))
```

## things i vaguely want to do

### add docs

i'm still getting used to writing godoc stuff; this would be decent practice.

### sort out if I can have sized iterators next to unsized ones

if i can see a size, I can preallocate in collection.

### other functions

#### natatime
trivially process n things from the iterator at a time.

there is a variety of this that does

    Iterator[T] -> Iterator[[]T]

for some given output size, and then one can `FlatMap` it back into `Iterator[R any]`, but this does
allocate a bunch of (probably small) extra slices that could be trivially avoided with something
like a `MapChunked[T, R any](uint, f([]T) -> []R` with the sacrifice of a bit more code.

#### some/has
"has anything matching predicate"

trivially implemented with filter+first and a nil check but i don't really want to do that.
