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
		context   *testContext
		a         string
		b         string
		diffLines []string
	}{
		{ctx, `[]`, `[]`, []string {}},
		{ctx, `[1]`, `[]`, []string {
			`@ [0]`,
	        `- 1`,
		  },
	    },
		{ctx, `[[]]`, `[[1]]`, []string {
			`@ [0,0]`,
		    `+ 1`,
		  },
	    },
		{ctx, `[1]`, `[2]`, []string {
			`@ [0]`,
		    `- 1`,
		    `+ 2`,
		  },
	    },
		{ctx, `[]`, `[2]`, []string {
			`@ [0]`,
		    `+ 2`,
		  },
	    },
		{ctx, `[[]]`, `[{}]`, []string {
			`@ [0]`,
		    `- []`,
		    `+ {}`,
		  },
	    },
		{ctx, `[{"a":[1]}]`, `[{"a":[2]}]`, []string {
			`@ [0,"a",0]`,
		    `- 1`,
		    `+ 2`,
		  },
	    },
		{ctx, `[1,2,3]`, `[1,2]`, []string {
			`@ [2]`,
		    `- 3`,
		  },
	    },
		{ctx, `[1,2,3]`, `[1,4,3]`, []string {
			`@ [1]`,
		    `- 2`,
		    `+ 4`,
		  },
	    },
		{ctx, `[1,2,3]`, `[1,null,3]`, []string {
			`@ [1]`,
		    `- 2`,
		    `+ null`,
		  },
	    },
		{ctx, `[]`, `[3,4,5]`, []string {
			`@ [0]`,
		    `+ 3`,
		    `@ [1]`,
		    `+ 4`,
		    `@ [2]`,
		    `+ 5`,
		  },
	    },
	}

	for _, tt := range tests {
		checkDiff(tt.context, tt.a, tt.b, tt.diffLines...)
	}
}

func TestArrayPatch(t *testing.T) {
	ctx := newTestContext(t)
	tests := []struct {
		context   *testContext
		a         string
		b         string
		diffLines []string
	}{
		{ctx, `[]`, `[]`, []string {}},
		{ctx, `[1]`, `[]`, []string {
			`@ [0]`,
	        `- 1`,
		  },
	    },
		{ctx, `[[]]`, `[[1]]`, []string {
			`@ [0, 0]`,
			`+ 1`,
		  },
	    },
		{ctx, `[1]`, `[2]`, []string {
			`@ [0]`,
			`- 1`,
			`+ 2`,
		  },
	    },
		{ctx, `[1]`, `[2]`, []string {
			`@ [0]`,
			`- 1`,
			`+ 2`,
		  },
	    },
		{ctx, `[]`, `[2]`, []string {
			`@ [0]`,
		    `+ 2`,
		  },
	    },
		{ctx, `[[]]`, `[{}]`, []string {
			`@ [0]`,
		    `- []`,
		    `+ {}`,
		  },
	    },
		{ctx, `[{"a":[1]}]`, `[{"a":[2]}]`, []string {
			`@ [0,"a",0]`,
		    `- 1`,
		    `+ 2`,
		  },
	    },
		{ctx, `[1,2,3]`, `[1,2]`, []string {
			`@ [2]`,
		    `- 3`,
		  },
	    },
		{ctx, `[1,2,3]`, `[1,4,3]`, []string {
			`@ [1]`,
		    `- 2`,
		    `+ 4`,
		  },
	    },
		{ctx, `[1,2,3]`, `[1,null,3]`, []string {
			`@ [1]`,
		    `- 2`,
		    `+ null`,
		  },
	    },
		{ctx, `[1,2,3]`, `[1,null,3]`, []string {
			`@ [1]`,
		    `- 2`,
		    `+ null`,
		  },
	    },
		{ctx, "[]", "[3,4,5]", []string {
			"@ [0]",
		    "+ 3",
		    "@ [1]",
		    "+ 4",
		    "@ [2]",
		    "+ 5",
		  },
	    },
		
	}

	for _, tt := range tests {
		checkPatch(tt.context, tt.a, tt.b, tt.diffLines...)
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
