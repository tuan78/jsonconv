package main

import (
	"os"
	"testing"
)

func TestMainFn_GoodArgs(t *testing.T) {
	os.Args = []string{"jsonconv", "--help"}
	main()
}

func TestMainFn_BadArgs(t *testing.T) {
	exitCode := 0
	exitFn = func(code int) { exitCode = 1 }
	os.Args = []string{"jsonconv", "bad-args"}
	main()
	if exitCode == 0 {
		t.Errorf("it should exit with code 1 for bad args")
	}
}
