package jd

import (
	"fmt"
)

// ============================================================================
// CORE EVENT-DRIVEN DIFF SYSTEM
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
// CORE DIFF PROCESSOR
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
