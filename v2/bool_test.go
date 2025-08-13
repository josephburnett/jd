package jd

import (
	"testing"
)

func TestBoolJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "true boolean",
			input:    `true`,
			expected: `true`,
		},
		{
			name:     "false boolean",
			input:    `false`,
			expected: `false`,
		},
		{
			name:     "true with whitespace",
			input:    ` true `,
			expected: `true`,
		},
		{
			name:     "false with whitespace",
			input:    ` false `,
			expected: `false`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkJson(ctx, tt.input, tt.expected)
		})
	}
}

func TestBoolEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "true equals true",
			a:    `true`,
			b:    `true`,
		},
		{
			name: "false equals false",
			a:    `false`,
			b:    `false`,
		},
		{
			name: "true with whitespace equals true",
			a:    ` true `,
			b:    `true`,
		},
		{
			name: "false with whitespace equals false",
			a:    ` false `,
			b:    `false`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestBoolNotEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "true not equal to false",
			a:    `true`,
			b:    `false`,
		},
		{
			name: "false not equal to true",
			a:    `false`,
			b:    `true`,
		},
		{
			name: "false not equal to empty array",
			a:    `false`,
			b:    `[]`,
		},
		{
			name: "true not equal to string true",
			a:    `true`,
			b:    `"true"`,
		},
		{
			name: "false not equal to string false",
			a:    `false`,
			b:    `"false"`,
		},
		{
			name: "true not equal to number 1",
			a:    `true`,
			b:    `1`,
		},
		{
			name: "false not equal to number 0",
			a:    `false`,
			b:    `0`,
		},
		{
			name: "true not equal to null",
			a:    `true`,
			b:    `null`,
		},
		{
			name: "false not equal to null",
			a:    `false`,
			b:    `null`,
		},
		{
			name: "true not equal to empty object",
			a:    `true`,
			b:    `{}`,
		},
		{
			name: "false not equal to empty object",
			a:    `false`,
			b:    `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkNotEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestBoolHash(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		wantSame bool
	}{
		{
			name:     "true hashes same as true",
			a:        `true`,
			b:        `true`,
			wantSame: true,
		},
		{
			name:     "false hashes same as false",
			a:        `false`,
			b:        `false`,
			wantSame: true,
		},
		{
			name:     "true hashes different from false",
			a:        `true`,
			b:        `false`,
			wantSame: false,
		},
		{
			name:     "false hashes different from true",
			a:        `false`,
			b:        `true`,
			wantSame: false,
		},
		{
			name:     "true hashes different from string true",
			a:        `true`,
			b:        `"true"`,
			wantSame: false,
		},
		{
			name:     "false hashes different from string false",
			a:        `false`,
			b:        `"false"`,
			wantSame: false,
		},
		{
			name:     "true hashes different from number 1",
			a:        `true`,
			b:        `1`,
			wantSame: false,
		},
		{
			name:     "false hashes different from number 0",
			a:        `false`,
			b:        `0`,
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

func TestBoolDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected []string
		options  []Option
	}{
		{
			name:     "no diff between same true values",
			a:        `true`,
			b:        `true`,
			expected: []string{},
		},
		{
			name:     "no diff between same false values",
			a:        `false`,
			b:        `false`,
			expected: []string{},
		},
		{
			name: "diff true to false",
			a:    `true`,
			b:    `false`,
			expected: []string{
				`@ []`,
				`- true`,
				`+ false`,
			},
		},
		{
			name: "diff false to true",
			a:    `false`,
			b:    `true`,
			expected: []string{
				`@ []`,
				`- false`,
				`+ true`,
			},
		},
		{
			name:    "diff true to false with merge option",
			a:       `true`,
			b:       `false`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ false`,
			},
		},
		{
			name:    "diff false to true with merge option",
			a:       `false`,
			b:       `true`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ true`,
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

func TestBoolPatch(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		expected string
		patch    []string
	}{
		{
			name:     "no patch on true",
			initial:  `true`,
			expected: `true`,
			patch:    []string{},
		},
		{
			name:     "no patch on false",
			initial:  `false`,
			expected: `false`,
			patch:    []string{},
		},
		{
			name:     "patch true to false",
			initial:  `true`,
			expected: `false`,
			patch: []string{
				`@ []`,
				`- true`,
				`+ false`,
			},
		},
		{
			name:     "patch false to true",
			initial:  `false`,
			expected: `true`,
			patch: []string{
				`@ []`,
				`- false`,
				`+ true`,
			},
		},
		{
			name:     "merge patch false to true",
			initial:  `false`,
			expected: `true`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ true`,
			},
		},
		{
			name:     "merge patch true to false",
			initial:  `true`,
			expected: `false`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ false`,
			},
		},
		{
			name:     "merge patch true to void",
			initial:  `true`,
			expected: ``,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:     "merge patch false to void",
			initial:  `false`,
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

func TestBoolPatchError(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		patch   []string
	}{
		{
			name:    "patch true with wrong remove value false",
			initial: `true`,
			patch: []string{
				`@ []`,
				`- false`,
			},
		},
		{
			name:    "patch false with wrong remove value true",
			initial: `false`,
			patch: []string{
				`@ []`,
				`- true`,
			},
		},
		{
			name:    "invalid merge patch true with remove and add",
			initial: `true`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`- true`,
				`+ false`,
			},
		},
		{
			name:    "invalid merge patch false with remove and add",
			initial: `false`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`- false`,
				`+ true`,
			},
		},
		{
			name:    "patch true with wrong value null",
			initial: `true`,
			patch: []string{
				`@ []`,
				`- null`,
			},
		},
		{
			name:    "patch false with wrong value number",
			initial: `false`,
			patch: []string{
				`@ []`,
				`- 0`,
			},
		},
		{
			name:    "patch true with wrong value string",
			initial: `true`,
			patch: []string{
				`@ []`,
				`- "true"`,
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
