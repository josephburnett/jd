package jd

import (
	"strconv"
	"testing"
)

var corpus = []string{
	``,  // void
	` `, // void
	`null`,
	`0`,
	`1`,
	`""`,
	`"foo"`,
	`"bar"`,
	`"null"`,
	`[]`,
	`[null]`,
	`[null,null,null]`,
	`[1]`,
	`[1,2,3]`,
	`[{},[],3]`,
	`[1,{},[]]`,
	`{}`,
	`{"foo":"bar"}`,
	`{"foo":null}`,
	`{"foo":1}`,
	`{"foo":[]}`,
	`{"foo":[null]}`,
	`{"foo":[1]}`,
	`{"foo":[1,2,3]}`,
	`{"foo":[1,null,3]}`,
	`{"foo":{}}`,
	`{"foo":{"bar":null}}`,
	`{"foo":{"bar":1}}`,
	`{"foo":{"bar":[]}}`,
	`{"foo":{"bar":[1,2,3]}}`,
	`{"foo":{"bar":{}}}`,
}

func FuzzJd(f *testing.F) {
	for _, a := range corpus {
		_, err := ReadJsonString(a)
		if err != nil {
			f.Errorf("corpus entry not valid JSON: %q", a)
		}
		for _, b := range corpus {
			f.Add(a, b)
		}
	}
	f.Fuzz(fuzz)
}

func fuzz(t *testing.T, aStr, bStr string) {
	// Only valid JSON input.
	a, err := ReadJsonString(aStr)
	if err != nil {
		return
	}
	if a == nil {
		t.Errorf("nil parsed input: %q", aStr)
		return
	}
	b, err := ReadJsonString(bStr)
	if err != nil {
		return
	}
	if b == nil {
		t.Errorf("nil parsed input: %q", bStr)
		return
	}
	for _, format := range [][2]string{{
		"jd", "list",
	}, {
		"jd", "set",
	}, {
		"jd", "mset",
	}, {
		"patch", "list",
	}, {
		"merge", "list",
	}} {
		a, _ = ReadJsonString(aStr) // Fresh parsed copy.
		if format[0] == "merge" {
			if hasUnsupportedNullValue(a) {
				continue
			}
			if hasUnsupportedNullValue(b) {
				continue
			}
			if b.Equals(jsonObject{}) {
				// An empty object is a JSON Merge patch noop
				continue
			}
		}
		var options []Option
		switch format[0] {
		case "jd":
			switch format[1] {
			case "set":
				options = append(options, setOption{})
			case "mset":
				options = append(options, multisetOption{})
			default: // list
			}
		case "merge":
			options = append(options, mergeOption{})
		default: // patch
		}

		// Diff A and B.
		d := a.Diff(b, options...)
		if d == nil {
			t.Errorf("nil diff of a and b")
			return
		}
		if format[0] == "patch" && hasUnsupportedObjectKey(d) {
			continue
		}
		var diffABStr string
		var diffAB Diff
		switch format[0] {
		case "jd":
			diffABStr = d.Render(options...)
			diffAB, err = ReadDiffString(diffABStr)
		case "patch":
			diffABStr, err = d.RenderPatch()
			if err != nil {
				t.Errorf("could not render diff %v as patch: %v", d, err)
				return
			}
			diffAB, err = ReadPatchString(diffABStr)
		case "merge":
			diffABStr, err = d.RenderMerge()
			if err != nil {
				t.Errorf("could not render diff %v as merge: %v", d, err)
				return
			}
			diffAB, err = ReadMergeString(diffABStr)
		}
		if err != nil {
			t.Errorf("error parsing diff string %q: %v", diffABStr, err)
			return
		}
		// Apply diff to A to get B.
		patchedA, err := a.Patch(diffAB)
		if err != nil {
			t.Errorf("applying patch %v to %v should give %v. Got err: %v", diffABStr, aStr, bStr, err)
			return
		}
		if !patchedA.Equals(b, options...) {
			t.Errorf("applying patch %v to %v should give %v. Got: %v", diffABStr, aStr, bStr, renderJson(patchedA))
			return
		}
	}

}

func hasUnsupportedObjectKey(diff Diff) bool {
	for _, d := range diff {
		for _, p := range d.Path {
			if s, ok := p.(PathKey); ok {
				// Object key that looks like number is interpretted incorrectly as array index.
				if _, err := strconv.Atoi(string(s)); err == nil {
					return true
				}
				// Object key "-" is interpretted incorrectly as append-to-array.
				if string(s) == "-" {
					return true
				}
			}
		}
	}
	return false
}

func hasUnsupportedNullValue(node JsonNode) bool {
	switch n := node.(type) {
	case jsonObject:
		for _, v := range n {
			if isNull(v) {
				return true
			}
			if hasUnsupportedNullValue(v) {
				return true
			}
		}
		return false
	case jsonArray, jsonList, jsonSet, jsonMultiset:
		for _, v := range n.(jsonArray) {
			if isNull(v) {
				return true
			}
			if hasUnsupportedNullValue(v) {
				return true
			}
		}
		return false
	case jsonNull:
		return true
	default:
		return false
	}
}
