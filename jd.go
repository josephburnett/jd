package jd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

type jsonNode interface {
	diff(b jsonNode, path path) []diffElement
	equals(b jsonNode) bool
}

// type jsonList []interface{}
type jsonNumber float64

func (n1 jsonNumber) equals(n jsonNode) bool {
	n2, ok := n.(jsonNumber)
	if !ok {
		return false
	}
	if n1 != n2 {
		return false
	}
	return true
}

func (n1 jsonNumber) diff(n jsonNode, path path) []diffElement {
	d := make([]diffElement, 0)
	if n1.equals(n) {
		return d
	}
	e := diffElement{
		path:     path.clone(),
		oldValue: n1,
		newValue: n,
	}
	return append(d, e)
}

type jsonStruct map[string]interface{}

func (s1 jsonStruct) equals(n jsonNode) bool {
	s2, ok := n.(jsonStruct)
	if !ok {
		return false
	}
	if len(s1) != len(s2) {
		return false
	}
	return reflect.DeepEqual(s1, s2)
}

func newJsonNode(n interface{}) jsonNode {
	switch t := n.(type) {
	case map[string]interface{}:
		return jsonStruct(t)
	// case []interface{}:
	// 	return jsonList(t)
	case float64:
		return jsonNumber(t)
	default:
		panic(fmt.Sprintf("Unexpected type %v", t))
	}
}

func (s1 jsonStruct) diff(n jsonNode, path path) []diffElement {
	d := make([]diffElement, 0)
	s2, ok := n.(jsonStruct)
	if !ok {
		// Different types
		e := diffElement{
			path:     path.clone(),
			oldValue: s1,
			newValue: n,
		}
		return append(d, e)
	}
	for k1, n1 := range s1 {
		v1 := newJsonNode(n1)
		if n2, ok := s2[k1]; ok {
			// Both keys are present
			v2 := newJsonNode(n2)
			subDiff := v1.diff(v2, append(path.clone(), k1))
			d = append(d, subDiff...)
		} else {
			// S2 missing key
			e := diffElement{
				path:     append(path.clone(), k1),
				oldValue: v1,
				newValue: nil,
			}
			d = append(d, e)
		}
	}
	for k2, n2 := range s2 {
		v2 := newJsonNode(n2)
		if _, ok := s1[k2]; !ok {
			// S1 missing key
			e := diffElement{
				path:     append(path.clone(), k2),
				oldValue: nil,
				newValue: v2,
			}
			d = append(d, e)
		}
	}
	return d
}

type pathElement interface{}
type path []pathElement

type diffElement struct {
	path     path
	oldValue jsonNode
	newValue jsonNode
}

func (p1 path) clone() path {
	p2 := make(path, len(p1), len(p1)+1)
	copy(p2, p1)
	return p2
}

func readFile(filename string) (jsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes)
}

func unmarshal(bytes []byte) (jsonNode, error) {
	node := make(jsonStruct)
	err := json.Unmarshal(bytes, &node)
	if err != nil {
		return nil, err
	}
	return node, nil
}
