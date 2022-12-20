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

	mapped := Map(Slice(xs), f)
	defer mapped.Close()
	out := []int{}
	for x := range mapped.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestMap: Map(xs, 2*x) = %v, want = %v", out, want)
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

	mapped := Filter(Slice(xs), f)
	defer mapped.Close()
	out := []int{}
	for x := range mapped.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestFilter: Map(xs, !x%%2) = %v, want = %v", out, want)
	}
}

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

	mapped := Map(Filter(Slice(xs), f), g)

	defer mapped.Close()
	out := []int{}
	for x := range mapped.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestChained: Map(Filter(xs, !x%%2), 2*x) = %v, want = %v", out, want)
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

func TestToSlice(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := make([]int, len(xs))
	copy(want, xs)

	out := ToSlice(Slice(xs))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestToSlice: out = %v, want %v", out, want)
	}
}

func TestToMap(t *testing.T) {
	want, _, _, pairs := testMap()

	out := ToMap(Slice(pairs))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestToMap: out = %v, want %v", out, want)
	}
}

func TestCollectToMap(t *testing.T) {
	want, _, _, pairs := testMap()

	out := Collect(Slice(pairs), MapCollector[string, int]())

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestCollectToMap: out = %v, want %v", out, want)
	}
}

func TestCollectToSlice(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := make([]int, len(xs))
	copy(want, xs)

	out := Collect(Slice(xs), SliceCollector[int]())

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestCollectToSlice: out = %v, want %v", out, want)
	}
}
