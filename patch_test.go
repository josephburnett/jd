package jd

import (
	"testing"
)

func TestPatch(t *testing.T) {
	checkPatch(t,
		`{"a":{"b":1}}`,
		`{"a":{"c":2}}`,
		`@ ["a"]`,
		`- {"b":1}`,
		`+ {"c":2}`)
}

func checkPatch(t *testing.T, a, e string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a)
	if err != nil {
		t.Errorf(err.Error())
	}
	diff, err := ReadDiffString(diffString)
	if err != nil {
		t.Errorf(err.Error())
	}
	expect, err := ReadJsonString(e)
	if err != nil {
		t.Errorf(err.Error())
	}
	b, err := initial.Patch(diff)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !expect.Equals(b) {
		t.Errorf("%v.Patch(%v) = %v. Want %v.",
			a, diffLines, renderJson(b), e)
	}
}
