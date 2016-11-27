package jd

import (
	"testing"
)

func TestVoidJson(t *testing.T) {
	checkJson(t, ``, ``)
}

func TestVoidEqual(t *testing.T) {
	checkEqual(t, ``, ``)
	checkEqual(t, `   `, ``)
}

func TestVoidNotEqual(t *testing.T) {
	checkNotEqual(t, ``, `null`)
	checkNotEqual(t, ``, `0`)
	checkNotEqual(t, ``, `[]`)
	checkNotEqual(t, ``, `{}`)
}

func TestVoidHash(t *testing.T) {
	checkHash(t, ``, ``, true)
	checkHash(t, ``, `null`, false)
}

func TestVoidDiff(t *testing.T) {
	checkDiff(t, ``, ``)
	checkDiff(t, ``, `1`,
		`@ []`,
		`+ 1`)
	checkDiff(t, `1`, ``,
		`@ []`,
		`- 1`)
}

func TestVoidPatch(t *testing.T) {
	checkPatch(t, ``, ``)
	checkPatch(t, ``, `1`,
		`@ []`,
		`+ 1`)
	checkPatch(t, `1`, ``,
		`@ []`,
		`- 1`)
}

func TestVoidPatchError(t *testing.T) {
	checkPatchError(t, ``,
		`@ []`,
		`- null`)
}
