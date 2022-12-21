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

	out := Collect(MapCollector[string, int](), Slice(pairs))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestCollectToMap: out = %v, want %v", out, want)
	}
}

func TestCollectToSlice(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := make([]int, len(xs))
	copy(want, xs)

	out := Collect(SliceCollector[int](), Slice(xs))

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

	out := Fold(0, func(x, cur int) int { return cur + x }, Slice(xs))

	if want != out {
		t.Errorf("TestScalarFold: out = %v, want %v", out, want)
	}
}

func TestSliceFold(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := make([]int, len(xs))
	copy(want, xs)

	out := Fold([]int{}, func(x int, cur []int) []int { return append(cur, x) }, Slice(xs))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestSliceFold: out = %v, want %v", out, want)
	}
}

func TestFirst(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	var want *int

	f := func(x int) bool { return x%2 == 0 }
	for _, x := range xs {
		if f(x) {
			want = &x
			break
		}
	}

	out := First(Filter(f, Slice(xs)))

	if !reflect.DeepEqual(out, want) {
		t.Errorf("TestFirst: First(Filter(xs, x %% 2 == 0)) = %v, want = %v", out, want)
	}
}

func ptrTargetEquals[T comparable](x, y *T) bool {
	if (x == nil) != (y == nil) {
		return false
	}

	return x == nil || *x == *y
}

func TestLast(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	var want *int

	f := func(x int) bool { return x%2 == 0 }
	for _, x := range xs {
		x := x
		if f(x) {
			want = &x
		}
	}

	out := Last(Filter(f, Slice(xs)))

	if !ptrTargetEquals(out, want) {
		nilify := func(xs ...*int) []interface{} {
			out := make([]interface{}, len(xs))
			for i, x := range xs {
				if x == nil {
					out[i] = "nil"
				} else {
					out[i] = *x
				}
			}
			return out
		}

		t.Errorf("TestLast: Last(Filter(xs, x %% 2 == 0)) = %v, want = %v", nilify(out, want)...)
	}
}

func TestAny(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	smol := func(x int) bool { return x < 10 }
	hueg := func(x int) bool { return x > 9000 }

	if !Any(smol, Slice(xs)) {
		t.Errorf("TestAny: Any(smol, xs) is false, should be true")
	}

	if Any(hueg, Slice(xs)) {
		t.Errorf("TestAny: Any(hueg, xs) is true, should be false")
	}
}
