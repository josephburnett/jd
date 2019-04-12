package jd

import (
	"fmt"
	"sort"
)

type jsonSet jsonArray

var _ JsonNode = jsonSet(nil)

func (s jsonSet) Json() string {
	sMap := make(map[[8]byte]JsonNode)
	for _, n := range s {
		hc := n.hashCode()
		sMap[hc] = n
	}
	hashes := make(hashCodes, 0, len(sMap))
	for hc := range sMap {
		hashes = append(hashes, hc)
	}
	sort.Sort(hashes)
	set := make(jsonSet, 0, len(sMap))
	for _, hc := range hashes {
		set = append(set, sMap[hc])
	}
	return renderJson(set)
}

func (s1 jsonSet) Equals(n JsonNode) bool {
	s2, ok := n.(jsonSet)
	if !ok {
		return false
	}
	if s1.hashCode() == s2.hashCode() {
		return true
	} else {
		return false
	}
}

func (s jsonSet) hashCode() [8]byte {
	sMap := make(map[[8]byte]bool)
	for _, v := range s {
		hc := v.hashCode()
		sMap[hc] = true
	}
	hashes := make(hashCodes, 0, len(sMap))
	for hc := range sMap {
		hashes = append(hashes, hc)
	}
	return hashes.combine()
}

func (s jsonSet) Diff(j JsonNode) Diff {
	return s.diff(j, Path{})
}

func (s1 jsonSet) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	s2, ok := n.(jsonSet)
	if !ok {
		// Different types
		e := DiffElement{
			Path:      path.clone(),
			OldValues: nodeList(s1),
			NewValues: nodeList(n),
		}
		return append(d, e)
	}
	s1Map := make(map[[8]byte]JsonNode)
	for _, v := range s1 {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			hc = o.ident()
		} else {
			hc = v.hashCode()
		}
		s1Map[hc] = v
	}
	s2Map := make(map[[8]byte]JsonNode)
	for _, v := range s2 {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			hc = o.ident()
		} else {
			hc = v.hashCode()
		}
		s2Map[hc] = v
	}
	s1Hashes := make(hashCodes, 0)
	for hc := range s1Map {
		s1Hashes = append(s1Hashes, hc)
	}
	sort.Sort(s1Hashes)
	s2Hashes := make(hashCodes, 0)
	for hc := range s2Map {
		s2Hashes = append(s2Hashes, hc)
	}
	sort.Sort(s2Hashes)
	e := DiffElement{
		Path:      append(path.clone(), map[string]interface{}{}),
		OldValues: nodeList(),
		NewValues: nodeList(),
	}
	for _, hc := range s1Hashes {
		n2, ok := s2Map[hc]
		if !ok {
			e.OldValues = append(e.OldValues, s1Map[hc])
		} else if o2, ok := n2.(jsonObject); ok {
			// Objects with the same identity may have changed.
			subDiff := o2.diff(s1Map[hc], append(path.clone(), o2.pathIdent()))
			for _, subElement := range subDiff {
				d = append(d, subElement)
			}
		}
	}
	for _, hc := range s2Hashes {
		n1, ok := s1Map[hc]
		if !ok {
			e.NewValues = append(e.NewValues, s2Map[hc])
		} else if o1, ok := n1.(jsonObject); ok {
			// Object with the same identity may have changed.
			subDiff := o1.diff(s2Map[hc], append(path.clone(), o1.pathIdent()))
			for _, subElement := range subDiff {
				d = append(d, subElement)
			}
		}
	}
	if len(e.OldValues) > 0 || len(e.NewValues) > 0 {
		d = append(d, e)
	}
	return d
}

func (s jsonSet) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonSet) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode) (JsonNode, error) {
	// Base case
	if len(pathAhead) == 0 {
		if len(oldValues) > 1 || len(newValues) > 1 {
			return patchErrNonSetDiff(oldValues, newValues, pathBehind)
		}
		oldValue := singleValue(oldValues)
		newValue := singleValue(newValues)
		if !s.Equals(oldValue) {
			return patchErrExpectValue(oldValue, s, pathBehind)
		}
		return newValue, nil
	}
	// Unrolled recursive case
	pe, ok := pathAhead[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"Invalid path element %v. Expected map[string]interface{}.",
			pathAhead[0])
	}
	if len(pe) != 0 {
		return nil, fmt.Errorf(
			"Invalid path element %v. Expected empty object.",
			pathAhead[0])
	}
	aMap := make(map[[8]byte]JsonNode)
	for _, v := range s {
		hc := v.hashCode()
		aMap[hc] = v
	}
	for _, v := range oldValues {
		hc := v.hashCode()
		if _, ok := aMap[hc]; !ok {
			return nil, fmt.Errorf(
				"Invalid diff. Expected %v at %v bug found nothing.",
				v.Json(), pathBehind)
		}
		delete(aMap, hc)
	}
	for _, v := range newValues {
		hc := v.hashCode()
		aMap[hc] = v
	}
	hashes := make(hashCodes, 0, len(aMap))
	for hc := range aMap {
		hashes = append(hashes, hc)
	}
	sort.Sort(hashes)
	newValue := make(jsonSet, 0, len(aMap))
	for _, hc := range hashes {
		newValue = append(newValue, aMap[hc])
	}
	return newValue, nil
}