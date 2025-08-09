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
