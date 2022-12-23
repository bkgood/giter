package giter

import (
	"reflect"
	"sort"
	"testing"
)

func TestChained(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	want := make([]int, 0, len(xs))

	f := func(x int) bool { return x%2 == 0 }
	g := func(x int) int { return 2 * x }

	for _, x := range xs {
		if f(x) {
			want = append(want, g(x))
		}
	}

	mapped := Map(g, Filter(f, Slice(xs)))

	defer mapped.Close()
	out := []int{}
	for x := range mapped.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestChained: Map(2*x, Filter(!x%%2, xs)) = %v, want = %v", out, want)
	}
}

func TestZipp(t *testing.T) {
	slices := [][]int{
		[]int{1, 3, 5},
		[]int{2, 4, 6},
	}

	want := make([]int, len(slices[0])*2)

	for i := range want {
		want[i] = slices[i%2][i/2]
	}

	zipped := Zip(Slice(slices[0]), Slice(slices[1]))

	defer zipped.Close()
	out := []int{}
	for x := range zipped.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestZip: Zip(odds, evens) = %v, want = %v", out, want)
	}
}

func slices[T any](slices ...[]T) []Iterator[T] {
	// XXX might add this to the api, it's nice
	return ToSlice(Map(func(xs []T) Iterator[T] { return Slice(xs) }, Slice(slices)))
}

func contains[T comparable](xs []T, x T) bool {
	// how does this exist
	for _, y := range xs {
		if y == x {
			return true
		}
	}

	return false
}

func TestMerge(t *testing.T) {
	in := [][]int{
		[]int{1, 2},
		[]int{3, 4, 5},
	}

	want := []int{1, 2, 3, 4, 5}

	out := ToSlice(Merge(slices(in...)...))

	sort.Ints(out)

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestMerge: Merge({1,2}, {3,4,5}) (unordered) = %v, want = %v", out, want)
	}

	// vaguely verify that closing doesn't bug out.
	// XXX maybe there is some decent way to make sure no goroutines that we spawned linger
	// after the Close()? that doesn't involve getting super invasive.
	again := Merge(slices(in...)...)

	oneOut := <-again.Each
	again.Close()

	if !contains(want, oneOut) {
		t.Errorf(
			"TestMerge: Merge({1,2}, {3,4,5}) one entry = %v, want one of = %v",
			oneOut, want)
	}
}

func TestConcat(t *testing.T) {
	slices := [][]int{
		[]int{1, 3, 5},
		[]int{2, 4, 6},
	}

	want := make([]int, len(slices[0])*2)

	for i := range want {
		slice := slices[0]
		j := i

		if i >= len(slice) {
			slice = slices[1]
			j -= len(slice)
		}

		want[i] = slice[j]
	}

	concatted := Concat(Slice(slices[0]), Slice(slices[1]))

	defer concatted.Close()
	out := []int{}
	for x := range concatted.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestConcat: Concat(odds, evens) = %v, want = %v", out, want)
	}
}
