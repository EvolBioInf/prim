package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/prim/util"
	"io"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

var penaltyRE = regexp.MustCompile(
	`PRIMER_PAIR_[0-9]+_PENALTY`)
var forwardRE = regexp.MustCompile(
	`PRIMER_LEFT_[0-9]+_SEQUENCE`)
var reverseRE = regexp.MustCompile(
	`PRIMER_RIGHT_[0-9]+_SEQUENCE`)
var internalRE = regexp.MustCompile(
	`PRIMER_INTERNAL_[0-9]+_SEQUENCE`)

func scan(r io.Reader, args ...interface{}) {
	w := args[0].(*tabwriter.Writer)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fields := strings.Split(sc.Text(), "=")
		if len(fields) == 2 {
			if penaltyRE.MatchString(fields[0]) {
				fmt.Fprintf(w, "%s", fields[1])
			}
			if forwardRE.MatchString(fields[0]) {
				fmt.Fprintf(w, "\t%s", fields[1])
			}
			if reverseRE.MatchString(fields[0]) {
				fmt.Fprintf(w, "\t%s", fields[1])
			}
			if internalRE.MatchString(fields[0]) {
				fmt.Fprintf(w, "\t%s\n", fields[1])
			}
		}
	}
}
func main() {
	util.SetName("prim2tab")
	u := "prim2tab [option]... [file]..."
	p := "Convert output of primer3 to table."
	e := "primt2tab primer3.out | sort -n"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.Version()
	}
	files := flag.Args()
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	fmt.Fprintf(w, "# Penalty\tForward\tReverse\tInternal\n")
	clio.ParseFiles(files, scan, w)
	w.Flush()
}
