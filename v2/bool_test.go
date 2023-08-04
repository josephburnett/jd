package jd

import (
	"testing"
)

func TestBoolJson(t *testing.T) {
	ctx := newTestContext(t)
	checkJson(ctx, `true`, `true`)
	checkJson(ctx, `false`, `false`)
}

func TestBoolEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkEqual(ctx, `true`, `true`)
	checkEqual(ctx, `false`, `false`)
}

func TestBoolNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkNotEqual(ctx, `true`, `false`)
	checkNotEqual(ctx, `false`, `true`)
	checkNotEqual(ctx, `false`, `[]`)
	checkNotEqual(ctx, `true`, `"true"`)
}

func TestBoolHash(t *testing.T) {
	ctx := newTestContext(t)
	checkHash(ctx, `true`, `true`, true)
	checkHash(ctx, `false`, `false`, true)
	checkHash(ctx, `true`, `false`, false)
}

func TestBoolDiff(t *testing.T) {
	ctx := newTestContext(t)
	checkDiff(ctx, `true`, `true`)
	checkDiff(ctx, `false`, `false`)
	checkDiff(ctx, `true`, `false`,
		`^ {"Version":2}`,
		`@ []`,
		`- true`,
		`+ false`)
	checkDiff(ctx, `false`, `true`,
		`^ {"Version":2}`,
		`@ []`,
		`- false`,
		`+ true`)
	ctx = ctx.withOptions(MERGE)
	checkDiff(ctx, `true`, `false`,
		`^ {"Version":2}`,
		`^ {"Merge":true}`,
		`@ []`,
		`+ false`)
}

func TestBoolPatch(t *testing.T) {
	ctx := newTestContext(t)
	checkPatch(ctx, `true`, `true`)
	checkPatch(ctx, `false`, `false`)
	checkPatch(ctx, `true`, `false`,
		`@ []`,
		`- true`,
		`+ false`)
	checkPatch(ctx, `false`, `true`,
		`@ []`,
		`- false`,
		`+ true`)
	checkPatch(ctx, `false`, `true`,
		`^ {"Merge":true}`,
		`@ []`,
		`+ true`)
	checkPatch(ctx, `true`, `false`,
		`^ {"Merge":true}`,
		`@ []`,
		`+ false`)
	checkPatch(ctx, `true`, ``,
		`^ {"Merge":true}`,
		`@ []`,
		`+`)
}

func TestBoolPatchError(t *testing.T) {
	ctx := newTestContext(t)
	checkPatchError(ctx, `true`,
		`@ []`,
		`- false`)
	checkPatchError(ctx, `false`,
		`@ []`,
		`- true`)
	checkPatchError(ctx, `true`,
		`@ [["MERGE"]]`,
		`- true`,
		`+ false`)
	checkPatchError(ctx, `false`,
		`@ [["MERGE"]]`,
		`- false`,
		`+ true`)
}
