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

func (o jsonObject) hashCode() [8]byte {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	a := make([]byte, 0, len(o)*16)
	for _, k := range keys {
		keyHash := hash([]byte(k))
		a = append(a, keyHash[:]...)
		valueHash := o[k].hashCode()
		a = append(a, valueHash[:]...)
	}
	return hash(a)
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
			Path:      path.clone(),
			OldValues: []JsonNode{o1},
			NewValues: []JsonNode{n},
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
				Path:      append(path.clone(), k1),
				OldValues: nodeList(v1),
				NewValues: nodeList(),
			}
			d = append(d, e)
		}
	}
	for _, k2 := range o2Keys {
		v2 := o2[k2]
		if _, ok := o1[k2]; !ok {
			// O1 missing key
			e := DiffElement{
				Path:      append(path.clone(), k2),
				OldValues: nodeList(),
				NewValues: nodeList(v2),
			}
			d = append(d, e)
		}
	}
	return d
}

func (o jsonObject) Patch(d Diff) (JsonNode, error) {
	return patchAll(o, d)
}

func (o jsonObject) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode) (JsonNode, error) {
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	// Base case
	if len(pathAhead) == 0 {
		if !o.Equals(oldValue) {
			return patchErrExpectValue(oldValue, o, pathBehind)
		}
		return newValue, nil
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
	patchedNode, err := nextNode.patch(append(pathBehind, pe), pathAhead[1:], oldValues, newValues)
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
