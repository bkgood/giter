package giter

type Iterator[T any] struct {
	Each  <-chan T
	Close func()
}

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
