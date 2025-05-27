package jd

import (
	"regexp"
	"strings"
	"testing"
)

func CheckJson(ctx *testContext, a, b string) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	nodeAJson := nodeA.Json(ctx.options...)
	if nodeAJson != b {
		ctx.t.Errorf("%v.Json() = %v. Want %v.", nodeA, nodeAJson, b)
	}
}

func CheckEqual(ctx *testContext, a, b string) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	if !nodeA.Equals(nodeB, ctx.options...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeB.Equals(nodeA, ctx.options...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeA.Equals(nodeA, ctx.options...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
	if !nodeB.Equals(nodeB, ctx.options...) {
		ctx.t.Errorf("%v.Equals(%v) == false. Want true.", nodeA, nodeB)
	}
}

func CheckNotEqual(ctx *testContext, a, b string) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	if nodeA.Equals(nodeB, ctx.options...) {
		ctx.t.Errorf("nodeA.Equals(nodeB) == true. Want false.")
	}
	if nodeB.Equals(nodeA, ctx.options...) {
		ctx.t.Errorf("nodeB.Equals(nodeA) == true. Want false.")
	}
}

func CheckHash(ctx *testContext, a, b string, wantSame bool) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Fatalf("%v", err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Fatalf("%v", err.Error())
	}
	o := refine(&options{retain: ctx.options}, nil)
	hashA := nodeA.hashCode(o)
	hashB := nodeB.hashCode(o)
	if wantSame && hashA != hashB {
		ctx.t.Errorf("%v.hashCode = %v. %v.hashCode = %v. Want the same.",
			a, hashA, b, hashB)
	}
	if !wantSame && hashA == hashB {
		ctx.t.Errorf("%v.hashCode = %v. %v.hashCode = %v. Want the different.",
			a, hashA, b, hashB)
	}
}

func CheckDiff(ctx *testContext, a, b string, diffLines ...string) {
	nodeA, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Fatalf("%v", err.Error())
	}
	nodeB, err := ReadJsonString(b)
	if err != nil {
		ctx.t.Fatalf("%v", err.Error())
	}
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	d := nodeA.Diff(nodeB, ctx.options...)
	got := d.Render()
	if got != diff {
		ctx.t.Errorf("%v.Diff(%v) = \n%v. Want %v.", nodeA, nodeB, got, diff)
	}
}

func CheckPatch(ctx *testContext, a, e string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	diff, err := ReadDiffString(diffString)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	expect, err := ReadJsonString(e)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	b, err := initial.Patch(diff)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	if !expect.Equals(b, ctx.options...) {
		ctx.t.Errorf("%v.Patch(%v) = %v. Want %v.",
			a, diffLines, b, e)
	}
}

func CheckPatchError(ctx *testContext, a string, diffLines ...string) {
	diffString := ""
	for _, dl := range diffLines {
		diffString += dl + "\n"
	}
	initial, err := ReadJsonString(a)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
	}
	diff, err := ReadDiffString(diffString)
	if err != nil {
		ctx.t.Errorf("%v", err.Error())
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

func m(m ...Option) []Option {
	return m
}

type testContext struct {
	t       *testing.T
	options []Option
}

func newTestContext(t *testing.T) *testContext {
	return &testContext{
		t:       t,
		options: make([]Option, 0),
	}
}

func (tc *testContext) withOptions(options ...Option) *testContext {
	tc.options = append(tc.options, options...)
	return tc
}

// stripAnsiCodes removes ANSI color escape sequences from a string
func stripAnsiCodes(input string) string {
	// Regular expression to match ANSI escape codes
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return re.ReplaceAllString(input, "")
}
