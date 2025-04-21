package jd

import (
	"fmt"
	"sort"
)

type jsonSet jsonArray

var _ JsonNode = jsonSet(nil)

func (s jsonSet) Json(_ ...Option) string {
	return renderJson(s.raw())
}

func (s jsonSet) Yaml(_ ...Option) string {
	return renderYaml(s.raw())
}

func (s jsonSet) raw() interface{} {
	sMap := make(map[[8]byte]JsonNode)
	for _, n := range s {
		hc := n.hashCode(&options{retain: []Option{setOption{}}})
		sMap[hc] = n
	}
	hashes := make(hashCodes, 0, len(sMap))
	for hc := range sMap {
		hashes = append(hashes, hc)
	}
	sort.Sort(hashes)
	set := make([]interface{}, 0, len(sMap))
	for _, hc := range hashes {
		set = append(set, sMap[hc].raw())
	}
	return set
}

func (s1 jsonSet) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return s1.equals(n, o)
}

func (s1 jsonSet) equals(n JsonNode, o *options) bool {
	n = dispatch(n, o)
	s2, ok := n.(jsonSet)
	if !ok {
		return false
	}
	if s1.hashCode(o) == s2.hashCode(o) {
		return true
	} else {
		return false
	}
}

func (s jsonSet) hashCode(opts *options) [8]byte {
	sMap := make(map[[8]byte]bool)
	for _, v := range s {
		v = dispatch(v, opts)
		hc := v.hashCode(opts)
		sMap[hc] = true
	}
	hashes := make(hashCodes, 0, len(sMap))
	for hc := range sMap {
		hashes = append(hashes, hc)
	}
	return hashes.combine()
}

func (s jsonSet) Diff(j JsonNode, opts ...Option) Diff {
	o := refine(&options{retain: opts}, nil)
	return s.diff(j, make(Path, 0), o, getPatchStrategy(o))
}

func (s1 jsonSet) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	d := make(Diff, 0)
	s2, ok := n.(jsonSet)
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
				Remove: nodeList(s1),
				Add:    nodeList(n),
			}
		}
		return append(d, e)
	}
	if strategy == mergePatchStrategy && !s1.Equals(n) {
		e := DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: path.clone(),
			Add:  nodeList(n),
		}
		return append(d, e)
	}
	s1Map := make(map[[8]byte]JsonNode)
	for _, v := range s1 {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identity.
			hc = o.ident(opts)
		} else {
			// Everything else by full content.
			hc = v.hashCode(opts)
		}
		s1Map[hc] = v
	}
	s2Map := make(map[[8]byte]JsonNode)
	for _, v := range s2 {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identity.
			hc = o.ident(opts)
		} else {
			// Everything else by full content.
			hc = v.hashCode(opts)
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
		Path:   append(path.clone(), PathSet{}),
		Remove: nodeList(),
		Add:    nodeList(),
	}
	for _, hc := range s1Hashes {
		n2, ok := s2Map[hc]
		if !ok {
			// Deleted value.
			e.Remove = append(e.Remove, s1Map[hc])
		} else {
			// Changed value.
			o1, isObject1 := s1Map[hc].(jsonObject)
			o2, isObject2 := n2.(jsonObject)
			if isObject1 && isObject2 {
				// Sub diff objects with same identity.
				p := append(path.clone(), newPathSetKeys(o1, opts))
				subDiff := o1.diff(o2, p, opts, strategy)
				d = append(d, subDiff...)
			}
		}
	}
	for _, hc := range s2Hashes {
		_, ok := s1Map[hc]
		if !ok {
			// Added value.
			e.Add = append(e.Add, s2Map[hc])
		}
	}
	if len(e.Remove) > 0 || len(e.Add) > 0 {
		d = append(d, e)
	}
	return d
}

func (s jsonSet) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonSet) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {

	// Merge patch strategy
	if strategy == mergePatchStrategy {
		return patch(s, pathBehind, pathAhead, before, oldValues, newValues, after, mergePatchStrategy)
	}

	// Strict patch strategy
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
	n, o, rest := pathAhead.next()
	opts := refine(&options{retain: o}, nil)

	pathSetKeys, ok := n.(PathSetKeys)
	if ok && len(rest) > 0 {
		// Recurse into a specific object.
		lookingFor := jsonObject(pathSetKeys).ident(opts)
		for _, v := range s {
			if o, ok := v.(jsonObject); ok {
				id := o.pathIdent(jsonObject(pathSetKeys), opts)
				if id == lookingFor {
					v.patch(append(pathBehind, n), rest, before, oldValues, newValues, after, strategy)
					return s, nil
				}
			}
		}
		return nil, fmt.Errorf("invalid diff: expected object with id %v but found none", jsonObject(pathSetKeys).Json())
	}
	_, ok = n.(PathSet)
	if !ok {
		return nil, fmt.Errorf(
			"invalid path element %v: expected jsonObject", n)
	}
	// Patch set
	aMap := make(map[[8]byte]JsonNode)
	for _, v := range s {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identitiy.
			hc = o.ident(opts)
		} else {
			// Everything else by full content.
			hc = v.hashCode(opts)
		}
		aMap[hc] = v
	}
	for _, v := range oldValues {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Find objects by their identitiy.
			hc = o.ident(opts)
		} else {
			// Everything else by full content.
			hc = v.hashCode(opts)
		}
		toDelete, ok := aMap[hc]
		if !ok {
			return nil, fmt.Errorf(
				"invalid diff: expected %v at %v but found nothing",
				v.Json(), pathBehind)
		}
		if !toDelete.equals(v, opts) {
			return nil, fmt.Errorf(
				"invalid diff: expected %v at %v but found %v",
				v.Json(), pathBehind, toDelete.Json())

		}
		delete(aMap, hc)
	}
	for _, v := range newValues {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identitiy.
			hc = o.ident(opts)
		} else {
			// Everything else by full content.
			hc = v.hashCode(opts)
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
