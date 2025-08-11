package jd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// strPtr returns a pointer to a string literal
func strPtr(s string) *string {
	return &s
}

func TestOptionJSON(t *testing.T) {
	cases := []struct {
		json   string
		option Option
	}{{
		json:   `["MERGE"]`,
		option: MERGE,
	}, {
		json:   `["SET"]`,
		option: SET,
	}, {
		json:   `["MULTISET"]`,
		option: MULTISET,
	}, {
		json:   `["COLOR"]`,
		option: COLOR,
	}, {
		json:   `[{"precision":1.01}]`,
		option: Precision(1.01),
	}, {
		json:   `[{"setkeys":["foo","bar"]}]`,
		option: SetKeys("foo", "bar"),
	}, {
		json:   `[{"@":["foo"],"^":["SET"]}]`,
		option: PathOption(Path{PathKey("foo")}, SET),
	}, {
		json:   `[{"@":["foo"],"^":[{"@":["bar"],"^":["SET"]}]}]`,
		option: PathOption(Path{PathKey("foo")}, PathOption(Path{PathKey("bar")}, SET)),
	}}
	for _, c := range cases {
		t.Run(c.json, func(t *testing.T) {
			opts, err := ReadOptionsString(c.json)
			require.NoError(t, err)
			s, err := json.Marshal(opts)
			require.NoError(t, err)
			require.Equal(t, c.json, string(s))
			gotOpts := []any{c.option}
			s, err = json.Marshal(gotOpts)
			require.NoError(t, err)
			require.Equal(t, c.json, string(s))
		})
	}
}

func TestRefine(t *testing.T) {
	cases := []struct {
		name      string
		opts      []Option
		element   PathElement
		wantApply []Option
		wantRest  []Option
	}{{
		name:      "recurse applies and keeps gloval options",
		opts:      []Option{SET, COLOR},
		element:   PathKey("foo"),
		wantApply: []Option{SET, COLOR},
		wantRest:  []Option{SET, COLOR},
	}, {
		name:      "recurse consumes on element of path",
		opts:      []Option{PathOption(Path{PathKey("foo"), PathKey("bar")}, SET)},
		element:   PathKey("foo"),
		wantApply: nil,
		wantRest:  []Option{PathOption(Path{PathKey("bar")}, SET)},
	}, {
		name:      "path option ending in set",
		opts:      []Option{PathOption(Path{PathKey("foo"), PathSet{}})},
		element:   PathKey("foo"),
		wantApply: []Option{SET},
		wantRest:  nil,
	}, {
		name:      "path option delivering set in payload",
		opts:      []Option{PathOption(Path{PathKey("foo")}, SET)},
		element:   PathKey("foo"),
		wantApply: []Option{SET},
		wantRest:  nil,
	}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := refine(&options{retain: c.opts}, c.element)
			require.Equal(t, c.wantApply, o.apply)
			require.Equal(t, c.wantRest, o.retain)
		})
	}
}

func TestPathOption(t *testing.T) {
	cases := []struct {
		name         string
		opts         string
		a, b         string
		expectedDiff *string // Optional: if provided, expects this exact diff instead of no diff
	}{{
		name: "Precision on a number in an object",
		opts: `[{"@":["foo"],"^":[{"precision":0.1}]}]`,
		a:    `{"foo":1.0}`,
		b:    `{"foo":1.001}`,
	}, {
		name: "Precision on a number alone",
		opts: `[{"@":[],"^":[{"precision":0.1}]}]`,
		a:    `1.0`,
		b:    `1.001`,
	}, {
		name: "Precision in a list index",
		opts: `[{"@":[0],"^":[{"precision":0.1}]}]`,
		a:    `[1.0]`,
		b:    `[1.001]`,
	}, {
		// SET option tests - treating arrays as sets ignores order and duplicates
		name: "SET option on array in object",
		opts: `[{"@":["items"],"^":["SET"]}]`,
		a:    `{"items":[1,2,3]}`,
		b:    `{"items":[3,1,2]}`, // same elements, different order
	}, {
		name: "SET option on nested array",
		opts: `[{"@":["data","values"],"^":["SET"]}]`,
		a:    `{"data":{"values":[1,2,2,3]}}`,
		b:    `{"data":{"values":[3,2,1]}}`, // same unique elements, different order, no duplicates
	}, {
		name: "SET option on root array",
		opts: `[{"@":[],"^":["SET"]}]`,
		a:    `[1,2,2,3]`,
		b:    `[3,2,1]`, // same unique elements
	}, {
		// MULTISET option tests - treats arrays as multisets (order doesn't matter, duplicates count)
		name: "MULTISET option on array in object",
		opts: `[{"@":["tags"],"^":["MULTISET"]}]`,
		a:    `{"tags":["red","blue","red"]}`,
		b:    `{"tags":["red","red","blue"]}`, // same counts, different order
	}, {
		name: "MULTISET option on nested array",
		opts: `[{"@":["config","flags"],"^":["MULTISET"]}]`,
		a:    `{"config":{"flags":[1,1,2,3]}}`,
		b:    `{"config":{"flags":[2,1,3,1]}}`, // same element counts
	}, {
		// Precision with nested paths
		name: "Precision on nested object field",
		opts: `[{"@":["data","temperature"],"^":[{"precision":0.01}]}]`,
		a:    `{"data":{"temperature":20.123}}`,
		b:    `{"data":{"temperature":20.125}}`, // within 0.01 precision
	}, {
		name: "Precision on array element in nested structure",
		opts: `[{"@":["measurements",2],"^":[{"precision":0.05}]}]`,
		a:    `{"measurements":[10.0, 20.0, 30.42]}`,
		b:    `{"measurements":[10.0, 20.0, 30.44]}`, // index 2 within precision
	}, {
		// SetKeys option tests - for objects as sets with specific key matching
		name: "SetKeys option with id field",
		opts: `[{"@":["users"],"^":[{"setkeys":["id"]}]}]`,
		a:    `{"users":[{"id":"1","name":"Alice"},{"id":"2","name":"Bob"}]}`,
		b:    `{"users":[{"id":"2","name":"Bob"},{"id":"1","name":"Alice"}]}`, // same objects by id, different order
	}, {
		name: "SetKeys option with multiple keys",
		opts: `[{"@":["items"],"^":[{"setkeys":["type","id"]}]}]`,
		a:    `{"items":[{"type":"A","id":"1","data":"x"},{"type":"B","id":"2","data":"y"}]}`,
		b:    `{"items":[{"type":"B","id":"2","data":"y"},{"type":"A","id":"1","data":"x"}]}`, // same by type+id
	}, {
		// Complex nested path options
		name: "SET option on deeply nested array",
		opts: `[{"@":["level1","level2","items"],"^":["SET"]}]`,
		a:    `{"level1":{"level2":{"items":["a","b","c"]}}}`,
		b:    `{"level1":{"level2":{"items":["c","a","b"]}}}`, // same set elements
	}, {
		name: "Multiple path options - SET on one path, precision on another",
		opts: `[{"@":["coords"],"^":["SET"]},{"@":["value"],"^":[{"precision":0.1}]}]`,
		a:    `{"coords":[1,2,3],"value":5.01}`,
		b:    `{"coords":[3,1,2],"value":5.05}`, // coords as set, value within precision
	}, {
		// Mixed data structure tests
		name: "Precision in array within object array",
		opts: `[{"@":["items",0,"score"],"^":[{"precision":0.1}]}]`,
		a:    `{"items":[{"score":85.1},{"score":90.0}]}`,
		b:    `{"items":[{"score":85.15},{"score":90.0}]}`, // first item's score within precision
	}, {
		name: "SET option on array of objects with nested arrays",
		opts: `[{"@":["groups",0,"members"],"^":["SET"]}]`,
		a:    `{"groups":[{"members":["alice","bob","charlie"]}]}`,
		b:    `{"groups":[{"members":["bob","alice","charlie"]}]}`, // members as set
	}, {
		// Edge cases
		name: "Precision on number in array of arrays",
		opts: `[{"@":[1,0],"^":[{"precision":0.01}]}]`,
		a:    `[[1.0],[2.123],[3.0]]`,
		b:    `[[1.0],[2.125],[3.0]]`, // second array's first element within precision
	}, {
		name: "MULTISET option on specific nested array path",
		opts: `[{"@":["data","values"],"^":["MULTISET"]}]`,
		a:    `{"data":{"values":[1,2,1]}, "other":"same"}`,
		b:    `{"data":{"values":[2,1,1]}, "other":"same"}`, // same multiset counts, different order
	}, {
		// Test cases with expected diffs - mixing targeted and non-targeted paths
		name:         "SET option on one array path, list semantics preserved elsewhere",
		opts:         `[{"@":["setArray"],"^":["SET"]}]`,
		a:            `{"setArray":[1,2,3], "listArray":[1,2,3]}`,
		b:            `{"setArray":[3,1,2], "listArray":[3,1,2]}`, // same elements, different order
		expectedDiff: strPtr("@ [\"listArray\",0]\n[\n+ 3\n  1\n@ [\"listArray\",3]\n  2\n- 3\n]\n"),
	}, {
		name:         "Precision option on one number path, exact comparison elsewhere",
		opts:         `[{"@":["approx"],"^":[{"precision":0.1}]}]`,
		a:            `{"approx":1.05, "exact":1.05}`,
		b:            `{"approx":1.07, "exact":1.07}`, // within precision vs exact difference
		expectedDiff: strPtr("@ [\"exact\"]\n- 1.05\n+ 1.07\n"),
	}, {
		name:         "SET option on nested path, preserves list semantics at other paths",
		opts:         `[{"@":["data","tags"],"^":["SET"]}]`,
		a:            `{"data":{"tags":[1,2,3], "items":[1,2,3]}, "other":[1,2,3]}`,
		b:            `{"data":{"tags":[3,1,2], "items":[3,1,2]}, "other":[3,1,2]}`, // reordered everywhere
		expectedDiff: strPtr("@ [\"data\",\"items\",0]\n[\n+ 3\n  1\n@ [\"data\",\"items\",3]\n  2\n- 3\n]\n@ [\"other\",0]\n[\n+ 3\n  1\n@ [\"other\",3]\n  2\n- 3\n]\n"),
	}, {
		name:         "Multiple PathOptions - SET on one path, MULTISET on another, list preserved elsewhere",
		opts:         `[{"@":["sets"],"^":["SET"]}, {"@":["multisets"],"^":["MULTISET"]}]`,
		a:            `{"sets":[1,2,2,3], "multisets":[1,1,2], "lists":[1,2,3]}`,
		b:            `{"sets":[3,2,1], "multisets":[2,1,1], "lists":[3,2,1]}`, // same unique for sets, same counts for multisets, different order for lists
		expectedDiff: strPtr("@ [\"lists\",0]\n[\n+ 3\n+ 2\n  1\n@ [\"lists\",3]\n  1\n- 2\n- 3\n]\n"),
	}, {
		name:         "Precision on array element, exact comparison on other array elements",
		opts:         `[{"@":["values",0],"^":[{"precision":0.1}]}]`,
		a:            `{"values":[1.05, 2.05, 3.05]}`,
		b:            `{"values":[1.07, 2.07, 3.07]}`, // first within precision, others exact difference
		expectedDiff: strPtr("@ [\"values\",1]\n  1.05\n- 2.05\n- 3.05\n+ 2.07\n+ 3.07\n]\n"),
	}, {
		name:         "SET option at root doesn't affect nested lists",
		opts:         `[{"@":[],"^":["SET"]}]`,
		a:            `[{"nested":[1,2,3]}, {"nested":[4,5,6]}]`,
		b:            `[{"nested":[4,5,6]}, {"nested":[1,2,3]}]`, // root as set (reordered), but nested lists should remain as lists
		expectedDiff: strPtr(""),                                 // Should be empty since root array as set has same elements
	}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.expectedDiff == nil {
				// Original behavior: expect control diff to be non-empty, and options diff to be empty
				controlA, err := ReadJsonString(c.a)
				require.NoError(t, err)
				controlB, err := ReadJsonString(c.b)
				require.NoError(t, err)
				controlDiff := controlA.Diff(controlB)
				require.NotEmpty(t, controlDiff)

				a, err := ReadJsonString(c.a)
				require.NoError(t, err)
				b, err := ReadJsonString(c.b)
				require.NoError(t, err)
				o, err := ReadOptionsString(c.opts)
				require.NoError(t, err)
				diff := a.Diff(b, o...)
				require.Empty(t, diff)
			} else {
				// New behavior: expect specific diff
				a, err := ReadJsonString(c.a)
				require.NoError(t, err)
				b, err := ReadJsonString(c.b)
				require.NoError(t, err)
				o, err := ReadOptionsString(c.opts)
				require.NoError(t, err)
				actualDiff := a.Diff(b, o...)

				// Compare jd diff format representations
				actualDiffString := actualDiff.Render()
				require.Equal(t, *c.expectedDiff, actualDiffString, "Diff should match expected")
			}
		})
	}
}
