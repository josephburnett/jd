package jd

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	checkUnmarshal(t, ``, voidNode{})
	checkUnmarshal(t, `null`, jsonNull{})
	checkUnmarshal(t, `true`, jsonBool(true))
	checkUnmarshal(t, `"a"`, jsonString("a"))
	checkUnmarshal(t, `1.0`, jsonNumber(1.0))
	checkUnmarshal(t, `1`, jsonNumber(1.0))
	checkUnmarshal(t, `{}`, jsonObject{})
	checkUnmarshal(t, `[]`, jsonArray{})
}

func checkUnmarshal(t *testing.T, s string, n JsonNode) {
	node, err := unmarshal([]byte(s))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !n.Equals(node) {
		t.Errorf("unmarshal(%v) = %v. Want %v.", s, node, n)
	}
}

func TestReadDiff(t *testing.T) {
	checkReadDiff(t,
		Diff{
			DiffElement{
				Path:      Path{"a"},
				OldValues: []JsonNode{jsonNumber(1)},
				NewValues: []JsonNode{jsonNumber(2)},
			},
		},
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkReadDiff(t,
		Diff{
			DiffElement{
				Path:      Path{"a", 1.0, "b"},
				OldValues: []JsonNode{jsonNumber(1)},
				NewValues: []JsonNode{jsonNumber(2)},
			},
		},
		`@ ["a", 1, "b"]`,
		`- 1`,
		`+ 2`)
	checkReadDiff(t,
		Diff{
			DiffElement{
				Path:      Path{},
				OldValues: []JsonNode{jsonNumber(1)},
				NewValues: []JsonNode{jsonNumber(2)},
			},
			DiffElement{
				Path:      Path{},
				OldValues: []JsonNode{jsonNumber(2)},
				NewValues: []JsonNode{jsonNumber(3)},
			},
		},
		`@ []`,
		`- 1`,
		`+ 2`,
		`@ []`,
		`- 2`,
		`+ 3`)
}

func TestReadDiffError(t *testing.T) {
	checkReadDiffError(t, `- 1`)
	checkReadDiffError(t, `+ 1`)
	checkReadDiffError(t,
		`@ ["a"]`,
		`@ ["b"]`)
	checkReadDiffError(t,
		`@ ["a"]`,
		`+ 1`,
		`- 2`)
	checkReadDiffError(t,
		`@ ["a"]`,
		`- 1`,
		`- 1`)
	checkReadDiffError(t,
		`@ ["a"]`,
		`+ 2`,
		`+ 2`)
	checkReadDiffError(t,
		`@ ["a"]`,
		`- 1`,
		`@ ["b"]`)
	checkReadDiffError(t,
		`@ `,
		`- 1`)
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

func checkReadDiffError(t *testing.T, diffLines ...string) {
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	actual, err := readDiff(diff)
	if actual != nil {
		t.Errorf("readDiff(%v) = %v. Want nil.", diff, actual)
	}
	if err == nil {
		t.Errorf("Expected error for readDiff(%v).", diff)
	}
}
