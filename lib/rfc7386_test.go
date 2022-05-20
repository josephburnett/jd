package jd

import "testing"

func TestRfc7386AppendixA(t *testing.T) {
	exampleTestCases := []struct {
		original string
		patch    string
		result   string
	}{{
		original: `{"a":"b"}`,
		patch:    `{"a":"c"}`,
		result:   `{"a":"c"}`,
	}, {
		original: `{"a":"b"}`,
		patch:    `{"b":"c"}`,
		result:   `{"a":"b","b":"c"}`,
	}, {
		original: `{"a":"b"}`,
		patch:    `{"a":null}`,
		result:   `{}`,
	}, {
		original: `{"a":"b","b":"c"}`,
		patch:    `{"a":null}`,
		result:   `{"b":"c"}`,
	}, {
		original: `{"a":["b"]}`,
		patch:    `{"a":"c"}`,
		result:   `{"a":"c"}`,
	}, {
		original: `{"a":"c"}`,
		patch:    `{"a":["b"]}`,
		result:   `{"a":["b"]}`,
	}, {
		original: `{"a":{"b":"c"}}`,
		patch:    `{"a":{"b":"d","c":null}}`,
		result:   `{"a": {"b": "d"}}`,
	}, {
		original: `{"a": [{"b":"c"}]}`,
		patch:    `{"a": [1]}`,
		result:   `{"a": [1]}`,
	}, {
		original: `["a","b"]`,
		patch:    `["c","d"]`,
		result:   `["c","d"]`,
	}, {
		original: `{"a":"b"}`,
		patch:    `["c"]`,
		result:   `["c"]`,
	}, {
		original: `{"a":"foo"}`,
		patch:    `null`,
		result:   `null`,
	}, {
		original: `{"a":"foo"}`,
		patch:    `"bar"`,
		result:   `"bar"`,
	}, {
		original: `{"e":null}`,
		patch:    `{"a":1}`,
		result:   `{"e":null,"a":1}`,
	}, {
		original: `[1,2]`,
		patch:    `{"a":"b","c":null}`,
		result:   `{"a":"b"}`,
	}, {
		original: `{}`,
		patch:    `{"a":{"bb":{"ccc":null}}}`,
		result:   `{"a":{"bb":{}}}`,
	}}

	for _, tc := range exampleTestCases {
		diff, err := ReadMergeString(tc.patch)
		if err != nil {
			t.Errorf("Error reading patch: %v", err)
		}
		original, err := ReadJsonString(tc.original)
		if err != nil {
			t.Errorf("Error reading original: %v", err)
		}
		_, err = original.Patch(diff)
		if err != nil {
			t.Errorf("Error patching %v with %v (%v): %v", tc.original, tc.patch, diff.Render(), err)
		}

	}
}
