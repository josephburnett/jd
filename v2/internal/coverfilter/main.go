// Command coverfilter reads a Go coverage profile and rewrites blocks
// annotated with //jd:nocover as covered. This allows excluding specific
// unreachable branches from coverage enforcement while keeping the rest
// of the function under the 100% coverage requirement.
//
// Usage: go run internal/coverfilter/main.go coverage.out
//
// Filtered profile is written to stdout.
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var blockRe = regexp.MustCompile(`^(.+):(\d+)\.\d+,(\d+)\.\d+ \d+ (\d+)$`)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: coverfilter <coverprofile>\n")
		os.Exit(1)
	}

	module := readModule()
	prefix := module + "/"

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "coverfilter: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	sources := make(map[string][]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m := blockRe.FindStringSubmatch(line)
		if m == nil {
			fmt.Println(line) // mode line
			continue
		}
		count, _ := strconv.Atoi(m[4])
		if count == 0 {
			startLine, _ := strconv.Atoi(m[2])
			endLine, _ := strconv.Atoi(m[3])
			relPath := strings.TrimPrefix(m[1], prefix)
			if hasNocover(sources, relPath, startLine, endLine) {
				line = line[:len(line)-1] + "1"
			}
		}
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "coverfilter: %v\n", err)
		os.Exit(1)
	}
}

func readModule() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Fprintf(os.Stderr, "coverfilter: cannot read go.mod: %v\n", err)
		os.Exit(1)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if mod, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(mod)
		}
	}
	fmt.Fprintf(os.Stderr, "coverfilter: no module directive in go.mod\n")
	os.Exit(1)
	return ""
}

func loadFile(cache map[string][]string, path string) []string {
	if lines, ok := cache[path]; ok {
		return lines
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	cache[path] = lines
	return lines
}

func hasNocover(cache map[string][]string, path string, start, end int) bool {
	lines := loadFile(cache, path)
	if lines == nil {
		return false
	}
	// Check from one line before the block through the block's end.
	// This finds annotations on the case/default line, on the line
	// above the block, or within the block itself.
	lo := max(start-2, 0) // 1-indexed to 0-indexed, then one line before
	hi := min(end, len(lines))
	for i := lo; i < hi; i++ {
		if strings.Contains(lines[i], "//jd:nocover") {
			return true
		}
	}
	return false
}
