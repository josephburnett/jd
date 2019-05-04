package jd

import (
	"fmt"
	"reflect"
)

type jsonArray []JsonNode

var _ JsonNode = jsonArray(nil)

func (a jsonArray) Json() string {
	return renderJson(a)
}

func (a1 jsonArray) Equals(n JsonNode, metadata ...Metadata) bool {
	a2, ok := n.(jsonArray)
	if !ok {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	return reflect.DeepEqual(a1, a2)
}

func (a jsonArray) hashCode() [8]byte {
	b := make([]byte, 0, len(a)*8)
	for _, el := range a {
		h := el.hashCode()
		b = append(b, h[:]...)
	}
	return hash(b)
}

func (a jsonArray) Diff(n JsonNode, metadata ...Metadata) Diff {
	return a.diff(n, Path{}, metadata)
}

func (a1 jsonArray) diff(n JsonNode, path Path, metadata []Metadata) Diff {
	if checkMetadata(SET, metadata) {
		// Use set semantics.
		if n2, ok := n.(jsonArray); ok {
			n = jsonSet(n2)
		}
		return jsonSet(a1).diff(n, path, metadata)
	}
	if checkMetadata(MULTISET, metadata) {
		// Use multiset semantics.
		if n2, ok := n.(jsonArray); ok {
			n = jsonMultiset(n2)
		}
		return jsonMultiset(a1).diff(n, path, metadata)
	}
	// Use list semantics.
	d := make(Diff, 0)
	a2, ok := n.(jsonArray)
	if !ok {
		// Different types
		e := DiffElement{
			Path:      path.clone(),
			OldValues: nodeList(a1),
			NewValues: nodeList(n),
		}
		return append(d, e)
	}
	maxLen := len(a1)
	if len(a1) < len(a2) {
		maxLen = len(a2)
	}
	for i := 0; i < maxLen; i++ {
		a1Has := i < len(a1)
		a2Has := i < len(a2)
		subPath := append(path.clone(), float64(i))
		if a1Has && a2Has {
			subDiff := a1[i].diff(a2[i], subPath, metadata)
			d = append(d, subDiff...)
		}
		if a1Has && !a2Has {
			e := DiffElement{
				Path:      subPath,
				OldValues: nodeList(a1[i]),
				NewValues: nodeList(),
			}
			d = append(d, e)
		}
		if !a1Has && a2Has {
			e := DiffElement{
				Path:      subPath,
				OldValues: nodeList(),
				NewValues: nodeList(a2[i]),
			}
			d = append(d, e)
		}
	}
	return d
}

func (a jsonArray) Patch(d Diff, metadata ...Metadata) (JsonNode, error) {
	return patchAll(a, d, metadata)
}

func (a jsonArray) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode, metadata []Metadata) (JsonNode, error) {
	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	// Base case
	if len(pathAhead) == 0 {
		if !a.Equals(oldValue) {
			return patchErrExpectValue(oldValue, a, pathBehind)
		}
		return newValue, nil
	}
	// Recursive case
	pe, ok := pathAhead[0].(float64)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid path element %v. Expected float64.",
			pathAhead[0])
	}
	i := int(pe)
	var nextNode JsonNode = voidNode{}
	if len(a) > i {
		nextNode = a[i]
	}
	patchedNode, err := nextNode.patch(append(pathBehind, pe), pathAhead[1:], oldValues, newValues, metadata)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		if i != len(a)-1 {
			return nil, fmt.Errorf(
				"Removal of a non-terminal element of an array.")
		}
		// Delete an element
		return a[:len(a)-1], nil
	}
	if i > len(a) {
		return nil, fmt.Errorf(
			"Addition beyond the terminal element of an array.")
	}
	if i == len(a) {
		// Add an element
		return append(a, patchedNode), nil
	}
	// Replace an element
	a[i] = patchedNode
	return a, nil
}
