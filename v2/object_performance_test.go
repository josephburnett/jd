package jd

import (
	"fmt"
	"testing"
)

func BenchmarkPatch_Object_KeyRestructuring(b *testing.B) {
	scenarios := []struct {
		name string
		keys int
		operation string
	}{
		{"AddKey_Small", 10, "add"},
		{"AddKey_Large", 1000, "add"},
		{"RemoveKey_Small", 10, "remove"},
		{"RemoveKey_Large", 1000, "remove"},
	}
	
	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			original := generateObject(scenario.keys)
			
			var diff Diff
			if scenario.operation == "add" {
				modifiedData := make(map[string]interface{})
				for i := 0; i < scenario.keys; i++ {
					modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
				}
				modifiedData["new_key"] = "new_value"
				modified, _ := NewJsonNode(modifiedData)
				diff = original.Diff(modified)
			} else {
				modifiedData := make(map[string]interface{})
				for i := 1; i < scenario.keys; i++ {
					modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
				}
				modified, _ := NewJsonNode(modifiedData)
				diff = original.Diff(modified)
			}
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkPatch_Object_SingleValueChange(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Keys_%d", size), func(b *testing.B) {
			original := generateObject(size)
			diff := createSingleValueChangeDiff()
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkPatch_Object_BulkOperations(b *testing.B) {
	sizes := []int{100, 1000}
	changePercentages := []float64{0.1, 0.5, 1.0}
	
	for _, size := range sizes {
		for _, changePercent := range changePercentages {
			b.Run(fmt.Sprintf("Size_%d_Change_%.0f%%", size, changePercent*100), func(b *testing.B) {
				original := generateObject(size)
				
				modifiedData := make(map[string]interface{})
				for i := 0; i < size; i++ {
					modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
				}
				
				numChanges := int(float64(size) * changePercent)
				for i := 0; i < numChanges; i++ {
					modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("modified_value_%d", i)
				}
				
				modified, _ := NewJsonNode(modifiedData)
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

func BenchmarkPatch_Object_EmptyToFull(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original, _ := NewJsonNode(map[string]interface{}{})
			target := generateObject(size)
			diff := original.Diff(target)
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkPatch_Object_FullToEmpty(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateObject(size)
			target, _ := NewJsonNode(map[string]interface{}{})
			diff := original.Diff(target)
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkDiff_Object_PartialUpdates(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	changePercentages := []float64{0.1, 0.25, 0.5}
	
	for _, size := range sizes {
		for _, changePercent := range changePercentages {
			b.Run(fmt.Sprintf("Size_%d_Changes_%.0f%%", size, changePercent*100), func(b *testing.B) {
				original := generateObject(size)
				
				modifiedData := make(map[string]interface{})
				for i := 0; i < size; i++ {
					modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
				}
				
				numChanges := int(float64(size) * changePercent)
				for i := 0; i < numChanges; i++ {
					keyIndex := i * (size / numChanges)
					if keyIndex < size {
						modifiedData[fmt.Sprintf("key_%d", keyIndex)] = fmt.Sprintf("updated_value_%d", keyIndex)
					}
				}
				modified, _ := NewJsonNode(modifiedData)
				
				b.ReportAllocs()
				b.ResetTimer()
				
				for i := 0; i < b.N; i++ {
					diff := original.Diff(modified)
					_ = diff
				}
			})
		}
	}
}

func BenchmarkDiff_Object_ApiUpdatePatterns(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d_ApiUpdate", size), func(b *testing.B) {
			original := generateObject(size)
			
			modifiedData := make(map[string]interface{})
			for i := 0; i < size; i++ {
				modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
			}
			
			modifiedData["updated_at"] = "2025-08-13T10:00:00Z"
			modifiedData["version"] = 2
			
			for i := 0; i < 5 && i < size; i++ {
				modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("new_value_%d", i)
			}
			
			numDeletes := min(3, size/10)
			for i := 0; i < numDeletes; i++ {
				delete(modifiedData, fmt.Sprintf("key_%d", size-1-i))
			}
			
			modified, _ := NewJsonNode(modifiedData)
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				diff := original.Diff(modified)
				_ = diff
			}
		})
	}
}

func BenchmarkDiff_Object_CompleteRewrite_Pathological(b *testing.B) {
	sizes := []int{100, 500}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateObject(size)
			
			modifiedData := make(map[string]interface{})
			for i := 0; i < size; i++ {
				modifiedData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("modified_value_%d", i)
			}
			modified, _ := NewJsonNode(modifiedData)
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				diff := original.Diff(modified)
				_ = diff
			}
		})
	}
}

func BenchmarkPatch_Object_KeyTypeChanges(b *testing.B) {
	sizes := []int{100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateObject(size)
			
			modifiedData := make(map[string]interface{})
			for i := 0; i < size; i++ {
				if i%2 == 0 {
					modifiedData[fmt.Sprintf("key_%d", i)] = i
				} else {
					modifiedData[fmt.Sprintf("key_%d", i)] = map[string]interface{}{"nested": i}
				}
			}
			
			modified, _ := NewJsonNode(modifiedData)
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