package fn

import "iter"

type SL[T any] []T

func (s SL[T]) TransForm(fn func(T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func (s SL[T]) Filter(fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if fn(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}
