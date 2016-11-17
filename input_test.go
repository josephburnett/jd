package jd

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	checkUnmarshal(t, ``, voidNode{})
	checkUnmarshal(t, `"a"`, jsonString("a"))
	checkUnmarshal(t, `1.0`, jsonNumber(1.0))
	checkUnmarshal(t, `1`, jsonNumber(1.0))
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
