package jd

import (
	"fmt"
	"reflect"
	"sort"
)

type jsonObject map[string]JsonNode

var _ JsonNode = jsonObject(nil)

func (o jsonObject) Json() string {
	return renderJson(o)
}

func (o1 jsonObject) Equals(n JsonNode) bool {
	o2, ok := n.(jsonObject)
	if !ok {
		return false
	}
	if len(o1) != len(o2) {
		return false
	}
	return reflect.DeepEqual(o1, o2)
}

func (o jsonObject) Diff(n JsonNode) Diff {
	return o.diff(n, Path{})
}

func (o1 jsonObject) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	o2, ok := n.(jsonObject)
	if !ok {
		// Different types
		e := DiffElement{
			Path:     path.clone(),
			OldValue: o1,
			NewValue: n,
		}
		return append(d, e)
	}
	o1Keys := make([]string, 0, len(o1))
	for k := range o1 {
		o1Keys = append(o1Keys, k)
	}
	sort.Strings(o1Keys)
	o2Keys := make([]string, 0, len(o2))
	for k := range o2 {
		o2Keys = append(o2Keys, k)
	}
	sort.Strings(o2Keys)
	for _, k1 := range o1Keys {
		v1 := o1[k1]
		if v2, ok := o2[k1]; ok {
			// Both keys are present
			subDiff := v1.diff(v2, append(path.clone(), k1))
			d = append(d, subDiff...)
		} else {
			// O2 missing key
			e := DiffElement{
				Path:     append(path.clone(), k1),
				OldValue: v1,
				NewValue: voidNode{},
			}
			d = append(d, e)
		}
	}
	for _, k2 := range o2Keys {
		v2 := o2[k2]
		if _, ok := o1[k2]; !ok {
			// O1 missing key
			e := DiffElement{
				Path:     append(path.clone(), k2),
				OldValue: voidNode{},
				NewValue: v2,
			}
			d = append(d, e)
		}
	}
	return d
}

func (o jsonObject) Patch(d Diff) (JsonNode, error) {
	return patchAll(o, d)
}

func (o jsonObject) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	// Base case
	if len(pathAhead) == 0 {
		if !o.Equals(oldValue) {
			return nil, fmt.Errorf(
				"Found %v at %v. Expected %v.",
				o.Json(), pathBehind, oldValue.Json())
		}
	}
	// Recursive case
	pe, ok := pathAhead[0].(string)
	if !ok {
		return nil, fmt.Errorf(
			"Found %v at %v. Expected JSON object.",
			o.Json(), pathBehind)
	}
	nextNode, ok := o[pe]
	if !ok {
		nextNode = voidNode{}
	}
	patchedNode, err := nextNode.patch(append(pathBehind, pe), pathAhead[1:], oldValue, newValue)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		// Delete a pair
		delete(o, pe)
	} else {
		// Add or replace a pair
		o[pe] = patchedNode
	}
	return o, nil
}
