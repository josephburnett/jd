package jd

import (
	"testing"
)

func TestNullJson(t *testing.T) {
	checkJson(t, `null`, `null`)
}

func TestNullEqual(t *testing.T) {
	checkEqual(t, `null`, `null`)
}

func TestNullNotEqual(t *testing.T) {
	checkNotEqual(t, `null`, `0`)
	checkNotEqual(t, `null`, `[]`)
	checkNotEqual(t, `null`, `{}`)
}

func TestNullHash(t *testing.T) {
	checkHash(t, `null`, `null`, true)
	checkHash(t, `null`, ``, false)
}

func TestNullDiff(t *testing.T) {
	checkDiff(t, `null`, `null`)
	checkDiff(t, `null`, ``,
		`@ []`,
		`- null`)
	checkDiff(t, ``, `null`,
		`@ []`,
		`+ null`)
}

func TestNullPatch(t *testing.T) {
	checkPatch(t, `null`, `null`)
	checkPatch(t, `null`, ``,
		`@ []`,
		`- null`)
	checkPatch(t, ``, `null`,
		`@ []`,
		`+ null`)
}

func TestNullPatchError(t *testing.T) {
	checkPatchError(t, `null`,
		`@ []`,
		`- 0`)
}
