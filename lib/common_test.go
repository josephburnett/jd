package jd

import (
	"strings"
	"testing"
)

func checkJson(ctx *testContext, a, b string) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	nodeAJson := nodeA.Json(ctx.metadata...)
	if nodeAJson != b {
		ctx.t.Errorf("%v.Json() = %v. Want %v.", nodeA, nodeAJson, b)
	}
}

func checkEqual(ctx *testContext, a, b string) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	if !nodeA.Equals(nodeB, ctx.metadata...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeB.Equals(nodeA, ctx.metadata...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeA.Equals(nodeA, ctx.metadata...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeB.Equals(nodeB, ctx.metadata...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
}

func checkNotEqual(ctx *testContext, a, b string, metadata ...Metadata) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	if nodeA.Equals(nodeB, ctx.metadata...) {
		ctx.t.Errorf("nodeA.Equals(nodeB) == true. Want false.")
	}
	if nodeB.Equals(nodeA, ctx.metadata...) {
		ctx.t.Errorf("nodeB.Equals(nodeA) == true. Want false.")
	}
}

func checkHash(ctx *testContext, a, b string, wantSame bool) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Fatalf(err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Fatalf(err.Error())
	}
	hashA := nodeA.hashCode(ctx.metadata)
	hashB := nodeB.hashCode(ctx.metadata)
	if wantSame && hashA != hashB {
		ctx.t.Errorf("%v.hashCode = %v. %v.hashCode = %v. Want the same.",
			a, hashA, b, hashB)
	}
	if !wantSame && hashA == hashB {
		ctx.t.Errorf("%v.hashCode = %v. %v.hashCode = %v. Want the different.",
			a, hashA, b, hashB)
	}
}

func checkDiff(ctx *testContext, a, b string, diffLines ...string) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Fatalf(err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Fatalf(err.Error())
	}
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	d := nodeA.Diff(nodeB, ctx.metadata...)
	got := d.Render()
	if got != diff {
		ctx.t.Errorf("%v.Diff(%v) = \n%v. Want %v.", nodeA, nodeB, got, diff)
	}
}

func checkPatch(ctx *testContext, a, e string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	diff, err := ReadDiffString(diffString)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	expect, err := ReadJsonString(e)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	b, err := initial.Patch(diff)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	if !expect.Equals(b, ctx.metadata...) {
		ctx.t.Errorf("%v.Patch(%v) = %v. Want %v.",
			a, diffLines, b, e)
	}
}

func checkPatchError(ctx *testContext, a string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	diff, err := ReadDiffString(diffString)
	if err != nil {
		ctx.t.Errorf(err.Error())
	}
	b, err := initial.Patch(diff)
	if b != nil {
		ctx.t.Errorf("%v.Patch(%v) = %v. Want nil.", initial, diff, b)
	}
	if err == nil {
		ctx.t.Errorf("Expected error. Got nil.")
	}
}

func s(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

func ss(s ...string) []string {
	return s
}

func m(m ...Metadata) []Metadata {
	return m
}

type testContext struct {
	t        *testing.T
	metadata []Metadata
}

func newTestContext(t *testing.T) *testContext {
	return &testContext{
		t:        t,
		metadata: make([]Metadata, 0),
	}
}

func (tc *testContext) withMetadata(metadata ...Metadata) *testContext {
	tc.metadata = append(tc.metadata, metadata...)
	return tc
}

func fuzz(t *testing.T, aStr, bStr string) {
	// Only valid JSON input.
	a, err := ReadJsonString(aStr)
	if err != nil {
		return
	}
	if a == nil {
		t.Errorf("nil parsed input: %q", aStr)
		return
	}
	b, err := ReadJsonString(bStr)
	if err != nil {
		return
	}
	if b == nil {
		t.Errorf("nil parsed input: %q", bStr)
		return
	}
	for _, format := range [][2]string{{
		"jd", "list",
	}, {
		"jd", "set",
	}, {
		"jd", "mset",
	}, {
		"patch", "list",
	}} {
		var metadata []Metadata
		switch format[0] {
		case "jd":
			switch format[1] {
			case "set":
				metadata = append(metadata, SET)
			case "mset":
				metadata = append(metadata, MULTISET)
			default: // list
			}
		default: // patch
		}

		// Diff A and B.
		d := a.Diff(b, metadata...)
		if d == nil {
			t.Errorf("nil diff of a and b")
			return
		}
		var diffABStr string
		var diffAB Diff
		switch format[0] {
		case "jd":
			diffABStr = d.Render()
			diffAB, err = ReadDiffString(diffABStr)
		case "patch":
			diffABStr, err = d.RenderPatch()
			if err != nil {
				t.Errorf("could not render diff %v as patch: %v", d, err)
				return
			}
			diffAB, err = ReadPatchString(diffABStr)
		}
		if err != nil {
			t.Errorf("error parsing diff string %q: %v", diffABStr, err)
			return
		}
		// Apply diff to A to get B.
		patchedA, err := a.Patch(diffAB)
		if err != nil {
			t.Errorf("applying patch %v to %v should give %v. Got err: %v", diffAB.Render(), aStr, bStr, err)
			return
		}
		if !patchedA.Equals(b) {
			t.Errorf("applying patch %v to %v should give %v. Got: %v", diffAB.Render(), aStr, bStr, patchedA)
			return
		}
	}

}
