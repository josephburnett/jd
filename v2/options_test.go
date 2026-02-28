package jd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadOptionsStringErrors(t *testing.T) {
	// Invalid JSON
	_, err := ReadOptionsString("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	// Not an array
	_, err = ReadOptionsString(`"string"`)
	if err == nil {
		t.Fatal("expected error for non-array")
	}
	// Unknown option in array
	_, err = ReadOptionsString(`["UNKNOWN"]`)
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestNewOptionEdgeCases(t *testing.T) {
	// Unrecognized string
	_, err := NewOption("UNKNOWN")
	if err == nil {
		t.Fatal("expected error")
	}
	// precision with wrong type
	_, err = NewOption(map[string]any{"precision": "not a number"})
	if err == nil {
		t.Fatal("expected error")
	}
	// keys with wrong type
	_, err = NewOption(map[string]any{"keys": "not an array"})
	if err == nil {
		t.Fatal("expected error")
	}
	// keys with non-string element
	_, err = NewOption(map[string]any{"keys": []any{42}})
	if err == nil {
		t.Fatal("expected error")
	}
	// setkeys backward compat
	opt, err := NewOption(map[string]any{"setkeys": []any{"id"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := opt.(setKeysOption); !ok {
		t.Errorf("expected setKeysOption, got %T", opt)
	}
	// Merge: true
	opt, err = NewOption(map[string]any{"Merge": true})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := opt.(mergeOption); !ok {
		t.Errorf("expected mergeOption, got %T", opt)
	}
	// Merge: false
	_, err = NewOption(map[string]any{"Merge": false})
	if err == nil {
		t.Fatal("expected error for Merge:false")
	}
	// Merge: wrong type
	_, err = NewOption(map[string]any{"Merge": "yes"})
	if err == nil {
		t.Fatal("expected error for Merge with wrong type")
	}
	// Unknown single-key object
	_, err = NewOption(map[string]any{"unknown": 1})
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
	// 2-key object with bad ^ type
	_, err = NewOption(map[string]any{"@": []any{"a"}, "^": "not array"})
	if err == nil {
		t.Fatal("expected error")
	}
	// 2-key object with bad option in ^
	_, err = NewOption(map[string]any{"@": []any{"a"}, "^": []any{"UNKNOWN"}})
	if err == nil {
		t.Fatal("expected error")
	}
	// 2-key object with unknown key
	_, err = NewOption(map[string]any{"@": []any{"a"}, "x": "y"})
	if err == nil {
		t.Fatal("expected error for unknown key in 2-key object")
	}
	// 3-key object
	_, err = NewOption(map[string]any{"a": 1, "b": 2, "c": 3})
	if err == nil {
		t.Fatal("expected error for 3-key object")
	}
	// Unsupported base type
	_, err = NewOption(42.0)
	if err == nil {
		t.Fatal("expected error for float64 type")
	}
	// Valid 2-key PathOption
	opt, err = NewOption(map[string]any{"@": []any{"users"}, "^": []any{"SET"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := opt.(pathOption); !ok {
		t.Errorf("expected pathOption, got %T", opt)
	}
}

func TestRefineDiffOnOff(t *testing.T) {
	// Global DIFF_OFF
	opts := newOptions([]Option{DIFF_OFF})
	refined := refine(opts, PathKey("anything"))
	if refined.diffingOn {
		t.Error("DIFF_OFF should set diffingOn to false")
	}
	// Global DIFF_ON
	opts = newOptions([]Option{DIFF_OFF, DIFF_ON})
	refined = refine(opts, PathKey("anything"))
	if !refined.diffingOn {
		t.Error("DIFF_ON should set diffingOn back to true")
	}
	// PathOption with DIFF_OFF in Then
	opts = newOptions([]Option{PathOption(Path{PathKey("a")}, DIFF_OFF)})
	refined = refine(opts, PathKey("a"))
	if refined.diffingOn {
		t.Error("PathOption DIFF_OFF should set diffingOn to false")
	}
	// PathOption with DIFF_ON in Then
	opts = newOptions([]Option{DIFF_OFF, PathOption(Path{PathKey("a")}, DIFF_ON)})
	refined = refine(opts, PathKey("a"))
	if !refined.diffingOn {
		t.Error("PathOption DIFF_ON should set diffingOn back to true")
	}
}

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
		json:   `[{"keys":["foo","bar"]}]`,
		option: SetKeys("foo", "bar"),
	}, {
		json:   `[{"@":["foo"],"^":["SET"]}]`,
		option: PathOption(Path{PathKey("foo")}, SET),
	}, {
		json:   `[{"@":["foo"],"^":[{"@":["bar"],"^":["SET"]}]}]`,
		option: PathOption(Path{PathKey("foo")}, PathOption(Path{PathKey("bar")}, SET)),
	}, {
		json:   `["COLOR_WORDS"]`,
		option: COLOR_WORDS,
	}, {
		json:   `["DIFF_ON"]`,
		option: DIFF_ON,
	}, {
		json:   `["DIFF_OFF"]`,
		option: DIFF_OFF,
	}, {
		json:   `[{"file":"example.json"}]`,
		option: File("example.json"),
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
		name: "Keys option with id field",
		opts: `[{"@":["users"],"^":[{"keys":["id"]}]}]`,
		a:    `{"users":[{"id":"1","name":"Alice"},{"id":"2","name":"Bob"}]}`,
		b:    `{"users":[{"id":"2","name":"Bob"},{"id":"1","name":"Alice"}]}`, // same objects by id, different order
	}, {
		name: "Keys option with multiple keys",
		opts: `[{"@":["items"],"^":[{"keys":["type","id"]}]}]`,
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

func TestDiffOnOffOption(t *testing.T) {
	cases := []struct {
		name         string
		opts         string
		a, b         string
		expectEmpty  bool    // true if diff should be empty
		expectedDiff *string // if specified, check exact diff output
	}{{
		name:        "DIFF_OFF at root ignores all changes",
		opts:        `[{"@":[],"^":["DIFF_OFF"]}]`,
		a:           `{"foo":1,"bar":"hello"}`,
		b:           `{"foo":2,"bar":"world","baz":true}`,
		expectEmpty: true,
	}, {
		name:        "DIFF_ON at root allows all changes (default behavior)",
		opts:        `[{"@":[],"^":["DIFF_ON"]}]`,
		a:           `{"foo":1}`,
		b:           `{"foo":2}`,
		expectEmpty: false,
	}, {
		name:         "DIFF_OFF for specific field ignores that field only",
		opts:         `[{"@":["timestamp"],"^":["DIFF_OFF"]}]`,
		a:            `{"data":"hello","timestamp":"2023-01-01"}`,
		b:            `{"data":"world","timestamp":"2023-01-02"}`,
		expectedDiff: strPtr("@ [\"data\"]\n- \"hello\"\n+ \"world\"\n"),
	}, {
		name:         "Allow-list approach: DIFF_OFF at root, DIFF_ON for specific paths",
		opts:         `[{"@":[],"^":["DIFF_OFF"]}, {"@":["userdata"],"^":["DIFF_ON"]}]`,
		a:            `{"userdata":"important","system":"ignore1","timestamp":"2023-01-01"}`,
		b:            `{"userdata":"changed","system":"ignore2","timestamp":"2023-01-02"}`,
		expectedDiff: strPtr("@ [\"userdata\"]\n- \"important\"\n+ \"changed\"\n"),
	}, {
		name:         "Nested PathOptions: override parent state",
		opts:         `[{"@":["config"],"^":["DIFF_OFF"]}, {"@":["config","user_settings"],"^":["DIFF_ON"]}]`,
		a:            `{"config":{"system":"val1","user_settings":"setting1"}}`,
		b:            `{"config":{"system":"val2","user_settings":"setting2"}}`,
		expectedDiff: strPtr("@ [\"config\",\"user_settings\"]\n- \"setting1\"\n+ \"setting2\"\n"),
	}, {
		name:        "DIFF_OFF with arrays - no diff generated",
		opts:        `[{"@":["tags"],"^":["DIFF_OFF"]}]`,
		a:           `{"tags":[1,2,3],"data":"same"}`,
		b:           `{"tags":[4,5,6],"data":"same"}`,
		expectEmpty: true,
	}, {
		name:        "DIFF_OFF with objects - no diff generated",
		opts:        `[{"@":["metadata"],"^":["DIFF_OFF"]}]`,
		a:           `{"metadata":{"id":1,"created":"2023"},"value":10}`,
		b:           `{"metadata":{"id":2,"created":"2024"},"value":10}`,
		expectEmpty: true,
	}, {
		name:        "DIFF_OFF combined with SET option",
		opts:        `[{"@":["ignored"],"^":["DIFF_OFF"]}, {"@":["tags"],"^":["SET"]}]`,
		a:           `{"ignored":[1,2],"tags":[1,2,3]}`,
		b:           `{"ignored":[3,4],"tags":[3,1,2]}`, // tags reordered but same set
		expectEmpty: true,
	}, {
		name:         "Multiple DIFF_OFF paths",
		opts:         `[{"@":["field1"],"^":["DIFF_OFF"]}, {"@":["field2"],"^":["DIFF_OFF"]}]`,
		a:            `{"field1":"change1","field2":"change2","field3":"old"}`,
		b:            `{"field1":"new1","field2":"new2","field3":"new"}`,
		expectedDiff: strPtr("@ [\"field3\"]\n- \"old\"\n+ \"new\"\n"),
	}, {
		name:        "Order precedence: last DIFF_OFF wins",
		opts:        `[{"@":["test"],"^":["DIFF_ON"]}, {"@":["test"],"^":["DIFF_OFF"]}]`,
		a:           `{"test":"old","other":"same"}`,
		b:           `{"test":"new","other":"same"}`,
		expectEmpty: true,
	}, {
		name:         "Order precedence: last DIFF_ON wins",
		opts:         `[{"@":["test"],"^":["DIFF_OFF"]}, {"@":["test"],"^":["DIFF_ON"]}]`,
		a:            `{"test":"old","other":"same"}`,
		b:            `{"test":"new","other":"same"}`,
		expectedDiff: strPtr("@ [\"test\"]\n- \"old\"\n+ \"new\"\n"),
	}, {
		name:         "Deep nesting with selective diffing",
		opts:         `[{"@":["level1","level2"],"^":["DIFF_OFF"]}, {"@":["level1","level2","important"],"^":["DIFF_ON"]}]`,
		a:            `{"level1":{"level2":{"ignore":"val1","important":"data1"}}}`,
		b:            `{"level1":{"level2":{"ignore":"val2","important":"data2"}}}`,
		expectedDiff: strPtr("@ [\"level1\",\"level2\",\"important\"]\n- \"data1\"\n+ \"data2\"\n"),
	}, {
		name:        "DIFF_OFF on array elements",
		opts:        `[{"@":[0],"^":["DIFF_OFF"]}]`,
		a:           `["ignore","keep"]`,
		b:           `["changed","keep"]`,
		expectEmpty: true,
	}, {
		name:        "DIFF_OFF with primitive types",
		opts:        `[{"@":[],"^":["DIFF_OFF"]}]`,
		a:           `"old string"`,
		b:           `"new string"`,
		expectEmpty: true,
	}, {
		name:        "DIFF_OFF with numbers",
		opts:        `[{"@":[],"^":["DIFF_OFF"]}]`,
		a:           `42`,
		b:           `100`,
		expectEmpty: true,
	}, {
		name:        "DIFF_OFF with booleans",
		opts:        `[{"@":[],"^":["DIFF_OFF"]}]`,
		a:           `true`,
		b:           `false`,
		expectEmpty: true,
	}, {
		name:        "DIFF_OFF with null values",
		opts:        `[{"@":[],"^":["DIFF_OFF"]}]`,
		a:           `null`,
		b:           `"not null"`,
		expectEmpty: true,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a, err := ReadJsonString(c.a)
			require.NoError(t, err)
			b, err := ReadJsonString(c.b)
			require.NoError(t, err)
			opts, err := ReadOptionsString(c.opts)
			require.NoError(t, err)

			diff := a.Diff(b, opts...)

			if c.expectEmpty {
				require.Empty(t, diff, "Expected empty diff but got: %s", diff.Render())
			} else if c.expectedDiff != nil {
				actualDiffString := diff.Render()
				require.Equal(t, *c.expectedDiff, actualDiffString, "Diff should match expected")
			} else {
				require.NotEmpty(t, diff, "Expected non-empty diff")
			}
		})
	}
}

func TestRefinePathMultiset(t *testing.T) {
	// PathOption ending in PathMultiset should infer MULTISET
	opts := newOptions([]Option{PathOption(Path{PathKey("foo"), PathMultiset{}})})
	refined := refine(opts, PathKey("foo"))
	_, found := getOption[multisetOption](refined)
	if !found {
		t.Error("expected multisetOption to be applied for PathMultiset")
	}
}

func TestFileOption(t *testing.T) {
	// NewOption recognizes {"file":"example.json"}
	opt, err := NewOption(map[string]any{"file": "example.json"})
	require.NoError(t, err)
	fo, ok := opt.(fileOption)
	require.True(t, ok)
	require.Equal(t, "example.json", fo.file)

	// MarshalJSON produces {"file":"example.json"}
	b, err := json.Marshal(opt)
	require.NoError(t, err)
	require.Equal(t, `{"file":"example.json"}`, string(b))

	// Round-trip: File() -> marshal -> unmarshal via NewOption -> equal
	opt2 := File("a.json")
	b2, err := json.Marshal(opt2)
	require.NoError(t, err)
	var raw any
	err = json.Unmarshal(b2, &raw)
	require.NoError(t, err)
	opt3, err := NewOption(raw)
	require.NoError(t, err)
	fo3, ok := opt3.(fileOption)
	require.True(t, ok)
	require.Equal(t, "a.json", fo3.file)

	// Wrong type for file value
	_, err = NewOption(map[string]any{"file": 42})
	require.Error(t, err)
}

func TestValidateOptions(t *testing.T) {
	cases := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{{
		name: "SET alone is valid",
		opts: []Option{SET},
	}, {
		name: "precision alone is valid",
		opts: []Option{Precision(0.01)},
	}, {
		name: "SET with SetKeys is valid",
		opts: []Option{SET, SetKeys("id")},
	}, {
		name: "precision with SetKeys is valid",
		opts: []Option{Precision(0.01), SetKeys("id")},
	}, {
		name:    "SET with precision is invalid",
		opts:    []Option{SET, Precision(0.01)},
		wantErr: true,
	}, {
		name:    "MULTISET with precision is invalid",
		opts:    []Option{MULTISET, Precision(0.01)},
		wantErr: true,
	}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateOptions(c.opts)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRefineEmptyAt(t *testing.T) {
	// PathOption with empty At and non-nil path element should be skipped
	opts := newOptions([]Option{PathOption(Path{}, SET)})
	refined := refine(opts, PathKey("foo"))
	if checkOption[setOption](refined) {
		t.Error("expected setOption to NOT be applied for empty At with non-nil path")
	}
}
