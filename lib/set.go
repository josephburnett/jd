package jd

import (
	"fmt"
	"sort"
)

type jsonSet jsonArray

var _ JsonNode = jsonSet(nil)

func (s jsonSet) Json(metadata ...Metadata) string {
	sMap := make(map[[8]byte]JsonNode)
	for _, n := range s {
		hc := n.hashCode(metadata)
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

func (s1 jsonSet) Equals(n JsonNode, metadata ...Metadata) bool {
	s2, ok := n.(jsonSet)
	if !ok {
		return false
	}
	if s1.hashCode(metadata) == s2.hashCode(metadata) {
		return true
	} else {
		return false
	}
}

func (s jsonSet) hashCode(metadata []Metadata) [8]byte {
	sMap := make(map[[8]byte]bool)
	for _, v := range s {
		v = dispatch(v, metadata)
		hc := v.hashCode(metadata)
		sMap[hc] = true
	}
	hashes := make(hashCodes, 0, len(sMap))
	for hc := range sMap {
		hashes = append(hashes, hc)
	}
	return hashes.combine()
}

func (s jsonSet) Diff(j JsonNode, metadata ...Metadata) Diff {
	return s.diff(j, make(path, 0), metadata)
}

func (s1 jsonSet) diff(n JsonNode, path path, metadata []Metadata) Diff {
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
			// Hash objects by their identity.
			hc = o.ident(metadata)
		} else {
			// Everything else by full content.
			hc = v.hashCode(metadata)
		}
		s1Map[hc] = v
	}
	s2Map := make(map[[8]byte]JsonNode)
	for _, v := range s2 {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identity.
			hc = o.ident(metadata)
		} else {
			// Everything else by full content.
			hc = v.hashCode(metadata)
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
	o, _ := NewJsonNode(map[string]interface{}{})
	e := DiffElement{
		Path:      path.appendIndex(o.(jsonObject), metadata).clone(),
		OldValues: nodeList(),
		NewValues: nodeList(),
	}
	for _, hc := range s1Hashes {
		n2, ok := s2Map[hc]
		if !ok {
			// Deleted value.
			e.OldValues = append(e.OldValues, s1Map[hc])
		} else {
			// Changed value.
			o1, isObject1 := s1Map[hc].(jsonObject)
			o2, isObject2 := n2.(jsonObject)
			if isObject1 && isObject2 {
				// Sub diff objects with same identity.
				p := path.appendIndex(o1, metadata)
				subDiff := o1.diff(o2, p, metadata)
				for _, subElement := range subDiff {
					d = append(d, subElement)
				}
			}
		}
	}
	for _, hc := range s2Hashes {
		_, ok := s1Map[hc]
		if !ok {
			// Added value.
			e.NewValues = append(e.NewValues, s2Map[hc])
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

func (s jsonSet) patch(pathBehind, pathAhead path, oldValues, newValues []JsonNode) (JsonNode, error) {
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
	n, metadata, rest := pathAhead.next()
	o, ok := n.(jsonObject)
	if !ok {
		return nil, fmt.Errorf(
			"Invalid path element %v. Expected jsonObject.", n)
	}
	if len(rest) > 0 {
		// Recurse into a specific object.
		lookingFor := o.ident(metadata)
		for _, v := range s {
			if o, ok := v.(jsonObject); ok {
				id := o.ident(metadata)
				if id == lookingFor {
					v.patch(append(pathBehind, n), rest, oldValues, newValues)
					return s, nil
				}
			}
		}
		return nil, fmt.Errorf("Invalid diff. Expected object with id %v but found none", o.Json(metadata...))
	}
	// Patch set
	aMap := make(map[[8]byte]JsonNode)
	for _, v := range s {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identitiy.
			hc = o.ident(metadata)
		} else {
			// Everything else by full content.
			hc = v.hashCode(metadata)
		}
		aMap[hc] = v
	}
	for _, v := range oldValues {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Find objects by their identitiy.
			hc = o.ident(metadata)
		} else {
			// Everything else by full content.
			hc = v.hashCode(metadata)
		}
		toDelete, ok := aMap[hc]
		if !ok {
			return nil, fmt.Errorf(
				"Invalid diff. Expected %v at %v but found nothing.",
				v.Json(metadata...), pathBehind)
		}
		if !toDelete.Equals(v, metadata...) {
			return nil, fmt.Errorf(
				"Invalid diff. Expected %v at %v but found %v.",
				v.Json(metadata...), pathBehind, toDelete.Json(metadata...))

		}
		delete(aMap, hc)
	}
	for _, v := range newValues {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identitiy.
			hc = o.ident(metadata)
		} else {
			// Everything else by full content.
			hc = v.hashCode(metadata)
		}
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
