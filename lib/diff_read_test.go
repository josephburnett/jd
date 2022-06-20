package jd

import "testing"

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
		diff, err := ReadPatchString(tc.patch)
		if err != nil && !tc.wantErr {
			t.Errorf("Wanted no error. Got %v", err)
		}
		if err == nil && tc.wantErr {
			t.Errorf("Wanted an error. Got nil")
		}
		if err != nil && tc.wantErr {
			// Everything is okay
			continue
		}
		got := diff.Render()
		if got != tc.diff {
			t.Errorf("Wanted \n%q. Got \n%q", tc.diff, got)
		}
	}
}

func TestReadMerge(t *testing.T) {
	cases := []struct {
		patch string
		diff  string
	}{{
		patch: `{"a":1}`,
		diff: s(
			`@ [["MERGE"],"a"]`,
			`+ 1`,
		),
	}, {
		patch: ``,
		diff: s(
			`@ [["MERGE"]]`,
			`+`,
		),
	}, {
		patch: `null`,
		diff: s(
			`@ [["MERGE"]]`,
			`+`,
		),
	}, {
		patch: `[1,2,3]`,
		diff: s(
			`@ [["MERGE"]]`,
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
