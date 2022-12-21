package giter

func Sum[T int | int32 | int64 | float32 | float64](iter Iterator[T]) T {
	return Fold(0, func(x, y T) T { return x + y }, iter)
}

// XXX gived 1 on empty. i am ok with this.
func Prod[T int | int32 | int64 | float32 | float64](iter Iterator[T]) T {
	return Fold(1, func(x, y T) T { return x * y }, iter)
}

func Range[T int | int32 | int64 | float32 | float64](from, until T) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			for v := from; v < until; v++ {
				select {
				case out <- v:
				case <-stopChan:
					return
				}
			}
		})
}

func RangeBy[T int | int32 | int64 | float32 | float64](from, until, by T) Iterator[T] {
	return Make(
		func(out chan<- T, stopChan <-chan interface{}) {
			for v := from; v < until; v += by {
				select {
				case out <- v:
				case <-stopChan:
					return
				}
			}
		})
}
