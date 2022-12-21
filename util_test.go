package giter

import (
	"reflect"
	"testing"
)

func TestClear(t *testing.T) {
	type foo struct {
		y int
	}

	xs := []foo{foo{5}}

	want := []foo{foo{}}

	clear(&xs)

	if len(xs) > 0 {
		t.Errorf("TestClear: clear(xs) len(xs) = %v, want %v", len(xs), 0)
	}

	// extend len so DeepEqual can look at what we should have zeroed
	xs = xs[:cap(xs)]

	if !reflect.DeepEqual(xs, want) {
		t.Errorf("TestClear: clear(xs) = %v, want %v", xs, want)
	}
}
