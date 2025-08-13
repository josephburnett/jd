package jd

import "testing"

func TestReadPointer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		output  string
		wantErr bool
	}{
		{
			name:   "empty pointer path",
			input:  ``,
			output: `[]`,
		},
		{
			name:   "single string key",
			input:  `/foo`,
			output: `["foo"]`,
		},
		{
			name:   "nested string keys",
			input:  `/foo/bar`,
			output: `["foo","bar"]`,
		},
		{
			name:   "string key with numeric index",
			input:  `/foo/0`,
			output: `["foo",0]`,
		},
		{
			name:   "mixed string and numeric path",
			input:  `/foo/0/bar`,
			output: `["foo",0,"bar"]`,
		},
		{
			name:   "numeric index first",
			input:  `/0/foo`,
			output: `[0,"foo"]`,
		},
		{
			name:   "array append indicator",
			input:  `/foo/-/bar`,
			output: `["foo",-1,"bar"]`,
		},
		{
			name:   "multiple numeric indices",
			input:  `/0/1/2`,
			output: `[0,1,2]`,
		},
		{
			name:   "deep nested path",
			input:  `/a/b/c/d/e`,
			output: `["a","b","c","d","e"]`,
		},
		{
			name:   "numeric with append indicator",
			input:  `/0/-`,
			output: `[0,-1]`,
		},
		{
			name:   "string with special characters",
			input:  `/foo~1bar`,
			output: `["foo/bar"]`,
		},
		{
			name:   "string with tilde escape",
			input:  `/foo~0bar`,
			output: `["foo~bar"]`,
		},
		{
			name:   "empty string key",
			input:  `//foo`,
			output: `["","foo"]`,
		},
		{
			name:   "key with spaces",
			input:  `/hello world`,
			output: `["hello world"]`,
		},
		{
			name:   "unicode key",
			input:  `/héllo`,
			output: `["héllo"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPointer(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("Wanted err. Got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Wanted no err. Got %v", err)
			}
			want, _ := ReadJsonString(tt.output)
			if !got.JsonNode().Equals(want) {
				t.Errorf("Wanted %v. Got %v", tt.output, got)
			}
			back, err := writePointer(got.JsonNode().(jsonArray))
			if err != nil {
				t.Errorf("Wanted no err on back translation. Got %v", err)
			}
			if back != tt.input {
				t.Errorf("Wanted %q. Got %q", tt.input, back)
			}
		})
	}
}
