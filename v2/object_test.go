package jd

import (
	"testing"
)

func TestObjectJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple object",
			input:    `{"a":1}`,
			expected: `{"a":1}`,
		},
		{
			name:     "object with whitespace",
			input:    ` { "a" : 1 } `,
			expected: `{"a":1}`,
		},
		{
			name:     "empty object",
			input:    `{}`,
			expected: `{}`,
		},
		{
			name:     "empty object with whitespace",
			input:    ` { } `,
			expected: `{}`,
		},
		{
			name:     "object with multiple keys",
			input:    `{"a":1,"b":2}`,
			expected: `{"a":1,"b":2}`,
		},
		{
			name:     "object with nested objects",
			input:    `{"a":{"b":1}}`,
			expected: `{"a":{"b":1}}`,
		},
		{
			name:     "object with array values",
			input:    `{"a":[1,2,3]}`,
			expected: `{"a":[1,2,3]}`,
		},
		{
			name:     "object with empty string key",
			input:    `{"":1}`,
			expected: `{"":1}`,
		},
		{
			name:     "object with mixed value types",
			input:    `{"a":1,"b":"hello","c":true,"d":null}`,
			expected: `{"a":1,"b":"hello","c":true,"d":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkJson(ctx, tt.input, tt.expected)
		})
	}
}

func TestObjectEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "same simple objects",
			a:    `{"a":1}`,
			b:    `{"a":1}`,
		},
		{
			name: "integer equals decimal in object",
			a:    `{"a":1}`,
			b:    `{"a":1.0}`,
		},
		{
			name: "objects with same array values",
			a:    `{"a":[1,2]}`,
			b:    `{"a":[1,2]}`,
		},
		{
			name: "objects with same string values",
			a:    `{"a":"b"}`,
			b:    `{"a":"b"}`,
		},
		{
			name: "empty objects equal",
			a:    `{}`,
			b:    `{}`,
		},
		{
			name: "objects with multiple same keys",
			a:    `{"a":1,"b":2}`,
			b:    `{"a":1,"b":2}`,
		},
		{
			name: "objects with keys in different order",
			a:    `{"a":1,"b":2}`,
			b:    `{"b":2,"a":1}`,
		},
		{
			name: "nested objects equal",
			a:    `{"a":{"b":1}}`,
			b:    `{"a":{"b":1}}`,
		},
		{
			name: "objects with null values equal",
			a:    `{"a":null}`,
			b:    `{"a":null}`,
		},
		{
			name: "objects with boolean values equal",
			a:    `{"a":true,"b":false}`,
			b:    `{"a":true,"b":false}`,
		},
		{
			name: "objects with empty string key equal",
			a:    `{"":1}`,
			b:    `{"":1}`,
		},
		{
			name: "objects with mixed types equal",
			a:    `{"a":1,"b":"hello","c":true,"d":null}`,
			b:    `{"a":1,"b":"hello","c":true,"d":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestObjectNotEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "different keys same values",
			a:    `{"a":1}`,
			b:    `{"b":1}`,
		},
		{
			name: "same key different array order",
			a:    `{"a":[1,2]}`,
			b:    `{"a":[2,1]}`,
		},
		{
			name: "same key different string values",
			a:    `{"a":"b"}`,
			b:    `{"a":"c"}`,
		},
		{
			name: "empty object not equal to non-empty",
			a:    `{}`,
			b:    `{"a":1}`,
		},
		{
			name: "different number of keys",
			a:    `{"a":1}`,
			b:    `{"a":1,"b":2}`,
		},
		{
			name: "same key different value types",
			a:    `{"a":1}`,
			b:    `{"a":"1"}`,
		},
		{
			name: "nested vs flat object",
			a:    `{"a":{"b":1}}`,
			b:    `{"a.b":1}`,
		},
		{
			name: "object not equal to null",
			a:    `{"a":1}`,
			b:    `null`,
		},
		{
			name: "object not equal to array",
			a:    `{"a":1}`,
			b:    `[1]`,
		},
		{
			name: "object not equal to string",
			a:    `{"a":1}`,
			b:    `"object"`,
		},
		{
			name: "object not equal to number",
			a:    `{"a":1}`,
			b:    `1`,
		},
		{
			name: "object not equal to boolean",
			a:    `{"a":1}`,
			b:    `true`,
		},
		{
			name: "object not equal to void",
			a:    `{"a":1}`,
			b:    ``,
		},
		{
			name: "empty object not equal to empty array",
			a:    `{}`,
			b:    `[]`,
		},
		{
			name: "different null vs missing key",
			a:    `{"a":null}`,
			b:    `{}`,
		},
		{
			name: "case sensitive keys",
			a:    `{"a":1}`,
			b:    `{"A":1}`,
		},
		{
			name: "whitespace in key names",
			a:    `{"a":1}`,
			b:    `{" a ":1}`,
		},
		{
			name: "different nested structure",
			a:    `{"a":{"b":{"c":1}}}`,
			b:    `{"a":{"b":1}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t)
			checkNotEqual(ctx, tt.a, tt.b)
		})
	}
}

func TestObjectHash(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		wantSame bool
	}{
		{
			name:     "empty objects hash same",
			a:        `{}`,
			b:        `{}`,
			wantSame: true,
		},
		{
			name:     "identical simple objects hash same",
			a:        `{"a":1}`,
			b:        `{"a":1}`,
			wantSame: true,
		},
		{
			name:     "same keys different values hash different",
			a:        `{"a":1}`,
			b:        `{"a":2}`,
			wantSame: false,
		},
		{
			name:     "different keys same values hash different",
			a:        `{"a":1}`,
			b:        `{"b":1}`,
			wantSame: false,
		},
		{
			name:     "same content different key order hash same",
			a:        `{"a":1,"b":2}`,
			b:        `{"b":2,"a":1}`,
			wantSame: true,
		},
		{
			name:     "nested objects hash same when identical",
			a:        `{"a":{"b":1}}`,
			b:        `{"a":{"b":1}}`,
			wantSame: true,
		},
		{
			name:     "nested objects hash different when different",
			a:        `{"a":{"b":1}}`,
			b:        `{"a":{"b":2}}`,
			wantSame: false,
		},
		{
			name:     "objects with arrays hash same when identical",
			a:        `{"a":[1,2,3]}`,
			b:        `{"a":[1,2,3]}`,
			wantSame: true,
		},
		{
			name:     "objects with arrays hash different when different",
			a:        `{"a":[1,2,3]}`,
			b:        `{"a":[1,3,2]}`,
			wantSame: false,
		},
		{
			name:     "object hashes different from null",
			a:        `{"a":1}`,
			b:        `null`,
			wantSame: false,
		},
		{
			name:     "object hashes different from array",
			a:        `{"a":1}`,
			b:        `[1]`,
			wantSame: false,
		},
		{
			name:     "object hashes different from string",
			a:        `{"a":1}`,
			b:        `"object"`,
			wantSame: false,
		},
		{
			name:     "object hashes different from number",
			a:        `{"a":1}`,
			b:        `1`,
			wantSame: false,
		},
		{
			name:     "object hashes different from boolean",
			a:        `{"a":1}`,
			b:        `true`,
			wantSame: false,
		},
		{
			name:     "object hashes different from void",
			a:        `{"a":1}`,
			b:        ``,
			wantSame: false,
		},
		{
			name:     "empty object hashes different from empty array",
			a:        `{}`,
			b:        `[]`,
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

func TestObjectDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected []string
		options  []Option
	}{
		{
			name:     "no diff between same empty objects",
			a:        `{}`,
			b:        `{}`,
			expected: []string{},
		},
		{
			name:     "no diff between same simple objects",
			a:        `{"a":1}`,
			b:        `{"a":1}`,
			expected: []string{},
		},
		{
			name: "diff object key value change",
			a:    `{"a":1}`,
			b:    `{"a":2}`,
			expected: []string{
				`@ ["a"]`,
				`- 1`,
				`+ 2`,
			},
		},
		{
			name:     "no diff between same empty string key objects",
			a:        `{"":1}`,
			b:        `{"":1}`,
			expected: []string{},
		},
		{
			name: "diff empty string key to named key",
			a:    `{"":1}`,
			b:    `{"a":2}`,
			expected: []string{
				`@ [""]`,
				`- 1`,
				`@ ["a"]`,
				`+ 2`,
			},
		},
		{
			name: "diff nested object expansion",
			a:    `{"a":{"b":{}}}`,
			b:    `{"a":{"b":{"c":1},"d":2}}`,
			expected: []string{
				`@ ["a","b","c"]`,
				`+ 1`,
				`@ ["a","d"]`,
				`+ 2`,
			},
		},
		{
			name: "regression test for issue #18 - array insertion",
			a:    `{"R": [{"I": [{"T": [{"V": "t","K": "N"},{"V": "T","K": "I"}]}]}]}`,
			b:    `{"R": [{"I": [{"T": [{"V": "t","K": "N"},{"V": "Q","K": "C"},{"V": "T","K": "I"}]}]}]}`,
			expected: []string{
				`@ ["R",0,"I",0,"T",1]`,
				`  {"K":"N","V":"t"}`,
				`+ {"K":"C","V":"Q"}`,
				`  {"K":"I","V":"T"}`,
			},
		},
		{
			name: "diff object to null",
			a:    `{"a":1}`,
			b:    `null`,
			expected: []string{
				`@ []`,
				`- {"a":1}`,
				`+ null`,
			},
		},
		{
			name: "diff object to array",
			a:    `{"a":1}`,
			b:    `[1]`,
			expected: []string{
				`@ []`,
				`- {"a":1}`,
				`+ [1]`,
			},
		},
		{
			name: "diff object to string",
			a:    `{"a":1}`,
			b:    `"object"`,
			expected: []string{
				`@ []`,
				`- {"a":1}`,
				`+ "object"`,
			},
		},
		{
			name: "diff object to number",
			a:    `{"a":1}`,
			b:    `42`,
			expected: []string{
				`@ []`,
				`- {"a":1}`,
				`+ 42`,
			},
		},
		{
			name: "diff object to boolean",
			a:    `{"a":1}`,
			b:    `true`,
			expected: []string{
				`@ []`,
				`- {"a":1}`,
				`+ true`,
			},
		},
		{
			name: "diff object to void",
			a:    `{"a":1}`,
			b:    ``,
			expected: []string{
				`@ []`,
				`- {"a":1}`,
			},
		},
		{
			name:    "merge diff object key value change",
			a:       `{"a":1}`,
			b:       `{"a":2}`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ ["a"]`,
				`+ 2`,
			},
		},
		{
			name:    "merge diff object value to null",
			a:       `{"a":1}`,
			b:       `{"a":null}`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ ["a"]`,
				`+ null`,
			},
		},
		{
			name:    "merge diff remove key",
			a:       `{"a":1}`,
			b:       `{}`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ ["a"]`,
				`+`,
			},
		},
		{
			name:    "merge diff change key",
			a:       `{"a":1}`,
			b:       `{"b":1}`,
			options: []Option{MERGE},
			expected: []string{
				`^ {"Merge":true}`,
				`@ ["a"]`,
				`+`,
				`^ {"Merge":true}`,
				`@ ["b"]`,
				`+ 1`,
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

func TestObjectPatch(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		expected string
		patch    []string
	}{
		{
			name:     "no patch on empty object",
			initial:  `{}`,
			expected: `{}`,
			patch:    []string{},
		},
		{
			name:     "no patch on same object",
			initial:  `{"a":1}`,
			expected: `{"a":1}`,
			patch:    []string{},
		},
		{
			name:     "patch object key value change",
			initial:  `{"a":1}`,
			expected: `{"a":2}`,
			patch: []string{
				`@ ["a"]`,
				`- 1`,
				`+ 2`,
			},
		},
		{
			name:     "no patch on empty string key object",
			initial:  `{"":1}`,
			expected: `{"":1}`,
			patch:    []string{},
		},
		{
			name:     "patch empty string key to named key",
			initial:  `{"":1}`,
			expected: `{"a":2}`,
			patch: []string{
				`@ [""]`,
				`- 1`,
				`@ ["a"]`,
				`+ 2`,
			},
		},
		{
			name:     "patch nested object expansion",
			initial:  `{"a":{"b":{}}}`,
			expected: `{"a":{"b":{"c":1},"d":2}}`,
			patch: []string{
				`@ ["a","b","c"]`,
				`+ 1`,
				`@ ["a","d"]`,
				`+ 2`,
			},
		},
		{
			name:     "patch object to null",
			initial:  `{"a":1}`,
			expected: `null`,
			patch: []string{
				`@ []`,
				`- {"a":1}`,
				`+ null`,
			},
		},
		{
			name:     "patch object to array",
			initial:  `{"a":1}`,
			expected: `[1]`,
			patch: []string{
				`@ []`,
				`- {"a":1}`,
				`+ [1]`,
			},
		},
		{
			name:     "patch object to string",
			initial:  `{"a":1}`,
			expected: `"object"`,
			patch: []string{
				`@ []`,
				`- {"a":1}`,
				`+ "object"`,
			},
		},
		{
			name:     "patch object to number",
			initial:  `{"a":1}`,
			expected: `42`,
			patch: []string{
				`@ []`,
				`- {"a":1}`,
				`+ 42`,
			},
		},
		{
			name:     "patch object to boolean",
			initial:  `{"a":1}`,
			expected: `true`,
			patch: []string{
				`@ []`,
				`- {"a":1}`,
				`+ true`,
			},
		},
		{
			name:     "patch object to void",
			initial:  `{"a":1}`,
			expected: ``,
			patch: []string{
				`@ []`,
				`- {"a":1}`,
			},
		},
		{
			name:     "merge patch object key value",
			initial:  `{"foo":1}`,
			expected: `{"foo":2}`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ ["foo"]`,
				`+ 2`,
			},
		},
		{
			name:     "merge patch object array replacement",
			initial:  `{"foo":[1,2,3]}`,
			expected: `{"foo":[4,5,6]}`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ ["foo"]`,
				`+ [4,5,6]`,
			},
		},
		{
			name:     "merge patch nested object creation",
			initial:  `{}`,
			expected: `{"foo":{"bar":1}}`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ ["foo","bar"]`,
				`+ 1`,
			},
		},
		{
			name:     "patch add new key",
			initial:  `{"a":1}`,
			expected: `{"a":1,"b":2}`,
			patch: []string{
				`@ ["b"]`,
				`+ 2`,
			},
		},
		{
			name:     "patch remove key",
			initial:  `{"a":1,"b":2}`,
			expected: `{"a":1}`,
			patch: []string{
				`@ ["b"]`,
				`- 2`,
			},
		},
		{
			name:     "patch nested key change",
			initial:  `{"a":{"b":1}}`,
			expected: `{"a":{"b":2}}`,
			patch: []string{
				`@ ["a","b"]`,
				`- 1`,
				`+ 2`,
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

func TestObjectPatchMerge(t *testing.T) {
	// Merge strategy at base case
	a, _ := ReadJsonString(`{"a":1,"b":2}`)
	b, _ := ReadJsonString(`{"a":1,"b":3}`)
	d := a.Diff(b, MERGE)
	result, err := a.Patch(d)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Equals(b) {
		t.Errorf("merge patch failed: got %v", result.Json())
	}
}

func TestObjectPatchError(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		patch   []string
	}{
		{
			name:    "patch empty object with remove non-existent key",
			initial: `{}`,
			patch: []string{
				`@ ["a"]`,
				`- 1`,
			},
		},
		{
			name:    "patch object with add to existing key",
			initial: `{"a":1}`,
			patch: []string{
				`@ ["a"]`,
				`+ 2`,
			},
		},
		{
			name:    "patch object with duplicate add",
			initial: `{"a":1}`,
			patch: []string{
				`@ ["a"]`,
				`+ 1`,
			},
		},
		{
			name:    "patch object with wrong remove value",
			initial: `{"a":1}`,
			patch: []string{
				`@ ["a"]`,
				`- 2`,
			},
		},
		{
			name:    "patch object with wrong remove type",
			initial: `{"a":1}`,
			patch: []string{
				`@ ["a"]`,
				`- "1"`,
			},
		},
		{
			name:    "patch object with remove non-existent nested key",
			initial: `{"a":{}}`,
			patch: []string{
				`@ ["a","b"]`,
				`- 1`,
			},
		},
		{
			name:    "patch object with wrong nested value",
			initial: `{"a":{"b":1}}`,
			patch: []string{
				`@ ["a","b"]`,
				`- 2`,
			},
		},
		{
			name:    "patch object with add to non-existent path",
			initial: `{}`,
			patch: []string{
				`@ ["a","b"]`,
				`+ 1`,
			},
		},
		{
			name:    "patch null with object patch",
			initial: `null`,
			patch: []string{
				`@ ["a"]`,
				`+ 1`,
			},
		},
		{
			name:    "patch array with object patch",
			initial: `[1,2,3]`,
			patch: []string{
				`@ ["a"]`,
				`+ 1`,
			},
		},
		{
			name:    "invalid merge patch object with remove and add",
			initial: `{"a":1}`,
			patch: []string{
				`^ {"Merge":true}`,
				`@ ["a"]`,
				`- 1`,
				`+ 2`,
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
