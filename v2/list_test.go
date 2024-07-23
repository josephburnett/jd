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
	tests := []struct {
		a       string
		b       string
		diff    []string
		options []Option
	}{{
		a:    `[]`,
		b:    `[]`,
		diff: ss(),
	}, {
		a: `[1]`,
		b: `[]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`]`,
		),
	}, {
		a: `[[]]`,
		b: `[[1]]`,
		diff: ss(
			`@ [0,0]`,
			`[`,
			`+ 1`,
			`]`,
		),
	}, {
		a: `[1]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`+ 2`,
			`]`,
		),
	}, {
		a: `[]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 2`,
			`]`,
		),
	}, {
		a: `[[]]`,
		b: `[{}]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- []`,
			`+ {}`,
			`]`,
		),
	}, {
		a: `[{"a":[1]}]`,
		b: `[{"a":[2]}]`,
		diff: ss(
			`@ [0,"a",0]`,
			`[`,
			`- 1`,
			`+ 2`,
			`]`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,2]`,
		diff: ss(
			`@ [2]`,
			`  2`,
			`- 3`,
			`]`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,4,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`+ 4`,
			`  3`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,null,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`+ null`,
			`  3`,
		),
	}, {
		a: `[1,2]`,
		b: `[1,2,3,4]`,
		diff: ss(
			`@ [2]`,
			`  2`,
			`+ 3`,
			`+ 4`,
			`]`,
		),
	}, {
		a: `[]`,
		b: `[3,4,5]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 3`,
			`+ 4`,
			`+ 5`,
			`]`,
		),
	}, {
		a: `[null,null,null]`,
		b: `[]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- null`,
			`- null`,
			`- null`,
			`]`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,4,3]`,
		diff: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+ [1,4,3]`,
		),
		options: []Option{MERGE},
	}, {
		a: `[1,2,3]`,
		b: `{}`,
		diff: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+ {}`,
		),
		options: []Option{MERGE},
	}, {
		a: `[1,2,3,4]`,
		b: `[2,3]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
			`@ [2]`,
			`  3`,
			`- 4`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8]`,
		b: `[2,3,6,8]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
			`@ [2]`,
			`  3`,
			`- 4`,
			`- 5`,
			`  6`,
			`@ [3]`,
			`  6`,
			`- 7`,
			`  8`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[3,2,1]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 3`,
			`+ 2`,
			`  1`,
			`@ [3]`,
			`  1`,
			`- 2`,
			`- 3`,
			`]`,
		),
	}, {
		a: `[[1],[2],[3]]`,
		b: `[[2]]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- [1]`,
			`  [2]`,
			`@ [1]`,
			`  [2]`,
			`- [3]`,
			`]`,
		),
	}, {
		a: `[1,2,3,[4,5],6]`,
		b: `[2,3,[4],6]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
			`@ [2,1]`,
			`  4`,
			`- 5`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9]`,
		b: `[1,5,9]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`- 3`,
			`- 4`,
			`  5`,
			`@ [2]`,
			`  5`,
			`- 6`,
			`- 7`,
			`- 8`,
			`  9`,
		),
	}, {
		a: `[[[1,2,5,6]]]`,
		b: `[[[1,2,3,4,5,6]]]`,
		diff: ss(
			`@ [0,0,2]`,
			`  2`,
			`+ 3`,
			`+ 4`,
			`  5`,
		),
	}}

	for _, tt := range tests {
		t.Run(tt.a+tt.b, func(t *testing.T) {
			ctx := newTestContext(t)
			if len(tt.options) > 0 {
				ctx = ctx.withOptions(tt.options...)
			}
			checkDiff(ctx, tt.a, tt.b, tt.diff...)
		})
	}
}

func TestListPatch(t *testing.T) {
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
			`[`,
			`- 1`,
			`]`,
		),
	}, {
		a: `[[]]`,
		b: `[[1]]`,
		diff: ss(
			`@ [0, 0]`,
			`[`,
			`+ 1`,
			`]`,
		),
	}, {
		a: `[[]]`,
		b: `[[1]]`,
		diff: ss(
			`@ [0, -1]`,
			`[`,
			`+ 1`,
			`]`,
		),
	}, {
		a: `[1]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`+ 2`,
			`]`,
		),
	}, {
		a: `[]`,
		b: `[2]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 2`,
			`]`,
		),
	}, {
		a: `[]`,
		b: `[2]`,
		diff: ss(
			`@ [-1]`,
			`[`,
			`+ 2`,
			`]`,
		),
	}, {
		a: `[[]]`,
		b: `[{}]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- []`,
			`+ {}`,
			`]`,
		),
	}, {
		a: `[{"a":[1]}]`,
		b: `[{"a":[2]}]`,
		diff: ss(
			`@ [0,"a",0]`,
			`[`,
			`- 1`,
			`+ 2`,
			`]`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,2]`,
		diff: ss(
			`@ [2]`,
			`  2`,
			`- 3`,
			`]`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,4,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`+ 4`,
			`  3`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,null,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`+ null`,
			`  3`,
		),
	}, {
		a: `[]`,
		b: `[3,4,5]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 3`,
			`+ 4`,
			`+ 5`,
			`]`,
		),
	}, {
		a: `[2]`,
		b: `[1,2]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 1`,
			`  2`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[1,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`  3`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[2,3]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
		),
	}, {
		a: `[1,3]`,
		b: `[1,2,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`+ 2`,
			`  3`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[4,5,6]`,
		diff: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+ [4,5,6]`,
		),
	}, {
		a: `[1,2,3]`,
		b: ``,
		diff: ss(
			`^ {"Merge":true}`,
			`@ []`,
			`+`,
		),
	}, {
		a: `[1,2]`,
		b: `[1,2,3,4]`,
		diff: ss(
			`@ [2]`,
			`  2`,
			`+ 3`,
			`+ 4`,
			`]`,
		),
	}, {
		a: `[1,2,3,4]`,
		b: `[2,3]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
			`@ [2]`,
			`  3`,
			`- 4`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8]`,
		b: `[2,3,6,8]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
			`@ [2]`,
			`  3`,
			`- 4`,
			`- 5`,
			`  6`,
			`@ [3]`,
			`  6`,
			`- 7`,
			`  8`,
		),
	}, {
		a: `[1,2,3]`,
		b: `[3,2,1]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 3`,
			`+ 2`,
			`  1`,
			`@ [3]`,
			`  1`,
			`- 2`,
			`- 3`,
			`]`,
		),
	}, {
		a: `[[1],[2],[3]]`,
		b: `[[2]]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- [1]`,
			`  [2]`,
			`@ [1]`,
			`  [2]`,
			`- [3]`,
			`]`,
		),
	}, {
		a: `[1,2,3,[4,5],6]`,
		b: `[2,3,[4],6]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
			`@ [2,1]`,
			`  4`,
			`- 5`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9]`,
		b: `[1,5,9]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`- 3`,
			`- 4`,
			`  5`,
			`@ [2]`,
			`  5`,
			`- 6`,
			`- 7`,
			`- 8`,
			`  9`,
		),
	}, {
		a: `[[[1,2,5,6]]]`,
		b: `[[[1,2,3,4,5,6]]]`,
		diff: ss(
			`@ [0,0,2]`,
			`  2`,
			`+ 3`,
			`+ 4`,
			`  5`,
		),
	}}

	for _, tt := range tests {
		t.Run(tt.a+tt.b, func(t *testing.T) {
			ctx := newTestContext(t)
			checkPatch(ctx, tt.a, tt.b, tt.diff...)
		})
	}
}

func TestListPatchError(t *testing.T) {
	tests := []struct {
		a    string
		diff []string
	}{{
		a: `[]`,
		diff: ss(
			`@ ["a"]`,
			`+ 1`,
		),
	}, {
		a: `[]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
		),
	}, {
		a: `[]`,
		diff: ss(
			`@ [0]`,
			`- null`,
		),
	}, {
		a: `[1,2,3]`,
		diff: ss(
			`^ {"Merge":true}`,
			`@ [1]`,
			`+ 4`,
		),
	}, {
		a: `[1,2,3]`,
		diff: ss(
			`^ {"Merge":true}`,
			`@ [1]`,
			`- 2`,
			`+ 4`,
		),
	}, {
		a: `[1,2,3]`,
		diff: ss(
			`^ {"Merge":true}`,
			`@ [-1]`,
			`+ 4`,
		),
	}, {
		a: `[1,2,3]`,
		diff: ss(
			`@ [1]`,
			`  0`, // wrong before context
			`- 2`,
			`  3`,
		),
	}, {
		a: `[1,2,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`  -2`,
			`  4`, // wrong after context
		),
	}, {
		a: `[1,2,3]`,
		diff: ss(
			`@ [1]`,
			`[`, // wrong before context
			`- 2`,
			`  3`,
		),
	}, {
		a: `[1,2,3]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`]`, // wrong after context
		),
	}}

	for _, tt := range tests {
		t.Run(tt.a, func(t *testing.T) {
			ctx := newTestContext(t)
			checkPatchError(ctx, tt.a, tt.diff...)
		})
	}
}
