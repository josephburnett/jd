package jd

import (
	"testing"
)

func TestEqual(t *testing.T) {
	checkEqual(t, `{"a":1}`, `{"a":1}`)
	checkEqual(t, `{"a":1}`, `{"a":1.0}`)
	checkEqual(t, `{"a":[1,2]}`, `{"a":[1,2]}`)
	checkEqual(t, `{"a":"b"}`, `{"a":"b"}`)
}

func TestNotEqual(t *testing.T) {
	checkNotEqual(t, `{"a":1}`, `{"b":1}`)
	checkNotEqual(t, `{"a":[1,2]}`, `{"a":[2,1]}`)
	checkNotEqual(t, `{"a":"b"}`, `{"a":"c"}`)
}
