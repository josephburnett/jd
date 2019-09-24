package jd

import (
	"testing"
)

// TODO: convert array tests to table tests.
func TestArrayJson(t *testing.T) {
	ctx := newTestContext(t)
	checkJson(ctx, `[]`, `[]`)
	checkJson(ctx, ` [ ] `, `[]`)
	checkJson(ctx, `[1,2,3]`, `[1,2,3]`)
	checkJson(ctx, ` [1, 2, 3] `, `[1,2,3]`)
}

func TestArrayEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkEqual(ctx, `[]`, `[]`)
	checkEqual(ctx, `[1,2,3]`, `[1,2,3]`)
	checkEqual(ctx, `[[]]`, `[[]]`)
	checkEqual(ctx, `[{"a":1}]`, `[{"a":1}]`)
	checkEqual(ctx, `[{"a":[]}]`, `[{"a":[]}]`)
}

func TestArrayNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkNotEqual(ctx, `[]`, `0`)
	checkNotEqual(ctx, `[]`, `{}`)
	checkNotEqual(ctx, `[]`, `[[]]`)
	checkNotEqual(ctx, `[1,2,3]`, `[3,2,1]`)
}

func TestArrayHash(t *testing.T) {
	ctx := newTestContext(t)
	checkHash(ctx, `[]`, `[]`, true)
	checkHash(ctx, `[1]`, `[]`, false)
	checkHash(ctx, `[1]`, `[1]`, true)
	checkHash(ctx, `[1]`, `[2]`, false)
	checkHash(ctx, `[[1]]`, `[[1]]`, true)
	checkHash(ctx, `[[1]]`, `[[[1]]]`, false)
}

func TestArrayDiff(t *testing.T) {
	ctx := newTestContext(t)
	checkDiff(ctx, `[]`, `[]`)
	checkDiff(ctx, `[1]`, `[]`,
		`@ [0]`,
		`- 1`)
	checkDiff(ctx, `[[]]`, `[[1]]`,
		`@ [0, 0]`,
		`+ 1`)
	checkDiff(ctx, `[1]`, `[2]`,
		`@ [0]`,
		`- 1`,
		`+ 2`)
	checkDiff(ctx, `[]`, `[2]`,
		`@ [0]`,
		`+ 2`)
	checkDiff(ctx, `[[]]`, `[{}]`,
		`@ [0]`,
		`- []`,
		`+ {}`)
	checkDiff(ctx, `[{"a":[1]}]`, `[{"a":[2]}]`,
		`@ [0,"a",0]`,
		`- 1`,
		`+ 2`)
	checkDiff(ctx, `[1,2,3]`, `[1,2]`,
		`@ [2]`,
		`- 3`)
	checkDiff(ctx, `[1,2,3]`, `[1,4,3]`,
		`@ [1]`,
		`- 2`,
		`+ 4`)
	checkDiff(ctx, `[1,2,3]`, `[1,null,3]`,
		`@ [1]`,
		`- 2`,
		`+ null`)
	checkDiff(ctx, `[]`, `[3,4,5]`,
		`@ [0]`,
		`+ 3`,
		`@ [1]`,
		`+ 4`,
		`@ [2]`,
		`+ 5`)
}

func TestArrayPatch(t *testing.T) {
	ctx := newTestContext(t)
	checkPatch(ctx, `[]`, `[]`)
	checkPatch(ctx, `[1]`, `[]`,
		`@ [0]`,
		`- 1`)
	checkPatch(ctx, `[[]]`, `[[1]]`,
		`@ [0, 0]`,
		`+ 1`)
	checkPatch(ctx, `[1]`, `[2]`,
		`@ [0]`,
		`- 1`,
		`+ 2`)
	checkPatch(ctx, `[]`, `[2]`,
		`@ [0]`,
		`+ 2`)
	checkPatch(ctx, `[[]]`, `[{}]`,
		`@ [0]`,
		`- []`,
		`+ {}`)
	checkPatch(ctx, `[{"a":[1]}]`, `[{"a":[2]}]`,
		`@ [0,"a",0]`,
		`- 1`,
		`+ 2`)
	checkPatch(ctx, `[1,2,3]`, `[1,2]`,
		`@ [2]`,
		`- 3`)
	checkPatch(ctx, `[1,2,3]`, `[1,4,3]`,
		`@ [1]`,
		`- 2`,
		`+ 4`)
	checkPatch(ctx, `[1,2,3]`, `[1,null,3]`,
		`@ [1]`,
		`- 2`,
		`+ null`)
	checkPatch(ctx, "[]", "[3,4,5]",
		"@ [0]",
		"+ 3",
		"@ [1]",
		"+ 4",
		"@ [2]",
		"+ 5")
}

func TestArrayPatchError(t *testing.T) {
	ctx := newTestContext(t)
	checkPatchError(ctx, `[]`,
		`@ ["a"]`,
		`+ 1`)
	checkPatchError(ctx, `[]`,
		`@ [0]`,
		`- 1`)
	checkPatchError(ctx, `[]`,
		`@ [0]`,
		`- null`)
	checkPatchError(ctx, `[1,2,3]`,
		`@ [1]`,
		`- 2`)
	checkPatchError(ctx, `[1,2,3]`,
		`@ [0]`,
		`- 1`)
	checkPatchError(ctx, `[1,3]`,
		`@ [1]`,
		`+ 2`)
}
