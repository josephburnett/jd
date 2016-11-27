package jd

import (
	"testing"
)

func TestArrayJson(t *testing.T) {
	checkJson(t, `[]`, `[]`)
	checkJson(t, ` [ ] `, `[]`)
	checkJson(t, `[1,2,3]`, `[1,2,3]`)
	checkJson(t, ` [1, 2, 3] `, `[1,2,3]`)
}

func TestArrayEqual(t *testing.T) {
	checkEqual(t, `[]`, `[]`)
	checkEqual(t, `[1,2,3]`, `[1,2,3]`)
	checkEqual(t, `[[]]`, `[[]]`)
	checkEqual(t, `[{"a":1}]`, `[{"a":1}]`)
	checkEqual(t, `[{"a":[]}]`, `[{"a":[]}]`)
}

func TestArrayNotEqual(t *testing.T) {
	checkNotEqual(t, `[]`, `0`)
	checkNotEqual(t, `[]`, `{}`)
	checkNotEqual(t, `[]`, `[[]]`)
	checkNotEqual(t, `[1,2,3]`, `[3,2,1]`)
}

func TestArrayHash(t *testing.T) {
	checkHash(t, `[]`, `[]`, true)
	checkHash(t, `[1]`, `[]`, false)
	checkHash(t, `[1]`, `[1]`, true)
	checkHash(t, `[1]`, `[2]`, false)
	checkHash(t, `[[1]]`, `[[1]]`, true)
	checkHash(t, `[[1]]`, `[[[1]]]`, false)
}

func TestArrayDiff(t *testing.T) {
	checkDiff(t, `[]`, `[]`)
	checkDiff(t, `[1]`, `[]`,
		`@ [0]`,
		`- 1`)
	checkDiff(t, `[[]]`, `[[1]]`,
		`@ [0, 0]`,
		`+ 1`)
	checkDiff(t, `[1]`, `[2]`,
		`@ [0]`,
		`- 1`,
		`+ 2`)
	checkDiff(t, `[]`, `[2]`,
		`@ [0]`,
		`+ 2`)
	checkDiff(t, `[[]]`, `[{}]`,
		`@ [0]`,
		`- []`,
		`+ {}`)
	checkDiff(t, `[{"a":[1]}]`, `[{"a":[2]}]`,
		`@ [0,"a",0]`,
		`- 1`,
		`+ 2`)
	checkDiff(t, `[1,2,3]`, `[1,2]`,
		`@ [2]`,
		`- 3`)
	checkDiff(t, `[1,2,3]`, `[1,4,3]`,
		`@ [1]`,
		`- 2`,
		`+ 4`)
	checkDiff(t, `[1,2,3]`, `[1,null,3]`,
		`@ [1]`,
		`- 2`,
		`+ null`)
}

func TestArrayPatch(t *testing.T) {
	checkPatch(t, `[]`, `[]`)
	checkPatch(t, `[1]`, `[]`,
		`@ [0]`,
		`- 1`)
	checkPatch(t, `[[]]`, `[[1]]`,
		`@ [0, 0]`,
		`+ 1`)
	checkPatch(t, `[1]`, `[2]`,
		`@ [0]`,
		`- 1`,
		`+ 2`)
	checkPatch(t, `[]`, `[2]`,
		`@ [0]`,
		`+ 2`)
	checkPatch(t, `[[]]`, `[{}]`,
		`@ [0]`,
		`- []`,
		`+ {}`)
	checkPatch(t, `[{"a":[1]}]`, `[{"a":[2]}]`,
		`@ [0,"a",0]`,
		`- 1`,
		`+ 2`)
	checkPatch(t, `[1,2,3]`, `[1,2]`,
		`@ [2]`,
		`- 3`)
	checkPatch(t, `[1,2,3]`, `[1,4,3]`,
		`@ [1]`,
		`- 2`,
		`+ 4`)
	checkPatch(t, `[1,2,3]`, `[1,null,3]`,
		`@ [1]`,
		`- 2`,
		`+ null`)
}

func TestArrayPatchError(t *testing.T) {
	checkPatchError(t, `[]`,
		`@ ["a"]`,
		`+ 1`)
	checkPatchError(t, `[]`,
		`@ [0]`,
		`- 1`)
	checkPatchError(t, `[]`,
		`@ [0]`,
		`- null`)
	checkPatchError(t, `[1,2,3]`,
		`@ [1]`,
		`- 2`)
	checkPatchError(t, `[1,2,3]`,
		`@ [0]`,
		`- 1`)
	checkPatchError(t, `[1,3]`,
		`@ [1]`,
		`+ 2`)
}
