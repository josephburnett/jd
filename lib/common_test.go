package jd

import (
	"reflect"
	"testing"
)

func checkJson(t *testing.T, a, b string, options ...option) {
	nodeA, err := ReadJsonString(a, options...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeAJson := nodeA.Json()
	if nodeAJson != b {
		t.Errorf("%v.Json() = %v. Want %v.", nodeA, nodeAJson, b)
	}
}

func checkEqual(t *testing.T, a, b string, options ...option) {
	nodeA, err := unmarshal([]byte(a), options...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b), options...)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !nodeA.Equals(nodeB) {
		t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeB.Equals(nodeA) {
		t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeA.Equals(nodeA) {
		t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeB.Equals(nodeB) {
		t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
}

func checkNotEqual(t *testing.T, a, b string, options ...option) {
	nodeA, err := unmarshal([]byte(a), options...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b), options...)
	if err != nil {
		t.Errorf(err.Error())
	}
	if nodeA.Equals(nodeB) {
		t.Errorf("nodeA.Equals(nodeB) == true. Want false.")
	}
	if nodeB.Equals(nodeA) {
		t.Errorf("nodeB.Equals(nodeA) == true. Want false.")
	}
}

func checkHash(t *testing.T, a, b string, wantSame bool) {
	nodeA, err := unmarshal([]byte(a))
	if err != nil {
		t.Fatalf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b))
	if err != nil {
		t.Fatalf(err.Error())
	}
	hashA := nodeA.hashCode()
	hashB := nodeB.hashCode()
	if wantSame && hashA != hashB {
		t.Errorf("%v.hashCode = %v. %v.hashCode = %v. Want the same.",
			a, hashA, b, hashB)
	}
	if !wantSame && hashA == hashB {
		t.Errorf("%v.hashCode = %v. %v.hashCode = %v. Want the different.",
			a, hashA, b, hashB)
	}
}

func checkDiff(t *testing.T, a, b string, diffLines ...string) {
	checkDiffOption(t, "", a, b, diffLines...)
}

func checkDiffOption(t *testing.T, o option, a, b string, diffLines ...string) {
	options := make([]option, 0)
	if o != "" {
		options = append(options, o)
	}
	nodeA, err := ReadJsonString(a, options...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := ReadJsonString(b, options...)
	if err != nil {
		t.Errorf(err.Error())
	}
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	d := nodeA.Diff(nodeB)
	expectedDiff, err := ReadDiffString(diff, o)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(d, expectedDiff) {
		t.Errorf("%v.Diff(%v) = %v. Want %v.", nodeA, nodeB, d, expectedDiff)
	}
}

func checkPatch(t *testing.T, a, e string, diffLines ...string) {
	checkPatchOption(t, "", a, e, diffLines...)
}

func checkPatchOption(t *testing.T, o option, a, e string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a, o)
	if err != nil {
		t.Errorf(err.Error())
	}
	diff, err := ReadDiffString(diffString, o)
	if err != nil {
		t.Errorf(err.Error())
	}
	expect, err := ReadJsonString(e, o)
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

func checkPatchError(t *testing.T, a string, diffLines ...string) {
	checkPatchErrorOption(t, "", a, diffLines...)
}

func checkPatchErrorOption(t *testing.T, o option, a string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a)
	if err != nil {
		t.Errorf(err.Error())
	}
	diff, err := ReadDiffString(diffString, o)
	if err != nil {
		t.Errorf(err.Error())
	}
	b, err := initial.Patch(diff)
	if b != nil {
		t.Errorf("%v.Patch(%v) = %v. Want nil.", initial, diff, b)
	}
	if err == nil {
		t.Errorf("Expected error. Got nil.")
	}
}
