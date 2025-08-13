# Immutable JsonNode Implementation Plan

## Overview

This document provides a detailed, systematic plan for implementing immutable JsonNode behavior in the v2 system. The goal is to ensure that all patch operations return new JsonNode instances while preserving structural sharing for performance, maintaining full backward compatibility.

## Critical Constraints

- **ONLY modify files in `/home/joseph/jd/v2/` directory**
- **NO public method signature changes allowed**
- **ALL 89 test functions across 22 test files must pass**
- **Maintain full backward compatibility**

## Context Loading for LLMs

### Essential Files to Read First

1. **Core Interfaces** (`/home/joseph/jd/v2/node.go:37-51`)
   - `JsonNode` interface with `Patch(d Diff) (JsonNode, error)` method  
   - `jsonNodeInternals` interface with `patch(...)` method
   - These signatures CANNOT change

2. **Current Mutation Points** (files with in-place modifications):
   - `/home/joseph/jd/v2/object.go:147-202` - Object patch method
   - `/home/joseph/jd/v2/list.go:324-400` - List patch method  
   - `/home/joseph/jd/v2/array.go` - Array patch method
   - `/home/joseph/jd/v2/set.go` - Set patch method
   - `/home/joseph/jd/v2/multiset.go` - Multiset patch method

3. **All Patch Method Locations** (search results from `func.*Patch.*Diff`):
   ```
   /home/joseph/jd/v2/bool.go:57
   /home/joseph/jd/v2/set.go:102  
   /home/joseph/jd/v2/array.go:99
   /home/joseph/jd/v2/null.go:53
   /home/joseph/jd/v2/object.go:143
   /home/joseph/jd/v2/string.go:53
   /home/joseph/jd/v2/list.go:320
   /home/joseph/jd/v2/number.go:68
   /home/joseph/jd/v2/void.go:74
   /home/joseph/jd/v2/multiset.go:90
   ```

4. **Test Files** (22 files with 89 test functions) - all must continue passing

### Key Patterns to Understand

1. **Current Mutation Pattern** (object.go:197-200):
   ```go
   delete(o, string(pe))         // MUTATES original
   o[string(pe)] = patchedNode   // MUTATES original
   return o, nil                 // Returns mutated original
   ```

2. **Current Mutation Pattern** (list.go:363):
   ```go
   l[i] = patchedNode   // MUTATES original slice
   return l, nil        // Returns mutated original
   ```

3. **Patch Call Chain**:
   - Public `Patch(d Diff)` calls `patchAll()` in `patch_common.go:7-20`
   - `patchAll()` calls internal `patch()` method for each DiffElement
   - Internal `patch()` methods are where mutations currently occur

## Implementation Phases

### Phase 1: Primitive Types (Already Immutable)
**Files**: `bool.go`, `string.go`, `number.go`, `null.go`, `void.go`

**Status**: ✅ **NO CHANGES NEEDED**
- These types are already immutable (strings, numbers, etc.)
- Their patch methods already return new instances
- Verify with tests that they don't mutate

**Validation**: Run type-specific tests
```bash
go test -v ./v2/ -run "TestBool|TestString|TestNumber|TestNull|TestVoid"
```

### Phase 2: Object Immutability
**File**: `/home/joseph/jd/v2/object.go`

**Current Problem** (lines 195-201):
```go
if isVoid(patchedNode) {
    delete(o, string(pe))     // MUTATION!
} else {
    o[string(pe)] = patchedNode  // MUTATION!
}
return o, nil  // Returns mutated original
```

**Solution**: Implement copy-with-modification
```go
// Create new object with structural sharing
newObj := make(jsonObject, len(o))
for k, v := range o {
    if k != string(pe) {  // Skip key being modified/deleted
        newObj[k] = v     // Structural sharing - same reference
    }
}
if !isVoid(patchedNode) {
    newObj[string(pe)] = patchedNode  // Add/replace in new object
}
return newObj, nil
```

**Testing Strategy**:
1. Create test to verify original object unchanged after patch
2. Verify all existing object tests continue passing
3. Add memory/performance test to ensure structural sharing works

### Phase 3: List Immutability  
**File**: `/home/joseph/jd/v2/list.go`

**Current Problem** (line 363):
```go
l[i] = patchedNode  // MUTATION!
return l, nil       // Returns mutated original
```

**Solution**: Implement copy-with-slice-sharing
```go
// Create new slice with same capacity
newList := make(jsonList, len(l))
copy(newList, l)  // Copy slice structure (pointers only)
newList[i] = patchedNode  // Modify in new slice
return newList, nil
```

**Advanced Optimization** (for large lists):
- For small changes (<10% of list), use copy approach above
- For large changes, consider more sophisticated persistent data structures

**Testing Strategy**:
1. Create test to verify original list unchanged after patch  
2. Verify all list tests continue passing
3. Test with large lists (>1000 elements) for performance

### Phase 4: Array Immutability
**File**: `/home/joseph/jd/v2/array.go`

**Current Analysis Needed**:
- Read current `array.go` patch implementation
- Likely similar mutation pattern to list.go
- Apply same copy-with-modification strategy

### Phase 5: Set and Multiset Immutability
**Files**: `/home/joseph/jd/v2/set.go`, `/home/joseph/jd/v2/multiset.go`

**Current Analysis Needed**:
- Read current set/multiset patch implementations  
- These likely use maps or slices internally
- Apply appropriate copy-with-modification strategy based on internal structure

### Phase 6: Integration Testing
**Critical Validation**:

1. **All Unit Tests Pass**:
   ```bash
   go test ./v2/
   ```

2. **Immutability Verification Tests** (new tests to add):
   ```go
   func TestPatchImmutability(t *testing.T) {
       original := /* create test object */
       diff := /* create test diff */
       
       // Keep reference to original
       originalCopy := /* deep comparison data */
       
       patched, err := original.Patch(diff)
       
       // Verify original unchanged
       assert.Equal(t, originalCopy, original)
       // Verify patched is different
       assert.NotEqual(t, original, patched)
       // Verify structural sharing where possible
   }
   ```

3. **Memory/Performance Tests**:
   ```go
   func TestStructuralSharing(t *testing.T) {
       // Verify that unchanged portions share memory
       // Measure memory usage before/after patches
       // Ensure no excessive copying
   }
   ```

4. **End-to-End Tests**:
   ```bash
   go test ./v2/ -run "E2E"
   ```

## Systematic Implementation Order

### Step 1: Analysis Phase
1. Read all current patch implementations
2. Identify exact mutation points in each type
3. Understand test coverage for each type

### Step 2: Implementation Phase  
1. Start with `object.go` (most complex)
2. Implement `list.go` (second most complex)
3. Implement remaining collection types (`array.go`, `set.go`, `multiset.go`)
4. Add immutability verification tests

### Step 3: Validation Phase
1. Run all existing tests after each change
2. Add new immutability tests
3. Performance benchmarking
4. Memory usage analysis

## Performance Considerations

### Structural Sharing Benefits
- **Memory**: Only modified path from root to changed node is copied
- **Time**: O(log n) complexity with persistent data structures
- **Correctness**: Eliminates accidental mutations

### Example Structural Sharing:
```
Original: {a: 1, b: {c: 2, d: 3}, e: 4}
Patch: set a = 10

Result:   {a: 10, b: {c: 2, d: 3}, e: 4}
          ↑new   ↑shared subtree  ↑shared

Only root object is new, all subtrees shared.
```

### Advanced Optimizations (Future)
- **HAMT** (Hash Array Mapped Trie) for large objects
- **Persistent Vectors** for large lists
- **Reference counting** for memory management

## Debugging and Validation

### Key Invariants to Test
1. **Immutability**: `original.Equals(originalCopy)` after any patch
2. **Correctness**: `patched.Equals(expected)` for all patch operations  
3. **Structural Sharing**: Memory usage doesn't grow linearly with patches
4. **Performance**: No significant slowdown vs current mutable version

### Common Pitfalls to Avoid
1. **Deep copying everything** - defeats performance benefits
2. **Forgetting to copy slice headers** - can still mutate through shared slices
3. **Not handling edge cases** - empty objects/lists, void nodes, etc.
4. **Breaking public interfaces** - all existing method signatures must remain

## Success Criteria

- ✅ All 89 existing tests pass
- ✅ New immutability tests pass  
- ✅ No public API changes
- ✅ Memory usage optimized through structural sharing
- ✅ Thread-safe JsonNode operations (bonus)

This plan provides the systematic approach needed for an LLM to implement immutable JsonNodes efficiently while maintaining full backward compatibility.
