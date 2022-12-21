package giter

func Map[T, TP any](f func(T) TP, iter Iterator[T]) Iterator[TP] {
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

func Filter[T any](f func(T) bool, iter Iterator[T]) Iterator[T] {
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

func FlatMap[T, R any](f func(T) []R, iter Iterator[T]) Iterator[R] {
	return Make(
		func(out chan<- R, stopChan <-chan interface{}) {
			defer iter.Close()
			for v := range iter.Each {
				mapped := f(v)

				for _, v := range mapped {
					select {
					case out <- v:
					case <-stopChan:
						return
					}
				}
			}
		})
}

func Chunk[T any](n int, iter Iterator[T]) Iterator[[]T] {
	return Make(
		func(out chan<- []T, stopChan <-chan interface{}) {
			defer iter.Close()

			buf := make([]T, 0, n)

			// returns false if signaled that we need to stop and bail out
			flush := func() bool {
				if len(buf) == 0 {
					return true
				}

				outs := make([]T, len(buf))
				copy(outs, buf)

				select {
				case out <- outs:
				case <-stopChan:
					return false
				}

				buf = buf[:0]

				return true
			}

			for v := range iter.Each {
				buf = buf[:len(buf)+1]
				buf[len(buf)-1] = v

				if len(buf) == cap(buf) {
					if !flush() {
						return
					}
				}
			}

			// flush anything remaining; ignore if we flushed everything since we're
			// going to return afterwards regardless.
			_ = flush()
		})
}
