package jd

import (
	"testing"
)

func TestVoidJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "void value",
			input:    ``,
			expected: ``,
		},
		{
			name:     "whitespace only treated as void",
			input:    `   `,
			expected: ``,
		},
		{
			name:     "tabs and spaces treated as void",
			input:    "	  	",
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkJson(ctx, tt.input, tt.expected)
		})
	}
}

func TestVoidEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "void equals void",
			a:    ``,
			b:    ``,
		},
		{
			name: "void equals whitespace",
			a:    ``,
			b:    `   `,
		},
		{
			name: "whitespace equals void",
			a:    `   `,
			b:    ``,
		},
		{
			name: "whitespace equals whitespace",
			a:    `   `,
			b:    `	 	`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestVoidNotEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "void not equal to null",
			a:    ``,
			b:    `null`,
		},
		{
			name: "void not equal to number 0",
			a:    ``,
			b:    `0`,
		},
		{
			name: "void not equal to empty array",
			a:    ``,
			b:    `[]`,
		},
		{
			name: "void not equal to empty object",
			a:    ``,
			b:    `{}`,
		},
		{
			name: "void not equal to false",
			a:    ``,
			b:    `false`,
		},
		{
			name: "void not equal to true",
			a:    ``,
			b:    `true`,
		},
		{
			name: "void not equal to empty string",
			a:    ``,
			b:    `""`,
		},
		{
			name: "void not equal to string with space",
			a:    ``,
			b:    `" "`,
		},
		{
			name: "void not equal to number 1",
			a:    ``,
			b:    `1`,
		},
		{
			name: "void not equal to negative number",
			a:    ``,
			b:    `-1`,
		},
		{
			name: "void not equal to decimal number",
			a:    ``,
			b:    `0.5`,
		},
		{
			name: "void not equal to string hello",
			a:    ``,
			b:    `"hello"`,
		},
		{
			name: "void not equal to array with null",
			a:    ``,
			b:    `[null]`,
		},
		{
			name: "void not equal to object with null value",
			a:    ``,
			b:    `{"key": null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkNotEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestVoidHash(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		wantSame bool
	}{
		{
			name:     "void hashes same as void",
			a:        ``,
			b:        ``,
			wantSame: true,
		},
		{
			name:     "void hashes same as whitespace",
			a:        ``,
			b:        `   `,
			wantSame: true,
		},
		{
			name:     "void hashes different from null",
			a:        ``,
			b:        `null`,
			wantSame: false,
		},
		{
			name:     "void hashes different from false",
			a:        ``,
			b:        `false`,
			wantSame: false,
		},
		{
			name:     "void hashes different from true",
			a:        ``,
			b:        `true`,
			wantSame: false,
		},
		{
			name:     "void hashes different from number 0",
			a:        ``,
			b:        `0`,
			wantSame: false,
		},
		{
			name:     "void hashes different from empty string",
			a:        ``,
			b:        `""`,
			wantSame: false,
		},
		{
			name:     "void hashes different from string hello",
			a:        ``,
			b:        `"hello"`,
			wantSame: false,
		},
		{
			name:     "void hashes different from empty array",
			a:        ``,
			b:        `[]`,
			wantSame: false,
		},
		{
			name:     "void hashes different from empty object",
			a:        ``,
			b:        `{}`,
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

func TestVoidDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected []string
		options  []Option
	}{
		{
			name:     "no diff between same void values",
			a:        ``,
			b:        ``,
			expected: []string{},
		},
		{
			name:     "no diff between void and whitespace",
			a:        ``,
			b:        `   `,
			expected: []string{},
		},
		{
			name: "diff void to number",
			a:    ``,
			b:    `1`,
			expected: []string{
				`@ []`,
				`+ 1`,
			},
		},
		{
			name: "diff number to void",
			a:    `1`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- 1`,
			},
		},
		{
			name: "diff void to null",
			a:    ``,
			b:    `null`,
			expected: []string{
				`@ []`,
				`+ null`,
			},
		},
		{
			name: "diff null to void",
			a:    `null`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- null`,
			},
		},
		{
			name: "diff void to boolean true",
			a:    ``,
			b:    `true`,
			expected: []string{
				`@ []`,
				`+ true`,
			},
		},
		{
			name: "diff boolean false to void",
			a:    `false`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- false`,
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
			name: "diff string to void",
			a:    `"hello"`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- "hello"`,
			},
		},
		{
			name: "diff void to array",
			a:    ``,
			b:    `[1,2,3]`,
			expected: []string{
				`@ []`,
				`+ [1,2,3]`,
			},
		},
		{
			name: "diff array to void",
			a:    `[1,2,3]`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- [1,2,3]`,
			},
		},
		{
			name: "diff void to object",
			a:    ``,
			b:    `{"key":"value"}`,
			expected: []string{
				`@ []`,
				`+ {"key":"value"}`,
			},
		},
		{
			name: "diff object to void",
			a:    `{"key":"value"}`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- {"key":"value"}`,
			},
		},
		{
			name:    "merge diff void to number",
			a:       ``,
			b:       `1`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ 1`,
			},
		},
		{
			name:    "merge diff number to void",
			a:       `1`,
			b:       ``,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:     "merge diff void to void",
			a:        ``,
			b:        ``,
			options:  []Option{MERGE},
			expected: []string{},
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

func TestVoidPatch(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		expected string
		patch    []string
	}{
		{
			name:     "no patch on void",
			initial:  ``,
			expected: ``,
			patch:    []string{},
		},
		{
			name:     "patch void to number",
			initial:  ``,
			expected: `1`,
			patch: []string{
				`@ []`,
				`+ 1`,
			},
		},
		{
			name:     "patch number to void",
			initial:  `1`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- 1`,
			},
		},
		{
			name:     "patch void to null",
			initial:  ``,
			expected: `null`,
			patch: []string{
				`@ []`,
				`+ null`,
			},
		},
		{
			name:     "patch null to void",
			initial:  `null`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- null`,
			},
		},
		{
			name:     "patch void to boolean",
			initial:  ``,
			expected: `true`,
			patch: []string{
				`@ []`,
				`+ true`,
			},
		},
		{
			name:     "patch boolean to void",
			initial:  `false`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- false`,
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
			name:     "patch string to void",
			initial:  `"hello"`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- "hello"`,
			},
		},
		{
			name:     "patch void to array",
			initial:  ``,
			expected: `[1,2,3]`,
			patch: []string{
				`@ []`,
				`+ [1,2,3]`,
			},
		},
		{
			name:     "patch array to void",
			initial:  `[1,2,3]`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- [1,2,3]`,
			},
		},
		{
			name:     "patch void to object",
			initial:  ``,
			expected: `{"key":"value"}`,
			patch: []string{
				`@ []`,
				`+ {"key":"value"}`,
			},
		},
		{
			name:     "patch object to void",
			initial:  `{"key":"value"}`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- {"key":"value"}`,
			},
		},
		{
			name:     "merge patch void to number",
			initial:  ``,
			expected: `1`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ 1`,
			},
		},
		{
			name:     "merge patch number to void",
			initial:  `1`,
			expected: ``,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:     "merge patch void to void",
			initial:  ``,
			expected: ``,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
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

func TestVoidPatchError(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		patch   []string
	}{
		{
			name:    "patch void with wrong remove value null",
			initial: ``,
			patch: []string{
				`@ []`,
				`- null`,
			},
		},
		{
			name:    "patch void with wrong remove value number",
			initial: ``,
			patch: []string{
				`@ []`,
				`- 0`,
			},
		},
		{
			name:    "patch void with wrong remove value boolean",
			initial: ``,
			patch: []string{
				`@ []`,
				`- false`,
			},
		},
		{
			name:    "patch void with wrong remove value string",
			initial: ``,
			patch: []string{
				`@ []`,
				`- ""`,
			},
		},
		{
			name:    "patch void with wrong remove value array",
			initial: ``,
			patch: []string{
				`@ []`,
				`- []`,
			},
		},
		{
			name:    "patch void with wrong remove value object",
			initial: ``,
			patch: []string{
				`@ []`,
				`- {}`,
			},
		},
		{
			name:    "patch number with add but no remove",
			initial: `1`,
			patch: []string{
				`@ []`,
				`+ 2`,
			},
		},
		{
			name:    "patch null with add but no remove",
			initial: `null`,
			patch: []string{
				`@ []`,
				`+ 1`,
			},
		},
		{
			name:    "invalid merge patch void with remove and add",
			initial: ``,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`- null`,
				`+ 1`,
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
