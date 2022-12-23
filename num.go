package giter

// Sum consumes an Iterator of numbers and returns their sum (or zero, if no
// values were emitted).
func Sum[T int | int32 | int64 | float32 | float64](iter Iterator[T]) T {
	return Fold(0, func(x, y T) T { return x + y }, iter)
}

// Prod consumes an Iterator of numbers and returns their product (or one, if
// no values were emitted).
func Prod[T int | int32 | int64 | float32 | float64](iter Iterator[T]) T {
	return Fold(1, func(x, y T) T { return x * y }, iter)
}

// Range returns an iterator emitting numeric values over a given range.
//
// The given range is half-open, inclusive on the left and exclusive on the right.
//
// Values are emitted in steps of 1; odd floating point values may cause unexpected results.
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

// Range returns an iterator emitting numeric values over a given range with
// a given step.
//
// The given range is half-open, inclusive on the left and exclusive on the right.
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
