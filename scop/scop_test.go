package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestScop(t *testing.T) {
	test := exec.Command("./scop", "-d", "../data/sample",
		"-t", "tarTax.txt", "prim.fasta")
	get, err := test.Output()
	if err != nil {
		t.Error(err)
	}
	want, err := os.ReadFile("r.txt")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(get, want) {
		t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
	}
}
