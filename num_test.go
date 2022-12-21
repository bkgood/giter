package giter

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := 0

	for _, x := range xs {
		want += x
	}

	out := Sum(Slice(xs))

	if want != out {
		t.Errorf("TestSum: out = %v, want %v", out, want)
	}
}

func TestProd(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := 1

	for _, x := range xs {
		want *= x
	}

	out := Prod(Slice(xs))

	if want != out {
		t.Errorf("TestProd: out = %v, want %v", out, want)
	}
}

func TestRange(t *testing.T) {
	want := []int{1, 2, 3, 4, 5}

	out := ToSlice(Range(1, 6))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestRange: out = %v, want %v", out, want)
	}
}

func TestFpRange(t *testing.T) {
	want := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	out := ToSlice(Range[float64](1, 6))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestFpRange: out = %v, want %v", out, want)
	}
}

func TestRangeBy(t *testing.T) {
	want := []int{10, 20, 30, 40, 50}

	out := ToSlice(RangeBy(10, 55, 10))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestRangeBy: out = %v, want %v", out, want)
	}
}
