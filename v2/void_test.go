package jd

import (
	"testing"
)

func TestVoidJson(t *testing.T) {
	ctx := newTestContext(t)
	checkJson(ctx, ``, ``)
}

func TestVoidEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkEqual(ctx, ``, ``)
	checkEqual(ctx, `   `, ``)
}

func TestVoidNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkNotEqual(ctx, ``, `null`)
	checkNotEqual(ctx, ``, `0`)
	checkNotEqual(ctx, ``, `[]`)
	checkNotEqual(ctx, ``, `{}`)
}

func TestVoidHash(t *testing.T) {
	ctx := newTestContext(t)
	checkHash(ctx, ``, ``, true)
	checkHash(ctx, ``, `null`, false)
}

func TestVoidDiff(t *testing.T) {
	ctx := newTestContext(t)
	checkDiff(ctx, ``, ``)
	checkDiff(ctx, ``, `1`,
		`^ {"Version":2}`,
		`@ []`,
		`+ 1`)
	checkDiff(ctx, `1`, ``,
		`^ {"Version":2}`,
		`@ []`,
		`- 1`)
}

func TestVoidPatch(t *testing.T) {
	ctx := newTestContext(t)
	checkPatch(ctx, ``, ``)
	checkPatch(ctx, ``, `1`,
		`@ []`,
		`+ 1`)
	checkPatch(ctx, `1`, ``,
		`@ []`,
		`- 1`)
	checkPatch(ctx, ``, `1`,
		`^ {"Merge":true}`,
		`@ []`,
		`+ 1`)
	checkPatch(ctx, ``, ``,
		`^ {"Merge":true}`,
		`@ []`,
		`+`)
}

func TestVoidPatchError(t *testing.T) {
	ctx := newTestContext(t)
	checkPatchError(ctx, ``,
		`@ []`,
		`- null`)
}
