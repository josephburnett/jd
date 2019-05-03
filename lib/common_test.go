package jd

import (
	"reflect"
	"testing"
)

func checkJson(ctx *testContext, a, b string) {
	nodeA, err := ReadJsonString(a, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeAJson := nodeA.Json()
	if nodeAJson != b {
		t.Errorf("%v.Json() = %v. Want %v.", nodeA, nodeAJson, b)
	}
}

func checkEqual(ctx *testContext, a, b string) {
	nodeA, err := unmarshal([]byte(a), ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b), ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	// TODO: plumb metadata into Equals.
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

func checkNotEqual(ctx *testContext, a, b string, metadata ...Metadata) {
	nodeA, err := unmarshal([]byte(a), ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b), ctx.readMetadata...)
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

func checkHash(ctx *testContext, a, b string, wantSame bool) {
	nodeA, err := unmarshal([]byte(a), ctx.readMetadata...)
	if err != nil {
		t.Fatalf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b), ctx.readMetadata...)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// TODO: plumb metadata into hashCode and get rid of ident method.
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

func checkDiff(ctx *testContext, a, b string, diffLines ...string) {
	nodeA, err := ReadJsonString(a, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := ReadJsonString(b, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	d := nodeA.Diff(nodeB, ctx.applyMetadata...)
	expectedDiff, err := ReadDiffString(diff, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(d, expectedDiff) {
		t.Errorf("%v.Diff(%v) = %v. Want %v.", nodeA, nodeB, d, expectedDiff)
	}
}

func checkPatch(ctx *testContext, a, e string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	diff, err := ReadDiffString(diffString, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	expect, err := ReadJsonString(e, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	b, err := initial.Patch(diff, ctx.applyMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !expect.Equals(b) {
		t.Errorf("%v.Patch(%v) = %v. Want %v.",
			a, diffLines, renderJson(b), e)
	}
}

func checkPatchError(ctx *testContext, a string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	diff, err := ReadDiffString(diffString, ctx.readMetadata...)
	if err != nil {
		t.Errorf(err.Error())
	}
	b, err := initial.Patch(diff, ctx.applyMetadata...)
	if b != nil {
		t.Errorf("%v.Patch(%v) = %v. Want nil.", initial, diff, b)
	}
	if err == nil {
		t.Errorf("Expected error. Got nil.")
	}
}

type testContext struct {
	t             *testing.T
	readMetadata  []Metadata
	applyMetadata []Metadata
}

func newTestContext(t *testing.T) *testContext {
	return &testContext{
		t:             t,
		readMetadata:  make([]Metadata, 0),
		applyMetadata: make([]Metadata, 0),
	}
}

func (tc *testContext) withReadMetadata(metadata ...Metadata) {
	tc.readMetadata = append(tc.readMetadata, metadata...)
}

func (tc *testContext) withApplyMetadata(metadata ...Metadata) {
	tc.applyMetadata = append(tc.applyMetadata, metadata...)
}
