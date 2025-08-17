package jd

import (
	"strings"
	"testing"
)

func TestDiffRender(t *testing.T) {
	tests := []struct {
		name      string
		a         string
		b         string
		diffLines []string
	}{
		{
			name: "simple object value change",
			a:    `{"a":1}`,
			b:    `{"a":2}`,
			diffLines: []string{
				`@ ["a"]`,
				`- 1`,
				`+ 2`,
			},
		},
		{
			name: "nested object value change",
			a:    `{"a":{"b":1}}`,
			b:    `{"a":{"b":2}}`,
			diffLines: []string{
				`@ ["a","b"]`,
				`- 1`,
				`+ 2`,
			},
		},
		{
			name: "nested object key change",
			a:    `{"a":{"b":1}}`,
			b:    `{"a":{"c":2}}`,
			diffLines: []string{
				`@ ["a","b"]`,
				`- 1`,
				`@ ["a","c"]`,
				`+ 2`,
			},
		},
		{
			name: "parent key change with same nested object",
			a:    `{"a":{"b":1}}`,
			b:    `{"c":{"b":1}}`,
			diffLines: []string{
				`@ ["a"]`,
				`- {"b":1}`,
				`@ ["c"]`,
				`+ {"b":1}`,
			},
		},
		{
			name: "string value change",
			a:    `{"a":"bar"}`,
			b:    `{"a":"baz"}`,
			diffLines: []string{
				`@ ["a"]`,
				`- "bar"`,
				`+ "baz"`,
			},
		},
		{
			name: "array string element change",
			a:    `{"qux":["foobar","foobaz"]}`,
			b:    `{"qux":["fooarrr","foobaz"]}`,
			diffLines: []string{
				`@ ["qux",0]`,
				`[`,
				`- "foobar"`,
				`+ "fooarrr"`,
				`  "foobaz"`,
			},
		},
		{
			name: "string addition from empty",
			a:    `{"str":""}`,
			b:    `{"str":"abc"}`,
			diffLines: []string{
				`@ ["str"]`,
				`- ""`,
				`+ "abc"`,
			},
		},
		{
			name: "string removal to empty",
			a:    `{"str":"abc"}`,
			b:    `{"str":""}`,
			diffLines: []string{
				`@ ["str"]`,
				`- "abc"`,
				`+ ""`,
			},
		},
		{
			name: "nested string change",
			a:    `{"a":{"b":"hello"}}`,
			b:    `{"a":{"b":"world"}}`,
			diffLines: []string{
				`@ ["a","b"]`,
				`- "hello"`,
				`+ "world"`,
			},
		},
		{
			name: "multiple string changes",
			a:    `{"a":"foo","b":"bar"}`,
			b:    `{"a":"baz","b":"qux"}`,
			diffLines: []string{
				`@ ["a"]`,
				`- "foo"`,
				`+ "baz"`,
				`@ ["b"]`,
				`- "bar"`,
				`+ "qux"`,
			},
		},
		{
			name: "key change with same value",
			a:    `{"a":"foo"}`,
			b:    `{"b":"foo"}`,
			diffLines: []string{
				`@ ["a"]`,
				`- "foo"`,
				`@ ["b"]`,
				`+ "foo"`,
			},
		},
		{
			name: "unicode string diff",
			a:    `{"a":"こんにちは"}`,
			b:    `{"a":"さようなら"}`,
			diffLines: []string{
				`@ ["a"]`,
				`- "こんにちは"`,
				`+ "さようなら"`,
			},
		},
		{
			name: "object to null",
			a:    `{"a":1}`,
			b:    `null`,
			diffLines: []string{
				`@ []`,
				`- {"a":1}`,
				`+ null`,
			},
		},
		{
			name: "null to object",
			a:    `null`,
			b:    `{"a":1}`,
			diffLines: []string{
				`@ []`,
				`- null`,
				`+ {"a":1}`,
			},
		},
		{
			name: "object to array",
			a:    `{"a":1}`,
			b:    `[1]`,
			diffLines: []string{
				`@ []`,
				`- {"a":1}`,
				`+ [1]`,
			},
		},
		{
			name: "array to object",
			a:    `[1]`,
			b:    `{"a":1}`,
			diffLines: []string{
				`@ []`,
				`- [1]`,
				`+ {"a":1}`,
			},
		},
		{
			name: "boolean to number",
			a:    `true`,
			b:    `1`,
			diffLines: []string{
				`@ []`,
				`- true`,
				`+ 1`,
			},
		},
		{
			name: "string to boolean",
			a:    `"true"`,
			b:    `true`,
			diffLines: []string{
				`@ []`,
				`- "true"`,
				`+ true`,
			},
		},
		{
			name: "void to value",
			a:    ``,
			b:    `42`,
			diffLines: []string{
				`@ []`,
				`+ 42`,
			},
		},
		{
			name: "value to void",
			a:    `42`,
			b:    ``,
			diffLines: []string{
				`@ []`,
				`- 42`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkDiffRender(t, tt.a, tt.b, tt.diffLines...)
		})
	}
}

func checkDiffRender(t *testing.T, a, b string, diffLines ...string) {
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	aJson, err := ReadJsonString(a)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	bJson, err := ReadJsonString(b)
	if err != nil {
		t.Errorf("%v", err.Error())
	}

	// Test without color
	d := aJson.diff(bJson, nil, newOptions([]Option{}), strictPatchStrategy).Render()
	if d != diff {
		t.Errorf("%v.diff(%v) = %v. Want %v.", a, b, d, diff)
	}

	// Test with color
	coloredDiff := aJson.diff(bJson, nil, newOptions([]Option{}), strictPatchStrategy).Render(COLOR)
	strippedDiff := stripAnsiCodes(coloredDiff)
	if strippedDiff != diff {
		t.Errorf("%v.diff(%v) with color (stripped) = %v. Want %v.", a, b, strippedDiff, diff)
	}

	// Verify that uncolored parts in string diffs match between + and - lines
	lines := strings.Split(coloredDiff, "\n")
	var minusLine, plusLine string
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '-' && strings.Contains(line, "\"") { // Only check string diffs
			minusLine = line
			if i+1 < len(lines) && len(lines[i+1]) > 0 && lines[i+1][0] == '+' {
				plusLine = lines[i+1]
				minusUncolored := removeColoredParts(minusLine[1:]) // Skip the "- " prefix
				plusUncolored := removeColoredParts(plusLine[1:])   // Skip the "+ " prefix
				if minusUncolored != plusUncolored {
					t.Errorf("Uncolored parts don't match:\n- %s\n+ %s", minusUncolored, plusUncolored)
				}
			}
		}
	}
}

// removeColoredParts returns the string with the colored parts (including the text between color codes) removed
func removeColoredParts(s string) string {
	var result strings.Builder
	inColor := false
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		// detect a color code (starts coloring)
		if !inColor && i+1 < len(runes) && runes[i] == '\033' && runes[i+1] == '[' {
			inColor = true
			i++ // skip '['
			continue
		}
		// if not colored, add the character to the result
		if !inColor {
			result.WriteRune(runes[i])
		}
		// detect the reset color code (ends coloring)
		if inColor && i+2 < len(runes) && runes[i] == '[' && runes[i+1] == '0' && runes[i+2] == 'm' {
			inColor = false
			i += 2
		}
	}
	return result.String()
}

func TestDiffRenderPatch(t *testing.T) {
	tests := []struct {
		name    string
		diff    string
		patch   string
		wantErr bool
	}{
		{
			name: "simple add operation",
			diff: s(`@ ["foo"]`,
				`+ 1`),
			patch: s(`[{"op":"add","path":"/foo","value":1}]`),
		},
		{
			name: "simple remove operation",
			diff: s(`@ ["foo"]`,
				`- 1`),
			patch: s(`[`,
				`{"op":"test","path":"/foo","value":1},`,
				`{"op":"remove","path":"/foo","value":1}`,
				`]`),
		},
		{
			name: "replace operation",
			diff: s(`@ ["foo"]`,
				`- 1`,
				`+ 2`),
			patch: s(`[`,
				`{"op":"test","path":"/foo","value":1},`,
				`{"op":"remove","path":"/foo","value":1},`,
				`{"op":"add","path":"/foo","value":2}`,
				`]`),
		},
		{
			name: "complex array operations",
			diff: s(`@ [0]`,
				`[`,
				`- {}`,
				`+ 0`,
				`  []`,
				`@ [2]`,
				`  []`,
				`- 0`),
			patch: s(`[`,
				`{"op":"test","path":"/1","value":[]},`,
				`{"op":"test","path":"/0","value":{}},`,
				`{"op":"remove","path":"/0","value":{}},`,
				`{"op":"add","path":"/0","value":0},`,
				`{"op":"test","path":"/1","value":[]},`,
				`{"op":"test","path":"/2","value":0},`,
				`{"op":"remove","path":"/2","value":0}`,
				`]`),
		},
		{
			name: "add to empty object",
			diff: s(`@ ["key"]`,
				`+ "value"`),
			patch: s(`[{"op":"add","path":"/key","value":"value"}]`),
		},
		{
			name: "remove from object",
			diff: s(`@ ["key"]`,
				`- "value"`),
			patch: s(`[`,
				`{"op":"test","path":"/key","value":"value"},`,
				`{"op":"remove","path":"/key","value":"value"}`,
				`]`),
		},
		{
			name: "nested add operation",
			diff: s(`@ ["a","b"]`,
				`+ 1`),
			patch: s(`[{"op":"add","path":"/a/b","value":1}]`),
		},
		{
			name: "array element replacement",
			diff: s(`@ [0]`,
				`- 1`,
				`+ 2`),
			patch: s(`[`,
				`{"op":"test","path":"/0","value":1},`,
				`{"op":"remove","path":"/0","value":1},`,
				`{"op":"add","path":"/0","value":2}`,
				`]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := ReadDiffString(tt.diff)
			if err != nil {
				t.Errorf("Error reading diff: %v", err)
			}
			gotJson, err := diff.RenderPatch()
			if err != nil && !tt.wantErr {
				t.Errorf("Want no err. Got %v", err)
			}
			if err == nil && tt.wantErr {
				t.Errorf("Want err. Got nil")
			}
			got, err := ReadJsonString(gotJson)
			if err != nil {
				t.Errorf("Error reading JSON Patch: %v", err)
			}
			want, err := ReadJsonString(tt.patch)
			if err != nil {
				t.Errorf("Error reading patch: %v", err)
			}
			if !want.Equals(got) {
				t.Errorf("Want %v. Got %v", tt.patch, gotJson)
			}
		})
	}
}

func TestDiffRenderMerge(t *testing.T) {
	tests := []struct {
		name  string
		diff  string
		merge string
	}{
		{
			name: "simple merge add value",
			diff: s(
				`^ {"Merge":true}`,
				`@ []`,
				`+ 1`,
			),
			merge: `1`,
		},
		{
			name: "merge add to object key",
			diff: s(
				`^ {"Merge":true}`,
				`@ ["foo"]`,
				`+ 1`,
			),
			merge: `{"foo":1}`,
		},
		{
			name: "merge add nested value",
			diff: s(
				`^ {"Merge":true}`,
				`@ ["a","b"]`,
				`+ 1`,
			),
			merge: `{"a":{"b":1}}`,
		},
		{
			name: "merge add string value",
			diff: s(
				`^ {"Merge":true}`,
				`@ ["name"]`,
				`+ "John"`,
			),
			merge: `{"name":"John"}`,
		},
		{
			name: "merge add boolean value",
			diff: s(
				`^ {"Merge":true}`,
				`@ ["active"]`,
				`+ true`,
			),
			merge: `{"active":true}`,
		},
		{
			name: "merge add null value",
			diff: s(
				`^ {"Merge":true}`,
				`@ ["value"]`,
				`+ null`,
			),
			merge: `{"value":null}`,
		},
		{
			name: "merge remove to void",
			diff: s(
				`^ {"Merge":true}`,
				`@ []`,
				`+`,
			),
			merge: `null`,
		},
		{
			name: "merge remove object key to void",
			diff: s(
				`^ {"Merge":true}`,
				`@ ["key"]`,
				`+`,
			),
			merge: `{"key":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := ReadDiffString(tt.diff)
			if err != nil {
				t.Errorf("Error reading diff: %v", err)
			}
			got, err := diff.RenderMerge()
			if err != nil {
				t.Errorf("Error rendering diff as merge patch: %v", err)
			}
			if got != tt.merge {
				t.Errorf("Want %v. Got %v", tt.merge, got)
			}
		})
	}
}
