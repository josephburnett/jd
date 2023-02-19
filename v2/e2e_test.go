package jd

import (
	"testing"
)

func TestDiffAndPatch(t *testing.T) {
	checkDiffAndPatchSuccess(t,
		`{"a":1}`,
		`{"a":2}`,
		`{"a":1,"c":3}`,
		`{"a":2,"c":3}`)
	checkDiffAndPatchSuccess(t,
		`[[]]`,
		`[[1]]`,
		`[[],[2]]`,
		`[[1],[2]]`)
	checkDiffAndPatchSuccess(t,
		`[{"a":1},{"a":1}]`,
		`[{"a":2},{"a":3}]`,
		`[{"a":1},{"a":1,"b":4},{"c":5}]`,
		`[{"a":2},{"a":3,"b":4},{"c":5}]`)
}

func TestDiffAndPatchSet(t *testing.T) {
	checkDiffAndPatchSuccessSet(t,
		`{"a":{"b" : ["3", "4" ],"c" : ["2", "1"]}}`,
		`{"a":{"b" : ["3", "4", "5", "6"],"c" : ["2", "1"]}}`,
		`{"a":{"b" : ["3", "4" ],"c" : ["2", "1"]}}`,
		`{"a":{"b" : ["3", "4", "5", "6"],"c" : ["2", "1"]}}`)
}

func TestDiffAndPatchError(t *testing.T) {
	checkDiffAndPatchError(t,
		`{"a":1}`,
		`{"a":2}`,
		`{"a":3}`)
	checkDiffAndPatchError(t,
		`{"a":1}`,
		`{"a":2}`,
		`{}`)
	checkDiffAndPatchError(t,
		`1`,
		`2`,
		``)
	checkDiffAndPatchError(t,
		`1`,
		``,
		`2`)
}

type format string

const (
	formatJd    format = "jd"
	formatPatch format = "patch"
	formatMerge format = "merge"
)

func checkDiffAndPatchSuccessSet(t *testing.T, a, b, c, expect string) {
	err := checkDiffAndPatch(t, formatJd, a, b, c, expect, SET)
	if err != nil {
		t.Errorf("error round-tripping jd format: %v", err)
	}
	// JSON Patch format does not support sets.
}

func checkDiffAndPatchSuccess(t *testing.T, a, b, c, expect string) {
	err := checkDiffAndPatch(t, formatJd, a, b, c, expect)
	if err != nil {
		t.Errorf("error round-tripping jd format: %v", err)
	}
	err = checkDiffAndPatch(t, formatPatch, a, b, c, expect)
	if err != nil {
		t.Errorf("error round-tripping patch format: %v", err)
	}
}

func checkDiffAndPatchError(t *testing.T, a, b, c string) {
	err := checkDiffAndPatch(t, formatJd, a, b, c, "")
	if err == nil {
		t.Errorf("expected error round-tripping jd format")
	}
	err = checkDiffAndPatch(t, formatPatch, a, b, c, "")
	if err == nil {
		t.Errorf("expected error rount-tripping patch format")
	}
}

func checkDiffAndPatch(t *testing.T, f format, a, b, c, expect string, metadata ...Metadata) error {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		return err
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		return err
	}
	nodeC, err := ReadJsonString(c)
	if err != nil {
		return err
	}
	expectNode, err := ReadJsonString(expect)
	if err != nil {
		return err
	}
	var diff Diff
	switch f {
	case formatJd:
		diffString := nodeA.Diff(nodeB).Render()
		diff, err = ReadDiffString(diffString)
	case formatPatch:
		patchString, err := nodeA.Diff(nodeB).RenderPatch()
		if err != nil {
			return nil
		}
		diff, err = ReadPatchString(patchString)
		if err != nil {
			return err
		}
	case formatMerge:
		// not yet implemented
	}
	if err != nil {
		return err
	}
	actualNode, err := nodeC.Patch(diff)
	if err != nil {
		return err
	}
	if !actualNode.Equals(expectNode) {
		t.Errorf("actual = %v. Want %v.", actualNode, expectNode)
	}
	return nil
}
