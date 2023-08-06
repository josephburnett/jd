package jd

import (
	"testing"
)

func TestNumberJson(t *testing.T) {
	ctx := newTestContext(t)
	checkJson(ctx, `0`, `0`)
	checkJson(ctx, `0.0`, `0`)
	checkJson(ctx, `0.01`, `0.01`)
}

func TestNumberEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkEqual(ctx, `0`, `0`)
	checkEqual(ctx, `0`, `0.0`)
	checkEqual(ctx, `0.0001`, `0.0001`)
	checkEqual(ctx, `123`, `123`)
}

func TestNumberNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkNotEqual(ctx, `0`, `1`)
	checkNotEqual(ctx, `0`, `0.0001`)
	checkNotEqual(ctx, `1234`, `1235`)
}

func TestNumberHash(t *testing.T) {
	ctx := newTestContext(t)
	checkHash(ctx, `0`, `0`, true)
	checkHash(ctx, `0`, `1`, false)
	checkHash(ctx, `1.0`, `1`, true)
	checkHash(ctx, `0.1`, `0.01`, false)
}

func TestNumberDiff(t *testing.T) {
	ctx := newTestContext(t)
	checkDiff(ctx, `0`, `0`)
	checkDiff(ctx, `0`, `1`,
		`^ {"Version":2}`,
		`@ []`,
		`- 0`,
		`+ 1`)
	checkDiff(ctx, `0`, ``,
		`^ {"Version":2}`,
		`@ []`,
		`- 0`)
	ctx = ctx.withOptions(MERGE)
	checkDiff(ctx, `1`, `2`,
		`^ {"Version":2}`,
		`^ {"Merge":true}`,
		`@ []`,
		`+ 2`)
}

func TestNumberPatch(t *testing.T) {
	ctx := newTestContext(t)
	checkPatch(ctx, `0`, `0`)
	checkPatch(ctx, `0`, `1`,
		`@ []`,
		`- 0`,
		`+ 1`)
	checkPatch(ctx, `0`, ``,
		`@ []`,
		`- 0`)
	checkPatch(ctx, `0`, `1`,
		`^ {"Merge":true}`,
		`@ []`,
		`+ 1`)
	checkPatch(ctx, `1`, ``,
		`^ {"Merge":true}`,
		`@ []`,
		`+`)
}

func TestNumberPatchError(t *testing.T) {
	ctx := newTestContext(t)
	checkPatchError(ctx, `0`,
		`@ []`,
		`- 1`)
	checkPatchError(ctx, ``,
		`@ []`,
		`- 0`)
	checkPatchError(ctx, `0`,
		`^ {"Merge":true}`,
		`@ []`,
		`- 0`,
		`+ 1`)
}
