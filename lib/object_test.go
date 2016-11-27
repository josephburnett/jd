package jd

import (
	"testing"
)

func TestObjectJson(t *testing.T) {
	checkJson(t, `{"a":1}`, `{"a":1}`)
	checkJson(t, ` { "a" : 1 } `, `{"a":1}`)
	checkJson(t, `{}`, `{}`)
}

func TestObjectEqual(t *testing.T) {
	checkEqual(t, `{"a":1}`, `{"a":1}`)
	checkEqual(t, `{"a":1}`, `{"a":1.0}`)
	checkEqual(t, `{"a":[1,2]}`, `{"a":[1,2]}`)
	checkEqual(t, `{"a":"b"}`, `{"a":"b"}`)
}

func TestObjectNotEqual(t *testing.T) {
	checkNotEqual(t, `{"a":1}`, `{"b":1}`)
	checkNotEqual(t, `{"a":[1,2]}`, `{"a":[2,1]}`)
	checkNotEqual(t, `{"a":"b"}`, `{"a":"c"}`)
}

func TestObjectHash(t *testing.T) {
	checkHash(t, `{}`, `{}`, true)
	checkHash(t, `{"a":1}`, `{"a":1}`, true)
	checkHash(t, `{"a":1}`, `{"a":2}`, false)
	checkHash(t, `{"a":1}`, `{"b":1}`, false)
	checkHash(t, `{"a":1,"b":2}`, `{"b":2,"a":1}`, true)
}

func TestObjectDiff(t *testing.T) {
	checkDiff(t, `{}`, `{}`)
	checkDiff(t, `{"a":1}`, `{"a":1}`)
	checkDiff(t, `{"a":1}`, `{"a":2}`,
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkDiff(t, `{"":1}`, `{"":1}`)
	checkDiff(t, `{"":1}`, `{"a":2}`,
		`@ [""]`,
		`- 1`,
		`@ ["a"]`,
		`+ 2`)
	checkDiff(t, `{"a":{"b":{}}}`, `{"a":{"b":{"c":1},"d":2}}`,
		`@ ["a","b","c"]`,
		`+ 1`,
		`@ ["a","d"]`,
		`+ 2`)
}

func testObjectPatch(t *testing.T) {
	checkPatch(t, `{}`, `{}`)
	checkPatch(t, `{"a":1}`, `{"a":1}`)
	checkPatch(t, `{"a":1}`, `{"a":2}`,
		`@ ["a"]`,
		`- 1`,
		`+ 2`)
	checkPatch(t, `{"":1}`, `{"":1}`)
	checkPatch(t, `{"":1}`, `{"a":2}`,
		`@ [""]`,
		`- 1`,
		`@ ["a"]`,
		`+ 2`)
	checkPatch(t, `{"a":{"b":{}}}`, `{"a":{"b":{"c":1},"d":2}}`,
		`@ ["a","b","c"]`,
		`+ 1`,
		`@ ["a","d"]`,
		`+ 2`)
}

func testObjectPatchError(t *testing.T) {
	checkPatchError(t, `{}`,
		`@ ["a"]`,
		`- 1`)
	checkPatchError(t, `{"a":1}`,
		`@ ["a"]`,
		`+ 2`)
	checkPatchError(t, `{"a":1}`,
		`@ ["a"]`,
		`+ 1`)
}
