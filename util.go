package giter

func clear[T any](xs *[]T) {
	var zero T

	for i := range *xs {
		(*xs)[i] = zero
	}

	*xs = (*xs)[:0]
}
