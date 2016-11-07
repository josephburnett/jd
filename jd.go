package jd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
)

func newJsonNode(n interface{}) JsonNode {
	switch t := n.(type) {
	case map[string]interface{}:
		return JsonStruct(t)
	case []interface{}:
		return JsonList(t)
	case float64:
		return JsonNumber(t)
	default:
		panic(fmt.Sprintf("Unexpected type %v", t))
	}
}

type JsonNode interface {
	diff(b JsonNode, path Path) Diff
	equals(b JsonNode) bool
}

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

type JsonList []interface{}

func (l1 JsonList) equals(n JsonNode) bool {
	l2, ok := n.(JsonList)
	if !ok {
		return false
	}
	if len(l1) != len(l2) {
		return false
	}
	return reflect.DeepEqual(l1, l2)
}

func (l1 JsonList) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	l2, ok := n.(JsonList)
	if !ok {
		// Different types
		e := DiffElement{
			Path:     path.clone(),
			OldValue: l1,
			NewValue: n,
		}
		return append(d, e)
	}
	maxLen := len(l1)
	if len(l1) < len(l2) {
		maxLen = len(l2)
	}
	for i := maxLen - 1; i >= 0; i-- {
		l1Has := i < len(l1)
		l2Has := i < len(l2)
		subPath := append(path.clone(), i)
		if l1Has && l2Has {
			subDiff := newJsonNode(l1[i]).diff(newJsonNode(l2[i]), subPath)
			d = append(d, subDiff...)
		}
		if l1Has && !l2Has {
			e := DiffElement{
				Path:     subPath,
				OldValue: newJsonNode(l1[i]),
				NewValue: nil,
			}
			d = append(d, e)
		}
		if !l1Has && l2Has {
			e := DiffElement{
				Path:     subPath,
				OldValue: nil,
				NewValue: newJsonNode(l2[i]),
			}
			d = append(d, e)
		}
	}
	return d
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
	s1Keys := make([]string, 0, len(s1))
	for k := range s1 {
		s1Keys = append(s1Keys, k)
	}
	sort.Strings(s1Keys)
	s2Keys := make([]string, 0, len(s2))
	for k := range s2 {
		s2Keys = append(s2Keys, k)
	}
	sort.Strings(s2Keys)
	for _, k1 := range s1Keys {
		v1 := newJsonNode(s1[k1])
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
	for _, k2 := range s2Keys {
		v2 := newJsonNode(s2[k2])
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
