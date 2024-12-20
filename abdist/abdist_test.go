package abdist

import (
	"testing"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func TestUniform(t *testing.T) {
	want := []float64{0.2, 0.2, 0.2, 0.2, 0.2}
	got := Uniform(5, 5)
	if !slices.Equal(got, want) {
		t.Fatalf("Uniform(5,5)=%v, want %v", got, want)
	}
}

func TestUniform_nz(t *testing.T) {
	want := map[float64]int{0: 3, 0.5: 2}
	got := map[float64]int{}
	for _, v := range Uniform(5, 2) {
		got[v]++
	}
	if !maps.Equal(got, want) {
		t.Fatalf("Uniform(5,2)=%v, want %v", got, want)
	}
}
