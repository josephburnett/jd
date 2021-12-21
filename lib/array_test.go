package jd

import (
	"testing"
)

func TestArrayJson(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		context *testContext
		a           string
		b           string
	}{
		{ctx, `[]`, `[]`},
		{ctx, ` [ ] `, `[]`},
		{ctx, `[1,2,3]`, `[1,2,3]`},
		{ctx, ` [1, 2, 3] `, `[1,2,3]`},
	}

	for _, tt := range tests {
		checkJson(tt.context, tt.a, tt.b)
	}
}

func TestArrayEqual(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		context *testContext
		a       string
		b       string
	}{
		{ctx, `[]`, `[]`},
		{ctx, `[1,2,3]`, `[1,2,3]`},
		{ctx, `[[]]`, `[[]]`},
		{ctx, `[{"a":1}]`, `[{"a":1}]`},
		{ctx, `[{"a":[]}]`, `[{"a":[]}]`},
	}

	for _, tt := range tests {
		checkEqual(tt.context, tt.a, tt.b)
	}
}

func TestArrayNotEqual(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		context *testContext
		a       string
		b       string
	}{
		{ctx, `[]`, `0`},
		{ctx, `[]`, `{}`},
		{ctx, `[]`, `[[]]`},
		{ctx, `[1,2,3]`, `[3,2,1]`},
	}

	for _, tt := range tests {
		checkNotEqual(tt.context, tt.a, tt.b)
	}
}

func TestArrayHash(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		context  *testContext
		a        string
		b        string
		wantSame bool
	}{
		{ctx, `[]`, `[]`, true},
		{ctx, `[1]`, `[]`, false},
		{ctx, `[1]`, `[1]`, true},
		{ctx, `[1]`, `[2]`, false},
		{ctx, `[[1]]`, `[[1]]`, true},
		{ctx, `[[1]]`, `[[[1]]]`, false},
	}

	for _, tt := range tests {
		checkHash(tt.context, tt.a, tt.b, tt.wantSame)
	}
}

func TestArrayDiff(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a         string
		b         string
		diff []string
	}{{
		a:`[]`,
		b:`[]`,
		diff: ss(),
	}, {
		a: `[1]`,
		b:`[]`,
		diff: ss(
			`@ [0]`,
			`- 1`,
		),
	}, {
		a: `[[]]`,
		b:`[[1]]`,
		diff: ss(
			`@ [0,-1]`,
			`+ 1`,
		),
	}, {
		a: `[1]`,
		b:`[2]`,
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
	}}

	for _, tt := range tests {
		checkDiff(ctx, tt.a, tt.b, tt.diff...)
	}
}

func TestArrayPatch(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		a         string
		b         string
		diff []string
	}{{
		a: `[]`,
		b: `[]`,
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
		b:`[[1]]`,
		diff: ss(
			`@ [0, 0]`,
			`+ 1`,
		),
	}, {
		a: `[[]]`,
		b:`[[1]]`,
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
	}}

	for _, tt := range tests {
		checkPatch(ctx, tt.a, tt.b, tt.diff...)
	}
}

func TestArrayPatchError(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		context   *testContext
		a         string
		diffLines []string
	}{
		{ctx, `[]`, []string {
			`@ ["a"]`,
		    `+ 1`,
		  },
		},
		{ctx, `[]`, []string {
			`@ [0]`,
		    `- 1`,
		  },
		},
		{ctx, `[]`, []string {
			`@ [0]`,
		    `- null`,
		  },
		},
		{ctx, `[1,2,3]`, []string {
			`@ [1]`,
		    `- 2`,
		  },
		},
		{ctx, `[1,2,3]`, []string {
			`@ [0]`,
		    `- 1`,
		  },
		},
		{ctx, `[1,3]`, []string {
			`@ [1]`,
		    `+ 2`,
		  },
		},
	}

	for _, tt := range tests {
		checkPatchError(tt.context, tt.a, tt.diffLines...)
	}
}
