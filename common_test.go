package jd

import (
	"testing"
)

func checkEqual(t *testing.T, a, b string) {
	nodeA, err := unmarshal([]byte(a))
	if err != nil {
		t.Errorf(err.Error())
	}
	nodeB, err := unmarshal([]byte(b))
	if err != nil {
		t.Errorf(err.Error())
	}
	if !nodeA.Equals(nodeB) {
		t.Errorf("nodeA.Equals(nodeB) == false. Want true.")
	}
	if !nodeB.Equals(nodeA) {
		t.Errorf("nodeB.Equals(nodeA) == false. Want true.")
	}
	if !nodeA.Equals(nodeA) {
		t.Errorf("nodeA.Equals(nodeA) == false. Want true.")
	}
	if !nodeB.Equals(nodeB) {
		t.Errorf("nodeB.Equals(nodeB) == false. Want true.")
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
	if nodeA.Equals(nodeB) {
		t.Errorf("nodeA.Equals(nodeB) == true. Want false.")
	}
	if nodeB.Equals(nodeA) {
		t.Errorf("nodeB.Equals(nodeA) == true. Want false.")
	}
}
