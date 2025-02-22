package jd

import (
	"fmt"
	"sort"
)

type jsonMultiset jsonArray

var _ JsonNode = jsonMultiset(nil)

func (a jsonMultiset) Json(_ ...Option) string {
	return renderJson(a.raw())
}

func (a jsonMultiset) Yaml(_ ...Option) string {
	return renderYaml(a.raw())
}

func (a jsonMultiset) raw() interface{} {
	return jsonArray(a).raw()
}

func (a1 jsonMultiset) Equals(n JsonNode, options ...Option) bool {
	n = dispatch(n, options)
	a2, ok := n.(jsonMultiset)
	if !ok {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	if a1.hashCode(options) == a2.hashCode(options) {
		return true
	} else {
		return false
	}
}

func (a jsonMultiset) hashCode(options []Option) [8]byte {
	h := make(hashCodes, 0, len(a))
	for _, v := range a {
		h = append(h, v.hashCode(options))
	}
	sort.Sort(h)
	b := make([]byte, 0, len(a)*8)
	for _, c := range h {
		b = append(b, c[:]...)
	}
	return hash(b)
}

func (a jsonMultiset) Diff(n JsonNode, options ...Option) Diff {
	return a.diff(n, nil, options, getPatchStrategy(options))
}

func (a1 jsonMultiset) diff(
	n JsonNode,
	path Path,
	options []Option,
	strategy patchStrategy,
) Diff {
	d := make(Diff, 0)
	a2, ok := n.(jsonMultiset)
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
				Add:  nodeList(n),
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
	if strategy == mergePatchStrategy && !a1.Equals(n) {
		e := DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: path.clone(),
			Add:  nodeList(n),
		}
		return append(d, e)
	}
	a1Counts := make(map[[8]byte]int)
	a1Map := make(map[[8]byte]JsonNode)
	for _, v := range a1 {
		hc := v.hashCode(options)
		a1Counts[hc]++
		a1Map[hc] = v
	}
	a2Counts := make(map[[8]byte]int)
	a2Map := make(map[[8]byte]JsonNode)
	for _, v := range a2 {
		hc := v.hashCode(options)
		a2Counts[hc]++
		a2Map[hc] = v
	}
	pathWithMultiset := append(path.clone(), PathMultiset{})
	e := DiffElement{
		Path:   pathWithMultiset,
		Remove: nodeList(),
		Add:    nodeList(),
	}
	a1Hashes := make(hashCodes, 0)
	for hc := range a1Counts {
		a1Hashes = append(a1Hashes, hc)
	}
	sort.Sort(a1Hashes)
	a2Hashes := make(hashCodes, 0)
	for hc := range a2Counts {
		a2Hashes = append(a2Hashes, hc)
	}
	sort.Sort(a2Hashes)
	for _, hc := range a1Hashes {
		a1Count := a1Counts[hc]
		a2Count, ok := a2Counts[hc]
		if !ok {
			a2Count = 0
		}
		removed := a1Count - a2Count
		if removed > 0 {
			for i := 0; i < removed; i++ {
				e.Remove = append(e.Remove, a1Map[hc])
			}
		}
	}
	for _, hc := range a2Hashes {
		a2Count := a2Counts[hc]
		a1Count, ok := a1Counts[hc]
		if !ok {
			a1Count = 0
		}
		added := a2Count - a1Count
		if added > 0 {
			for i := 0; i < added; i++ {
				e.Add = append(e.Add, a2Map[hc])
			}
		}
	}
	if len(e.Remove) > 0 || len(e.Add) > 0 {
		d = append(d, e)
	}
	return d
}

func (a jsonMultiset) Patch(d Diff) (JsonNode, error) {
	return patchAll(a, d)
}

func (a jsonMultiset) patch(pathBehind, pathAhead Path, before, oldValues, newValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {

	// Merge patch strategy
	if strategy == mergePatchStrategy {
		return patch(a, pathBehind, pathAhead, before, oldValues, newValues, after, mergePatchStrategy)
	}

	// Strict patch strategy
	// Base case
	if len(pathAhead) == 0 {
		if len(oldValues) > 1 || len(newValues) > 1 {
			return patchErrNonSetDiff(oldValues, newValues, pathBehind)
		}
		oldValue := singleValue(oldValues)
		newValue := singleValue(newValues)
		if !a.Equals(oldValue) {
			return patchErrExpectValue(oldValue, a, pathBehind)
		}
		return newValue, nil
	}
	// Unrolled recursive case
	n, metadata, _ := pathAhead.next()
	_, ok := n.(PathMultiset)
	if !ok {
		return nil, fmt.Errorf(
			"invalid path element %v: expected map[string]interface{}", n)
	}
	aCounts := make(map[[8]byte]int)
	aMap := make(map[[8]byte]JsonNode)
	for _, v := range a {
		hc := v.hashCode(metadata)
		aCounts[hc]++
		aMap[hc] = v
	}
	for _, v := range oldValues {
		hc := v.hashCode(metadata)
		aCounts[hc]--
		aMap[hc] = v
	}
	for hc, count := range aCounts {
		if count < 0 {
			return nil, fmt.Errorf(
				"invalid diff: expected %v at %v but found nothing",
				aMap[hc].Json(), pathBehind)
		}
	}
	for _, v := range newValues {
		hc := v.hashCode(metadata)
		aCounts[hc]++
		aMap[hc] = v
	}
	aHashes := make(hashCodes, 0)
	for hc := range aCounts {
		if aCounts[hc] > 0 {
			for i := 0; i < aCounts[hc]; i++ {
				aHashes = append(aHashes, hc)
			}
		}
	}
	sort.Sort(aHashes)
	newValue := make(jsonMultiset, 0)
	for _, hc := range aHashes {
		newValue = append(newValue, aMap[hc])
	}
	return newValue, nil
}
