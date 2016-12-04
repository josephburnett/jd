package jd

import (
	"testing"
)

func TestMultisetJson(t *testing.T) {
	checkJson(t, `[]`, `[]`, MULTISET)
	checkJson(t, ` [ ] `, `[]`, MULTISET)
	checkJson(t, `[1,2,3]`, `[1,2,3]`, MULTISET)
	checkJson(t, ` [1, 2, 3] `, `[1,2,3]`, MULTISET)
	checkJson(t, `[1,1,1]`, `[1,1,1]`, MULTISET)
}

func TestMultisetEquals(t *testing.T) {
	checkEqual(t, `[]`, `[]`, MULTISET)
	checkEqual(t, `[1,2,3]`, `[3,2,1]`, MULTISET)
	checkEqual(t, `[1,2,3]`, `[2,3,1]`, MULTISET)
	checkEqual(t, `[1,2,3]`, `[1,3,2]`, MULTISET)
	checkEqual(t, `[{},{}]`, `[{},{}]`, MULTISET)
	checkEqual(t, `[[1,2],[3,4]]`, `[[2,1],[4,3]]`, MULTISET)
}

func TestMultisetNotEquals(t *testing.T) {
	checkNotEqual(t, `[]`, `[1]`, MULTISET)
	checkNotEqual(t, `[1,2,3]`, `[1,2,2]`, MULTISET)
	checkNotEqual(t, `[1,2,3]`, `[1,2]`, MULTISET)
	checkNotEqual(t, `[[],[1]]`, `[[],[2]]`, MULTISET)
}

func TestMultisetDiff(t *testing.T) {
	checkDiffOption(t, MULTISET, `[]`, `[]`)
	checkDiffOption(t, MULTISET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkDiffOption(t, MULTISET, `[1,2]`, `[1,2]`)
	checkDiffOption(t, MULTISET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`,
		`+ 2`)
	checkDiffOption(t, MULTISET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkDiffOption(t, MULTISET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffOption(t, MULTISET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffOption(t, MULTISET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`- "foo"`,
		`+ "baz"`)
	checkDiffOption(t, MULTISET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "bar"`,
		`+ "baz"`)
	checkDiffOption(t, MULTISET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestMultisetPatch(t *testing.T) {
	checkPatchOption(t, MULTISET, `[]`, `[]`)
	checkPatchOption(t, MULTISET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkPatchOption(t, MULTISET, `[1,2]`, `[1,2]`)
	checkPatchOption(t, MULTISET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`,
		`+ 2`)
	checkPatchOption(t, MULTISET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkPatchOption(t, MULTISET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchOption(t, MULTISET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchOption(t, MULTISET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`- "foo"`,
		`+ "baz"`)
	checkPatchOption(t, MULTISET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "bar"`,
		`+ "baz"`)
	checkPatchOption(t, MULTISET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestMultisetPatchError(t *testing.T) {
	checkPatchErrorOption(t, MULTISET, `[]`,
		`@ [{}]`,
		`- 1`)
	checkPatchErrorOption(t, MULTISET, `[1]`,
		`@ [{}]`,
		`- 1`,
		`- 1`)
	checkPatchErrorOption(t, MULTISET, `[]`,
		`@ [{}]`,
		`- 1`,
		`+ 1`)
	checkPatchErrorOption(t, MULTISET, `[]`,
		`@ []`,
		`- {}`)
}
