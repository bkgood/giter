package giter

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	want := make([]int, 0, len(xs))

	f := func(x int) int { return 2 * x }
	for _, x := range xs {
		want = append(want, f(x))
	}

	mapped := Map(f, Slice(xs))
	defer mapped.Close()
	out := []int{}
	for x := range mapped.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestMap: Map(2*x, xs) = %v, want = %v", out, want)
	}
}

func TestFilter(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	want := make([]int, 0, len(xs))

	f := func(x int) bool { return x%2 == 0 }
	for _, x := range xs {
		if f(x) {
			want = append(want, x)
		}
	}

	mapped := Filter(f, Slice(xs))
	defer mapped.Close()
	out := []int{}
	for x := range mapped.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestFilter: Map(!x%%2, xs) = %v, want = %v", out, want)
	}
}

func TestFlatMap(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	want := make([]int, 0, len(xs))

	f := func(x int) []int { return []int{x, x / 2} }
	for _, x := range xs {
		want = append(want, f(x)...)
	}

	out := ToSlice(FlatMap(f, Slice(xs)))

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestFlatMap: FlatMap(x -> [ x, x / 2 ], xs) = %v, want = %v", out, want)
	}
}

func TestChunk(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	want := [][]int{
		[]int{xs[0], xs[1]},
		[]int{xs[2], xs[3]},
		[]int{xs[4]},
	}

	out := ToSlice(Chunk(2, Slice(xs)))

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestChunk: Chunk(2, {1..=5}) = %v, want = %v", out, want)
	}
}
