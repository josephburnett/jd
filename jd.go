package jd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

type JsonNode interface {
	diff(b JsonNode, path Path) Diff
	equals(b JsonNode) bool
}

// type jsonList []interface{}
type JsonNumber float64

func (n1 JsonNumber) equals(n JsonNode) bool {
	n2, ok := n.(JsonNumber)
	if !ok {
		return false
	}
	if n1 != n2 {
		return false
	}
	return true
}

func (n1 JsonNumber) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	if n1.equals(n) {
		return d
	}
	e := DiffElement{
		Path:     path.clone(),
		OldValue: n1,
		NewValue: n,
	}
	return append(d, e)
}

type JsonStruct map[string]interface{}

func (s1 JsonStruct) equals(n JsonNode) bool {
	s2, ok := n.(JsonStruct)
	if !ok {
		return false
	}
	if len(s1) != len(s2) {
		return false
	}
	return reflect.DeepEqual(s1, s2)
}

func newJsonNode(n interface{}) JsonNode {
	switch t := n.(type) {
	case map[string]interface{}:
		return JsonStruct(t)
	// case []interface{}:
	// 	return jsonList(t)
	case float64:
		return JsonNumber(t)
	default:
		panic(fmt.Sprintf("Unexpected type %v", t))
	}
}

func (s1 JsonStruct) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	s2, ok := n.(JsonStruct)
	if !ok {
		// Different types
		e := DiffElement{
			Path:     path.clone(),
			OldValue: s1,
			NewValue: n,
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
			e := DiffElement{
				Path:     append(path.clone(), k1),
				OldValue: v1,
				NewValue: nil,
			}
			d = append(d, e)
		}
	}
	for k2, n2 := range s2 {
		v2 := newJsonNode(n2)
		if _, ok := s1[k2]; !ok {
			// S1 missing key
			e := DiffElement{
				Path:     append(path.clone(), k2),
				OldValue: nil,
				NewValue: v2,
			}
			d = append(d, e)
		}
	}
	return d
}

type PathElement interface{}
type Path []PathElement

type DiffElement struct {
	Path     Path
	OldValue JsonNode
	NewValue JsonNode
}
type Diff []DiffElement

func (p1 Path) clone() Path {
	p2 := make(Path, len(p1), len(p1)+1)
	copy(p2, p1)
	return p2
}

func readFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes)
}

func unmarshal(bytes []byte) (JsonNode, error) {
	node := make(JsonStruct)
	err := json.Unmarshal(bytes, &node)
	if err != nil {
		return nil, err
	}
	return node, nil
}
