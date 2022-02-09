package jd

import "testing"

// TestFuzzBackports is a backport of fuzz testing failures to golang
// versions before 1.18. This is only necessary until go1.18beta2 is GA
// and we can depend on it.
//
// To run fuzzing use `go1.18beta2 test -tags=test_fuzz ./lib/ -fuzz=FuzzJd`
func TestFuzzBackport(t *testing.T) {
	for _, backport := range [][2]string{{
		"[]", "0", // FuzzJd/e193f6c4bfd5b8d3c12e1ac42162b2ccd7a31f9aafd466066c1ec7a95da48e1e
	}, {
		"{}", " ", // FuzzJd/868060b2021521d32933f40415c6f95b38fda5f5c6bdb7fa6664d046c637c03c
	}} {
		fuzz(t, backport[0], backport[1])
	}
}
