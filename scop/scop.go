package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/prim/util"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Hit struct {
	length, qlen, mismatch int
	saccver                string
	sstart, send           int
}
type HitSlice []*Hit

func (h HitSlice) Len() int {
	return len(h)
}
func (h HitSlice) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h HitSlice) Less(i, j int) bool {
	if h[i].saccver == h[j].saccver {
		return h[i].sstart < h[j].sstart
	}
	return h[i].saccver < h[j].saccver
}
func main() {
	util.SetName("scop")
	u := "scop -d <db> -t <taxids.txt> [option]... [foo.fa]..."
	p := "Score primers by comparing them to a Blast database."
	e := "scop -d nt -t targets.txt prim.fa"
	clio.Usage(u, p, e)
	optD := flag.String("d", "", "Blast database")
	optT := flag.String("t", "", "file of target taxon IDs")
	optN := flag.String("n", "", "file of negative taxon IDs "+
		"(-negative_taxidlist for blastn)")
	optP := flag.String("p", "", "file of positive taxon IDs "+
		"(-taxidlist for blastn)")
	optI := flag.Int("i", 5, "maximum number of mismatches")
	optL := flag.Int("l", 4000, "maximum length of amplicon")
	optE := flag.Float64("e", 1000.0, "E-value")
	nt := runtime.NumCPU()
	optTT := flag.Int("T", nt, "number of threads (default CPUs)")
	optM := flag.Float64("m", 3.0, "set the maximum number "+
		"of target sequences in Blast to the expected "+
		"number of accessions times -m")
	optV := flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.Version()
	}
	if *optD == "" {
		log.Fatal("please supply a Blast datbase")
	} else {
		fn := *optD + ".ndb"
		_, err := os.Stat(fn)
		if err != nil {
			m := "couldn't find Blast database %q\n"
			log.Fatalf(m, *optD)
		}
	}
	if *optT == "" {
		log.Fatal("please supply a file of target taxon IDs")
	} else {
		_, err := os.Stat(*optT)
		if err != nil {
			m := "couldn't find file %q\n"
			log.Fatalf(m, *optT)
		}
	}
	etacc := make(map[string]bool)
	ptacc := make(map[string]bool)
	tmpl := "blastdbcmd -db %s -taxidlist %s -outfmt "
	cs := fmt.Sprintf(tmpl, *optD, *optT)
	args := strings.Fields(cs)
	args = append(args, "%a %t")
	cmd := exec.Command("blastdbcmd")
	cmd.Args = args
	out, err := cmd.CombinedOutput()
	util.Check(err)
	cg := []byte("complete genome")
	pl := []byte("plasmid")
	entries := bytes.Split(out, []byte("\n"))
	for _, entry := range entries {
		acc := ""
		if len(entry) > 0 {
			acc = string(bytes.Fields(entry)[0])
			ptacc[acc] = true
		}
		if bytes.Contains(entry, cg) &&
			!bytes.Contains(entry, pl) {
			etacc[acc] = true
		}
	}
	if *optTT == 0 {
		(*optTT) = runtime.NumCPU()
	}
	primerSets := flag.Args()
	if len(primerSets) == 0 {
		ps, err := os.CreateTemp("", "prim*.fasta")
		util.Check(err)
		defer ps.Close()
		defer os.Remove(ps.Name())
		sc := fasta.NewScanner(os.Stdin)
		for sc.ScanSequence() {
			seq := sc.Sequence()
			fmt.Fprintf(ps, "%s\n", seq)
		}
		primerSets = append(primerSets, ps.Name())
	}
	for _, primerSet := range primerSets {
		tmpl = "blastn -task blastn-short -query %s -db %s -evalue " +
			"%g -num_threads %d -max_target_seqs %d"
		mts := int(*optM * float64(len(ptacc)))
		cs = fmt.Sprintf(tmpl, primerSet, *optD, *optE, *optTT, mts)
		if *optN != "" && *optP != "" {
			log.Fatal("please use either a positive or " +
				"negative taxon ID list")
		}

		if *optN != "" {
			tmpl = "-negative_taxidlist %s"
			cs = cs + " " + fmt.Sprintf(tmpl, *optN)
		}

		if *optP != "" {
			tmpl = "-taxidlist %s"
			cs = cs + " " + fmt.Sprintf(tmpl, *optP)
		}
		args = strings.Fields(cs)
		args = append(args, "-outfmt", "6 length qlen mismatch "+
			"saccver sstart send")
		cmd = exec.Command("blastn")
		cmd.Args = args
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Fatal(string(out))
		}
		otacc := make(map[string]bool)
		hits := make([]*Hit, 0)
		lines := bytes.Split(out, []byte("\n"))
		for _, line := range lines {
			fields := bytes.Fields(line)
			if len(fields) == 6 {
				hit := new(Hit)
				hit.length, err = strconv.Atoi(string(fields[0]))
				util.Check(err)
				hit.qlen, err = strconv.Atoi(string(fields[1]))
				util.Check(err)
				hit.mismatch, err = strconv.Atoi(string(fields[2]))
				util.Check(err)
				hit.saccver = string(fields[3])
				hit.sstart, err = strconv.Atoi(string(fields[4]))
				util.Check(err)
				hit.send, err = strconv.Atoi(string(fields[5]))
				util.Check(err)
				hits = append(hits, hit)
			}
		}
		i := 0
		for _, hit := range hits {
			if hit.qlen == hit.length &&
				hit.mismatch <= *optI {
				hits[i] = hit
				i++
			}
		}
		hits = hits[:i]
		sort.Sort(HitSlice(hits))
		for i, hit := range hits {
			if !otacc[hit.saccver] && hit.sstart < hit.send {
				fp := hit
				for j := i + 1; j < len(hits); j++ {
					if hits[j].sstart > hits[j].send &&
						fp.saccver == hits[j].saccver {
						rp := hits[j]
						if rp.send-fp.sstart+1 <= *optL {
							otacc[fp.saccver] = true
							break
						}
					}
				}
			}
		}
		tp := 0
		truePositives := make([]string, 0)
		for o, _ := range otacc {
			if ptacc[o] {
				truePositives = append(truePositives, o)
				tp++
			}
		}
		fp := 0
		falsePositives := make([]string, 0)
		for o, _ := range otacc {
			if !ptacc[o] {
				fp++
				falsePositives = append(falsePositives, o)
			}
		}
		fn := 0
		falseNegatives := make([]string, 0)
		for e, _ := range etacc {
			if !otacc[e] {
				fn++
				falseNegatives = append(falseNegatives, e)
			}
		}
		sn := float64(tp) / (float64(tp) + float64(fn))
		sp := float64(tp) / (float64(tp) + float64(fp))
		fmt.Printf("PrimerSet:\t%s\n", primerSet)
		fmt.Printf("Sensitivity:\t%.3g\n", sn)
		fmt.Printf("Specificity:\t%.3g\n", sp)
		if len(truePositives) > 0 {
			sort.Strings(truePositives)
			fmt.Printf("TruePositives:\t%s", truePositives[0])
			for i := 1; i < tp; i++ {
				fmt.Printf(" %s", truePositives[i])
			}
			fmt.Printf("\n")
		}
		if len(falsePositives) > 0 {
			sort.Strings(falsePositives)
			fmt.Printf("FalsePositives:\t%s", falsePositives[0])
			for i := 1; i < fp; i++ {
				fmt.Printf(" %s", falsePositives[i])
			}
			fmt.Printf("\n")
		}
		if len(falseNegatives) > 0 {
			sort.Strings(falseNegatives)
			fmt.Printf("FalseNegatives:\t%s", falseNegatives[0])
			for i := 1; i < fn; i++ {
				fmt.Printf(" %s", falseNegatives[i])
			}
			fmt.Printf("\n")
		}
	}
}
