package jd

import (
	"strings"
	"testing"
)

func TestDiffRender(t *testing.T) {
	checkDiffRender(t, `{"a":1}`, `{"a":2}`,
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkDiffRender(t, `{"a":{"b":1}}`, `{"a":{"b":2}}`,
		`@ ["a","b"]`,
		`- 1`,
		`+ 2`)
	checkDiffRender(t, `{"a":{"b":1}}`, `{"a":{"c":2}}`,
		`@ ["a","b"]`,
		`- 1`,
		`@ ["a","c"]`,
		`+ 2`)
	checkDiffRender(t, `{"a":{"b":1}}`, `{"c":{"b":1}}`,
		`@ ["a"]`,
		`- {"b":1}`,
		`@ ["c"]`,
		`+ {"b":1}`)
	// String changes
	checkDiffRender(t, `{"a":"bar"}`, `{"a":"baz"}`,
		`@ ["a"]`,
		`- "bar"`,
		`+ "baz"`)
	// Array of strings
	checkDiffRender(t, `{"qux":["foobar","foobaz"]}`, `{"qux":["fooarrr","foobaz"]}`,
		`@ ["qux",0]`,
		`[`,
		`- "foobar"`,
		`+ "fooarrr"`,
		`  "foobaz"`,
	)
	// Addition only
	checkDiffRender(t, `{"str":""}`, `{"str":"abc"}`,
		`@ ["str"]`,
		`- ""`,
		`+ "abc"`)
	// Removal only
	checkDiffRender(t, `{"str":"abc"}`, `{"str":""}`,
		`@ ["str"]`,
		`- "abc"`,
		`+ ""`)
	// Nested strings
	checkDiffRender(t, `{"a":{"b":"hello"}}`, `{"a":{"b":"world"}}`,
		`@ ["a","b"]`,
		`- "hello"`,
		`+ "world"`)
	// Multiple string changes
	checkDiffRender(t, `{"a":"foo","b":"bar"}`, `{"a":"baz","b":"qux"}`,
		`@ ["a"]`,
		`- "foo"`,
		`+ "baz"`,
		`@ ["b"]`,
		`- "bar"`,
		`+ "qux"`)
	// Key change
	checkDiffRender(t, `{"a":"foo"}`, `{"b":"foo"}`,
		`@ ["a"]`,
		`- "foo"`,
		`@ ["b"]`,
		`+ "foo"`)
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
	d := aJson.diff(bJson, nil, []Option{}, strictPatchStrategy).Render()
	if d != diff {
		t.Errorf("%v.diff(%v) = %v. Want %v.", a, b, d, diff)
	}

	// Test with color
	coloredDiff := aJson.diff(bJson, nil, []Option{}, strictPatchStrategy).Render(COLOR)
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
	result := ""
	inColor := false
	for i := 0; i < len(s); i++ {
		// detect a color code (starts coloring)
		if !inColor && i+1 < len(s) && s[i] == '\033' && s[i+1] == '[' {
			inColor = true
			i++ // skip '['
			continue
		}
		// if not colored, add the character to the result
		if !inColor {
			result += string(s[i])
		}
		// detect the reset color code (ends coloring)
		if inColor && i+1 < len(s) && s[i] == '[' && s[i+1] == '0' && i+2 < len(s) && s[i+2] == 'm' {
			inColor = false
			i += 2
		}
	}
	return result
}

// stripAnsiCodes removes ANSI color escape sequences from a string
func stripAnsiCodes(s string) string {
	result := ""
	inEscape := false

	for i := 0; i < len(s); i++ {
		if !inEscape && i+1 < len(s) && s[i] == '\033' && s[i+1] == '[' {
			inEscape = true
			i++ // skip the '['
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(s[i])
	}
	return result
}

func TestDiffRenderPatch(t *testing.T) {
	testCases := []struct {
		diff    string
		patch   string
		wantErr bool
	}{{
		diff: s(`@ ["foo"]`,
			`+ 1`),
		patch: s(`[{"op":"add","path":"/foo","value":1}]`),
	}, {
		diff: s(`@ ["foo"]`,
			`- 1`),
		patch: s(`[`,
			`{"op":"test","path":"/foo","value":1},`,
			`{"op":"remove","path":"/foo","value":1}`,
			`]`),
	}, {
		diff: s(`@ ["foo"]`,
			`- 1`,
			`+ 2`),
		patch: s(`[`,
			`{"op":"test","path":"/foo","value":1},`,
			`{"op":"remove","path":"/foo","value":1},`,
			`{"op":"add","path":"/foo","value":2}`,
			`]`),
	}, {
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
	}}

	for _, tc := range testCases {
		diff, err := ReadDiffString(tc.diff)
		if err != nil {
			t.Errorf("Error reading diff: %v", err)
		}
		gotJson, err := diff.RenderPatch()
		if err != nil && !tc.wantErr {
			t.Errorf("Want no err. Got %v", err)
		}
		if err == nil && tc.wantErr {
			t.Errorf("Want err. Got nil")
		}
		got, err := ReadJsonString(gotJson)
		if err != nil {
			t.Errorf("Error reading JSON Patch: %v", err)
		}
		want, err := ReadJsonString(tc.patch)
		if err != nil {
			t.Errorf("Error reading patch: %v", err)
		}
		if !want.Equals(got) {
			t.Errorf("Want %v. Got %v", tc.patch, gotJson)
		}
	}
}

func TestDiffRenderMerge(t *testing.T) {
	cases := []struct {
		diff  string
		merge string
	}{{
		diff: s(
			`^ {"Merge":true}`,
			`@ []`,
			`+ 1`,
		),
		merge: `1`,
	}, {
		diff: s(
			`^ {"Merge":true}`,
			`@ ["foo"]`,
			`+ 1`,
		),
		merge: `{"foo":1}`,
	}}

	for _, c := range cases {
		d, err := ReadDiffString(c.diff)
		if err != nil {
			t.Errorf("Error reading diff: %v", err)
		}
		s, err := d.RenderMerge()
		if err != nil {
			t.Errorf("Error rendering diff as merge patch: %v", err)
		}
		if s != c.merge {
			t.Errorf("Want %v. Got %v", c.merge, s)
		}
	}
}
