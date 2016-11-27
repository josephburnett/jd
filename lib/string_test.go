package jd

import (
	"testing"
)

func TestStringJson(t *testing.T) {
	checkJson(t, `""`, `""`)
	checkJson(t, ` "" `, `""`)
	checkJson(t, `"\""`, `"\""`)
}

func TestStringEqual(t *testing.T) {
	checkEqual(t, `""`, `""`)
	checkEqual(t, `"a"`, `"a"`)
	checkEqual(t, `"123"`, `"123"`)
}

func TestStringNotEqual(t *testing.T) {
	checkNotEqual(t, `""`, `"a"`)
	checkNotEqual(t, `""`, `[]`)
	checkNotEqual(t, `""`, `{}`)
	checkNotEqual(t, `""`, `0`)
}

func TestStringHash(t *testing.T) {
	checkHash(t, `""`, `""`, true)
	checkHash(t, `"abc"`, `"abc"`, true)
	checkHash(t, `""`, `" "`, false)
	checkHash(t, `"abc"`, `"123"`, false)
}

func TestStringDiff(t *testing.T) {
	checkDiff(t, `""`, `""`)
	checkDiff(t, `""`, `1`,
		`@ []`,
		`- ""`,
		`+ 1`)
	checkDiff(t, `null`, `"abc"`,
		`@ []`,
		`- null`,
		`+ "abc"`)
}

func TestStringPatch(t *testing.T) {
	checkPatch(t, `""`, `""`)
	checkPatch(t, `""`, `1`,
		`@ []`,
		`- ""`,
		`+ 1`)
	checkPatch(t, `null`, `"abc"`,
		`@ []`,
		`- null`,
		`+ "abc"`)
}

func TestStringPatchError(t *testing.T) {
	checkPatchError(t, `""`,
		`@ []`,
		`- "a"`,
		`+ ""`)
	checkPatchError(t, `null`,
		`@ []`,
		`+ "a"`)
}
