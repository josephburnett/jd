package jd

import (
	"fmt"
	"testing"
)

func generateSequentialArray(size int) JsonNode {
	data := make([]interface{}, size)
	for i := 0; i < size; i++ {
		data[i] = i
	}
	node, _ := NewJsonNode(data)
	return node
}

func generateObject(numKeys int) JsonNode {
	data := make(map[string]interface{})
	for i := 0; i < numKeys; i++ {
		data[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	node, _ := NewJsonNode(data)
	return node
}

func generateDeeplyNestedObject(depth int) JsonNode {
	var current interface{} = "leaf_value"
	
	for i := 0; i < depth; i++ {
		current = map[string]interface{}{
			"nested": current,
		}
	}
	
	node, _ := NewJsonNode(current)
	return node
}

func createInsertAtBeginningDiff() Diff {
	original := []interface{}{0, 1, 2}
	modified := []interface{}{-1, 0, 1, 2}
	
	originalNode, _ := NewJsonNode(original)
	modifiedNode, _ := NewJsonNode(modified)
	
	return originalNode.Diff(modifiedNode)
}

func createSingleValueChangeDiff() Diff {
	original := map[string]interface{}{"key_0": "value_0", "key_1": "value_1"}
	modified := map[string]interface{}{"key_0": "new_value", "key_1": "value_1"}
	
	originalNode, _ := NewJsonNode(original)
	modifiedNode, _ := NewJsonNode(modified)
	
	return originalNode.Diff(modifiedNode)
}

func createAppendDiff() Diff {
	original := []interface{}{0, 1, 2}
	modified := []interface{}{0, 1, 2, 3}
	
	originalNode, _ := NewJsonNode(original)
	modifiedNode, _ := NewJsonNode(modified)
	
	return originalNode.Diff(modifiedNode)
}

func BenchmarkPatch_NoOp(b *testing.B) {
	original := generateObject(100)
	diff := Diff{}
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result, _ := original.Patch(diff)
		_ = result
	}
}

func BenchmarkDiff_IdenticalObjects(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateObject(size)
			identical := generateObject(size)
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				diff := original.Diff(identical)
				_ = diff
			}
		})
	}
}

func BenchmarkDiff_CompletelyDifferent(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateObject(size)
			different, _ := NewJsonNode(map[string]interface{}{"completely": "different"})
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				diff := original.Diff(different)
				_ = diff
			}
		})
	}
}