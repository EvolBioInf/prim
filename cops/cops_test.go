package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestCops(t *testing.T) {
	var tests []*exec.Cmd
	d := "../data/sample"
	r := "AE005174"
	x := "2e-3"
	i := "scop.out"
	test := exec.Command("./cops", "-d", d, "-r", r, "-t", x, i)
	tests = append(tests, test)
	test = exec.Command("./cops", "-d", d, "-r", r, "-t", x,
		"-D", i)
	tests = append(tests, test)
	test = exec.Command("./cops", "-d", d, "-r", r, "-t", x,
		"-D", "-p", i)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := os.ReadFile(f)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
