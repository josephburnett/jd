package jd

import (
	"testing"
)

func TestNullJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "null value",
			input:    `null`,
			expected: `null`,
		},
		{
			name:     "null with whitespace",
			input:    ` null `,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkJson(ctx, tt.input, tt.expected)
		})
	}
}

func TestNullEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "null equals null",
			a:    `null`,
			b:    `null`,
		},
		{
			name: "null with whitespace equals null",
			a:    ` null `,
			b:    `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestNullNotEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "null not equal to number 0",
			a:    `null`,
			b:    `0`,
		},
		{
			name: "null not equal to empty array",
			a:    `null`,
			b:    `[]`,
		},
		{
			name: "null not equal to empty object",
			a:    `null`,
			b:    `{}`,
		},
		{
			name: "null not equal to false",
			a:    `null`,
			b:    `false`,
		},
		{
			name: "null not equal to true",
			a:    `null`,
			b:    `true`,
		},
		{
			name: "null not equal to string null",
			a:    `null`,
			b:    `"null"`,
		},
		{
			name: "null not equal to empty string",
			a:    `null`,
			b:    `""`,
		},
		{
			name: "null not equal to void",
			a:    `null`,
			b:    ``,
		},
		{
			name: "null not equal to negative number",
			a:    `null`,
			b:    `-1`,
		},
		{
			name: "null not equal to positive number",
			a:    `null`,
			b:    `1`,
		},
		{
			name: "null not equal to decimal number",
			a:    `null`,
			b:    `0.5`,
		},
		{
			name: "null not equal to array with null",
			a:    `null`,
			b:    `[null]`,
		},
		{
			name: "null not equal to object with null value",
			a:    `null`,
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

func TestNullHash(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		wantSame bool
	}{
		{
			name:     "null hashes same as null",
			a:        `null`,
			b:        `null`,
			wantSame: true,
		},
		{
			name:     "null hashes different from void",
			a:        `null`,
			b:        ``,
			wantSame: false,
		},
		{
			name:     "null hashes different from false",
			a:        `null`,
			b:        `false`,
			wantSame: false,
		},
		{
			name:     "null hashes different from true",
			a:        `null`,
			b:        `true`,
			wantSame: false,
		},
		{
			name:     "null hashes different from number 0",
			a:        `null`,
			b:        `0`,
			wantSame: false,
		},
		{
			name:     "null hashes different from string null",
			a:        `null`,
			b:        `"null"`,
			wantSame: false,
		},
		{
			name:     "null hashes different from empty array",
			a:        `null`,
			b:        `[]`,
			wantSame: false,
		},
		{
			name:     "null hashes different from empty object",
			a:        `null`,
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

func TestNullDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected []string
		options  []Option
	}{
		{
			name:     "no diff between same null values",
			a:        `null`,
			b:        `null`,
			expected: []string{},
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
			name: "diff void to null",
			a:    ``,
			b:    `null`,
			expected: []string{
				`@ []`,
				`+ null`,
			},
		},
		{
			name: "diff null to true",
			a:    `null`,
			b:    `true`,
			expected: []string{
				`@ []`,
				`- null`,
				`+ true`,
			},
		},
		{
			name: "diff null to false",
			a:    `null`,
			b:    `false`,
			expected: []string{
				`@ []`,
				`- null`,
				`+ false`,
			},
		},
		{
			name: "diff null to number",
			a:    `null`,
			b:    `42`,
			expected: []string{
				`@ []`,
				`- null`,
				`+ 42`,
			},
		},
		{
			name: "diff null to string",
			a:    `null`,
			b:    `"hello"`,
			expected: []string{
				`@ []`,
				`- null`,
				`+ "hello"`,
			},
		},
		{
			name: "diff null to array",
			a:    `null`,
			b:    `[1,2,3]`,
			expected: []string{
				`@ []`,
				`- null`,
				`+ [1,2,3]`,
			},
		},
		{
			name: "diff null to object",
			a:    `null`,
			b:    `{"key": "value"}`,
			expected: []string{
				`@ []`,
				`- null`,
				`+ {"key":"value"}`,
			},
		},
		{
			name:    "merge diff true to null",
			a:       `true`,
			b:       `null`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ null`,
			},
		},
		{
			name:    "merge diff null to true",
			a:       `null`,
			b:       `true`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ true`,
			},
		},
		{
			name:    "merge diff null to void",
			a:       `null`,
			b:       ``,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:    "merge diff null to number",
			a:       `null`,
			b:       `0`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ 0`,
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

func TestNullPatch(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		expected string
		patch    []string
	}{
		{
			name:     "no patch on null",
			initial:  `null`,
			expected: `null`,
			patch:    []string{},
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
			name:     "patch void to null",
			initial:  ``,
			expected: `null`,
			patch: []string{
				`@ []`,
				`+ null`,
			},
		},
		{
			name:     "patch null to true",
			initial:  `null`,
			expected: `true`,
			patch: []string{
				`@ []`,
				`- null`,
				`+ true`,
			},
		},
		{
			name:     "patch null to false",
			initial:  `null`,
			expected: `false`,
			patch: []string{
				`@ []`,
				`- null`,
				`+ false`,
			},
		},
		{
			name:     "patch null to number",
			initial:  `null`,
			expected: `42`,
			patch: []string{
				`@ []`,
				`- null`,
				`+ 42`,
			},
		},
		{
			name:     "patch null to string",
			initial:  `null`,
			expected: `"hello"`,
			patch: []string{
				`@ []`,
				`- null`,
				`+ "hello"`,
			},
		},
		{
			name:     "patch null to array",
			initial:  `null`,
			expected: `[1,2,3]`,
			patch: []string{
				`@ []`,
				`- null`,
				`+ [1,2,3]`,
			},
		},
		{
			name:     "patch null to object",
			initial:  `null`,
			expected: `{"key":"value"}`,
			patch: []string{
				`@ []`,
				`- null`,
				`+ {"key":"value"}`,
			},
		},
		{
			name:     "merge patch null to void",
			initial:  `null`,
			expected: ``,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:     "merge patch null to true",
			initial:  `null`,
			expected: `true`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ true`,
			},
		},
		{
			name:     "merge patch null to number",
			initial:  `null`,
			expected: `0`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ 0`,
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

func TestNullPatchError(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		patch   []string
	}{
		{
			name:    "patch null with wrong remove value 0",
			initial: `null`,
			patch: []string{
				`@ []`,
				`- 0`,
			},
		},
		{
			name:    "patch null with wrong remove value false",
			initial: `null`,
			patch: []string{
				`@ []`,
				`- false`,
			},
		},
		{
			name:    "patch null with wrong remove value true",
			initial: `null`,
			patch: []string{
				`@ []`,
				`- true`,
			},
		},
		{
			name:    "patch null with wrong remove value empty string",
			initial: `null`,
			patch: []string{
				`@ []`,
				`- ""`,
			},
		},
		{
			name:    "patch null with wrong remove value string null",
			initial: `null`,
			patch: []string{
				`@ []`,
				`- "null"`,
			},
		},
		{
			name:    "patch null with wrong remove value empty array",
			initial: `null`,
			patch: []string{
				`@ []`,
				`- []`,
			},
		},
		{
			name:    "patch null with wrong remove value empty object",
			initial: `null`,
			patch: []string{
				`@ []`,
				`- {}`,
			},
		},
		{
			name:    "invalid merge patch null with remove and add",
			initial: `null`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`- null`,
				`+ true`,
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
