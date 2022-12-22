package giter

import (
	"reflect"
	"sort"
	"testing"
)

func TestSlice(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := make([]int, len(xs))
	copy(want, xs)

	iter := Slice(xs)

	out := []int{}

	defer iter.Close()
	for x := range iter.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestSlice: out = %v, want %v", out, want)
	}
}

func TestConsumeSlice(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}

	want := make([]int, len(xs))
	copy(want, xs)

	// I don't have any way to confirm shrinking is happening atm; I can only tell it to shrink
	// and verify things don't bug out on our example input.
	iter := ConsumeSlice(func(l, c int) bool { return true }, xs)

	out := []int{}

	defer iter.Close()
	for x := range iter.Each {
		out = append(out, x)
	}

	for i, x := range xs {
		if x != 0 {
			// we should zero every element we emit
			t.Errorf("TestConsumeSlice: xs[%v] = %v, want %v sigh %v", i, x, 0, xs)
		}
	}

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestConsumeSlice: out = %v, want %v", out, want)
	}
}

func TestStop(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	want := []int{xs[0]}

	iter := Slice(xs)

	out := []int{}

	defer iter.Close()
	for x := range iter.Each {
		out = append(out, x)
		iter.Close()
	}

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestStop: out = %v, want %v", out, want)
	}
}

func testMap() (m map[string]int, keys []string, values []int, pairs []KVPair[string, int]) {
	keys = []string{"foo", "bar", "baz"}
	values = []int{}
	pairs = []KVPair[string, int]{}

	m = make(map[string]int, len(keys))

	for i, x := range keys {
		m[x] = i

		values = append(values, i)
		pairs = append(pairs, KVPair[string, int]{x, i})
	}

	return m, keys, values, pairs
}

func TestMapKeys(t *testing.T) {
	m, want, _, _ := testMap()

	iter := MapKeys(m)

	defer iter.Close()

	out := []string{}

	for x := range iter.Each {
		out = append(out, x)
	}

	sort.Strings(out)
	sort.Strings(want)

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestMapKeys: out = %v, want %v", out, want)
	}
}

func TestMapValues(t *testing.T) {
	m, _, want, _ := testMap()

	iter := MapValues(m)

	defer iter.Close()

	out := []int{}

	for x := range iter.Each {
		out = append(out, x)
	}

	sort.Ints(out)
	sort.Ints(want)

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestMapValues: out = %v, want %v", out, want)
	}
}

func TestMapPairs(t *testing.T) {
	m, _, _, want := testMap()

	iter := MapPairs(m)

	defer iter.Close()

	out := []KVPair[string, int]{}

	for x := range iter.Each {
		out = append(out, x)
	}

	less := func(xs []KVPair[string, int]) func(int, int) bool {
		return func(a, b int) bool {
			i, j := xs[a], xs[b]

			if i.Key != j.Key {
				return i.Key < j.Key
			}

			return i.Value < j.Value
		}
	}

	sort.Slice(out, less(out))
	sort.Slice(want, less(want))

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestMapPairs: out = %v, want %v", out, want)
	}
}

func TestOne(t *testing.T) {
	x := 1337
	want := []int{x}

	iter := One(x)
	defer iter.Close()

	out := []int{}

	for x := range iter.Each {
		out = append(out, x)
	}

	if !reflect.DeepEqual(want, out) {
		t.Errorf("TestOne: out = %v, want %v", out, want)
	}
}
