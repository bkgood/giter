package giter

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
				// XXX is this bugged? is it always closing the final struct pointed
				// to by iter?
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
