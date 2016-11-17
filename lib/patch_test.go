package jd

import (
	"testing"
)

func TestPatch(t *testing.T) {
	checkPatch(t,
		`{"a":1}`,
		`{"a":2}`,
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkPatch(t,
		`1`,
		`2`,
		`@ []`,
		`- 1`,
		`+ 2`)
	checkPatch(t,
		`{"a":1}`,
		`{}`,
		`@ ["a"]`,
		`- 1`)
}
