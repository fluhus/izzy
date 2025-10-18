package cdf

import (
	"math/rand/v2"
	"testing"
)

func TestChoose_oneValue(t *testing.T) {
	tests := []struct {
		cdf  CDF
		want int
	}{
		{[]float64{1}, 0},
		{[]float64{0, 1}, 1},
		{[]float64{1, 1}, 0},
		{[]float64{0, 0, 1}, 2},
		{[]float64{0, 1, 1}, 1},
		{[]float64{1, 1, 1}, 0},
	}
	rng := rand.New(rand.NewPCG(0, 0))
	for _, test := range tests {
		test.cdf.Check()
		for range 10 {
			if got := test.cdf.Choose(rng); got != test.want {
				t.Fatalf("Choose(%v)=%d, want %d", test.cdf, got, test.want)
			}
		}
	}
}
