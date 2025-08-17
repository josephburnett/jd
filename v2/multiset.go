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
	o := refine(newOptions(opts), nil)
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
		processor := newmultisetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Handle merge patch strategy for same types
	if strategy == mergePatchStrategy && !a1.Equals(n) {
		events := generateSimpleEvents(a1, n, opts)
		processor := newmultisetDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Same type - use multiset-specific event generation
	events := generateMultisetdiffEvents(a1, a2, opts)
	processor := newmultisetDiffProcessor(path, opts, strategy)
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

// ============================================================================
// MULTISET-SPECIFIC EVENT-DRIVEN DIFF SYSTEM
// ============================================================================

// multisetElementEvent represents operations on multiset elements with counts
type multisetElementEvent struct {
	Operation string // "ADD" or "REMOVE"
	Element   JsonNode
	Count     int     // How many instances to add/remove
	Hash      [8]byte // For identity tracking
}

func (e multisetElementEvent) String() string {
	return fmt.Sprintf("MULTISET_%s(%s x%d)", e.Operation, e.Element.Json(), e.Count)
}

func (e multisetElementEvent) GetType() string { return "MULTISET_ELEMENT" }

// multisetDiffProcessor processes multiset diff events
type multisetDiffProcessor struct {
	*baseDiffProcessor
}

func newmultisetDiffProcessor(path Path, opts *options, strategy patchStrategy) *multisetDiffProcessor {
	return &multisetDiffProcessor{
		baseDiffProcessor: newBaseDiffProcessor(path, opts, strategy),
	}
}

func (p *multisetDiffProcessor) ProcessEvents(events []diffEvent) Diff {
	p.debugLog("Starting to process %d multiset events", len(events))

	// Collect all add/remove events for the multiset operation
	var multisetElement *DiffElement

	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())

		switch e := event.(type) {
		case multisetElementEvent:
			if multisetElement == nil {
				multisetElement = &DiffElement{
					Path:   append(p.path.clone(), PathMultiset{}),
					Remove: []JsonNode{},
					Add:    []JsonNode{},
				}
			}

			// Add the element the specified number of times
			for i := 0; i < e.Count; i++ {
				if e.Operation == "REMOVE" {
					multisetElement.Remove = append(multisetElement.Remove, e.Element)
				} else if e.Operation == "ADD" {
					multisetElement.Add = append(multisetElement.Add, e.Element)
				}
			}

		case simpleReplaceEvent:
			p.processsimpleReplaceEvent(e)

		default:
			p.debugLog("WARNING: Unknown event type for multisetDiffProcessor: %T", event)
		}
	}

	// Add the accumulated multiset diff element if it has changes
	if multisetElement != nil && (len(multisetElement.Remove) > 0 || len(multisetElement.Add) > 0) {
		p.finalDiff = append(p.finalDiff, *multisetElement)
	}

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *multisetDiffProcessor) processsimpleReplaceEvent(event simpleReplaceEvent) {
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

// generateMultisetdiffEvents analyzes two multisets and generates appropriate diff events
func generateMultisetdiffEvents(a1, a2 jsonMultiset, opts *options) []diffEvent {
	if !opts.diffingOn {
		return []diffEvent{} // No events when diffing is off
	}

	var events []diffEvent

	// Count elements in both multisets
	a1Counts := make(map[[8]byte]int)
	a1Map := make(map[[8]byte]JsonNode)
	for _, v := range a1 {
		hc := v.hashCode(opts)
		a1Counts[hc]++
		a1Map[hc] = v
	}

	a2Counts := make(map[[8]byte]int)
	a2Map := make(map[[8]byte]JsonNode)
	for _, v := range a2 {
		hc := v.hashCode(opts)
		a2Counts[hc]++
		a2Map[hc] = v
	}

	// Get sorted hash codes for deterministic ordering (matches original implementation)
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

	// Process removals first (sorted by hash)
	for _, hc := range a1Hashes {
		a1Count := a1Counts[hc]
		a2Count, ok := a2Counts[hc]
		if !ok {
			a2Count = 0
		}
		removed := a1Count - a2Count
		if removed > 0 {
			events = append(events, multisetElementEvent{
				Operation: "REMOVE",
				Element:   a1Map[hc],
				Count:     removed,
				Hash:      hc,
			})
		}
	}

	// Process additions (sorted by hash)
	for _, hc := range a2Hashes {
		a2Count := a2Counts[hc]
		a1Count, ok := a1Counts[hc]
		if !ok {
			a1Count = 0
		}
		added := a2Count - a1Count
		if added > 0 {
			events = append(events, multisetElementEvent{
				Operation: "ADD",
				Element:   a2Map[hc],
				Count:     added,
				Hash:      hc,
			})
		}
	}

	return events
}
