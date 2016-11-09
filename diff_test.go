package jd

import (
	"testing"
)

func TestDiffSimple(t *testing.T) {
	checkDiffRender(t, `{"a":1}`, `{"a":2}`, "@ [\"a\"]\n- 1\n+ 2\n")
}

func checkDiffRender(t *testing.T, a, b, diff string) {
	aJson, err := unmarshal([]byte(a))
	if err != nil {
		t.Errorf(err.Error())
	}
	bJson, err := unmarshal([]byte(b))
	if err != nil {
		t.Errorf(err.Error())
	}
	d, err := aJson.diff(bJson, Path{}).Render()
	if err != nil {
		t.Errorf(err.Error())
	}
	if d != diff {
		t.Errorf("%v.diff(%v) = %v. Want %v.", a, b, d, diff)
	}
}
