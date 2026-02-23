package jd

// LCS provides functions to calculate Longest Common Subsequence (LCS)
// values from two arbitrary arrays.

import (
	"context"
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
}

// indexPair represents an pair of indeices in the Left and Right arrays found in the LCS value.
type indexPair struct {
	Left  int
	Right int
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

func (lcs *lcsImpl) calcTable(ctx context.Context) (table [][]int, err error) {
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

// Length implements lcs.Length()
func (lcs *lcsImpl) Length() (length int) {
	length, _ = lcs.LengthContext(context.Background())
	return length
}

// LengthContext implements lcs.LengthContext()
func (lcs *lcsImpl) LengthContext(ctx context.Context) (length int, err error) {
	table, err := lcs.calcTable(ctx)
	if err != nil {
		return 0, err
	}
	return table[len(lcs.left)][len(lcs.right)], nil
}

// IndexPairs implements lcs.IndexPairs()
func (lcs *lcsImpl) IndexPairs() (pairs []indexPair) {
	pairs, _ = lcs.IndexPairsContext(context.Background())
	return pairs
}

// IndexPairsContext implements lcs.IndexPairsContext()
func (lcs *lcsImpl) IndexPairsContext(ctx context.Context) (pairs []indexPair, err error) {
	if lcs.indexPairs != nil {
		return lcs.indexPairs, nil
	}

	table, err := lcs.calcTable(ctx)
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

// Values implements lcs.Values()
func (lcs *lcsImpl) Values() (values []JsonNode) {
	values, _ = lcs.ValuesContext(context.Background())
	return values
}

// ValuesContext implements lcs.ValuesContext()
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
