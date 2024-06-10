package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/dist"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/prim/util"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func parse(r io.Reader, args ...interface{}) {
	db := args[0].(string)
	re := args[1].(string)
	dt := args[2].(float64)
	nt := args[3].(int)
	pd := args[4].(bool)
	truePos := args[5].(bool)
	d, e := io.ReadAll(r)
	util.Check(e)
	reports := bytes.Split(d, []byte("PrimerSet"))
	reports = reports[1:]
	for _, report := range reports {
		lines := bytes.Split(report, []byte("\n"))
		if len(lines[len(lines)-1]) == 0 {
			lines = lines[:len(lines)-1]
		}
		l := len(lines)
		if l < 3 || l > 6 {
			log.Fatalf("mal-formed report:\n%s\n", string(report))
		}
		primerSet := string(bytes.Fields(lines[0])[1])
		lines = lines[3:]
		accessions := make(map[string][]string)
		for _, line := range lines {
			fields := bytes.Fields(line)
			if len(fields) > 0 {
				k := string(fields[0][:len(fields[0])-1])
				for i := 1; i < len(fields); i++ {
					v := string(fields[i])
					accessions[k] = append(accessions[k],
						v)
				}
			}
		}
		td, err := os.MkdirTemp("", "temp*")
		util.Check(err)
		defer os.RemoveAll(td)
		f, err := os.CreateTemp(td, "acc*.txt")
		util.Check(err)
		defer f.Close()
		cmd := exec.Command("blastdbcmd", "-db", db, "-entry", re,
			"-outfmt", "%a")
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("%s\n", out)
		}
		re = string(out[:len(out)-1])
		fmt.Fprintf(f, "%s\n", re)
		keys := []string{"FalsePositives", "FalseNegatives"}
		if truePos {
			keys = append(keys, "TruePositives")
		}
		for _, key := range keys {
			accs := accessions[key]
			if accs != nil {
				for _, acc := range accs {
					fmt.Fprintf(f, "%s\n", acc)
				}
			}
		}
		cmd = exec.Command("blastdbcmd", "-db", db,
			"-entry_batch", f.Name())
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("%s\n", out)
		}
		sc := fasta.NewScanner(bytes.NewReader(out))
		for sc.ScanSequence() {
			seq := sc.Sequence()
			fn := strings.Fields(seq.Header())[0]
			f, err := os.Create(td + "/" + fn + ".fasta")
			util.Check(err)
			fmt.Fprintf(f, "%s\n", seq)
			f.Close()
		}
		cmd = exec.Command("phylonium")
		nts := strconv.Itoa(nt)
		args := []string{"phylonium", "-t", nts, "-r",
			td + "/" + re + ".fasta"}
		p, err := filepath.Glob(td + "/*.fasta")
		util.Check(err)
		args = append(args, p...)
		cmd.Args = args
		out, _ = cmd.Output()
		var mat *dist.DistMat
		r := bytes.NewReader(out)
		scanner := dist.NewScanner(r)
		if scanner.Scan() {
			mat = scanner.DistanceMatrix()
		} else {
			log.Fatal("couldn't read distance matrix")
		}
		accMap := make(map[string]int)
		for i, name := range mat.Names {
			accMap[name] = i
		}
		ri := accMap[re]
		ntp := make([]string, 0)
		n := 0
		accs := accessions["FalsePositives"]
		for _, acc := range accs {
			j := accMap[acc]
			di := mat.Matrix[ri][j]
			if di > dt || math.IsNaN(di) {
				accs[n] = acc
				n++
			} else {
				ntp = append(ntp, acc)
			}
		}
		accessions["FalsePositives"] = accs[:n]
		n = 0
		accs = accessions["FalseNegatives"]
		for _, acc := range accs {
			i := accMap[re]
			j := accMap[acc]
			di := mat.Matrix[i][j]
			if di <= dt {
				accs[n] = acc
				n++
			}
		}
		accessions["FalseNegatives"] = accs[:n]
		nfp := make([]string, 0)
		if truePos {
			n = 0
			accs = accessions["TruePositives"]
			for _, acc := range accs {
				j := accMap[acc]
				di := mat.Matrix[ri][j]
				if di <= dt {
					accs[n] = acc
					n++
				} else {
					nfp = append(nfp, acc)
				}
			}
			accessions["TruePositives"] = accs[:n]
		}
		tps := accessions["TruePositives"]
		tps = append(tps, ntp...)
		sort.Strings(tps)
		fps := accessions["FalsePositives"]
		fps = append(fps, nfp...)
		sort.Strings(fps)
		fns := accessions["FalseNegatives"]
		sort.Strings(fns)
		tp := float64(len(tps))
		fn := float64(len(fns))
		sn := tp / (tp + fn)
		fp := float64(len(fps))
		sp := tp / (tp + fp)
		fmt.Printf("PrimerSet:\t%s\n", primerSet)
		fmt.Printf("Sensitivity:\t%.3g\n", sn)
		fmt.Printf("Specificity:\t%.3g\n", sp)
		if len(tps) > 0 {
			fmt.Print("TruePositives:\t")
			if pd && truePos {
				printAccDist(tps, re, mat, accMap)
			} else {
				printAcc(tps)
			}
		}
		if len(fps) > 0 {
			fmt.Print("FalsePositives:\t")
			if pd {
				printAccDist(fps, re, mat, accMap)
			} else {
				printAcc(fps)
			}
		}
		if len(fns) > 0 {
			fmt.Print("FalseNegatives:\t")
			if pd {
				printAccDist(fns, re, mat, accMap)
			} else {
				printAcc(fns)
			}
		}
	}
}
func printAccDist(accs []string, re string,
	mat *dist.DistMat,
	accMap map[string]int) {
	dists := mat.Matrix
	i := accMap[re]
	acc := accs[0]
	j := accMap[acc]
	d := dists[i][j]
	fmt.Printf("%s %.3g", acc, d)
	accs = accs[1:]
	for _, acc := range accs {
		j := accMap[acc]
		d := dists[i][j]
		fmt.Printf(" %s %.3g", acc, d)
	}
	fmt.Printf("\n")
}
func printAcc(accs []string) {
	fmt.Printf("%s", accs[0])
	accs = accs[1:]
	for _, acc := range accs {
		fmt.Printf(" %s", acc)
	}
	fmt.Printf("\n")
}
func main() {
	util.SetName("cops")
	u := "cops -d <db> -r <ref> -t <d> [option]... [foo.txt]..."
	p := "Correct primer scores calculated with scop."
	e := "cops -d nt -r AE005174 -t 2e-3 scop.out"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optD := flag.String("d", "", "Blast datbase")
	optR := flag.String("r", "",
		"accession of reference target strain")
	optT := flag.Float64("t", 0, "threshold distance to reference")
	numThreads := runtime.NumCPU()
	optTT := flag.Int("T", numThreads, "number of threads")
	optDD := flag.Bool("D", false, "include distances in output")
	optP := flag.Bool("p", false, "also check true positives "+
		"(default only check false positives and false negatives)")
	flag.Parse()
	if *optV {
		util.Version()
	}
	if *optD == "" {
		log.Fatal("please supply a Blast database")
	}
	if *optR == "" {
		log.Fatal("please supply a reference strain")
	}
	if *optT == 0 {
		log.Fatal("please supply the threshold " +
			"distance to the reference")
	}
	files := flag.Args()
	clio.ParseFiles(files, parse, *optD, *optR,
		(*optT), *optTT, *optDD, *optP)
}
