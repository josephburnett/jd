package jd

import (
	"fmt"

	lcs "github.com/yudai/golcs"
)

type jsonList []JsonNode

var _ JsonNode = jsonList(nil)

func (l jsonList) Json(_ ...Metadata) string {
	return renderJson(l.raw())
}

func (l jsonList) Yaml(...Metadata) string {
	return renderYaml(l.raw())
}

func (l jsonList) raw() interface{} {
	return jsonArray(l).raw()
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
	return l.diff(n, make(path, 0), metadata, getPatchStrategy(metadata))
}

func (a1 jsonList) diff(n JsonNode, path path, metadata []Metadata, strategy patchStrategy) Diff {
	d := make(Diff, 0)
	a2, ok := n.(jsonList)
	if !ok {
		// Different types
		var e DiffElement
		switch strategy {
		case mergePatchStrategy:
			e = DiffElement{
				Path:      path.prependMetadataMerge(),
				NewValues: jsonArray{n},
			}
		default:
			e = DiffElement{
				Path:      path.clone(),
				OldValues: nodeList(a1),
				NewValues: nodeList(n),
			}
		}
		return append(d, e)
	}
	if strategy == mergePatchStrategy {
		// Merge patches do not recurse into lists
		if !a1.Equals(a2, metadata...) {
			e := DiffElement{
				Path:      path.prependMetadataMerge(),
				NewValues: nodeList(n),
			}
			return append(d, e)
		}
	}

	a1Hashes := make([]interface{}, len(a1))
	a2Hashes := make([]interface{}, len(a2))
	for i, v := range a1 {
		a1Hashes[i] = v.hashCode(metadata)
	}
	for i, v := range a2 {
		a2Hashes[i] = v.hashCode(metadata)
	}
	sequence := lcs.New([]interface{}(a1Hashes), []interface{}(a2Hashes)).Values()

	a1Ptr, a2Ptr := 0, 0
	for _, hash := range sequence {
		// Advance to the next common element and add to diff.
		e := DiffElement{
			Path: append(path.clone(), jsonNumber(a1Ptr)),
		}
		for a1Hashes[a1Ptr] != hash {
			e.OldValues = append(e.OldValues, a1[a1Ptr])
			a1Ptr++
		}
		for a2Hashes[a2Ptr] != hash {
			e.NewValues = append(e.NewValues, a2[a2Ptr])
			a2Ptr++
		}
		d = append(d, e)
	}
	// Add all remaining elements to the diff.
	e := DiffElement{
		Path: append(path.clone(), jsonNumber(a1Ptr)),
	}
	for a1Ptr < len(a1) {
		e.OldValues = append(e.OldValues, a1[a1Ptr])
		a1Ptr++
	}
	for a2Ptr < len(a2) {
		e.NewValues = append(e.NewValues, a2[a2Ptr])
		a2Ptr++
	}
	d = append(d, e)

	// maxLen := len(a1)
	// if len(a1) < len(a2) {
	// 	maxLen = len(a2)
	// }
	// from, to, by := maxLen-1, -1, -1
	// if len(a1) < len(a2) {
	// 	from, to, by = 0, maxLen, 1
	// }
	// for i := from; i != to; i = i + by {
	// 	a1Has := i < len(a1)
	// 	a2Has := i < len(a2)
	// 	subPath := append(path, jsonNumber(i))
	// 	if a1Has && a2Has {
	// 		n1 := dispatch(a1[i], metadata)
	// 		n2 := dispatch(a2[i], metadata)
	// 		subDiff := n1.diff(n2, subPath, metadata, strategy)
	// 		d = append(d, subDiff...)
	// 	}
	// 	if a1Has && !a2Has {
	// 		e := DiffElement{
	// 			Path:      subPath.clone(),
	// 			OldValues: nodeList(a1[i]),
	// 			NewValues: nodeList(),
	// 		}
	// 		d = append(d, e)
	// 	}
	// 	if !a1Has && a2Has {
	// 		appendPath := append(path, jsonNumber(-1))
	// 		e := DiffElement{
	// 			Path:      appendPath.clone(),
	// 			OldValues: nodeList(),
	// 			NewValues: nodeList(a2[i]),
	// 		}
	// 		d = append(d, e)
	// 	}
	// }
	return d
}

func (l jsonList) Patch(d Diff) (JsonNode, error) {
	return patchAll(l, d)
}

func (l jsonList) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode, strategy patchStrategy) (JsonNode, error) {

	if len(oldValues) > 1 || len(newValues) > 1 {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}

	if strategy == mergePatchStrategy {
		return patch(l, pathBehind, pathAhead, oldValues, newValues, mergePatchStrategy)
	}

	oldValue := singleValue(oldValues)
	newValue := singleValue(newValues)

	// Strict patch strategy
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
			"invalid path element %T: expected float64", n)
	}
	i := int(jn)

	if i == -1 {
		// Append at end of list
		i = len(l)
	}

	switch {
	case isVoid(newValue):
		var nextNode JsonNode = voidNode{}
		if len(l) > i {
			nextNode = l[i]
		}
		patchedNode, err := nextNode.patch(append(pathBehind, n), rest, oldValues, newValues, strategy)
		if err != nil {
			return nil, err
		}
		if i < 0 || i >= len(l) {
			return nil, fmt.Errorf(
				"deletion of element outside of array bounds")
		}
		if len(rest) == 0 {
			// Delete an element (base case).
			return append(l[:i], l[i+1:]...), nil
		} else {
			l[i] = patchedNode
			return l, nil
		}
	case isVoid(oldValue):
		var nextNode JsonNode = voidNode{}
		if len(l) > i && len(rest) != 0 {
			// Replacing an element.
			nextNode = l[i]
		}
		patchedNode, err := nextNode.patch(append(pathBehind, n), rest, oldValues, newValues, strategy)
		if err != nil {
			return nil, err
		}
		if i < 0 || i > len(l) {
			return nil, fmt.Errorf(
				"addition of element outside of array bounds +1")
		}
		if i == len(l) {
			// Append an element.
			return append(l, patchedNode), nil
		}
		if len(rest) == 0 {
			// Insert an element (base case).
			l = append(l[:i+1], l[i:]...)
			l[i] = patchedNode
		} else {
			// Replace an element after recursion.
			l[i] = patchedNode
		}
		return l, nil
	default:
		var nextNode JsonNode = voidNode{}
		if len(l) > i {
			nextNode = l[i]
		}
		patchedNode, err := nextNode.patch(append(pathBehind, n), rest, oldValues, newValues, strategy)
		if err != nil {
			return nil, err
		}
		// Replace an element (base case).
		l[i] = patchedNode
		return l, nil
	}
}
