package main

import "math/bits"

func Combinations[T any](set []T, n int) func(func([]T) bool) {
	return func(yield func([]T) bool) {
		length := uint(len(set))

		if n > len(set) {
			n = len(set)
		}

		for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
			if n > 0 && bits.OnesCount(uint(subsetBits)) != n {
				continue
			}

			var subset []T

			for object := uint(0); object < length; object++ {
				if (subsetBits>>object)&1 == 1 {
					subset = append(subset, set[object])
				}
			}

			if !yield(subset) {
				return
			}
		}
	}
}

func Permutations[T any](set []T) func(func([]T) bool) {
	return func(yield func([]T) bool) {
		state := make([]int, len(set))

		for i := 0; i < len(state); i++ {
			state[i] = 0
		}

		if !yield(set) {
			return
		}

		i := 1
		for i < len(set) {
			if state[i] < i {
				if i%2 == 0 {
					set[0], set[i] = set[i], set[0]
				} else {
					set[state[i]], set[i] = set[i], set[state[i]]
				}

				if !yield(set) {
					return
				}

				state[i] += 1
				i = 1
			} else {
				state[i] = 0
				i += 1
			}
		}
	}
}
