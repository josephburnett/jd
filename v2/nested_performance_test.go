package jd

import (
	"fmt"
	"testing"
)

func createDeepLeafModification(depth int) Diff {
	original := generateDeeplyNestedObject(depth)

	var current interface{} = "modified_leaf_value"
	for i := 0; i < depth; i++ {
		current = map[string]interface{}{
			"nested": current,
		}
	}

	modified, _ := NewJsonNode(current)
	return original.Diff(modified)
}

func createDeepObjectAddition(depth int) Diff {
	original := generateDeeplyNestedObject(depth)

	var current interface{} = map[string]interface{}{
		"leaf_value": "original",
		"new_branch": map[string]interface{}{
			"added": "value",
		},
	}

	for i := 0; i < depth; i++ {
		current = map[string]interface{}{
			"nested": current,
		}
	}

	modified, _ := NewJsonNode(current)
	return original.Diff(modified)
}

func BenchmarkPatch_DeepNesting_PathCopying(b *testing.B) {
	depths := []int{5, 10, 15}

	for _, depth := range depths {
		b.Run(fmt.Sprintf("Depth_%d", depth), func(b *testing.B) {
			original := generateDeeplyNestedObject(depth)
			diff := createDeepLeafModification(depth)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkPatch_DeepNesting_StructuralChanges(b *testing.B) {
	depths := []int{5, 10, 15}

	for _, depth := range depths {
		b.Run(fmt.Sprintf("AddBranch_Depth_%d", depth), func(b *testing.B) {
			original := generateDeeplyNestedObject(depth)
			diff := createDeepObjectAddition(depth)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkPatch_DeepNesting_ArrayWithinObject(b *testing.B) {
	depths := []int{5, 10, 15}
	arraySizes := []int{10, 100, 1000}

	for _, depth := range depths {
		for _, arraySize := range arraySizes {
			b.Run(fmt.Sprintf("Depth_%d_ArraySize_%d", depth, arraySize), func(b *testing.B) {
				arrayData := make([]interface{}, arraySize)
				for i := 0; i < arraySize; i++ {
					arrayData[i] = i
				}

				var current interface{} = arrayData
				for i := 0; i < depth; i++ {
					current = map[string]interface{}{
						"nested": current,
					}
				}

				original, _ := NewJsonNode(current)

				modifiedArrayData := make([]interface{}, arraySize+1)
				modifiedArrayData[0] = -1
				for i := 0; i < arraySize; i++ {
					modifiedArrayData[i+1] = i
				}

				var modifiedCurrent interface{} = modifiedArrayData
				for i := 0; i < depth; i++ {
					modifiedCurrent = map[string]interface{}{
						"nested": modifiedCurrent,
					}
				}

				modified, _ := NewJsonNode(modifiedCurrent)
				diff := original.Diff(modified)

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					result, _ := original.Patch(diff)
					_ = result
				}
			})
		}
	}
}

func BenchmarkDiff_DeepNesting_ComplexStructures(b *testing.B) {
	depths := []int{5, 10, 15}

	for _, depth := range depths {
		b.Run(fmt.Sprintf("Depth_%d", depth), func(b *testing.B) {
			var current1 interface{} = map[string]interface{}{
				"value": "original",
				"array": []interface{}{1, 2, 3},
				"object": map[string]interface{}{
					"nested_value": "original_nested",
				},
			}

			var current2 interface{} = map[string]interface{}{
				"value": "modified",
				"array": []interface{}{1, 2, 3, 4},
				"object": map[string]interface{}{
					"nested_value": "modified_nested",
					"new_key":      "new_value",
				},
			}

			for i := 0; i < depth; i++ {
				current1 = map[string]interface{}{
					"nested": current1,
				}
				current2 = map[string]interface{}{
					"nested": current2,
				}
			}

			original, _ := NewJsonNode(current1)
			modified, _ := NewJsonNode(current2)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				diff := original.Diff(modified)
				_ = diff
			}
		})
	}
}

func BenchmarkPatch_DeepNesting_MultiplePathsModified(b *testing.B) {
	depths := []int{5, 10, 15}

	for _, depth := range depths {
		b.Run(fmt.Sprintf("Depth_%d", depth), func(b *testing.B) {
			var current interface{} = map[string]interface{}{
				"branch1": "value1",
				"branch2": "value2",
				"branch3": "value3",
			}

			for i := 0; i < depth; i++ {
				current = map[string]interface{}{
					"nested": current,
				}
			}

			original, _ := NewJsonNode(current)

			var modifiedCurrent interface{} = map[string]interface{}{
				"branch1": "modified_value1",
				"branch2": "modified_value2",
				"branch3": "modified_value3",
			}

			for i := 0; i < depth; i++ {
				modifiedCurrent = map[string]interface{}{
					"nested": modifiedCurrent,
				}
			}

			modified, _ := NewJsonNode(modifiedCurrent)
			diff := original.Diff(modified)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}
