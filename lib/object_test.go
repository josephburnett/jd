package jd

import (
	"testing"
)

func TestObjectJson(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct{
		a string
		b string
	}{{
		a:`{"a":1}`,
		b: `{"a":1}`,
	}, {
		a: ` { "a" : 1 } `,
		b:`{"a":1}`,
	}, {
		a:`{}`,
		b:`{}`,
	}}

	for _, tt := range tests {
		checkJson(ctx, tt.a, tt.b)
	}
}

func TestObjectEqual(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct{
		a string
		b string
	}{{
		a:`{"a":1}`,
		b:`{"a":1}`,
	}, {
		a: `{"a":1}`,
		b:`{"a":1.0}`,
	}, {
		a: `{"a":[1,2]}`,
		b: `{"a":[1,2]}`,
	}, {
		a: `{"a":"b"}`,
		b: `{"a":"b"}`,
	}}

	for _, tt := range tests {
		checkEqual(ctx, tt.a, tt.b)
	}
}

func TestObjectNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct{
		a string
		b string
	}{{
		a: `{"a":1}`,
		b: `{"b":1}`,
	}, {
		a: `{"a":[1,2]}`,
		b:`{"a":[2,1]}`,
	}, {
		a: `{"a":"b"}`,
		b:`{"a":"c"}`,
	}}

	for _, tt := range tests {
		checkNotEqual(ctx, tt.a, tt.b)
	}
}

// TODO: add unit test for object identity with setkeys metadata.
func TestObjectHash(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct{
		a string
		b string
		wantEqual bool
	}{{
		a: `{}`,
		b:`{}`,
		wantEqual: true,
	}, {
		a: `{"a":1}`,
		b:`{"a":1}`,
		wantEqual: true,
	}, {
		a: `{"a":1}`,
		b:`{"a":2}`,
		wantEqual: false,
	}, {
		a: `{"a":1}`,
		b:`{"b":1}`,
		wantEqual: false,
	}, {
		a: `{"a":1,"b":2}`,
		b: `{"b":2,"a":1}`,
		wantEqual: true,
	}}

	for _, tt := range tests {
		checkHash(ctx, tt.a, tt.b, tt.wantEqual)
	}
}

func TestObjectDiff(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct{
		a string
		b string
		diff []string
	}{{
		a: `{}`,
		b: `{}`,
	}, {
		a: `{"a":1}`,
		b: `{"a":1}`,
	}, {
		a: `{"a":1}`,
		b: `{"a":2}`,
		diff: ss(
			`@ ["a"]`,
			`- 1`,
			`+ 2`,
		),
	}, {
		a: `{"":1}`,
		b: `{"":1}`,
	}, {
		a: `{"":1}`,
		b: `{"a":2}`,
		diff: ss(
			`@ [""]`,
			`- 1`,
			`@ ["a"]`,
			`+ 2`,
		),
	}, {
		a: `{"a":{"b":{}}}`,
		b: `{"a":{"b":{"c":1},"d":2}}`,
		diff: ss(
			`@ ["a","b","c"]`,
			`+ 1`,
			`@ ["a","d"]`,
			`+ 2`,
		),
	}, {
		// regression test for issue #18
		a: `{"R": [{"I": [{"T": [{"V": "t","K": "N"},{"V": "T","K": "I"}]}]}]}`,
		b: `{"R": [{"I": [{"T": [{"V": "t","K": "N"},{"V": "Q","K": "C"},{"V": "T","K": "I"}]}]}]}`,
		diff: ss(
			`@ ["R",0,"I",0,"T",1,"K"]`,
			`- "I"`,
			`+ "C"`,
			`@ ["R",0,"I",0,"T",1,"V"]`,
			`- "T"`,
			`+ "Q"`,
			`@ ["R",0,"I",0,"T",-1]`,
			`+ {"K":"I","V":"T"}`,
		),
	}}

	for _, tt := range tests {
		checkDiff(ctx, tt.a, tt.b, tt.diff...)
	}
}

func testObjectPatch(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct{
		a string
		b string
		diff []string
	}{{
		a: `{}`,
		b: `{}`,
	}, {
		a: `{"a":1}`,
		b: `{"a":1}`,
	}, {
		a: `{"a":1}`,
		b: `{"a":2}`,
		diff: ss(
			`@ ["a"]`,
			`- 1`,
			`+ 2`,
		),
	}, {
		a: `{"":1}`,
		b: `{"":1}`,
	}, {
		a: `{"":1}`,
		b: `{"a":2}`,
		diff: ss(
			`@ [""]`,
			`- 1`,
			`@ ["a"]`,
			`+ 2`,
		),
	}, {
		a: `{"a":{"b":{}}}`,
		b: `{"a":{"b":{"c":1},"d":2}}`,
		diff: ss(
			`@ ["a","b","c"]`,
			`+ 1`,
			`@ ["a","d"]`,
			`+ 2`,
		),
	}}

	for _, tt := range tests {
		checkPatch(ctx, tt.a, tt.b, tt.diff...)
	}
}

func testObjectPatchError(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct{
		a string
		diff []string
	}{{
		a: `{}`,
		diff: ss(
			`@ ["a"]`,
			`- 1`,
		),
	}, {
		a: `{"a":1}`,
		diff: ss(
			`@ ["a"]`,
			`+ 2`,
		),
	}, {
		a: `{"a":1}`,
		diff: ss(
			`@ ["a"]`,
			`+ 1`,
		),
	}}

	for _, tt := range tests {
		checkPatchError(ctx, tt.a, tt.diff...)
	}
}
