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

func BenchmarkDiff_Array_PartialChanges(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	changePercentages := []float64{0.1, 0.25, 0.5}
	
	for _, size := range sizes {
		for _, changePercent := range changePercentages {
			b.Run(fmt.Sprintf("Size_%d_Changes_%.0f%%", size, changePercent*100), func(b *testing.B) {
				original := generateSequentialArray(size)
				
				modifiedData := make([]interface{}, size)
				for i := 0; i < size; i++ {
					modifiedData[i] = i
				}
				
				numChanges := int(float64(size) * changePercent)
				for i := 0; i < numChanges; i++ {
					pos := i * (size / numChanges)
					if pos < size {
						modifiedData[pos] = -modifiedData[pos].(int)
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

func BenchmarkDiff_Array_BlockMoves(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	blockSizes := []int{5, 10, 50}
	
	for _, size := range sizes {
		for _, blockSize := range blockSizes {
			if blockSize >= size/2 {
				continue
			}
			b.Run(fmt.Sprintf("Size_%d_MoveBlock_%d", size, blockSize), func(b *testing.B) {
				original := generateSequentialArray(size)
				
				modifiedData := make([]interface{}, size)
				for i := 0; i < size; i++ {
					modifiedData[i] = i
				}
				
				srcStart := size / 4
				destStart := 3 * size / 4
				for i := 0; i < blockSize && srcStart+i < size && destStart+i < size; i++ {
					modifiedData[srcStart+i], modifiedData[destStart+i] = modifiedData[destStart+i], modifiedData[srcStart+i]
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

func BenchmarkDiff_Array_MixedOperations(b *testing.B) {
	sizes := []int{100, 1000, 3000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d_Insert_Modify_Delete", size), func(b *testing.B) {
			original := generateSequentialArray(size)
			
			modifiedData := make([]interface{}, 0, size+10)
			numInserts := size / 20
			numModifies := size / 10
			numDeletes := size / 30
			
			i := 0
			for pos := 0; pos < size; pos++ {
				if pos%20 == 0 && numInserts > 0 {
					modifiedData = append(modifiedData, -1000-pos)
					numInserts--
				}
				
				if pos%30 == 0 && numDeletes > 0 {
					numDeletes--
					i++
					continue
				}
				
				value := i
				if pos%10 == 0 && numModifies > 0 {
					value = value + 10000
					numModifies--
				}
				
				modifiedData = append(modifiedData, value)
				i++
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

func BenchmarkDiff_Array_CompleteReverse_Pathological(b *testing.B) {
	sizes := []int{100, 500}
	
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