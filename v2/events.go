package jd

import (
	"fmt"
)

// ============================================================================
// CORE EVENT-DRIVEN DIFF SYSTEM
// ============================================================================

// diffEvent represents an operation needed to transform one structure into another
type diffEvent interface {
	String() string // For debugging
	GetType() string
}

// matchEvent represents elements that are identical between A and B
type matchEvent struct {
	AIndex, BIndex int
	Element        JsonNode
}

func (e matchEvent) String() string {
	return fmt.Sprintf("MATCH(A[%d]=B[%d]: %s)", e.AIndex, e.BIndex, e.Element.Json())
}

func (e matchEvent) GetType() string { return "MATCH" }

// containerDiffEvent represents containers that are compatible and need recursive diffing
type containerDiffEvent struct {
	AIndex, BIndex     int
	AElement, BElement JsonNode
}

func (e containerDiffEvent) String() string {
	return fmt.Sprintf("CONTAINER_DIFF(A[%d] vs B[%d])", e.AIndex, e.BIndex)
}

func (e containerDiffEvent) GetType() string { return "CONTAINER_DIFF" }

// removeEvent represents an element that exists only in A (needs to be removed)
type removeEvent struct {
	AIndex  int
	Element JsonNode
}

func (e removeEvent) String() string {
	return fmt.Sprintf("REMOVE(A[%d]: %s)", e.AIndex, e.Element.Json())
}

func (e removeEvent) GetType() string { return "REMOVE" }

// addEvent represents an element that exists only in B (needs to be added)
type addEvent struct {
	BIndex  int
	Element JsonNode
}

func (e addEvent) String() string {
	return fmt.Sprintf("ADD(B[%d]: %s)", e.BIndex, e.Element.Json())
}

func (e addEvent) GetType() string { return "ADD" }

// replaceEvent represents elements at the same position that are different
type replaceEvent struct {
	AIndex, BIndex     int
	AElement, BElement JsonNode
}

func (e replaceEvent) String() string {
	return fmt.Sprintf("REPLACE(A[%d]: %s -> B[%d]: %s)",
		e.AIndex, e.AElement.Json(), e.BIndex, e.BElement.Json())
}

func (e replaceEvent) GetType() string { return "REPLACE" }

// simpleReplaceEvent represents a simple replacement between two different values
type simpleReplaceEvent struct {
	OldValue, NewValue JsonNode
}

func (e simpleReplaceEvent) String() string {
	return fmt.Sprintf("SIMPLE_REPLACE(%s -> %s)", e.OldValue.Json(), e.NewValue.Json())
}

func (e simpleReplaceEvent) GetType() string { return "SIMPLE_REPLACE" }

// ============================================================================
// CORE DIFF PROCESSOR
// ============================================================================

// diffProcessorState represents the current state of diff processing
type diffProcessorState int

const (
	idle diffProcessorState = iota
	accumulatingChanges
	processingMatch
	processingContainerDiff
)

func (s diffProcessorState) String() string {
	switch s {
	case idle:
		return "IDLE"
	case accumulatingChanges:
		return "ACCUMULATING_CHANGES"
	case processingMatch:
		return "PROCESSING_MATCH"
	case processingContainerDiff:
		return "PROCESSING_CONTAINER_DIFF"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

// baseDiffProcessor provides common functionality for all diff processors
type baseDiffProcessor struct {
	state       diffProcessorState
	currentDiff DiffElement
	finalDiff   Diff
	opts        *options
	strategy    patchStrategy
	path        Path
	debug       bool
}

func newBaseDiffProcessor(path Path, opts *options, strategy patchStrategy) *baseDiffProcessor {
	return &baseDiffProcessor{
		state:     idle,
		finalDiff: Diff{}, // Initialize to empty slice, not nil
		opts:      opts,
		strategy:  strategy,
		path:      path,
		debug:     false,
	}
}

func (p *baseDiffProcessor) SetDebug(debug bool) {
	p.debug = debug
}

func (p *baseDiffProcessor) debugLog(format string, args ...interface{}) {
	if p.debug {
		fmt.Printf("[BaseDiffProcessor:%s] "+format+"\n", append([]interface{}{p.state}, args...)...)
	}
}

// ============================================================================
// SIMPLE DIFF PROCESSOR (for primitive types)
// ============================================================================

// simpleDiffProcessor handles diff processing for simple types (bool, string, number, null, void)
type simpleDiffProcessor struct {
	*baseDiffProcessor
}

func newSimpleDiffProcessor(path Path, opts *options, strategy patchStrategy) *simpleDiffProcessor {
	return &simpleDiffProcessor{
		baseDiffProcessor: newBaseDiffProcessor(path, opts, strategy),
	}
}

func (p *simpleDiffProcessor) ProcessEvents(events []diffEvent) Diff {
	p.debugLog("Starting to process %d events", len(events))

	for i, event := range events {
		p.debugLog("Processing event %d: %s", i, event.String())
		p.processEvent(event)
	}

	p.debugLog("Final diff has %d elements", len(p.finalDiff))
	return p.finalDiff
}

func (p *simpleDiffProcessor) processEvent(event diffEvent) {
	switch e := event.(type) {
	case simpleReplaceEvent:
		p.processSimpleReplaceEvent(e)
	default:
		p.debugLog("WARNING: Unknown event type for simpleDiffProcessor: %T", event)
	}
}

func (p *simpleDiffProcessor) processSimpleReplaceEvent(event simpleReplaceEvent) {
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
func generateSimpleEvents(a, b JsonNode, opts *options) []diffEvent {
	if !opts.diffingOn {
		return []diffEvent{} // No events when diffing is off
	}
	if a.equals(b, opts) {
		return []diffEvent{}
	}
	return []diffEvent{simpleReplaceEvent{OldValue: a, NewValue: b}}
}
