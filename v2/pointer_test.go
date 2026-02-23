package jd

import (
	"strings"
	"testing"
)

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

func TestWritePointerEdgeCases(t *testing.T) {
	// -1 index -> "-"
	path := jsonArray{jsonNumber(-1)}
	s, err := writePointer(path)
	if err != nil {
		t.Fatal(err)
	}
	if s != "/-" {
		t.Errorf("expected /-. got %v", s)
	}
	// Numeric string key -> error
	path = jsonArray{jsonString("123")}
	_, err = writePointer(path)
	if err == nil {
		t.Fatal("expected error for numeric string key")
	}
	// "-" string key -> error
	path = jsonArray{jsonString("-")}
	_, err = writePointer(path)
	if err == nil {
		t.Fatal("expected error for dash key")
	}
	// jsonArray in path -> error
	path = jsonArray{jsonArray{}}
	_, err = writePointer(path)
	if err == nil {
		t.Fatal("expected error for array in path")
	}
}

func TestWritePointerSetPathError(t *testing.T) {
	// Set-based paths contain jsonObject elements which cannot be
	// expressed as JSON Pointers (RFC 6901).
	path := jsonArray{
		jsonString("moves"),
		jsonObject{
			"move": jsonObject{"name": jsonString("mimic")},
		},
	}
	_, err := writePointer(path)
	if err == nil {
		t.Fatal("Wanted err for set-based path. Got nil")
	}
	if !strings.Contains(err.Error(), "set-based paths") {
		t.Errorf("Wanted set-based path error. Got %v", err)
	}
}

func TestWritePointerDefaultType(t *testing.T) {
	// unsupported type like jsonBool should error
	_, err := writePointer(jsonArray{jsonBool(true)})
	if err == nil {
		t.Fatal("expected error for unsupported type in pointer path")
	}
}
