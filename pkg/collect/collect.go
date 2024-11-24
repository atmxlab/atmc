package collect

import "github.com/samber/lo"

type Collection[T any] []T

func (c Collection[T]) Append(item ...T) Collection[T] {
	return append(c, item...)
}

func Map[T any, R any](items []T, fn func(item T) R) []R {
	return lo.Map(items, func(item T, index int) R {
		return fn(item)
	})
}
