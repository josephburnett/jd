package jd

import "testing"

func TestReadPatch(t *testing.T) {
	cases := []struct {
		patch   string
		diff    string
		wantErr bool
	}{{
		patch: s(`[{"op":"add","path":"/foo","value":1}]`),
		diff: s(`@ ["foo"]`,
			`+ 1`),
	}, {
		patch: s(`[{"op":"test","path":"/foo","value":1},`,
			`{"op":"remove","path":"/foo","value":1}]`),
		diff: s(`@ ["foo"]`,
			`- 1`),
	}, {
		patch: s(`[{"op":"add","path":"/foo","value":1},`,
			`{"op":"test","path":"/foo","value":1},`,
			`{"op":"remove","path":"/foo","value":1}]`),
		diff: s(`@ ["foo"]`,
			`+ 1`,
			`@ ["foo"]`,
			`- 1`),
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
