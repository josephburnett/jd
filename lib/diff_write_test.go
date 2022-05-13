package jd

import (
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
}

func checkDiffRender(t *testing.T, a, b string, diffLines ...string) {
	diff := ""
	for _, dl := range diffLines {
		diff += dl + "\n"
	}
	aJson, err := ReadJsonString(a)
	if err != nil {
		t.Errorf(err.Error())
	}
	bJson, err := ReadJsonString(b)
	if err != nil {
		t.Errorf(err.Error())
	}
	d := aJson.diff(bJson, nil, []Metadata{}, strictPatchStrategy).Render()
	if d != diff {
		t.Errorf("%v.diff(%v) = %v. Want %v.", a, b, d, diff)
	}
}

func TestDiffRenderPatch(t *testing.T) {
	testCases := []struct {
		diff    string
		patch   string
		wantErr bool
	}{{
		diff: `@ ["foo"]` + "\n" +
			`+ 1`,
		patch: `[{"op":"add","path":"/foo","value":1}]`,
	}, {
		diff: `@ ["foo"]` + "\n" +
			`- 1`,
		patch: `[{"op":"test","path":"/foo","value":1},` +
			`{"op":"remove","path":"/foo","value":1}]`,
	}, {
		diff: `@ ["foo"]` + "\n" +
			`- 1` + "\n" +
			`+ 2`,
		patch: `[{"op":"test","path":"/foo","value":1},` +
			`{"op":"remove","path":"/foo","value":1},` +
			`{"op":"add","path":"/foo","value":2}]`,
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
			t.Errorf("Want %v. Got %v", want, got)
		}
	}
}
