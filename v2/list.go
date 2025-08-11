package jd

import (
	"fmt"
)

type jsonList []JsonNode

var _ JsonNode = jsonList(nil)

func (l jsonList) Json(_ ...Option) string {
	return renderJson(l.raw())
}

func (l jsonList) Yaml(_ ...Option) string {
	return renderYaml(l.raw())
}

func (l jsonList) raw() interface{} {
	return jsonArray(l).raw()
}

func (l1 jsonList) Equals(n JsonNode, opts ...Option) bool {
	o := refine(&options{retain: opts}, nil)
	return l1.equals(n, o)
}

func (l1 jsonList) equals(n JsonNode, o *options) bool {
	n2 := dispatch(n, o)
	l2, ok := n2.(jsonList)
	if !ok {
		return false
	}
	if len(l1) != len(l2) {
		return false
	}
	for i, v1 := range l1 {
		v2 := l2[i]
		if !v1.equals(v2, o) {
			return false
		}
	}
	return true
}

func (l jsonList) hashCode(opts *options) [8]byte {
	b := []byte{0xF5, 0x18, 0x0A, 0x71, 0xA4, 0xC4, 0x03, 0xF3} // random bytes
	for _, n := range l {
		h := n.hashCode(opts)
		b = append(b, h[:]...)
	}
	return hash(b)
}

func (l jsonList) Diff(n JsonNode, opts ...Option) Diff {
	o := &options{retain: opts}
	return l.diff(n, make(Path, 0), o, getPatchStrategy(o))
}

func (a jsonList) diff(
	n JsonNode,
	path Path,
	opts *options,
	strategy patchStrategy,
) Diff {
	b, ok := n.(jsonList)
	if !ok {
		return a.diffDifferentTypes(n, path, strategy)
	}
	if strategy == mergePatchStrategy {
		return a.diffMergePatchStrategy(b, path, opts)
	}

	// NEW EVENT-DRIVEN ALGORITHM
	return a.diffEventDriven(b, path, opts, strategy)
}

// diffEventDriven implements the new event-driven diff algorithm
func (a jsonList) diffEventDriven(b jsonList, path Path, opts *options, strategy patchStrategy) Diff {
	// Step 1: Generate diff events using LCS analysis
	events := a.generateDiffEvents(b, opts)

	// Step 2: Create processor with path structure that matches original diffRest
	// The original algorithm used append(path, PathIndex(0)), so we need:
	// - basePath = the provided path (e.g., ["qux"])
	// - pathIndex = 0 (starting index within the list)
	basePath := path
	pathIndex := PathIndex(0)

	processor := NewListDiffProcessor(basePath, pathIndex, voidNode{}, a, b, opts, strategy)

	// Step 3: Process events to generate final diff
	return processor.ProcessEvents(events)
}

// ============================================================================
// EVENT-DRIVEN DIFF SYSTEM
// ============================================================================

// DiffEvent represents an operation needed to transform list A into list B
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

// generateDiffEvents converts LCS analysis into a sequence of diff events
func (a jsonList) generateDiffEvents(b jsonList, opts *options) []DiffEvent {
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
// LIST DIFF PROCESSOR - STATE MACHINE
// ============================================================================

// ListDiffProcessorState represents the current state of diff processing
type ListDiffProcessorState int

const (
	IDLE ListDiffProcessorState = iota
	ACCUMULATING_CHANGES
	PROCESSING_MATCH
	PROCESSING_CONTAINER_DIFF
)

func (s ListDiffProcessorState) String() string {
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
		state:                IDLE,
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
	if p.state == ACCUMULATING_CHANGES {
		// The matched element becomes the After context for the previous operation
		p.nextAfterElement = event.Element
		p.finalizeCurrentDiff()
	}

	p.state = PROCESSING_MATCH
	p.pathCalc.AdvanceForMatch()
	p.contextMgr.SetCursors(event.AIndex+1, event.BIndex+1) // Move past the match
	p.lastProcessedElement = event.Element                  // Track the matched element
	p.state = IDLE
}

func (p *ListDiffProcessor) processContainerDiffEvent(event ContainerDiffEvent) {
	p.debugLog("Processing container diff: A[%d] vs B[%d]", event.AIndex, event.BIndex)

	// If we were accumulating changes, finalize them first
	if p.state == ACCUMULATING_CHANGES {
		// The After context should be the element from original array A that comes after the previous operation
		p.nextAfterElement = event.AElement // Use A element, not B element
		p.finalizeCurrentDiff()
	}

	p.state = PROCESSING_CONTAINER_DIFF

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
	p.state = IDLE
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
	if p.state != ACCUMULATING_CHANGES {
		p.startAccumulating()
	}
}

func (p *ListDiffProcessor) startAccumulating() {
	p.debugLog("Starting to accumulate changes at path %s", p.pathCalc.CurrentPath())
	p.state = ACCUMULATING_CHANGES

	// Use the last processed element as Before context
	beforeContext := []JsonNode{p.lastProcessedElement}

	p.currentDiff = DiffElement{
		Path:   p.pathCalc.CurrentPath(),
		Before: beforeContext,
	}
}

func (p *ListDiffProcessor) finalizeCurrentDiff() {
	if p.state != ACCUMULATING_CHANGES {
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
	p.state = IDLE
}

func (p *ListDiffProcessor) finalize() {
	if p.state == ACCUMULATING_CHANGES {
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

func (a jsonList) diffRest(
	pathIndex PathIndex,
	b jsonList,
	path Path,
	aHashes, bHashes, commonSequence []interface{},
	previous JsonNode,
	opts *options,
	strategy patchStrategy,
) Diff {
	var aCursor, bCursor, commonSequenceCursor int
	pathCursor := pathIndex
	pathNow := func() Path {
		return append(path.clone().drop(), pathCursor)
	}
	endA := func() bool {
		return aCursor == len(a)
	}
	endB := func() bool {
		return bCursor == len(b)
	}
	atCommonA := func() bool {
		if endA() || len(commonSequence) == 0 {
			return false
		}
		return aHashes[aCursor] == commonSequence[0]
	}
	atCommonB := func() bool {
		if endB() || len(commonSequence) == 0 {
			return false
		}
		return bHashes[bCursor] == commonSequence[0]
	}
	d := Diff{{
		Path:   pathNow(),
		Before: []JsonNode{previous},
	}}
	haveDiff := func() bool {
		if len(d) == 0 {
			return false
		}
		if len(d[0].Add) > 0 || len(d[0].Remove) > 0 {
			return true
		}
		return false
	}
	after := func() []JsonNode {
		i := aCursor - commonSequenceCursor
		if i+1 > len(a) {
			return []JsonNode{voidNode{}}
		}
		return []JsonNode{a[i]}
	}

accumulatingDiff:
	for {
		switch {
		case endA():
			// We are at the end of A so there are no more
			// common elements. So we accumulate the rest
			// of B as additions. The path cursor advances
			// by 2 because the result is getting longer
			// by 1 and we are moving to the next element.
			for !endB() {
				d[0].Add = append(d[0].Add, b[bCursor])
				bCursor++
				pathCursor += 2
			}
			break accumulatingDiff
		case endB():
			// We are at the end of B so there are no more
			// common elements. So we accumulate the rest
			// of A as removals. The path cursor stays the
			// same because the result is getting shorter
			// by 1 but we are also moving to the next
			// element.
			for !endA() {
				d[0].Remove = append(d[0].Remove, a[aCursor])
				aCursor++
			}
			break accumulatingDiff
		case atCommonA() && atCommonB():
			// We are at a common element of A and B.
			// All cursors advance because we are moving
			// past a common element.
			aCursor++
			bCursor++
			commonSequenceCursor++
			pathCursor++
			break accumulatingDiff
		case atCommonA():
			// We are at a common element in A. We need to
			// catch up B. Add elements of B until we do.
			for !atCommonB() {
				d[0].Add = append(d[0].Add, b[bCursor])
				bCursor++
				pathCursor++
			}
		case atCommonB():
			// We are at a common element in B. We need to
			// catch up A. Remove elements of A until we
			// do.
			for !atCommonA() {
				d[0].Remove = append(d[0].Remove, a[aCursor])
				aCursor++
			}
		case sameContainerType(a[aCursor], b[bCursor], opts):
			// We are at compatible containers which
			// contain additional differences. If we've
			// accumulated differences at this level then
			// keep them before the sub-diff.
			subDiff := a[aCursor].diff(b[bCursor], pathNow(), opts, strategy)
			if haveDiff() {
				d[0].After = after()
				d = append(d, subDiff...)
			} else {
				d = subDiff
			}
			aCursor++
			bCursor++
			pathCursor++
			break accumulatingDiff
		default:
			// We are at elements of A and B which are
			// different. Add them to the accumulated diff
			// and continue.
			d[0].Remove = append(d[0].Remove, a[aCursor])
			d[0].Add = append(d[0].Add, b[bCursor])
			aCursor++
			bCursor++
			pathCursor++
		}
	}

	if !haveDiff() {
		// Throw away temporary diff because we didn't
		// accumulate anything.
		d = Diff{}
	} else {
		if len(d[0].Path) > len(path) {
			// This is a subdiff. Don't touch it.
		} else {
			// Record context of accumulated diff. If we appended
			// a sub-diff then it already has context.
			if len(d) < 2 {
				d[0].After = after()
			}
		}
	}
	if endA() && endB() {
		return d
	}
	// Cursors point to the next elements.
	return append(d, a[aCursor:].diffRest(
		pathCursor,
		b[bCursor:],
		pathNow(),
		aHashes[aCursor:], bHashes[bCursor:], commonSequence[commonSequenceCursor:],
		b[bCursor-1],
		opts,
		strategy,
	)...)
}

func (a jsonList) diffDifferentTypes(n JsonNode, path Path, strategy patchStrategy) Diff {
	var e DiffElement
	switch strategy {
	case mergePatchStrategy:
		e = DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: path.clone(),
			Add:  jsonArray{n},
		}
	default:
		e = DiffElement{
			Path:   path.clone(),
			Remove: nodeList(a),
			Add:    nodeList(n),
		}
	}
	return Diff{e}
}

func (a jsonList) diffMergePatchStrategy(b jsonList, path Path, opts *options) Diff {
	if !a.equals(b, opts) {
		e := DiffElement{
			Metadata: Metadata{
				Merge: true,
			},
			Path: path.clone(),
			Add:  nodeList(b),
		}
		return Diff{e}
	}
	return Diff{}
}

func sameContainerType(n1, n2 JsonNode, opts *options) bool {
	c1 := dispatch(n1, opts)
	c2 := dispatch(n2, opts)
	switch c1.(type) {
	case jsonObject:
		if _, ok := c2.(jsonObject); ok {
			return true
		}
	case jsonList:
		if _, ok := c2.(jsonList); ok {
			return true
		}
	case jsonSet:
		if _, ok := c2.(jsonSet); ok {
			return true
		}
	case jsonMultiset:
		if _, ok := c2.(jsonMultiset); ok {
			return true
		}
	default:
		return false
	}
	return false
}

func (l jsonList) Patch(d Diff) (JsonNode, error) {
	return patchAll(l, d)
}

func (l jsonList) patch(pathBehind, pathAhead Path, before, removeValues, addValues, after []JsonNode, strategy patchStrategy) (JsonNode, error) {

	if strategy == mergePatchStrategy {
		return patch(l, pathBehind, pathAhead, before, removeValues, addValues, after, mergePatchStrategy)
	}

	// Special case for replacing the whole list
	if len(pathAhead) == 0 {
		if len(removeValues) > 1 || len(addValues) > 1 {
			return nil, fmt.Errorf("cannot replace list with multiple values")
		}
		if len(removeValues) == 0 && strategy == strictPatchStrategy {
			return nil, fmt.Errorf("invalid diff. must declare list to replace it")
		}
		if !l.Equals(removeValues[0]) {
			return nil, fmt.Errorf("wanted %v. found %v", removeValues[0], l)
		}
		if len(addValues) == 0 {
			return voidNode{}, nil
		} else {
			return addValues[0], nil
		}
	}

	n, _, rest := pathAhead.next()
	i, ok := n.(PathIndex)
	if !ok {
		return nil, fmt.Errorf("invalid path element %T: expected float64", n)
	}

	// Recursive case
	if len(rest) > 0 {
		if int(i) > len(l)-1 {
			return nil, fmt.Errorf("patch index out of bounds: %v", i)
		}
		patchedNode, err := l[i].patch(append(pathBehind, n), rest, nil, removeValues, addValues, nil, strategy)
		if err != nil {
			return nil, err
		}
		l[i] = patchedNode
		return l, nil
	}

	// Special case for appending to the end of list
	if int(i) == -1 {
		if len(removeValues) > 0 {
			return nil, fmt.Errorf("invalid patch. appending to -1 index. but want to remove values")
		}
		l = append(l, addValues...)
		return l, nil
	}

	// Check context before
	for j, b := range before {
		bIndex := int(i) - (len(before) - j)
		switch {
		case bIndex < 0:
			if bIndex == -1 && isVoid(b) {
				continue
			}
			return nil, fmt.Errorf("invalid patch. before context %v out of bounds: %v", b, bIndex)
		case !b.Equals(l[bIndex]):
			return nil, fmt.Errorf("invalid patch. expected %v before. got %v", b, l[bIndex])
		}
	}

	// Patch list
	for len(removeValues) > 0 {
		if int(i) > len(l)-1 {
			return nil, fmt.Errorf("remove values out bounds: %v", i)
		}
		if !l[i].Equals(removeValues[0]) {
			return nil, fmt.Errorf("invalid patch. wanted %v. found %v", removeValues[0], l[i])
		}
		l = append(l[:i], l[i+1:]...)
		removeValues = removeValues[1:]
	}
	l2 := make(jsonList, i)
	copy(l2, l[:i])
	l2 = append(l2, addValues...)
	if int(i) < len(l) {
		l2 = append(l2, l[i:]...)
	}

	// Check context after
	for j, a := range after {
		aIndex := int(i) + j
		if aIndex > len(l)-1 {
			if aIndex == len(l) && isVoid(a) {
				continue
			}
			return nil, fmt.Errorf("invalid patch. after context %v out of bounds: %v", a, aIndex)
		}
		if !a.Equals(l[aIndex]) {
			return nil, fmt.Errorf("invalid patch. expected %v after. got %v", a, l[aIndex])
		}
	}

	return l2, nil
}
