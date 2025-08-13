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

func (a1 jsonMultiset) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return a1.equals(n, o)
}

func (a1 jsonMultiset) equals(n JsonNode, o *options) bool {
	n = dispatch(n, o)
	a2, ok := n.(jsonMultiset)
	if !ok {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	if a1.hashCode(o) == a2.hashCode(o) {
		return true
	} else {
		return false
	}
}

func (a jsonMultiset) hashCode(opts *options) [8]byte {
	h := make(hashCodes, 0, len(a))
	for _, v := range a {
		h = append(h, v.hashCode(opts))
	}
	sort.Sort(h)
	b := make([]byte, 0, len(a)*8)
	for _, c := range h {
		b = append(b, c[:]...)
	}
	return hash(b)
}

func (a jsonMultiset) Diff(n JsonNode, opts ...Option) Diff {
	o := refine(&options{retain: opts}, nil)
	return a.diff(n, nil, o, getPatchStrategy(o))
}

func (a1 jsonMultiset) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	a2, ok := n.(jsonMultiset)
	if !ok {
		// Different types - use simple replace event
		events := generateSimpleEvents(a1, n, opts)
		processor := NewMultisetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Handle merge patch strategy for same types
	if strategy == mergePatchStrategy && !a1.Equals(n) {
		events := generateSimpleEvents(a1, n, opts)
		processor := NewMultisetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Same type - use multiset-specific event generation
	events := generateMultisetDiffEvents(a1, a2, opts)
	processor := NewMultisetDiffProcessor(path, opts, strategy)
	return processor.ProcessEvents(events)
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
	n, opts, _ := pathAhead.next()
	o := refine(&options{retain: opts}, nil)
	_, ok := n.(PathMultiset)
	if !ok {
		return nil, fmt.Errorf(
			"invalid path element %v: expected map[string]interface{}", n)
	}
	aCounts := make(map[[8]byte]int)
	aMap := make(map[[8]byte]JsonNode)
	for _, v := range a {
		hc := v.hashCode(o)
		aCounts[hc]++
		aMap[hc] = v
	}
	for _, v := range oldValues {
		hc := v.hashCode(o)
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
		hc := v.hashCode(o)
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
