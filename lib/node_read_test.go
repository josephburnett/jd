package jd

import (
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
	node, err := ReadJsonString(s)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	if !n.Equals(node) {
		t.Errorf("ReadJsonString(%v) = %v. Want %v.", s, node, n)
	}
	node, err = ReadYamlString(s)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	if !n.Equals(node) {
		t.Errorf("ReadYamlString(%v) = %v. Want %v.", s, node, n)
	}
}

func TestReadDiff(t *testing.T) {
	checkReadDiff(t,
		Diff{
			DiffElement{
				Path:      p("a"),
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
				Path:      p("a", 1.0, "b"),
				OldValues: []JsonNode{jsonNumber(1)},
				NewValues: []JsonNode{jsonNumber(2)},
			},
		},
		`@ ["a",1,"b"]`,
		`- 1`,
		`+ 2`)
	checkReadDiff(t,
		Diff{
			DiffElement{
				Path:      p(),
				OldValues: []JsonNode{jsonNumber(1)},
				NewValues: []JsonNode{jsonNumber(2)},
			},
			DiffElement{
				Path:      p(),
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
	want := ""
	for _, dl := range diffLines {
		want += dl + "\n"
	}
	actual, err := readDiff(want)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	got := actual.Render()
	if got != want {
		t.Errorf("readDiff got %v. Want %v.", got, want)
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

func p(elements ...interface{}) path {
	var path path
	for _, e := range elements {
		n, err := NewJsonNode(e)
		if err != nil {
			panic(err)
		}
		path = append(path, n)
	}
	return path
}
