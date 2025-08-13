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
	s2, ok := n.(jsonSet)
	if !ok {
		// Different types - use simple replace event
		events := generateSimpleEvents(s1, n, opts)
		processor := NewSetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Handle merge patch strategy for same types
	if strategy == mergePatchStrategy && !s1.Equals(n) {
		events := generateSimpleEvents(s1, n, opts)
		processor := NewSetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Same type - use set-specific event generation
	events := generateSetDiffEvents(s1, s2, opts)
	processor := NewSetDiffProcessor(path, opts, strategy)
	return processor.ProcessEvents(events)
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
