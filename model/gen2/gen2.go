package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"

	"github.com/fluhus/gostuff/jio"
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
	fmt.Fprintf(buf, "var %sModel = %#v\n", name, m)
}
