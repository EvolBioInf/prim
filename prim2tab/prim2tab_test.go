package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestPrim2tab(t *testing.T) {
	test := exec.Command("./prim2tab", "prim.out")
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
