package internal

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"sort"
)

func WrapErr(name string, errptr *error) {
	if *errptr != nil {
		*errptr = fmt.Errorf(name+": %w", *errptr)
	}
}

func Uniq[T constraints.Ordered](arr []T) []T {
	var result []T
	uniq := map[T]struct{}{}
	for _, item := range arr {
		if _, exists := uniq[item]; exists {
			continue
		}

		uniq[item] = struct{}{}
		result = append(result, item)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}
