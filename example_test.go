package giter

import "testing"

func TestExample(_ *testing.T) {
	// sum the even numbers in [0, 100)
	_ = Sum(Filter(func(x int) bool { return x%2 == 0 }, Range[int](0, 100)))

	// build an index of some structs by an id
	type foo struct {
		id   int
		name string
	}

	_ = Collect(
		MapCollector[int, foo](),
		Map(
			func(f foo) KVPair[int, foo] {
				return KVPair[int, foo]{f.id, f}
			}, Slice([]foo{foo{1, "hi"}})))
}
