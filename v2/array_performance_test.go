package jd

import (
	"fmt"
	"testing"
)

func BenchmarkPatch_Array_InsertAtBeginning(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateSequentialArray(size)
			diff := createInsertAtBeginningDiff()
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkPatch_Array_DeleteAtBeginning(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateSequentialArray(size)
			
			modifiedData := make([]interface{}, size-1)
			for i := 0; i < size-1; i++ {
				modifiedData[i] = i + 1
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

func BenchmarkPatch_Array_PrependVsAppend(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Prepend_Size_%d", size), func(b *testing.B) {
			original := generateSequentialArray(size)
			
			modifiedData := make([]interface{}, size+1)
			modifiedData[0] = -1
			for i := 0; i < size; i++ {
				modifiedData[i+1] = i
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
		
		b.Run(fmt.Sprintf("Append_Size_%d", size), func(b *testing.B) {
			original := generateSequentialArray(size)
			diff := createAppendDiff()
			
			b.ReportAllocs()
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				result, _ := original.Patch(diff)
				_ = result
			}
		})
	}
}

func BenchmarkPatch_Array_MultipleBeginningInsertions(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	insertions := []int{1, 5, 10}
	
	for _, size := range sizes {
		for _, numInserts := range insertions {
			b.Run(fmt.Sprintf("Size_%d_Inserts_%d", size, numInserts), func(b *testing.B) {
				original := generateSequentialArray(size)
				
				modifiedData := make([]interface{}, size+numInserts)
				for i := 0; i < numInserts; i++ {
					modifiedData[i] = -1 - i
				}
				for i := 0; i < size; i++ {
					modifiedData[i+numInserts] = i
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

func BenchmarkPatch_Array_MiddleInsertions(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateSequentialArray(size)
			
			modifiedData := make([]interface{}, size+1)
			mid := size / 2
			for i := 0; i < mid; i++ {
				modifiedData[i] = i
			}
			modifiedData[mid] = -999
			for i := mid; i < size; i++ {
				modifiedData[i+1] = i
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

func BenchmarkDiff_Array_LargeChanges(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			original := generateSequentialArray(size)
			
			modifiedData := make([]interface{}, size)
			for i := 0; i < size; i++ {
				modifiedData[i] = size - i - 1
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