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
	checkDiffOption(t, SET, `[]`, `[]`)
	checkDiffOption(t, SET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkDiffOption(t, SET, `[1,2]`, `[1,2]`)
	checkDiffOption(t, SET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`)
	checkDiffOption(t, SET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkDiffOption(t, SET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffOption(t, SET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffOption(t, SET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`+ "baz"`)
	checkDiffOption(t, SET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "baz"`)
	checkDiffOption(t, SET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestSetPatch(t *testing.T) {
	checkPatchOption(t, SET, `[]`, `[]`)
	checkPatchOption(t, SET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkPatchOption(t, SET, `[1,2]`, `[1,2]`)
	checkPatchOption(t, SET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`)
	checkPatchOption(t, SET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkPatchOption(t, SET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchOption(t, SET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchOption(t, SET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`+ "baz"`)
	checkPatchOption(t, SET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "baz"`)
	checkPatchOption(t, SET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestSetPatchError(t *testing.T) {
	checkPatchErrorOption(t, SET, `[]`,
		`@ [{}]`,
		`- 1`)
	checkPatchErrorOption(t, SET, `[1]`,
		`@ [{}]`,
		`- 1`,
		`- 1`)
	checkPatchErrorOption(t, SET, `[]`,
		`@ [{}]`,
		`- 1`,
		`+ 1`)
	checkPatchErrorOption(t, SET, `[]`,
		`@ []`,
		`- {}`)
	checkPatchErrorOption(t, SET, `[]`,
		`@ [{}]`,
		`+ 1`,
		`+ 1`,
		`@ [{}]`,
		`- 1`,
		`- 1`)
}
