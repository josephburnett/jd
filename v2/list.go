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
	o := newOptions(opts)
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
	events := generateListdiffEvents(a, b, opts)

	// Step 2: Create processor with path structure that matches original diffRest
	// The original algorithm used append(path, PathIndex(0)), so we need:
	// - basePath = the provided path (e.g., ["qux"])
	// - pathIndex = 0 (starting index within the list)
	basePath := path
	pathIndex := PathIndex(0)

	processor := newlistDiffProcessor(basePath, pathIndex, voidNode{}, a, b, opts, strategy)

	// Step 3: Process events to generate final diff
	return processor.ProcessEvents(events)
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
		if len(removeValues) > 0 && !l.Equals(removeValues[0]) {
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
		case bIndex >= len(l) || !b.Equals(l[bIndex]):
			if bIndex >= len(l) {
				return nil, fmt.Errorf("invalid patch. before context %v index %v out of bounds (list length %v)", b, bIndex, len(l))
			}
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
	if int(i) <= len(l) {
		copy(l2, l[:i])
	}
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

// ============================================================================
// LIST-SPECIFIC EVENT-DRIVEN DIFF SYSTEM
// ============================================================================

// listDiffProcessorState represents the current state of diff processing
type listDiffProcessorState int

const (
	listIdle listDiffProcessorState = iota
	listAccumulatingChanges
	listProcessingMatch
	listProcessingContainerDiff
)

func (s listDiffProcessorState) String() string {
	switch s {
	case listIdle:
		return "listIdle"
	case listAccumulatingChanges:
		return "listAccumulatingChanges"
	case listProcessingMatch:
		return "listProcessingMatch"
	case listProcessingContainerDiff:
		return "listProcessingContainerDiff"
	default:
		return fmt.Sprintf("LIST_UNKNOWN(%d)", s)
	}
}

// listDiffProcessor processes diff events using a state machine
type listDiffProcessor struct {
	state       listDiffProcessorState
	currentDiff DiffElement
	finalDiff   Diff
	opts        *options
	strategy    patchStrategy
	previous    JsonNode

	// Helper components
	pathCalc   *pathCalculator
	contextMgr *contextManager

	// Context tracking
	lastProcessedElement JsonNode
	nextAfterElement     JsonNode // The element that should appear as After context

	// Debug info
	debug bool
}

func newlistDiffProcessor(basePath Path, pathIndex PathIndex, previous JsonNode, a, b jsonList, opts *options, strategy patchStrategy) *listDiffProcessor {
	return &listDiffProcessor{
		state:                listIdle,
		finalDiff:            Diff{}, // Initialize to empty slice, not nil
		opts:                 opts,
		strategy:             strategy,
		previous:             previous,
		pathCalc:             newpathCalculator(basePath, pathIndex),
		contextMgr:           newcontextManager(a, b),
		lastProcessedElement: previous,
		nextAfterElement:     voidNode{}, // Default to void
		debug:                false,      // Disable debugging
	}
}

func (p *listDiffProcessor) SetDebug(debug bool) {
	p.debug = debug
	p.pathCalc.SetDebug(debug)
	p.contextMgr.SetDebug(debug)
}

func (p *listDiffProcessor) debugLog(format string, args ...interface{}) {
	if p.debug {
		fmt.Printf("[listDiffProcessor:%s] "+format+"\n", append([]interface{}{p.state}, args...)...)
	}
}

func (p *listDiffProcessor) ProcessEvents(events []diffEvent) Diff {
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

func (p *listDiffProcessor) processEvent(event diffEvent) {
	switch e := event.(type) {
	case matchEvent:
		p.processmatchEvent(e)
	case containerDiffEvent:
		p.processcontainerDiffEvent(e)
	case removeEvent:
		p.processremoveEvent(e)
	case addEvent:
		p.processaddEvent(e)
	case replaceEvent:
		p.processreplaceEvent(e)
	default:
		p.debugLog("WARNING: Unknown event type: %T", event)
	}
}

func (p *listDiffProcessor) processmatchEvent(event matchEvent) {
	p.debugLog("Processing match: A[%d]=B[%d]", event.AIndex, event.BIndex)

	// If we were accumulating changes, finalize them first
	if p.state == listAccumulatingChanges {
		// The matched element becomes the After context for the previous operation
		p.nextAfterElement = event.Element
		p.finalizeCurrentDiff()
	}

	p.state = listProcessingMatch
	p.pathCalc.AdvanceForMatch()
	p.contextMgr.SetCursors(event.AIndex+1, event.BIndex+1) // Move past the match
	p.lastProcessedElement = event.Element                  // Track the matched element
	p.state = listIdle
}

func (p *listDiffProcessor) processcontainerDiffEvent(event containerDiffEvent) {
	p.debugLog("Processing container diff: A[%d] vs B[%d]", event.AIndex, event.BIndex)

	// If we were accumulating changes, finalize them first
	if p.state == listAccumulatingChanges {
		// The After context should be the element from original array A that comes after the previous operation
		p.nextAfterElement = event.AElement // Use A element, not B element
		p.finalizeCurrentDiff()
	}

	p.state = listProcessingContainerDiff

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
	p.state = listIdle
}

func (p *listDiffProcessor) processremoveEvent(event removeEvent) {
	p.debugLog("Processing remove: A[%d]", event.AIndex)

	// Check if diffing is enabled for the current index
	currentIndex := p.pathCalc.GetCurrentIndex()
	refinedOpts := refine(p.opts, currentIndex)
	if !refinedOpts.diffingOn {
		p.debugLog("Skipping remove at index %v - diffing is off", currentIndex)
		p.pathCalc.AdvanceForRemove()
		return
	}

	p.ensureAccumulatingState()
	p.currentDiff.Remove = append(p.currentDiff.Remove, event.Element)
	p.pathCalc.AdvanceForRemove()
	// For removes, advance A cursor but not B cursor
}

func (p *listDiffProcessor) processaddEvent(event addEvent) {
	p.debugLog("Processing add: B[%d]", event.BIndex)

	// Check if diffing is enabled for the current index
	currentIndex := p.pathCalc.GetCurrentIndex()
	refinedOpts := refine(p.opts, currentIndex)
	if !refinedOpts.diffingOn {
		p.debugLog("Skipping add at index %v - diffing is off", currentIndex)
		p.pathCalc.AdvanceForAdd()
		return
	}

	p.ensureAccumulatingState()
	p.currentDiff.Add = append(p.currentDiff.Add, event.Element)
	p.pathCalc.AdvanceForAdd()
	// For adds, advance B cursor but not A cursor
}

func (p *listDiffProcessor) processreplaceEvent(event replaceEvent) {
	p.debugLog("Processing replace: A[%d] -> B[%d]", event.AIndex, event.BIndex)

	// Check if diffing is enabled for the current index
	currentIndex := p.pathCalc.GetCurrentIndex()
	refinedOpts := refine(p.opts, currentIndex)
	if !refinedOpts.diffingOn {
		p.debugLog("Skipping replace at index %v - diffing is off", currentIndex)
		p.pathCalc.AdvanceForReplace()
		return
	}

	p.ensureAccumulatingState()
	p.currentDiff.Remove = append(p.currentDiff.Remove, event.AElement)
	p.currentDiff.Add = append(p.currentDiff.Add, event.BElement)
	p.pathCalc.AdvanceForReplace()
	// For replaces, both cursors advance
}

func (p *listDiffProcessor) ensureAccumulatingState() {
	if p.state != listAccumulatingChanges {
		p.startAccumulating()
	}
}

func (p *listDiffProcessor) startAccumulating() {
	p.debugLog("Starting to accumulate changes at path %s", p.pathCalc.CurrentPath())
	p.state = listAccumulatingChanges

	// Use the last processed element as Before context
	beforeContext := []JsonNode{p.lastProcessedElement}

	p.currentDiff = DiffElement{
		Path:   p.pathCalc.CurrentPath(),
		Before: beforeContext,
	}
}

func (p *listDiffProcessor) finalizeCurrentDiff() {
	if p.state != listAccumulatingChanges {
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
	p.state = listIdle
}

func (p *listDiffProcessor) finalize() {
	if p.state == listAccumulatingChanges {
		// When finalizing at the end, After context should be void
		p.nextAfterElement = voidNode{}
		p.finalizeCurrentDiff()
	}
}

// pathCalculator handles path calculations and advancement during diff processing
type pathCalculator struct {
	basePath  Path
	pathIndex PathIndex
	debug     bool
}

func newpathCalculator(basePath Path, pathIndex PathIndex) *pathCalculator {
	return &pathCalculator{
		basePath:  basePath,
		pathIndex: pathIndex,
		debug:     false,
	}
}

func (pc *pathCalculator) SetDebug(debug bool) {
	pc.debug = debug
}

func (pc *pathCalculator) debugLog(format string, args ...interface{}) {
	if pc.debug {
		fmt.Printf("[pathCalculator] "+format+"\n", args...)
	}
}

func (pc *pathCalculator) CurrentPath() Path {
	return append(pc.basePath.clone(), pc.pathIndex)
}

func (pc *pathCalculator) GetCurrentIndex() PathIndex {
	return pc.pathIndex
}

func (pc *pathCalculator) AdvanceForMatch() {
	pc.debugLog("Advancing path for match: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex++
}

func (pc *pathCalculator) AdvanceForAdd() {
	pc.debugLog("Advancing path for add: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex += 1 // Path advances by 1 per added element
}

func (pc *pathCalculator) AdvanceForRemove() {
	pc.debugLog("Advancing path for remove: %d (no change)", pc.pathIndex)
	// Path stays same because list is getting shorter but we're moving to next element
}

func (pc *pathCalculator) AdvanceForReplace() {
	pc.debugLog("Advancing path for replace: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex++
}

func (pc *pathCalculator) AdvanceForContainer() {
	pc.debugLog("Advancing path for container: %d -> %d", pc.pathIndex, pc.pathIndex+1)
	pc.pathIndex++
}

// contextManager handles Before/After context calculation for diff elements
type contextManager struct {
	originalA jsonList
	originalB jsonList
	aIndex    int
	bIndex    int
	debug     bool
}

func newcontextManager(a, b jsonList) *contextManager {
	return &contextManager{
		originalA: a,
		originalB: b,
		aIndex:    0,
		bIndex:    0,
		debug:     false,
	}
}

func (cm *contextManager) SetDebug(debug bool) {
	cm.debug = debug
}

func (cm *contextManager) debugLog(format string, args ...interface{}) {
	if cm.debug {
		fmt.Printf("[contextManager] "+format+"\n", args...)
	}
}

func (cm *contextManager) SetCursors(aIndex, bIndex int) {
	cm.debugLog("Setting cursors: A=%d, B=%d", aIndex, bIndex)
	cm.aIndex = aIndex
	cm.bIndex = bIndex
}

func (cm *contextManager) GetBeforeContext(previous JsonNode) []JsonNode {
	cm.debugLog("Getting before context with previous: %s", previous.Json())
	return []JsonNode{previous}
}

func (cm *contextManager) GetAfterContext() []JsonNode {
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

func (cm *contextManager) AdvanceForMatch() {
	cm.debugLog("Advancing context for match: A=%d->%d, B=%d->%d",
		cm.aIndex, cm.aIndex+1, cm.bIndex, cm.bIndex+1)
	cm.aIndex++
	cm.bIndex++
}

func (cm *contextManager) AdvanceForAdd() {
	cm.debugLog("Advancing context for add: B=%d->%d", cm.bIndex, cm.bIndex+1)
	cm.bIndex++
}

func (cm *contextManager) AdvanceForRemove() {
	cm.debugLog("Advancing context for remove: A=%d->%d", cm.aIndex, cm.aIndex+1)
	cm.aIndex++
}

func (cm *contextManager) AdvanceForReplace() {
	cm.debugLog("Advancing context for replace: A=%d->%d, B=%d->%d",
		cm.aIndex, cm.aIndex+1, cm.bIndex, cm.bIndex+1)
	cm.aIndex++
	cm.bIndex++
}

func (cm *contextManager) AdvanceForContainer() {
	cm.debugLog("Advancing context for container: A=%d->%d, B=%d->%d",
		cm.aIndex, cm.aIndex+1, cm.bIndex, cm.bIndex+1)
	cm.aIndex++
	cm.bIndex++
}

// hasCommonElements quickly checks if two lists share any common elements
func hasCommonElements(a, b jsonList, opts *options) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	// For small arrays, just do direct comparison
	if len(a) <= 10 && len(b) <= 10 {
		for _, aItem := range a {
			for _, bItem := range b {
				if aItem.equals(bItem, opts) {
					return true
				}
			}
		}
		return false
	}

	// For larger arrays, use a map for faster lookup
	bSet := make(map[string]bool)
	for _, bItem := range b {
		key := bItem.Json() // Use JSON representation as key
		bSet[key] = true
	}

	for _, aItem := range a {
		key := aItem.Json()
		if bSet[key] {
			return true
		}
	}

	return false
}

// shouldUseFastPath determines if we should skip expensive LCS computation
func shouldUseFastPath(a, b jsonList, opts *options) bool {
	// If either array is empty, use fast path
	if len(a) == 0 || len(b) == 0 {
		return true
	}

	// If arrays are identical, use fast path
	if len(a) == len(b) {
		identical := true
		for i := 0; i < len(a); i++ {
			if !a[i].equals(b[i], opts) {
				identical = false
				break
			}
		}
		if identical {
			return true
		}
	}

	// For large arrays with no common elements, use fast path
	if len(a) > 100 && len(b) > 100 && !hasCommonElements(a, b, opts) {
		return true
	}

	return false
}

// generateFastPathEvents creates diff events for obvious cases without LCS
func generateFastPathEvents(a, b jsonList, opts *options) []diffEvent {
	// Empty lists
	if len(a) == 0 && len(b) == 0 {
		return []diffEvent{}
	}
	if len(a) == 0 {
		// All additions
		events := make([]diffEvent, len(b))
		for i, elem := range b {
			events[i] = addEvent{BIndex: i, Element: elem}
		}
		return events
	}
	if len(b) == 0 {
		// All removals
		events := make([]diffEvent, len(a))
		for i, elem := range a {
			events[i] = removeEvent{AIndex: i, Element: elem}
		}
		return events
	}

	// Identical arrays
	if len(a) == len(b) {
		identical := true
		for i := 0; i < len(a); i++ {
			if !a[i].equals(b[i], opts) {
				identical = false
				break
			}
		}
		if identical {
			// No differences - all matches
			events := make([]diffEvent, len(a))
			for i, elem := range a {
				events[i] = matchEvent{AIndex: i, BIndex: i, Element: elem}
			}
			return events
		}
	}

	// Completely different arrays - replace all
	events := make([]diffEvent, 0, len(a)+len(b))
	for i, elem := range a {
		events = append(events, removeEvent{AIndex: i, Element: elem})
	}
	for i, elem := range b {
		events = append(events, addEvent{BIndex: i, Element: elem})
	}
	return events
}

// myersEdit represents a single edit operation in Myers' algorithm
type myersEdit struct {
	Type  string   // "insert", "delete", "match"
	AIdx  int      // Index in array A (for delete/match)
	BIdx  int      // Index in array B (for insert/match)
	Value JsonNode // The value being operated on
}

// myersDiff implements Myers' O(ND) diff algorithm
func myersDiff(a, b jsonList, opts *options) []myersEdit {
	N := len(a)
	M := len(b)

	// Handle trivial cases
	if N == 0 && M == 0 {
		return []myersEdit{}
	}
	if N == 0 {
		edits := make([]myersEdit, M)
		for i := 0; i < M; i++ {
			edits[i] = myersEdit{Type: "insert", BIdx: i, Value: b[i]}
		}
		return edits
	}
	if M == 0 {
		edits := make([]myersEdit, N)
		for i := 0; i < N; i++ {
			edits[i] = myersEdit{Type: "delete", AIdx: i, Value: a[i]}
		}
		return edits
	}

	// Myers' algorithm implementation
	MAX := N + M
	V := make([]int, 2*MAX+1)

	var trace [][]int

	for D := 0; D <= MAX; D++ {
		// Copy current V for trace
		currentV := make([]int, len(V))
		copy(currentV, V)
		trace = append(trace, currentV)

		for k := -D; k <= D; k += 2 {
			var x int
			if k == -D || (k != D && V[k-1+MAX] < V[k+1+MAX]) {
				x = V[k+1+MAX]
			} else {
				x = V[k-1+MAX] + 1
			}

			y := x - k

			// Follow diagonal (matches)
			for x < N && y < M && a[x].equals(b[y], opts) {
				x++
				y++
			}

			V[k+MAX] = x

			if x >= N && y >= M {
				// Found the optimal path, backtrack to build edits
				return buildMyersEdits(a, b, trace, opts)
			}
		}
	}

	// Should never reach here for valid inputs
	return []myersEdit{}
}

// buildMyersEdits reconstructs the edit sequence from Myers' algorithm trace
func buildMyersEdits(a, b jsonList, trace [][]int, opts *options) []myersEdit {
	var edits []myersEdit
	N := len(a)
	M := len(b)
	MAX := N + M

	x, y := N, M

	// Backtrack through the trace
	for D := len(trace) - 1; D > 0; D-- {
		prevV := trace[D]

		k := x - y

		var prevK int
		if k == -D || (k != D && prevV[k-1+MAX] < prevV[k+1+MAX]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX := prevV[prevK+MAX]
		prevY := prevX - prevK

		// Add any diagonal moves (matches) first
		for x > prevX && y > prevY {
			x--
			y--
			edits = append([]myersEdit{{Type: "match", AIdx: x, BIdx: y, Value: a[x]}}, edits...)
		}

		// Add the edit operation
		if x > prevX {
			// Deletion
			x--
			edits = append([]myersEdit{{Type: "delete", AIdx: x, Value: a[x]}}, edits...)
		} else if y > prevY {
			// Insertion
			y--
			edits = append([]myersEdit{{Type: "insert", BIdx: y, Value: b[y]}}, edits...)
		}
	}

	// Add any remaining matches at the beginning
	for x > 0 && y > 0 && a[x-1].equals(b[y-1], opts) {
		x--
		y--
		edits = append([]myersEdit{{Type: "match", AIdx: x, BIdx: y, Value: a[x]}}, edits...)
	}

	return edits
}

// convertMyersToEvents converts Myers' edits to the internal diff event format
func convertMyersToEvents(myersEdits []myersEdit) []diffEvent {
	events := make([]diffEvent, 0, len(myersEdits))

	for _, edit := range myersEdits {
		switch edit.Type {
		case "match":
			events = append(events, matchEvent{
				AIndex:  edit.AIdx,
				BIndex:  edit.BIdx,
				Element: edit.Value,
			})
		case "delete":
			events = append(events, removeEvent{
				AIndex:  edit.AIdx,
				Element: edit.Value,
			})
		case "insert":
			events = append(events, addEvent{
				BIndex:  edit.BIdx,
				Element: edit.Value,
			})
		}
	}

	return events
}

// shouldUseMyers determines if we should use Myers algorithm instead of LCS
func shouldUseMyers(a, b jsonList) bool {
	// Use Myers for medium-sized arrays where LCS might be expensive
	// but fast path doesn't apply
	return len(a) > 10 && len(b) > 10 && (len(a) < 1000 || len(b) < 1000)
}

// generateListdiffEvents converts LCS analysis into a sequence of diff events for lists
func generateListdiffEvents(a, b jsonList, opts *options) []diffEvent {
	// Check if we should use fast path optimization
	if shouldUseFastPath(a, b, opts) {
		return generateFastPathEvents(a, b, opts)
	}

	// Check if we should use Myers algorithm
	if shouldUseMyers(a, b) {
		myersEdits := myersDiff(a, b, opts)
		return convertMyersToEvents(myersEdits)
	}

	// Use LCS to find matches with options awareness
	lcsResult := newLcsWithOptions(a, b, opts)
	matches := lcsResult.IndexPairs()

	// Build events by walking through both arrays
	var events []diffEvent
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
				events = append(events, containerDiffEvent{
					AIndex:   aIndex,
					BIndex:   bIndex,
					AElement: a[aIndex],
					BElement: b[bIndex],
				})
			} else {
				// Perfect match (check equality with refined options for precision)
				events = append(events, matchEvent{
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
				events = append(events, containerDiffEvent{
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
				events = append(events, replaceEvent{
					AIndex:   aIndex,
					BIndex:   bIndex,
					AElement: a[aIndex],
					BElement: b[bIndex],
				})
				aIndex++
				bIndex++
			} else if aIndex < nextMatchA && aIndex < len(a) {
				// Element only in A - removal
				events = append(events, removeEvent{
					AIndex:  aIndex,
					Element: a[aIndex],
				})
				aIndex++
			} else if bIndex < nextMatchB && bIndex < len(b) {
				// Element only in B - addition
				events = append(events, addEvent{
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
