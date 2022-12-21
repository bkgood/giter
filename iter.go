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
	Each  <-chan T
	Close func()
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

	stop, stopChan := func() (func(), <-chan interface{}) {
		bidiStopChan := make(chan interface{})
		stopChan := bidiStopChan

		stop := func() {
			select {
			case bidiStopChan <- nil:
			default:
			}
		}

		return stop, stopChan
	}()

	go func() {
		impl(values, stopChan)
		close(values)
	}()

	return Iterator[T]{Each: values, Close: stop}
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
