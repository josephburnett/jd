package jd

import (
	"testing"
)

func TestBoolJson(t *testing.T) {
	checkJson(t, `true`, `true`)
	checkJson(t, `false`, `false`)
}

func TestBoolEqual(t *testing.T) {
	checkEqual(t, `true`, `true`)
	checkEqual(t, `false`, `false`)
}

func TestBoolNotEqual(t *testing.T) {
	checkNotEqual(t, `true`, `false`)
	checkNotEqual(t, `false`, `true`)
	checkNotEqual(t, `false`, `[]`)
	checkNotEqual(t, `true`, `"true"`)
}

func TestBoolHash(t *testing.T) {
	checkHash(t, `true`, `true`, true)
	checkHash(t, `false`, `false`, true)
	checkHash(t, `true`, `false`, false)
}

func TestBoolDiff(t *testing.T) {
	checkDiff(t, `true`, `true`)
	checkDiff(t, `false`, `false`)
	checkDiff(t, `true`, `false`,
		`@ []`,
		`- true`,
		`+ false`)
	checkDiff(t, `false`, `true`,
		`@ []`,
		`- false`,
		`+ true`)
}

func TestBoolPatch(t *testing.T) {
	checkPatch(t, `true`, `true`)
	checkPatch(t, `false`, `false`)
	checkPatch(t, `true`, `false`,
		`@ []`,
		`- true`,
		`+ false`)
	checkPatch(t, `false`, `true`,
		`@ []`,
		`- false`,
		`+ true`)
}

func TestBoolPatchError(t *testing.T) {
	checkPatchError(t, `true`,
		`@ []`,
		`- false`)
	checkPatchError(t, `false`,
		`@ []`,
		`- true`)
}
