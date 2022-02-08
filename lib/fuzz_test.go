// +build test_fuzz

package jd

import (
	"testing"
)

func FuzzJd(f *testing.F) {
	a := `{"foo":{},"bar":[1,null,3]}`
	b := a
	f.Add(a, b)
	f.Fuzz(func(t *testing.T, aStr, bStr string) {
		// Only valid JSON input.
		a, err := ReadJsonString(aStr)
		if err != nil {
			return
		}
		if a == nil {
			t.Errorf("nil parsed input: %v", aStr)
			return
		}
		b, err := ReadJsonString(bStr)
		if err != nil {
			return
		}
		if b == nil {
			t.Errorf("nil parsed input: %v", bStr)
			return
		}
		diffAB := a.Diff(b)
		if diffAB == nil {
			t.Errorf("nil diff of a and b")
			return
		}
	})
}
