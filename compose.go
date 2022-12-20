package giter

// XXX it would be nice to make these methods but you can't have type params on methods.

func Map[T, TP any](iter Iterator[T], f func(T) TP) Iterator[TP] {
	return Make(
		func(out chan<- TP, stopChan <-chan interface{}) {
			defer iter.Close()
			for v := range iter.Each {
				select {
				case out <- f(v):
				case <-stopChan:
					return
				}
			}
		})
}

func Filter[T any](iter Iterator[T], f func(T) bool) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			defer iter.Close()
			for v := range iter.Each {
				if f(v) {
					select {
					case out <- v:
					case <-stopChan:
						return
					}
				}
			}
		})
}

// Zip takes n iterators and gives n elements, one from each, until one iterator stops.
// If the iterators give a different number of results from the given iterators, unless it is told
// to stop prematurely.
func Zip[T any](iters ...Iterator[T]) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			for _, i := range iters {
				defer i.Close()
			}

			buf := make([]T, 0, len(iters))

			for {
				for len(buf) < len(iters) {
					select {
					case x, ok := <-iters[len(buf)].Each:
						if !ok {
							return
						}
						buf = append(buf, x)
					case <-stopChan:
						return
					}
				}

				for _, x := range buf {
					select {
					case out <- x:
					case <-stopChan:
						return
					}
				}

				buf = buf[:0]
			}
		})
}

// undefined output order
func Merge[T any](iters ...Iterator[T]) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			// TODO
		})
}

func Concat[T any](iters ...Iterator[T]) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			for _, iter := range iters {
				defer iter.Close()
			}

			for _, iter := range iters {
				for x := range iter.Each {
					select {
					case out <- x:
					case <-stopChan:
						return
					}
				}
			}
		})
}

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
