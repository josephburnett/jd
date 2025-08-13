package jd

import (
	"fmt"
	"sort"
)

// ============================================================================
// GENERALIZED EVENT-DRIVEN DIFF SYSTEM
// ============================================================================

// DiffEvent represents an operation needed to transform one structure into another
type DiffEvent interface {
	String() string // For debugging
	GetType() string
}

// MatchEvent represents elements that are identical between A and B
type MatchEvent struct {
	AIndex, BIndex int
	Element        JsonNode
}

func (e MatchEvent) String() string {
	return fmt.Sprintf("MATCH(A[%d]=B[%d]: %s)", e.AIndex, e.BIndex, e.Element.Json())
}

func (e MatchEvent) GetType() string { return "MATCH" }

// ContainerDiffEvent represents containers that are compatible and need recursive diffing
type ContainerDiffEvent struct {
	AIndex, BIndex     int
	AElement, BElement JsonNode
}

func (e ContainerDiffEvent) String() string {
	return fmt.Sprintf("CONTAINER_DIFF(A[%d] vs B[%d])", e.AIndex, e.BIndex)
}

func (e ContainerDiffEvent) GetType() string { return "CONTAINER_DIFF" }

// RemoveEvent represents an element that exists only in A (needs to be removed)
type RemoveEvent struct {
	AIndex  int
	Element JsonNode
}

func (e RemoveEvent) String() string {
	return fmt.Sprintf("REMOVE(A[%d]: %s)", e.AIndex, e.Element.Json())
}

func (e RemoveEvent) GetType() string { return "REMOVE" }

// AddEvent represents an element that exists only in B (needs to be added)
type AddEvent struct {
	BIndex  int
	Element JsonNode
}

func (e AddEvent) String() string {
	return fmt.Sprintf("ADD(B[%d]: %s)", e.BIndex, e.Element.Json())
}

func (e AddEvent) GetType() string { return "ADD" }

// ReplaceEvent represents elements at the same position that are different
type ReplaceEvent struct {
	AIndex, BIndex     int
	AElement, BElement JsonNode
}

func (e ReplaceEvent) String() string {
	return fmt.Sprintf("REPLACE(A[%d]: %s -> B[%d]: %s)",
		e.AIndex, e.AElement.Json(), e.BIndex, e.BElement.Json())
}

func (e ReplaceEvent) GetType() string { return "REPLACE" }

// SimpleReplaceEvent represents a simple replacement between two different values
type SimpleReplaceEvent struct {
	OldValue, NewValue JsonNode
}

func (e SimpleReplaceEvent) String() string {
	return fmt.Sprintf("SIMPLE_REPLACE(%s -> %s)", e.OldValue.Json(), e.NewValue.Json())
}

func (e SimpleReplaceEvent) GetType() string { return "SIMPLE_REPLACE" }

// ============================================================================
// GENERALIZED DIFF PROCESSOR - STATE MACHINE
// ============================================================================

// DiffProcessorState represents the current state of diff processing
type DiffProcessorState int

const (
	IDLE DiffProcessorState = iota
	ACCUMULATING_CHANGES
	PROCESSING_MATCH
	PROCESSING_CONTAINER_DIFF
)

func (s DiffProcessorState) String() string {
	switch s {
	case IDLE:
		return "IDLE"
	case ACCUMULATING_CHANGES:
		return "ACCUMULATING_CHANGES"
	case PROCESSING_MATCH:
		return "PROCESSING_MATCH"
	case PROCESSING_CONTAINER_DIFF:
		return "PROCESSING_CONTAINER_DIFF"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

// BaseDiffProcessor provides common functionality for all diff processors
type BaseDiffProcessor struct {
	state       DiffProcessorState
	currentDiff DiffElement
	finalDiff   Diff
	opts        *options
	strategy    patchStrategy
	path        Path
	debug       bool
}

func NewBaseDiffProcessor(path Path, opts *options, strategy patchStrategy) *BaseDiffProcessor {
	return &BaseDiffProcessor{
		state:     IDLE,
		finalDiff: Diff{}, // Initialize to empty slice, not nil
		opts:      opts,
		strategy:  strategy,
		path:      path,
		debug:     false,
	}
}

func (p *BaseDiffProcessor) SetDebug(debug bool) {
	p.debug = debug
}

func (p *BaseDiffProcessor) debugLog(format string, args ...interface{}) {
	if p.debug {
		fmt.Printf("[BaseDiffProcessor:%s] "+format+"\n", append([]interface{}{p.state}, args...)...)
	}
}

// ============================================================================
// SIMPLE DIFF PROCESSOR (for primitive types)
// ============================================================================

// SimpleDiffProcessor handles diff processing for simple types (bool, string, number, null, void)
type SimpleDiffProcessor struct {
	*BaseDiffProcessor
}

func NewSimpleDiffProcessor(path Path, opts *options, strategy patchStrategy) *SimpleDiffProcessor {
	return &SimpleDiffProcessor{
		BaseDiffProcessor: NewBaseDiffProcessor(path, opts, strategy),
	}
}

func (p *SimpleDiffProcessor) ProcessEvents(events []DiffEvent) Diff {
	p.debugLog("Starting to process %d events", len(events))

	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())
		p.processEvent(event)
	}

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *SimpleDiffProcessor) processEvent(event DiffEvent) {
	switch e := event.(type) {
	case SimpleReplaceEvent:
		p.processSimpleReplaceEvent(e)
	default:
		p.debugLog("WARNING: Unknown event type for SimpleDiffProcessor: %T", event)
	}
}

func (p *SimpleDiffProcessor) processSimpleReplaceEvent(event SimpleReplaceEvent) {
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

// generateSimpleEvents creates events for simple type differences
func generateSimpleEvents(a, b JsonNode, opts *options) []DiffEvent {
	if a.equals(b, opts) {
		return []DiffEvent{}
	}
	return []DiffEvent{SimpleReplaceEvent{OldValue: a, NewValue: b}}
}

// ============================================================================
// LIST DIFF PROCESSOR - STATE MACHINE
// ============================================================================

// ListDiffProcessorState represents the current state of diff processing
type ListDiffProcessorState int

const (
	LIST_IDLE ListDiffProcessorState = iota
	LIST_ACCUMULATING_CHANGES
	LIST_PROCESSING_MATCH
	LIST_PROCESSING_CONTAINER_DIFF
)

func (s ListDiffProcessorState) String() string {
	switch s {
	case LIST_IDLE:
		return "LIST_IDLE"
	case LIST_ACCUMULATING_CHANGES:
		return "LIST_ACCUMULATING_CHANGES"
	case LIST_PROCESSING_MATCH:
		return "LIST_PROCESSING_MATCH"
	case LIST_PROCESSING_CONTAINER_DIFF:
		return "LIST_PROCESSING_CONTAINER_DIFF"
	default:
		return fmt.Sprintf("LIST_UNKNOWN(%d)", s)
	}
}

// ListDiffProcessor processes diff events using a state machine
type ListDiffProcessor struct {
	state       ListDiffProcessorState
	currentDiff DiffElement
	finalDiff   Diff
	opts        *options
	strategy    patchStrategy
	previous    JsonNode

	// Helper components
	pathCalc   *PathCalculator
	contextMgr *ContextManager

	// Context tracking
	lastProcessedElement JsonNode
	nextAfterElement     JsonNode // The element that should appear as After context

	// Debug info
	debug bool
}

func NewListDiffProcessor(basePath Path, pathIndex PathIndex, previous JsonNode, a, b jsonList, opts *options, strategy patchStrategy) *ListDiffProcessor {
	return &ListDiffProcessor{
		state:                LIST_IDLE,
		finalDiff:            Diff{}, // Initialize to empty slice, not nil
		opts:                 opts,
		strategy:             strategy,
		previous:             previous,
		pathCalc:             NewPathCalculator(basePath, pathIndex),
		contextMgr:           NewContextManager(a, b),
		lastProcessedElement: previous,
		nextAfterElement:     voidNode{}, // Default to void
		debug:                false,      // Disable debugging
	}
}

func (p *ListDiffProcessor) SetDebug(debug bool) {
	p.debug = debug
	p.pathCalc.SetDebug(debug)
	p.contextMgr.SetDebug(debug)
}

func (p *ListDiffProcessor) debugLog(format string, args ...interface{}) {
	if p.debug {
		fmt.Printf("[ListDiffProcessor:%s] "+format+"\n", append([]interface{}{p.state}, args...)...)
	}
}

func (p *ListDiffProcessor) ProcessEvents(events []DiffEvent) Diff {
	p.debugLog("Starting to process %d events", len(events))

	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())
		p.processEvent(event)
	}

	// Finalize any accumulated changes
	p.finalize()

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *ListDiffProcessor) processEvent(event DiffEvent) {
	switch e := event.(type) {
	case MatchEvent:
		p.processMatchEvent(e)
	case ContainerDiffEvent:
		p.processContainerDiffEvent(e)
	case RemoveEvent:
		p.processRemoveEvent(e)
	case AddEvent:
		p.processAddEvent(e)
	case ReplaceEvent:
		p.processReplaceEvent(e)
	default:
		p.debugLog("WARNING: Unknown event type: %T", event)
	}
}

func (p *ListDiffProcessor) processMatchEvent(event MatchEvent) {
	p.debugLog("Processing match: A[%d]=B[%d]", event.AIndex, event.BIndex)

	// If we were accumulating changes, finalize them first
	if p.state == LIST_ACCUMULATING_CHANGES {
		// The matched element becomes the After context for the previous operation
		p.nextAfterElement = event.Element
		p.finalizeCurrentDiff()
	}

	p.state = LIST_PROCESSING_MATCH
	p.pathCalc.AdvanceForMatch()
	p.contextMgr.SetCursors(event.AIndex+1, event.BIndex+1) // Move past the match
	p.lastProcessedElement = event.Element                  // Track the matched element
	p.state = LIST_IDLE
}

func (p *ListDiffProcessor) processContainerDiffEvent(event ContainerDiffEvent) {
	p.debugLog("Processing container diff: A[%d] vs B[%d]", event.AIndex, event.BIndex)

	// If we were accumulating changes, finalize them first
	if p.state == LIST_ACCUMULATING_CHANGES {
		// The After context should be the element from original array A that comes after the previous operation
		p.nextAfterElement = event.AElement // Use A element, not B element
		p.finalizeCurrentDiff()
	}

	p.state = LIST_PROCESSING_CONTAINER_DIFF

	// Perform recursive diff
	subPath := p.pathCalc.CurrentPath()
	p.debugLog("Container diff sub-path: %v", subPath)
	// Refine options for the current index
	currentIndex := p.pathCalc.GetCurrentIndex()
	refinedOpts := refine(p.opts, currentIndex)
	subDiff := event.AElement.diff(event.BElement, subPath, refinedOpts, p.strategy)

	p.debugLog("Sub-diff has %d elements with paths:", len(subDiff))
	for i, elem := range subDiff {
		p.debugLog("  Sub-diff %d: path=%v", i, elem.Path)
	}
	p.finalDiff = append(p.finalDiff, subDiff...)

	p.pathCalc.AdvanceForContainer()
	p.contextMgr.SetCursors(event.AIndex+1, event.BIndex+1) // Move past the container
	p.lastProcessedElement = event.BElement                 // After container diff, the B element is what remains
	p.state = LIST_IDLE
}

func (p *ListDiffProcessor) processRemoveEvent(event RemoveEvent) {
	p.debugLog("Processing remove: A[%d]", event.AIndex)
	p.ensureAccumulatingState()
	p.currentDiff.Remove = append(p.currentDiff.Remove, event.Element)
	p.pathCalc.AdvanceForRemove()
	// For removes, advance A cursor but not B cursor
}

func (p *ListDiffProcessor) processAddEvent(event AddEvent) {
	p.debugLog("Processing add: B[%d]", event.BIndex)
	p.ensureAccumulatingState()
	p.currentDiff.Add = append(p.currentDiff.Add, event.Element)
	p.pathCalc.AdvanceForAdd()
	// For adds, advance B cursor but not A cursor
}

func (p *ListDiffProcessor) processReplaceEvent(event ReplaceEvent) {
	p.debugLog("Processing replace: A[%d] -> B[%d]", event.AIndex, event.BIndex)
	p.ensureAccumulatingState()
	p.currentDiff.Remove = append(p.currentDiff.Remove, event.AElement)
	p.currentDiff.Add = append(p.currentDiff.Add, event.BElement)
	p.pathCalc.AdvanceForReplace()
	// For replaces, both cursors advance
}

func (p *ListDiffProcessor) ensureAccumulatingState() {
	if p.state != LIST_ACCUMULATING_CHANGES {
		p.startAccumulating()
	}
}

func (p *ListDiffProcessor) startAccumulating() {
	p.debugLog("Starting to accumulate changes at path %s", p.pathCalc.CurrentPath())
	p.state = LIST_ACCUMULATING_CHANGES

	// Use the last processed element as Before context
	beforeContext := []JsonNode{p.lastProcessedElement}

	p.currentDiff = DiffElement{
		Path:   p.pathCalc.CurrentPath(),
		Before: beforeContext,
	}
}

func (p *ListDiffProcessor) finalizeCurrentDiff() {
	if p.state != LIST_ACCUMULATING_CHANGES {
		return
	}

	if len(p.currentDiff.Remove) > 0 || len(p.currentDiff.Add) > 0 {
		p.debugLog("Finalizing accumulated diff with %d removes, %d adds at path %s, pathIndex=%d",
			len(p.currentDiff.Remove), len(p.currentDiff.Add), p.pathCalc.CurrentPath(), p.pathCalc.pathIndex)

		// Use the pre-calculated After context
		p.currentDiff.After = []JsonNode{p.nextAfterElement}

		p.finalDiff = append(p.finalDiff, p.currentDiff)
	}

	p.currentDiff = DiffElement{}
	p.state = LIST_IDLE
}

func (p *ListDiffProcessor) finalize() {
	if p.state == LIST_ACCUMULATING_CHANGES {
		// When finalizing at the end, After context should be void
		p.nextAfterElement = voidNode{}
		p.finalizeCurrentDiff()
	}
}

// ============================================================================
// PATH CALCULATOR AND CONTEXT MANAGER
// ============================================================================

// PathCalculator handles path calculations and advancement during diff processing
type PathCalculator struct {
	basePath  Path
	pathIndex PathIndex
	debug     bool
}

func NewPathCalculator(basePath Path, pathIndex PathIndex) *PathCalculator {
	return &PathCalculator{
		basePath:  basePath,
		pathIndex: pathIndex,
		debug:     false,
	}
}

func (pc *PathCalculator) SetDebug(debug bool) {
	pc.debug = debug
}

func (pc *PathCalculator) debugLog(format string, args ...interface{}) {
	if pc.debug {
		fmt.Printf("[PathCalculator] "+format+"\n", args...)
	}
}

func (pc *PathCalculator) CurrentPath() Path {
	return append(pc.basePath.clone(), pc.pathIndex)
}

func (pc *PathCalculator) GetCurrentIndex() PathIndex {
	return pc.pathIndex
}

func (pc *PathCalculator) AdvanceForMatch() {
	pc.debugLog("Advancing path for match: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex++
}

func (pc *PathCalculator) AdvanceForAdd() {
	pc.debugLog("Advancing path for add: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex += 1 // Path advances by 1 per added element
}

func (pc *PathCalculator) AdvanceForRemove() {
	pc.debugLog("Advancing path for remove: %d (no change)", pc.pathIndex)
	// Path stays same because list is getting shorter but we're moving to next element
}

func (pc *PathCalculator) AdvanceForReplace() {
	pc.debugLog("Advancing path for replace: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex++
}

func (pc *PathCalculator) AdvanceForContainer() {
	pc.debugLog("Advancing path for container: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex++
}

// ContextManager handles Before/After context calculation for diff elements
type ContextManager struct {
	originalA jsonList
	originalB jsonList
	aIndex    int
	bIndex    int
	debug     bool
}

func NewContextManager(a, b jsonList) *ContextManager {
	return &ContextManager{
		originalA: a,
		originalB: b,
		aIndex:    0,
		bIndex:    0,
		debug:     false,
	}
}

func (cm *ContextManager) SetDebug(debug bool) {
	cm.debug = debug
}

func (cm *ContextManager) debugLog(format string, args ...interface{}) {
	if cm.debug {
		fmt.Printf("[ContextManager] "+format+"\n", args...)
	}
}

func (cm *ContextManager) SetCursors(aIndex, bIndex int) {
	cm.debugLog("Setting cursors: A=%d, B=%d", aIndex, bIndex)
	cm.aIndex = aIndex
	cm.bIndex = bIndex
}

func (cm *ContextManager) GetBeforeContext(previous JsonNode) []JsonNode {
	cm.debugLog("Getting before context with previous: %s", previous.Json())
	return []JsonNode{previous}
}

func (cm *ContextManager) GetAfterContext() []JsonNode {
	// Calculate After context similar to original diffRest logic
	// This is the element that comes after the current operation

	if cm.aIndex >= len(cm.originalA) {
		cm.debugLog("After context: void (A cursor beyond end)")
		return []JsonNode{voidNode{}}
	}

	afterElement := cm.originalA[cm.aIndex]
	cm.debugLog("After context: A[%d] = %s", cm.aIndex, afterElement.Json())
	return []JsonNode{afterElement}
}

func (cm *ContextManager) AdvanceForMatch() {
	cm.debugLog("Advancing context for match: A=%d->%d, B=%d->%d",
		cm.aIndex, cm.aIndex+1, cm.bIndex, cm.bIndex+1)
	cm.aIndex++
	cm.bIndex++
}

func (cm *ContextManager) AdvanceForAdd() {
	cm.debugLog("Advancing context for add: B=%d->%d", cm.bIndex, cm.bIndex+1)
	cm.bIndex++
}

func (cm *ContextManager) AdvanceForRemove() {
	cm.debugLog("Advancing context for remove: A=%d->%d", cm.aIndex, cm.aIndex+1)
	cm.aIndex++
}

func (cm *ContextManager) AdvanceForReplace() {
	cm.debugLog("Advancing context for replace: A=%d->%d, B=%d->%d",
		cm.aIndex, cm.aIndex+1, cm.bIndex, cm.bIndex+1)
	cm.aIndex++
	cm.bIndex++
}

func (cm *ContextManager) AdvanceForContainer() {
	cm.debugLog("Advancing context for container: A=%d->%d, B=%d->%d",
		cm.aIndex, cm.aIndex+1, cm.bIndex, cm.bIndex+1)
	cm.aIndex++
	cm.bIndex++
}

// generateListDiffEvents converts LCS analysis into a sequence of diff events for lists
func generateListDiffEvents(a, b jsonList, opts *options) []DiffEvent {
	// Handle empty lists
	if len(a) == 0 && len(b) == 0 {
		return []DiffEvent{}
	}
	if len(a) == 0 {
		// All additions
		events := make([]DiffEvent, len(b))
		for i, elem := range b {
			events[i] = AddEvent{BIndex: i, Element: elem}
		}
		return events
	}
	if len(b) == 0 {
		// All removals
		events := make([]DiffEvent, len(a))
		for i, elem := range a {
			events[i] = RemoveEvent{AIndex: i, Element: elem}
		}
		return events
	}

	// Use LCS to find matches with options awareness
	lcsResult := NewLcsWithOptions(a, b, opts)
	matches := lcsResult.IndexPairs()

	// Build events by walking through both arrays
	var events []DiffEvent
	var aIndex, bIndex, matchIndex int

	for aIndex < len(a) || bIndex < len(b) {
		// Check if we're at a match point
		atMatch := matchIndex < len(matches) &&
			aIndex == matches[matchIndex].Left &&
			bIndex == matches[matchIndex].Right

		if atMatch {
			// We have a match - but first check if it's a container that needs diffing
			// Refine options for this specific list index
			refinedOpts := refine(opts, PathIndex(aIndex))
			if sameContainerType(a[aIndex], b[bIndex], opts) &&
				!a[aIndex].equals(b[bIndex], refinedOpts) {
				// Compatible containers with differences
				events = append(events, ContainerDiffEvent{
					AIndex:   aIndex,
					BIndex:   bIndex,
					AElement: a[aIndex],
					BElement: b[bIndex],
				})
			} else {
				// Perfect match (check equality with refined options for precision)
				events = append(events, MatchEvent{
					AIndex:  aIndex,
					BIndex:  bIndex,
					Element: a[aIndex],
				})
			}
			aIndex++
			bIndex++
			matchIndex++
		} else {
			// No match - determine what kind of gap this is
			nextMatchA := len(a)
			nextMatchB := len(b)
			if matchIndex < len(matches) {
				nextMatchA = matches[matchIndex].Left
				nextMatchB = matches[matchIndex].Right
			}

			// Check if we should try container diffing even without LCS match
			if aIndex < len(a) && bIndex < len(b) &&
				aIndex < nextMatchA && bIndex < nextMatchB &&
				sameContainerType(a[aIndex], b[bIndex], opts) {
				// Compatible containers - try container diff
				events = append(events, ContainerDiffEvent{
					AIndex:   aIndex,
					BIndex:   bIndex,
					AElement: a[aIndex],
					BElement: b[bIndex],
				})
				aIndex++
				bIndex++
			} else if aIndex < len(a) && bIndex < len(b) &&
				aIndex < nextMatchA && bIndex < nextMatchB {
				// Different elements at same position - replacement
				events = append(events, ReplaceEvent{
					AIndex:   aIndex,
					BIndex:   bIndex,
					AElement: a[aIndex],
					BElement: b[bIndex],
				})
				aIndex++
				bIndex++
			} else if aIndex < nextMatchA && aIndex < len(a) {
				// Element only in A - removal
				events = append(events, RemoveEvent{
					AIndex:  aIndex,
					Element: a[aIndex],
				})
				aIndex++
			} else if bIndex < nextMatchB && bIndex < len(b) {
				// Element only in B - addition
				events = append(events, AddEvent{
					BIndex:  bIndex,
					Element: b[bIndex],
				})
				bIndex++
			} else {
				// Should not happen, but avoid infinite loop
				break
			}
		}
	}

	return events
}

// ============================================================================
// OBJECT-SPECIFIC EVENTS AND PROCESSOR
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
	Key        string
	OldValue   JsonNode
	NewValue   JsonNode
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

// ============================================================================
// SET-SPECIFIC EVENTS AND PROCESSOR
// ============================================================================

// SetElementEvent represents operations on set elements
type SetElementEvent struct {
	Operation string    // "ADD" or "REMOVE"
	Element   JsonNode
	Hash      [8]byte   // For identity tracking
}

func (e SetElementEvent) String() string {
	return fmt.Sprintf("SET_%s(%s)", e.Operation, e.Element.Json())
}

func (e SetElementEvent) GetType() string { return "SET_ELEMENT" }

// SetObjectDiffEvent represents an object in a set that needs recursive diffing
type SetObjectDiffEvent struct {
	OldObject JsonNode
	NewObject JsonNode
	Hash      [8]byte // Identity hash
}

func (e SetObjectDiffEvent) String() string {
	return fmt.Sprintf("SET_OBJECT_DIFF(%s -> %s)", e.OldObject.Json(), e.NewObject.Json())
}

func (e SetObjectDiffEvent) GetType() string { return "SET_OBJECT_DIFF" }

// SetDiffProcessor processes set diff events
type SetDiffProcessor struct {
	*BaseDiffProcessor
}

func NewSetDiffProcessor(path Path, opts *options, strategy patchStrategy) *SetDiffProcessor {
	return &SetDiffProcessor{
		BaseDiffProcessor: NewBaseDiffProcessor(path, opts, strategy),
	}
}

func (p *SetDiffProcessor) ProcessEvents(events []DiffEvent) Diff {
	p.debugLog("Starting to process %d set events", len(events))

	// Collect all add/remove events for the set operation
	var setElement *DiffElement
	
	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())
		
		switch e := event.(type) {
		case SetElementEvent:
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
			
		case SetObjectDiffEvent:
			p.processSetObjectDiffEvent(e)
			
		case SimpleReplaceEvent:
			p.processSimpleReplaceEvent(e)
			
		default:
			p.debugLog("WARNING: Unknown event type for SetDiffProcessor: %T", event)
		}
	}

	// Add the accumulated set diff element if it has changes
	if setElement != nil && (len(setElement.Remove) > 0 || len(setElement.Add) > 0) {
		p.finalDiff = append(p.finalDiff, *setElement)
	}

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *SetDiffProcessor) processSetObjectDiffEvent(event SetObjectDiffEvent) {
	p.debugLog("Processing set object diff: %s -> %s", event.OldObject.Json(), event.NewObject.Json())

	// For set object diffs, we need to create a path with PathSetKeys
	o1, _ := event.OldObject.(jsonObject)
	setKeysPath := newPathSetKeys(o1, p.opts)
	subPath := append(p.path.clone(), setKeysPath)
	
	subDiff := event.OldObject.diff(event.NewObject, subPath, p.opts, p.strategy)
	p.finalDiff = append(p.finalDiff, subDiff...)
}

func (p *SetDiffProcessor) processSimpleReplaceEvent(event SimpleReplaceEvent) {
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

// generateSetDiffEvents analyzes two sets and generates appropriate diff events
func generateSetDiffEvents(s1, s2 jsonSet, opts *options) []DiffEvent {
	var events []DiffEvent

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
			events = append(events, SetElementEvent{
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
				events = append(events, SetObjectDiffEvent{
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
			events = append(events, SetElementEvent{
				Operation: "ADD",
				Element:   s2Map[hc],
				Hash:      hc,
			})
		}
	}

	return events
}

// ============================================================================
// MULTISET-SPECIFIC EVENTS AND PROCESSOR
// ============================================================================

// MultisetElementEvent represents operations on multiset elements with counts
type MultisetElementEvent struct {
	Operation string    // "ADD" or "REMOVE"
	Element   JsonNode
	Count     int       // How many instances to add/remove
	Hash      [8]byte   // For identity tracking
}

func (e MultisetElementEvent) String() string {
	return fmt.Sprintf("MULTISET_%s(%s x%d)", e.Operation, e.Element.Json(), e.Count)
}

func (e MultisetElementEvent) GetType() string { return "MULTISET_ELEMENT" }

// MultisetDiffProcessor processes multiset diff events
type MultisetDiffProcessor struct {
	*BaseDiffProcessor
}

func NewMultisetDiffProcessor(path Path, opts *options, strategy patchStrategy) *MultisetDiffProcessor {
	return &MultisetDiffProcessor{
		BaseDiffProcessor: NewBaseDiffProcessor(path, opts, strategy),
	}
}

func (p *MultisetDiffProcessor) ProcessEvents(events []DiffEvent) Diff {
	p.debugLog("Starting to process %d multiset events", len(events))

	// Collect all add/remove events for the multiset operation
	var multisetElement *DiffElement
	
	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())
		
		switch e := event.(type) {
		case MultisetElementEvent:
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
			
		case SimpleReplaceEvent:
			p.processSimpleReplaceEvent(e)
			
		default:
			p.debugLog("WARNING: Unknown event type for MultisetDiffProcessor: %T", event)
		}
	}

	// Add the accumulated multiset diff element if it has changes
	if multisetElement != nil && (len(multisetElement.Remove) > 0 || len(multisetElement.Add) > 0) {
		p.finalDiff = append(p.finalDiff, *multisetElement)
	}

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *MultisetDiffProcessor) processSimpleReplaceEvent(event SimpleReplaceEvent) {
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

// generateMultisetDiffEvents analyzes two multisets and generates appropriate diff events
func generateMultisetDiffEvents(a1, a2 jsonMultiset, opts *options) []DiffEvent {
	var events []DiffEvent

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
			events = append(events, MultisetElementEvent{
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
			events = append(events, MultisetElementEvent{
				Operation: "ADD",
				Element:   a2Map[hc],
				Count:     added,
				Hash:      hc,
			})
		}
	}

	return events
}