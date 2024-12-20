package cdf

import (
	"math/rand/v2"
	"testing"
)

func TestChoose_oneValue(t *testing.T) {
	tests := []struct {
		cdf  []float64
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
		cdf, err := New(test.cdf)
		if err != nil {
			t.Fatalf("New(%v) failed: %v", test.cdf, err)
		}
		for i := 0; i < 10; i++ {
			if got := cdf.Choose(rng); got != test.want {
				t.Fatalf("Choose(%v)=%d, want %d", test.cdf, got, test.want)
			}
		}
	}
}
