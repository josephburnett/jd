package jd

import (
	"testing"
)

func TestListJson(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a string
		b string
	}{
		{`[]`, `[]`},
		{` [ ] `, `[]`},
		{`[1,2,3]`, `[1,2,3]`},
		{` [1, 2, 3] `, `[1,2,3]`},
	}

	for _, tt := range tests {
		checkJson(ctx, tt.a, tt.b)
	}
}

func TestListEqual(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a string
		b string
	}{
		{`[]`, `[]`},
		{`[1,2,3]`, `[1,2,3]`},
		{`[[]]`, `[[]]`},
		{`[{"a":1}]`, `[{"a":1}]`},
		{`[{"a":[]}]`, `[{"a":[]}]`},
	}

	for _, tt := range tests {
		checkEqual(ctx, tt.a, tt.b)
	}
}

func TestListNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a string
		b string
	}{
		{`[]`, `0`},
		{`[]`, `{}`},
		{`[]`, `[[]]`},
		{`[1,2,3]`, `[3,2,1]`},
	}

	for _, tt := range tests {
		checkNotEqual(ctx, tt.a, tt.b)
	}
}

func TestListHash(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a        string
		b        string
		wantSame bool
	}{
		{`[]`, `[]`, true},
		{`[1]`, `[]`, false},
		{`[1]`, `[1]`, true},
		{`[1]`, `[2]`, false},
		{`[[1]]`, `[[1]]`, true},
		{`[[1]]`, `[[[1]]]`, false},
	}

	for _, tt := range tests {
		checkHash(ctx, tt.a, tt.b, tt.wantSame)
	}
}

func TestListDiff(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a    string
		b    string
		diff []string
	}{{
		a:    `[]`,
		b:    `[]`,
		diff: ss(),
	}, {
		a: `[1]`,
		b: `[]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
		),
	}, {
		a: `[[]]`,
		b: `[[1]]`,
		diff: ss(
			`@ [0,-1]`,
			`+ 1`,
		),
	}, {
		a: `[1]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
			`+ 2`,
		),
	}, {
		a: `[]`,
		b: `[2]`,
		diff: ss(
			`@ [-1]`,
			`+ 2`,
		),
	}, {
		a: `[[]]`,
		b: `[{}]`,
		diff: ss(
			`@ [0]`,
			`- []`,
			`+ {}`,
		),
	}, {
		a: `[{"a":[1]}]`,
		b: `[{"a":[2]}]`,
		diff: ss(
			`@ [0,"a",0]`,
			`- 1`,
			`+ 2`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,2]`,
		diff: ss(
			`@ [2]`,
			`- 3`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,4,3]`,
		diff: ss(
			`@ [1]`,
			`- 2`,
			`+ 4`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,null,3]`,
		diff: ss(
			`@ [1]`,
			`- 2`,
			`+ null`,
		),
	}, {
		a: `[1,2]`,
		b: `[1,2,3,4]`,
		diff: ss(
			`@ [-1]`,
			`+ 3`,
			`@ [-1]`,
			`+ 4`,
		),
	}, {
		a: `[]`,
		b: `[3,4,5]`,
		diff: ss(
			`@ [-1]`,
			`+ 3`,
			`@ [-1]`,
			`+ 4`,
			`@ [-1]`,
			`+ 5`,
		),
	}, {
		a: `[null,null,null]`,
		b: `[]`,
		diff: ss(
			`@ [2]`,
			`- null`,
			`@ [1]`,
			`- null`,
			`@ [0]`,
			`- null`,
		),
	}}

	for _, tt := range tests {
		checkDiff(ctx, tt.a, tt.b, tt.diff...)
	}
}

func TestListPatch(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a    string
		b    string
		diff []string
	}{{
		a:    `[]`,
		b:    `[]`,
		diff: ss(),
	}, {
		a: `[1]`,
		b: `[]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
		),
	}, {
		a: `[[]]`,
		b: `[[1]]`,
		diff: ss(
			`@ [0, 0]`,
			`+ 1`,
		),
	}, {
		a: `[[]]`,
		b: `[[1]]`,
		diff: ss(
			`@ [0, -1]`,
			`+ 1`,
		),
	}, {
		a: `[1]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
			`+ 2`,
		),
	}, {
		a: `[1]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
			`+ 2`,
		),
	}, {
		a: `[]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`+ 2`,
		),
	}, {
		a: `[]`,
		b: `[2]`,
		diff: ss(
			`@ [-1]`,
			`+ 2`,
		),
	}, {
		a: `[[]]`,
		b: `[{}]`,
		diff: ss(
			`@ [0]`,
			`- []`,
			`+ {}`,
		),
	}, {
		a: `[{"a":[1]}]`,
		b: `[{"a":[2]}]`,
		diff: ss(
			`@ [0,"a",0]`,
			`- 1`,
			`+ 2`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,2]`,
		diff: ss(
			`@ [2]`,
			`- 3`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,4,3]`,
		diff: ss(
			`@ [1]`,
			`- 2`,
			`+ 4`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,null,3]`,
		diff: ss(
			`@ [1]`,
			`- 2`,
			`+ null`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,null,3]`,
		diff: ss(
			`@ [1]`,
			`- 2`,
			`+ null`,
		),
	}, {
		a: "[]",
		b: "[3,4,5]",
		diff: ss(
			"@ [0]",
			"+ 3",
			"@ [1]",
			"+ 4",
			"@ [2]",
			"+ 5",
		),
	}, {
		a: "[]",
		b: "[3,4,5]",
		diff: ss(
			"@ [-1]",
			"+ 3",
			"@ [-1]",
			"+ 4",
			"@ [-1]",
			"+ 5",
		),
	}, {
		a: "[2]",
		b: "[1,2]",
		diff: ss(
			"@ [0]",
			"+ 1",
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,3]`,
		diff: ss(
			`@ [1]`,
			`- 2`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[2,3]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
		),
	}, {
		a: `[1,3]`,
		b: `[1,2,3]`,
		diff: ss(
			`@ [1]`,
			`+ 2`,
		),
	}}

	for _, tt := range tests {
		checkPatch(ctx, tt.a, tt.b, tt.diff...)
	}
}

func TestListPatchError(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a         string
		diffLines []string
	}{{
		`[]`,
		[]string{
			`@ ["a"]`,
			`+ 1`,
		},
	}, {
		`[]`,
		[]string{
			`@ [0]`,
			`- 1`,
		},
	}, {
		`[]`,
		[]string{
			`@ [0]`,
			`- null`,
		},
	}}

	for _, tt := range tests {
		checkPatchError(ctx, tt.a, tt.diffLines...)
	}
}
