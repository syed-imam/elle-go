package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/syed-imam/elle-go/anomaly"
	"github.com/syed-imam/elle-go/checker"
	"github.com/syed-imam/elle-go/history"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "pg" {
		return runPG(args[1:], stdout, stderr)
	}
	in := stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		defer f.Close()
		in = f
	}
	var h history.History
	if err := json.NewDecoder(in).Decode(&h); err != nil {
		fmt.Fprintln(stderr, "invalid history:", err)
		return 2
	}
	v, err := checker.Check(h)
	if err != nil {
		fmt.Fprintln(stderr, "invalid history:", err)
		return 2
	}
	fmt.Fprintln(stdout, v)
	if v.Anomaly != anomaly.None {
		return 1
	}
	return 0
}
