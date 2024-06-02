package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/prim/util"
	"io"
)

type Parameters struct {
	primMinSize, primMaxSize,
	prodMinSize, prodMaxSize int
	primMinTm, primOptTm, primMaxTm,
	inMinTm, inOptTm, inMaxTm float64
}

func parse(r io.Reader, args ...interface{}) {
	p := args[0].(*Parameters)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		s := string(sc.Sequence().Data())
		fmt.Println("PRIMER_TASK=generic")
		fmt.Println("PRIMER_PICK_LEFT_PRIMER=1")
		fmt.Println("PRIMER_PICK_RIGHT_PRIMER=1")
		fmt.Println("PRIMER_PICK_INTERNAL_OLIGO=1")
		fmt.Printf("PRIMER_MIN_SIZE=%d\n", p.primMinSize)
		fmt.Printf("PRIMER_MAX_SIZE=%d\n", p.primMaxSize)
		fmt.Printf("PRIMER_PRODUCT_SIZE_RANGE=%d-%d\n",
			p.prodMinSize, p.prodMaxSize)
		fmt.Printf("PRIMER_MIN_TM=%.1f\n", p.primMinTm)
		fmt.Printf("PRIMER_OPT_TM=%.1f\n", p.primOptTm)
		fmt.Printf("PRIMER_MAX_TM=%.1f\n", p.primMaxTm)
		fmt.Printf("PRIMER_INTERNAL_MIN_TM=%.1f\n", p.inMinTm)
		fmt.Printf("PRIMER_INTERNAL_OPT_TM=%.1f\n", p.inOptTm)
		fmt.Printf("PRIMER_INTERNAL_MAX_TM=%.1f\n", p.inMaxTm)
		fmt.Printf("SEQUENCE_TEMPLATE=%s\n", s)
		fmt.Println("=")
	}
}
func main() {
	util.SetName("fa2prim")
	u := "fa2prim [option]... [template.fasta]..."
	p := "Convert FASTA sequences to primer3 input."
	e := "fa2prim foo.fasta | primer3_core"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	primMinSize := flag.Int("primMinSize", 15,
		"minimum primer size")
	primMaxSize := flag.Int("primMaxSize", 25,
		"maximum primer size")
	prodMinSize := flag.Int("prodMinSize", 70,
		"minimum product size")
	prodMaxSize := flag.Int("prodMaxSize", 150,
		"maximum product size")
	primMinTm := flag.Float64("primMinTm", 54,
		"minimum primer T_m")
	primOptTm := flag.Float64("primOptTm", 56,
		"optimal primer T_m")
	primMaxTm := flag.Float64("primMaxTm", 58,
		"maximum primer T_m")
	inMinTm := flag.Float64("inMinTm", 43,
		"minimum internal oligo T_m")
	inOptTm := flag.Float64("inOptTm", 45,
		"optimal internal oligo T_m")
	inMaxTm := flag.Float64("inMaxTm", 47,
		"maximum internal oligo T_m")
	flag.Parse()
	if *optV {
		util.Version()
	}
	pa := new(Parameters)
	pa.primMinSize = *primMinSize
	pa.primMaxSize = *primMaxSize
	pa.prodMinSize = *prodMinSize
	pa.prodMaxSize = *prodMaxSize
	pa.primMinTm = *primMinTm
	pa.primOptTm = *primOptTm
	pa.primMaxTm = *primMaxTm
	pa.inMinTm = *inMinTm
	pa.inOptTm = *inOptTm
	pa.inMaxTm = *inMaxTm
	files := flag.Args()
	clio.ParseFiles(files, parse, pa)
}
