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
			`[{"op":"add","path":"/foo","value":1},`,
			`{"op":"test","path":"/foo","value":1},`,
			`{"op":"remove","path":"/foo","value":1}]`,
		),
		diff: s(
			`@ ["foo"]`,
			`+ 1`,
			`@ ["foo"]`,
			`- 1`,
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
					require.True(t, d.Before[0].Equals(c.wantBefore))
				}
				if c.wantAfter == nil {
					require.Nil(t, d.After)
				} else {
					require.Len(t, d.After, 1)
					require.True(t, d.After[0].Equals(c.wantAfter))
				}
				require.Equal(t, c.patch[c.wantConsumed:], rest)
			}
		})
	}
}
