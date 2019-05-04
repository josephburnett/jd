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
		`@ []`,
		`- true`,
		`+ false`)
	checkDiff(ctx, `false`, `true`,
		`@ []`,
		`- false`,
		`+ true`)
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
}

func TestBoolPatchError(t *testing.T) {
	ctx := newTestContext(t)
	checkPatchError(ctx, `true`,
		`@ []`,
		`- false`)
	checkPatchError(ctx, `false`,
		`@ []`,
		`- true`)
}
