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
			a, diffLines, b.Json(ctx.metadata...), e)
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

func s(s ...string) []string {
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

func mustParseJson(ss ...string) JsonNode {
	s := strings.Join(ss, "\n")
	n, err := ReadJsonString(s)
	if err != nil {
		panic(err)
	}
	return n
}

func mustParseJsonArray(ss ...string) JsonNode {
	return mustParseJson(ss...).(jsonArray)
}

func mustParseMask(ss ...string) Mask {
	s := strings.Join(ss, "\n")
	m, err := ReadMaskString(s)
	if err != nil {
		panic(err)
	}
	return m
}

func mustParseDiff(ss ...string) Diff {
	s := strings.Join(ss, "\n")
	d, err := ReadDiffString(s)
	if err != err {
		panic(err)
	}
	return d
}
