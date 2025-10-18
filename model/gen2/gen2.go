// Generates code variables holding each model from the JSON files.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"

	"github.com/fluhus/gostuff/jio"
	"github.com/fluhus/izzy/cdf"
	"github.com/fluhus/izzy/model"
)

func main() {
	buf := bytes.NewBufferString(`// Auto-generated models.
	
	package model
	import "github.com/fluhus/izzy/cdf"
	`)
	models := []string{"Basic", "Perfect", "HiSeq", "MiSeq", "NovaSeq"}
	for _, name := range models {
		writeModelVar(name, buf)
	}
	b := bytes.ReplaceAll(buf.Bytes(), []byte("model.Model"), []byte("Model"))
	src, err := format.Source(b)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("models.go", src, 0o644)
	if err != nil {
		panic(err)
	}
}

func writeModelVar(name string, buf *bytes.Buffer) {
	m := &model.Model{}
	if err := jio.Read("gen2/"+name+".json", m); err != nil {
		panic(err)
	}
	checkModelCDFs(m)
	fmt.Fprintf(buf, "var %sModel = %#v\n", name, m)
}

func checkModelCDFs(m *model.Model) {
	m.InsertLen.Check()
	m.MeanCountForward.Check()
	m.MeanCountReverse.Check()
	checkCDFss(m.QualityHistForward)
	checkCDFss(m.QualityHistReverse)
	checkCDFs4(m.SubstChoicesForward)
	checkCDFs4(m.SubstChoicesReverse)
}

func checkCDFs(cc []cdf.CDF) {
	for _, c := range cc {
		c.Check()
	}
}

func checkCDFss(cc [][]cdf.CDF) {
	for _, c := range cc {
		checkCDFs(c)
	}
}

func checkCDFs4(ccc [][4]cdf.CDF) {
	for _, cc := range ccc {
		for _, c := range cc {
			c.Check()
		}
	}
}
