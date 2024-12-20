// Package abdist provides abundance distributions.
//
// Each function returns a normalized vector of length n,
// with nz non-zero values.
package abdist

import (
	"math"
	"math/rand/v2"

	"github.com/fluhus/gostuff/gnum"
)

const (
	lognormalScale = 1.5 // STD of the normal distribution in lognormal. Measured in 10K samples.
)

// BUG(amit): Add zero-inflated lognormal?

// LogNormal returns lognormal values (exp(normal)).
func LogNormal(n, nz int) []float64 {
	return abndnc(n, nz, func() float64 {
		return math.Exp(rand.NormFloat64() * lognormalScale)
	})
}

// Uniform returns a uniform distribution.
func Uniform(n, nz int) []float64 {
	return abndnc(n, nz, func() float64 { return 1 })
}

// HalfNormal returns half-normal values (abs(normal)).
func HalfNormal(n, nz int) []float64 {
	return abndnc(n, nz, func() float64 {
		return math.Abs(rand.NormFloat64())
	})
}

// Exponential returns an exponential distribution.
func Exponential(n, nz int) []float64 {
	return abndnc(n, nz, rand.ExpFloat64)
}

// Returns a normalized vector of size n with nz non-zeros,
// each non-zero is generated with p.
func abndnc(n, nz int, p func() float64) []float64 {
	a := make([]float64, n)
	for _, i := range rand.Perm(n)[:nz] {
		a[i] = p()
	}
	gnum.Mul1(a, 1.0/gnum.Sum(a))
	return a
}
