package jd

import (
	"reflect"
	"sort"
)

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
