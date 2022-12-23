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

// Merge returns an iterator emitting all the values of the given iterators
//
// The output order is undefined.
func Merge[T any](iters ...Iterator[T]) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			// no way to mux reading from n channels, so...

			// we launch 0..len(iters) goroutine, each watches a stop channel to know
			// when to bail out
			stops := make([]chan interface{}, len(iters))

			// after we spawn our goroutines, we wait to see len(iters) messages on done
			done := make(chan interface{})

			for i := range iters {
				go func(iter Iterator[T], stop <-chan interface{}) {
					defer iter.Close()

					for {
						select {
						case x, ok := <-iter.Each:
							if !ok {
								done <- nil
								return
							}
							out <- x
						case <-stop:
							return
						}
					}
				}(iters[i], stops[i])
			}

			waiting := len(iters)
			for {
				select {
				case <-done:
					waiting--

					if waiting == 0 {
						return
					}
				case <-stopChan:
					for _, stop := range stops {
						stop <- nil
					}
				}
			}
		})
}

// Concat emits all the values of each the given iterators, one iterator after another (i.e. first
// the elements of the first one, then the second, and so on).
func Concat[T any](iters ...Iterator[T]) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			for _, iter := range iters {
				iter := iter // sigh
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
