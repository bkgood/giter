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

// ChunkedFlatMap maps n elements of a given iterator at a time into a new iterator.
// The mapping function receives input and output slices, and is expected to return 0 or more values
// via the output slice.
// The output slice has len 0 and cap >= n.
// The resulting iterator produces the elements returned by the mapping function in the order they
// were produced.
// The contents of the input and output slices will be cleared at least as often as every chunk is
// emitted into the output iterator, to avoid retaining excessive heap space.
func ChunkedFlatMap[T, R any](n int, f func([]T, []R) []R, iter Iterator[T]) Iterator[R] {
	// XXX maybe this should just be MapChunked? don't know that explicitly calling it FlatMap
	// is necessary, since it's obvious that we're mapping one chunk at a time and we just
	// receive a slice from the mapping function because we have to receive some sort of type
	// capable of holding multiple values.
	return Make(
		func(out chan<- R, stopChan <-chan interface{}) {
			defer iter.Close()

			buf := make([]T, 0, n)
			outBuf := make([]R, 0, cap(buf))

			// returns false if signaled that we need to stop and bail out
			flush := func() bool {
				// once we leave this function we are done with the contents of
				// these buffers and their contents (and anything it points to)
				// should be allowed to go away.
				defer clear(&buf)
				defer clear(&outBuf)

				if len(buf) == 0 {
					return true
				}

				outs := f(buf, outBuf)
				clear(&buf)
				var zero R

				for i, x := range outs {
					select {
					case out <- x:
						outs[i] = zero
					case <-stopChan:
						return false
					}
				}

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
