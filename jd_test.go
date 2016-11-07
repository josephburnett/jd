package jd

import (
	"reflect"
	"testing"
)

func TestEquals(t *testing.T) {
	checkEqual(t, `{"a":1}`, `{"a":1}`)
	checkEqual(t, `{"a":1}`, `{"a":1.0}`)
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

func TestDiff(t *testing.T) {
	checkDiff(t, `{"a":1}`, `{"a":2}`,
		[]diffElement{
			diffElement{
				path:     path{"a"},
				oldValue: jsonNumber(1.0),
				newValue: jsonNumber(2.0),
			},
		})
}

func checkDiff(t *testing.T, a, b string, diff []diffElement) {
	jsonA, err := unmarshal([]byte(a))
	if err != nil {
		t.Error(err.Error())
	}
	jsonB, err := unmarshal([]byte(b))
	if err != nil {
		t.Error(err.Error())
	}
	path := make([]pathElement, 0)
	d := jsonA.diff(jsonB, path)
	if !reflect.DeepEqual(d, diff) {
		t.Errorf("!reflect.DeepEqual(d, diff). d was %v. diff was %v.", d, diff)
	}
}
