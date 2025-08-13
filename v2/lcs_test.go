package jd

import (
	"context"
	"reflect"
	"testing"
	"time"
)

// Helper function to convert interface{} slice to JsonNode slice
func toJsonNodes(values []interface{}) []JsonNode {
	nodes := make([]JsonNode, len(values))
	for i, v := range values {
		node, _ := NewJsonNode(v)
		nodes[i] = node
	}
	return nodes
}

func TestLCS(t *testing.T) {
	cases := []struct {
		left       []JsonNode
		right      []JsonNode
		indexPairs []indexPair
		values     []JsonNode
		length     int
	}{
		{
			left:       toJsonNodes([]interface{}{1, 2, 3}),
			right:      toJsonNodes([]interface{}{2, 3}),
			indexPairs: []indexPair{{1, 0}, {2, 1}},
			values:     toJsonNodes([]interface{}{2, 3}),
			length:     2,
		},
		{
			left:       toJsonNodes([]interface{}{2, 3}),
			right:      toJsonNodes([]interface{}{1, 2, 3}),
			indexPairs: []indexPair{{0, 1}, {1, 2}},
			values:     toJsonNodes([]interface{}{2, 3}),
			length:     2,
		},
		{
			left:       toJsonNodes([]interface{}{2, 3}),
			right:      toJsonNodes([]interface{}{2, 5, 3}),
			indexPairs: []indexPair{{0, 0}, {1, 2}},
			values:     toJsonNodes([]interface{}{2, 3}),
			length:     2,
		},
		{
			left:       toJsonNodes([]interface{}{2, 3, 3}),
			right:      toJsonNodes([]interface{}{2, 5, 3}),
			indexPairs: []indexPair{{0, 0}, {2, 2}},
			values:     toJsonNodes([]interface{}{2, 3}),
			length:     2,
		},
		{
			left:       toJsonNodes([]interface{}{1, 2, 5, 3, 1, 1, 5, 8, 3}),
			right:      toJsonNodes([]interface{}{1, 2, 3, 3, 4, 4, 5, 1, 6}),
			indexPairs: []indexPair{{0, 0}, {1, 1}, {2, 6}, {4, 7}},
			values:     toJsonNodes([]interface{}{1, 2, 5, 1}),
			length:     4,
		},
		{
			left:       toJsonNodes([]interface{}{}),
			right:      toJsonNodes([]interface{}{2, 5, 3}),
			indexPairs: []indexPair{},
			values:     toJsonNodes([]interface{}{}),
			length:     0,
		},
		{
			left:       toJsonNodes([]interface{}{3, 4}),
			right:      toJsonNodes([]interface{}{}),
			indexPairs: []indexPair{},
			values:     toJsonNodes([]interface{}{}),
			length:     0,
		},
		{
			left:       toJsonNodes([]interface{}{"foo"}),
			right:      toJsonNodes([]interface{}{"baz", "foo"}),
			indexPairs: []indexPair{{0, 1}},
			values:     toJsonNodes([]interface{}{"foo"}),
			length:     1,
		},
		{
			left:       toJsonNodes([]interface{}{int(byte('T')), int(byte('G')), int(byte('A')), int(byte('G')), int(byte('T')), int(byte('A'))}),
			right:      toJsonNodes([]interface{}{int(byte('G')), int(byte('A')), int(byte('T')), int(byte('A'))}),
			indexPairs: []indexPair{{1, 0}, {2, 1}, {4, 2}, {5, 3}},
			values:     toJsonNodes([]interface{}{int(byte('G')), int(byte('A')), int(byte('T')), int(byte('A'))}),
			length:     4,
		},
	}

	for i, c := range cases {
		lcs := newLcs(c.left, c.right)

		actualPairs := lcs.IndexPairs()
		if !reflect.DeepEqual(actualPairs, c.indexPairs) {
			t.Errorf("test case %d failed at index pair, actual: %#v, expected: %#v", i, actualPairs, c.indexPairs)
		}

		actualValues := lcs.Values()
		if !reflect.DeepEqual(actualValues, c.values) {
			t.Errorf("test case %d failed at values, actual: %#v, expected: %#v", i, actualValues, c.values)
		}

		actualLength := lcs.Length()
		if actualLength != c.length {
			t.Errorf("test case %d failed at length, actual: %d, expected: %d", i, actualLength, c.length)
		}
	}
}

func TestContextCancel(t *testing.T) {
	leftRaw := make([]interface{}, 100000) // takes over 1 sec
	rightRaw := make([]interface{}, 100000)
	rightRaw[0] = 1
	rightRaw[len(rightRaw)-1] = 1
	left := toJsonNodes(leftRaw)
	right := toJsonNodes(rightRaw)
	lcs := newLcs(left, right)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	_, err := lcs.LengthContext(ctx)
	if err != context.Canceled {
		t.Fatalf("unexpected err: %s", err)
	}
}
