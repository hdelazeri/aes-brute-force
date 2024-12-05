package main

import "unicode"

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateKeys(base string) func(func(string) bool) {
	return func(yield func(string) bool) {
		remaining := 16 - len(base)
		for combination := range Combinations([]rune(characters), remaining) {
			for suffix := range Permutations(combination) {
				if !yield(base + string(suffix)) {
					return
				}
			}
		}
	}
}

func IsASCIIBytes(s []byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
