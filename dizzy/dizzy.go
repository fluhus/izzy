// Generates random reads.
package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/fluhus/biostuff/formats/fastq"
	"github.com/fluhus/gostuff/aio"
	"github.com/fluhus/gostuff/flagx"
	"github.com/fluhus/gostuff/ptimer"
	"github.com/fluhus/gostuff/snm"
)

var (
	nReads  = flag.Uint("n", 0, "Number of reads")
	readLen = flag.Uint("l", 100, "Read length")
	outFile = flag.String("o", "", "Output file")
	gcPerc  = flagx.IntBetween("gc", 50, "GC percent (default 50)", 0, 100)
)

func main() {
	flag.Parse()

	fout, err := aio.Create(*outFile)
	die(err)

	dist := makeDistribution()
	fq := &fastq.Fastq{Sequence: make([]byte, *readLen),
		Quals: snm.Slice(int(*readLen), func(i int) byte { return 'I' })}
	pt := ptimer.New()

	for i := range *nReads {
		fq.Name = []byte(fmt.Sprint(i + 1))
		for j := range fq.Sequence {
			fq.Sequence[j] = dist[rand.IntN(len(dist))]
		}
		die(fq.Write(fout))
		pt.Inc()
	}
	fout.Close()
	pt.Done()
}

// Prints the error and exits if the error is non-nil.
func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(2)
	}
}

// Creates the base distribution according to the GC percent.
func makeDistribution() []byte {
	dist := make([]byte, 0, 200)
	for range *gcPerc {
		dist = append(dist, 'G', 'C')
	}
	for range 100 - *gcPerc {
		dist = append(dist, 'A', 'T')
	}
	if len(dist) != 200 {
		panic(fmt.Sprintf("bad length: %d, want 200", len(dist)))
	}
	return dist
}
