package jd

import (
	"testing"
)

func TestDiffRender(t *testing.T) {
	checkDiffRender(t, `{"a":1}`, `{"a":2}`,
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkDiffRender(t, `{"a":{"b":1}}`, `{"a":{"b":2}}`,
		`@ ["a","b"]`,
		`- 1`,
		`+ 2`)
	checkDiffRender(t, `{"a":{"b":1}}`, `{"a":{"c":2}}`,
		`@ ["a","b"]`,
		`- 1`,
		`@ ["a","c"]`,
		`+ 2`)
	checkDiffRender(t, `{"a":{"b":1}}`, `{"c":{"b":1}}`,
		`@ ["a"]`,
		`- {"b":1}`,
		`@ ["c"]`,
		`+ {"b":1}`)
}

func checkDiffRender(t *testing.T, a, b string, diffLines ...string) {
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	aJson, err := unmarshal([]byte(a))
	if err != nil {
		t.Errorf(err.Error())
	}
	bJson, err := unmarshal([]byte(b))
	if err != nil {
		t.Errorf(err.Error())
	}
	d := aJson.diff(bJson, Path{}).Render()
	if d != diff {
		t.Errorf("%v.diff(%v) = %v. Want %v.", a, b, d, diff)
	}
}
