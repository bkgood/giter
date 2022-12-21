package giter

func ToSlice[T any](iter Iterator[T]) []T {
	return Collect(iter, SliceCollector[T]())
}

func ToMap[K comparable, V any](iter Iterator[KVPair[K, V]]) map[K]V {
	return Collect(iter, MapCollector[K, V]())
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

func Collect[T, R any](iter Iterator[T], collector Collector[T, R]) R {
	defer iter.Close()
	return collector(iter.Each)
}

func Fold[T, R any](iter Iterator[T], initial R, f func(next T, current R) R) R {
	defer iter.Close()

	for x := range iter.Each {
		initial = f(x, initial)
	}

	return initial
}
