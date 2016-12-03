package jd

import (
	"sort"
)

type jsonMultiset jsonArray

var _ JsonNode = jsonMultiset(nil)

func (a jsonMultiset) Json() string {
	return renderJson(a)
}

func (a1 jsonMultiset) Equals(n JsonNode) bool {
	a2, ok := n.(jsonMultiset)
	if !ok {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	if a1.hashCode() == a2.hashCode() {
		return true
	} else {
		return false
	}
}

func (a jsonMultiset) hashCode() [8]byte {
	h := make(hashCodes, 0, len(a))
	for _, v := range a {
		h = append(h, v.hashCode())
	}
	sort.Sort(h)
	b := make([]byte, 0, len(a)*8)
	for _, c := range h {
		b = append(b, c[:]...)
	}
	return hash(b)
}

func (a jsonMultiset) Diff(n JsonNode) Diff {
	return a.diff(n, Path{})
}

func (a1 jsonMultiset) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	a2, ok := n.(jsonMultiset)
	if !ok {
		// Different types
		e := DiffElement{
			Path:      path.clone(),
			OldValues: []JsonNode{a1},
			NewValues: []JsonNode{n},
		}
		return append(d, e)
	}
	a1HashCodes := make(map[[8]byte]int)
	a1Map := make(map[[8]byte]JsonNode)
	for _, v := range a1 {
		hc := v.hashCode()
		a1HashCodes[hc]++
		a1Map[hc] = v
	}
	a2HashCodes := make(map[[8]byte]int)
	a2Map := make(map[[8]byte]JsonNode)
	for _, v := range a2 {
		hc := v.hashCode()
		a2HashCodes[hc]++
		a2Map[hc] = v
	}
	e := DiffElement{
		Path:      append(path.clone(), map[string]interface{}{}),
		OldValues: []JsonNode{},
		NewValues: []JsonNode{},
	}
	for hc, a1Count := range a1HashCodes {
		a2Count, ok := a2HashCodes[hc]
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
	for hc, a2Count := range a2HashCodes {
		a1Count, ok := a1HashCodes[hc]
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

func (a jsonMultiset) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode) (JsonNode, error) {
	return nil, nil
}
