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
	out := []T{}

	defer iter.Close()
	for x := range iter.Each {
		out = append(out, x)
	}

	return out
}

func ToMap[K comparable, V any](iter Iterator[KVPair[K, V]]) map[K]V {
	out := map[K]V{}

	defer iter.Close()
	for x := range iter.Each {
		out[x.Key] = x.Value
	}

	return out
}
