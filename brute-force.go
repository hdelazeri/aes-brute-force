package main

import (
	"unicode/utf8"
)

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateKeys(base string) func(func(string) bool) {
	return func(yield func(string) bool) {
		if len(base) == 16 {
			yield(base)
			return
		}

		for i := 0; i < len(characters); i++ {
			next := base + string(characters[i])

			for k := range GenerateKeys(next) {
				if !yield(k) {
					return
				}
			}
		}
	}
}

func IsASCIIBytes(s []byte) bool {
	return utf8.Valid(s)
}
