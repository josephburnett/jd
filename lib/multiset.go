package jd

import (
	"fmt"
	"sort"
)

type jsonMultiset jsonArray

var _ JsonNode = jsonMultiset(nil)

func (a jsonMultiset) Json(metadata ...Metadata) string {
	return renderJson(a)
}

func (a1 jsonMultiset) Equals(n JsonNode, metadata ...Metadata) bool {
	a2, ok := n.(jsonMultiset)
	if !ok {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	if a1.hashCode(metadata) == a2.hashCode(metadata) {
		return true
	} else {
		return false
	}
}

func (a jsonMultiset) hashCode(metadata []Metadata) [8]byte {
	h := make(hashCodes, 0, len(a))
	for _, v := range a {
		h = append(h, v.hashCode(metadata))
	}
	sort.Sort(h)
	b := make([]byte, 0, len(a)*8)
	for _, c := range h {
		b = append(b, c[:]...)
	}
	return hash(b)
}

func (a jsonMultiset) Diff(n JsonNode, metadata ...Metadata) Diff {
	return a.diff(n, nil, metadata)
}

func (a1 jsonMultiset) diff(n JsonNode, path path, metadata []Metadata) Diff {
	d := make(Diff, 0)
	a2, ok := n.(jsonMultiset)
	if !ok {
		// Different types
		e := DiffElement{
			Path:      path.clone(),
			OldValues: nodeList(a1),
			NewValues: nodeList(n),
		}
		return append(d, e)
	}
	a1Counts := make(map[[8]byte]int)
	a1Map := make(map[[8]byte]JsonNode)
	for _, v := range a1 {
		hc := v.hashCode(metadata)
		a1Counts[hc]++
		a1Map[hc] = v
	}
	a2Counts := make(map[[8]byte]int)
	a2Map := make(map[[8]byte]JsonNode)
	for _, v := range a2 {
		hc := v.hashCode(metadata)
		a2Counts[hc]++
		a2Map[hc] = v
	}
	// TODO: cast directly to jsonObject when jsonObject drops idKeys.
	o, _ := NewJsonNode(map[string]interface{}{})
	e := DiffElement{
		Path:      path.appendIndex(o.(jsonObject), metadata).clone(),
		OldValues: nodeList(),
		NewValues: nodeList(),
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
				e.OldValues = append(e.OldValues, a1Map[hc])
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
				e.NewValues = append(e.NewValues, a2Map[hc])
			}
		}
	}
	if len(e.OldValues) > 0 || len(e.NewValues) > 0 {
		d = append(d, e)
	}
	return d
}

func (a jsonMultiset) Patch(d Diff) (JsonNode, error) {
	return patchAll(a, d)
}

func (a jsonMultiset) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode) (JsonNode, error) {
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
	o, ok := n.(jsonObject)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid path element %v. Expected map[string]interface{}.", n)
	}
	if len(o.properties) != 0 {
		return nil, fmt.Errorf(
			"Invalid path element %v. Expected empty object.", n)
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
				"Invalid diff. Expected %v at %v but found nothing.",
				aMap[hc].Json(metadata...), pathBehind)
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
