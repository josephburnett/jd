package jd

import (
	"testing"
)

func TestDiffAndPatch(t *testing.T) {
	cases := []struct {
		a      string
		b      string
		c      string
		expect string
	}{{
		a:      `{"a":1}`,
		b:      `{"a":2}`,
		c:      `{"a":1,"c":3}`,
		expect: `{"a":2,"c":3}`,
	}, {
		a:      `[[]]`,
		b:      `[[1]]`,
		c:      `[[],[2]]`,
		expect: `[[1],[2]]`,
	}, {
		a:      `[{"a":1},{"a":1}]`,
		b:      `[{"a":2},{"a":3}]`,
		c:      `[{"a":1},{"a":1,"b":4},{"c":5}]`,
		expect: `[{"a":2},{"a":3,"b":4},{"c":5}]`,
	}}
	for _, c := range cases {
		t.Run(c.a+c.b+c.c+c.expect, func(t *testing.T) {
			checkDiffAndPatchSuccess(
				t,
				c.a,
				c.b,
				c.c,
				c.expect,
			)
		})
	}
}

func TestDiffAndPatchSet(t *testing.T) {
	checkDiffAndPatchSuccessSet(t,
		`{"a":{"b" : ["3", "4" ],"c" : ["2", "1"]}}`,
		`{"a":{"b" : ["3", "4", "5", "6"],"c" : ["2", "1"]}}`,
		`{"a":{"b" : ["3", "4" ],"c" : ["2", "1"]}}`,
		`{"a":{"b" : ["3", "4", "5", "6"],"c" : ["2", "1"]}}`)
}

func TestDiffAndPatchError(t *testing.T) {
	cases := []struct {
		a string
		b string
		c string
	}{{
		a: `{"a":1}`,
		b: `{"a":2}`,
		c: `{"a":3}`,
	}, {
		a: `{"a":1}`,
		b: `{"a":2}`,
		c: `{}`,
	}, {
		a: `1`,
		b: `2`,
		c: ``,
	}, {
		a: `1`,
		b: ``,
		c: `2`,
	}}
	for _, c := range cases {
		t.Run(c.a+c.b+c.c, func(t *testing.T) {
			checkDiffAndPatchError(
				t,
				c.a,
				c.b,
				c.c,
			)
		})
	}
}

type format string

const (
	formatJd    format = "jd"
	formatPatch format = "patch"
	formatMerge format = "merge"
)

func checkDiffAndPatchSuccessSet(t *testing.T, a, b, c, expect string) {
	err := checkDiffAndPatch(t, formatJd, a, b, c, expect, setOption{})
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

func checkDiffAndPatch(t *testing.T, f format, a, b, c, expect string, options ...Option) error {
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
