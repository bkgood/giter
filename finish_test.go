package giter

import (
	"reflect"
	"testing"
)

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

func TestScalarFold(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := 0

	for _, x := range xs {
		want += x
	}

	out := Fold(Slice(xs), 0, func(x, cur int) int { return cur + x })

	if want != out {
		t.Errorf("TestScalarFold: out = %v, want %v", out, want)
	}
}

func TestSliceFold(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := make([]int, len(xs))
	copy(want, xs)

	out := Fold(Slice(xs), []int{}, func(x int, cur []int) []int { return append(cur, x) })

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestSliceFold: out = %v, want %v", out, want)
	}
}
