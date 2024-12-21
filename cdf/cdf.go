// Package cdf provides a CDF type for discrete sampling.
package cdf

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"golang.org/x/exp/slices"
)

type CDF struct {
	a []float64
}

func New(cdf []float64) (CDF, error) {
	if len(cdf) == 0 {
		return CDF{}, fmt.Errorf("got empty cdf")
	}
	if cdf[len(cdf)-1] != 1 {
		return CDF{}, fmt.Errorf("last element is %f, want 1",
			cdf[len(cdf)-1])
	}
	for i := range cdf {
		if cdf[i] < 0 {
			return CDF{}, fmt.Errorf("cdf[%d]=%f, want >=0",
				i, cdf[len(cdf)-1])
		}
		if i > 0 && cdf[i-1] > cdf[i] {
			return CDF{}, fmt.Errorf("cdf[%d]>cdf[%d]: %f>%f",
				i-1, i, cdf[i-1], cdf[i])
		}
	}
	return CDF{slices.Clone(cdf)}, nil
}

func Must(cdf []float64) CDF {
	c, err := New(cdf)
	if err != nil {
		panic(err)
	}
	return c
}

func (c CDF) Choose(rng *rand.Rand) int {
	p := rng.Float64()
	i, _ := slices.BinarySearch(c.a, p)
	return i
}

func (c *CDF) UnmarshalJSON(data []byte) error {
	var a []float64
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	cdf, err := New(a)
	if err != nil {
		return err
	}
	*c = cdf
	return nil
}

func (c *CDF) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.a)
}

func (c CDF) GoString() string {
	return fmt.Sprintf("cdf.Must(%#v)", c.a)
}
