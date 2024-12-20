package model

import (
	"bytes"
	"math/rand/v2"
	"testing"
)

func TestPerfectModel(t *testing.T) {
	m := PerfectModel
	seq := bytes.Repeat([]byte("G"), 450)
	wantFwd := bytes.Repeat([]byte("G"), 125)
	wantBwd := bytes.Repeat([]byte("C"), 125)
	wantQual := bytes.Repeat([]byte{73}, 125)
	wantFwdName := []byte("1")
	wantBwdName := []byte("326")
	rng := rand.New(rand.NewPCG(0, 0))

	for i := 0; i < 10; i++ {
		r1, r2 := m.SimulateRead(seq, rng)
		if !bytes.Equal(wantFwd, r1.Sequence) {
			t.Errorf("r1.Sequence=%q, want %q", r1.Sequence, wantFwd)
		}
		if !bytes.Equal(wantQual, r1.Quals) {
			t.Errorf("r1.Quals=%q, want %q", r1.Quals, wantQual)
		}
		if !bytes.Equal(wantFwdName, r1.Name) {
			t.Errorf("r1.Name=%q, want %q", r1.Name, wantFwdName)
		}
		if !bytes.Equal(wantBwd, r2.Sequence) {
			t.Errorf("r2.Sequence=%q, want %q", r2.Sequence, wantBwd)
		}
		if !bytes.Equal(wantQual, r2.Quals) {
			t.Errorf("r2.Quals=%q, want %q", r2.Quals, wantQual)
		}
		if !bytes.Equal(wantBwdName, r2.Name) {
			t.Errorf("r2.Name=%q, want %q", r2.Name, wantBwdName)
		}
	}
}

func TestPerfectModel_pos(t *testing.T) {
	m := PerfectModel
	seq := bytes.Repeat([]byte("A"), 451)
	wantFwd := map[string]int{"1": 5, "2": 5}
	wantBwd := map[string]int{"326": 5, "327": 5}
	gotFwd := map[string]int{}
	gotBwd := map[string]int{}
	rng := rand.New(rand.NewPCG(0, 0))

	for i := 0; i < 20; i++ {
		r1, r2 := m.SimulateRead(seq, rng)
		gotFwd[string(r1.Name)]++
		gotBwd[string(r2.Name)]++
	}
	if !mapAtLeast(gotFwd, wantFwd) {
		t.Errorf("pos count=%v, want at least %v", gotFwd, wantFwd)
	}
	if !mapAtLeast(gotBwd, wantBwd) {
		t.Errorf("pos count=%v, want at least %v", gotBwd, wantBwd)
	}
}

func mapAtLeast(m1, m2 map[string]int) bool {
	for k, v := range m2 {
		if v > m1[k] {
			return false
		}
	}
	return true
}
