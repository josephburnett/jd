package jd

import "testing"

// TestFuzzBackports is a backport of fuzz testing failures to golang
// versions before 1.18. This is only necessary until go1.18beta2 is GA
// and we can depend on it.
//
// To run fuzzing use `go1.18beta2 test -tags=test_fuzz ./lib/-run=FuzzJd`
func TestFuzzBackport(t *testing.T) {
	for _, backport := range [][2]string{{
		"[]", "0",
	}} {
		fuzz(t, backport[0], backport[1])
	}
}
