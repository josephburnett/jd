package jd

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	checkNodeDiff(t, `{"a":1}`, `{"a":2}`,
		Diff{DiffElement{Path{"a"}, jsonNumber(1.0), jsonNumber(2.0)}})
	checkNodeDiff(t, `{"a":1}`, `{}`,
		Diff{DiffElement{Path{"a"}, jsonNumber(1.0), nil}})
	checkNodeDiff(t, `{}`, `{"a":2}`,
		Diff{DiffElement{Path{"a"}, nil, jsonNumber(2.0)}})
	checkNodeDiff(t, `{"a":1}`, `{"a":1}`, Diff{})
	checkNodeDiff(t, `{"a":{"b":1}}`, `{"a":{"c":2}}`,
		Diff{
			DiffElement{Path{"a", "b"}, jsonNumber(1.0), nil},
			DiffElement{Path{"a", "c"}, nil, jsonNumber(2.0)}})
	checkNodeDiff(t, `{"a":[1,2]}`, `{"a":[2,1]}`,
		Diff{
			DiffElement{Path{"a", 1.0}, jsonNumber(2.0), jsonNumber(1.0)},
			DiffElement{Path{"a", 0.0}, jsonNumber(1.0), jsonNumber(2.0)}})
	checkNodeDiff(t, `{"a":[1]}`, `{"a":[1,2]}`,
		Diff{DiffElement{Path{"a", 1.0}, nil, jsonNumber(2.0)}})
	checkNodeDiff(t, `{"a":[1,2]}`, `{"a":[1]}`,
		Diff{DiffElement{Path{"a", 1.0}, jsonNumber(2.0), nil}})
	checkNodeDiff(t, `{"a":[1,2]}`, `{"a":[3,4]}`,
		Diff{
			DiffElement{Path{"a", 1.0}, jsonNumber(2.0), jsonNumber(4.0)},
			DiffElement{Path{"a", 0.0}, jsonNumber(1.0), jsonNumber(3.0)}})
	checkNodeDiff(t, `{"a":[{"b":1}]}`, `{"a":[{"b":2}]}`,
		Diff{DiffElement{Path{"a", 0.0, "b"}, jsonNumber(1.0), jsonNumber(2.0)}})
}

func checkNodeDiff(t *testing.T, a, b string, diff Diff) {
	jsonA, err := unmarshal([]byte(a))
	if err != nil {
		t.Error(err.Error())
	}
	jsonB, err := unmarshal([]byte(b))
	if err != nil {
		t.Error(err.Error())
	}
	path := make(Path, 0)
	d := jsonA.diff(jsonB, path)
	if !reflect.DeepEqual(d, diff) {
		t.Errorf("Got %v. Want %v.", d, diff)
	}
}

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
