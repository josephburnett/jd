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
}

func testObjectPatch(t *testing.T) {

}

func testObjectPatchError(t *testing.T) {

}
