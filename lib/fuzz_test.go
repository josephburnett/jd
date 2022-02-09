// +build test_fuzz

package jd

import (
	"testing"
)

func FuzzJd(f *testing.F) {
	a := `{"foo":{},"bar":[1,null,3]}`
	b := a
	f.Add(a, b)
	f.Fuzz(fuzz)
}
