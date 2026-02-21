package jd

import (
	"fmt"
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
	},
	// >10 element arrays to exercise Myers diff code path.
	{
		// Single substitution at end (the #112 case).
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,2,3,4,5,6,7,8,9,10,99]`,
		diff: ss(
			`@ [10]`,
			`  10`,
			`- 11`,
			`+ 99`,
			`]`,
		),
	}, {
		// Single substitution at beginning.
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[99,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`+ 99`,
			`  2`,
		),
	}, {
		// Single substitution in middle.
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,2,3,4,5,99,7,8,9,10,11]`,
		diff: ss(
			`@ [5]`,
			`  5`,
			`- 6`,
			`+ 99`,
			`  7`,
		),
	}, {
		// Multiple scattered substitutions.
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,99,3,4,5,6,7,98,9,10,97]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`+ 99`,
			`  3`,
			`@ [7]`,
			`  7`,
			`- 8`,
			`+ 98`,
			`  9`,
			`@ [10]`,
			`  10`,
			`- 11`,
			`+ 97`,
			`]`,
		),
	}, {
		// Deletion from end.
		a: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		b: `[1,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(
			`@ [11]`,
			`  11`,
			`- 12`,
			`]`,
		),
	}, {
		// Deletion from beginning.
		a: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		b: `[2,3,4,5,6,7,8,9,10,11,12]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
		),
	}, {
		// Insertion at end.
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		diff: ss(
			`@ [11]`,
			`  11`,
			`+ 12`,
			`]`,
		),
	}, {
		// Insertion at beginning.
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[0,1,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 0`,
			`  1`,
		),
	}, {
		// Mixed insert + delete (different lengths).
		a: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		b: `[1,3,4,5,6,7,8,9,10,11,12,13]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`  3`,
			`@ [11]`,
			`  12`,
			`+ 13`,
			`]`,
		),
	}, {
		// Identical arrays (no diff expected).
		a:    `[1,2,3,4,5,6,7,8,9,10,11]`,
		b:    `[1,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(),
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
	},
	// >10 element arrays to exercise Myers diff code path.
	{
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,2,3,4,5,6,7,8,9,10,99]`,
		diff: ss(
			`@ [10]`,
			`  10`,
			`- 11`,
			`+ 99`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[99,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`+ 99`,
			`  2`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,2,3,4,5,99,7,8,9,10,11]`,
		diff: ss(
			`@ [5]`,
			`  5`,
			`- 6`,
			`+ 99`,
			`  7`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,99,3,4,5,6,7,98,9,10,97]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`+ 99`,
			`  3`,
			`@ [7]`,
			`  7`,
			`- 8`,
			`+ 98`,
			`  9`,
			`@ [10]`,
			`  10`,
			`- 11`,
			`+ 97`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		b: `[1,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(
			`@ [11]`,
			`  11`,
			`- 12`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		b: `[2,3,4,5,6,7,8,9,10,11,12]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`- 1`,
			`  2`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		diff: ss(
			`@ [11]`,
			`  11`,
			`+ 12`,
			`]`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11]`,
		b: `[0,1,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(
			`@ [0]`,
			`[`,
			`+ 0`,
			`  1`,
		),
	}, {
		a: `[1,2,3,4,5,6,7,8,9,10,11,12]`,
		b: `[1,3,4,5,6,7,8,9,10,11,12,13]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`  3`,
			`@ [11]`,
			`  12`,
			`+ 13`,
			`]`,
		),
	}, {
		a:    `[1,2,3,4,5,6,7,8,9,10,11]`,
		b:    `[1,2,3,4,5,6,7,8,9,10,11]`,
		diff: ss(),
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
	}, {
		// Additional edge cases for context verification vulnerability
		a: `[1,2,3,4]`,
		diff: ss(
			`@ [1]`,
			`  99`, // wrong before context - should be 1
			`- 2`,
			`+ 5`,
			`  3`,
		),
	}, {
		a: `[1,2,3,4]`,
		diff: ss(
			`@ [1]`,
			`  1`,
			`- 2`,
			`+ 5`,
			`  99`, // wrong after context - should be 3
		),
	}, {
		a: `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`,
		diff: ss(
			`@ [0]`,
			`  {"id":1,"name":"Charlie"}`, // wrong before context - name should be Alice
			`- {"id":1,"name":"Alice"}`,
			`+ {"id":3,"name":"Alice"}`,
			`  {"id":2,"name":"Bob"}`,
		),
	}, {
		a: `[10,20,30]`,
		diff: ss(
			`@ [1]`,
			`  10`,
			`- 20`,
			`+ 25`,
			`  999`, // fabricated after context - should be 30
		),
	}}

	for _, tt := range tests {
		t.Run(tt.a, func(t *testing.T) {
			ctx := newTestContext(t)
			checkPatchError(ctx, tt.a, tt.diff...)
		})
	}
}

func TestMyersRoundtrip(t *testing.T) {
	// makeArray builds a jsonList of sequential numbers [1, 2, ..., n].
	makeArray := func(n int) jsonList {
		a := make(jsonList, n)
		for i := range a {
			a[i] = jsonNumber(i + 1)
		}
		return a
	}

	// copyArray returns a shallow copy of a jsonList.
	copyArray := func(a jsonList) jsonList {
		c := make(jsonList, len(a))
		copy(c, a)
		return c
	}

	type scenario struct {
		name string
		a, b jsonList
	}

	var scenarios []scenario

	for _, size := range []int{11, 25, 50, 100} {
		base := makeArray(size)

		// Identical arrays.
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("identical_%d", size),
			a:    base,
			b:    copyArray(base),
		})

		// Single substitution at end.
		b := copyArray(base)
		b[size-1] = jsonNumber(-1)
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("sub_end_%d", size),
			a:    base,
			b:    b,
		})

		// Single substitution at beginning.
		b = copyArray(base)
		b[0] = jsonNumber(-1)
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("sub_begin_%d", size),
			a:    base,
			b:    b,
		})

		// Single substitution in middle.
		b = copyArray(base)
		b[size/2] = jsonNumber(-1)
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("sub_mid_%d", size),
			a:    base,
			b:    b,
		})

		// 10% scattered substitutions.
		b = copyArray(base)
		for i := 0; i < size; i += 10 {
			b[i] = jsonNumber(-i)
		}
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("sub_10pct_%d", size),
			a:    base,
			b:    b,
		})

		// 50% scattered substitutions.
		b = copyArray(base)
		for i := 0; i < size; i += 2 {
			b[i] = jsonNumber(-i)
		}
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("sub_50pct_%d", size),
			a:    base,
			b:    b,
		})

		// Deletion from end.
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("del_end_%d", size),
			a:    base,
			b:    copyArray(base[:size-1]),
		})

		// Deletion from beginning.
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("del_begin_%d", size),
			a:    base,
			b:    copyArray(base[1:]),
		})

		// Deletion from middle.
		b = append(copyArray(base[:size/2]), base[size/2+1:]...)
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("del_mid_%d", size),
			a:    base,
			b:    b,
		})

		// Insertion at end.
		b = append(copyArray(base), jsonNumber(-1))
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("ins_end_%d", size),
			a:    base,
			b:    b,
		})

		// Insertion at beginning.
		b = append(jsonList{jsonNumber(-1)}, base...)
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("ins_begin_%d", size),
			a:    base,
			b:    b,
		})

		// Insertion in middle.
		b = make(jsonList, 0, size+1)
		b = append(b, base[:size/2]...)
		b = append(b, jsonNumber(-1))
		b = append(b, base[size/2:]...)
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("ins_mid_%d", size),
			a:    base,
			b:    b,
		})

		// Mixed: delete from beginning, insert at end.
		b = append(copyArray(base[1:]), jsonNumber(-1))
		scenarios = append(scenarios, scenario{
			name: fmt.Sprintf("mixed_%d", size),
			a:    base,
			b:    b,
		})
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			d := sc.a.Diff(sc.b)
			patched, err := sc.a.Patch(d)
			if err != nil {
				t.Fatalf("patch failed: %v", err)
			}
			if !patched.Equals(sc.b) {
				t.Errorf("roundtrip failed: got %v, want %v",
					renderJson(patched), renderJson(sc.b))
			}
		})
	}
}
