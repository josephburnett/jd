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
	b, ok := n.(jsonList)
	if !ok {
		return a.diffDifferentTypes(n, path, strategy)
	}
	if strategy == mergePatchStrategy {
		return a.diffMergePatchStrategy(b, path, options)
	}
	aHashes := make([]interface{}, len(a))
	bHashes := make([]interface{}, len(b))
	for i, v := range a {
		aHashes[i] = v.hashCode(options)
	}
	for i, v := range b {
		bHashes[i] = v.hashCode(options)
	}
	commonSequence := lcs.New([]interface{}(aHashes), []interface{}(bHashes)).Values()
	return a.diffRest(
		0,
		b,
		append(path, PathIndex(0)),
		aHashes, bHashes, commonSequence,
		voidNode{},
		options,
		strategy,
	)
}

func (a jsonList) diffRest(
	pathIndex PathIndex,
	b jsonList,
	path Path,
	aHashes, bHashes, commonSequence []interface{},
	previous JsonNode,
	options []Option,
	strategy patchStrategy,
) Diff {
	var aCursor, bCursor, commonSequenceCursor int
	pathCursor := pathIndex
	pathNow := func() Path {
		return append(path.clone().drop(), pathCursor)
	}
	endA := func() bool {
		return aCursor == len(a)
	}
	endB := func() bool {
		return bCursor == len(b)
	}
	atCommonA := func() bool {
		if endA() || len(commonSequence) == 0 {
			return false
		}
		return aHashes[aCursor] == commonSequence[0]
	}
	atCommonB := func() bool {
		if endB() || len(commonSequence) == 0 {
			return false
		}
		return bHashes[bCursor] == commonSequence[0]
	}
	d := Diff{{
		Path:   pathNow(),
		Before: []JsonNode{previous},
	}}
	haveDiff := func() bool {
		if len(d) == 0 {
			return false
		}
		if len(d[0].Add) > 0 || len(d[0].Remove) > 0 {
			return true
		}
		return false
	}
	after := func() []JsonNode {
		i := aCursor - commonSequenceCursor
		if i+1 > len(a) {
			return []JsonNode{voidNode{}}
		}
		return []JsonNode{a[i]}
	}

accumulatingDiff:
	for {
		switch {
		case endA():
			// We are at the end of A so there are no more
			// common elements. So we accumulate the rest
			// of B as additions. The path cursor advances
			// by 2 because the result is getting longer
			// by 1 and we are moving to the next element.
			for !endB() {
				d[0].Add = append(d[0].Add, b[bCursor])
				bCursor++
				pathCursor += 2
			}
			break accumulatingDiff
		case endB():
			// We are at the end of B so there are no more
			// common elements. So we accumulate the rest
			// of A as removals. The path cursor stays the
			// same because the result is getting shorter
			// by 1 but we are also moving to the next
			// element.
			for !endA() {
				d[0].Remove = append(d[0].Remove, a[aCursor])
				aCursor++
			}
			break accumulatingDiff
		case atCommonA() && atCommonB():
			// We are at a common element of A and B.
			// All cursors advance because we are moving
			// past a common element.
			aCursor++
			bCursor++
			commonSequenceCursor++
			pathCursor++
			break accumulatingDiff
		case atCommonA():
			// We are at a common element in A. We need to
			// catch up B. Add elements of B until we do.
			for !atCommonB() {
				d[0].Add = append(d[0].Add, b[bCursor])
				bCursor++
				pathCursor++
			}
		case atCommonB():
			// We are at a common element in B. We need to
			// catch up A. Remove elements of A until we
			// do.
			for !atCommonA() {
				d[0].Remove = append(d[0].Remove, a[aCursor])
				aCursor++
			}
		case sameContainerType(a[aCursor], b[bCursor], options):
			// We are at compatible containers which
			// contain additional differences. If we've
			// accumulated differences at this level then
			// keep them before the sub-diff.
			subDiff := a[aCursor].diff(b[bCursor], pathNow(), options, strategy)
			if haveDiff() {
				d[0].After = after()
				d = append(d, subDiff...)
			} else {
				d = subDiff
			}
			aCursor++
			bCursor++
			pathCursor++
			break accumulatingDiff
		default:
			// We are at elements of A and B which are
			// different. Add them to the accumulated diff
			// and continue.
			d[0].Remove = append(d[0].Remove, a[aCursor])
			d[0].Add = append(d[0].Add, b[bCursor])
			aCursor++
			bCursor++
			pathCursor++
		}
	}

	if !haveDiff() {
		// Throw away temporary diff because we didn't
		// accumulate anything.
		d = Diff{}
	} else {
		if len(d[0].Path) > len(path) {
			// This is a subdiff. Don't touch it.
		} else {
			// Record context of accumulated diff. If we appended
			// a sub-diff then it already has context.
			if len(d) < 2 {
				d[0].After = after()
			}
		}
	}
	if endA() && endB() {
		return d
	}
	// Cursors point to the next elements.
	return append(d, a[aCursor:].diffRest(
		pathCursor,
		b[bCursor:],
		pathNow(),
		aHashes[aCursor:], bHashes[bCursor:], commonSequence[commonSequenceCursor:],
		b[bCursor-1],
		options,
		strategy,
	)...)
}

func (a jsonList) diffDifferentTypes(n JsonNode, path Path, strategy patchStrategy) Diff {
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
	return Diff{e}
}

func (a jsonList) diffMergePatchStrategy(b jsonList, path Path, options []Option) Diff {
	if !a.Equals(b, options...) {
		e := DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: path.clone(),
			Add:  nodeList(b),
		}
		return Diff{e}
	}
	return Diff{}
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

func (l jsonList) patch(pathBehind, pathAhead Path, before, removeValues, addValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {

	if strategy == mergePatchStrategy {
		return patch(l, pathBehind, pathAhead, before, removeValues, addValues, after, mergePatchStrategy)
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
		patchedNode, err := l[i].patch(append(pathBehind, n), rest, nil, removeValues, addValues, nil, strategy)
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

	// Check context before
	for j, b := range before {
		bIndex := int(i) - (len(before) - j)
		switch {
		case bIndex < 0:
			if bIndex == -1 && isVoid(b) {
				continue
			}
			return nil, fmt.Errorf("invalid patch. before context %v out of bounds: %v", b, bIndex)
		case !b.Equals(l[bIndex]):
			return nil, fmt.Errorf("invalid patch. expected %v before. got %v", b, l[bIndex])
		}
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

	// Check context after
	for j, a := range after {
		aIndex := int(i) + j
		if aIndex > len(l)-1 {
			if aIndex == len(l) && isVoid(a) {
				continue
			}
			return nil, fmt.Errorf("invalid patch. after context %v out of bounds: %v", a, aIndex)
		}
		if !a.Equals(l[aIndex]) {
			return nil, fmt.Errorf("invalid patch. expected %v after. got %v", a, l[aIndex])
		}
	}

	return l2, nil
}
