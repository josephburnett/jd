package jd

import (
	"reflect"
	"testing"
)

func TestEquals(t *testing.T) {
	checkEqual(t, `{"a":1}`, `{"a":1}`)
	checkEqual(t, `{"a":1}`, `{"a":1.0}`)
	checkEqual(t, `{"a":[1,2]}`, `{"a":[1,2]}`)
	checkEqual(t, `{"a":"b"}`, `{"a":"b"}`)
	checkNotEqual(t, `{"a":1}`, `{"b":1}`)
	checkNotEqual(t, `{"a":[1,2]}`, `{"a":[2,1]}`)
	checkNotEqual(t, `{"a":"b"}`, `{"a":"c"}`)
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

func TestDiff(t *testing.T) {
	checkDiff(t, `{"a":1}`, `{"a":2}`,
		Diff{DiffElement{Path{"a"}, jsonNumber(1.0), jsonNumber(2.0)}})
	checkDiff(t, `{"a":1}`, `{}`,
		Diff{DiffElement{Path{"a"}, jsonNumber(1.0), nil}})
	checkDiff(t, `{}`, `{"a":2}`,
		Diff{DiffElement{Path{"a"}, nil, jsonNumber(2.0)}})
	checkDiff(t, `{"a":1}`, `{"a":1}`, Diff{})
	checkDiff(t, `{"a":{"b":1}}`, `{"a":{"c":2}}`,
		Diff{
			DiffElement{Path{"a", "b"}, jsonNumber(1.0), nil},
			DiffElement{Path{"a", "c"}, nil, jsonNumber(2.0)}})
	checkDiff(t, `{"a":[1,2]}`, `{"a":[2,1]}`,
		Diff{
			DiffElement{Path{"a", 1}, jsonNumber(2.0), jsonNumber(1.0)},
			DiffElement{Path{"a", 0}, jsonNumber(1.0), jsonNumber(2.0)}})
	checkDiff(t, `{"a":[1]}`, `{"a":[1,2]}`,
		Diff{DiffElement{Path{"a", 1}, nil, jsonNumber(2.0)}})
	checkDiff(t, `{"a":[1,2]}`, `{"a":[1]}`,
		Diff{DiffElement{Path{"a", 1}, jsonNumber(2.0), nil}})
	checkDiff(t, `{"a":[1,2]}`, `{"a":[3,4]}`,
		Diff{
			DiffElement{Path{"a", 1}, jsonNumber(2.0), jsonNumber(4.0)},
			DiffElement{Path{"a", 0}, jsonNumber(1.0), jsonNumber(3.0)}})
	checkDiff(t, `{"a":[{"b":1}]}`, `{"a":[{"b":2}]}`,
		Diff{DiffElement{Path{"a", 0, "b"}, jsonNumber(1.0), jsonNumber(2.0)}})
}

func checkDiff(t *testing.T, a, b string, diff Diff) {
	jsonA, err := unmarshal([]byte(a))
	if err != nil {
		t.Error(err.Error())
	}
	jsonB, err := unmarshal([]byte(b))
	if err != nil {
		t.Error(err.Error())
	}
	path := make(Path, 0)
	d := jsonA.diff(jsonB, path)
	if !reflect.DeepEqual(d, diff) {
		t.Errorf("Got %v. Want %v.", d, diff)
	}
}
