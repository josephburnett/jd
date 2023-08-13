package jd

import "fmt"

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

func (a1 jsonList) diff(
	n JsonNode,
	path Path,
	options []Option,
	strategy patchStrategy,
) Diff {
	d := make(Diff, 0)
	a2, ok := n.(jsonList)
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
				Remove: nodeList(a1),
				Add:    nodeList(n),
			}
		}
		return append(d, e)
	}
	if strategy == mergePatchStrategy {
		// Merge patches do not recurse into lists
		if !a1.Equals(a2, options...) {
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
	maxLen := len(a1)
	if len(a1) < len(a2) {
		maxLen = len(a2)
	}
	from, to, by := maxLen-1, -1, -1
	if len(a1) < len(a2) {
		from, to, by = 0, maxLen, 1
	}
	for i := from; i != to; i = i + by {
		a1Has := i < len(a1)
		a2Has := i < len(a2)
		subPath := append(path, PathIndex(i))
		if a1Has && a2Has {
			n1 := dispatch(a1[i], options)
			n2 := dispatch(a2[i], options)
			subDiff := n1.diff(n2, subPath, options, strategy)
			d = append(d, subDiff...)
		}
		if a1Has && !a2Has {
			e := DiffElement{
				Path:   subPath.clone(),
				Remove: nodeList(a1[i]),
				Add:    nodeList(),
			}
			d = append(d, e)
		}
		if !a1Has && a2Has {
			appendPath := append(path, PathIndex(-1))
			e := DiffElement{
				Path:   appendPath.clone(),
				Remove: nodeList(),
				Add:    nodeList(a2[i]),
			}
			d = append(d, e)
		}
	}
	return d
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
