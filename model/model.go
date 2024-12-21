// Package model implements a probabilistic model for generating reads.
package model

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/fluhus/biostuff/formats/fastq"
	"github.com/fluhus/biostuff/sequtil"
	"github.com/fluhus/gostuff/snm"
	"github.com/fluhus/izzy/cdf"
	"golang.org/x/exp/slices"
)

//go:generate go run ./gen/gen.go
//go:generate go run ./gen2/gen2.go

// BUG(amit): Add unit tests for model.

const (
	// Use the original indel logic, for testing.
	originalIndel = false
)

// Model holds probabilities for randomizing reads.
type Model struct {
	Name                string
	ReadLen             int
	InsertLen           cdf.CDF
	MeanCountForward    cdf.CDF
	MeanCountReverse    cdf.CDF
	QualityHistForward  [][]cdf.CDF
	QualityHistReverse  [][]cdf.CDF
	SubstChoicesForward [][4]cdf.CDF
	SubstChoicesReverse [][4]cdf.CDF
	InsForward          [][4]float64
	InsReverse          [][4]float64
	DelForward          [][4]float64
	DelReverse          [][4]float64
}

// Returns a slice of ReadLen random phred scores.
func (m *Model) genPhredScores(forward bool, rng *rand.Rand) []int {
	mean := m.MeanCountForward
	if !forward {
		mean = m.MeanCountReverse
	}
	qbin := mean.Choose(rng)

	cdfss := m.QualityHistForward
	if !forward {
		cdfss = m.QualityHistReverse
	}
	cdfs := cdfss[qbin]

	return snm.Slice(m.ReadLen, func(i int) int {
		return cdfs[i].Choose(rng)
	})
}

// Returns a random insert size between the forward read and its reverse
// end.
func (m *Model) randomInsertSize(rng *rand.Rand) int {
	return m.InsertLen.Choose(rng)
}

// Applies indels to seq and returns the new sequence.
func (m *Model) introduceIndels(seq []byte, forward bool, rng *rand.Rand,
) []byte {
	ins, del := m.InsForward, m.DelForward
	if !forward {
		ins, del = m.InsReverse, m.DelReverse
	}
	result := make([]byte, 0, len(seq)*11/10)
	if !originalIndel {
		// BUG(amit): In ISS i runs up to len-1, not sure why.
		for i, b := range seq {
			// Deletion - skip if rand < p.
			ntoi := sequtil.Ntoi(b)
			if ntoi == -1 {
				continue
			}
			if rng.Float64() > del[i][ntoi] {
				result = append(result, b)
			}
			// BUG(amit): Not sure about this insertion logic. Taken from ISS.
			for ii, p := range ins[i] {
				if rng.Float64() < p {
					result = append(result, sequtil.Iton(ii))
				}
			}
		}
	} else {
		result = slices.Clone(seq)
		// Original logic from ISS.
		pos := 0
		for i := range make([]struct{}, m.ReadLen-1) {
			// Deletion - skip if rand < p.
			ntoi := sequtil.Ntoi(seq[i])
			if ntoi == -1 {
				pos++
				continue
			}
			for ii, p := range ins[pos] {
				if rng.Float64() < p {
					// Insert at pos+1.
					if pos >= len(result) {
						result = append(result, sequtil.Iton(ii))
					} else {
						result = append(result[:pos+2], slices.Clone(result[pos+1:])...)
						result[pos+1] = sequtil.Iton(ii)
					}
				}
			}
			if rng.Float64() < del[pos][ntoi] {
				if pos >= len(result) {
					continue
				}
				result = append(result[:pos], result[pos+1:]...)
			}
			pos++
		}
	}
	return result
}

// Applies SNPs to seq according to the given phred scores.
func (m *Model) introduceSNPs(seq []byte, phreds []int, forward bool,
	rng *rand.Rand) {
	subst := m.SubstChoicesForward
	if !forward {
		subst = m.SubstChoicesReverse
	}
	for i := range seq {
		p := phredToProb[phreds[i]]
		if rng.Float64() < p {
			ntoi := sequtil.Ntoi(seq[i])
			if ntoi == -1 {
				continue
			}
			cdf := subst[i][ntoi]
			seq[i] = sequtil.Iton(cdf.Choose(rng))
		}
	}
}

// SimulateRead randomizes a pair of reads from seq.
// Returns nil of seq is too short.
func (m *Model) SimulateRead(seq []byte, rng *rand.Rand,
) (*fastq.Fastq, *fastq.Fastq) {
	if len(seq) < 2*m.ReadLen {
		return nil, nil
	}
	intervalLen := 2*m.ReadLen + m.randomInsertSize(rng)
	// BUG(amit): Check if this is the best thing to do in this case.
	if intervalLen > len(seq) {
		intervalLen = len(seq)
	}
	i := rng.IntN(len(seq) - intervalLen + 1)
	bwdStart := i + intervalLen - m.ReadLen
	fwd := slices.Clone(seq[i : i+m.ReadLen])
	bwd := sequtil.ReverseComplement(nil, seq[bwdStart:bwdStart+m.ReadLen])

	fwd = m.introduceIndels(fwd, true, rng)
	if len(fwd) > m.ReadLen {
		fwd = fwd[:m.ReadLen]
	}
	if len(fwd) < m.ReadLen {
		d := m.ReadLen - len(fwd)
		fwd = append(fwd, seq[i+m.ReadLen:i+m.ReadLen+d]...)
	}

	bwd = m.introduceIndels(bwd, true, rng)
	if len(bwd) > m.ReadLen {
		bwd = bwd[:m.ReadLen]
	}
	if len(bwd) < m.ReadLen {
		if !originalIndel {
			d := m.ReadLen - len(bwd)
			bwd = sequtil.ReverseComplement(bwd, seq[bwdStart-d:bwdStart])
			bwdStart -= d
		} else {
			d := m.ReadLen - len(bwd)
			for i := 0; i < d; i++ {
				ii := i + intervalLen
				if ii >= len(seq) {
					bwd = append(bwd, 'A')
				} else {
					bwd = sequtil.ReverseComplement(bwd, seq[ii:ii+1])
				}
			}
		}
	}

	fwdQuals := m.genPhredScores(true, rng)
	bwdQuals := m.genPhredScores(false, rng)

	m.introduceSNPs(fwd, fwdQuals, true, rng)
	m.introduceSNPs(bwd, bwdQuals, false, rng)

	// +1 to convert positions to 1-based.
	fwdq := &fastq.Fastq{
		Name:     []byte(fmt.Sprint(i + 1)),
		Sequence: fwd,
		Quals:    phredsToASCII(fwdQuals),
	}
	bwdq := &fastq.Fastq{
		Name:     []byte(fmt.Sprint(bwdStart + 1)),
		Sequence: bwd,
		Quals:    phredsToASCII(bwdQuals),
	}

	return fwdq, bwdq
}

// Encodes the given phred scores as ASCII for text output.
func phredsToASCII(phreds []int) []byte {
	return snm.Slice(len(phreds), func(i int) byte {
		return 33 + byte(phreds[i])
	})
}

// From phred score to error probability.
var phredToProb = snm.Slice(100, func(i int) float64 {
	return math.Pow(10, -float64(i)/10)
})
