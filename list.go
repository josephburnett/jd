package jd

import (
	"reflect"
)

type jsonList []JsonNode

var _ JsonNode = jsonList(nil)

func (l1 jsonList) Equals(n JsonNode) bool {
	l2, ok := n.(jsonList)
	if !ok {
		return false
	}
	if len(l1) != len(l2) {
		return false
	}
	return reflect.DeepEqual(l1, l2)
}

func (l1 jsonList) Diff(n JsonNode) Diff {
	return l1.diff(n, Path{})
}

func (l1 jsonList) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	l2, ok := n.(jsonList)
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
			subDiff := l1[i].diff(l2[i], subPath)
			d = append(d, subDiff...)
		}
		if l1Has && !l2Has {

			e := DiffElement{
				Path:     subPath,
				OldValue: l1[i],
				NewValue: nil,
			}
			d = append(d, e)
		}
		if !l1Has && l2Has {
			e := DiffElement{
				Path:     subPath,
				OldValue: nil,
				NewValue: l2[i],
			}
			d = append(d, e)
		}
	}
	return d
}
