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
	o := refine(newOptions(opts), nil)
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
		processor := newsetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Handle merge patch strategy for same types
	if strategy == mergePatchStrategy && !s1.Equals(n) {
		events := generateSimpleEvents(s1, n, opts)
		processor := newsetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Same type - use set-specific event generation
	events := generateSetdiffEvents(s1, s2, opts)
	processor := newsetDiffProcessor(path, opts, strategy)
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
		if !toDelete.equals(v, opts) { //jd:nocover — requires ident hash collision
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

// ============================================================================
// SET-SPECIFIC EVENT-DRIVEN DIFF SYSTEM
// ============================================================================

// setElementEvent represents operations on set elements
type setElementEvent struct {
	Operation string // "ADD" or "REMOVE"
	Element   JsonNode
	Hash      [8]byte // For identity tracking
}

func (e setElementEvent) String() string {
	return fmt.Sprintf("SET_%s(%s)", e.Operation, e.Element.Json())
}

func (e setElementEvent) GetType() string { return "SET_ELEMENT" }

// setObjectDiffEvent represents an object in a set that needs recursive diffing
type setObjectDiffEvent struct {
	OldObject JsonNode
	NewObject JsonNode
	Hash      [8]byte // Identity hash
}

func (e setObjectDiffEvent) String() string {
	return fmt.Sprintf("SET_OBJECT_DIFF(%s -> %s)", e.OldObject.Json(), e.NewObject.Json())
}

func (e setObjectDiffEvent) GetType() string { return "SET_OBJECT_DIFF" }

// setDiffProcessor processes set diff events
type setDiffProcessor struct {
	*baseDiffProcessor
}

func newsetDiffProcessor(path Path, opts *options, strategy patchStrategy) *setDiffProcessor {
	return &setDiffProcessor{
		baseDiffProcessor: newBaseDiffProcessor(path, opts, strategy),
	}
}

func (p *setDiffProcessor) ProcessEvents(events []diffEvent) Diff {
	p.debugLog("Starting to process %d set events", len(events))

	// Collect all add/remove events for the set operation
	var setElement *DiffElement

	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())

		switch e := event.(type) {
		case setElementEvent:
			if setElement == nil {
				setElement = &DiffElement{
					Path:   append(p.path.clone(), PathSet{}),
					Remove: []JsonNode{},
					Add:    []JsonNode{},
				}
			}

			if e.Operation == "REMOVE" {
				setElement.Remove = append(setElement.Remove, e.Element)
			} else if e.Operation == "ADD" {
				setElement.Add = append(setElement.Add, e.Element)
			}

		case setObjectDiffEvent:
			p.processsetObjectDiffEvent(e)

		case simpleReplaceEvent:
			p.processsimpleReplaceEvent(e)

		default:
			p.debugLog("WARNING: Unknown event type for setDiffProcessor: %T", event)
		}
	}

	// Add the accumulated set diff element if it has changes
	if setElement != nil && (len(setElement.Remove) > 0 || len(setElement.Add) > 0) {
		p.finalDiff = append(p.finalDiff, *setElement)
	}

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *setDiffProcessor) processsetObjectDiffEvent(event setObjectDiffEvent) {
	p.debugLog("Processing set object diff: %s -> %s", event.OldObject.Json(), event.NewObject.Json())

	// For set object diffs, we need to create a path with PathSetKeys
	o1, _ := event.OldObject.(jsonObject)
	setKeysPath := newPathSetKeys(o1, p.opts)
	subPath := append(p.path.clone(), setKeysPath)

	subDiff := event.OldObject.diff(event.NewObject, subPath, p.opts, p.strategy)
	p.finalDiff = append(p.finalDiff, subDiff...)
}

func (p *setDiffProcessor) processsimpleReplaceEvent(event simpleReplaceEvent) {
	p.debugLog("Processing simple replace: %s -> %s", event.OldValue.Json(), event.NewValue.Json())

	var e DiffElement
	switch p.strategy {
	case mergePatchStrategy:
		e = DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: p.path.clone(),
			Add:  []JsonNode{event.NewValue},
		}
	default:
		e = DiffElement{
			Path:   p.path.clone(),
			Remove: []JsonNode{event.OldValue},
			Add:    []JsonNode{event.NewValue},
		}
	}

	p.finalDiff = append(p.finalDiff, e)
}

// generateSetdiffEvents analyzes two sets and generates appropriate diff events
func generateSetdiffEvents(s1, s2 jsonSet, opts *options) []diffEvent {
	if !opts.diffingOn {
		return []diffEvent{} // No events when diffing is off
	}

	var events []diffEvent

	// Create hash maps for identity-based comparison
	s1Map := make(map[[8]byte]JsonNode)
	for _, v := range s1 {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identity
			hc = o.ident(opts)
		} else {
			// Everything else by full content
			hc = v.hashCode(opts)
		}
		s1Map[hc] = v
	}

	s2Map := make(map[[8]byte]JsonNode)
	for _, v := range s2 {
		var hc [8]byte
		if o, ok := v.(jsonObject); ok {
			// Hash objects by their identity
			hc = o.ident(opts)
		} else {
			// Everything else by full content
			hc = v.hashCode(opts)
		}
		s2Map[hc] = v
	}

	// Get sorted hash codes for deterministic ordering
	s1Hashes := make(hashCodes, 0, len(s1Map))
	for hc := range s1Map {
		s1Hashes = append(s1Hashes, hc)
	}
	sort.Sort(s1Hashes)

	s2Hashes := make(hashCodes, 0, len(s2Map))
	for hc := range s2Map {
		s2Hashes = append(s2Hashes, hc)
	}
	sort.Sort(s2Hashes)

	// Process removes first (sorted by hash)
	for _, hc := range s1Hashes {
		v1 := s1Map[hc]
		if v2, ok := s2Map[hc]; !ok {
			// Deleted value
			events = append(events, setElementEvent{
				Operation: "REMOVE",
				Element:   v1,
				Hash:      hc,
			})
		} else {
			// Check for object diffs with same identity
			o1, isObject1 := v1.(jsonObject)
			o2, isObject2 := v2.(jsonObject)
			if isObject1 && isObject2 && !o1.equals(o2, opts) {
				// Sub diff objects with same identity
				events = append(events, setObjectDiffEvent{
					OldObject: o1,
					NewObject: o2,
					Hash:      hc,
				})
			}
		}
	}

	// Process adds (sorted by hash)
	for _, hc := range s2Hashes {
		if _, ok := s1Map[hc]; !ok {
			// Added value
			events = append(events, setElementEvent{
				Operation: "ADD",
				Element:   s2Map[hc],
				Hash:      hc,
			})
		}
	}

	return events
}
