package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestFa2prim(t *testing.T) {
	var tests []*exec.Cmd
	f := "./test.fasta"
	test := exec.Command("./fa2prim", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-inMaxTm", "48", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-inMinTm", "44", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-inOptTm", "46", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-primMaxTm", "59", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-primMinTm", "55", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-primOptTm", "57", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-primMaxSize", "26", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-primMinSize", "16", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-prodMaxSize", "151", f)
	tests = append(tests, test)
	test = exec.Command("./fa2prim", "-prodMinSize", "71", f)
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
