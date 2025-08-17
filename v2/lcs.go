package jd

// LCS provides functions to calculate Longest Common Subsequence (LCS)
// values from two arbitrary arrays.

import (
	"context"
	"fmt"
	"strings"
)

// lcs is the interface to calculate the LCS of two arrays.
type lcs interface {
	// Values calculates the LCS value of the two arrays.
	Values() (values []JsonNode)
	// ValueContext is a context aware version of Values()
	ValuesContext(ctx context.Context) (values []JsonNode, err error)
	// IndexPairs calculates paris of indices which have the same value in LCS.
	IndexPairs() (pairs []indexPair)
	// IndexPairsContext is a context aware version of IndexPairs()
	IndexPairsContext(ctx context.Context) (pairs []indexPair, err error)
	// Length calculates the length of the LCS.
	Length() (length int)
	// LengthContext is a context aware version of Length()
	LengthContext(ctx context.Context) (length int, err error)
	// Left returns one of the two arrays to be compared.
	Left() (leftValues []JsonNode)
	// Right returns the other of the two arrays to be compared.
	Right() (righttValues []JsonNode)

	// ============================================================================
	// DEBUGGING METHODS
	// ============================================================================

	// Table returns the LCS dynamic programming table for debugging
	Table() (table [][]int)
	// TableContext is a context aware version of Table()
	TableContext(ctx context.Context) (table [][]int, err error)
	// DebugString returns a human-readable debug representation of the LCS analysis
	DebugString() string
	// DebugMatches returns detailed information about each match found
	DebugMatches() []debugMatch
	// DebugTable returns a formatted string representation of the LCS table
	DebugTable() string
}

// indexPair represents an pair of indeices in the Left and Right arrays found in the LCS value.
type indexPair struct {
	Left  int
	Right int
}

// debugMatch provides detailed information about each match found during LCS analysis
type debugMatch struct {
	LeftIndex  int
	RightIndex int
	Value      JsonNode
	Position   int // Position in the LCS sequence
}

type lcsImpl struct {
	left    []JsonNode
	right   []JsonNode
	options *options
	/* for caching */
	table      [][]int
	indexPairs []indexPair
	values     []JsonNode
}

// newLcs creates a new LCS calculator from two arrays.
func newLcs(left, right []JsonNode) lcs {
	return newLcsWithOptions(left, right, newOptions([]Option{}))
}

// newLcsWithOptions creates a new LCS calculator from two arrays with options.
func newLcsWithOptions(left, right []JsonNode, opts *options) lcs {
	return &lcsImpl{
		left:       left,
		right:      right,
		options:    opts,
		table:      nil,
		indexPairs: nil,
		values:     nil,
	}
}

// Table implements lcs.Table()
func (lcs *lcsImpl) Table() (table [][]int) {
	table, _ = lcs.TableContext(context.Background())
	return table
}

// Table implements lcs.TableContext()
func (lcs *lcsImpl) TableContext(ctx context.Context) (table [][]int, err error) {
	if lcs.table != nil {
		return lcs.table, nil
	}

	sizeX := len(lcs.left) + 1
	sizeY := len(lcs.right) + 1

	table = make([][]int, sizeX)
	for x := 0; x < sizeX; x++ {
		table[x] = make([]int, sizeY)
	}

	for y := 1; y < sizeY; y++ {
		select { // check in each y to save some time
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// nop
		}
		for x := 1; x < sizeX; x++ {
			increment := 0
			// Use options-aware equality check with path-specific refinement
			leftIndex := x - 1
			rightIndex := y - 1

			// Refine options for the specific left index to handle PathOptions
			refinedOpts := refine(lcs.options, PathIndex(leftIndex))
			allOpts := append(refinedOpts.apply, refinedOpts.retain...)

			if lcs.left[leftIndex].Equals(lcs.right[rightIndex], allOpts...) {
				increment = 1
			}
			table[x][y] = max(table[x-1][y-1]+increment, table[x-1][y], table[x][y-1])
		}
	}

	lcs.table = table
	return table, nil
}

// Table implements lcs.Length()
func (lcs *lcsImpl) Length() (length int) {
	length, _ = lcs.LengthContext(context.Background())
	return length
}

// Table implements lcs.LengthContext()
func (lcs *lcsImpl) LengthContext(ctx context.Context) (length int, err error) {
	table, err := lcs.TableContext(ctx)
	if err != nil {
		return 0, err
	}
	return table[len(lcs.left)][len(lcs.right)], nil
}

// Table implements lcs.IndexPairs()
func (lcs *lcsImpl) IndexPairs() (pairs []indexPair) {
	pairs, _ = lcs.IndexPairsContext(context.Background())
	return pairs
}

// Table implements lcs.IndexPairsContext()
func (lcs *lcsImpl) IndexPairsContext(ctx context.Context) (pairs []indexPair, err error) {
	if lcs.indexPairs != nil {
		return lcs.indexPairs, nil
	}

	table, err := lcs.TableContext(ctx)
	if err != nil {
		return nil, err
	}

	pairs = make([]indexPair, table[len(table)-1][len(table[0])-1])

	for x, y := len(lcs.left), len(lcs.right); x > 0 && y > 0; {
		if lcs.left[x-1].Equals(lcs.right[y-1]) {
			pairs[table[x][y]-1] = indexPair{Left: x - 1, Right: y - 1}
			x--
			y--
		} else {
			if table[x-1][y] >= table[x][y-1] {
				x--
			} else {
				y--
			}
		}
	}

	lcs.indexPairs = pairs

	return pairs, nil
}

// Table implements lcs.Values()
func (lcs *lcsImpl) Values() (values []JsonNode) {
	values, _ = lcs.ValuesContext(context.Background())
	return values
}

// Table implements lcs.ValuesContext()
func (lcs *lcsImpl) ValuesContext(ctx context.Context) (values []JsonNode, err error) {
	if lcs.values != nil {
		return lcs.values, nil
	}

	pairs, err := lcs.IndexPairsContext(ctx)
	if err != nil {
		return nil, err
	}

	values = make([]JsonNode, len(pairs))
	for i, pair := range pairs {
		values[i] = lcs.left[pair.Left]
	}
	lcs.values = values

	return values, nil
}

// Table implements lcs.Left()
func (lcs *lcsImpl) Left() (leftValues []JsonNode) {
	leftValues = lcs.left
	return
}

// Table implements lcs.Right()
func (lcs *lcsImpl) Right() (rightValues []JsonNode) {
	rightValues = lcs.right
	return
}

// ============================================================================
// DEBUGGING METHOD IMPLEMENTATIONS
// ============================================================================

// DebugString returns a human-readable debug representation of the LCS analysis
func (lcs *lcsImpl) DebugString() string {
	pairs := lcs.IndexPairs()
	values := lcs.Values()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("LCS Analysis:\n"))
	sb.WriteString(fmt.Sprintf("  Left array length: %d\n", len(lcs.left)))
	sb.WriteString(fmt.Sprintf("  Right array length: %d\n", len(lcs.right)))
	sb.WriteString(fmt.Sprintf("  LCS length: %d\n", len(values)))
	sb.WriteString(fmt.Sprintf("  Matches found: %d\n", len(pairs)))

	if len(pairs) > 0 {
		sb.WriteString("  Match details:\n")
		for i, pair := range pairs {
			value := values[i]
			sb.WriteString(fmt.Sprintf("    %d: L[%d]=R[%d] = %s\n",
				i, pair.Left, pair.Right, value.Json()))
		}
	}

	return sb.String()
}

// DebugMatches returns detailed information about each match found
func (lcs *lcsImpl) DebugMatches() []debugMatch {
	pairs := lcs.IndexPairs()
	values := lcs.Values()

	matches := make([]debugMatch, len(pairs))
	for i, pair := range pairs {
		matches[i] = debugMatch{
			LeftIndex:  pair.Left,
			RightIndex: pair.Right,
			Value:      values[i],
			Position:   i,
		}
	}

	return matches
}

// DebugTable returns a formatted string representation of the LCS table
func (lcs *lcsImpl) DebugTable() string {
	table := lcs.Table()
	if len(table) == 0 {
		return "Empty LCS table"
	}

	var sb strings.Builder

	// Header with right array values
	sb.WriteString("LCS Table (Left=rows, Right=columns):\n")
	sb.WriteString("    ")
	for j := 0; j < len(lcs.right); j++ {
		sb.WriteString(fmt.Sprintf("%3s ", lcs.right[j].Json()[:min(3, len(lcs.right[j].Json()))]))
	}
	sb.WriteString("\n")

	// Table rows
	for i, row := range table {
		if i == 0 {
			sb.WriteString("  ")
		} else {
			leftVal := lcs.left[i-1].Json()
			sb.WriteString(fmt.Sprintf("%3s", leftVal[:min(3, len(leftVal))]))
		}

		for _, val := range row {
			sb.WriteString(fmt.Sprintf("%4d", val))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(first int, rest ...int) (max int) {
	max = first
	for _, value := range rest {
		if value > max {
			max = value
		}
	}
	return
}
