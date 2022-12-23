// Package giter implements lazy, type-safe iterators in Go.
//
// Iterators are implemented via a struct Iterator, which exposes a channel that emits the
// iterator's values.
//
// Make creates Iterators via a provided implementation function that produces the iterator's
// values.
//
// These iterators can be effectively created and used with just these two elements; however, many
// common use-cases are supported by more convenient functions that produce iterators from different
// values, transform the contents of iterators to produce new iterators, compose multiple iterators
// into a single iterator, or produce a non-iterator value from an iterator.
package giter

// An Iterator that produces 0 or more values
//
// To consume the iterator, range-loop over the Each channel.
//
// An iterator must be closed, otherwise any resources held by the iterator's producer will not be
// released. It may be closed repeatedly without ill effect.
//
// Thus a trivial iterator use that consumes all elements and discards them might look like:
//
//	iter := ...
//	defer iter.Close()
//	for _ := range iter.Each {
//	}
type Iterator[T any] struct {
	Each <-chan T

	// stopChan is used to coordinate stopping of the Each producer goroutine: Close() sends
	// a message on the channel, which the Each producer takes to mean it should stop producing
	// and exit.
	stopChan chan<- interface{}
}

// Close stops production to Each and releases goroutines and any other resources held for producing
// to this Iterator.
func (iter *Iterator[T]) Close() {
	select {
	case iter.stopChan <- nil:
	default:
	}
}

// Make creates an Iterator via a given function that produces values.
//
// This is the most basic interface to creating Iterators. There are many functions already provided
// in this package that create iterators from various common inputs that may preclude the need to
// use this function.
//
// The function receives a function that is given two channels: one to which the function must
// produce the values of the iterator, and the other which signals that no more values should be
// produced and the implementation should return.
//
// The resulting iterator will produce values in the order that they are produced by the passed
// implementation function, in the order they were produced in.
func Make[T any](impl func(chan<- T, <-chan interface{})) (i Iterator[T]) {
	values := make(chan T)

	stopChan := make(chan interface{})

	go func() {
		impl(values, stopChan)
		close(values)
	}()

	return Iterator[T]{
		Each:     values,
		stopChan: stopChan,
	}
}

// Slice creates an iterator that emits the values of a given slice.
//
// Modification of the provided slice can impact the values produced, so caution is advised.
//
// The given slice will be held until all values are consumed or the iterator is closed.
func Slice[T any](xs []T) (i Iterator[T]) {
	return Make(
		func(values chan<- T, stopChan <-chan interface{}) {
			for _, x := range xs {
				select {
				case values <- x:
				case <-stopChan:
					return
				}
			}
		})
}

// NeverShrink can be used to indicate to ConsumeSlice to never shrink the slice it's consuming.
func NeverShrink(_, _ int) bool { return false }

// ConsumeSlice produces an iterator emitting the values of a slice, consuming the values of the
// slice as they are emitted.
//
// Values in the source slice are overwritten with zero when they are emitted, so any pointers
// within are released.
//
// The slice can be optionally shrunk (a smaller slice is allocated, the remaining elements are
// copied into it and the old, larger slice is released) if the source slice itself needs to be
// released before all values are emitted. shrink receives the current len (remaining) and cap of
// the slice and can return true to trigger shrinking.
//
// A convenience shrink indicator function NeverShrink is provided to disable shrinking.
func ConsumeSlice[T any](shrink func(l, c int) bool, xs []T) (i Iterator[T]) {
	return Make(
		func(values chan<- T, stopChan <-chan interface{}) {
			var zero T
			defer clear(&xs)

		READ_XS:
			for {
				for i, x := range xs {
					xs[i] = zero // zero out any pointers

					select {
					case values <- x:
					case <-stopChan:
						return
					}

					// don't bother making an slice of cap 0 in the last
					// iteration.
					if i < len(xs)-1 && shrink(len(xs)-i-1, cap(xs)) {
						xs = xs[i+1:]

						prime := make([]T, len(xs))

						copy(prime, xs)

						clear(&xs)

						xs = prime

						// need to restart range so that we release xs and
						// our indexing makes sense.
						continue READ_XS
					}
				}

				// we made it to the end of the range, thus we consumed everything.
				break
			}
		})
}

// MapKeys returns an iterator that emits the keys of a given map.
func MapKeys[K comparable, V any](m map[K]V) (i Iterator[K]) {
	return Make(
		func(keys chan<- K, stopChan <-chan interface{}) {
			for k, _ := range m {
				select {
				case keys <- k:
				case <-stopChan:
					return
				}
			}
		})
}

func MapValues[K comparable, V any](m map[K]V) (i Iterator[V]) {
	return Make(
		func(values chan<- V, stopChan <-chan interface{}) {
			for _, v := range m {
				select {
				case values <- v:
				case <-stopChan:
					return
				}
			}
		})
}

type KVPair[K comparable, V any] struct {
	Key   K
	Value V
}

func MapPairs[K comparable, V any](m map[K]V) (i Iterator[KVPair[K, V]]) {
	return Make(
		func(keys chan<- KVPair[K, V], stopChan <-chan interface{}) {
		LOOP:
			for k, v := range m {
				select {
				case keys <- KVPair[K, V]{k, v}:
				case <-stopChan:
					break LOOP
				}
			}
		})
}

func One[V any](x V) Iterator[V] {
	return Make(
		func(values chan<- V, stopChan <-chan interface{}) {
			select {
			case values <- x:
			case <-stopChan:
				return
			}
		})
}
