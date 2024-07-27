package utils

import "github.com/samber/lo"

func Select[T any, Y any](in []T, trans func(p T, index int) Y) (out []Y) {
	return lo.Map(in, trans)
}

func Where[T any](in []T, filter func(p T, index int) bool) []T {
	return lo.Filter(in, filter)
}

func Distinct[T any, Y comparable](src []T, dKey func(item T) Y) []T {
	return lo.UniqBy(src, dKey)
}

func Contains[T comparable](src T, dst ...T) bool {
	return lo.Contains(dst, src)
}

func GroupBy[T any, Y comparable](src []T, do func(in T) Y) map[Y][]T {
	return lo.GroupBy(src, do)
}
