package differential

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"elle-go/history"
)

type ElleResult struct {
	ValidField   string
	AnomalyTypes []string
}

func (r ElleResult) Violation() bool { return r.ValidField == "false" }

func elleDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "..", "elle-upstream")
}

func ElleAvailable() bool {
	if _, err := exec.LookPath("lein"); err != nil {
		return false
	}
	if _, err := os.Stat(elleDir()); err != nil {
		return false
	}
	return true
}

var (
	validRe   = regexp.MustCompile(`:valid\?\s+([^,}\s]+)`)
	anomalyRe = regexp.MustCompile(`:anomaly-types\s+\[([^\]]*)\]`)
)

func RunElle(h history.History) (ElleResult, error) {
	f, err := os.CreateTemp("", "elle-history-*.edn")
	if err != nil {
		return ElleResult{}, err
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(ToEDN(h)); err != nil {
		f.Close()
		return ElleResult{}, err
	}
	f.Close()

	shim, err := os.Getwd()
	if err != nil {
		return ElleResult{}, err
	}
	cmd := exec.Command("lein", "update-in", ":source-paths", "conj",
		strconv.Quote(shim), "--", "run", "-m", "elle-bridge", f.Name())
	cmd.Dir = elleDir()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ElleResult{}, fmt.Errorf("lein run failed: %v\n%s", err, out)
	}
	return parseElle(string(out))
}

func parseElle(out string) (ElleResult, error) {
	var r ElleResult
	m := validRe.FindStringSubmatch(out)
	if m == nil {
		return r, fmt.Errorf("no :valid? found in Elle output:\n%s", out)
	}
	r.ValidField = m[1]
	if a := anomalyRe.FindStringSubmatch(out); a != nil {
		for _, tok := range strings.Fields(a[1]) {
			r.AnomalyTypes = append(r.AnomalyTypes, strings.TrimPrefix(tok, ":"))
		}
	}
	return r, nil
}
