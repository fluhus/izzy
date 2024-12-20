// Generates JSON for the basic and perfect models.
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/fluhus/gostuff/gnum"
	"github.com/fluhus/gostuff/snm"
	"github.com/fluhus/izzy/cdf"
)

func main() {
	m := basicModel()
	j, _ := json.Marshal(m)
	os.WriteFile("Basic.json", j, 0o644)
	m = perfectModel()
	j, _ = json.Marshal(m)
	os.WriteFile("Perfect.json", j, 0o644)
}

func basicModel() map[string]any {
	const (
		readLen   = 125
		insertLen = 200
	)

	m := map[string]any{}
	m["name"] = "basic"
	m["readLen"] = readLen
	m["insertLen"] = append(make([]float64, insertLen), 1)
	m["meanCountForward"] = []float64{1}
	m["meanCountReverse"] = []float64{1}

	phreds := phredCDF(0.999, 0.01, 0.9999)
	m["qualityHistForward"] = [][]cdf.CDF{snm.Slice(readLen, func(i int) cdf.CDF {
		return phreds
	})}
	m["qualityHistReverse"] = m["qualityHistForward"]

	subst := basicSubstCDF()
	m["substChoicesForward"] = snm.Slice(readLen, func(i int) [4]cdf.CDF {
		return subst
	})
	m["substChoicesReverse"] = m["substChoicesForward"]

	indel := [4]float64{} // All zeros
	m["insForward"] = snm.Slice(readLen, func(i int) [4]float64 {
		return indel
	})
	m["insReverse"] = m["insForward"]
	m["delForward"] = m["insForward"]
	m["delReverse"] = m["insForward"]

	return m
}

func perfectModel() map[string]any {
	const (
		readLen   = 125
		insertLen = 200
		phred     = 40
	)

	m := map[string]any{}
	m["name"] = "perfect"
	m["readLen"] = readLen
	m["insertLen"] = append(make([]float64, insertLen), 1)
	m["meanCountForward"] = []float64{1}
	m["meanCountReverse"] = []float64{1}

	phreds := toCDF(append(make([]float64, phred), 1))
	m["qualityHistForward"] = [][]cdf.CDF{snm.Slice(readLen, func(i int) cdf.CDF {
		return phreds
	})}
	m["qualityHistReverse"] = m["qualityHistForward"]

	subst := perfectSubstCDF()
	m["substChoicesForward"] = snm.Slice(readLen, func(i int) [4]cdf.CDF {
		return subst
	})
	m["substChoicesReverse"] = m["substChoicesForward"]

	indel := [4]float64{} // All zeros
	m["insForward"] = snm.Slice(readLen, func(i int) [4]float64 {
		return indel
	})
	m["insReverse"] = m["insForward"]
	m["delForward"] = m["insForward"]
	m["delReverse"] = m["insForward"]

	return m
}

func phredCDF(mean, std, max float64) cdf.CDF {
	var ps []float64
	for i := 0; i < 5000000; i++ {
		n := min(rand.NormFloat64()*std+mean, max)
		p := int(math.Round(-10 * math.Log10(1-n)))
		for p >= len(ps) {
			ps = append(ps, 0)
		}
		ps[p]++
	}
	return toCDF(ps)
}

func basicSubstCDF() [4]cdf.CDF {
	return [4]cdf.CDF{
		toCDF([]float64{0, 1, 1, 1}),
		toCDF([]float64{1, 0, 1, 1}),
		toCDF([]float64{1, 1, 0, 1}),
		toCDF([]float64{1, 1, 1, 0}),
	}
}

func perfectSubstCDF() [4]cdf.CDF {
	return [4]cdf.CDF{
		toCDF([]float64{1, 0, 0, 0}),
		toCDF([]float64{0, 1, 0, 0}),
		toCDF([]float64{0, 0, 1, 0}),
		toCDF([]float64{0, 0, 0, 1}),
	}
}

func toCDF(ps []float64) cdf.CDF {
	for i := range ps[1:] {
		ps[i+1] += ps[i]
	}
	gnum.Mul1(ps, 1/ps[len(ps)-1])
	cdf, err := cdf.New(ps)
	if err != nil {
		panic(fmt.Sprintf("creating CDF: %v", err))
	}
	return cdf
}
