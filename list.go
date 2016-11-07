package jd

import (
	"reflect"
)

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
