package jd

import "testing"

func TestReadPointer(t *testing.T) {
	testCases := []struct {
		input   string
		output  string
		wantErr bool
	}{{
		input:  ``,
		output: `[]`,
	}, {
		input:  `/foo`,
		output: `["foo"]`,
	}, {
		input:  `/foo/bar`,
		output: `["foo","bar"]`,
	}, {
		input:  `/foo/0`,
		output: `["foo",0]`,
	}, {
		input:  `/foo/0/bar`,
		output: `["foo",0,"bar"]`,
	}, {
		input:  `/0/foo`,
		output: `[0,"foo"]`,
	}, {
		input:  `/foo/-/bar`,
		output: `["foo",-1,"bar"]`,
	}}

	for _, tc := range testCases {
		got, err := readPointer(tc.input)
		if tc.wantErr && err == nil {
			t.Errorf("Wanted err. Got nil")
		}
		if !tc.wantErr && err != nil {
			t.Errorf("Wanted no err. Got %v", err)
		}
		want, _ := ReadJsonString(tc.output)
		if !got.JsonNode().Equals(want) {
			t.Errorf("Wanted %v. Got %v", tc.output, got)
		}
		back, err := writePointer(got.JsonNode().(jsonArray))
		if err != nil {
			t.Errorf("Wanted no err on back translation. Got %v", err)
		}
		if back != tc.input {
			t.Errorf("Wanted %q. Got %q", tc.input, back)
		}
	}
}
