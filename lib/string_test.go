package jd

import (
	"testing"
)

func TestStringJson(t *testing.T) {
	ctx := newTestContext(t)
	checkJson(ctx, `""`, `""`)
	checkJson(ctx, ` "" `, `""`)
	checkJson(ctx, `"\""`, `"\""`)
}

func TestStringEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkEqual(ctx, `""`, `""`)
	checkEqual(ctx, `"a"`, `"a"`)
	checkEqual(ctx, `"123"`, `"123"`)
}

func TestStringNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkNotEqual(ctx, `""`, `"a"`)
	checkNotEqual(ctx, `""`, `[]`)
	checkNotEqual(ctx, `""`, `{}`)
	checkNotEqual(ctx, `""`, `0`)
}

func TestStringHash(t *testing.T) {
	ctx := newTestContext(t)
	checkHash(ctx, `""`, `""`, true)
	checkHash(ctx, `"abc"`, `"abc"`, true)
	checkHash(ctx, `""`, `" "`, false)
	checkHash(ctx, `"abc"`, `"123"`, false)
}

func TestStringDiff(t *testing.T) {
	ctx := newTestContext(t)
	checkDiff(ctx, `""`, `""`)
	checkDiff(ctx, `""`, `1`,
		`@ []`,
		`- ""`,
		`+ 1`)
	checkDiff(ctx, `null`, `"abc"`,
		`@ []`,
		`- null`,
		`+ "abc"`)
}

func TestStringPatch(t *testing.T) {
	ctx := newTestContext(t)
	checkPatch(ctx, `""`, `""`)
	checkPatch(ctx, `""`, `1`,
		`@ []`,
		`- ""`,
		`+ 1`)
	checkPatch(ctx, `null`, `"abc"`,
		`@ []`,
		`- null`,
		`+ "abc"`)
	checkPatch(ctx, `"def"`, `"abc"`,
		`@ [["merge"]]`,
		`+ "abc"`)

	// Null deletes a node
	checkPatch(ctx, `"abc"`, ``,
		`@ [["merge"]]`,
		`+ null`)
}

func TestStringPatchError(t *testing.T) {
	ctx := newTestContext(t)
	checkPatchError(ctx, `""`,
		`@ []`,
		`- "a"`,
		`+ ""`)
	checkPatchError(ctx, `null`,
		`@ []`,
		`+ "a"`)
}
