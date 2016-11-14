package jd

import (
	"testing"
)

func TestListEqual(t *testing.T) {
	checkEqual(t, `[]`, `[]`)
	checkEqual(t, `[1,2,3]`, `[1,2,3]`)
	checkEqual(t, `[[]]`, `[[]]`)
	checkEqual(t, `[{"a":1}]`, `[{"a":1}]`)
	checkEqual(t, `[{"a":[]}]`, `[{"a":[]}]`)
}

func TestListNotEqual(t *testing.T) {
	checkNotEqual(t, `[]`, `0`)
	checkNotEqual(t, `[]`, `{}`)
	checkNotEqual(t, `[]`, `[[]]`)
	checkNotEqual(t, `[1,2,3]`, `[3,2,1]`)
}
