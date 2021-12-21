package jd

import (
	"testing"
)

func TestObjectJson(t *testing.T) {
	ctx := newTestContext(t)
	checkJson(ctx, `{"a":1}`, `{"a":1}`)
	checkJson(ctx, ` { "a" : 1 } `, `{"a":1}`)
	checkJson(ctx, `{}`, `{}`)
}

func TestObjectEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkEqual(ctx, `{"a":1}`, `{"a":1}`)
	checkEqual(ctx, `{"a":1}`, `{"a":1.0}`)
	checkEqual(ctx, `{"a":[1,2]}`, `{"a":[1,2]}`)
	checkEqual(ctx, `{"a":"b"}`, `{"a":"b"}`)
}

func TestObjectNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	checkNotEqual(ctx, `{"a":1}`, `{"b":1}`)
	checkNotEqual(ctx, `{"a":[1,2]}`, `{"a":[2,1]}`)
	checkNotEqual(ctx, `{"a":"b"}`, `{"a":"c"}`)
}

// TODO: add unit test for object identity with setkeys metadata.
func TestObjectHash(t *testing.T) {
	ctx := newTestContext(t)
	checkHash(ctx, `{}`, `{}`, true)
	checkHash(ctx, `{"a":1}`, `{"a":1}`, true)
	checkHash(ctx, `{"a":1}`, `{"a":2}`, false)
	checkHash(ctx, `{"a":1}`, `{"b":1}`, false)
	checkHash(ctx, `{"a":1,"b":2}`, `{"b":2,"a":1}`, true)
}

func TestObjectDiff(t *testing.T) {
	ctx := newTestContext(t)
	checkDiff(ctx, `{}`, `{}`)
	checkDiff(ctx, `{"a":1}`, `{"a":1}`)
	checkDiff(ctx, `{"a":1}`, `{"a":2}`,
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkDiff(ctx, `{"":1}`, `{"":1}`)
	checkDiff(ctx, `{"":1}`, `{"a":2}`,
		`@ [""]`,
		`- 1`,
		`@ ["a"]`,
		`+ 2`)
	checkDiff(ctx, `{"a":{"b":{}}}`, `{"a":{"b":{"c":1},"d":2}}`,
		`@ ["a","b","c"]`,
		`+ 1`,
		`@ ["a","d"]`,
		`+ 2`)
	// regression test for issue #18
	checkDiff(ctx,
		`{"R": [{"I": [{"T": [{"V": "t","K": "N"},{"V": "T","K": "I"}]}]}]}`,
		`{"R": [{"I": [{"T": [{"V": "t","K": "N"},{"V": "Q","K": "C"},{"V": "T","K": "I"}]}]}]}`,
		`@ ["R",0,"I",0,"T",1,"K"]`,
		`- "I"`,
		`+ "C"`,
		`@ ["R",0,"I",0,"T",1,"V"]`,
		`- "T"`,
		`+ "Q"`,
		`@ ["R",0,"I",0,"T",-1]`,
		`+ {"K":"I","V":"T"}`)
}

func testObjectPatch(t *testing.T) {
	ctx := newTestContext(t)
	checkPatch(ctx, `{}`, `{}`)
	checkPatch(ctx, `{"a":1}`, `{"a":1}`)
	checkPatch(ctx, `{"a":1}`, `{"a":2}`,
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkPatch(ctx, `{"":1}`, `{"":1}`)
	checkPatch(ctx, `{"":1}`, `{"a":2}`,
		`@ [""]`,
		`- 1`,
		`@ ["a"]`,
		`+ 2`)
	checkPatch(ctx, `{"a":{"b":{}}}`, `{"a":{"b":{"c":1},"d":2}}`,
		`@ ["a","b","c"]`,
		`+ 1`,
		`@ ["a","d"]`,
		`+ 2`)
}

func testObjectPatchError(t *testing.T) {
	ctx := newTestContext(t)
	checkPatchError(ctx, `{}`,
		`@ ["a"]`,
		`- 1`)
	checkPatchError(ctx, `{"a":1}`,
		`@ ["a"]`,
		`+ 2`)
	checkPatchError(ctx, `{"a":1}`,
		`@ ["a"]`,
		`+ 1`)
}
