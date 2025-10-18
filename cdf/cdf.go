// Package cdf provides a CDF type for discrete sampling.
package cdf

import (
	"fmt"
	"math/rand/v2"

	"golang.org/x/exp/slices"
)

// CDF is a cumulative distribution.
type CDF []float64

// Check checks that a CDF is non-empty, non-decreasing,
// and ends in 1. Panics if not.
func (c CDF) Check() {
	if len(c) == 0 {
		panic("got empty cdf")
	}
	if c[len(c)-1] != 1 {
		panic(fmt.Sprintf("last element is %f, want 1", c[len(c)-1]))
	}
	for i := range c {
		if c[i] < 0 {
			panic(fmt.Sprintf("cdf[%d]=%f, want >=0", i, c[len(c)-1]))
		}
		if i > 0 && c[i-1] > c[i] {
			panic(fmt.Sprintf("cdf[%d]>cdf[%d]: %f>%f",
				i-1, i, c[i-1], c[i]))
		}
	}
}

// Choose picks an element from the CDF according to the distribution.
func (c CDF) Choose(rng *rand.Rand) int {
	p := rng.Float64()
	i, _ := slices.BinarySearch(c, p)
	return i
}

// GoString returns a string for generating code that includes building CDFs.
func (c CDF) GoString() string {
	return fmt.Sprintf("%#v", []float64(c))
}
