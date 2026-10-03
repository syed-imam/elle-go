package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name  string
		input string
		code  int
	}{
		{"serializable", `[{"process":1,"type":"ok","mops":[{"type":"append","key":"x","app":1}]},{"process":2,"type":"ok","mops":[{"type":"r","key":"x","read":[1]}]}]`, 0},
		{"write skew", `[{"process":1,"type":"ok","mops":[{"type":"r","key":"y","read":[]},{"type":"append","key":"x","app":1}]},{"process":2,"type":"ok","mops":[{"type":"r","key":"x","read":[]},{"type":"append","key":"y","app":2}]}]`, 1},
		{"malformed json", `[{`, 2},
		{"unwritten read", `[{"process":1,"type":"ok","mops":[{"type":"r","key":"x","read":[99]}]}]`, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			if got := run(nil, strings.NewReader(c.input), &out, &errOut); got != c.code {
				t.Errorf("exit %d, want %d (stdout %q, stderr %q)", got, c.code, out.String(), errOut.String())
			}
		})
	}
}

func TestRunPGBadArgs(t *testing.T) {
	cases := [][]string{
		{"pg", "-isolation", "snapshot", "-dsn", "x"},
		{"pg", "-dsn", ""},
		{"pg", "-nope"},
	}
	for _, args := range cases {
		var out, errOut bytes.Buffer
		if got := run(args, strings.NewReader(""), &out, &errOut); got != 2 {
			t.Errorf("%v: exit %d, want 2 (stderr %q)", args, got, errOut.String())
		}
	}
}
