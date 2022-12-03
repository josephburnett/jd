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

	a1Ptr, a2Ptr, pathPtr := 0, 0, 0
	for _, hash := range sequence {
		// Advance to the next common element accumulating diff elements.
		currentDiffElement := DiffElement{
			Path: append(path.clone(), jsonNumber(pathPtr)),
		}
		for a1Hashes[a1Ptr] != hash || a2Hashes[a2Ptr] != hash {

			switch {
			case a1Hashes[a1Ptr] == hash:
				// A1 is done. The rest of A2 are new values.
				for a2Hashes[a2Ptr] != hash {
					currentDiffElement.NewValues = append(currentDiffElement.NewValues, a2[a2Ptr])
					a2Ptr++
					pathPtr++
				}
			case a2Hashes[a2Ptr] == hash:
				// A2 is done. The rest of A1 are old values.
				for a1Hashes[a1Ptr] != hash {
					currentDiffElement.OldValues = append(currentDiffElement.OldValues, a1[a1Ptr])
					a1Ptr++
					pathPtr--
				}
			case sameContainerType(a1[a1Ptr], a2[a2Ptr], metadata):
				// Add what we have.
				if len(currentDiffElement.NewValues) != 0 || len(currentDiffElement.OldValues) != 0 {
					d = append(d, currentDiffElement)
				}
				// Recurse and add the subdiff.
				subDiff := a1[a1Ptr].diff(a2[a2Ptr], append(path.clone(), jsonNumber(pathPtr)), metadata, strategy)
				if len(subDiff) > 0 {
					d = append(d, subDiff...)
				}
				// Continue after subdiff.
				a1Ptr++
				a2Ptr++
				pathPtr++
				currentDiffElement = DiffElement{
					Path: append(path.clone(), jsonNumber(pathPtr)),
				}
			default:
				currentDiffElement.OldValues = append(currentDiffElement.OldValues, a1[a1Ptr])
				currentDiffElement.NewValues = append(currentDiffElement.NewValues, a2[a2Ptr])
			}

		}
		if len(currentDiffElement.NewValues) > 0 || len(currentDiffElement.OldValues) > 0 {
			d = append(d, currentDiffElement)
		}
		// Advance past common element
		a1Ptr++
		a2Ptr++
		pathPtr++
	}
	// Add all remaining elements to the diff.
	e := DiffElement{
		Path: append(path.clone(), jsonNumber(pathPtr)),
	}
	for a1Ptr < len(a1) {
		e.OldValues = append(e.OldValues, a1[a1Ptr])
		a1Ptr++
	}
	for a2Ptr < len(a2) {
		e.NewValues = append(e.NewValues, a2[a2Ptr])
		a2Ptr++
	}
	if len(e.NewValues) != 0 || len(e.OldValues) != 0 {
		d = append(d, e)
	}

	return d
}

func sameContainerType(n1, n2 JsonNode, metadata []Metadata) bool {
	c1 := dispatch(n1, metadata)
	c2 := dispatch(n2, metadata)
	switch c1.(type) {
	case jsonObject:
		if _, ok := c2.(jsonObject); ok {
			return true
		}
	case jsonList:
		if _, ok := c2.(jsonList); ok {
			return true
		}
	case jsonSet:
		if _, ok := c2.(jsonSet); ok {
			return true
		}
	case jsonMultiset:
		if _, ok := c2.(jsonMultiset); ok {
			return true
		}
	default:
		return false
	}
	return false
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
