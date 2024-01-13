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
	b := []byte{0xF5, 0x18, 0x0A, 0x71, 0xA4, 0xC4, 0x03, 0xF3} // random bytes
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
		// Advance past common element
		aCursor++
		bCursor++
		pathCursor++
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

func (l jsonList) patch(pathBehind, pathAhead Path, removeValues, addValues []JsonNode, strategy patchStrategy) (JsonNode, error) {

	if strategy == mergePatchStrategy {
		return patch(l, pathBehind, pathAhead, removeValues, addValues, mergePatchStrategy)
	}

	// Special case for replacing the whole list
	if len(pathAhead) == 0 {
		if len(removeValues) > 1 || len(addValues) > 1 {
			return nil, fmt.Errorf("cannot replace list with multiple values")
		}
		if len(removeValues) == 0 && strategy == strictPatchStrategy {
			return nil, fmt.Errorf("invalid diff. must declare list to replace it")
		}
		if !l.Equals(removeValues[0]) {
			return nil, fmt.Errorf("wanted %v. found %v", removeValues[0], l)
		}
		if len(addValues) == 0 {
			return voidNode{}, nil
		} else {
			return addValues[0], nil
		}
	}

	n, _, rest := pathAhead.next()
	i, ok := n.(PathIndex)
	if !ok {
		return nil, fmt.Errorf("invalid path element %T: expected float64", n)
	}

	// Recursive case
	if len(rest) > 0 {
		if int(i) > len(l)-1 {
			return nil, fmt.Errorf("patch index out of bounds: %v", i)
		}
		patchedNode, err := l[i].patch(append(pathBehind, n), rest, removeValues, addValues, strategy)
		if err != nil {
			return nil, err
		}
		l[i] = patchedNode
		return l, nil
	}

	// Special case for appending to the end of list
	if int(i) == -1 {
		if len(removeValues) > 0 {
			return nil, fmt.Errorf("invalid patch. appending to -1 index. but want to remove values")
		}
		l = append(l, addValues...)
		return l, nil
	}

	// Patch list
	for len(removeValues) > 0 {
		if int(i) > len(l)-1 {
			return nil, fmt.Errorf("remove values out bounds: %v", i)
		}
		if !l[i].Equals(removeValues[0]) {
			return nil, fmt.Errorf("invalid patch. wanted %v. found %v", removeValues[0], l[i])
		}
		l = append(l[:i], l[i+1:]...)
		removeValues = removeValues[1:]
	}
	l2 := make(jsonList, i)
	copy(l2, l[:i])
	l2 = append(l2, addValues...)
	if int(i) < len(l) {
		l2 = append(l2, l[i:]...)
	}
	return l2, nil
}
