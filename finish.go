package giter

func ToSlice[T any](iter Iterator[T]) []T {
	return Collect(SliceCollector[T](), iter)
}

func ToMap[K comparable, V any](iter Iterator[KVPair[K, V]]) map[K]V {
	return Collect(MapCollector[K, V](), iter)
}

type Collector[T, R any] func(<-chan T) R

func MapCollector[K comparable, V any]() Collector[KVPair[K, V], map[K]V] {
	return func(each <-chan KVPair[K, V]) map[K]V {
		out := map[K]V{}

		for x := range each {
			out[x.Key] = x.Value
		}
		return out
	}
}

func SliceCollector[V any]() Collector[V, []V] {
	return func(each <-chan V) []V {
		out := []V{}

		for x := range each {
			out = append(out, x)
		}
		return out
	}
}

func Collect[T, R any](collector Collector[T, R], iter Iterator[T]) R {
	defer iter.Close()
	return collector(iter.Each)
}

func Fold[T, R any](initial R, f func(next T, current R) R, iter Iterator[T]) R {
	defer iter.Close()

	for x := range iter.Each {
		initial = f(x, initial)
	}

	return initial
}

func First[T any](iter Iterator[T]) *T {
	defer iter.Close()

	for x := range iter.Each {
		return &x
	}

	return nil
}

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

func Any[T any](pred func(T) bool, iter Iterator[T]) bool {
	defer iter.Close()

	for x := range iter.Each {
		if pred(x) {
			return true
		}
	}

	return false
}
