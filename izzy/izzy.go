// Read simulator.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fluhus/biostuff/formats/fasta"
	"github.com/fluhus/biostuff/formats/fastq"
	"github.com/fluhus/biostuff/sequtil"
	"github.com/fluhus/gostuff/aio"
	"github.com/fluhus/gostuff/csvdec"
	"github.com/fluhus/gostuff/flagx"
	"github.com/fluhus/gostuff/gnum"
	"github.com/fluhus/gostuff/ptimer"
	"github.com/fluhus/gostuff/snm"
	"github.com/fluhus/izzy/abdist"
	"github.com/fluhus/izzy/model"
	"golang.org/x/exp/maps"
)

// BUG(amit): Add flag for grouping by file name.
// BUG(amit): Add flag for group RE.
// BUG(amit): Add flag for distribution type.

var (
	inGlob       = flag.String("i", "", "Input file glob pattern")
	outFile      = flag.String("o", "", "Output file prefix")
	nReads       = flag.Int("n", 0, "Number of reads")
	nGenomes     = flag.Int("u", 0, "Number of genomes to simulate from (default: all)")
	modelName    = flag.String("m", "", "Model name, one of "+modelNames)
	ignoreLength = flag.Bool("l", false, "Ignore genome lengths for read counts")
	singleOutput = flag.Bool("s", false, "Output one file instead of two")
	abndFile     = flag.String("a", "", "Use abundances from a file")
	re           = flagx.Regexp("g", regexp.MustCompile(".*"), "Pattern by which to group contigs of the same species")

	modelNameToModel = map[string]*model.Model{
		"basic":   model.BasicModel,
		"perfect": model.PerfectModel,
		"hiseq":   model.HiSeqModel,
		"miseq":   model.MiSeqModel,
		"novaseq": model.NovaSeqModel,
	}
	modelNames = fmt.Sprint(snm.Sorted(maps.Keys(modelNameToModel)))

	rng     = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	inFiles []string
)

func main() {
	flag.Parse()
	die(checkArgs())

	m := modelNameToModel[*modelName]

	fmt.Println("Reading sequence lengths")
	// re := regexp.MustCompile(`Rep_\d+|[iu]vig_\d+|SGB_\d+|sgb_\d+|taxid\|[^|]+|v\d+|NC_\d+`)
	// re = regexp.MustCompile(`.`)
	lens, err := readSequenceLens(inFiles, *re)
	die(err)

	groupLens := map[string]int{}
	for _, x := range lens {
		groupLens[x.g] += x.n
	}
	fmt.Println(len(groupLens), "groups")
	if *nGenomes == 0 {
		*nGenomes = len(groupLens)
	}
	if *nGenomes > len(groupLens) {
		die(fmt.Errorf("%d genomes were requested but only %d "+
			"genomes were found", *nGenomes, len(groupLens)))
	}

	var groupRatios map[string]float64
	if *abndFile != "" {
		fmt.Println("Loading abundance from file")
		groupRatios, err = readAbundanceFile(*abndFile, groupLens)
	} else {
		fmt.Println("Creating abundance distribution")
		groupRatios, err = createAbundance(groupLens, *outFile+".abundances.tsv")
	}
	die(err)

	fmt.Println("Generating reads")
	var fout1, fout2 *aio.Writer
	if *singleOutput {
		fout1, err = aio.Create(*outFile + ".fastq.gz")
		die(err)
		fout2 = fout1
	} else {
		fout1, err = aio.Create(*outFile + "_R1.fastq.gz")
		die(err)
		fout2, err = aio.Create(*outFile + "_R2.fastq.gz")
		die(err)
	}

	pt := ptimer.NewMessage("{} reads generated")
	for _, f := range inFiles {
		for fa, err := range fasta.File(f) {
			die(err)
			if !isNucs(fa.Sequence) {
				continue
			}
			gl := lens[0]
			lens = lens[1:]
			groupReads := groupRatios[gl.g] * float64(*nReads)
			seqRatio := float64(gl.n) / float64(groupLens[gl.g])

			// Convert fraction to whole number.
			nreadsf := groupReads * seqRatio
			nreads := int(math.Floor(nreadsf))
			if rand.Float64() < nreadsf-math.Floor(nreadsf) {
				nreads++
			}

			simulateReads(fa.Sequence, m, nreads,
				func(fwd, bwd *fastq.Fastq) error {
					if len(fwd.Sequence) != m.ReadLen {
						return fmt.Errorf("bad read length: %d, want %d",
							len(fwd.Sequence), m.ReadLen)
					}
					if len(bwd.Sequence) != m.ReadLen {
						return fmt.Errorf("bad read length: %d, want %d",
							len(bwd.Sequence), m.ReadLen)
					}
					fwd.Name = []byte(fmt.Sprintf(
						"%d.%s.%s", pt.N+1, fwd.Name, fa.Name))
					bwd.Name = []byte(fmt.Sprintf(
						"%d.%s.%s", pt.N+2, bwd.Name, fa.Name))
					txt, _ := fwd.MarshalText()
					fout1.Write(txt)
					txt, _ = bwd.MarshalText()
					fout2.Write(txt)
					pt.Inc()
					pt.Inc() // Each pair is 2 reads.
					return nil
				})
		}
	}
	fout1.Close()
	fout2.Close()
	pt.Done()
}

func checkArgs() error {
	if *outFile == "" {
		return fmt.Errorf("no output file")
	}
	if *nReads < 1 {
		return fmt.Errorf("number of reads needs to be at least 1")
	}
	*nReads = (*nReads + 1) / 2 // We will create nreads/2 pairs.
	if *nGenomes < 0 {
		return fmt.Errorf("bad number of genomes: %d", *nGenomes)
	}
	files, err := filepath.Glob(*inGlob)
	if err != nil {
		return err
	}
	inFiles = files
	if len(files) == 0 {
		return fmt.Errorf("found 0 input files")
	}
	if modelNameToModel[*modelName] == nil {
		return fmt.Errorf("bad model name: %q, need one of %v",
			*modelName, snm.Sorted(maps.Keys(modelNameToModel)))
	}
	return nil
}

func readSequenceLens(files []string, grouper *regexp.Regexp) ([]lenGroup, error) {
	pt := ptimer.NewMessage("{} sequences read")
	var result []lenGroup
	for _, f := range files {
		for fa, err := range fasta.File(f) {
			if err != nil {
				return nil, err
			}
			g := string(fa.Name)
			if grouper != nil {
				g = grouper.FindString(g)
			}
			if isNucs(fa.Sequence) {
				result = append(result, lenGroup{g, len(fa.Sequence)})
			}
			pt.Inc()
		}
	}
	pt.Done()
	return result, nil
}

type lenGroup struct {
	g string // Group name
	n int    // Length of sequence
}

func simulateReads(seq []byte, m *model.Model, n int,
	forEach func(fwd, bwd *fastq.Fastq) error) error {
	for i := 0; i < n; i++ {
		fwd, bwd := m.SimulateRead(seq, rng)
		if fwd == nil { // Sequence is too short.
			continue
		}
		if err := forEach(fwd, bwd); err != nil {
			return err
		}
	}
	return nil
}

func createAbundance(groupLens map[string]int, file string) (map[string]float64, error) {
	fout, err := aio.Create(file)
	if err != nil {
		return nil, err
	}
	defer fout.Close()

	abnd := abdist.LogNormal(len(groupLens), *nGenomes)
	groupRatios := map[string]float64{}
	for k := range groupLens {
		ab := abnd[0]
		abnd = abnd[1:]
		if ab == 0 {
			continue
		}
		if _, err := fmt.Fprintf(fout, "%s\t%.10f\n", k, ab); err != nil {
			return nil, err
		}
		groupRatios[k] = ab
		if !*ignoreLength {
			groupRatios[k] *= float64(groupLens[k])
		}
	}
	sum := gnum.Sum(maps.Values(groupRatios))
	for k := range groupRatios {
		groupRatios[k] /= sum
	}
	return groupRatios, nil
}

func readAbundanceFile(file string, groupLens map[string]int) (map[string]float64, error) {
	type entry struct {
		Name string
		Abnd float64
	}
	result := map[string]float64{}
	for row, err := range csvdec.File[entry](file, toTSV) {
		if err != nil {
			return nil, err
		}
		if _, ok := result[row.Name]; ok {
			return nil, fmt.Errorf("duplicate name: %s, values: %v %v",
				row.Name, result[row.Name], row.Abnd)
		}
		if row.Abnd <= 0 {
			return nil, fmt.Errorf("bad abundance value: %v", row.Abnd)
		}
		ln, ok := groupLens[row.Name]
		if !ok {
			return nil, fmt.Errorf("unrecognized name: %s", row.Name)
		}
		result[row.Name] = row.Abnd
		if !*ignoreLength {
			result[row.Name] *= float64(ln)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no abundances in file")
	}
	sum := gnum.Sum(maps.Values(result))
	for k := range result {
		result[k] /= sum
	}
	return result, nil
}

func isNucs(seq []byte) bool {
	for _, b := range seq {
		if sequtil.Ntoi(b) == -1 {
			return false
		}
	}
	return true
}

// Prints the error and exits if the error is non-nil.
func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(2)
	}
}

func toTSV(r *csv.Reader) {
	r.Comma = '\t'
}
