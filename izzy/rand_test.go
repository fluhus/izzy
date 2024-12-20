package main

import (
	rand1 "math/rand"
	"math/rand/v2"
	"testing"
)

func BenchmarkRand(b *testing.B) {
	b.Run("PCG", func(b *testing.B) {
		rng := rand.New(rand.NewPCG(0, 0))
		for range b.N {
			rng.Int64()
		}
	})
	b.Run("ChaCha8", func(b *testing.B) {
		rng := rand.New(rand.NewChaCha8([32]byte{}))
		for range b.N {
			rng.Int64()
		}
	})
	b.Run("Rand1", func(b *testing.B) {
		rng := rand1.New(rand1.NewSource(0))
		for range b.N {
			rng.Int63()
		}
	})
}
