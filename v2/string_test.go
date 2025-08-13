package jd

import (
	"testing"
)

func TestStringJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    `""`,
			expected: `""`,
		},
		{
			name:     "empty string with whitespace",
			input:    ` "" `,
			expected: `""`,
		},
		{
			name:     "string with escaped quote",
			input:    `"\""`,
			expected: `"\""`,
		},
		{
			name:     "simple string",
			input:    `"hello"`,
			expected: `"hello"`,
		},
		{
			name:     "string with whitespace around",
			input:    ` "hello" `,
			expected: `"hello"`,
		},
		{
			name:     "string with numbers",
			input:    `"123"`,
			expected: `"123"`,
		},
		{
			name:     "string with special characters",
			input:    `"hello\nworld"`,
			expected: `"hello\nworld"`,
		},
		{
			name:     "string with unicode",
			input:    `"héllo wørld"`,
			expected: `"héllo wørld"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkJson(ctx, tt.input, tt.expected)
		})
	}
}

func TestStringEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "empty strings equal",
			a:    `""`,
			b:    `""`,
		},
		{
			name: "same single character strings",
			a:    `"a"`,
			b:    `"a"`,
		},
		{
			name: "same number strings",
			a:    `"123"`,
			b:    `"123"`,
		},
		{
			name: "same multi-word strings",
			a:    `"hello world"`,
			b:    `"hello world"`,
		},
		{
			name: "same strings with whitespace around input",
			a:    ` "hello" `,
			b:    `"hello"`,
		},
		{
			name: "same strings with escape sequences",
			a:    `"hello\nworld"`,
			b:    `"hello\nworld"`,
		},
		{
			name: "same strings with unicode",
			a:    `"héllo wørld"`,
			b:    `"héllo wørld"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestStringNotEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "empty string not equal to single character",
			a:    `""`,
			b:    `"a"`,
		},
		{
			name: "empty string not equal to empty array",
			a:    `""`,
			b:    `[]`,
		},
		{
			name: "empty string not equal to empty object",
			a:    `""`,
			b:    `{}`,
		},
		{
			name: "empty string not equal to number 0",
			a:    `""`,
			b:    `0`,
		},
		{
			name: "empty string not equal to null",
			a:    `""`,
			b:    `null`,
		},
		{
			name: "empty string not equal to false",
			a:    `""`,
			b:    `false`,
		},
		{
			name: "empty string not equal to true",
			a:    `""`,
			b:    `true`,
		},
		{
			name: "empty string not equal to void",
			a:    `""`,
			b:    ``,
		},
		{
			name: "string not equal to number with same digits",
			a:    `"123"`,
			b:    `123`,
		},
		{
			name: "string not equal to boolean with same text",
			a:    `"true"`,
			b:    `true`,
		},
		{
			name: "string not equal to boolean false with same text",
			a:    `"false"`,
			b:    `false`,
		},
		{
			name: "string not equal to null with same text",
			a:    `"null"`,
			b:    `null`,
		},
		{
			name: "different strings",
			a:    `"hello"`,
			b:    `"world"`,
		},
		{
			name: "case sensitive strings",
			a:    `"Hello"`,
			b:    `"hello"`,
		},
		{
			name: "strings with different whitespace",
			a:    `"hello world"`,
			b:    `"hello  world"`,
		},
		{
			name: "string not equal to array containing same string",
			a:    `"hello"`,
			b:    `["hello"]`,
		},
		{
			name: "string not equal to object with same string value",
			a:    `"hello"`,
			b:    `{"key": "hello"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkNotEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestStringHash(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		wantSame bool
	}{
		{
			name:     "empty strings hash same",
			a:        `""`,
			b:        `""`,
			wantSame: true,
		},
		{
			name:     "identical strings hash same",
			a:        `"abc"`,
			b:        `"abc"`,
			wantSame: true,
		},
		{
			name:     "empty string hashes different from space",
			a:        `""`,
			b:        `" "`,
			wantSame: false,
		},
		{
			name:     "different strings hash different",
			a:        `"abc"`,
			b:        `"123"`,
			wantSame: false,
		},
		{
			name:     "string hashes different from number",
			a:        `"123"`,
			b:        `123`,
			wantSame: false,
		},
		{
			name:     "string hashes different from boolean",
			a:        `"true"`,
			b:        `true`,
			wantSame: false,
		},
		{
			name:     "string hashes different from null",
			a:        `"null"`,
			b:        `null`,
			wantSame: false,
		},
		{
			name:     "string hashes different from void",
			a:        `"hello"`,
			b:        ``,
			wantSame: false,
		},
		{
			name:     "string hashes different from empty array",
			a:        `"[]"`,
			b:        `[]`,
			wantSame: false,
		},
		{
			name:     "string hashes different from empty object",
			a:        `"{}"`,
			b:        `{}`,
			wantSame: false,
		},
		{
			name:     "case sensitive strings hash different",
			a:        `"Hello"`,
			b:        `"hello"`,
			wantSame: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkHash(ctx, tt.a, tt.b, tt.wantSame)
		})
	}
}

func TestStringDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected []string
		options  []Option
	}{
		{
			name:     "no diff between same strings",
			a:        `""`,
			b:        `""`,
			expected: []string{},
		},
		{
			name:     "no diff between identical non-empty strings",
			a:        `"hello"`,
			b:        `"hello"`,
			expected: []string{},
		},
		{
			name: "diff empty string to number",
			a:    `""`,
			b:    `1`,
			expected: []string{
				`@ []`,
				`- ""`,
				`+ 1`,
			},
		},
		{
			name: "diff null to string",
			a:    `null`,
			b:    `"abc"`,
			expected: []string{
				`@ []`,
				`- null`,
				`+ "abc"`,
			},
		},
		{
			name: "diff string to null",
			a:    `"abc"`,
			b:    `null`,
			expected: []string{
				`@ []`,
				`- "abc"`,
				`+ null`,
			},
		},
		{
			name: "diff string to boolean",
			a:    `"true"`,
			b:    `true`,
			expected: []string{
				`@ []`,
				`- "true"`,
				`+ true`,
			},
		},
		{
			name: "diff string to number",
			a:    `"123"`,
			b:    `123`,
			expected: []string{
				`@ []`,
				`- "123"`,
				`+ 123`,
			},
		},
		{
			name: "diff string to array",
			a:    `"hello"`,
			b:    `["hello"]`,
			expected: []string{
				`@ []`,
				`- "hello"`,
				`+ ["hello"]`,
			},
		},
		{
			name: "diff string to object",
			a:    `"hello"`,
			b:    `{"key": "hello"}`,
			expected: []string{
				`@ []`,
				`- "hello"`,
				`+ {"key":"hello"}`,
			},
		},
		{
			name: "diff string to void",
			a:    `"hello"`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- "hello"`,
			},
		},
		{
			name: "diff void to string",
			a:    ``,
			b:    `"hello"`,
			expected: []string{
				`@ []`,
				`+ "hello"`,
			},
		},
		{
			name: "diff between different strings",
			a:    `"hello"`,
			b:    `"world"`,
			expected: []string{
				`@ []`,
				`- "hello"`,
				`+ "world"`,
			},
		},
		{
			name:    "merge diff string to string",
			a:       `"def"`,
			b:       `"abc"`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ "abc"`,
			},
		},
		{
			name:    "merge diff string to void",
			a:       `"abc"`,
			b:       ``,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:    "merge diff string to number",
			a:       `"123"`,
			b:       `456`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ 456`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			if len(tt.options) > 0 {
				ctx = ctx.withOptions(tt.options...)
			}
			checkDiff(ctx, tt.a, tt.b, tt.expected...)
		})
	}
}

func TestStringPatch(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		expected string
		patch    []string
	}{
		{
			name:     "no patch on empty string",
			initial:  `""`,
			expected: `""`,
			patch:    []string{},
		},
		{
			name:     "no patch on string",
			initial:  `"hello"`,
			expected: `"hello"`,
			patch:    []string{},
		},
		{
			name:     "patch empty string to number",
			initial:  `""`,
			expected: `1`,
			patch: []string{
				`@ []`,
				`- ""`,
				`+ 1`,
			},
		},
		{
			name:     "patch null to string",
			initial:  `null`,
			expected: `"abc"`,
			patch: []string{
				`@ []`,
				`- null`,
				`+ "abc"`,
			},
		},
		{
			name:     "patch string to null",
			initial:  `"abc"`,
			expected: `null`,
			patch: []string{
				`@ []`,
				`- "abc"`,
				`+ null`,
			},
		},
		{
			name:     "patch string to boolean",
			initial:  `"true"`,
			expected: `true`,
			patch: []string{
				`@ []`,
				`- "true"`,
				`+ true`,
			},
		},
		{
			name:     "patch string to number",
			initial:  `"123"`,
			expected: `123`,
			patch: []string{
				`@ []`,
				`- "123"`,
				`+ 123`,
			},
		},
		{
			name:     "patch string to array",
			initial:  `"hello"`,
			expected: `["hello"]`,
			patch: []string{
				`@ []`,
				`- "hello"`,
				`+ ["hello"]`,
			},
		},
		{
			name:     "patch string to object",
			initial:  `"hello"`,
			expected: `{"key":"hello"}`,
			patch: []string{
				`@ []`,
				`- "hello"`,
				`+ {"key":"hello"}`,
			},
		},
		{
			name:     "patch string to void",
			initial:  `"hello"`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- "hello"`,
			},
		},
		{
			name:     "patch void to string",
			initial:  ``,
			expected: `"hello"`,
			patch: []string{
				`@ []`,
				`+ "hello"`,
			},
		},
		{
			name:     "patch between different strings",
			initial:  `"hello"`,
			expected: `"world"`,
			patch: []string{
				`@ []`,
				`- "hello"`,
				`+ "world"`,
			},
		},
		{
			name:     "merge patch string to string",
			initial:  `"def"`,
			expected: `"abc"`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ "abc"`,
			},
		},
		{
			name:     "merge patch string to void",
			initial:  `"abc"`,
			expected: ``,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:     "merge patch string to number",
			initial:  `"123"`,
			expected: `456`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ 456`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkPatch(ctx, tt.initial, tt.expected, tt.patch...)
		})
	}
}

func TestStringPatchError(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		patch   []string
	}{
		{
			name:    "patch empty string with wrong remove value",
			initial: `""`,
			patch: []string{
				`@ []`,
				`- "a"`,
				`+ ""`,
			},
		},
		{
			name:    "patch string with wrong remove value",
			initial: `"hello"`,
			patch: []string{
				`@ []`,
				`- "world"`,
				`+ "hello"`,
			},
		},
		{
			name:    "patch null with add string but no remove",
			initial: `null`,
			patch: []string{
				`@ []`,
				`+ "a"`,
			},
		},
		{
			name:    "patch string with wrong type remove value",
			initial: `"123"`,
			patch: []string{
				`@ []`,
				`- 123`,
			},
		},
		{
			name:    "patch string with boolean remove value",
			initial: `"true"`,
			patch: []string{
				`@ []`,
				`- true`,
			},
		},
		{
			name:    "patch string with null remove value",
			initial: `"hello"`,
			patch: []string{
				`@ []`,
				`- null`,
			},
		},
		{
			name:    "patch string with array remove value",
			initial: `"hello"`,
			patch: []string{
				`@ []`,
				`- []`,
			},
		},
		{
			name:    "patch string with object remove value",
			initial: `"hello"`,
			patch: []string{
				`@ []`,
				`- {}`,
			},
		},
		{
			name:    "invalid merge patch string with remove and add",
			initial: `"hello"`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`- "hello"`,
				`+ "world"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkPatchError(ctx, tt.initial, tt.patch...)
		})
	}
}
