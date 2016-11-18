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

func checkDiffAndPatchSuccess(t *testing.T, a, b, c, expect string) {
	err := checkDiffAndPatch(t, a, b, c, expect)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func checkDiffAndPatchError(t *testing.T, a, b, c string) {
	err := checkDiffAndPatch(t, a, b, c, "")
	if err == nil {
		t.Errorf("Expected error.")
	}
}

func checkDiffAndPatch(t *testing.T, a, b, c, expect string) error {
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
	diffString := nodeA.Diff(nodeB).Render()
	diff, err := ReadDiffString(diffString)
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
