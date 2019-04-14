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
	checkDiffMetadata(t, MULTISET, `[]`, `[]`)
	checkDiffMetadata(t, MULTISET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkDiffMetadata(t, MULTISET, `[1,2]`, `[1,2]`)
	checkDiffMetadata(t, MULTISET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`,
		`+ 2`)
	checkDiffMetadata(t, MULTISET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkDiffMetadata(t, MULTISET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffMetadata(t, MULTISET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkDiffMetadata(t, MULTISET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`- "foo"`,
		`+ "baz"`)
	checkDiffMetadata(t, MULTISET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "bar"`,
		`+ "baz"`)
	checkDiffMetadata(t, MULTISET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestMultisetPatch(t *testing.T) {
	checkPatchMetadata(t, MULTISET, `[]`, `[]`)
	checkPatchMetadata(t, MULTISET, `[1]`, `[1,2]`,
		`@ [{}]`,
		`+ 2`)
	checkPatchMetadata(t, MULTISET, `[1,2]`, `[1,2]`)
	checkPatchMetadata(t, MULTISET, `[1]`, `[1,2,2]`,
		`@ [{}]`,
		`+ 2`,
		`+ 2`)
	checkPatchMetadata(t, MULTISET, `[1,2,3]`, `[1,3]`,
		`@ [{}]`,
		`- 2`)
	checkPatchMetadata(t, MULTISET, `[{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchMetadata(t, MULTISET, `[{"a":1},{"a":1}]`, `[{"a":2}]`,
		`@ [{}]`,
		`- {"a":1}`,
		`- {"a":1}`,
		`+ {"a":2}`)
	checkPatchMetadata(t, MULTISET, `["foo","foo","bar"]`, `["baz"]`,
		`@ [{}]`,
		`- "bar"`,
		`- "foo"`,
		`- "foo"`,
		`+ "baz"`)
	checkPatchMetadata(t, MULTISET, `["foo"]`, `["bar","baz","bar"]`,
		`@ [{}]`,
		`- "foo"`,
		`+ "bar"`,
		`+ "bar"`,
		`+ "baz"`)
	checkPatchMetadata(t, MULTISET, `{}`, `[]`,
		`@ []`,
		`- {}`,
		`+ []`)
}

func TestMultisetPatchError(t *testing.T) {
	checkPatchErrorMetadata(t, MULTISET, `[]`,
		`@ [{}]`,
		`- 1`)
	checkPatchErrorMetadata(t, MULTISET, `[1]`,
		`@ [{}]`,
		`- 1`,
		`- 1`)
	checkPatchErrorMetadata(t, MULTISET, `[]`,
		`@ [{}]`,
		`- 1`,
		`+ 1`)
	checkPatchErrorMetadata(t, MULTISET, `[]`,
		`@ []`,
		`- {}`)
}
