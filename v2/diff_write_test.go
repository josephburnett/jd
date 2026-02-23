package jd

import (
	"fmt"
	"strings"
	"testing"
	"time"
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

	// Test with color (line-level only, no character-level LCS)
	coloredDiff := aJson.diff(bJson, nil, newOptions([]Option{}), strictPatchStrategy).Render(COLOR)
	strippedDiff := stripAnsiCodes(coloredDiff)
	expectedDiffWithColorHeader := `^ "COLOR"` + "\n" + diff
	if strippedDiff != expectedDiffWithColorHeader {
		t.Errorf("%v.diff(%v) with color (stripped) = %v. Want %v.", a, b, strippedDiff, expectedDiffWithColorHeader)
	}

	// Test with color-words (character-level diff).
	// Verify that uncolored parts in string diffs match between + and - lines.
	colorWordsDiff := aJson.diff(bJson, nil, newOptions([]Option{}), strictPatchStrategy).Render(COLOR_WORDS)
	lines := strings.Split(colorWordsDiff, "\n")
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

func TestDiffRenderColorNoCharDiff(t *testing.T) {
	// COLOR alone should produce line-level red/green coloring without
	// character-level ANSI codes inside the string values.
	a, _ := ReadJsonString(`{"a":"foobar"}`)
	b, _ := ReadJsonString(`{"a":"foobaz"}`)
	d := a.Diff(b)
	rendered := d.Render(COLOR)

	lines := strings.Split(rendered, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '-' && strings.Contains(line, `"foo`) {
			// The minus line should be wrapped in red at the line level,
			// not contain inline color codes within the quoted string.
			// With line-level coloring: \033[31m- "foobar"\n\033[0m
			// Character-level would have color codes between individual chars.
			stripped := stripAnsiCodes(line)
			if stripped != `- "foobar"` {
				t.Errorf("COLOR minus line stripped = %q, want %q", stripped, `- "foobar"`)
			}
			// Should NOT have color codes inside the quotes (character-level)
			quoteStart := strings.Index(line, `"`)
			quoteEnd := strings.LastIndex(line, `"`)
			if quoteStart >= 0 && quoteEnd > quoteStart {
				inside := line[quoteStart : quoteEnd+1]
				if strings.Contains(inside, "\033[31m") || strings.Contains(inside, "\033[32m") {
					t.Errorf("COLOR should not have character-level coloring inside strings, got: %q", inside)
				}
			}
		}
	}
}

func TestDiffRenderColorWordsCharDiff(t *testing.T) {
	// COLOR_WORDS should produce character-level ANSI highlighting inside
	// the string values.
	a, _ := ReadJsonString(`{"a":"foobar"}`)
	b, _ := ReadJsonString(`{"a":"foobaz"}`)
	d := a.Diff(b)
	rendered := d.Render(COLOR_WORDS)

	lines := strings.Split(rendered, "\n")
	foundCharLevel := false
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '-' && strings.Contains(line, "foo") {
			// Should have color codes inside the string (character-level diff)
			quoteStart := strings.Index(line, `"`)
			quoteEnd := strings.LastIndex(line, `"`)
			if quoteStart >= 0 && quoteEnd > quoteStart {
				inside := line[quoteStart : quoteEnd+1]
				if strings.Contains(inside, "\033[31m") {
					foundCharLevel = true
				}
			}
		}
	}
	if !foundCharLevel {
		t.Errorf("COLOR_WORDS should produce character-level coloring inside strings, got:\n%s", rendered)
	}
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

func TestDiffElementOptionsRendering(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		want    string
	}{
		{
			name:    "no options",
			options: []Option{},
			want:    "",
		},
		{
			name:    "single SET option",
			options: []Option{SET},
			want:    `^ "SET"` + "\n",
		},
		{
			name:    "single MERGE option",
			options: []Option{MERGE},
			want:    `^ "MERGE"` + "\n",
		},
		{
			name:    "single MULTISET option",
			options: []Option{MULTISET},
			want:    `^ "MULTISET"` + "\n",
		},
		{
			name:    "single COLOR option",
			options: []Option{COLOR},
			want:    `^ "COLOR"` + "\n",
		},
		{
			name:    "precision option",
			options: []Option{Precision(0.01)},
			want:    `^ {"precision":0.01}` + "\n",
		},
		{
			name:    "setkeys option",
			options: []Option{SetKeys("id", "name")},
			want:    `^ {"setkeys":["id","name"]}` + "\n",
		},
		{
			name:    "multiple simple options",
			options: []Option{MERGE, COLOR},
			want:    `^ "MERGE"` + "\n" + `^ "COLOR"` + "\n",
		},
		{
			name:    "multiple mixed options",
			options: []Option{SET, Precision(0.001), COLOR},
			want:    `^ "SET"` + "\n" + `^ {"precision":0.001}` + "\n" + `^ "COLOR"` + "\n",
		},
		{
			name:    "path option",
			options: []Option{PathOption(Path{PathKey("users")}, SET)},
			want:    `^ {"@":["users"],"^":["SET"]}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a DiffElement with the options to test rendering
			element := DiffElement{
				Options: tt.options,
				Path:    Path{PathKey("test")},
				Add:     []JsonNode{jsonString("value")},
			}
			rendered := element.Render()

			// Extract just the options part (everything before "@ ")
			parts := strings.Split(rendered, "@ ")
			got := parts[0]

			if got != tt.want {
				t.Errorf("DiffElement options rendering = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDiffRenderWithOptions(t *testing.T) {
	tests := []struct {
		name      string
		a         string
		b         string
		options   []Option
		wantLines []string
	}{
		{
			name:    "simple diff with SET option",
			a:       `{"items":[1,2,3]}`,
			b:       `{"items":[2,1,4]}`,
			options: []Option{SET},
			wantLines: []string{
				`^ "SET"`,
				`@ ["items",{}]`,
				`- 3`,
				`+ 4`,
			},
		},
		{
			name:    "simple diff with MERGE option",
			a:       `{"a":1}`,
			b:       `{"a":2}`,
			options: []Option{MERGE},
			wantLines: []string{
				`^ "MERGE"`,
				`@ ["a"]`,
				`+ 2`,
			},
		},
		{
			name:    "diff with multiple options",
			a:       `{"price":10.99}`,
			b:       `{"price":11.05}`,
			options: []Option{SET, Precision(0.001)},
			wantLines: []string{
				`^ "SET"`,
				`^ {"precision":0.001}`,
				`@ ["price"]`,
				`- 10.99`,
				`+ 11.05`,
			},
		},
		{
			name:    "diff with no options (no header)",
			a:       `{"a":1}`,
			b:       `{"a":2}`,
			options: []Option{},
			wantLines: []string{
				`@ ["a"]`,
				`- 1`,
				`+ 2`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aJson, err := ReadJsonString(tt.a)
			if err != nil {
				t.Fatalf("Error reading a: %v", err)
			}
			bJson, err := ReadJsonString(tt.b)
			if err != nil {
				t.Fatalf("Error reading b: %v", err)
			}

			diff := aJson.Diff(bJson, tt.options...)
			got := diff.Render(tt.options...)

			want := ""
			for _, line := range tt.wantLines {
				want += line + "\n"
			}

			if got != want {
				t.Errorf("Diff.Render() = %q, want %q", got, want)
			}
		})
	}
}

func TestDiffRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		options []Option
	}{
		{
			name:    "round-trip with SET option",
			a:       `{"items":[1,2,3]}`,
			b:       `{"items":[2,1,4]}`,
			options: []Option{SET},
		},
		{
			name:    "round-trip with MERGE option",
			a:       `{"a":1}`,
			b:       `{"a":2}`,
			options: []Option{MERGE},
		},
		{
			name:    "round-trip with multiple options",
			a:       `{"price":10.99}`,
			b:       `{"price":11.05}`,
			options: []Option{SET, Precision(0.001)},
		},
		{
			name:    "round-trip with path option",
			a:       `{"users":[{"id":1,"name":"alice"}]}`,
			b:       `{"users":[{"id":1,"name":"bob"}]}`,
			options: []Option{PathOption(Path{PathKey("users")}, SET)},
		},
		{
			name:    "round-trip with no options",
			a:       `{"a":1}`,
			b:       `{"a":2}`,
			options: []Option{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original diff
			aJson, err := ReadJsonString(tt.a)
			if err != nil {
				t.Fatalf("Error reading a: %v", err)
			}
			bJson, err := ReadJsonString(tt.b)
			if err != nil {
				t.Fatalf("Error reading b: %v", err)
			}

			originalDiff := aJson.Diff(bJson, tt.options...)

			// Render diff to string
			renderedDiff := originalDiff.Render(tt.options...)

			// Parse diff back from string
			parsedDiff, err := ReadDiffString(renderedDiff)
			if err != nil {
				t.Fatalf("Error parsing rendered diff: %v", err)
			}

			// Render parsed diff again (without passing options since they're now stored in DiffElement.Options)
			reRenderedDiff := parsedDiff.Render()

			// Should be identical (round-trip) unless MERGE option causes normalization
			expectedRerendered := renderedDiff
			// Special case: MERGE option causes {"Merge":true} to be normalized to "MERGE"
			if len(tt.options) > 0 {
				for _, opt := range tt.options {
					if _, isMerge := opt.(mergeOption); isMerge {
						// Replace the legacy format with modern format for comparison
						expectedRerendered = strings.ReplaceAll(expectedRerendered, `^ {"Merge":true}`+"\n", "")
						break
					}
				}
			}

			if expectedRerendered != reRenderedDiff {
				t.Errorf("Round-trip failed.\nOriginal:\n%s\nRe-rendered:\n%s\nExpected:\n%s", renderedDiff, reRenderedDiff, expectedRerendered)
			}
		})
	}
}

func TestLegacyMetadataRoundTrip(t *testing.T) {
	// Test that legacy {"Merge":true} format gets normalized to "MERGE"
	legacyDiff := `^ {"Merge":true}
@ ["a"]
+ 2
`

	// Parse legacy format
	parsedDiff, err := ReadDiffString(legacyDiff)
	if err != nil {
		t.Fatalf("Error parsing legacy diff: %v", err)
	}

	// Should have both Metadata.Merge=true and Options=[MERGE]
	if len(parsedDiff) != 1 {
		t.Fatalf("Expected 1 diff element, got %d", len(parsedDiff))
	}

	element := parsedDiff[0]
	if !element.Metadata.Merge {
		t.Error("Expected Metadata.Merge to be true")
	}

	if len(element.Options) != 1 {
		t.Fatalf("Expected 1 option, got %d", len(element.Options))
	}

	if _, ok := element.Options[0].(mergeOption); !ok {
		t.Errorf("Expected MERGE option, got %T", element.Options[0])
	}

	// Render should normalize to modern format
	rendered := parsedDiff.Render()
	expectedModern := `^ "MERGE"
@ ["a"]
+ 2
`

	if rendered != expectedModern {
		t.Errorf("Legacy format should normalize to modern.\nGot:\n%s\nExpected:\n%s", rendered, expectedModern)
	}
}

func TestDiffRenderPatchError(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		opts []Option
	}{
		{
			name: "set-based diff returns error",
			a:    `[1,2,3]`,
			b:    `[1,3,4]`,
			opts: []Option{SET},
		},
		{
			name: "multiset-based diff returns error",
			a:    `[1,2,3]`,
			b:    `[1,3,4]`,
			opts: []Option{MULTISET},
		},
		{
			name: "setkeys-based diff returns error",
			a:    `[{"id":"a","v":1}]`,
			b:    `[{"id":"a","v":2}]`,
			opts: []Option{SET, SetKeys("id")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aJson, err := ReadJsonString(tt.a)
			if err != nil {
				t.Fatalf("Error reading a: %v", err)
			}
			bJson, err := ReadJsonString(tt.b)
			if err != nil {
				t.Fatalf("Error reading b: %v", err)
			}
			d := aJson.Diff(bJson, tt.opts...)
			_, err = d.RenderPatch()
			if err == nil {
				t.Errorf("Expected error from RenderPatch with %s, got nil", tt.name)
			}
		})
	}
}

func TestRenderMergeVoidAdd(t *testing.T) {
	// Render a merge deletion (void add) with and without color
	de := DiffElement{
		Metadata: Metadata{Merge: true},
		Path:     Path{PathKey("foo")},
		Add:      []JsonNode{voidNode{}},
	}
	// Without color
	rendered := de.Render()
	if !strings.Contains(rendered, "+\n") {
		t.Errorf("expected void add line, got: %q", rendered)
	}
	// With color
	rendered = de.Render(COLOR)
	stripped := stripAnsiCodes(rendered)
	if !strings.Contains(stripped, "+\n") {
		t.Errorf("expected void add in color mode, got: %q", stripped)
	}
}

func TestDiffRenderWithFileOption(t *testing.T) {
	a, err := ReadJsonString(`{"a":1}`)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ReadJsonString(`{"a":2}`)
	if err != nil {
		t.Fatal(err)
	}
	d := a.Diff(b)
	rendered := d.Render(File("a.json"))
	want := "^ {\"file\":\"a.json\"}\n@ [\"a\"]\n- 1\n+ 2\n"
	if rendered != want {
		t.Errorf("got %q, want %q", rendered, want)
	}
}

func TestDiffRenderEmpty(t *testing.T) {
	// Empty diff with options should produce no output
	d := Diff{}
	got := d.Render(SET)
	if got != "" {
		t.Errorf("expected empty string for empty diff. got %q", got)
	}
}

func TestRenderMergeError(t *testing.T) {
	// Non-merge diff element should error
	d := Diff{DiffElement{
		Path: Path{PathKey("a")},
		Add:  []JsonNode{jsonString("b")},
	}}
	_, err := d.RenderMerge()
	if err == nil {
		t.Fatal("expected error for non-merge diff")
	}
}

func TestRenderPatchErrors(t *testing.T) {
	tests := []struct {
		name string
		diff Diff
	}{
		{
			name: "empty diff element",
			diff: Diff{DiffElement{Path: Path{PathKey("a")}}},
		},
		{
			name: "too many before context lines",
			diff: Diff{DiffElement{
				Path:   Path{PathIndex(1)},
				Before: []JsonNode{jsonNumber(1), jsonNumber(2)},
				Remove: []JsonNode{jsonNumber(3)},
			}},
		},
		{
			name: "too many after context lines",
			diff: Diff{DiffElement{
				Path:  Path{PathIndex(1)},
				After: []JsonNode{jsonNumber(1), jsonNumber(2)},
				Add:   []JsonNode{jsonNumber(3)},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.diff.RenderPatch()
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestRenderColorLargeStringNoOOM(t *testing.T) {
	// COLOR alone must not run O(n^2) LCS on large strings.
	old := strings.Repeat("a", 10000)
	new := strings.Repeat("b", 10000)
	a, _ := ReadJsonString(fmt.Sprintf(`"%s"`, old))
	b, _ := ReadJsonString(fmt.Sprintf(`"%s"`, new))
	d := a.Diff(b)
	start := time.Now()
	d.Render(COLOR)
	if time.Since(start) > time.Second {
		t.Fatal("Render(COLOR) on large string took too long — LCS may be running")
	}
}

func TestRenderPatchContextErrors(t *testing.T) {
	tests := []struct {
		name string
		diff Diff
	}{
		{
			name: "before context with empty path",
			diff: Diff{DiffElement{
				Path:   Path{},
				Before: []JsonNode{jsonNumber(1)},
				Remove: []JsonNode{jsonNumber(2)},
			}},
		},
		{
			name: "before context path not PathIndex",
			diff: Diff{DiffElement{
				Path:   Path{PathKey("foo")},
				Before: []JsonNode{jsonNumber(1)},
				Remove: []JsonNode{jsonNumber(2)},
			}},
		},
		{
			name: "after context with empty path",
			diff: Diff{DiffElement{
				Path:   Path{},
				After:  []JsonNode{jsonNumber(1)},
				Remove: []JsonNode{jsonNumber(2)},
			}},
		},
		{
			name: "after context path not PathIndex",
			diff: Diff{DiffElement{
				Path:   Path{PathKey("foo")},
				After:  []JsonNode{jsonNumber(1)},
				Remove: []JsonNode{jsonNumber(2)},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.diff.RenderPatch()
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestRenderPatchVoidSkip(t *testing.T) {
	// void remove should be skipped
	d := Diff{DiffElement{
		Path:   Path{PathKey("a")},
		Remove: []JsonNode{voidNode{}},
		Add:    []JsonNode{jsonNumber(1)},
	}}
	got, err := d.RenderPatch()
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"op":"add","path":"/a","value":1}]`
	wantNode, _ := ReadJsonString(want)
	gotNode, _ := ReadJsonString(got)
	if !wantNode.Equals(gotNode) {
		t.Errorf("want %s, got %s", want, got)
	}
	// void add should also be skipped
	d2 := Diff{DiffElement{
		Path:   Path{PathKey("a")},
		Remove: []JsonNode{jsonNumber(1)},
		Add:    []JsonNode{voidNode{}},
	}}
	got2, err := d2.RenderPatch()
	if err != nil {
		t.Fatal(err)
	}
	want2 := `[{"op":"test","path":"/a","value":1},{"op":"remove","path":"/a","value":1}]`
	wantNode2, _ := ReadJsonString(want2)
	gotNode2, _ := ReadJsonString(got2)
	if !wantNode2.Equals(gotNode2) {
		t.Errorf("want %s, got %s", want2, got2)
	}
}

func TestDiffElementRenderSkipMetadata(t *testing.T) {
	// When global MERGE option is passed and DiffElement has Metadata.Merge=true
	// but no Options, the metadata rendering should be skipped
	de := DiffElement{
		Metadata: Metadata{Merge: true},
		Path:     Path{PathKey("a")},
		Add:      []JsonNode{jsonNumber(1)},
	}
	rendered := de.Render(MERGE)
	if strings.Contains(rendered, "Merge") {
		t.Errorf("expected metadata to be skipped when global MERGE covers it, got: %q", rendered)
	}
}

func TestDiffElementRenderColorWordsNonString(t *testing.T) {
	// COLOR_WORDS with non-string values should fall back to line-level coloring
	de := DiffElement{
		Path:   Path{PathKey("a")},
		Remove: []JsonNode{jsonNumber(1)},
		Add:    []JsonNode{jsonNumber(2)},
	}
	rendered := de.Render(COLOR_WORDS)
	stripped := stripAnsiCodes(rendered)
	if !strings.Contains(stripped, "- 1") || !strings.Contains(stripped, "+ 2") {
		t.Errorf("expected line-level diff for non-string COLOR_WORDS, got: %q", stripped)
	}
}

func TestRenderMergeWithPatchError(t *testing.T) {
	// A merge diff with a Remove value will fail because merge patch can't have old values
	d := Diff{DiffElement{
		Metadata: Metadata{Merge: true},
		Path:     Path{PathKey("a")},
		Remove:   []JsonNode{jsonNumber(1)},
		Add:      []JsonNode{jsonNumber(2)},
	}}
	_, err := d.RenderMerge()
	if err == nil {
		t.Fatal("expected error from RenderMerge with remove value in merge diff")
	}
}
