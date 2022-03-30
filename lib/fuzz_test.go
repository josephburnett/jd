// +build test_fuzz

package jd

import (
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
