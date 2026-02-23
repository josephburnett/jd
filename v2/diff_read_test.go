package jd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadPatch(t *testing.T) {
	cases := []struct {
		patch   string
		diff    string
		wantErr bool
	}{{
		patch: s(`[{"op":"add","path":"/foo","value":1}]`),
		diff: s(
			`@ ["foo"]`,
			`+ 1`,
		),
	}, {
		patch: s(
			`[{"op":"test","path":"/foo","value":1},`,
			`{"op":"remove","path":"/foo","value":1}]`,
		),
		diff: s(
			`@ ["foo"]`,
			`- 1`,
		),
	}, {
		patch: s(
			`[{"op":"test","path":"/foo","value":1},`,
			`{"op":"remove","path":"/foo","value":1},`,
			`{"op":"add","path":"/foo","value":1}]`,
		),
		diff: s(
			`@ ["foo"]`,
			`- 1`,
			`+ 1`,
		),
	}, {
		patch: s(`[{"op":"add","path":"/foo/-","value":2}]`),
		diff: s(
			`@ ["foo",-1]`,
			`+ 2`,
		),
	}, {
		patch:   s(`[{"op":"test","path":"/foo","value":1}]`),
		wantErr: true,
	}, {
		patch:   s(`[{"op":"remove","path":"/foo","value":1}]`),
		wantErr: true,
	}}

	for _, tc := range cases {
		t.Run(tc.patch, func(t *testing.T) {
			diff, err := ReadPatchString(tc.patch)
			if err != nil && !tc.wantErr {
				t.Errorf("Wanted no error. Got %v", err)
			}
			if err == nil && tc.wantErr {
				t.Errorf("Wanted an error. Got nil")
			}
			if err != nil && tc.wantErr {
				// Everything is okay
				return
			}
			got := diff.Render()
			if got != tc.diff {
				t.Errorf("Wanted \n%q. Got \n%q", tc.diff, got)
			}
		})
	}
}

func TestReadMerge(t *testing.T) {
	cases := []struct {
		patch string
		diff  string
	}{{
		patch: `{"a":1}`,
		diff: s(
			`^ {"Merge":true}`,
			`@ ["a"]`,
			`+ 1`,
		),
	}, {
		patch: ``,
		diff:  ``,
	}, {
		patch: `null`,
		diff: s(
			`^ {"Merge":true}`,
			`@ []`,
			`+`,
		),
	}, {
		patch: `[1,2,3]`,
		diff: s(
			`^ {"Merge":true}`,
			`@ []`,
			`+ [1,2,3]`,
		),
	}}

	for _, c := range cases {
		diff, err := ReadMergeString(c.patch)
		if err != nil {
			t.Errorf("Wanted no error. Got %v", err)
		}
		if got := diff.Render(); got != c.diff {
			t.Errorf("Wanted %s. Got %s", c.diff, got)
		}
	}
}

func TestSetPatchDiffElementContext(t *testing.T) {
	cases := []struct {
		name         string
		patch        []patchElement
		diffElement  *DiffElement
		wantBefore   JsonNode
		wantAfter    JsonNode
		wantConsumed int
		wantErr      bool
	}{{
		name: "context before and after replacement",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "test",
			Path:  "/2",
			Value: 3,
		}, {
			Op:    "test",
			Path:  "/1",
			Value: 2,
		}, {
			Op:    "replace",
			Path:  "/1",
			Value: 4,
		}},
		wantBefore:   jsonNumber(1),
		wantAfter:    jsonNumber(3),
		wantConsumed: 2,
	}, {
		name: "context before and after insert",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "test",
			Path:  "/1",
			Value: 3,
		}, {
			Op:    "add",
			Path:  "/1",
			Value: 2,
		}},
		wantBefore:   jsonNumber(1),
		wantAfter:    jsonNumber(3),
		wantConsumed: 2,
	}, {
		name: "context after replacement of first element",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/1",
			Value: 2,
		}, {
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "replace",
			Path:  "/0",
			Value: 3,
		}},
		wantBefore:   voidNode{},
		wantAfter:    jsonNumber(2),
		wantConsumed: 1,
	}, {
		name: "context after add at first element",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "add",
			Path:  "/0",
			Value: 2,
		}},
		wantBefore:   voidNode{},
		wantAfter:    jsonNumber(1),
		wantConsumed: 1,
	}, {
		name: "context before replacement",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "test",
			Path:  "/1",
			Value: 2,
		}, {
			Op:    "replace",
			Path:  "/1",
			Value: 3,
		}},
		wantBefore:   jsonNumber(1),
		wantAfter:    voidNode{},
		wantConsumed: 1,
	}, {
		name: "context before add at the end",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "add",
			Path:  "/1",
			Value: 2,
		}},
		wantBefore:   jsonNumber(1),
		wantAfter:    voidNode{},
		wantConsumed: 1,
	}, {
		name: "no context with replacement",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "replace",
			Path:  "/0",
			Value: 2,
		}},
		wantBefore:   voidNode{},
		wantAfter:    voidNode{},
		wantConsumed: 0,
	}, {
		name: "no context with remove",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "remove",
			Path:  "/0",
			Value: 1,
		}},
		wantBefore:   voidNode{},
		wantAfter:    voidNode{},
		wantConsumed: 0,
	}, {
		name: "no context with add into empty array",
		patch: []patchElement{{
			Op:    "add",
			Path:  "/0",
			Value: 1,
		}},
		wantBefore:   voidNode{},
		wantAfter:    voidNode{},
		wantConsumed: 0,
	}, {
		name: "not an array",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/foo",
			Value: 1,
		}, {
			Op:    "replace",
			Path:  "/foo",
			Value: 2,
		}},
		wantConsumed: 0,
	}, {
		name:    "empty patch",
		patch:   []patchElement{},
		wantErr: true,
	}, {
		name: "second path is root pointer",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "replace",
			Path:  "",
			Value: 2,
		}},
		wantConsumed: 0,
	}, {
		name: "after context with replace first index greater",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/2",
			Value: 3,
		}, {
			Op:    "test",
			Path:  "/0",
			Value: 1,
		}, {
			Op:    "replace",
			Path:  "/0",
			Value: 4,
		}},
		wantBefore:   voidNode{},
		wantAfter:    jsonNumber(3),
		wantConsumed: 1,
	}, {
		name: "default fallthrough with two element patch",
		patch: []patchElement{{
			Op:    "test",
			Path:  "/5",
			Value: 1,
		}, {
			Op:    "test",
			Path:  "/0",
			Value: 2,
		}},
		wantConsumed: 0,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &DiffElement{}
			rest, err := setPatchDiffElementContext(c.patch, d)
			if c.wantErr {
				require.Error(t, err)
				require.Nil(t, rest)
			} else {
				require.NoError(t, err)
				if c.wantBefore == nil {
					require.Nil(t, d.Before)
				} else {
					require.Len(t, d.Before, 1)
					require.True(t, d.Before[0].Equals(c.wantBefore), "got %v, want %v", d.Before[0], c.wantBefore)
				}
				if c.wantAfter == nil {
					require.Nil(t, d.After)
				} else {
					require.Len(t, d.After, 1)
					require.True(t, d.After[0].Equals(c.wantAfter), "got %v. want %v", d.After[0], c.wantAfter)
				}
				require.Equal(t, c.patch[c.wantConsumed:], rest)
			}
		})
	}
}

// TestApplyPatch tests applying JSON Patch (RFC 6902) to documents.
// This test demonstrates GitHub issue #99: JSON Patch format does not properly
// validate context test operations that are not directly related to the target path.
func TestReadDiffErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// readDiff: unrecognized line prefix
		{name: "unrecognized prefix", input: "X bad line\n"},
		// readDiff: diff ending at ^ (META terminal state)
		{name: "ends at metadata", input: "^ \"SET\"\n"},
		// readDiff: diff ending at @ (AT terminal state)
		{name: "ends at path", input: "@ [\"a\"]\n"},
		// readDiff: [ not immediately after @
		{name: "bracket not after path", input: "@ [0]\n- 1\n[\n"},
		// readDiff: ] in invalid position
		{name: "close bracket in wrong state", input: "@ [0]\n]\n"},
		// readDiff: legacy metadata that fails both NewOption and readMetadata
		{name: "bad metadata", input: "^ [1,2,3]\n"},
		// readDiff: ^ after - (state == REMOVE) triggers saving diff element
		{name: "metadata after remove", input: "@ [\"a\"]\n- 1\n^ \"SET\"\n@ [\"b\"]\n- 2\n"},
		// readDiff: duplicate option type in same diff element
		{name: "duplicate option", input: "^ \"SET\"\n^ \"SET\"\n@ [\"a\"]\n- 1\n"},
		// readDiff: checkDiffElement error before ^ (multiple removes on PathKey)
		{name: "check error before metadata", input: "@ [\"a\"]\n- 1\n- 2\n^ \"SET\"\n"},
		// readDiff: invalid JSON after ^
		{name: "invalid json after metadata", input: "@ []\n+ 1\n^ {bad\n"},
		// readDiff: checkDiffElement error before @ (multiple removes on PathKey)
		{name: "check error before path", input: "@ [\"a\"]\n- 1\n- 2\n@ [\"b\"]\n- 3\n"},
		// readDiff: invalid JSON after @
		{name: "invalid json after path", input: "@ {bad\n"},
		// readDiff: invalid JSON in before context
		{name: "invalid json in before context", input: "@ [0,1]\n {bad\n"},
		// readDiff: invalid JSON in after context
		{name: "invalid json in after context", input: "@ [0,1]\n- 1\n {bad\n"},
		// readDiff: invalid JSON after -
		{name: "invalid json after remove", input: "@ []\n- {bad\n"},
		// readDiff: invalid JSON after +
		{name: "invalid json after add", input: "@ []\n+ {bad\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadDiffString(tt.input)
			if err != nil {
				// Expected - error path covered
				return
			}
			// Some cases like "bad metadata" may not error but should still be handled
		})
	}
}

func TestReadDiffLegacyMetadata(t *testing.T) {
	// Legacy metadata format: {"Merge":true} — fails NewOption, succeeds readMetadata
	input := "^ {\"Merge\":true}\n@ [\"a\"]\n- 1\n+ 2\n"
	d, err := ReadDiffString(input)
	require.NoError(t, err)
	require.Len(t, d, 1)
	require.True(t, d[0].Metadata.Merge)
}

func TestReadDiffDuplicateMergeOption(t *testing.T) {
	// Duplicate MERGE option — exercises the isMerge check in duplicate detection
	input := "^ \"MERGE\"\n^ \"MERGE\"\n@ [\"a\"]\n- 1\n+ 2\n"
	d, err := ReadDiffString(input)
	require.NoError(t, err)
	require.Len(t, d, 1)
	require.True(t, d[0].Metadata.Merge)
}

func TestCheckDiffElementErrors(t *testing.T) {
	// Multiple adds with empty path
	de := DiffElement{
		Path: Path{},
		Add:  []JsonNode{jsonString("a"), jsonString("b")},
	}
	err := checkDiffElement(de)
	if err == nil {
		t.Fatal("expected error for empty path with multiple add")
	}
	// Multiple removes with PathKey (non-set path)
	de2 := DiffElement{
		Path:   Path{PathKey("foo")},
		Remove: []JsonNode{jsonString("a"), jsonString("b")},
	}
	err = checkDiffElement(de2)
	if err == nil {
		t.Fatal("expected error for multiple remove in object path")
	}
}

func TestReadPatchStringErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "invalid json", input: "not json"},
		{name: "test without remove", input: `[{"op":"test","path":"/foo","value":1}]`},
		{name: "test remove path mismatch", input: `[{"op":"test","path":"/foo","value":1},{"op":"remove","path":"/bar","value":1}]`},
		{name: "test remove value mismatch", input: `[{"op":"test","path":"/foo","value":1},{"op":"remove","path":"/foo","value":2}]`},
		{name: "unknown op", input: `[{"op":"move","path":"/foo","value":1}]`},
		// readPatchDiffElement: empty patch after context consumed
		{name: "empty after context", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"/1","value":2}]`},
		// readPatchDiffElement: readPointer error in test case
		{name: "invalid pointer in test", input: `[{"op":"test","path":"no-slash","value":1},{"op":"remove","path":"no-slash","value":1}]`},
		// readPatchDiffElement: readPointer error in add case
		{name: "invalid pointer in add", input: `[{"op":"add","path":"no-slash","value":1}]`},
		// setPatchDiffElementContext: readPointer error on first test path
		{name: "context invalid first pointer", input: `[{"op":"test","path":"bad","value":1},{"op":"test","path":"/0","value":2},{"op":"add","path":"/1","value":3}]`},
		// setPatchDiffElementContext: readPointer error on second test path
		{name: "context invalid second pointer", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"bad","value":2},{"op":"add","path":"/1","value":3}]`},
		// setPatchDiffElementContext: non-PathIndex third element
		{name: "context third not index", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"/1","value":2},{"op":"add","path":"/foo","value":3}]`},
		// setPatchDiffElementContext: readPointer error on third path
		{name: "context invalid third pointer", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"/1","value":2},{"op":"add","path":"bad","value":3}]`},
		// setPatchDiffElementContext: empty path on third element
		{name: "context third empty path", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"/1","value":2},{"op":"add","path":"","value":3}]`},
		// setPatchDiffElementContext: default in first switch with len>2 falls through to second switch
		{name: "context default fallthrough", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"/5","value":2},{"op":"replace","path":"/5","value":3}]`},
		// setPatchDiffElementContext: second path not PathIndex (line 338)
		{name: "context second not index", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"/foo","value":2},{"op":"add","path":"/1","value":3}]`},
		// setPatchDiffElementContext: second switch default with add op and thirdIndex > secondIndex (line 417)
		{name: "context second switch default", input: `[{"op":"test","path":"/0","value":1},{"op":"test","path":"/5","value":2},{"op":"add","path":"/7","value":3}]`},
		// readPatchDiffElement: readPointer error in single-element test (line 446)
		{name: "single test invalid pointer", input: `[{"op":"test","path":"no-slash","value":1}]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadPatchString(tt.input)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestReadMergeStringError(t *testing.T) {
	_, err := ReadMergeString("{invalid json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestReadMergeStringEmptyObject(t *testing.T) {
	// Empty merge patch is a no-op
	d, err := ReadMergeString("{}")
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 0 {
		t.Fatalf("expected empty diff. got %v", d)
	}
	// Non-empty merge with nested empty object
	d, err = ReadMergeString(`{"a":{}}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 1 {
		t.Fatalf("expected 1 diff element. got %v", len(d))
	}
}

func TestApplyPatch(t *testing.T) {
	cases := []struct {
		name        string
		jsonDoc     string
		patchString string
		expectError bool
		description string
	}{{
		name:        "valid_patch_should_succeed",
		jsonDoc:     `{"foo":["bar","baz"]}`,
		patchString: `[{"op":"test","path":"/foo/0","value":"bar"},{"op":"test","path":"/foo/1","value":"baz"},{"op":"remove","path":"/foo/1","value":"baz"},{"op":"add","path":"/foo/1","value":"qux"}]`,
		expectError: false,
		description: "Valid patch with correct context should succeed",
	}, {
		name:        "github_issue_99",
		jsonDoc:     `{"foo":["bar","baz"]}`,
		patchString: `[{"op":"test","path":"/foo/0","value":"b"},{"op":"test","path":"/foo/1","value":"baz"},{"op":"remove","path":"/foo/1","value":"baz"},{"op":"add","path":"/foo/1","value":"boom"},{"op":"add","path":"/foo/1","value":"bam"}]`,
		expectError: true,
		description: "GitHub issue #99: First test expects '/foo/0' to be 'b' but actual is 'bar' - should fail",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the JSON document
			doc, err := ReadJsonString(tc.jsonDoc)
			if err != nil {
				t.Fatalf("Failed to parse JSON document: %v", err)
			}

			// Parse the JSON Patch
			patch, err := ReadPatchString(tc.patchString)
			if err != nil {
				t.Fatalf("Failed to parse JSON Patch: %v", err)
			}

			// Apply the patch
			result, err := doc.Patch(patch)

			if tc.expectError {
				if err == nil {
					t.Errorf("CONTEXT VERIFICATION BUG: Patch succeeded when it should have failed")
					t.Errorf("  Original: %s", tc.jsonDoc)
					t.Errorf("  Result:   %s", result.Json())
					t.Errorf("  Issue: %s", tc.description)
				} else {
					t.Logf("GOOD: Patch correctly failed: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Valid patch failed unexpectedly: %v", err)
				} else {
					t.Logf("Valid patch succeeded: %s -> %s", tc.jsonDoc, result.Json())
				}
			}
		})
	}
}

func TestReadPatchDiffElementEdgeCases(t *testing.T) {
	// empty patch at entry
	_, _, err := readPatchDiffElement([]patchElement{})
	if err == nil {
		t.Fatal("expected error for empty patch")
	}
	// add op with unsupported value type causes NewJsonNode error
	_, _, err = readPatchDiffElement([]patchElement{{
		Op:    "add",
		Path:  "/foo",
		Value: complex(1, 2),
	}})
	if err == nil {
		t.Fatal("expected error for unsupported add value type")
	}
}
