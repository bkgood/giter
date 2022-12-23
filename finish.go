package giter

// ToSlice consumes an iterator and returns the values in a slice.
func ToSlice[T any](iter Iterator[T]) []T {
	return Collect(SliceCollector[T](), iter)
}

// ToMap consumes an iterator of KVPair key-value pairs and returns a map.
func ToMap[K comparable, V any](iter Iterator[KVPair[K, V]]) map[K]V {
	return Collect(MapCollector[K, V](), iter)
}

// A Collector consumes the values of an Iterator and returns some aggregated value.
type Collector[T, R any] func(<-chan T) R

// MapCollector returns a Collectors that creates a map from an Iterator of KVPairs.
func MapCollector[K comparable, V any]() Collector[KVPair[K, V], map[K]V] {
	return func(each <-chan KVPair[K, V]) map[K]V {
		out := map[K]V{}

		for x := range each {
			out[x.Key] = x.Value
		}
		return out
	}
}

// SliceCollector returns a Collector that creates a slice from an Iterator's values.
func SliceCollector[V any]() Collector[V, []V] {
	return func(each <-chan V) []V {
		out := []V{}

		for x := range each {
			out = append(out, x)
		}
		return out
	}
}

// Collect creates a value resulting from consuming an Iterator's values via a Collector.
func Collect[T, R any](collector Collector[T, R], iter Iterator[T]) R {
	defer iter.Close()
	return collector(iter.Each)
}

// Fold returns the value resulting from calling a given function with an initial value and each
// value emitted by the iterator, updating the initial value with each invocation.
func Fold[T, R any](initial R, f func(next T, current R) R, iter Iterator[T]) R {
	defer iter.Close()

	for x := range iter.Each {
		initial = f(x, initial)
	}

	return initial
}

// First returns the first value emitted by an Iterator, if any.
func First[T any](iter Iterator[T]) *T {
	defer iter.Close()

	for x := range iter.Each {
		return &x
	}

	return nil
}

// Last returns the last value emitted by an Iterator, if any.
func Last[T any](iter Iterator[T]) *T {
	defer iter.Close()

	var out *T
	for x := range iter.Each {
		out = &x
	}

	return out
}

// i should choose one of these :/

func Some[T any](pred func(T) bool, iter Iterator[T]) bool {
	return Any(pred, iter)
}

func Has[T any](pred func(T) bool, iter Iterator[T]) bool {
	return Any(pred, iter)
}

// Any returns true if some value emitted by a given Iterator matches a given predicate.
func Any[T any](pred func(T) bool, iter Iterator[T]) bool {
	defer iter.Close()

	for x := range iter.Each {
		if pred(x) {
			return true
		}
	}

	return false
}
