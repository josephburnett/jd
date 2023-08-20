package jd

import (
	"fmt"

	lcs "github.com/yudai/golcs"
)

type jsonList []JsonNode

var _ JsonNode = jsonList(nil)

func (l jsonList) Json(_ ...Option) string {
	return renderJson(l.raw())
}

func (l jsonList) Yaml(_ ...Option) string {
	return renderYaml(l.raw())
}

func (l jsonList) raw() interface{} {
	return jsonArray(l).raw()
}

func (l1 jsonList) Equals(n JsonNode, options ...Option) bool {
	n2 := dispatch(n, options)
	l2, ok := n2.(jsonList)
	if !ok {
		return false
	}
	if len(l1) != len(l2) {
		return false
	}
	for i, v1 := range l1 {
		v2 := l2[i]
		if !v1.Equals(v2, options...) {
			return false
		}
	}
	return true
}

func (l jsonList) hashCode(options []Option) [8]byte {
	b := make([]byte, 0, len(l)*8)
	for _, n := range l {
		h := n.hashCode(options)
		b = append(b, h[:]...)
	}
	return hash(b)
}

func (l jsonList) Diff(n JsonNode, options ...Option) Diff {
	return l.diff(n, make(Path, 0), options, getPatchStrategy(options))
}

func (a jsonList) diff(
	n JsonNode,
	path Path,
	options []Option,
	strategy patchStrategy,
) Diff {
	d := make(Diff, 0)
	b, ok := n.(jsonList)
	if !ok {
		// Different types
		var e DiffElement
		switch strategy {
		case mergePatchStrategy:
			e = DiffElement{
				Metadata: Metadata{
					Merge: true,
				},
				Path: path.clone(),
				Add:  jsonArray{n},
			}
		default:
			e = DiffElement{
				Path:   path.clone(),
				Remove: nodeList(a),
				Add:    nodeList(n),
			}
		}
		return append(d, e)
	}
	if strategy == mergePatchStrategy {
		// Merge patches do not recurse into lists
		if !a.Equals(b, options...) {
			e := DiffElement{
				Metadata: Metadata{
					Merge: true,
				},
				Path: path.clone(),
				Add:  nodeList(n),
			}
			return append(d, e)
		}
	}
	aHashes := make([]interface{}, len(a))
	bHashes := make([]interface{}, len(b))
	for i, v := range a {
		aHashes[i] = v.hashCode(options)
	}
	for i, v := range b {
		bHashes[i] = v.hashCode(options)
	}
	sequence := lcs.New([]interface{}(aHashes), []interface{}(bHashes)).Values()
	aCursor, bCursor, pathCursor := 0, 0, 0
	lastACursor, lastBCursor := -1, -1
	for _, hash := range sequence {
		// Advanced to the next common element accumulating diff elements.
		currentDiffElement := DiffElement{
			Path: append(path.clone(), PathIndex(pathCursor)),
		}
		for aHashes[aCursor] != hash || bHashes[bCursor] != hash {
			if aCursor == lastACursor && bCursor == lastBCursor {
				panic("a and b cursors are not advancing")
			}
			lastACursor = aCursor
			lastBCursor = bCursor
			switch {
			case aHashes[aCursor] == hash:
				// A is done. The rest of B are new values.
				for bHashes[bCursor] != hash {
					currentDiffElement.Add = append(currentDiffElement.Add, b[bCursor])
					bCursor++
					pathCursor++
				}
			case bHashes[bCursor] == hash:
				// B is done. The rest of A are old values.
				for aHashes[aCursor] != hash {
					currentDiffElement.Remove = append(currentDiffElement.Remove, a[aCursor])
					aCursor++
					pathCursor--
				}
			case sameContainerType(a[aCursor], b[bCursor], options):
				// Add what we have.
				if len(currentDiffElement.Add) != 0 || len(currentDiffElement.Remove) != 0 {
					d = append(d, currentDiffElement)
				}
				// Recurse and add the subdiff.
				subDiff := a[aCursor].diff(b[bCursor], append(path.clone(), PathIndex(pathCursor)), options, strategy)
				if len(subDiff) > 0 {
					d = append(d, subDiff...)
				}
				// Continue after subdiff.
				aCursor++
				bCursor++
				pathCursor++
				currentDiffElement = DiffElement{
					Path: append(path.clone(), PathIndex(pathCursor)),
				}
			default:
				currentDiffElement.Remove = append(currentDiffElement.Remove, a[aCursor])
				currentDiffElement.Add = append(currentDiffElement.Add, b[bCursor])
				aCursor++
				bCursor++
			}
		}
		if len(currentDiffElement.Add) > 0 || len(currentDiffElement.Remove) > 0 {
			d = append(d, currentDiffElement)
		}
		// // Advance past common element
		// aCursor++
		// bCursor++
		// pathCursor++
	}
	// Recurse into remaining containers
	isSameContainerType := true
	for aCursor < len(a) && bCursor < len(b) && isSameContainerType {
		isSameContainerType = sameContainerType(a[aCursor], b[bCursor], options)
		if isSameContainerType {
			// Recurse and add the subdiff.
			subDiff := a[aCursor].diff(b[bCursor], append(path.clone(), PathIndex(pathCursor)), options, strategy)
			if len(subDiff) > 0 {
				d = append(d, subDiff...)
			}
			// Continue after subdiff.
			aCursor++
			bCursor++
			pathCursor++
		}
	}
	// Add all remaining elements to the diff.
	e := DiffElement{
		Path: append(path.clone(), PathIndex(pathCursor)),
	}
	for aCursor < len(a) {
		e.Remove = append(e.Remove, a[aCursor])
		aCursor++
	}
	for bCursor < len(b) {
		e.Add = append(e.Add, b[bCursor])
		bCursor++
	}
	if len(e.Add) != 0 || len(e.Remove) != 0 {
		d = append(d, e)
	}

	return d
}

func sameContainerType(n1, n2 JsonNode, options []Option) bool {
	c1 := dispatch(n1, options)
	c2 := dispatch(n2, options)
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

func (l jsonList) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode, strategy patchStrategy) (JsonNode, error) {

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
	jn, ok := n.(PathIndex)
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
