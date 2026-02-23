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

func TestListDiffLargeArraysNoCommon(t *testing.T) {
	// Two large arrays (>100 elements) with no common elements
	// This triggers shouldUseFastPath's large-array branch
	aStr := "["
	bStr := "["
	for i := 0; i < 110; i++ {
		if i > 0 {
			aStr += ","
			bStr += ","
		}
		aStr += fmt.Sprintf("%d", i)
		bStr += fmt.Sprintf("%d", i+1000)
	}
	aStr += "]"
	bStr += "]"
	aNode, _ := ReadJsonString(aStr)
	bNode, _ := ReadJsonString(bStr)
	d := aNode.Diff(bNode)
	// Should produce a diff (fast path: complete replacement)
	if len(d) == 0 {
		t.Fatal("expected non-empty diff for completely different large arrays")
	}
	// Verify round-trip: patch should produce b
	result, err := aNode.Patch(d)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equals(bNode) {
		t.Error("patch did not produce expected result")
	}
}

func TestListDiffContainerAfterChange(t *testing.T) {
	// Tests processcontainerDiffEvent: container diff after accumulated changes
	// [1, {"a":1}] vs [2, {"a":2}]: 1->2 is a replacement, then {"a":1}->{"a":2} is a container diff
	a, _ := ReadJsonString(`[1,{"a":1}]`)
	b, _ := ReadJsonString(`[2,{"a":2}]`)
	d := a.Diff(b)
	if len(d) == 0 {
		t.Fatal("expected non-empty diff")
	}
	result, err := a.Patch(d)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equals(b) {
		t.Error("patch did not produce expected result")
	}
}

func TestListDiffIdentical(t *testing.T) {
	// Tests finalizeCurrentDiff when processor is idle (identical lists)
	a, _ := ReadJsonString(`[1,2,3]`)
	b, _ := ReadJsonString(`[1,2,3]`)
	d := a.Diff(b)
	if len(d) != 0 {
		t.Errorf("expected empty diff for identical lists, got %v", d.Render())
	}
}

func TestListMultisetSameContainerType(t *testing.T) {
	// Triggers the multiset branch in sameContainerType
	a, _ := ReadJsonString(`[[1,2],[3,4]]`)
	b, _ := ReadJsonString(`[[1,3],[3,5]]`)
	d := a.Diff(b, MULTISET)
	// Should produce a diff
	result, err := a.Patch(d)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equals(b, MULTISET) {
		t.Error("multiset patch did not produce expected result")
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

func TestListPatchErrorPaths(t *testing.T) {
	tests := []struct {
		name string
		node jsonList
		diff Diff
	}{
		{
			name: "multiple values for root replace",
			node: jsonList{jsonNumber(1)},
			diff: Diff{{
				Path:   Path{},
				Remove: []JsonNode{jsonNumber(1), jsonNumber(2)},
			}},
		},
		{
			name: "recursive index out of bounds",
			node: jsonList{jsonNumber(1)},
			diff: Diff{{
				Path:   Path{PathIndex(5), PathIndex(0)},
				Remove: []JsonNode{jsonNumber(1)},
				Add:    []JsonNode{jsonNumber(2)},
			}},
		},
		{
			name: "no remove with strict strategy",
			node: jsonList{jsonNumber(1)},
			diff: Diff{{
				Path: Path{},
				Add:  []JsonNode{jsonNumber(2)},
			}},
		},
		{
			name: "recursive patch error",
			node: jsonList{jsonList{jsonNumber(1)}},
			diff: Diff{{
				Path:   Path{PathIndex(0), PathIndex(99)},
				Remove: []JsonNode{jsonNumber(1)},
				Add:    []JsonNode{jsonNumber(2)},
			}},
		},
		{
			name: "remove with append index",
			node: jsonList{jsonNumber(1)},
			diff: Diff{{
				Path:   Path{PathIndex(-1)},
				Remove: []JsonNode{jsonNumber(1)},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.node.Patch(tt.diff)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestSameContainerTypeSet(t *testing.T) {
	a, _ := ReadJsonString(`[1,2]`)
	b, _ := ReadJsonString(`[3,4]`)
	opts := refine(newOptions([]Option{SET}), nil)
	if !sameContainerType(a, b, opts) {
		t.Error("expected same container type for two arrays with SET")
	}
}

func TestSameContainerTypeMultiset(t *testing.T) {
	a, _ := ReadJsonString(`[1,2]`)
	b, _ := ReadJsonString(`[3,4]`)
	opts := refine(newOptions([]Option{MULTISET}), nil)
	if !sameContainerType(a, b, opts) {
		t.Error("expected same container type for two arrays with MULTISET")
	}
}

func TestListPatchEdgeCases(t *testing.T) {
	// Multiple remove/add at root path
	l := jsonList{jsonNumber(1), jsonNumber(2)}
	_, err := l.patch(nil, Path{}, nil,
		[]JsonNode{jsonNumber(1), jsonNumber(2)},
		[]JsonNode{jsonNumber(3), jsonNumber(4)},
		nil, strictPatchStrategy)
	if err == nil {
		t.Fatal("expected error for multiple remove/add at root path")
	}
	// Before context with bIndex == -1 and void (should continue)
	l2 := jsonList{jsonNumber(1), jsonNumber(2)}
	result, err := l2.patch(nil, Path{PathIndex(0)}, []JsonNode{voidNode{}},
		[]JsonNode{jsonNumber(1)},
		[]JsonNode{jsonNumber(3)},
		nil, strictPatchStrategy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// After context with aIndex == len(l) and void (should continue)
	l3 := jsonList{jsonNumber(1), jsonNumber(2)}
	result, err = l3.patch(nil, Path{PathIndex(1)}, nil,
		[]JsonNode{jsonNumber(2)},
		nil,
		[]JsonNode{voidNode{}}, strictPatchStrategy)
	if err != nil {
		t.Fatalf("unexpected error for after void context at end: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// After context out of bounds (non-void) should error
	l4 := jsonList{jsonNumber(1)}
	_, err = l4.patch(nil, Path{PathIndex(0)}, nil,
		[]JsonNode{jsonNumber(1)},
		nil,
		[]JsonNode{jsonNumber(99)}, strictPatchStrategy)
	if err == nil {
		t.Fatal("expected error for after context out of bounds")
	}
}
