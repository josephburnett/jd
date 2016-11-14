package jd

import (
	"reflect"
	"testing"
)

func TestReadDiff(t *testing.T) {
	checkReadDiff(t,
		Diff{
			DiffElement{
				Path:     Path{"a"},
				OldValue: jsonObject{"b": jsonNumber(1)},
				NewValue: jsonObject{"c": jsonNumber(2)},
			},
		},
		`@ ["a"]`,
		`- {"b":1}`,
		`+ {"c":2}`)
}

func checkReadDiff(t *testing.T, d Diff, diffLines ...string) {
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	actual, err := readDiff(diff)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(d, actual) {
		t.Errorf("readDiff(%v) = %v. Want %v.", diff, actual, d)
	}
}
