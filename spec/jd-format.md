# Structural JSON Diff Format Specification

  
**Status:** Draft Specification  
**Date:** 2025

## Abstract

The structural format is a human-readable diff format for JSON and YAML documents. It provides context-aware diffing with support for complex data structures, configurable comparison semantics, and translation between standard patch formats. This specification defines the complete syntax, semantics, and behavior required for implementing structural format tools.

## Table of Contents

1. [Introduction](#introduction)
2. [Format Overview](#format-overview)  
3. [Syntax Specification](#syntax-specification)
4. [Semantic Behavior](#semantic-behavior)
5. [Options System](#options-system)
6. [Error Handling](#error-handling)
7. [Implementation Requirements](#implementation-requirements)
8. [Interoperability](#interoperability)
9. [Security Considerations](#security-considerations)
10. [Examples](#examples)

## 1. Introduction

### 1.1 Purpose

The structural format addresses limitations in existing JSON diff formats:

- **Human Readability**: Unlike RFC 6902 (JSON Patch), the structural format provides unified diff-style output with familiar syntax
- **Context Preservation**: Shows surrounding elements in arrays to clarify change locations
- **Flexible Semantics**: Supports set/multiset semantics, precision-based numeric comparison, and path-specific comparison options
- **Format Interoperability**: Translates to and from standard patch formats (RFC 6902, RFC 7386)

### 1.2 Scope

This specification defines:
- Full syntax grammar for parsing and generation
- Semantic behavior for diff generation and patch application  
- Options system for controlling comparison behavior
- Error conditions and handling procedures
- Implementation requirements and testing guidelines

### 1.3 Terminology

- **Document**: A JSON or YAML data structure  
- **Path**: A sequence of elements identifying a location within a document
- **Diff**: A representation of changes between two documents
- **Patch**: Application of a diff to transform one document into another
- **Context**: Unchanged elements shown to provide understanding of changes
- **Void**: Absence of a value (distinct from JSON null)

### 1.4 Notational Conventions

This specification uses ABNF grammar notation (RFC 5234) and includes normative examples marked with "MUST" and informative examples marked with "Example:".

## 2. Format Overview

### 2.1 Document Structure

A structural diff consists of:
1. **Optional metadata header**: Lines starting with `^` that specify options
2. **Diff elements**: Blocks specifying location and changes
3. **Path lines**: Starting with `@` followed by JSON array path
4. **Change lines**: Starting with `-` (removal) or `+` (addition)  
5. **Context lines**: Starting with two spaces, showing unchanged values

### 2.2 Basic Example

```diff
@ ["user","name"]
- "Alice"
+ "Bob"
```

This shows changing the `user.name` property from "Alice" to "Bob".

### 2.3 Array Context Example

```diff
@ ["items",1]
  "apple"
- "banana"  
+ "blueberry"
  "cherry"
]
```

This shows changing `items[1]` from "banana" to "blueberry", with minimal context showing the surrounding array elements. Note the contextual `]` bracket indicating the change extends to the array end.

## 3. Syntax Specification

The complete ABNF grammar is provided in [grammar.md](grammar.md). Key elements:

### 3.1 Document Structure

```abnf
JdDiff = *MetadataLine *DiffElement
DiffElement = PathLine [ArrayOpen] *ContextLine *ChangeLine [*ContextLine] [ArrayClose]
```

### 3.2 Path Specification

```abnf
PathLine = "@" SP JsonArray CRLF
```

Paths use JSON array syntax to specify document locations:
- `["key"]` - Object property access
- `[0]` - Array index access  
- `[{}]` - Set operation on array
- `[{"id":"value"}]` - Object matching by key

### 3.3 Content Lines

```abnf
ChangeLine = (AddLine / RemoveLine)
AddLine = "+" SP [JsonValue] CRLF
RemoveLine = "-" SP JsonValue CRLF
ContextLine = SP SP JsonValue CRLF
```

### 3.4 Metadata Lines

```abnf
MetadataLine = "^" SP JsonValue CRLF
```

Examples:
- `^ "SET"` - Use set semantics
- `^ {"precision":0.01}` - Numeric precision tolerance
- `^ {"@":["path"],"^":["SET"]}` - PathOption applying SET to specific path

## 4. Semantic Behavior

Complete semantic definitions are provided in [semantics.md](semantics.md). Key behaviors:

### 4.1 Path Resolution

Paths navigate document structure using these rules:
1. **String elements** access object properties
2. **Numeric elements** access array indices (0-based, -1 for append)
3. **Special markers** modify comparison behavior:
   - `{}` enables set semantics
   - `[]` indicates list/multiset operations  
   - `{"key":"value"}` matches objects by key values

### 4.2 Array Diffing

**Default behavior** uses Longest Common Subsequence (LCS) algorithm:
1. Find longest sequence of unchanged elements
2. Generate minimal insertions/deletions
3. Show context around changes
4. Maintain array order semantics

**Set behavior** (when `{}` marker or SET option used):
1. Ignore element order
2. Remove duplicates  
3. Show additions/removals only
4. Use hash-based comparison

**Multiset behavior** (when MULTISET option used):
1. Ignore element order
2. Count duplicate frequencies
3. Show frequency changes
4. Track element occurrences

### 4.3 Context Preservation

For array modifications:
1. **Minimal Context**: Show exactly one unchanged element before/after when available
2. **Contextual Boundary Markers**: 
   - `[` only when showing changes at/near array beginning
   - `]` only when showing changes at/near array end
3. **Formatting**: Context lines indented with two spaces, changes with one space after `+`/`-`
4. **Consistency**: Same minimal approach for all array sizes

### 4.4 Value Comparison

**Primitives**: Exact equality except for numbers when precision option set
**Objects**: Recursive key-by-key comparison  
**Arrays**: Algorithm determined by options (LCS, set, or multiset)
**Numbers**: Absolute difference comparison when precision specified

## 5. Options System

### 5.1 Global Options

Applied to entire diff operation:

- **SET**: `^ "SET"` - All arrays treated as mathematical sets
- **MULTISET**: `^ "MULTISET"` - All arrays treated as multisets (bags)  
- **MERGE**: `^ "MERGE"` - Use JSON Merge Patch semantics (RFC 7386)
- **COLOR**: `^ "COLOR"` - Add ANSI color codes to output
- **Precision**: `^ {"precision":N}` - Numbers within N absolute difference considered equal
- **SetKeys**: `^ {"setkeys":["key1","key2"]}` - Object matching keys for arrays

### 5.2 PathOptions

Apply options to specific document paths:

```diff
^ {"@":["users"],"^":["SET"]}
```

This applies SET semantics only to the "users" path.

**Inheritance Rules:**
1. Child paths inherit parent PathOptions
2. More specific paths override general ones
3. Multiple options on same path are combined
4. Global options apply unless overridden

### 5.3 DIFF_ON/DIFF_OFF

Control which document parts are compared:

```diff
^ {"@":["metadata"],"^":["DIFF_OFF"]}
^ {"@":["data"],"^":["DIFF_ON"]}
```

**DIFF_OFF** ignores changes at specified paths (useful for timestamps, auto-generated fields)
**DIFF_ON** enables diffing at specified paths (overrides parent DIFF_OFF)

## 6. Error Handling

Complete error specifications are provided in [errors.md](errors.md). Key categories:

### 6.1 Parse Errors
- Invalid diff syntax
- Malformed JSON values
- Invalid path arrays
- Unrecognized line prefixes

### 6.2 Path Resolution Errors
- Non-existent paths
- Type mismatches (accessing array index on object)
- Index out of bounds
- Invalid path element types

### 6.3 Patch Application Errors  
- Context mismatches
- Precondition failures
- Conflicting operations
- Type conversion errors

### 6.4 Option Errors
- Conflicting options (precision with sets)
- Invalid option values
- Malformed PathOptions
- Unknown options

### 6.5 Resource Errors
- Memory limits exceeded
- Recursion depth exceeded
- Document size limits
- Operation timeouts

## 7. Implementation Requirements

### 7.1 Compliance Levels

**Core Compliance:**
- Basic diff generation and patch application
- Simple path navigation (object keys, array indices)
- Context preservation in array diffs
- Standard error handling

**Extended Compliance:**
- Full options support (SET, MULTISET, precision, setkeys)
- PathOptions with inheritance
- DIFF_ON/DIFF_OFF functionality  
- Options header rendering

**Format Compliance:**
- Translation between jd, RFC 6902, and RFC 7386 formats
- Preservation of semantic equivalence
- Round-trip conversion accuracy

### 7.2 Performance Requirements

- **LCS Algorithm**: O(m×n) time complexity for arrays of size m and n
- **Memory Usage**: Implementations SHOULD handle documents up to 100MB
- **Nesting Depth**: SHOULD support at least 1000 levels of nesting
- **Path Length**: SHOULD support paths with at least 1000 elements

### 7.3 Character Encoding

- **Input/Output**: UTF-8 encoding required
- **Line Endings**: LF (`\n`) only, not CRLF
- **Unicode**: Full Unicode support with proper JSON string escaping
- **Normalization**: Implementations MAY normalize equivalent Unicode representations

### 7.4 Security Considerations

- **Input Validation**: Thoroughly validate all JSON and diff inputs
- **Resource Limits**: Implement reasonable limits to prevent DoS attacks
- **Memory Management**: Handle large documents without unbounded memory growth
- **Recursion Control**: Prevent stack overflow from deeply nested structures

## 8. Interoperability

### 8.1 JSON Patch Translation (RFC 6902)

The structural format can translate to/from JSON Patch:

**Structural to JSON Patch:**
- Remove operations include test operations for validation
- Array indices use JSON Pointer format (`/path/0`)
- Context information is lost in translation

**JSON Patch to Structural:**
- Test operations become context validation
- Path pointers convert to jd path arrays
- Operations grouped by target path when possible

### 8.2 JSON Merge Patch Translation (RFC 7386)

**Structural to JSON Merge Patch:**
- Only works with MERGE option diffs
- Null values indicate deletions
- Nested objects merge recursively
- Arrays replace entirely

**JSON Merge Patch to Structural:**
- Automatically adds MERGE option header
- Null values become void additions (`+`)
- Object merging shown as targeted additions

### 8.3 Format Detection

Implementations SHOULD auto-detect input formats:
- Structural format: Starts with `@` or `^` lines
- JSON Patch: Array of operation objects
- JSON Merge Patch: Single JSON object (with restrictions)
- Regular JSON: Standard JSON documents

## 9. Security Considerations

### 9.1 Resource Exhaustion

- **Large Documents**: Implement size limits (suggested: 100MB maximum)
- **Deep Nesting**: Limit recursion depth (suggested: 1000 levels maximum)
- **Long Paths**: Limit path element count (suggested: 1000 elements maximum)
- **Memory Usage**: Use streaming parsers when possible

### 9.2 Input Validation

- **JSON Parsing**: Use secure JSON parsers that handle edge cases
- **Path Validation**: Ensure path elements are valid JSON values
- **Option Validation**: Validate all option parameters
- **Character Encoding**: Validate UTF-8 encoding

### 9.3 Information Disclosure

- **Error Messages**: Avoid exposing sensitive data in error messages
- **Path Leakage**: Be careful not to expose private document structure
- **Memory Dumps**: Clear sensitive data from memory after use

### 9.4 Injection Attacks

- **Path Injection**: Validate that paths don't access unexpected document areas
- **Option Injection**: Sanitize all option values
- **Command Injection**: When used in shell contexts, properly escape arguments

## 10. Examples

Comprehensive examples are provided in [examples.md](examples.md), including:

- Basic object and array changes
- Complex nested structure diffs
- All option types and combinations
- PathOptions with inheritance
- Format translation examples
- Real-world use cases
- Edge cases and error scenarios

## Appendices

### Appendix A: Grammar Reference

See [grammar.md](grammar.md) for complete ABNF grammar specification.

### Appendix B: Semantic Reference  

See [semantics.md](semantics.md) for detailed semantic definitions.

### Appendix C: Error Reference

See [errors.md](errors.md) for complete error handling specification.

### Appendix D: Test Suite

The test suite in `test/` provides comprehensive validation for implementations.

## Conformance

An implementation conforms to this specification if it:

1. **Parses** all valid jd format inputs according to the grammar
2. **Generates** syntactically correct jd format output
3. **Implements** semantic behavior as specified
4. **Handles** all documented error conditions appropriately  
5. **Passes** the compliance test suite for its declared compliance level

Implementations SHOULD declare their compliance level (core, extended, or format) and document any extensions or limitations.

## References

- RFC 5234: Augmented BNF for Syntax Specifications
- RFC 6902: JavaScript Object Notation (JSON) Patch  
- RFC 7386: JSON Merge Patch
- RFC 7159: The JavaScript Object Notation (JSON) Data Interchange Format
- RFC 3629: UTF-8, a transformation format of ISO 10646

---

*This specification is designed to enable independent, interoperable implementations of the structural diff format while maintaining compatibility with the reference Go implementation at https://github.com/josephburnett/jd.*