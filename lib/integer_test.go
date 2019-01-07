package jd

import (
	"testing"
)

func TestIntegerJson(t *testing.T) {
	checkJson(t, `0`, `0`)
}

func TestIntegerEqual(t *testing.T) {
	checkEqual(t, `0`, `0`)
	checkEqual(t, `123`, `123`)
}

func TestIntegerNotEqual(t *testing.T) {
	checkNotEqual(t, `0`, `1`)
	checkNotEqual(t, `1234`, `1235`)
}

func TestIntegerHash(t *testing.T) {
	checkHash(t, `0`, `0`, true)
	checkHash(t, `0`, `1`, false)
}

func TestIntegerDiff(t *testing.T) {
	checkDiff(t, `0`, `0`)
	checkDiff(t, `0`, `1`,
		`@ []`,
		`- 0`,
		`+ 1`)
	checkDiff(t, `0`, ``,
		`@ []`,
		`- 0`)
}

func TestIntegerPatch(t *testing.T) {
	checkPatch(t, `0`, `0`)
	checkPatch(t, `0`, `1`,
		`@ []`,
		`- 0`,
		`+ 1`)
	checkPatch(t, `0`, ``,
		`@ []`,
		`- 0`)
}

func TestIntegerPatchError(t *testing.T) {
	checkPatchError(t, `0`,
		`@ []`,
		`- 1`)
	checkPatchError(t, ``,
		`@ []`,
		`- 0`)
}
