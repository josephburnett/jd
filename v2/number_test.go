package jd

import (
	"testing"
)

func TestNumberJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "zero integer",
			input:    `0`,
			expected: `0`,
		},
		{
			name:     "zero decimal",
			input:    `0.0`,
			expected: `0`,
		},
		{
			name:     "small decimal",
			input:    `0.01`,
			expected: `0.01`,
		},
		{
			name:     "positive integer",
			input:    `123`,
			expected: `123`,
		},
		{
			name:     "negative integer",
			input:    `-123`,
			expected: `-123`,
		},
		{
			name:     "positive decimal",
			input:    `123.456`,
			expected: `123.456`,
		},
		{
			name:     "negative decimal",
			input:    `-123.456`,
			expected: `-123.456`,
		},
		{
			name:     "scientific notation positive",
			input:    `1e5`,
			expected: `100000`,
		},
		{
			name:     "scientific notation negative",
			input:    `1e-5`,
			expected: `0.00001`,
		},
		{
			name:     "number with whitespace",
			input:    ` 42 `,
			expected: `42`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkJson(ctx, tt.input, tt.expected)
		})
	}
}

func TestNumberEqual(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		options []Option
	}{
		{
			name: "zero equals zero",
			a:    `0`,
			b:    `0`,
		},
		{
			name: "zero integer equals zero decimal",
			a:    `0`,
			b:    `0.0`,
		},
		{
			name: "same small decimal",
			a:    `0.0001`,
			b:    `0.0001`,
		},
		{
			name: "same positive integer",
			a:    `123`,
			b:    `123`,
		},
		{
			name: "same negative integer",
			a:    `-123`,
			b:    `-123`,
		},
		{
			name: "same positive decimal",
			a:    `123.456`,
			b:    `123.456`,
		},
		{
			name: "same negative decimal",
			a:    `-123.456`,
			b:    `-123.456`,
		},
		{
			name: "integer equals equivalent decimal",
			a:    `1`,
			b:    `1.0`,
		},
		{
			name: "scientific notation equals regular",
			a:    `1e2`,
			b:    `100`,
		},
		{
			name:    "precision tolerance - within range",
			a:       `1.0`,
			b:       `1.09`,
			options: []Option{Precision(0.1)},
		},
		{
			name:    "precision tolerance - exact boundary",
			a:       `42.42`,
			b:       `42.420001`,
			options: []Option{Precision(0.01)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			if len(tt.options) > 0 {
				ctx = ctx.withOptions(tt.options...)
			}
			checkEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestNumberNotEqual(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		options []Option
	}{
		{
			name: "different integers",
			a:    `0`,
			b:    `1`,
		},
		{
			name: "small precision difference",
			a:    `0`,
			b:    `0.0001`,
		},
		{
			name: "consecutive integers",
			a:    `1234`,
			b:    `1235`,
		},
		{
			name: "positive vs negative",
			a:    `123`,
			b:    `-123`,
		},
		{
			name: "integer vs decimal",
			a:    `123`,
			b:    `123.001`,
		},
		{
			name: "different scientific notation",
			a:    `1e2`,
			b:    `1e3`,
		},
		{
			name: "number not equal to null",
			a:    `0`,
			b:    `null`,
		},
		{
			name: "number not equal to false",
			a:    `0`,
			b:    `false`,
		},
		{
			name: "number not equal to true",
			a:    `1`,
			b:    `true`,
		},
		{
			name: "number not equal to string number",
			a:    `123`,
			b:    `"123"`,
		},
		{
			name: "number not equal to empty string",
			a:    `0`,
			b:    `""`,
		},
		{
			name: "number not equal to void",
			a:    `0`,
			b:    ``,
		},
		{
			name: "number not equal to empty array",
			a:    `0`,
			b:    `[]`,
		},
		{
			name: "number not equal to empty object",
			a:    `0`,
			b:    `{}`,
		},
		{
			name: "number not equal to array with same number",
			a:    `123`,
			b:    `[123]`,
		},
		{
			name: "number not equal to object with same number value",
			a:    `123`,
			b:    `{"key": 123}`,
		},
		{
			name:    "precision tolerance - outside range",
			a:       `1`,
			b:       `1.2`,
			options: []Option{Precision(0.1)},
		},
		{
			name:    "precision tolerance - large difference",
			a:       `1`,
			b:       `2`,
			options: []Option{Precision(0.1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			if len(tt.options) > 0 {
				ctx = ctx.withOptions(tt.options...)
			}
			checkNotEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestNumberHash(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		wantSame bool
	}{
		{
			name:     "same number hashes same",
			a:        `0`,
			b:        `0`,
			wantSame: true,
		},
		{
			name:     "different numbers hash different",
			a:        `0`,
			b:        `1`,
			wantSame: false,
		},
		{
			name:     "integer and equivalent decimal hash same",
			a:        `1.0`,
			b:        `1`,
			wantSame: true,
		},
		{
			name:     "different decimals hash different",
			a:        `0.1`,
			b:        `0.01`,
			wantSame: false,
		},
		{
			name:     "positive and negative hash different",
			a:        `123`,
			b:        `-123`,
			wantSame: false,
		},
		{
			name:     "scientific notation and regular hash same",
			a:        `1e2`,
			b:        `100`,
			wantSame: true,
		},
		{
			name:     "number hashes different from string",
			a:        `123`,
			b:        `"123"`,
			wantSame: false,
		},
		{
			name:     "number hashes different from boolean",
			a:        `1`,
			b:        `true`,
			wantSame: false,
		},
		{
			name:     "zero hashes different from false",
			a:        `0`,
			b:        `false`,
			wantSame: false,
		},
		{
			name:     "number hashes different from null",
			a:        `0`,
			b:        `null`,
			wantSame: false,
		},
		{
			name:     "number hashes different from void",
			a:        `123`,
			b:        ``,
			wantSame: false,
		},
		{
			name:     "number hashes different from empty array",
			a:        `0`,
			b:        `[]`,
			wantSame: false,
		},
		{
			name:     "number hashes different from empty object",
			a:        `0`,
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

func TestNumberDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected []string
		options  []Option
	}{
		{
			name:     "no diff between same numbers",
			a:        `0`,
			b:        `0`,
			expected: []string{},
		},
		{
			name:     "no diff between equivalent forms",
			a:        `1`,
			b:        `1.0`,
			expected: []string{},
		},
		{
			name: "diff different numbers",
			a:    `0`,
			b:    `1`,
			expected: []string{
				`@ []`,
				`- 0`,
				`+ 1`,
			},
		},
		{
			name: "diff number to void",
			a:    `0`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- 0`,
			},
		},
		{
			name: "diff void to number",
			a:    ``,
			b:    `42`,
			expected: []string{
				`@ []`,
				`+ 42`,
			},
		},
		{
			name: "diff number to null",
			a:    `123`,
			b:    `null`,
			expected: []string{
				`@ []`,
				`- 123`,
				`+ null`,
			},
		},
		{
			name: "diff number to boolean",
			a:    `1`,
			b:    `true`,
			expected: []string{
				`@ []`,
				`- 1`,
				`+ true`,
			},
		},
		{
			name: "diff number to string",
			a:    `123`,
			b:    `"123"`,
			expected: []string{
				`@ []`,
				`- 123`,
				`+ "123"`,
			},
		},
		{
			name: "diff number to array",
			a:    `123`,
			b:    `[123]`,
			expected: []string{
				`@ []`,
				`- 123`,
				`+ [123]`,
			},
		},
		{
			name: "diff number to object",
			a:    `123`,
			b:    `{"key": 123}`,
			expected: []string{
				`@ []`,
				`- 123`,
				`+ {"key":123}`,
			},
		},
		{
			name: "diff positive to negative",
			a:    `123`,
			b:    `-123`,
			expected: []string{
				`@ []`,
				`- 123`,
				`+ -123`,
			},
		},
		{
			name: "diff integer to decimal",
			a:    `123`,
			b:    `123.456`,
			expected: []string{
				`@ []`,
				`- 123`,
				`+ 123.456`,
			},
		},
		{
			name:    "merge diff numbers",
			a:       `1`,
			b:       `2`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ 2`,
			},
		},
		{
			name:    "merge diff number to void",
			a:       `123`,
			b:       ``,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			},
		},
		{
			name:    "precision tolerance - no diff within range",
			a:       `42.42`,
			b:       `42.420001`,
			options: []Option{Precision(0.01)},
			expected: []string{},
		},
		{
			name:    "precision tolerance - diff outside range",
			a:       `42.42`,
			b:       `42.43`,
			options: []Option{Precision(0.005)},
			expected: []string{
				`@ []`,
				`- 42.42`,
				`+ 42.43`,
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

func TestNumberPatch(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		expected string
		patch    []string
	}{
		{
			name:     "no patch on number",
			initial:  `0`,
			expected: `0`,
			patch:    []string{},
		},
		{
			name:     "patch number to different number",
			initial:  `0`,
			expected: `1`,
			patch: []string{
				`@ []`,
				`- 0`,
				`+ 1`,
			},
		},
		{
			name:     "patch number to void",
			initial:  `0`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- 0`,
			},
		},
		{
			name:     "patch void to number",
			initial:  ``,
			expected: `42`,
			patch: []string{
				`@ []`,
				`+ 42`,
			},
		},
		{
			name:     "patch number to null",
			initial:  `123`,
			expected: `null`,
			patch: []string{
				`@ []`,
				`- 123`,
				`+ null`,
			},
		},
		{
			name:     "patch number to boolean",
			initial:  `1`,
			expected: `true`,
			patch: []string{
				`@ []`,
				`- 1`,
				`+ true`,
			},
		},
		{
			name:     "patch number to string",
			initial:  `123`,
			expected: `"123"`,
			patch: []string{
				`@ []`,
				`- 123`,
				`+ "123"`,
			},
		},
		{
			name:     "patch number to array",
			initial:  `123`,
			expected: `[123]`,
			patch: []string{
				`@ []`,
				`- 123`,
				`+ [123]`,
			},
		},
		{
			name:     "patch number to object",
			initial:  `123`,
			expected: `{"key":123}`,
			patch: []string{
				`@ []`,
				`- 123`,
				`+ {"key":123}`,
			},
		},
		{
			name:     "patch positive to negative",
			initial:  `123`,
			expected: `-123`,
			patch: []string{
				`@ []`,
				`- 123`,
				`+ -123`,
			},
		},
		{
			name:     "patch integer to decimal",
			initial:  `123`,
			expected: `123.456`,
			patch: []string{
				`@ []`,
				`- 123`,
				`+ 123.456`,
			},
		},
		{
			name:     "merge patch number",
			initial:  `0`,
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
			name:     "merge patch number to string",
			initial:  `123`,
			expected: `"hello"`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`+ "hello"`,
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

func TestNumberPatchError(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		patch   []string
	}{
		{
			name:    "patch number with wrong remove value",
			initial: `0`,
			patch: []string{
				`@ []`,
				`- 1`,
			},
		},
		{
			name:    "patch void with remove number",
			initial: ``,
			patch: []string{
				`@ []`,
				`- 0`,
			},
		},
		{
			name:    "invalid merge patch number with remove and add",
			initial: `0`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ []`,
				`- 0`,
				`+ 1`,
			},
		},
		{
			name:    "patch number with wrong remove type string",
			initial: `123`,
			patch: []string{
				`@ []`,
				`- "123"`,
			},
		},
		{
			name:    "patch number with wrong remove type boolean",
			initial: `1`,
			patch: []string{
				`@ []`,
				`- true`,
			},
		},
		{
			name:    "patch number with wrong remove type null",
			initial: `0`,
			patch: []string{
				`@ []`,
				`- null`,
			},
		},
		{
			name:    "patch number with wrong remove type array",
			initial: `123`,
			patch: []string{
				`@ []`,
				`- []`,
			},
		},
		{
			name:    "patch number with wrong remove type object",
			initial: `123`,
			patch: []string{
				`@ []`,
				`- {}`,
			},
		},
		{
			name:    "patch number with close but wrong value",
			initial: `123.456`,
			patch: []string{
				`@ []`,
				`- 123.457`,
			},
		},
		{
			name:    "patch positive with negative remove",
			initial: `123`,
			patch: []string{
				`@ []`,
				`- -123`,
			},
		},
		{
			name:    "patch null with add number but no remove",
			initial: `null`,
			patch: []string{
				`@ []`,
				`+ 123`,
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
