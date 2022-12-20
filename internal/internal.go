package internal

type Iterator[T any] struct {
	Each  <-chan T
	Close func()
}

func Make[T any](impl func(chan<- T, <-chan interface{})) (i Iterator[T]) {
	values := make(chan T)
	stop, stopChan := makeStop()

	go func() {
		impl(values, stopChan)
		close(values)
	}()

	return Iterator[T]{Each: values, Close: stop}
}

func makeStop() (stop func(), stopChan <-chan interface{}) {
	bidiStopChan := make(chan interface{})
	stopChan = bidiStopChan

	stop = func() {
		select {
		case bidiStopChan <- nil:
		default:
		}
	}

	return stop, stopChan
}
