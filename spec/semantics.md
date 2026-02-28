# Structural Format Semantics

This document defines the semantic behavior of structural diff operations, data structures, and algorithms.

## Core Concepts

### 1. Document Model

The structural format operates on a **hierarchical document model** where:
- **Documents** are trees of JSON values
- **Paths** identify specific locations within documents  
- **Values** can be primitives (string, number, boolean, null), arrays, or objects
- **Void** represents the absence of a value (different from null)

### 2. Path Resolution

Paths are sequences of path elements that navigate the document tree:

```
["users", 0, "name"] → document.users[0].name
["config", "timeout"] → document.config.timeout  
[] → document (root)
```

#### Path Element Types:

1. **String**: Object property key
2. **Number**: Array index (0-based, -1 for append)
3. **{}**: Set operation marker
4. **[]**: List/multiset operation marker
5. **{"key":"value"}**: Object matching by specific keys

## Diff Generation

### Basic Algorithm

1. **Traverse** both documents recursively
2. **Compare** values at each path
3. **Generate** minimal change set using context
4. **Apply options** to modify comparison behavior

### Value Comparison Rules

#### Primitive Values
- **Exact equality** for strings, booleans, null
- **Precision-aware equality** for numbers when precision option is set
- **Type mismatch** always produces change

#### Objects
- **Key-by-key comparison** of properties
- **Missing keys** produce additions/removals
- **Value changes** produce nested diffs

#### Arrays  
- **LCS algorithm** for minimal diff (by default)
- **Set semantics** when SET option applied (order ignored)  
- **Multiset semantics** when MULTISET option applied (order ignored, duplicates counted)

### Context Preservation

For array modifications, the structural format shows minimal surrounding context:

```
Original: ["a", "b", "c", "d"]
Modified: ["a", "x", "y", "d"] 

Diff:
@ [1]
  "a"        <- context before (minimal)
- "b"        <- removal  
- "c"        <- removal
+ "x"        <- addition
+ "y"        <- addition
  "d"        <- context after (minimal)
```

**Context Rules:**
1. **Minimal Context Strategy**: Show exactly one element before/after when available
2. **Contextual Boundary Markers**: 
   - `[` appears only when showing changes at/near array beginning
   - `]` appears only when showing changes at/near array end
   - Middle changes don't need brackets (array indices provide context)
3. **Formatting**: Context elements use two-space indentation, changes use `+`/`-` with single space
4. **Consistency**: Same minimal context approach for all array sizes

## Options Processing

### Global Options

#### SET
```
^ "SET"
```
- Treats all arrays as mathematical sets
- **Order ignored**: `[1,2,3]` equals `[3,1,2]`
- **Duplicates ignored**: `[1,1,2]` equals `[1,2]`
- Uses hash-based comparison for efficiency

#### MULTISET  
```
^ "MULTISET"
```  
- Treats arrays as multisets (bags)
- **Order ignored**: `[1,2,3]` equals `[3,1,2]`
- **Duplicates counted**: `[1,1,2]` differs from `[1,2]`
- Tracks element frequency

#### MERGE
```
^ "MERGE"
```
- Enables merge-patch semantics (RFC 7386)
- **Null removes** object properties
- **Objects merge** recursively
- **Arrays replace** entirely

#### Precision
```
^ {"precision": 0.001}
```
- Sets numeric comparison tolerance
- Numbers within tolerance are considered equal
- **Absolute difference**: `|a - b| <= precision`
- Incompatible with SET/MULTISET (uses hashing)

#### Keys
```
^ {"keys": ["id", "name"]}
```
- Defines object matching keys for arrays
- Objects with same key values are considered identical
- Enables object-level diffing within arrays

### PathOptions

PathOptions apply options to specific document paths:

```
^ {"@": ["users"], "^": ["SET"]}
```

#### Syntax
- `"@"`: Array of path elements (JSON path)
- `"^"`: Array of options to apply at that path

#### Inheritance Rules
1. **Child paths inherit** parent PathOptions
2. **More specific paths override** general ones
3. **Multiple options** on same path are combined
4. **Global options** apply everywhere unless overridden

#### Path Matching
```
Path: ["users", 0, "tags"]
PathOption: {"@": ["users"], "^": ["SET"]}
Result: SET applies to users[0].tags (inherited)

PathOption: {"@": ["users", 0], "^": ["MULTISET"]}  
Result: MULTISET overrides SET for users[0] and children
```

### DIFF_ON/DIFF_OFF Options

Control which parts of documents are compared:

#### DIFF_OFF
```
^ {"@": ["metadata"], "^": ["DIFF_OFF"]}
```
- **Ignores changes** at specified path
- Useful for timestamps, auto-generated fields
- Children also ignored unless overridden

#### DIFF_ON
```  
^ {"@": [], "^": ["DIFF_OFF"]}        # Ignore everything
^ {"@": ["data"], "^": ["DIFF_ON"]}   # Except data
```
- **Enables diffing** at specified path  
- Overrides parent DIFF_OFF settings
- Allows allow-list approach

## Array Diffing Algorithms

### 1. List Diffing (Default)

Uses **Longest Common Subsequence (LCS)** algorithm:

1. Find longest sequence of unchanged elements
2. Generate minimal insertions/deletions
3. Preserve context around changes
4. Maintain array order semantics

**Example:**
```
A: [1, 2, 3, 4]
B: [1, 5, 6, 4]

LCS: [1, 4] (common elements)
Operations: remove 2,3 at index 1, add 5,6 at index 1
```

### 2. Set Diffing

When SET option is applied:

1. **Convert to sets**: Remove duplicates, ignore order
2. **Find additions**: Elements in B but not A  
3. **Find removals**: Elements in A but not B
4. **Use hash comparison** for efficiency

**Example:**
```
^ "SET"
A: [3, 1, 2, 1] → Set{1, 2, 3}
B: [2, 4, 1] → Set{1, 2, 4}

Removals: {3}
Additions: {4}
```

### 3. Multiset Diffing

When MULTISET option is applied:

1. **Count frequencies**: Track element occurrences
2. **Compare counts**: Find frequency differences
3. **Generate changes**: Add/remove based on count differences

**Example:**
```
^ "MULTISET"  
A: [1, 1, 2, 3] → {1:2, 2:1, 3:1}
B: [1, 2, 2, 4] → {1:1, 2:2, 4:1}

Changes:
- 1     (reduce frequency 2→1)  
- 3     (remove entirely)
+ 2     (increase frequency 1→2)
+ 4     (add new)
```

## Object Matching with Keys

For arrays containing objects, Keys enables object-level comparison:

```
^ {"keys": ["id"]}

A: [{"id": "user1", "name": "Alice", "age": 25}]
B: [{"id": "user1", "name": "Alice", "age": 26}]

Path: [{"id": "user1"}, "age"]  # Match by id, diff age property
Result:
@ [{"id":"user1"},"age"]
- 25
+ 26
```

### Matching Algorithm

1. **Extract matching keys** from each object
2. **Group by key values** in both arrays
3. **Compare matched objects** recursively
4. **Handle unmatched objects** as additions/removals

### Multiple Keys
```
^ {"keys": ["type", "id"]}
```
Objects match when ALL specified keys have equal values.

### Duplicate Identity Keys

When multiple objects in the same array share the same identity key values under set semantics, the behavior is undefined. Implementations MAY reject this as an error, silently pick one object, or handle it in any other way. Users should ensure identity keys are unique within each array.

## Patch Application

### Application Algorithm

1. **Parse diff** into structured operations
2. **Validate paths** exist in target document
3. **Apply changes** in path order (depth-first)
4. **Verify context** when specified
5. **Handle conflicts** according to mode

### Context Validation

For array operations with context:
- **Before context** must match elements preceding the change
- **After context** must match elements following the change  
- **Mismatched context** produces application error

### Merge Semantics

When MERGE option is present:
- **Null values remove** object properties
- **Objects merge recursively** rather than replacing
- **Arrays replace entirely** (no element-wise merging)
- **Void values** (empty +) set properties to null

## Error Conditions

### Path Resolution Errors
- **Invalid path**: Path element doesn't exist
- **Type mismatch**: Path expects object but finds array
- **Out of bounds**: Array index exceeds bounds

### Value Errors
- **Invalid JSON**: Malformed JSON values
- **Type conflicts**: Cannot convert between incompatible types

### Option Conflicts
Options fall into two groups that are mutually exclusive:
- **Equivalence modifiers**: Options that alter value equality (e.g., `precision`, future case insensitivity)
- **Set semantics**: Options that use hash-based comparison (`SET`, `MULTISET`)

Implementations MUST reject combinations across these groups because set operations require hash-stable equivalence. The `keys` option belongs to neither group and is compatible with both.

Conflict detection within nested PathOptions is optional; implementations that return `(Diff, error)` can detect conflicts lazily during traversal.

### Context Errors
- **Context mismatch**: Expected context doesn't match actual values
- **Missing context**: Required context elements not found

## Implementation Requirements

### Performance Characteristics

- **LCS Algorithm**: O(m×n) time complexity for arrays of size m and n
- **Hash-based sets**: O(n) expected time for set operations
- **Deep recursion**: May require stack management for deeply nested structures

### Memory Considerations

- **Context preservation**: Requires storing surrounding elements
- **Path tracking**: Must maintain full path context during traversal
- **Option inheritance**: Requires efficient option lookup by path

### Precision and Accuracy

- **Floating-point comparison**: Must handle IEEE 754 edge cases
- **Unicode normalization**: Should handle equivalent Unicode representations
- **JSON canonicalization**: Numbers should be normalized (e.g., 1.0 → 1)

This semantic specification defines the full behavior of structural diff operations. Implementations following these semantics will produce consistent, interoperable results.