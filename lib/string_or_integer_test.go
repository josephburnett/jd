package jd

import (
	"strings"
	"testing"
)

func TestPatchJsonStringOrInteger(t *testing.T) {
	tests := []struct {
		a         string
		b         string
		diff      []string
		wantError bool
	}{{
		a: `{"0":{}}`,
		b: `{"0":{"foo":"bar"}}`,
		diff: ss(
			`@ ["0", "foo"]`,
			`+ "bar"`,
		),
	}, {
		a: `{"0":{}}`,
		b: `{"0":{"foo":"bar"}}`,
		diff: ss(
			`@ [0, "foo"]`,
			`+ "bar"`,
		),
	}, {
		a: `[]`,
		b: `[1]`,
		diff: ss(
			`@ ["0"]`,
			`+ 1`,
		),
	}, {
		a: `[]`,
		b: `[1]`,
		diff: ss(
			`@ [0]`,
			`+ 1`,
		),
	}}

	for _, tt := range tests {
		diffString := strings.Join(tt.diff, "\n")
		initial, err := ReadJsonString(tt.a)
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		diff, err := ReadDiffString(diffString)
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		expect, err := ReadJsonString(tt.b)
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		// Coerce to patch format so we'll create a jsonStringOrNumber object when reading the diff.
		patchString, err := diff.RenderPatch()
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		patchDiff, err := ReadPatchString(patchString)
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		b, err := initial.Patch(patchDiff)
		if tt.wantError && err == nil {
			t.Errorf("wanted error but got none")
		}
		if !tt.wantError && err != nil {
			t.Errorf("wanted no error but got %v", err)
		}
		if !tt.wantError && !expect.Equals(b) {
			t.Errorf("%v.Patch(%v) = %v. Want %v.",
				tt.a, diffString, b, tt.b)
		}
	}

}
