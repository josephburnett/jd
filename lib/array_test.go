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

	for _, test := range tests {
		checkJson(test.context, test.a, test.b)
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

	for _, test := range tests {
		checkEqual(test.context, test.a, test.b)
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

	for _, test := range tests {
		checkNotEqual(test.context, test.a, test.b)
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

	for _, test := range tests {
		checkHash(test.context, test.a, test.b, test.wantSame)
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

	for _, test := range tests {
		checkDiff(test.context, test.a, test.b, test.diffLines...)
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

	for _, test := range tests {
		checkPatch(test.context, test.a, test.b, test.diffLines...)
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

	for _, test := range tests {
		checkPatchError(test.context, test.a, test.diffLines...)
	}
}
