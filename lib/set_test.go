package jd

import (
	"testing"
)

func TestSetJson(t *testing.T) {
	checkJson(t, `[]`, `[]`, SET)
	checkJson(t, ` [ ] `, `[]`, SET)
	checkJson(t, `[1,2,3]`, `[3,2,1]`, SET)
	checkJson(t, `[3,2,1]`, `[3,2,1]`, SET)
	checkJson(t, ` [1, 2, 3] `, `[3,2,1]`, SET)
	checkJson(t, `[1,1,1]`, `[1]`, SET)
}

func TestSetEquals(t *testing.T) {
	checkEqual(t, `[]`, `[]`, SET)
	checkEqual(t, `[1,2,3]`, `[3,2,1]`, SET)
	checkEqual(t, `[1,2,3]`, `[2,3,1]`, SET)
	checkEqual(t, `[1,2,3]`, `[1,3,2]`, SET)
	checkEqual(t, `[{},{}]`, `[{},{}]`, SET)
	checkEqual(t, `[[1,2],[3,4]]`, `[[2,1],[4,3]]`, SET)
	checkEqual(t, `[1,1,1]`, `[1]`, SET)
	checkEqual(t, `[1,2,1]`, `[2,1,2]`, SET)
}

func TestSetNotEquals(t *testing.T) {
	checkNotEqual(t, `[]`, `[1]`, MULTISET)
	checkNotEqual(t, `[1,2,3]`, `[1,2,2]`, MULTISET)
	checkNotEqual(t, `[1,2,3]`, `[1,2]`, MULTISET)
	checkNotEqual(t, `[[],[1]]`, `[[],[2]]`, MULTISET)
}

func TestSetDiff(t *testing.T) {
	checkDiffMetadata(t, SET, `[]`, `[]`)
	checkDiffMetadata(t, SET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkDiffMetadata(t, SET, `[1,2]`, `[1,2]`)
	checkDiffMetadata(t, SET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`)
	checkDiffMetadata(t, SET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkDiffMetadata(t, SET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffMetadata(t, SET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffMetadata(t, SET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`+ "baz"`)
	checkDiffMetadata(t, SET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "baz"`)
	checkDiffMetadata(t, SET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestSetPatch(t *testing.T) {
	checkPatchMetadata(t, SET, `[]`, `[]`)
	checkPatchMetadata(t, SET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkPatchMetadata(t, SET, `[1,2]`, `[1,2]`)
	checkPatchMetadata(t, SET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`)
	checkPatchMetadata(t, SET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkPatchMetadata(t, SET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchMetadata(t, SET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchMetadata(t, SET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`+ "baz"`)
	checkPatchMetadata(t, SET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "baz"`)
	checkPatchMetadata(t, SET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestSetPatchError(t *testing.T) {
	checkPatchErrorMetadata(t, SET, `[]`,
		`@ [{}]`,
		`- 1`)
	checkPatchErrorMetadata(t, SET, `[1]`,
		`@ [{}]`,
		`- 1`,
		`- 1`)
	checkPatchErrorMetadata(t, SET, `[]`,
		`@ [{}]`,
		`- 1`,
		`+ 1`)
	checkPatchErrorMetadata(t, SET, `[]`,
		`@ []`,
		`- {}`)
	checkPatchErrorMetadata(t, SET, `[]`,
		`@ [{}]`,
		`+ 1`,
		`+ 1`,
		`@ [{}]`,
		`- 1`,
		`- 1`)
}
