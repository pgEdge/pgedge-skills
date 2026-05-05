package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	profile := flag.String("profile", "coverage.out", "Path to go cover profile")
	threshold := flag.Float64("threshold", 90.0, "Minimum % required")
	excludeCSV := flag.String("exclude", "", "Comma-separated path suffixes to exclude")
	flag.Parse()

	var excludes []string
	if *excludeCSV != "" {
		excludes = strings.Split(*excludeCSV, ",")
	}
	pct, err := computeCoverage(*profile, excludes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	fmt.Printf("Coverage: %.1f%% (threshold %.1f%%)\n", pct, *threshold)
	if pct+0.0001 < *threshold {
		os.Exit(1)
	}
}

func computeCoverage(path string, excludes []string) (float64, error) {
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var total, covered int64
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		file := line[:colon]
		skip := false
		for _, e := range excludes {
			if strings.HasSuffix(file, e) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		stmts, err := strconv.ParseInt(fields[len(fields)-2], 10, 64)
		if err != nil {
			continue
		}
		count, err := strconv.ParseInt(fields[len(fields)-1], 10, 64)
		if err != nil {
			continue
		}
		total += stmts
		if count > 0 {
			covered += stmts
		}
	}
	if err := sc.Err(); err != nil {
		return 0, err
	}
	if total == 0 {
		return 100, nil
	}
	return float64(covered) * 100 / float64(total), nil
}
