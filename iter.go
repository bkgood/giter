package giter

import "github.com/bkgood/giter/internal"

type Iterator[T any] internal.Iterator[T]

func Make[T any](impl func(chan<- T, <-chan interface{})) (i Iterator[T]) {
	return Iterator[T](internal.Make(impl))
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
