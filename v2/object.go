package jd

import (
	"fmt"
	"sort"
)

type jsonObject map[string]JsonNode

var _ JsonNode = jsonObject{}

func newJsonObject() jsonObject {
	return jsonObject{}
}

func (o jsonObject) Json(_ ...Option) string {
	return renderJson(o.raw())
}

func (o jsonObject) MarshalJSON() ([]byte, error) {
	return []byte(o.Json()), nil
}

func (o jsonObject) Yaml(_ ...Option) string {
	return renderYaml(o.raw())
}

func (o jsonObject) raw() interface{} {
	j := make(map[string]interface{})
	for k, v := range o {
		j[k] = v.raw()
	}
	return j
}

func (o1 jsonObject) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return o1.equals(n, o)
}

func (o1 jsonObject) equals(n JsonNode, o *options) bool {
	o2, ok := n.(jsonObject)
	if !ok {
		return false
	}
	if len(o1) != len(o2) {
		return false
	}

	for key1, val1 := range o1 {
		val2, ok := o2[key1]
		if !ok {
			return false
		}
		ret := val1.equals(val2, o)
		if !ret {
			return false
		}
	}
	return true
}

func (o jsonObject) hashCode(opts *options) [8]byte {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	a := []byte{0x00, 0x5D, 0x39, 0xA4, 0x18, 0x10, 0xEA, 0xD5} // random bytes
	for _, k := range keys {
		keyHash := hash([]byte(k))
		a = append(a, keyHash[:]...)
		valueHash := o[k].hashCode(opts)
		a = append(a, valueHash[:]...)
	}
	return hash(a)
}

// ident is the identity of the json object based on either the hash of a
// given set of keys or the full object if no keys are present.
func (o jsonObject) ident(opts *options) [8]byte {
	keys, ok := getOption[setKeysOption](opts)
	if !ok {
		return o.hashCode(opts)
	}
	hashes := hashCodes{
		// We start with a constant hash to distinguish between
		// an empty object and an empty array.
		[8]byte{0x4B, 0x08, 0xD2, 0x0F, 0xBD, 0xC8, 0xDE, 0x9A}, // randomly chosen bytes
	}
	for _, key := range []string(*keys) {
		v, ok := o[key]
		if ok {
			hashes = append(hashes, v.hashCode(opts))
		}
	}
	if len(hashes) == 0 {
		return o.hashCode(opts)
	}
	return hashes.combine()
}

func (o jsonObject) pathIdent(pathObject jsonObject, opts *options) [8]byte {
	keys := []string{}
	for k := range pathObject {
		keys = append(keys, k)
	}
	id := make(map[string]interface{})
	for _, key := range keys {
		if value, ok := o[key]; ok {
			id[key] = value
		}
	}
	e, _ := NewJsonNode(id)
	return e.hashCode(&options{})
}

func (o jsonObject) Diff(n JsonNode, opts ...Option) Diff {
	op := &options{retain: opts}
	return o.diff(n, make(Path, 0), op, getPatchStrategy(op))
}

func (o1 jsonObject) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	o2, ok := n.(jsonObject)
	if !ok {
		// Different types - use simple replace event
		events := generateSimpleEvents(o1, n, opts)
		processor := NewObjectDiffProcessor(path, opts, strategy)
		return processor.ProcessEvents(events)
	}

	// Same type - use object-specific event generation
	events := generateObjectDiffEvents(o1, o2, opts)
	processor := NewObjectDiffProcessor(path, opts, strategy)
	return processor.ProcessEvents(events)
}

func (o jsonObject) Patch(d Diff) (JsonNode, error) {
	return patchAll(o, d)
}

func (o jsonObject) patch(
	pathBehind, pathAhead Path,
	before, oldValues, newValues, after []JsonNode,
	strategy patchStrategy,
) (JsonNode, error) {
	if (len(pathAhead) == 0) && (len(oldValues) > 1 || len(newValues) > 1) {
		return patchErrNonSetDiff(oldValues, newValues, pathBehind)
	}
	// Base case
	if len(pathAhead) == 0 {
		newValue := singleValue(newValues)
		if strategy == mergePatchStrategy {
			return newValue, nil
		}
		oldValue := singleValue(oldValues)
		if !o.Equals(oldValue) {
			return patchErrExpectValue(oldValue, o, pathBehind)
		}
		return newValue, nil
	}
	// Recursive case
	n, _, rest := pathAhead.next()
	pe, ok := n.(PathKey)
	if !ok {
		return nil, fmt.Errorf(
			"found %v at %v: expected JSON object",
			o.Json(), pathBehind)
	}
	nextNode, ok := o[string(pe)]
	if !ok {
		switch strategy {
		case mergePatchStrategy:
			// Create objects
			if len(rest) == 0 {
				nextNode = voidNode{}
			} else {
				nextNode = newJsonObject()
			}
		case strictPatchStrategy:
			nextNode = voidNode{}
		default:
			return patchErrUnsupportedPatchStrategy(pathBehind, strategy)
		}
	}
	patchedNode, err := nextNode.patch(append(pathBehind, pe), rest, before, oldValues, newValues, after, strategy)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		// Delete a pair
		delete(o, string(pe))
	} else {
		// Add or replace a pair
		o[string(pe)] = patchedNode
	}
	return o, nil
}

// ============================================================================
// OBJECT-SPECIFIC EVENT-DRIVEN DIFF SYSTEM
// ============================================================================

// ObjectKeyEvent represents operations on object keys
type ObjectKeyEvent interface {
	DiffEvent
	GetKey() string
}

// ObjectKeyAddEvent represents adding a new key to an object
type ObjectKeyAddEvent struct {
	Key   string
	Value JsonNode
}

func (e ObjectKeyAddEvent) String() string {
	return fmt.Sprintf("OBJECT_KEY_ADD(%s: %s)", e.Key, e.Value.Json())
}

func (e ObjectKeyAddEvent) GetType() string { return "OBJECT_KEY_ADD" }
func (e ObjectKeyAddEvent) GetKey() string  { return e.Key }

// ObjectKeyRemoveEvent represents removing a key from an object
type ObjectKeyRemoveEvent struct {
	Key   string
	Value JsonNode
}

func (e ObjectKeyRemoveEvent) String() string {
	return fmt.Sprintf("OBJECT_KEY_REMOVE(%s: %s)", e.Key, e.Value.Json())
}

func (e ObjectKeyRemoveEvent) GetType() string { return "OBJECT_KEY_REMOVE" }
func (e ObjectKeyRemoveEvent) GetKey() string  { return e.Key }

// ObjectKeyDiffEvent represents a key that exists in both objects but with different values
type ObjectKeyDiffEvent struct {
	Key         string
	OldValue    JsonNode
	NewValue    JsonNode
	IsRecursive bool // true if values are compatible containers
}

func (e ObjectKeyDiffEvent) String() string {
	if e.IsRecursive {
		return fmt.Sprintf("OBJECT_KEY_DIFF_RECURSIVE(%s: %s -> %s)", e.Key, e.OldValue.Json(), e.NewValue.Json())
	}
	return fmt.Sprintf("OBJECT_KEY_DIFF(%s: %s -> %s)", e.Key, e.OldValue.Json(), e.NewValue.Json())
}

func (e ObjectKeyDiffEvent) GetType() string { return "OBJECT_KEY_DIFF" }
func (e ObjectKeyDiffEvent) GetKey() string  { return e.Key }

// ObjectDiffProcessor processes object diff events
type ObjectDiffProcessor struct {
	*BaseDiffProcessor
}

func NewObjectDiffProcessor(path Path, opts *options, strategy patchStrategy) *ObjectDiffProcessor {
	return &ObjectDiffProcessor{
		BaseDiffProcessor: NewBaseDiffProcessor(path, opts, strategy),
	}
}

func (p *ObjectDiffProcessor) ProcessEvents(events []DiffEvent) Diff {
	p.debugLog("Starting to process %d object events", len(events))

	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())
		p.processEvent(event)
	}

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *ObjectDiffProcessor) processEvent(event DiffEvent) {
	switch e := event.(type) {
	case ObjectKeyAddEvent:
		p.processKeyAddEvent(e)
	case ObjectKeyRemoveEvent:
		p.processKeyRemoveEvent(e)
	case ObjectKeyDiffEvent:
		p.processKeyDiffEvent(e)
	case SimpleReplaceEvent:
		p.processSimpleReplaceEvent(e)
	default:
		p.debugLog("WARNING: Unknown event type for ObjectDiffProcessor: %T", event)
	}
}

func (p *ObjectDiffProcessor) processKeyAddEvent(event ObjectKeyAddEvent) {
	p.debugLog("Processing key add: %s = %s", event.Key, event.Value.Json())

	var e DiffElement
	keyPath := append(p.path, PathKey(event.Key))
	switch p.strategy {
	case mergePatchStrategy:
		e = DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path:   keyPath.clone(),
			Remove: []JsonNode{},
			Add:    []JsonNode{event.Value},
		}
	default:
		e = DiffElement{
			Path:   keyPath.clone(),
			Remove: []JsonNode{},
			Add:    []JsonNode{event.Value},
		}
	}
	p.finalDiff = append(p.finalDiff, e)
}

func (p *ObjectDiffProcessor) processKeyRemoveEvent(event ObjectKeyRemoveEvent) {
	p.debugLog("Processing key remove: %s = %s", event.Key, event.Value.Json())

	var e DiffElement
	keyPath := append(p.path, PathKey(event.Key))
	switch p.strategy {
	case mergePatchStrategy:
		e = DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: keyPath.clone(),
			Add:  []JsonNode{voidNode{}},
		}
	default:
		e = DiffElement{
			Path:   keyPath.clone(),
			Remove: []JsonNode{event.Value},
			Add:    []JsonNode{},
		}
	}
	p.finalDiff = append(p.finalDiff, e)
}

func (p *ObjectDiffProcessor) processKeyDiffEvent(event ObjectKeyDiffEvent) {
	p.debugLog("Processing key diff: %s = %s -> %s (recursive=%t)",
		event.Key, event.OldValue.Json(), event.NewValue.Json(), event.IsRecursive)

	keyPath := append(p.path, PathKey(event.Key))

	if event.IsRecursive {
		// Recursive diff for compatible containers
		refinedOpts := refine(p.opts, PathKey(event.Key))
		subDiff := event.OldValue.diff(event.NewValue, keyPath, refinedOpts, p.strategy)
		p.finalDiff = append(p.finalDiff, subDiff...)
	} else {
		// Simple replacement
		var e DiffElement
		switch p.strategy {
		case mergePatchStrategy:
			e = DiffElement{
				Metadata: Metadata{
					Merge: true,
				},
				Path: keyPath.clone(),
				Add:  []JsonNode{event.NewValue},
			}
		default:
			e = DiffElement{
				Path:   keyPath.clone(),
				Remove: []JsonNode{event.OldValue},
				Add:    []JsonNode{event.NewValue},
			}
		}
		p.finalDiff = append(p.finalDiff, e)
	}
}

func (p *ObjectDiffProcessor) processSimpleReplaceEvent(event SimpleReplaceEvent) {
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

// generateObjectDiffEvents analyzes two objects and generates appropriate diff events
func generateObjectDiffEvents(o1, o2 jsonObject, opts *options) []DiffEvent {
	var events []DiffEvent

	// Get all unique keys and sort them for deterministic processing
	allKeys := make(map[string]bool)
	for k := range o1 {
		allKeys[k] = true
	}
	for k := range o2 {
		allKeys[k] = true
	}

	sortedKeys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// Process all keys in sorted order
	for _, k := range sortedKeys {
		v1, existsInO1 := o1[k]
		v2, existsInO2 := o2[k]

		if existsInO1 && existsInO2 {
			// Both keys are present - check if they're different
			o := refine(opts, PathKey(k))
			if !v1.equals(v2, o) {
				// Check if compatible containers for recursive diff
				isRecursive := sameContainerType(v1, v2, opts)
				events = append(events, ObjectKeyDiffEvent{
					Key:         k,
					OldValue:    v1,
					NewValue:    v2,
					IsRecursive: isRecursive,
				})
			}
			// If equal, no event needed
		} else if existsInO1 {
			// Key only in o1 - removal
			events = append(events, ObjectKeyRemoveEvent{
				Key:   k,
				Value: v1,
			})
		} else {
			// Key only in o2 - addition
			events = append(events, ObjectKeyAddEvent{
				Key:   k,
				Value: v2,
			})
		}
	}

	return events
}
