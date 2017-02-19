package main

import "math/rand"

var kRandLetters = []byte("01234567890abcdefghijklmnopqrstuvwxzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = kRandLetters[rand.Intn(len(kRandLetters))]
	}
	return b
}
