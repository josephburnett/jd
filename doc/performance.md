# Performance Testing Implementation Guide

## Overview

This document provides a systematic implementation plan for performance testing v2 JsonNode operations **before implementing immutability**. The primary goal is establishing a performance baseline and focusing on operations that will be most challenging to optimize with structural sharing.

## Implementation Objectives

- **Establish Baseline**: Measure current mutable operation performance 
- **Focus on Challenging Cases**: Test operations where immutability will be hardest to optimize
- **CPU Efficiency**: Measure diff and patch operations in nanoseconds/CPU cycles
- **Memory Patterns**: Understand current allocation behavior
- **Regression Prevention**: Enable performance tracking during immutability implementation

## Essential Context for LLMs

### Key Files to Read First

1. **Core API** (`/home/joseph/jd/v2/node.go:36-41`)
   - `JsonNode.Diff(n JsonNode, options ...Option) Diff`
   - `JsonNode.Patch(d Diff) (JsonNode, error)`

2. **Current Operations** (search for existing Diff/Patch methods):
   ```bash
   grep -r "func.*Diff\|func.*Patch" /home/joseph/jd/v2/*.go
   ```

3. **Existing Makefile** (`/home/joseph/jd/Makefile`) - Add benchmark targets here

4. **No Current Benchmarks**: This is greenfield implementation

## Critical Focus: Structural Sharing Challenges

### 1. **Worst-Case Scenarios** (High Priority)

These operations will be most impacted by immutability and require baseline measurement:

#### **Array Operations - Beginning Insertions**
```go
// Inserting at array start forces copying entire array with immutability
func BenchmarkPatch_Array_InsertAtBeginning(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
            // Create array [0, 1, 2, ..., size-1]
            original := generateSequentialArray(size)
            // Create patch to insert element at position 0
            diff := createInsertAtBeginningDiff()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                result, _ := original.Patch(diff)
                _ = result
            }
        })
    }
}
```

#### **Object Key Restructuring**
```go
// Adding/removing keys may require full object reconstruction
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
            diff := createKeyOperationDiff(scenario.operation)
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                result, _ := original.Patch(diff)
                _ = result
            }
        })
    }
}
```

#### **Deep Nested Modifications**
```go
// Deep changes require copying the entire path from root to leaf
func BenchmarkPatch_DeepNesting_PathCopying(b *testing.B) {
    depths := []int{5, 10, 20, 50}
    
    for _, depth := range depths {
        b.Run(fmt.Sprintf("Depth_%d", depth), func(b *testing.B) {
            // Create object nested 'depth' levels deep
            original := generateDeeplyNestedObject(depth)
            // Change leaf value at maximum depth
            diff := createDeepLeafModification(depth)
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                result, _ := original.Patch(diff)
                _ = result
            }
        })
    }
}
```

### 2. **Best-Case Scenarios** (For Comparison)

#### **Single Value Changes**
```go
// Changing one value should have maximum structural sharing
func BenchmarkPatch_Object_SingleValueChange(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("Keys_%d", size), func(b *testing.B) {
            original := generateObject(size)
            // Only change one value, not structure
            diff := createSingleValueChangeDiff()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                result, _ := original.Patch(diff)
                _ = result
            }
        })
    }
}
```

#### **Append-Only Operations**
```go
// Array appends should be efficient with immutability
func BenchmarkPatch_Array_AppendOnly(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
            original := generateArray(size)
            diff := createAppendDiff() // Add element at end
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                result, _ := original.Patch(diff)
                _ = result
            }
        })
    }
}
```

## Benchmark Implementation Structure

### File Organization

```
/home/joseph/jd/v2/
├── performance_baseline_test.go    # Main baseline benchmarks
├── array_performance_test.go       # Array operation focus
├── object_performance_test.go      # Object operation focus  
└── nested_performance_test.go      # Deep nesting scenarios
```

### Standard Patterns

#### Memory Allocation Tracking
```go
func BenchmarkWithAllocs_Operation(b *testing.B) {
    // Setup
    data := generateTestData()
    
    b.ReportAllocs() // Track memory allocations
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        result := performOperation(data)
        _ = result // Prevent optimization
    }
}
```

#### Test Data Generation Helpers
```go
// Helper functions for consistent test data
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
    
    // Build from leaf up to root
    for i := 0; i < depth; i++ {
        current = map[string]interface{}{
            "nested": current,
        }
    }
    
    node, _ := NewJsonNode(current)
    return node
}
```

## Critical Measurements

### 1. **Operation Time Scaling**
Focus on how performance degrades with size:
- **Linear scaling**: Expected for most operations
- **Quadratic scaling**: Warning signs for algorithms
- **Constant time**: Ideal for simple operations

### 2. **Memory Allocation Patterns**
Current mutable patterns to establish baseline:
- Allocations per operation
- Memory scaling with data size
- Peak memory usage during operations

### 3. **Specific Problem Operations**

#### **Array Beginning Operations**
```go
// These will be most expensive with immutability
- Insert at index 0
- Delete at index 0  
- Multiple beginning insertions
- Prepend vs append performance comparison
```

#### **Object Restructuring**
```go
// Key operations that may require full object copying
- Adding first key to empty object
- Removing keys from different positions
- Bulk key additions/removals
- Object merging operations
```

#### **Deep Path Modifications**
```go
// Operations requiring long copy paths
- Modify leaf in deeply nested structure
- Add/remove nested objects at depth
- Array operations within nested objects
```

## Makefile Integration

Add to `/home/joseph/jd/Makefile`:

```makefile
.PHONY: benchmark-baseline benchmark-memory benchmark-save

# Run baseline performance tests before immutability
benchmark-baseline:
	@echo "Running pre-immutability baseline benchmarks..."
	cd v2 && go test -bench=BenchmarkPatch -benchmem -count=3 -timeout=10m

# Focus on memory allocation patterns  
benchmark-memory:
	cd v2 && go test -bench=BenchmarkWithAllocs -benchmem -count=5 -timeout=10m

# Save results for immutability comparison
benchmark-save:
	@mkdir -p benchmarks
	@timestamp=$$(date +%Y%m%d_%H%M%S); \
	cd v2 && go test -bench=. -benchmem -count=3 > ../benchmarks/baseline_$$timestamp.txt; \
	echo "Baseline saved to benchmarks/baseline_$$timestamp.txt"

# Quick verification that benchmarks run
benchmark-quick:
	cd v2 && go test -bench=BenchmarkPatch_Object_SingleValue -benchtime=100ms
```

## Success Criteria

- **Comprehensive Baseline**: All challenging operations benchmarked
- **Clear Problem Cases**: Identify operations that will be most expensive with immutability
- **Consistent Results**: Reproducible benchmark runs
- **Memory Awareness**: Understand current allocation patterns
- **Ready for Comparison**: Established baseline for post-immutability comparison

## Implementation Priority

### Phase 1: Core Challenging Operations (Essential)
1. Array beginning insertions/deletions
2. Object key additions/removals  
3. Deep nested modifications
4. Memory allocation tracking

### Phase 2: Comparison Scenarios  
1. Best-case operations (single value changes, appends)
2. No-op patches (should be very fast)
3. Large vs small data structure scaling

### Phase 3: Analysis Setup
1. Makefile targets for easy execution
2. Result saving for comparison
3. Clear documentation of findings

This baseline will provide the foundation needed to measure immutability implementation impact and ensure performance doesn't regress during the transition.

## Baseline Performance Results

**Date:** August 13, 2025  
**Platform:** Linux (Intel Core m3-8100Y @ 1.10GHz)  
**Benchmark File:** `benchmarks/baseline_20250813_215048.txt`

### Critical Performance Issues Identified

#### 1. Exponential Diff Operations
**Worst performance bottlenecks** that will be severely impacted by immutability:

```
BenchmarkDiff_Array_LargeChanges:
- Size 100:  1.76 seconds,  1.1GB allocated,  21K allocs
- Size 1000: 224 seconds, 104GB allocated,   2M allocs  
- Size 3000: 1897 seconds, 939GB allocated, 18M allocs
```

```
BenchmarkDiff_Object_LargeObjects:
- Size 100:  358ms,   110MB allocated,  2.7K allocs
- Size 1000: 3.7s,    1GB allocated,   28K allocs
- Size 3000: 14.5s,   3.7GB allocated, 84K allocs
```

**Analysis:** These show **quadratic/exponential scaling** - the diff algorithms become prohibitively expensive with large changes. This is the highest priority area for optimization when implementing immutability.

#### 2. Array Beginning Operations - Linear Scaling
**Expected challenge areas** showing reasonable linear performance:

```
BenchmarkPatch_Array_InsertAtBeginning:
- Size 100:  1.4μs,  2KB allocated,   8 allocs
- Size 1000: 7.0μs, 16KB allocated,   8 allocs  
- Size 3000: 24.4μs, 49KB allocated,  8 allocs
```

**Analysis:** Linear scaling (17x slower for 30x data) with predictable memory growth. While this will require full array copying with immutability, the baseline performance is reasonable.

### Surprisingly Good Performance

#### 1. Object Operations - Size Independent
**Operations that won't be hurt by immutability:**

```
BenchmarkPatch_Object_KeyRestructuring:
- AddKey (Small):   1.5μs, 248B allocated, 11 allocs
- AddKey (Large):   1.8μs, 248B allocated, 11 allocs
- RemoveKey (3000): 1.5μs, 232B allocated, 10 allocs
```

**Analysis:** Performance is essentially constant regardless of object size (10 vs 3000 keys). This suggests excellent hash table implementation that won't degrade with structural sharing.

#### 2. Deep Nesting - Reasonable Scaling
**Path copying operations scale well:**

```
BenchmarkPatch_DeepNesting_PathCopying:
- Depth 5:  3.0μs,  664B allocated, 22 allocs
- Depth 10: 5.0μs, 1032B allocated, 28 allocs
- Depth 15: 9.4μs, 1144B allocated, 33 allocs
```

**Analysis:** Only 3x slower for 3x depth increase. Deep path copying for immutability will be manageable.

#### 3. No-Op Operations - Optimal Baseline
```
BenchmarkPatch_NoOp: 4.6ns, 0B allocated, 0 allocs
```

**Analysis:** Excellent baseline showing minimal overhead for unchanged structures.

### Memory Allocation Patterns

1. **Array Operations:** Linear, predictable memory growth
2. **Diff Operations:** Massive allocations (1GB+ for large changes) 
3. **Object Operations:** Constant small allocations (~300B)
4. **Deep Nesting:** Linear growth with depth

### Recommendations for Immutability Implementation

#### High Priority Optimization Areas:
1. **Diff generation algorithms** - Focus on structural sharing for large change detection
2. **Array middle insertions** - Consider specialized data structures
3. **Object diffing optimization** - Leverage size-independent key operations

#### Operations Ready for Immutability:
1. **Single value changes** - Already fast, will benefit from structural sharing
2. **Key add/remove operations** - Size-independent performance proven
3. **Deep nesting modifications** - Linear scaling indicates good foundation

#### Success Metrics:
- **Critical:** Keep diff operations under exponential complexity
- **Important:** Maintain linear scaling for array operations  
- **Baseline:** Preserve sub-microsecond performance for simple operations

This baseline establishes clear performance expectations and identifies the operations most at risk during the immutability transition.