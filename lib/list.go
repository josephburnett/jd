package jd

import "fmt"

type jsonList []JsonNode

var _ JsonNode = jsonList(nil)

func (l jsonList) Json(metadata ...Metadata) string {
	return renderJson(l)
}

func (l1 jsonList) Equals(n JsonNode, metadata ...Metadata) bool {
	n2 := dispatch(n, metadata)
	l2, ok := n2.(jsonList)
	if !ok {
		return false
	}
	if len(l1) != len(l2) {
		return false
	}
	for i, v1 := range l1 {
		v2 := l2[i]
		if !v1.Equals(v2, metadata...) {
			return false
		}
	}
	return true
}

func (l jsonList) hashCode(metadata []Metadata) [8]byte {
	b := make([]byte, 0, len(l)*8)
	for _, n := range l {
		h := n.hashCode(metadata)
		b = append(b, h[:]...)
	}
	return hash(b)
}

func (l jsonList) Diff(n JsonNode, metadata ...Metadata) Diff {
	return l.diff(n, make(path, 0), metadata)
}

func (a1 jsonList) diff(n JsonNode, path path, metadata []Metadata) Diff {
	d := make(Diff, 0)
	a2, ok := n.(jsonList)
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
		subPath := append(path, jsonNumber(i))
		if a1Has && a2Has {
			n1 := dispatch(a1[i], metadata)
			n2 := dispatch(a2[i], metadata)
			subDiff := n1.diff(n2, subPath, metadata)
			d = append(d, subDiff...)
		}
		if a1Has && !a2Has {
			e := DiffElement{
				Path:      subPath.clone(),
				OldValues: nodeList(a1[i]),
				NewValues: nodeList(),
			}
			d = append(d, e)
		}
		if !a1Has && a2Has {
			e := DiffElement{
				Path:      subPath.clone(),
				OldValues: nodeList(),
				NewValues: nodeList(a2[i]),
			}
			d = append(d, e)
		}
	}
	return d
}

func (l jsonList) Patch(d Diff) (JsonNode, error) {
	return patchAll(l, d)
}

func (l jsonList) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode) (JsonNode, error) {

	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)
	// Base case
	if len(pathAhead) == 0 {
		if !l.Equals(oldValue) {
			return patchErrExpectValue(oldValue, l, pathBehind)
		}
		return newValue, nil
	}
	// Recursive case
	n, _, rest := pathAhead.next()
	jn, ok := n.(jsonNumber)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid path element %T. Expected float64.", n)
	}
	i := int(jn)
	var nextNode JsonNode = voidNode{}
	if len(l) > i {
		nextNode = l[i]
	}
	patchedNode, err := nextNode.patch(append(pathBehind, n), rest, oldValues, newValues)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		if i != len(l)-1 {
			return nil, fmt.Errorf(
				"Removal of a non-terminal element of an array.")
		}
		// Delete an element
		return l[:len(l)-1], nil
	}
	if i > len(l) {
		return nil, fmt.Errorf(
			"Addition beyond the terminal element of an array.")
	}
	if i == len(l) {
		// Add an element
		return append(l, patchedNode), nil
	}
	// Replace an element
	l[i] = patchedNode
	return l, nil
}
