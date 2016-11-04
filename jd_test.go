package jd

import (
	"testing"
)

func testEquals(t *testing.T) {
	checkEqual(t, `{"a":1}`, `{"a":1}`)
	checkEqual(t, `{"a":1}`, `{"a":1.0`)
	checkEqual(t, `{"a":[1,2]}`, `{"a":[1,2]}`)
	checkNotEqual(t, `{"a":1}`, `{"b":1}`)
	checkNotEqual(t, `{"a":[1,2]}`, `{"a":[2,1]}`)
}

func checkEqual(t *testing.T, a, b string) {
	nodeA, err := unmarshal([]byte(a))
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b))
	if err != nil {
		t.Errorf(err.Error())
	}
	if !nodeA.equals(nodeB) {
		t.Errorf("nodeA.equals(nodeB) == false. Want true.")
	}
	if !nodeB.equals(nodeA) {
		t.Errorf("nodeB.equals(nodeA) == false. Want true.")
	}
	if !nodeA.equals(nodeA) {
		t.Errorf("nodeA.equals(nodeA) == false. Want true.")
	}
	if !nodeB.equals(nodeB) {
		t.Errorf("nodeB.equals(nodeB) == false. Want true.")
	}
}

func checkNotEqual(t *testing.T, a, b string) {
	nodeA, err := unmarshal([]byte(a))
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b))
	if err != nil {
		t.Errorf(err.Error())
	}
	if nodeA.equals(nodeB) {
		t.Errorf("nodeA.equals(nodeB) == true. Want false.")
	}
	if nodeB.equals(nodeA) {
		t.Errorf("nodeB.equals(nodeA) == true. Want false.")
	}
}
